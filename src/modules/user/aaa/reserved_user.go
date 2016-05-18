package aaa

import (
	"fmt"
	"math/rand"
	"modules/misc/trans"
	"modules/user/utils"
	"time"
)

// ReservedUser model
// @Model {
//		table = reserved_users
//		schema = aaa
//		primary = true, id
//		find_by = id,contact
//		list = yes
// }
type ReservedUser struct {
	ID        int64     `db:"id" json:"id"`
	Contact   string    `db:"contact" json:"contact"`
	Token     string    `db:"token" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Initialize try to set the token if the token is not available
func (ru *ReservedUser) Initialize() {
	if ru.Token == "" {
		ru.Token = fmt.Sprintf("%d", rand.Intn(8000)+1000)
	}
}

// ReserveToken for a registration, return the reserved token, if this is the first time or send is done
// already, and error in case of error
func (m *Manager) ReserveToken(contact string) (*ReservedUser, bool, error) {
	t, err := utils.DetectContactType(contact)
	if err != nil {
		return nil, false, err
	}
	// Remove the leading zero
	if t == utils.TypePhone && contact[0:1] == "0" {
		contact = contact[1:]
	}
	_, err = m.FindUserByContact(contact)
	if err == nil {
		return nil, false, fmt.Errorf(trans.T("this is already registered"))
	}

	ru, err := m.FindReservedUserByContact(contact)
	if err == nil {
		_ = m.UpdateReservedUser(ru)
		return ru, true, nil
	}

	ru = &ReservedUser{
		Contact: contact,
	}
	err = m.CreateReservedUser(ru)
	if err != nil {
		return nil, false, err
	}

	return ru, false, nil
}
