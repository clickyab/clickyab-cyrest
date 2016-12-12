package base

import (
	"common/utils"
	"database/sql/driver"
	"errors"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// IsValid try to validate enum value on ths type
func (e UserScope) IsValid() bool {
	return utils.StringInArray(
		string(e),
		string(ScopeSelf),
		string(ScopeParent),
		string(ScopeGlobal),
	)
}

// Scan convert the json array ino string slice
func (e *UserScope) Scan(src interface{}) error {
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
	if !UserScope(b).IsValid() {
		return errors.New("invaid value")
	}
	*e = UserScope(b)
	return nil
}

// Value try to get the string slice representation in database
func (e UserScope) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("invaid status")
	}
	return string(e), nil
}
