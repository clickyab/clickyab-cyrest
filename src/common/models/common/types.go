package common

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const sqlNull = "null"

// Initializer is for model when the have need extra initialize on save/update
type Initializer interface {
	// Initialize is the method to call att save/update
	Initialize()
}

// Int64Slice is simple slice to handle it for json field
type Int64Slice []int64

// Int64Array is used to handle real array in database
type Int64Array []int64

// NullTime is null-time for json in null
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// ZeroNullInt64 is ZeroNullInt64 zero if null
type ZeroNullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if num is not NULL
}

// NullInt64 is null int64 for json in null
type NullInt64 struct {
	Int64 int64
	Valid bool // Valid is true if Int64 is not NULL
}

// NullString is the json friendly of sql.NullString
type NullString struct {
	Valid  bool
	String string
}

// MB4String is the emoji ready string
type MB4String []byte

// GenericJSONField is used to handle generic json data in postgres
type GenericJSONField map[string]interface{}

// StringJSONArray is use to handle string to string map in postgres
type StringJSONArray map[string]string

// Scan convert the json array ino string slice
func (is *Int64Slice) Scan(src interface{}) error {
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

	return json.Unmarshal(b, is)
}

// Value try to get the string slice representation in database
func (is Int64Array) Value() (driver.Value, error) {
	b, err := json.Marshal(is)
	if err != nil {
		return nil, err
	}
	// Its time to change [] to {}
	b = bytes.Replace(b, []byte("["), []byte("{"), 1)
	b = bytes.Replace(b, []byte("]"), []byte("}"), 1)

	return b, nil
}

// Scan convert the json array ino string slice
func (is *Int64Array) Scan(src interface{}) error {
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
	b = bytes.Replace(b, []byte("{"), []byte("["), 1)
	b = bytes.Replace(b, []byte("}"), []byte("]"), 1)

	return json.Unmarshal(b, is)
}

// Value try to get the string slice representation in database
func (is Int64Slice) Value() (driver.Value, error) {
	return json.Marshal(is)
}

// MarshalJSON try to marshaling to json
func (nt NullTime) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return nt.Time.MarshalJSON()
	}

	return []byte(sqlNull), nil
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// UnmarshalJSON try to unmarshal dae from input
func (nt *NullTime) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == sqlNull {
		nt.Valid = false
		nt.Time = time.Time{}
		return nil
	}

	err := json.Unmarshal(b, &nt.Time)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// MarshalJSON try to marshaling to json
func (nt NullInt64) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(fmt.Sprintf(`%d`, nt.Int64)), nil
	}

	return []byte(sqlNull), nil
}

// UnmarshalJSON try to unmarshal dae from input
func (nt *NullInt64) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == sqlNull {
		nt.Valid = false

		return nil
	}

	err := json.Unmarshal(b, &nt.Int64)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// Scan implements the Scanner interface.
func (nt *NullInt64) Scan(value interface{}) error {
	inn := &sql.NullInt64{}
	err := inn.Scan(value)
	if err != nil {
		return err
	}
	nt.Int64 = inn.Int64
	nt.Valid = inn.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (nt NullInt64) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Int64, nil
}

// MarshalJSON try to marshaling to json
func (nt ZeroNullInt64) MarshalJSON() ([]byte, error) {
	if nt.Valid {
		return []byte(fmt.Sprint(nt.Int64)), nil
	}

	return []byte("0"), nil
}

// UnmarshalJSON try to unmarshal dae from input
func (nt *ZeroNullInt64) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == "0" {
		nt.Valid = false

		return nil
	}

	err := json.Unmarshal(b, &nt.Int64)
	if err != nil {
		return err
	}

	nt.Valid = true
	return nil
}

// Scan implements the Scanner interface.
func (nt *ZeroNullInt64) Scan(value interface{}) error {
	inn := &sql.NullInt64{}
	err := inn.Scan(value)
	if err != nil {
		return err
	}
	nt.Int64 = inn.Int64
	nt.Valid = inn.Valid
	return nil
}

// Value implements the driver Valuer interface.
func (nt ZeroNullInt64) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Int64, nil
}

// Scan convert the json array ino string slice
func (gjf *GenericJSONField) Scan(src interface{}) error {
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

	return json.Unmarshal(b, gjf)
}

// Value try to get the string slice representation in database
func (gjf GenericJSONField) Value() (driver.Value, error) {
	return json.Marshal(gjf)
}

// Scan convert the json array ino string slice
func (ms *MB4String) Scan(src interface{}) error {
	var b []byte
	var err error
	switch src.(type) {
	case []byte:
		b = src.([]byte)
	case string:
		b = []byte(src.(string))
	case nil:
		*ms = make([]byte, 0)
		return nil
	default:
		return errors.New("unsupported type")
	}
	*ms, err = base64.StdEncoding.DecodeString(string(b))
	return err
}

// Value try to get the string slice representation in database
func (ms MB4String) Value() (driver.Value, error) {
	if ms == nil || len(ms) == 0 {
		return nil, nil
	}
	tmp := base64.StdEncoding.EncodeToString(ms)
	return []byte(tmp), nil
}

// MarshalJSON try to marshaling to json
func (ms *MB4String) MarshalJSON() ([]byte, error) {
	return json.Marshal([]byte(*ms))
}

// UnmarshalJSON try to unmarshal dae from input
func (ms *MB4String) UnmarshalJSON(b []byte) error {
	tmp := []byte{}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	*ms = tmp
	return nil
}

// Scan convert the json array ino string slice
func (ss *StringJSONArray) Scan(src interface{}) error {
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

	return json.Unmarshal(b, ss)
}

// Value try to get the string slice representation in database
func (ss StringJSONArray) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	tmp := &sql.NullString{}
	err := tmp.Scan(value)
	if err != nil {
		return err
	}
	ns.Valid = tmp.Valid
	ns.String = tmp.String
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON try to marshaling to json
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}

	return []byte(sqlNull), nil
}

// UnmarshalJSON try to unmarshal dae from input
func (ns NullString) UnmarshalJSON(b []byte) error {
	text := strings.ToLower(string(b))
	if text == sqlNull {
		ns.Valid = false
		ns.String = ""
		return nil
	}

	err := json.Unmarshal(b, &ns.String)
	if err != nil {
		return err
	}

	ns.Valid = true
	return nil
}

// MakeNullString create a new null string
func MakeNullString(s string) NullString {
	return NullString{Valid: s != "", String: s}
}

// MakeNullTime create a new null time
func MakeNullTime(t time.Time) NullTime {
	return NullTime{Valid: !t.IsZero(), Time: t}
}

// MakeNullInt64 create a new null int64
func MakeNullInt64(t int64) NullInt64 {
	return NullInt64{Valid: t != 0, Int64: t}
}
