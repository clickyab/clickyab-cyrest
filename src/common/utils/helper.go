package utils

import (
	"os"
	"regexp"
	"strings"
)

var (
	spaceMatch = regexp.MustCompile(`\s+`)
	emailMath  = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

// PrefixMatch return the matched items in array
func PrefixMatch(in string, data ...string) []string {
	var res []string
	for i := range data {
		if len(in) < len(data[i]) {
			if data[i][:len(in)] == in {
				res = append(res, data[i][len(in):])
			}
		}
	}

	return res
}

// CleanSplit replace all multiple strings with one and then split them using the space as delimiter
func CleanSplit(line string) []string {
	str := strings.Trim(spaceMatch.ReplaceAllString(line, " "), " ")
	return strings.Split(str, " ")
}

// Exists returns whether the given file or directory exists or not
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// StringInArray check for a string in other strings
func StringInArray(q string, arr ...string) bool {
	for i := range arr {
		if arr[i] == q {
			return true
		}
	}

	return false
}

// Int64InArray check for a int in other ints
func Int64InArray(q int64, arr ...int64) bool {
	for i := range arr {
		if arr[i] == q {
			return true
		}
	}

	return false
}

// ValidateEmail try to validate email
func ValidateEmail(email string) bool {
	return emailMath.MatchString(email)
}

// DBImplode implode ids
func DBImplode(IDs []int64) string {
	return strings.Trim(strings.Repeat("?,", len(IDs)), ",")
}
