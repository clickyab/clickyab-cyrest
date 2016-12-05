package aaa

import (
	"common/utils"
	"database/sql/driver"
	"errors"
)

// AUTO GENERATED CODE. DO NOT EDIT!

// IsValid try to validate enum value on ths type
func (e ScopePerm) IsValid() bool {
	return utils.StringInArray(
		string(e),
		string(ScopePermGlobal),
		string(ScopePermParent),
		string(ScopePermOwn),
	)
}

// Scan convert the json array ino string slice
func (e *ScopePerm) Scan(src interface{}) error {
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
	if !ScopePerm(b).IsValid() {
		return errors.New("invaid value")
	}
	*e = ScopePerm(b)
	return nil
}

// Value try to get the string slice representation in database
func (e ScopePerm) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("invaid status")
	}
	return string(e), nil
}

// IsValid try to validate enum value on ths type
func (e UserStatus) IsValid() bool {
	return utils.StringInArray(
		string(e),
		string(UserStatusRegistered),
		string(UserStatusVerified),
		string(UserStatusBlocked),
	)
}

// Scan convert the json array ino string slice
func (e *UserStatus) Scan(src interface{}) error {
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
		return errors.New("invaid value")
	}
	*e = UserStatus(b)
	return nil
}

// Value try to get the string slice representation in database
func (e UserStatus) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("invaid status")
	}
	return string(e), nil
}

// IsValid try to validate enum value on ths type
func (e UserSource) IsValid() bool {
	return utils.StringInArray(
		string(e),
		string(UserSourceCRM),
		string(UserSourceClickyab),
	)
}

// Scan convert the json array ino string slice
func (e *UserSource) Scan(src interface{}) error {
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
	if !UserSource(b).IsValid() {
		return errors.New("invaid value")
	}
	*e = UserSource(b)
	return nil
}

// Value try to get the string slice representation in database
func (e UserSource) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("invaid status")
	}
	return string(e), nil
}

// IsValid try to validate enum value on ths type
func (e UserType) IsValid() bool {
	return utils.StringInArray(
		string(e),
		string(UserTypePersonal),
		string(UserTypeCorporation),
	)
}

// Scan convert the json array ino string slice
func (e *UserType) Scan(src interface{}) error {
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
	if !UserType(b).IsValid() {
		return errors.New("invaid value")
	}
	*e = UserType(b)
	return nil
}

// Value try to get the string slice representation in database
func (e UserType) Value() (driver.Value, error) {
	if !e.IsValid() {
		return nil, errors.New("invaid status")
	}
	return string(e), nil
}
