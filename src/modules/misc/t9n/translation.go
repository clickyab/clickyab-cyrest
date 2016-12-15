package t9n

import (
	"database/sql"
	"sync"
	"time"
)

// Translation model
// @Model {
//		table = translations
//		primary = true, id
//		list = yes
// }
type Translation struct {
	ID        int64          `db:"id" json:"id"`
	String    string         `db:"string" json:"string"`
	Single    sql.NullString `db:"single" json:"single"`
	Plural    sql.NullString `db:"plural" json:"plural"`
	CreatedAt time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt time.Time      `db:"updated_at" json:"updated_at"`
}

var (
	lock    = sync.Mutex{}
	allData map[string]bool
)

// LoadAllInMap try to load all translation from database into memory
func (m *Manager) LoadAllInMap(force bool) map[string]bool {
	lock.Lock()
	defer lock.Unlock()

	if allData != nil && !force {
		return allData
	}

	tmp := m.ListTranslations()

	res := make(map[string]bool)
	for i := range tmp {
		res[tmp[i].String] = true
	}

	return res
}

// AddMissing Add missing translation
func (m *Manager) AddMissing(txt string) error {
	lock.Lock()
	defer lock.Unlock()

	tmp := Translation{
		String: txt,
	}

	err := m.CreateTranslation(&tmp)
	if err != nil {
		return err
	}

	return nil
}
