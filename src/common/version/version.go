package version

import (
	"strconv"
	"time"
)

// The following variables, are for compile time set, do not edit
var (
	hash  string
	short string
	date  string
	build string
	count string
)

// Version is the application version in detail
type Version struct {
	Hash      string    `json:"hash"`
	Short     string    `json:"short_hash"`
	Date      time.Time `json:"commit_date"`
	Count     int64     `json:"build_number"`
	BuildDate time.Time `json:"build_date"`
}

// GetVersion return the application version in detail
func GetVersion() Version {
	c := Version{}
	c.Count, _ = strconv.ParseInt(count, 10, 64)
	c.Date, _ = time.Parse("Mon-Jan-2-15:04:05-2006", date)
	c.Hash = hash
	c.Short = short
	c.BuildDate, _ = time.Parse("2006-01-02-15:04:05Z07:00", build)

	return c
}
