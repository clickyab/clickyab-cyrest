package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"common/assert"

	"github.com/okian/go-git"
)

const (
	signPattern = "(?i) -----END PGP SIGNATURE-----"
)

var (
	minIssueCode = -1
	signExp      = regexp.MustCompile(signPattern)
	authorList   = make(map[string]author)
)

type author struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type commitInfo struct {
	Date     time.Time `json:"time"`
	Hash     string    `json:"hash"`
	Message  string    `json:"message"`
	Tags     []tag     `json:"tags"`
	AuthorID int       `json:"author_id"`
}

type tag struct {
	Code int    `json:"c"`
	Type string `json:"t"`
}

func commits(repo *git.Repository) []commitInfo {
	head, _ := repo.Head()
	commit, _ := repo.Commit(head.Hash())

	h, _ := commit.History()

	commits := make([]commitInfo, 0)
	var lmt int
	if len(h) < *limit {
		lmt = len(h)
	} else {
		lmt = *limit
	}
	for i := 0; i < lmt; i++ {
		c := h[i]
		commits = append(commits, commitInfo{
			c.Author.When,
			c.Hash.String(),
			pureMessage(c.Message),
			tagFinder(c.Message),
			authorID(author{
				0,
				c.Author.Name,
				c.Author.Email,
				emailHash(c.Author.Email),
			}),
		})
	}

	return commits
}

func isOldestIssue(code int) {
	if minIssueCode == -1 || minIssueCode > code {
		minIssueCode = code
	}
}

func authorID(a author) int {
	if i, exist := authorList[a.Hash]; exist {
		return i.ID
	}
	a.ID = len(authorList) + 1
	authorList[a.Hash] = a
	return a.ID
}

func pureMessage(m string) string {
	if signExp.Match([]byte(m)) {
		i := signExp.FindStringIndex(m)
		m = m[i[0]+len(signPattern)+2:]
	}
	return m
}

func emailHash(e string) string {
	h := md5.New()
	_, err := io.WriteString(h, e)
	assert.Nil(err)
	return hex.EncodeToString(h.Sum(nil))
}

func tagFinder(m string) []tag {
	m = strings.ToLower(m)

	rawTags := tagPattern.FindAllStringSubmatch(m, -1)
	result := make([]tag, 0)
	if len(rawTags) > 0 {
		for _, c := range rawTags {
			code, _ := strconv.Atoi(c[2])
			isOldestIssue(code)
			result = append(result, tag{code, c[1]})
		}
	}
	return result
}
