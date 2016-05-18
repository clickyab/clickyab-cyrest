package templates

import (
	"bytes"
	"html/template"
)

// RenderMail try to render a mail template from the template
func RenderMail(tpl string, context interface{}) (string, error) {
	data := MustAsset("mail/" + tpl + ".tmpl.html")
	t := template.Must(template.New(tpl).Parse(string(data)))
	b := &bytes.Buffer{}
	err := t.Execute(b, context)
	if err != nil {
		return "", err
	}

	return b.String(), nil
}
