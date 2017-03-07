package main

import (
	"fmt"

	"github.com/fzerorubigd/go-redmine"
)

var redClient = redmine.NewClient("http://tracker.clickyab.com", "062d796f163a9c7b0490c281297498b4247ba627")

var (
	endpoint, apikey string
)

func redmineIssue() ([]redmine.Issue, error) {
	//endpoint, e := redmineEndpoint()
	//apikey, er := redmineApikey()
	//if e != nil || er != nil {
	//	return nil, e
	//}
	//redClient := redmine.NewClient(endpoint, apikey)
	s, e := redClient.FilterIssues(redmine.IssueFilter{"created_on", findStartDate()})
	if e != nil {
		return nil, e
	}
	return s, nil
}

func findStartDate() string {
	//redClient := redmine.NewClient(endpoint, apikey)
	s, e := redClient.Issue(minIssueCode)
	if e != nil {
		return ""
	}
	return ">=" + s.CreatedOn
}
