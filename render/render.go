package render

import (
	"context"
	"github.com/dixonwille/skywalker"
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
	r.mu.Lock()
	defer r.mu.Unlock()
	for k, v := range fileSet {
		r.fileSet[k] = v
	}
}

func (r *Renderer) AddFile(filePath string, tmpl Template) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fileSet[filePath] = tmpl
}

func (r *Renderer) GetFile(filePath string) Template {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.fileSet[filePath]
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
		r.AddFile(cleansedPath, NewTemplate(string(bits)))
		return nil
	}
}

type walkFunc func(path string)

func (w walkFunc) Work(path string) {
	w(path)
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
		pathRewrites := map[string]string{
			tmpdir: dest,
		}
		walker := skywalker.New(tmpdir, walkFunc(func(path string) {
			bits, err := ioutil.ReadFile(path)
			if err != nil {
				return
			}
			var cleansedPath = path
			if len(pathRewrites) > 0 {
				for old, newPath := range pathRewrites {
					cleansedPath = strings.ReplaceAll(cleansedPath, old, newPath)
				}
			}
			r.AddFile(cleansedPath, NewTemplate(string(bits)))
		}))

		if err := walker.Walk(); err != nil {
			return err
		}
	}
	return nil
}

func (r *Renderer) Compile(data interface{}) error {
	wg := &sync.WaitGroup{}
	var errs []error
	for path, tmpl := range r.fileSet {
		wg.Add(1)
		go func(filePath string, content Template) {
			defer wg.Done()
			dirPath := filepath.Dir(filePath)
			if dirPath != "." {
				os.MkdirAll(dirPath, os.ModePerm)
			}
			file, err := os.Create(filePath)
			if err != nil {
				errs = append(errs, err)
				return
			}
			defer file.Close()
			if err := content.Compile(file, data); err != nil {
				errs = append(errs, err)
				return
			}
		}(path, tmpl)
	}
	wg.Wait()
	if len(errs) > 0 {
		var err = errors.New("compilation error!")
		for _, e := range errs {
			err = errors.Wrap(err, e.Error())
		}
		return err
	}
	return nil
}

func (r *Renderer) FileSet() map[string]Template {
	return r.fileSet
}
