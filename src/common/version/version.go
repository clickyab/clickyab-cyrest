package version

import (
	"strconv"
	"time"

	"github.com/Sirupsen/logrus"
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
	c.Date, _ = time.Parse("01-02-06-15-04-05", date)
	c.Hash = hash
	c.Short = short
	c.BuildDate, _ = time.Parse("01-02-06-15-04-05", build)

	return c
}

// LogVersion return an logrus entry with version information attached
func LogVersion() *logrus.Entry {
	ver := GetVersion()
	return logrus.WithFields(
		logrus.Fields{
			"Commit hash":       ver.Hash,
			"Commit short hash": ver.Short,
			"Commit date":       ver.Date.Format(time.RFC3339),
			"Build date":        ver.BuildDate.Format(time.RFC3339),
		},
	)
}
