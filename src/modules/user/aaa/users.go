package aaa

import (
	"common/utils"
	"regexp"
	"time"

	"database/sql"

	"common/assert"

	"strings"

	"common/redis"
	"modules/user/config"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v1"
)

// UserStatus is the registered user type
type (
	// UserStatus is the user status for a single use
	// @Enum{
	// }
	UserStatus string
	// UserSource is the source of user
	// @Enum{
	// }
	UserSource string
	// UserType is the type of user
	// @Enum{
	// }
	UserType string
)

const (
	// UserStatusRegistered is the registered user, normal one
	UserStatusRegistered UserStatus = "registered"
	// UserStatusVerified for verified users
	UserStatusVerified UserStatus = "verified"
	// UserStatusBlocked for banned user
	UserStatusBlocked UserStatus = "blocked"

	// UserSourceCRM is the crm source
	UserSourceCRM UserSource = "crm"
	// UserSourceClickyab is the clickyab source
	UserSourceClickyab UserSource = "clickyab"

	// UserTypePersonal is the personal profile
	UserTypePersonal UserType = "personal"
	// UserTypeCorporation is the corp profile
	UserTypeCorporation UserType = "corpartion"
)

// User model
// @Model {
//		table = users
//		primary = true, id
//		find_by = id,email,access_token
//		transaction = insert
//		list = yes
// }
type User struct {
	ID          int64          `db:"id" json:"id"`
	Email       string         `db:"email" json:"email"`
	Password    sql.NullString `db:"password" json:"-"`
	OldPassword sql.NullString `db:"old_password" json:"-"`
	AccessToken string         `db:"access_token" json:"-"`
	Source      UserSource     `db:"source" json:"source"`
	Type        UserType       `db:"user_type" json:"user_type"`
	ParentID    sql.NullInt64  `db:"parent_id" json:"parent_id"`
	Avatar      sql.NullString `db:"avatar" json:"avatar"`
	Status      UserStatus     `db:"status" json:"status"`
	CreatedAt   time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at" json:"updated_at"`
	resources   map[string][]string     `db:"-"`
	roles           []Role                  `db:"-"`
	//LastLogin   common.NullTime `db:"last_login" json:"last_login"`

	refreshToken bool `db:"-"`
}

// CreateUserHook is the hook for create a user
type CreateUserHook func(gorp.SqlExecutor, *User) error

// From the bcrypt package
const (
	minHashSize  = 59
	noPassString = "NO" // Size must be less than 6 character
)

var (
	isBcrypt = regexp.MustCompile(`^\$[^$]+\$[0-9]+\$`)
	//hooks    []CreateUserHook
	//lock     = &sync.RWMutex{}
)

// Initialize the user on save
func (u *User) Initialize() {
	u.Email = strings.ToLower(u.Email)
	if u.refreshToken || u.AccessToken == "" || u.Status == UserStatusBlocked {
		u.AccessToken = <-utils.ID
		u.refreshToken = false
	}

	//if u.updateLastLogin {
	//	u.LastLogin = common.NullTime{
	//		Time:  time.Now(),
	//		Valid: true,
	//	}
	//	u.updateLastLogin = false
	//}

	// TODO : Watch it if this creepy code is dangerous :)
	if (len(u.Password.String) < minHashSize ||
		!isBcrypt.MatchString(u.Password.String)) &&
		u.Password.String != noPassString {
		p, err := bcrypt.GenerateFromPassword([]byte(u.Password.String), bcrypt.DefaultCost)
		assert.Nil(err)
		u.Password.String = string(p)
		u.Password.Valid = true
	}
}

// GetResources for this user
func (u *User) GetPermission() map[string][]string {
	if u.resources == nil {
		r := u.GetRoles()
		u.resources = NewAaaManager().GetPermissionMap(r...)
	}
	return u.resources
}

// GetRoles on this user
func (u *User) GetRoles() []Role {
	if u.roles == nil {
		m := NewAaaManager()
		u.roles = m.GetUserRoles(u)
	}

	return u.roles
}

