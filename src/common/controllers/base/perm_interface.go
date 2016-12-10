package base

// PermInterface is the perm interface
type PermInterface interface {
	// HasPermString is the has perm check
	HasPermString(scope string, perm string) (string, bool)
	// HasPermStringOn is the has perm on check
	HasPermStringOn(perm string, ownerID, parentID int64, scopes ...string) (string, bool)
}
