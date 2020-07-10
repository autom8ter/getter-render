package render

import (
	"context"
	"github.com/hashicorp/go-getter"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

func (r *Renderer) AddFiles(fileSet map[string]Template) {
	for k, v := range fileSet {
		r.fileSet[k] = v
	}
}

func (r *Renderer) LoadFunc(pathRewrites map[string]string) filepath.WalkFunc {
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
		var cleansedPath = path
		if len(pathRewrites) > 0 {
			for old, newPath := range pathRewrites {
				cleansedPath = strings.ReplaceAll(cleansedPath, old, newPath)
			}
		}
		r.mu.Lock()
		r.fileSet[cleansedPath] = NewTemplate(string(bits))
		r.mu.Unlock()
		return nil
	}
}

func (r *Renderer) LoadSources(ctx context.Context, destSource map[string]string) error {
	if len(destSource) == 0 {
		return errors.New("empty source mapping")
	}
	for dest, source := range destSource {
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
		if err := filepath.Walk(tmpdir, r.LoadFunc(map[string]string{
			tmpdir: dest,
		})); err != nil {
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

func (r *Renderer) FileSet() map[string]Template {
	return r.fileSet
}