// check if user has specified permission
func (u *User) HasPerm(resource string,scope string) bool{
	//get user permissions
	m:=u.GetPermission()
	if scope==string(ScopePermOwn){
		if !utils.StringInArray(resource,m[string(ScopePermOwn)]...) && !utils.StringInArray(resource,m[string(ScopePermParent)]...) && !utils.StringInArray(resource,m[string(ScopePermGlobal)]...){
			return false
		}
	}else if scope==string(ScopePermParent){
		if !utils.StringInArray(resource,m[string(ScopePermParent)]...) && !utils.StringInArray(resource,m[string(ScopePermGlobal)]...){
			return false
		}
	}else if scope==string(ScopePermGlobal){
		if !utils.StringInArray(resource,m[string(ScopePermGlobal)]...){
			return false
		}
	}
	return true
}

// GetUserRoles return all Roles belong to User (many to many with UserRole)
func (m *Manager) GetUserRoles(u *User) []Role {
	var res []Role
	query:="SELECT roles.* FROM roles INNER JOIN user_role ON user_role.role_id=roles.id WHERE user_role.user_id=?"
	_, err := m.GetDbMap().Select(
		&res,
		query,
		u.ID,
	)

	assert.Nil(err)
	return res
}
//
//// VerifyPassword try to verify password for given hash
//func (u *User) VerifyPassword(password string) bool {
//	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
//}
//
//// HasPassword check if user set password or not
//func (u *User) HasPassword() bool {
//	return u.Password != noPassString
//}
//
//// PostInsert is a gorp hook. need to add any other manager to use this hook,
//// for when the user is created in the system
//func (u *User) PostInsert(q gorp.SqlExecutor) error {
//	lock.RLock()
//	defer lock.RUnlock()
//	for i := range hooks {
//		err := hooks[i](q, u)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
//
//// GetRoles on this user
//func (u *User) GetRoles() []Role {
//	if u.roles == nil {
//		m := NewAaaManager()
//		u.roles = m.GetUserRoles(u)
//	}
//
//	return u.roles
//}
//
//// GetResources for this user
//func (u *User) GetResources() []string {
//	if u.resources == nil {
//		u.resources = make([]string, 0)
//		r := u.GetRoles()
//		for i := range r {
//			// TODO : what if this is not unique? no harm is done! let it be :/
//			u.resources = append(u.resources, r[i].Resources...)
//			//for _, resource := range r[i].Resources {
//			//	u.resources = utils.AppendIfMissing(u.resources, resource)
//			//}
//		}
//	}
//	return u.resources
//}
//
//// AddInsertUserHook adding a hook to run after the insert of new user
//func (m *Manager) AddInsertUserHook(f ...CreateUserHook) {
//	lock.Lock()
//	defer lock.Unlock()
//
//	hooks = append(hooks, f...)
//}
//

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

//
//// EraseToken remove the current token from the redis
//func (m *Manager) EraseToken(token string) error {
//	return aredis.RemoveKey(token)
//}
//
//// LoginUserByPassword try to login user with username and password
//func (m *Manager) LoginUserByPassword(username, password string) (string, *User, error) {
//	u, err := m.FindUserByUsername(username)
//	if err != nil {
//		return "", nil, err
//	}
//
//	if u.Status == UserStatusBanned {
//		return "", nil, errors.New("sorry, but you are banned")
//	}
//
//	if u.VerifyPassword(password) {
//		return m.GetNewToken(u.Token), u, nil
//	}
//
//	return "", nil, errors.New("wrong password")
//}
//
//// LoginUserByPassword try to login user with username and password
//func (m *Manager) LoginUserByOAuth(email string) (string, *User, error) {
//	u, err := m.FindUserByContact(email)
//	if err != nil {
//		return "", nil, err
//	}
//
//	if u.Status == UserStatusBanned {
//		return "", nil, errors.New("sorry, but you are banned")
//	}
//
//	return m.GetNewToken(u.Token), u, nil
//}
//
//// FindUserByIndirectToken try to find a user by its indirect token in database
//func (m *Manager) FindUserByIndirectToken(token string) (*User, error) {
//	t, err := aredis.GetKey(token, true, ucfg.Cfg.TokenTimeout)
//	if err != nil {
//		return nil, err
//	}
//
//	return m.FindUserByToken(t)
//}
//
//// LogoutAllSession login from all user session
//func (m *Manager) LogoutAllSession(u *User) error {
//	u.refreshToken = true
//	return m.UpdateUser(u)
//}
//
//// UpdateLastLogin try to update last login
//func (m *Manager) UpdateLastLogin(u *User) error {
//	u.updateLastLogin = true
//	return m.UpdateUser(u)
//}
//

