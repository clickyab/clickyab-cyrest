package aaa

import (
	"common/assert"
	"common/models/common"
	"common/redis"
	"common/utils"
	"database/sql/driver"
	"errors"
	"fmt"
	"modules/misc/trans"
	"modules/user/config"
	"regexp"
	"strings"
	"time"

	"sync"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v1"
)

// UserStatus is the registered user type
type UserStatus string

const (
	// UserStatusRegistered is the registered user, normal one
	UserStatusRegistered UserStatus = "registered"
	// UserStatusVerified for verified users
	UserStatusVerified = "verified"
	// UserStatusBanned for banned user
	UserStatusBanned = "banned"
)

// User model
// @Model {
//		table = users
//		schema = aaa
//		primary = true, id
//		find_by = id,username,token,contact
//		transaction = insert
//		list = yes
// }
type User struct {
	ID              int64                   `db:"id" json:"id"`
	Username        string                  `db:"username" json:"username"`
	Password        string                  `db:"password" json:"-"`
	Contact         string                  `db:"contact" json:"contact"`
	Token           string                  `db:"token" json:"-"`
	CreatedAt       time.Time               `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time               `db:"updated_at" json:"updated_at"`
	LastLogin       common.NullTime         `db:"last_login" json:"last_login"`
	Attributes      common.GenericJSONField `db:"attributes" json:"attributes"`
	Status          UserStatus              `db:"status" json:"status"`
	roles           []Role                  `db:"-"`
	resources       []string                `db:"-"`
	refreshToken    bool                    `db:"-"`
	updateLastLogin bool                    `db:"-"`
}

type CreateUserHook func(gorp.SqlExecutor, *User) error

// From the bcrypt package
const (
	minHashSize = 59
	noPassString = "NO" // Size must be less than 6 character
)

var (
	isBcrypt = regexp.MustCompile(`^\$[^$]+\$[0-9]+\$`)
	hooks    []CreateUserHook
	lock     = &sync.RWMutex{}
)

// IsValid try to validate enum value on ths type
func (is UserStatus) IsValid() bool {
	return utils.StringInArray(string(is), string(UserStatusBanned), string(UserStatusVerified), string(UserStatusRegistered))
}

// Scan convert the json array ino string slice
func (is *UserStatus) Scan(src interface{}) error {
	var b []byte
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		b = make([]byte, 0)
	default:
		return errors.New("unsupported type")
	}
	if !UserStatus(b).IsValid() {
		return fmt.Errorf("invaid status, valids are : %s, %s, %s", UserStatusBanned, UserStatusRegistered, UserStatusVerified)
	}
	*is = UserStatus(b)
	return nil
}

// Value try to get the string slice representation in database
func (is UserStatus) Value() (driver.Value, error) {
	if !is.IsValid() {
		return nil, fmt.Errorf("invaid status, valids are : %s, %s, %s", UserStatusBanned, UserStatusRegistered, UserStatusVerified)
	}
	return string(is), nil
}

// Initialize the user on save
func (u *User) Initialize() {
	u.Username = strings.ToLower(strings.Trim(u.Username, " \n\t"))
	if u.refreshToken || u.Token == "" || u.Status == UserStatusBanned {
		u.Token = <-utils.ID
		u.refreshToken = false
	}

	if u.updateLastLogin {
		u.LastLogin = common.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		u.updateLastLogin = false
	}

	// TODO : Watch it if this creepy code is dangerous :)
	if (len(u.Password) < minHashSize || !isBcrypt.MatchString(u.Password)) && u.Password != noPassString {
		p, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		assert.Nil(err)
		u.Password = string(p)
	}
}

// VerifyPassword try to verify password for given hash
func (u *User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// HasPassword check if user set password or not
func (u *User) HasPassword() bool {
	return u.Password != noPassString
}

// PostInsert is a gorp hook. need to add any other manager to use this hook,
// for when the user is created in the system
func (u *User) PostInsert(q gorp.SqlExecutor) error {
	lock.RLock()
	defer lock.RUnlock()
	for i := range hooks {
		err := hooks[i](q, u)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetRoles on this user
func (u *User) GetRoles() []Role {
	if u.roles == nil {
		m := NewAaaManager()
		u.roles = m.GetUserRoles(u)
	}

	return u.roles
}

// GetResources for this user
func (u *User) GetResources() []string {
	if u.resources == nil {
		u.resources = make([]string, 0)
		r := u.GetRoles()
		for i := range r {
			// TODO : what if this is not unique? no harm is done! let it be :/
			u.resources = append(u.resources, r[i].Resources...)
			//for _, resource := range r[i].Resources {
			//	u.resources = utils.AppendIfMissing(u.resources, resource)
			//}
		}
	}
	return u.resources
}

// AddInsertUserHook adding a hook to run after the insert of new user
func (m *Manager) AddInsertUserHook(f ...CreateUserHook) {
	lock.Lock()
	defer lock.Unlock()

	hooks = append(hooks, f...)
}

// GetNewToken try to create a time base token in redis
func (m *Manager) GetNewToken(baseToken string) string {
	t := <-utils.ID
	assert.Nil(
		aredis.StoreKey(
			t,
			baseToken,
			ucfg.Cfg.TokenTimeout,
		),
	)

	return t
}

// EraseToken remove the current token from the redis
func (m *Manager) EraseToken(token string) error {
	return aredis.RemoveKey(token)
}

// LoginUserByPassword try to login user with username and password
func (m *Manager) LoginUserByPassword(username, password string) (string, *User, error) {
	u, err := m.FindUserByUsername(username)
	if err != nil {
		return "", nil, err
	}

	if u.Status == UserStatusBanned {
		return "", nil, errors.New("sorry, but you are banned")
	}

	if u.VerifyPassword(password) {
		return m.GetNewToken(u.Token), u, nil
	}

	return "", nil, errors.New("wrong password")
}

// LoginUserByPassword try to login user with username and password
func (m *Manager) LoginUserByOAuth(email string) (string, *User, error) {
	u, err := m.FindUserByContact(email)
	if err != nil {
		return "", nil, err
	}
	
	if u.Status == UserStatusBanned {
		return "", nil, errors.New("sorry, but you are banned")
	}
	
	return m.GetNewToken(u.Token), u, nil
}

// FindUserByIndirectToken try to find a user by its indirect token in database
func (m *Manager) FindUserByIndirectToken(token string) (*User, error) {
	t, err := aredis.GetKey(token, true, ucfg.Cfg.TokenTimeout)
	if err != nil {
		return nil, err
	}

	return m.FindUserByToken(t)
}

// LogoutAllSession login from all user session
func (m *Manager) LogoutAllSession(u *User) error {
	u.refreshToken = true
	return m.UpdateUser(u)
}

// UpdateLastLogin try to update last login
func (m *Manager) UpdateLastLogin(u *User) error {
	u.updateLastLogin = true
	return m.UpdateUser(u)
}

// RegisterUser is try for the user registration
func (m *Manager) RegisterUser(contact, username, password string) (u *User, err error) {
	u = &User{
		Contact:         contact,
		Username:        username,
		Password:        password,
		Attributes:      make(common.GenericJSONField),
		Status:          UserStatusRegistered,
		updateLastLogin: true, // in this case, we need to update it since it means a login
	}
	err = m.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}
		
		if err != nil {
			u = nil
		}
	}()
	err = m.CreateUser(u)
	if err != nil {
		u = nil
	}
	return
}

// RegisterUser is try for the user registration
func (m *Manager) RegisterUserByContact(contact string) (u *User, err error) {
	anonUser, err := m.GetDbMap().SelectInt("SELECT nextval('aaa.anon_user')")
	if err != nil {
		return nil, err
	}
	u = &User{
		Contact:         contact,
		Username:        fmt.Sprintf("user_%d",anonUser),
		Password:        noPassString,
		Attributes:      make(common.GenericJSONField),
		Status:          UserStatusRegistered,
		updateLastLogin: true, // in this case, we need to update it since it means a login
	}
	err = m.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}
		
		if err != nil {
			u = nil
		}
	}()
	err = m.CreateUser(u)
	if err != nil {
		u = nil
	}
	return
}

// ListUserFilterByUsername and pagination
func (m *Manager) ListUserFilterByUsername(offset, perPage int, username string, status UserStatus) ([]User, int64) {
	username = strings.Trim(username, " ")
	var where string
	var params []interface{}
	if username != "" {
		where = "username ILIKE $1"
		// TODO : like by pass all indexes :) fix this
		params = append(params, "%"+username+"%")
	}
	if status.IsValid() {
		if where != "" {
			where += " AND "
		}
		where += fmt.Sprintf("status=$%d", len(params)+1)
		params = append(params, status)
	}
	return m.ListUsersWithPaginationFilter(offset, perPage, where, params...), m.CountUsersWithFilter(where, params...)
}

// RegisterUserByToken is try for the user registration
func (m *Manager) RegisterUserByToken(token, contact, username, password string) (u *User, err error) {
	ru, err := m.FindReservedUserByContact(contact)
	if err != nil {
		return nil, err
	}

	if ru.Token != token {
		return nil, fmt.Errorf(trans.T("invalid token"))
	}

	u = &User{
		Contact:         contact,
		Username:        username,
		Password:        password,
		Status:          UserStatusRegistered,
		Attributes:      make(common.GenericJSONField),
		updateLastLogin: true, // in this case, we need to update it since it means a login
	}
	err = m.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			assert.Nil(m.Rollback())
		} else {
			err = m.Commit()
		}

		if err != nil {
			u = nil
		}
	}()
	err = m.CreateUser(u)
	if err == nil {
		_, err = m.GetDbMap().Delete(ru)
	}

	return
}

// ListUserByIDs try to load users by id
func (m *Manager) ListUserByIDs(ids ...int64) []User {
	params := make([]interface{}, len(ids))
	for i := range ids {
		params[i] = ids[i]
	}
	qs, params := utils.BuildPgPlaceHolder(1, params...)

	return m.ListUsersWithFilter("IN("+strings.Join(qs, ",")+")", params...)
}
