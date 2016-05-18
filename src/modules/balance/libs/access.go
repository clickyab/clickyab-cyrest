package libs

// Access is the structure represent the access model
type Access interface {
	// IsOwner return if the current user is the owner of this unit or not
	IsOwner() bool
	// IsAccessToAccount return if the current user has access to this account or not
	// the first result means its visible, the last means its writable
	IsAccessToObject(int64) (bool, bool)
	// VisibleAccountIDs return list of all visible account for this user, the first bool means all
	VisibleObjectIDs() (bool, []int64)
	// WritableAccountIDs return the account ids that this access can write into, the first result means all
	WritableObjectIDs() (bool, []int64)
}
