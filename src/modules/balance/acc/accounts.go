package acc

import (
	"common/assert"
	"common/models/common"
	"common/utils"
	"fmt"
	"modules/misc/trans"
	"modules/user/aaa"
	"strings"
	"sync"
	"time"

	"gopkg.in/gorp.v1"
)

// Accounts represent the account model in database
// @Model {
//		table = accounts
//		schema = acc
//		primary = true, id
//		find_by = id
//		list = yes
//		belong_to = aaa.User:owner_id
//		transaction = insert, update, delete
// }
type Account struct {
	ID          int64                   `db:"id" json:"id"`
	OwnerID     int64                   `db:"owner_id" json:"owner_id"`
	Title       string                  `db:"title" json:"title"`
	Description common.NullString       `db:"description" json:"description"`
	CreatedAt   time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time               `db:"updated_at" json:"updated_at"`
	Attributes  common.GenericJSONField `db:"attributes" json:"attributes"`
	Disabled    bool                    `db:"disabled" json:"disabled"`
	Current     int64                   `db:"current" json:"current"`
	user        *aaa.User               `db:"-"` // A hook required data
}

type AccountHookFunc func(gorp.SqlExecutor, string, *Account, *aaa.User) error

var (
	accountHooks []AccountHookFunc
	accountLock  = &sync.RWMutex{}
)

func (a *Account) runHooks(s gorp.SqlExecutor, action string) error {
	accountLock.RLock()
	defer accountLock.RUnlock()
	if a.user == nil {
		return fmt.Errorf("user record is not available")
	}
	for _, f := range accountHooks {
		err := f(s, action, a, a.user)
		if err != nil {
			return err
		}
	}
	return nil
}

// PostInsert is a hook to run tha actual hooks :)
func (a *Account) PostInsert(s gorp.SqlExecutor) error {
	return a.runHooks(s, "INSERT")
}

// PostInsert is a hook to run tha actual hooks :)
func (a *Account) PostUpdate(s gorp.SqlExecutor) error {
	return a.runHooks(s, "UPDATE")
}

// PostInsert is a hook to run tha actual hooks :)
func (a *Account) PostDelete(s gorp.SqlExecutor) error {
	return a.runHooks(s, "DELETE")
}

// AddAccount is for creating an account. this action must be synced
func (m *Manager) AddAccount(
	u *aaa.User,
	initial Money,
	title string,
	description string,
	disabled bool) (acc *Account, err error) {
	err = m.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			acc = nil
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}
	}()
	t := time.Now()
	acc = &Account{
		OwnerID:     u.ID,
		Title:       title,
		Description: common.NullString{Valid: true, String: description},
		CreatedAt:   t,
		UpdatedAt:   t,
		Disabled:    disabled,
		user:        u,
		Attributes:  common.GenericJSONField{"initial": initial},
	}
	err = m.CreateAccount(acc)
	if err != nil && initial != 0 {
		_, err = m.AddTransaction(
			u,
			acc,
			initial,
			acc.CreatedAt,
			trans.T("create_account_initial"),
			trans.T("income"),
			trans.T("initial"),
		)
	}
	return
}

// EditAccount try to edit the account
func (m *Manager) EditAccount(
	u *aaa.User,
	acc *Account,
	title string,
	description string,
	disabled bool,
) (err error) {
	acc.user = u
	assert.True(u.ID == acc.OwnerID, "[BUG] not the owner of account")
	err = m.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}
	}()
	t := time.Now()

	acc.Title = title
	acc.Description = common.NullString{Valid: true, String: description}
	acc.UpdatedAt = t
	acc.Disabled = disabled
	return m.UpdateAccount(acc)
}

// EditAccount try to delete the account. if there is any transaction related to that account is exists,
// simply disable it
func (m *Manager) DeleteAccount(u *aaa.User, acc *Account) (fully bool, err error) {
	acc.user = u
	assert.True(u.ID == acc.OwnerID, "[BUG] not the owner of account")
	err = m.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}
	}()

	// Count the transactions belong to this account
	cnt := m.CountAccountTransactions(acc)
	if cnt == 0 {
		// we can remove it safely
		_, err = m.GetDbMap().Delete(acc)
		fully = true
		return
	}
	t := time.Now()
	acc.Disabled = true
	acc.UpdatedAt = t
	err = m.UpdateAccount(acc)
	return
}

// ListAccountInUnit return the list of all account in this unit. filter with disable, read only for this user
// and also pagination support
func (m *Manager) ListUserAccounts(u *aaa.User, dis bool, all bool, offset, page int) ([]Account, int64) {
	filter, parameter := " owner_id=$1 ", []interface{}{u.ID}
	if !dis {
		filter += " AND NOT disabled "
	}
	var list []Account
	if all {
		list = m.ListAccountsWithFilter(filter, parameter...)
	} else {
		list = m.ListAccountsWithPaginationFilter(offset, page, filter, parameter...)
	}

	return list, m.CountAccountsWithFilter(filter, parameter...)
}

// ListUserAccountsWithID try to translate some id to account object
func (m *Manager) ListUserAccountsWithID(u *aaa.User, accounts ...int64) []Account {
	var res []Account
	sql, params := []string{" a.owner_id=$1 "}, []interface{}{u.ID}
	if len(accounts) > 0 {
		tmp := make([]interface{}, len(accounts))
		for i := range accounts {
			tmp[i] = accounts[i]
		}
		str, data := utils.BuildPgPlaceHolder(2, tmp...)
		params = append(params, data...)
		sql = append(sql, fmt.Sprintf(" a.id IN (%s) ", strings.Join(str, ",")))
	} else {
		return res
	}
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s a WHERE %s ",
			AccountTableFull,
			strings.Join(sql, " AND "),
		),
		params...,
	)
	assert.Nil(err)
	return res
}

// AddAccountHook is for adding hook to insert/update hook
func (m *Manager) AddAccountHook(f ...AccountHookFunc) {
	accountLock.Lock()
	defer accountLock.Unlock()

	accountHooks = append(accountHooks, f...)
}
