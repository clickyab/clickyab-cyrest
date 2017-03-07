package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"encoding/hex"

	"time"

	"github.com/okian/go-git"
)

const (
	signPattern = "(?i) -----END PGP SIGNATURE-----"
)

var (
	minIssueCode int = -1
	signExp          = regexp.MustCompile(signPattern)
	AuthorList       = make(map[string]Author)
)

type Author struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type CommitInfo struct {
	Date     time.Time `json:"date"`
	Hash     string    `json:"hash"`
	Message  string    `json:"message"`
	Tags     []Tag     `json:"tags"`
	AuthorId int       `json:"author_id"`
}

type Tag struct {
	Code int    `json:"c"`
	Type string `json:"t"`
}

func commits() []CommitInfo {
	repo, err := git.PlainOpen(getPath())
	fError(err, fmt.Sprintf("No git repository in %s!", getPath()))
	head, _ := repo.Head()
	commit, _ := repo.Commit(head.Hash())
	h, _ := commit.History()

	commits := make([]CommitInfo, 0)
	var lmt int
	if len(h) < limit() {
		lmt = len(h)
	} else {
		lmt = limit()
	}
	for i := 0; i < lmt; i++ {
		c := h[i]
		commits = append(commits, CommitInfo{
			c.Author.When,
			c.Hash.String(),
			pureMessage(c.Message),
			TagFinder(c.Message),
			authorId(Author{
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

func authorId(a Author) int {
	if i, exist := AuthorList[a.Hash]; exist == true {
		return i.Id
	}
	a.Id = len(AuthorList) + 1
	AuthorList[a.Hash] = a
	return a.Id
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
	io.WriteString(h, e)
	return hex.EncodeToString(h.Sum(nil))
}

func TagFinder(m string) []Tag {
	m = strings.ToLower(m)

	rawTags := tagPattern.FindAllStringSubmatch(m, -1)
	result := make([]Tag, 0)
	if len(rawTags) > 0 {
		for _, c := range rawTags {
			code, _ := strconv.Atoi(c[2])
			isOldestIssue(code)
			result = append(result, Tag{code, c[1]})
		}
	}
	return result
}
