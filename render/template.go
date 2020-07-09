package render

import (
	"io"
	"text/template"
)

type Template string

func NewTemplate(tmpl string) Template {
	return Template(tmpl)
}

func (t Template) String() string {
	return string(t)
}

func (t Template) Compile(w io.Writer, data interface{}) error {
	tmpl, err := template.New("").Funcs(funcMap()).Parse(t.String())
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}
