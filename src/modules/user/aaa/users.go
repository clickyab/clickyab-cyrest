package aaa

import (
	"common/assert"
	"common/controllers/base"
	"common/models/common"
	"common/redis"
	"common/utils"
	"fmt"
	"modules/user/config"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gorp.v1"
)

// UserStatus is the registered user type
type (
	// UserStatus is the user status for a single use
	// @Enum{
	// }
	UserStatus string
	// UserType is the type of user
	// @Enum{
	// }
	UserType string
)

const (
	// UserStatusRegistered is the registered user, normal one
	// Registered
	UserStatusRegistered UserStatus = "registered"
	// UserStatusVerified for verified users
	UserStatusVerified UserStatus = "verified"
	// UserStatusBlocked for banned user
	UserStatusBlocked UserStatus = "blocked"

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
	ID          int64                              `db:"id" json:"id" sort:"true" title:"ID"`
	Email       string                             `db:"email" json:"email" search:"true" title:"Email"`
	Password    string                             `db:"password" json:"-"`
	OldPassword common.NullString                  `db:"old_password" json:"-"`
	AccessToken string                             `db:"access_token" json:"-"`
	Type        UserType                           `db:"user_type" json:"user_type" filter:"true" title:"User type"`
	ParentID    common.NullInt64                   `db:"parent_id" json:"-"`
	Avatar      common.NullString                  `db:"avatar" json:"avatar" visible:"false"`
	Status      UserStatus                         `db:"status" json:"status" filter:"true" title:"User status"`
	CreatedAt   time.Time                          `db:"created_at" json:"created_at" sort:"true" title:"Created at"`
	UpdatedAt   time.Time                          `db:"updated_at" json:"updated_at" sort:"true" title:"Created at"`
	resources   map[base.UserScope]map[string]bool `db:"-"`
	roles       []Role                             `db:"-"`
	//LastLogin   common.NullTime `db:"last_login" json:"last_login"`

	refreshToken bool `db:"-"`
}

//UserDataTable is the user full data in data table, after join with other field
// @DataTable {
//		url = /users
//		entity = user
//		view = user_list:parent
//		controller = modules/user/controllers
//		fill = FillUserDataTableArray
//		_edit = user_edit:global
// }
type UserDataTable struct {
	User
	ParentID int64 `db:"parent_id_dt" json:"parent_id" visible:"false"`
	OwnerID  int64 `db:"owner_id_dt" json:"owner_id" visible:"false"`
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
)

// Initialize the user on save
func (u *User) Initialize() {
	// TODO : Watch it if this creepy code is dangerous :)
	if (len(u.Password) < minHashSize ||
		!isBcrypt.MatchString(u.Password)) &&
		u.Password != noPassString {
		p, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		assert.Nil(err)
		u.Password = string(p)
		u.refreshToken = true
	}

	u.Email = strings.ToLower(u.Email)
	if u.refreshToken || u.AccessToken == "" || u.Status == UserStatusBlocked {
		u.AccessToken = <-utils.ID
		u.refreshToken = false
	}
}

