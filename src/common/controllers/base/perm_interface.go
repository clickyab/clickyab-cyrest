package base

import "common/assert"

// PermInterface is the perm interface
type PermInterface interface {
	// HasPermString is the has perm check
	HasPermString(scope string, perm string) (string, bool)
	// HasPermStringOn is the has perm on check
	HasPermStringOn(perm string, ownerID, parentID int64, scopes ...string) (string, bool)
}

// PermInterfaceComplete is the complete version of the interface to use
type PermInterfaceComplete interface {
	PermInterface
	// GetID return the id of holder
	GetID() int64
	// GetCurrentPerm return the current permission that this object is built on
	GetCurrentPerm() string
	// GetCurrentScope return the current scope for this object (maximum)
	GetCurrentScope() string
}

type permComplete struct {
	inner PermInterface

	id    int64
	perm  string
	scope string
}

// HasPermString is the has perm check
func (pc permComplete) HasPermString(scope string, perm string) (string, bool) {
	return pc.inner.HasPermString(scope, perm)
}

// HasPermStringOn is the has perm on check
func (pc permComplete) HasPermStringOn(perm string, ownerID, parentID int64, scopes ...string) (string, bool) {
	return pc.HasPermStringOn(perm, ownerID, parentID, scopes...)
}

// GetID return the id of holder
func (pc permComplete) GetID() int64 {
	return pc.id
}

// GetCurrentPerm return the current permission that this object is built on
func (pc permComplete) GetCurrentPerm() string {
	return pc.perm
}

// GetCurrentScope return the current scope for this object (maximum)
func (pc permComplete) GetCurrentScope() string {
	return pc.scope
}

// NewPermInterfaceComplete return a new object base on the minimum object
func NewPermInterfaceComplete(inner PermInterface, id int64, perm, scope string) PermInterfaceComplete {
	s, ok := inner.HasPermString(scope, perm)
	assert.True(ok, "[BUG] probably there is some thing wrong with code generation")
	pc := &permComplete{
		inner: inner,
		id:    id,
		perm:  perm,
		scope: s,
	}

	return pc
}
