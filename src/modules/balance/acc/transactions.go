package acc

import (
	"common/assert"
	"common/models/common"
	"fmt"
	"modules/misc/trans"
	"modules/user/aaa"
	"sync"
	"time"

	"common/redis"

	"strings"

	"common/utils"

	"gopkg.in/gorp.v1"
)

// Transaction is the core table of the system. all transactions are in this table
// so far
// @Model {
//		table = transactions
//		schema = acc
//		primary = true, id
//		find_by = id
//		list = yes
//		filter_by = account_id
//		transaction = insert, update
//		belong_to = Account:account_id
// }
type Transaction struct {
	ID          int64     `db:"id" json:"id"`
	AccountID   int64     `db:"account_id" json:"account_id"`
	UserID      int64     `db:"user_id" json:"user_id"`
	Amount      Money     `db:"amount" json:"amount"`
	Description string    `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// UnitTags is my try to achieve tags in units
// @Model {
//		table = transaction_tags
//		schema = acc
//		primary = false, transaction_id
//		find_by = transaction_id
//		transaction = insert, update
// }
type TransactionTags struct {
	TransactionID int64              `db:"transaction_id" json:"transaction_id"`
	Tags          common.StringSlice `db:"tags" json:"tags"`
	transaction   *Transaction       `db:"-"` // This is simply the one for firing the hook func
	account       *Account           `db:"-"`
}

type TransactionHookFunc func(gorp.SqlExecutor, string, *Account, *Transaction, *TransactionTags) error

var (
	transactionHook []TransactionHookFunc
	transactionLock = &sync.RWMutex{}
)

// PostInsert is a hook for running hooks :)
func (tt *TransactionTags) PostInsert(s gorp.SqlExecutor) error {
	transactionLock.RLock()
	defer transactionLock.RUnlock()
	if tt.transaction == nil {
		return fmt.Errorf("transaction record is not available")
	}
	for i := range transactionHook {
		err := transactionHook[i](s, "INSERT", tt.account, tt.transaction, tt)
		if err != nil {
			return err
		}
	}

	return nil
}

// DoTransaction is the key part of transaction insert in the system. any transaction must pass this
// and must not accept any other path.
func (m *Manager) AddTransaction(
	u *aaa.User,
	acc *Account,
	amount Money,
	created time.Time,
	desc string,
	correlationID string,
	tags ...string,
) (t *Transaction, err error) {
	defer func() {
		// make sure the transaction is null if there is an error
		if err != nil {
			t = nil
		}
	}()
	if len(tags) < 1 {
		err = fmt.Errorf(trans.T("atleast_one_tag"))
		return
	}

	if amount == 0 {
		err = fmt.Errorf(trans.T("empty_amount"))
		return
	}

	// Its time to check for correlation id in the past few days transaction for this user
	key := fmt.Sprint("transaction_%s_%d_%d_%d", correlationID, u.ID, acc.ID, amount)
	data := desc + strings.Join(tags, "_")
	str, err := aredis.GetKey(
		key,
		false,
		time.Second,
	)
	if err == nil {
		// Oh, heck. possible duplicate entry?
		if str == data {
			return nil, fmt.Errorf(trans.T("possible_duplicate"))
		}
	}
	defer func() {
		// Set the redis key on success
		if err == nil {
			// TODO : Get this from config?
			assert.Nil(
				aredis.StoreKey(
					key,
					data,
					time.Hour*48,
				),
			)
		}
	}()

	already := m.InTransaction()
	if !already {
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

			if err != nil {
				t = nil
			}
		}()
	}
	t = &Transaction{
		AccountID:   acc.ID,
		UserID:      u.ID,
		Amount:      amount,
		Description: desc,
		CreatedAt:   created,
		UpdatedAt:   created,
	}

	if err = m.CreateTransaction(t); err != nil {
		return
	}

	tt := &TransactionTags{
		TransactionID: t.ID,
		Tags:          make(common.StringSlice, len(tags)),
		transaction:   t,
		account:       acc,
	}
	for i := range tags {
		tt.Tags[i] = tags[i]
	}
	err = m.CreateTransactionTags(tt)
	return
}

// FindTransactionByID return the Transaction base on its id
func (m *Manager) FindTransactionByAccountAndID(a *Account, id int64) (*Transaction, error) {
	var res Transaction
	err := m.GetDbMap().SelectOne(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s t JOIN %s a ON a.id=t.account_id WHERE t.id=$1 AND a.id=$2",
			TransactionTableFull,
			AccountTableFull,
		),
		id,
		a.ID,
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

// ListTransactionsByUserAndAccounts return all transactions in a list, with pagination and also with
// accounts filter
func (m *Manager) ListTransactionsByUserAndAccounts(u *aaa.User, offset, count int, accounts ...Account) []Transaction {
	sql, params := []string{" a.owner_id=$1 "}, []interface{}{u.ID}
	if len(accounts) > 0 {
		tmp := make([]interface{}, len(accounts))
		for i := range accounts {
			tmp[i] = accounts[i].ID
		}
		str, data := utils.BuildPgPlaceHolder(2, tmp...)
		params = append(params, data...)
		sql = append(sql, fmt.Sprintf(" a.id IN (%s) ", strings.Join(str, ",")))
	}
	params = append(params, offset, count)
	var res []Transaction
	_, err := m.GetDbMap().Select(
		&res,
		fmt.Sprintf(
			"SELECT * FROM %s t JOIN %s a WHERE %s ORDER BY t.created_at ASC OFFSET $%d LIMIT $%d",
			TransactionTableFull,
			AccountTableFull,
			strings.Join(sql, " AND "),
			len(params)-1,
			len(params),
		),
		params...,
	)
	assert.Nil(err)
	return res
}

// CountTransactionsByUserAndAccounts count all transactions in the accounts
func (m *Manager) CountTransactionsByUserAndAccounts(u *aaa.User, accounts ...Account) int64 {
	sql, params := []string{" a.owner_id=$1 "}, []interface{}{u.ID}
	if len(accounts) > 0 {
		tmp := make([]interface{}, len(accounts))
		for i := range accounts {
			tmp[i] = accounts[i].ID
		}
		str, data := utils.BuildPgPlaceHolder(2, tmp...)
		params = append(params, data...)
		sql = append(sql, fmt.Sprintf(" a.id IN (%s) ", strings.Join(str, ",")))
	}
	cnt, err := m.GetDbMap().SelectInt(
		fmt.Sprintf(
			"SELECT COUNT(*) FROM %s t JOIN %s a WHERE %s ",
			TransactionTableFull,
			AccountTableFull,
			strings.Join(sql, " AND "),
		),
		params...,
	)
	assert.Nil(err)
	return cnt
}

// AddTransactionHook add a new transaction hook to the hook list
func (m *Manager) AddTransactionHook(f ...TransactionHookFunc) {
	transactionLock.Lock()
	defer transactionLock.Unlock()

	transactionHook = append(transactionHook, f...)
}
