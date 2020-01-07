package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gobuffalo/packr/v2"
)

// project project config
type project struct {
	// project name
	Name string
	// mod prefix
	ModPrefix string
	// project dir
	path     string
	none     bool
	onlyGRPC bool
	onlyHTTP bool
}

var p project

//go:generate packr2
func create() (err error) {
	box := packr.New("all", "./templates/all")
	if p.onlyHTTP {
		box = packr.New("http", "./templates/http")
	} else if p.onlyGRPC {
		box = packr.New("grpc", "./templates/grpc")
	}
	if err = os.MkdirAll(p.path, 0755); err != nil {
		return
	}
	for _, name := range box.List() {
		if p.ModPrefix != "" && name == "go.mod.tmpl" {
			continue
		}
		tmpl, _ := box.FindString(name)
		i := strings.LastIndex(name, string(os.PathSeparator))
		if i > 0 {
			dir := name[:i]
			if err = os.MkdirAll(filepath.Join(p.path, dir), 0755); err != nil {
				return
			}
		}
		if strings.HasSuffix(name, ".tmpl") {
			name = strings.TrimSuffix(name, ".tmpl")
		}
		if err = write(filepath.Join(p.path, name), tmpl); err != nil {
			return
		}
	}

	if err = generate("./..."); err != nil {
		return
	}
	if err = generate("./internal/dao/wire.go"); err != nil {
		return
	}
	return
}

func generate(path string) error {
	cmd := exec.Command("go", "generate", "-x", path)
	cmd.Dir = p.path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func write(path, tpl string) (err error) {
	data, err := parse(tpl)
	if err != nil {
		return
	}
	return ioutil.WriteFile(path, data, 0644)
}

func parse(s string) ([]byte, error) {
	t, err := template.New("").Parse(s)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, p); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
