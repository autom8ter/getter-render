package render

import (
	"context"
	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
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

type FileSet map[string]Template

func NewFileSet(files ...map[string]Template) FileSet {
	fileset := map[string]Template{}
	for _, set := range files {
		for k, v := range set {
			fileset[k] = v
		}
	}
	return fileset
}

func (f FileSet) Compile(data interface{}) error {
	for filePath, content := range f {
		dirPath := filepath.Dir(filePath)
		if dirPath != "." {
			os.MkdirAll(dirPath, os.ModePerm)
		}
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close()
		if err := content.Compile(file, data); err != nil {
			return err
		}
	}
	return nil
}

func (f FileSet) LoadFunc() filepath.WalkFunc {
	var mu = &sync.Mutex{}
	return func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		bits, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		mu.Lock()
		f[path] = NewTemplate(string(bits))
		mu.Unlock()
		return nil
	}
}

func (f FileSet) Load(ctx context.Context, sources []string) error {
	for _, source := range sources {
		tmpdir, err := ioutil.TempDir("", "")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpdir)
		var mode = getter.ClientModeAny
		// Get the pwd
		pwd, err := os.Getwd()
		if err != nil {
			return errors.Wrap(err, "failed to os.Getwd() before loading files")
		}
		opts := []getter.ClientOption{}

		// Build the client
		client := &getter.Client{
			Ctx:     ctx,
			Src:     source,
			Dst:     tmpdir,
			Pwd:     pwd,
			Mode:    mode,
			Options: opts,
		}
		if err := client.Get(); err != nil {
			return errors.Wrapf(err, "failed to load files. source: %s", source)
		}
		if err := filepath.Walk(tmpdir, f.LoadFunc()); err != nil {
			return err
		}
	}
	return nil
}
