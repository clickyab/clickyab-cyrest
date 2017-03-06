package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"text/template"
	"time"

	"github.com/labstack/gommon/log"
)

var (
	input           = flag.String("i", ".", "Refer to git repository in your drive")
	limit           = flag.Int("l", 100, "Specify how many log should be exported")
	redmindEndpoint = flag.String("re", "", `redmindEndpoint = flag.String("re","",Redmine host (ex: http://redmine.example.com)`)
	redmindAPIKey   = flag.String("ra", "", `Redmine APIKey. you need to enable REST API. you can find more information about how
	     to enable it on http://www.redmine.org/projects/redmine/wiki/Rest_api#Authentication`)
)

func main() {

	flag.Parse()
	if *limit < 1 {
		log.Fatal()
	}

	//config()
	c := commits()
	red := make([]trackerData, 0)
	r, e := redmineIssue()
	if e == nil {
		red = r
	}

	authors := make([]author, 0)
	for _, a := range authorList {
		authors = append(authors, a)
	}
	report := report{
		time.Now(),
		len(c),
		authors,
		c,
		red,
	}

	j, e := json.Marshal(report)

	fError(e, "way!!!!")
	b := bytes.Buffer{}
	template.JSEscape(&b, j)
	fmt.Println(string(templateBuilder(b.Bytes())))
}

type templateInfo struct {
	Date  time.Time
	Data  string
	Js    string
	Style string
}

func templateBuilder(data []byte) []byte {
	master, _ := Asset("src/cmd/reporter/resource/template/report.html")
	app, _ := Asset("src/cmd/reporter/resource/template/app.js")
	style, _ := Asset("src/cmd/reporter/resource/template/style.css")

	t := templateInfo{
		time.Now(),
		string(data),
		string(app),
		string(style),
	}

	result, e := template.New("report").Parse(string(master))

	if e != nil {
		fmt.Println(e)
	}

	buf := bytes.Buffer{}

	result.Execute(&buf, t)
	return buf.Bytes()

}

type report struct {
	CreationTime time.Time     `json:"time"`
	Count        int           `json:"count"`
	Authors      []author      `json:"authors"`
	Commits      []commitInfo  `json:"commits"`
	Redmine      []trackerData `json:"redmine"`
}
