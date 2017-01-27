package base

import (
	"common/assert"
	"sync"

	"github.com/Sirupsen/logrus"
)

// UserScope is the permission level for a role
// @Enum {
// }
type UserScope string

// Permission is the type to handle permission
type Permission string

const (
	// ScopeSelf means the user him self, no additional parameter
	ScopeSelf UserScope = "self"
	// ScopeParent means the user child, need id of all child as parameter
	ScopeParent UserScope = "parent"
	// ScopeGlobal means the entire perm, no param is required
	ScopeGlobal UserScope = "global"
)

const (
	// PermissionGod is the god for perms
	PermissionGod Permission = "god"
)

var (
	registeredPerms = make(map[Permission]string)
	lock            = &sync.RWMutex{}
)

// PermInterface is the perm interface
type PermInterface interface {
	// HasPerm is the has perm check
	HasPerm(scope UserScope, perm Permission) (UserScope, bool)
	// HasPermOn is the has perm on check
	HasPermOn(perm Permission, ownerID, parentID int64, scopes ...UserScope) (UserScope, bool)
}

// PermInterfaceComplete is the complete version of the interface to use
type PermInterfaceComplete interface {
	PermInterface
	// GetID return the id of holder
	GetID() int64
	// GetCurrentPerm return the current permission that this object is built on
	GetCurrentPerm() Permission
	// GetCurrentScope return the current scope for this object (maximum)
	GetCurrentScope() UserScope
}

type permComplete struct {
	inner PermInterface

	id    int64
	perm  Permission
	scope UserScope
}

// HasPermString is the has perm check
func (pc permComplete) HasPerm(scope UserScope, perm Permission) (UserScope, bool) {
	return pc.inner.HasPerm(scope, perm)
}

// HasPermStringOn is the has perm on check
func (pc permComplete) HasPermOn(perm Permission, ownerID, parentID int64, scopes ...UserScope) (UserScope, bool) {
	return pc.inner.HasPermOn(perm, ownerID, parentID, scopes...)
}

// GetID return the id of holder
func (pc permComplete) GetID() int64 {
	return pc.id
}

// GetCurrentPerm return the current permission that this object is built on
func (pc permComplete) GetCurrentPerm() Permission {
	return pc.perm
}

// GetCurrentScope return the current scope for this object (maximum)
func (pc permComplete) GetCurrentScope() UserScope {
	return pc.scope
}

// RegisterPermission register a permission
func RegisterPermission(perm Permission, description string) {
	lock.Lock()
	defer lock.Unlock()

	registeredPerms[perm] = description
}

// PermissionRegistered check if the permission is registered in system or not
// and just log it
func PermissionRegistered(perm Permission) {
	lock.RLock()
	defer lock.RUnlock()

	if _, ok := registeredPerms[perm]; !ok {
		logrus.Panicf("The permission is not registered %s", perm)
	}

}

// GetAllPermission return the permission list in system
func GetAllPermission() map[Permission]string {
	lock.RLock()
	defer lock.RUnlock()

	return registeredPerms
}

// NewPermInterfaceComplete return a new object base on the minimum object
func NewPermInterfaceComplete(inner PermInterface, id int64, perm Permission, scope UserScope) PermInterfaceComplete {
	s, ok := inner.HasPerm(scope, perm)
	if !ok {
		s, ok = inner.HasPerm(ScopeGlobal, PermissionGod)
	}
	assert.True(ok, "[BUG] probably there is some thing wrong with code generation")
	pc := &permComplete{
		inner: inner,
		id:    id,
		perm:  perm,
		scope: s,
	}

	return pc
}

func init() {
	RegisterPermission(PermissionGod, "the god of all permissions")
}
