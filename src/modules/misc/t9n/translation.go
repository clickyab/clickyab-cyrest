package t9n

import (
	"database/sql"
	"time"
)

// Translation model
// @Model {
//		table = translations
//		schema = t9n
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

// LoadAllInMap try to load all translation from database into memory
func (m *Manager) LoadAllInMap() map[string]Translation {
	tmp := m.ListTranslations()

	res := make(map[string]Translation)
	for i := range tmp {
		res[tmp[i].String] = tmp[i]
	}

	return res
}

// AddMissing translation
func (m *Manager) AddMissing(txt string) (Translation, error) {
	tmp := Translation{
		String: txt,
	}

	err := m.CreateTranslation(&tmp)
	if err != nil {
		return Translation{}, err
	}

	return tmp, nil
}
