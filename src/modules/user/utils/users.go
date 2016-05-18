package utils

import (
	"fmt"
	"modules/misc/trans"
	"regexp"
)

// ContactType is the type of contact
type ContactType int

const (
	// TypeEmail is the email type in contact
	TypeEmail ContactType = iota
	// TypePhone is the phone type in the system
	TypePhone
)

var (
	phone = regexp.MustCompile("^[0]{0,1}9[0-9]{9}$")
	email = regexp.MustCompile(`(?i)^(([a-zA-Z]|[0-9])|([-]|[_]|[.]))+[@](([a-zA-Z0-9])|([-])){2,63}[.](([a-zA-Z0-9]){2,63})+$`)
)

// DetectContactType is the contact type detection
func DetectContactType(c string) (ContactType, error) {
	if phone.MatchString(c) {
		return TypePhone, nil
	}

	if email.MatchString(c) {
		return TypeEmail, nil
	}

	return 0, fmt.Errorf(trans.T("not a phone number, nor an email address"))
}
