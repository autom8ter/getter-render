package render_test

import (
	"bytes"
	"github.com/autom8ter/getter-tmpl/render"
	"testing"
)

var listObj = &List{Elements: []string{"http://localhost:3000", "http://localhost:3000/callback"}}

type List struct {
	Elements []string
}

func TestTemplate(t *testing.T) {
	list := render.NewTemplate(`list = [{{range $i, $v := .Elements}}{{if $i}}, {{end}}"{{.}}"{{end}}]`)
	buffer := bytes.NewBuffer(nil)
	if err := list.Compile(buffer, listObj); err != nil {
		t.Fatal(err.Error())
	}
	newList := render.NewTemplate(buffer.String())
	t.Logf("compiled: %s\n", newList.String())

	fileSet := render.NewFileSet(map[string]render.Template{
		"tmp/list.txt": list,
	})

	if err := fileSet.Compile(listObj); err != nil {
		t.Fatal(err.Error())
	}
}
