package t9n

import (
	"common/assert"
	"database/sql"
	"fmt"
	"sync"
)

// Strings model
// @Model {
//		table = strings
//		primary = true, id
//		list = yes
// }
type Strings struct {
	ID   int64  `db:"id" json:"id"`
	Text string `db:"text" json:"text"`
}

// Translations model
// @Model {
//		table = translations
//		primary = true, id
//		list = yes
// }
type Translations struct {
	ID         int64  `db:"id" json:"id"`
	StringID   int64  `db:"string_id" json:"string_id"`
	Lang       string `db:"lang" json:"lang"`
	Translated string `db:"translated" json:"translated"`
}

type mixed struct {
	Text       string         `db:"text"`
	Translated sql.NullString `db:"translated"`
}

var (
	lock    = sync.Mutex{}
	allData = make(map[string]map[string]string)
)

// LoadAllInMap try to load all translation from database into memory
func (m *Manager) LoadAllInMap(lang string) map[string]string {
	lock.Lock()
	defer lock.Unlock()

	query := fmt.Sprintf(
		"SELECT s.text, t.translated FROM %s AS s LEFT JOIN %s AS t ON t.string_id = s.id AND t.lang=?",
		StringsTableFull,
		TranslationsTableFull,
	)
	var tmp []mixed
	_, err := m.GetDbMap().Select(&tmp, query, lang)
	assert.Nil(err)

	res := make(map[string]string)
	for i := range tmp {
		if tmp[i].Translated.Valid {
			res[tmp[i].Text] = tmp[i].Translated.String
		} else {
			res[tmp[i].Text] = tmp[i].Text
		}
	}
	allData[lang] = res
	return res
}

// AddMissing Add missing translation
func (m *Manager) AddMissing(txt string) error {
	lock.Lock()
	defer lock.Unlock()

	tmp := Strings{
		Text: txt,
	}

	err := m.CreateStrings(&tmp)
	if err != nil {
		return err
	}
	for i := range allData {
		if allData[i] != nil {
			allData[i][txt] = txt
		}

	}
	return nil
}

// CreateOnDuplicateUpdateTranslations try to save a new Translations in database
func (m *Manager) CreateOnDuplicateUpdateTranslations(t *Translations) error {
	q := fmt.Sprintf(
		"INSERT INTO %s (string_id,lang,translated) VALUES (?,?,?) ON DUPLICATE KEY UPDATE translated=VALUES(translated)",
		TranslationsTableFull,
	)
	_, err := m.GetDbMap().Exec(
		q,
		t.StringID,
		t.Lang,
		t.Translated,
	)
	return err
}
