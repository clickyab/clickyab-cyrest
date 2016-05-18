package acc

import (
	"common/models/common"
	"database/sql"
	"database/sql/driver"
	"modules/misc/trans"
	"modules/user/aaa"

	"github.com/Sirupsen/logrus"
	"gopkg.in/gorp.v1"
)

// Money is the type for money. I can use int64 in Iran, but I prefer to cast it so I
// can simply change it in any case
type Money int64

// Scan convert the json array ino string slice
func (is *Money) Scan(src interface{}) error {
	tmp := &sql.NullInt64{}
	err := tmp.Scan(src)
	if err != nil {
		return err
	}

	*is = Money(tmp.Int64)
	return nil
}

// Value try to get the string slice representation in database
func (is Money) Value() (driver.Value, error) {
	return int64(is), nil
}

// insertUserHook is called after creating a user
func insertUserHook(s gorp.SqlExecutor, u *aaa.User) (err error) {
	defer func() {
		if err != nil {
			logrus.Warn(err)
		}
	}()
	m, err := NewAccManagerFromTransaction(s)
	if err != nil {
		return err
	}

	// OK its time to create new Personal Wallet in this unit
	pw := &Account{
		OwnerID:     u.ID,
		Title:       trans.T("personal_wallet"),
		Description: common.NullString{Valid: true, String: trans.T("default_account")},
		user:        u, // this is required for the hooks
	}

	if err = m.CreateAccount(pw); err != nil {
		return err
	}

	// TODO : Create default tag-ing structure

	return nil
}

// AddHook is called after initialize all modules
func (m *Manager) AddHook() {
	um := aaa.NewAaaManager()
	um.AddInsertUserHook(insertUserHook)
}
