package t9n

import (
	"common/assert"
	"common/models/common"
	"database/sql"
	"fmt"
	"strings"
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

// Mixed returns text and its translated
type Mixed struct {
	Text       string         `db:"text" json:"text"`
	Translated sql.NullString `db:"translated" json:"translated"`
}

// DumpAll dump all the translated data from repo
func (m *Manager) DumpAll(lang string) []Mixed {
	query := fmt.Sprintf(
		"SELECT text, translated from %[1]s LEFT JOIN strings ON %[1]s.string_id = %[2]s.id "+
			"WHERE lang=?",
		TranslationsTableFull,
		StringsTableFull)
	var result []Mixed
	_, err := m.GetDbMap().Select(&result, query, lang)
	if err != nil {
		assert.Nil(err)
	}
	return result
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
	var tmp []Mixed
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

// TranslateDataTable is data table version of translation
type TranslateDataTable struct {
	ID         int64             `db:"id" json:"id"`
	Text       string            `db:"text" json:"text"`
	Lang       common.NullString `db:"lang" json:"lang"`
	Translated common.NullString `db:"translated" json:"translated"`
	Actions    string            `db:"-" json:"_actions" visible:"false"`
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

// FillTranslateDataTableArray is the function to handle
func (m *Manager) FillTranslateDataTableArray(
	lang string,
	filters map[string]string,
	search map[string]string,
	sort, order string,
	p, c int) (TranslateDataTableArray, int64) {
	var params []interface{}
	var res TranslateDataTableArray
	var where []string

	countQuery := fmt.Sprintf("SELECT COUNT(%[1]s.id) FROM %[1]s  ", StringsTableFull)

	query := fmt.Sprintf(
		`SELECT s.id, s.text, t.translated, t.lang FROM %[1]s AS s LEFT JOIN %[2]s AS t ON s.id = t.string_id AND t.lang = "%s" `,
		StringsTableFull,
		TranslationsTableFull,
		lang,
	)

	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s=?", field))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s LIKE ?", column))
		params = append(params, "%"+val+"%")
	}
	//check for perm
	if len(where) > 0 {
		query += " WHERE "
		countQuery += " WHERE "
	}
	query += strings.Join(where, " AND ")
	countQuery += strings.Join(where, " AND ")
	limit := c
	offset := (p - 1) * c
	if sort != "" {
		query += fmt.Sprintf(" ORDER BY %s %s ", sort, order)
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)
	count, err := m.GetDbMap().SelectInt(countQuery, params...)
	assert.Nil(err)
	fmt.Println(query)
	_, err = m.GetDbMap().Select(&res, query, params...)
	assert.Nil(err)
	return res, count
}

type (
	// TranslateDataTableArray is the array of translation struct
	TranslateDataTableArray []TranslateDataTable
)

// Filter is for filtering base on permission
func (tdt TranslateDataTable) Filter() TranslateDataTable {
	res := TranslateDataTable{}

	res.ID = tdt.ID

	res.Text = tdt.Text

	res.Lang = tdt.Lang

	res.Translated = tdt.Translated

	res.Actions = tdt.Actions

	action := []string{}

	res.Actions = strings.Join(action, ",")
	return res
}
