package main

import (
	"github.com/fzerorubigd/go-redmine"

	"encoding/json"
	"fmt"
	"time"
)

func main() {

	config()
	c := commits()
	red := make([]redmine.Issue, 0)
	r, e := redmineIssue()
	if e != nil {
		red = append(red, r...)
	}

	authors := make([]Author, 0)
	for _, a := range AuthorList {
		authors = append(authors, a)
	}
	report := Report{
		time.Now(),
		len(c),
		authors,
		c,
		red,
	}

	j, _ := json.Marshal(report)
	fmt.Println(string(j))
}

type Report struct {
	CreationTime time.Time       `json:"time"`
	Count        int             `json:"count"`
	Authors      []Author        `json:"authors"`
	Commits      []CommitInfo    `json:"commits"`
	Redmine      []redmine.Issue `json:"redmine"`
}
