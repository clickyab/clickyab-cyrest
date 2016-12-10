package aaa

import (
	"common/utils"
	"regexp"
	"time"

	"common/assert"

	"strings"

	"common/redis"
	"modules/user/config"

	"common/models/common"

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
	ID          int64               `db:"id" json:"id"`
	Email       string              `db:"email" json:"email"`
	Password    common.NullString   `db:"password" json:"-"`
	OldPassword common.NullString   `db:"old_password" json:"-"`
	AccessToken string              `db:"access_token" json:"-"`
	Source      UserSource          `db:"source" json:"source"`
	Type        UserType            `db:"user_type" json:"user_type"`
	ParentID    common.NullInt64    `db:"parent_id" json:"parent_id"`
	Avatar      common.NullString   `db:"avatar" json:"avatar"`
	Status      UserStatus          `db:"status" json:"status"`
	CreatedAt   time.Time           `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time           `db:"updated_at" json:"updated_at"`
	resources   map[ScopePerm]map[string]bool `db:"-"`
	roles       []Role              `db:"-"`
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
	// TODO : Watch it if this creepy code is dangerous :)
	if (len(u.Password.String) < minHashSize ||
		!isBcrypt.MatchString(u.Password.String)) &&
		u.Password.String != noPassString {
		p, err := bcrypt.GenerateFromPassword([]byte(u.Password.String), bcrypt.DefaultCost)
		assert.Nil(err)
		u.Password.String = string(p)
		u.Password.Valid = true
		u.refreshToken = true
	}

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

}

// GetResources for this user
func (u *User) GetPermission() map[ScopePerm]map[string]bool {
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
func (u *User) HasPerm(scope ScopePerm, perm string) (ScopePerm, bool) {
	if scope.IsValid() {
		return ScopePermOwn, false
	}
	var (
		rScope ScopePerm = ScopePermOwn
		rHas   bool
	)
	res := u.GetPermission()
	switch scope {
	case ScopePermOwn:
		if res[scope][perm] {
			rScope = scope
			rHas = true
		}
		fallthrough
	case ScopePermParent:
		if res[scope][perm] {
			rScope = scope
			rHas = true
		}
		fallthrough
	case ScopePermGlobal:
		if res[scope][perm] {
			rScope = scope
			rHas = true
		}
	}
	return rScope, rHas
}

// HasPermOn check if user has permission on an object based on its owner id and its
// parent id
func (u *User) HasPermOn(perm string, ownerID, parentID int64, scopes ...ScopePerm) (ScopePerm, bool) {
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
			if scopes[i] == ScopePermOwn {
				self = true
			} else if scopes[i] == ScopePermParent {
				parent = true
			} else if scopes[i] == ScopePermGlobal {
				global = true
			}
		}
	}

	if self {
		if ownerID == u.ID {
			if res[ScopePermOwn][perm] {
				return ScopePermOwn, true
			}
		}
	}
	if parent {
		if parentID == u.ID {
			if res[ScopePermParent][perm] {
				return ScopePermParent, true
			}
		}
	}

	if global {
		if res[ScopePermGlobal][perm] {
			return ScopePermGlobal, true
		}
	}
	return ScopePermOwn, false
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

// RegisterUser is try for the user registration
func (m *Manager) RegisterUser(email, password string, profile interface{}) (u *User, err error) {
	u = &User{
		Email:    email,
		Password: common.NullString{String: password, Valid: true},
		Status:   UserStatusRegistered,
		Source:   UserSourceClickyab,
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

	switch p := profile.(type) {
	case *UserProfileCorporation:
		p.UserID = u.ID
		err = m.CreateUserProfileCorporation(p)
	case *UserProfilePersonal:
		p.UserID = u.ID
		err = m.CreateUserProfilePersonal(p)
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
	return bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(password)) == nil
}
