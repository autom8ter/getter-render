package render_test

import (
	"context"
	"github.com/autom8ter/getter-render/render"
	"testing"
)

func TestRenderer(t *testing.T) {
	renderer := render.NewRenderer()
	if err := renderer.LoadSources(context.Background(), map[string]string{
		"tmp2": "git@github.com:autom8ter/getter-render.git/LICENSE",
	}); err != nil {
		t.Fatal(err.Error())
	}
	for name, _ := range renderer.FileSet() {
		t.Logf("filename: %s\n", name)
	}
}
