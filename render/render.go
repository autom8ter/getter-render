package render

import (
	"context"
	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type Renderer struct {
	fileSet map[string]Template
	mu      *sync.Mutex
}

func NewRenderer() *Renderer {
	return &Renderer{
		fileSet: map[string]Template{},
		mu:      &sync.Mutex{},
	}
}

func (r *Renderer) LoadFunc() filepath.WalkFunc {
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
		r.mu.Lock()
		r.fileSet[path] = NewTemplate(string(bits))
		r.mu.Unlock()
		return nil
	}
}

func (r *Renderer) LoadSources(ctx context.Context, sources []string) error {
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
		if err := filepath.Walk(tmpdir, r.LoadFunc()); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) Compile(data interface{}) error {
	for filePath, content := range r.fileSet {
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
