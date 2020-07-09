package render_test

import (
	"context"
	"github.com/autom8ter/getter-render/render"
	"github.com/spf13/viper"
	"testing"
)

func TestRenderer(t *testing.T) {
	renderer := render.NewRenderer()
	if err := renderer.LoadSources(context.Background(), map[string]string{
		"tmp": "https://raw.githubusercontent.com/autom8ter/getter-render/master/test.txt",
	}); err != nil {
		t.Fatal(err.Error())
	}
	for name, _ := range renderer.FileSet() {
		t.Logf("filename: %s\n", name)
	}
	viper.Set("name", "Coleman Word")
	if err := renderer.Compile(viper.AllSettings()); err != nil {
		t.Fatal(err.Error())
	}
}