// GetResources for this user
func (u *User) GetPermission() map[base.UserScope]map[string]bool {
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

// HasPerm try to return if the user has permission with some scope
// if requesting for lower scope and user has upper scope, the the maximum scope
// is returned, since the user with scope global, can do the scope self too
// this is different when the check is done with ids included
func (u *User) HasPerm(scope base.UserScope, perm string) (base.UserScope, bool) {
	if scope.IsValid() {
		return base.ScopeSelf, false
	}
	var (
		rScope base.UserScope = base.ScopeSelf
		rHas   bool
	)
	res := u.GetPermission()
	switch scope {
	case base.ScopeSelf:
		if res[scope][perm] {
			rScope = scope
			rHas = true
		}
		fallthrough
	case base.ScopeParent:
		if res[scope][perm] {
			rScope = scope
			rHas = true
		}
		fallthrough
	case base.ScopeGlobal:
		if res[scope][perm] {
			rScope = scope
			rHas = true
		}
	}
	return rScope, rHas
}

// HasPermOn check if user has permission on an object based on its owner id and its
// parent id
func (u *User) HasPermOn(perm string, ownerID, parentID int64, scopes ...base.UserScope) (base.UserScope, bool) {
	res := u.GetPermission()
	var (
		self, parent, global bool
	)
	if len(scopes) == 0 {
		self = true
		parent = true
		global = true
	} else {
		for i := range scopes {
			if scopes[i] == base.ScopeSelf {
				self = true
			} else if scopes[i] == base.ScopeParent {
				parent = true
			} else if scopes[i] == base.ScopeGlobal {
				global = true
			}
		}
	}

	if self {
		if ownerID == u.ID {
			if res[base.ScopeSelf][perm] {
				return base.ScopeSelf, true
			}
		}
	}
	if parent {
		if parentID == u.ID {
			if res[base.ScopeParent][perm] {
				return base.ScopeParent, true
			}
		}
	}

	if global {
		if res[base.ScopeGlobal][perm] {
			return base.ScopeGlobal, true
		}
	}
	return base.ScopeSelf, false
}

// FormatStatus is the example status formatter
func (u UserDataTable) FormatStatus() string {
	return string(u.Status)
}

// FillUserDataTableArray is the function to fill user data table array
func (m *Manager) FillUserDataTableArray(u base.PermInterfaceComplete, filters map[string]string, search map[string]string, sort, order string, p, c int) (UserDataTableArray, int64) {
	var params []interface{}
	var count int64
	var res UserDataTableArray
	var where []string

	countQuery := "SELECT COUNT(id) FROM users"
	query := "SELECT users.*,users.id AS owner_id_dt,users.parent_id as parent_id_dt FROM users"
	for field, value := range filters {
		where = append(where, fmt.Sprintf("%s=%s", field, "?"))
		params = append(params, value)
	}

	for column, val := range search {
		where = append(where, fmt.Sprintf("%s=%s", column, "?"))
		params = append(params, val)
	}

	currentUserID := u.GetID()
	highestScope := u.GetCurrentScope()

	if highestScope == base.ScopeSelf {
		where = append(where, "users.id=?")
		params = append(params, currentUserID)
	} else if highestScope == base.ScopeParent {
		where = append(where, "users.parent_id=?")
		params = append(params, currentUserID)
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

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d", sort, order, limit, offset)
	fmt.Println(query)
	fmt.Println(countQuery)
	/*_,err:=m.GetDbMap().Select(
		&res,
		query,
		params...,
	)
	assert.Nil(err)
	count,err=m.GetDbMap().SelectInt(
		query,
		params...,
	)
	assert.Nil(err)*/
	return res, count

}

// GetUserRoles return all Roles belong to User (many to many with UserRole)
func (m *Manager) GetUserRoles(u *User) []Role {
	var res []Role
	query := "SELECT roles.* FROM roles INNER JOIN user_role ON user_role.role_id=roles.id WHERE user_role.user_id=?"
	_, err := m.GetDbMap().Select(
		&res,
		query,
		u.ID,
	)

	assert.Nil(err)
	return res
}

// GetNewToken try to create a time base token in redis
func (m *Manager) GetNewToken(user *User, ua, ip string) string {
	t := fmt.Sprintf("%d:%s", user.ID, <-utils.ID)
	// TODO set at once
	assert.Nil(
		aredis.StoreHashKey(
			t,
			"token",
			user.AccessToken,
			ucfg.Cfg.TokenTimeout,
		),
	)
	assert.Nil(
		aredis.StoreHashKey(
			t,
			"ua",
			ua,
			ucfg.Cfg.TokenTimeout,
		),
	)
	assert.Nil(
		aredis.StoreHashKey(
			t,
			"ip",
			ip,
			ucfg.Cfg.TokenTimeout,
		),
	)
	assert.Nil(
		aredis.StoreHashKey(
			t,
			"date",
			time.Now().Format(time.RFC3339),
			ucfg.Cfg.TokenTimeout,
		),
	)
	return t
}

// RegisterUser is try for the user registration
func (m *Manager) RegisterUser(email, password string) (u *User, err error) {
	u = &User{
		Email:    email,
		Password: password,
		Status:   UserStatusRegistered,
		Type:     UserTypePersonal,
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
		return
	}

	return
}

// FetchByToken find user by its token in db
func (m *Manager) FetchByToken(accessToken string) (*User, error) {
	var res = User{}
	query := "SELECT * FROM users WHERE access_token=?"
	err := m.GetDbMap().SelectOne(
		&res,
		query,
		accessToken,
	)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// VerifyPassword try to verify password for given hash
func (u *User) VerifyPassword(password string) bool {
	// TODO : verify old password
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) == nil
}

// EraseToken remove a token from redis
func (m *Manager) EraseToken(t string) error {
	return aredis.RemoveKey(t)
}

// LogoutAllSession login from all user session
func (m *Manager) LogoutAllSession(u *User) error {
	u.refreshToken = true
	return m.UpdateUser(u)
}