// RegisterUser is try for the user registration
func (m *Manager) RegisterUser(email, password string, personal bool) (u *User, err error) {
	u = &User{
		Email:    email,
		Password: sql.NullString{String: password, Valid: true},
		Status:   UserStatusRegistered,
		Source: UserSourceClickyab,
		Type: UserTypePersonal,

		//updateLastLogin: true, // in this case, we need to update it since it means a login
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

// FetchByToken find user by its token in db
func (m *Manager)  FetchByToken(accessToken string)(*User ,error){
	var res =User{}
	query:="SELECT * FROM users WHERE access_token=?"
	err:=m.GetDbMap().SelectOne(
		&res,
		query,
		accessToken,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

//
//// RegisterUser is try for the user registration
//func (m *Manager) RegisterUserByContact(contact string) (u *User, err error) {
//	anonUser, err := m.GetDbMap().SelectInt("SELECT nextval('aaa.anon_user')")
//	if err != nil {
//		return nil, err
//	}
//	u = &User{
//		Contact:         contact,
//		Username:        fmt.Sprintf("user_%d", anonUser),
//		Password:        noPassString,
//		Attributes:      make(common.GenericJSONField),
//		Status:          UserStatusRegistered,
//		updateLastLogin: true, // in this case, we need to update it since it means a login
//	}
//	err = m.Begin()
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if err != nil {
//			assert.Nil(m.Rollback())
//		} else {
//			err = m.Commit()
//		}
//
//		if err != nil {
//			u = nil
//		}
//	}()
//	err = m.CreateUser(u)
//	if err != nil {
//		u = nil
//	}
//	return
//}
//
//// ListUserFilterByUsername and pagination
//func (m *Manager) ListUserFilterByUsername(offset, perPage int, username string, status UserStatus) ([]User, int64) {
//	username = strings.Trim(username, " ")
//	var where string
//	var params []interface{}
//	if username != "" {
//		where = "username ILIKE $1"
//		// TODO : like by pass all indexes :) fix this
//		params = append(params, "%"+username+"%")
//	}
//	if status.IsValid() {
//		if where != "" {
//			where += " AND "
//		}
//		where += fmt.Sprintf("status=$%d", len(params)+1)
//		params = append(params, status)
//	}
//	return m.ListUsersWithPaginationFilter(offset, perPage, where, params...), m.CountUsersWithFilter(where, params...)
//}
//
//// RegisterUserByToken is try for the user registration
//func (m *Manager) RegisterUserByToken(token, contact, username, password string) (u *User, err error) {
//	ru, err := m.FindReservedUserByContact(contact)
//	if err != nil {
//		return nil, err
//	}
//
//	if ru.Token != token {
//		return nil, fmt.Errorf(trans.T("invalid token"))
//	}
//
//	u = &User{
//		Contact:         contact,
//		Username:        username,
//		Password:        password,
//		Status:          UserStatusRegistered,
//		Attributes:      make(common.GenericJSONField),
//		updateLastLogin: true, // in this case, we need to update it since it means a login
//	}
//	err = m.Begin()
//	if err != nil {
//		return nil, err
//	}
//	defer func() {
//		if err != nil {
//			assert.Nil(m.Rollback())
//		} else {
//			err = m.Commit()
//		}
//
//		if err != nil {
//			u = nil
//		}
//	}()
//	err = m.CreateUser(u)
//	if err == nil {
//		_, err = m.GetDbMap().Delete(ru)
//	}
//
//	return
//}
//
//// ListUserByIDs try to load users by id
//func (m *Manager) ListUserByIDs(ids ...int64) []User {
//	params := make([]interface{}, len(ids))
//	for i := range ids {
//		params[i] = ids[i]
//	}
//	qs, params := utils.BuildPgPlaceHolder(1, params...)
//
//	return m.ListUsersWithFilter("IN("+strings.Join(qs, ",")+")", params...)
//}
