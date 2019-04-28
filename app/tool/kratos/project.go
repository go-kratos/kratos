package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// project project config
type project struct {
	Department string
	Type       string
	Name       string
	Owner      string
	Path       string
	WithGRPC   bool
}

const (
	_tplTypeDao = iota
	_tplTypeHTTPServer
	_tplTypeAPIProto
	_tplTypeService
	_tplTypeMain
	_tplTypeChangeLog
	_tplTypeContributors
	_tplTypeReadme
	_tplTypeAppToml
	_tplTypeMySQLToml
	_tplTypeMCToml
	_tplTypeRedisToml
	_tplTypeHTTPToml
	_tplTypeGRPCToml
	_tplTypeModel
	_tplTypeGRPCServer
	_tplTypeAPIGenerate
)

var (
	p project
	// files type => path
	files = map[int]string{
		// init doc
		_tplTypeChangeLog:    "/CHANGELOG.md",
		_tplTypeContributors: "/CONTRIBUTORS.md",
		_tplTypeReadme:       "/README.md",
		// init project
		_tplTypeMain:       "/cmd/main.go",
		_tplTypeDao:        "/internal/dao/dao.go",
		_tplTypeHTTPServer: "/internal/server/http/http.go",
		_tplTypeService:    "/internal/service/service.go",
		_tplTypeModel:      "/internal/model/model.go",
		// init config
		_tplTypeAppToml:   "/configs/application.toml",
		_tplTypeMySQLToml: "/configs/mysql.toml",
		_tplTypeMCToml:    "/configs/memcache.toml",
		_tplTypeRedisToml: "/configs/redis.toml",
		_tplTypeHTTPToml:  "/configs/http.toml",
		_tplTypeGRPCToml:  "/configs/grpc.toml",
	}
	// tpls type => content
	tpls = map[int]string{
		_tplTypeDao:          _tplDao,
		_tplTypeHTTPServer:   _tplHTTPServer,
		_tplTypeAPIProto:     _tplAPIProto,
		_tplTypeAPIGenerate:  _tplAPIGenerate,
		_tplTypeMain:         _tplMain,
		_tplTypeChangeLog:    _tplChangeLog,
		_tplTypeContributors: _tplContributors,
		_tplTypeReadme:       _tplReadme,
		_tplTypeMySQLToml:    _tplMySQLToml,
		_tplTypeMCToml:       _tplMCToml,
		_tplTypeRedisToml:    _tplRedisToml,
		_tplTypeAppToml:      _tplAppToml,
		_tplTypeHTTPToml:     _tplHTTPToml,
		_tplTypeModel:        _tplModel,
	}
)

func validate() (ok bool) {
	if !depts[p.Department] {
		println("[-d] Invalid department name.")
		return
	}
	if !types[p.Type] {
		println("[-t] Invalid project type.")
		return
	}
	if p.Name == "" {
		println("[-n] Invalid project name.")
		return
	}
	// TODO: check project name
	return true
}

func create() (err error) {
	if p.WithGRPC {
		files[_tplTypeGRPCServer] = "/internal/server/grpc/server.go"
		files[_tplTypeAPIProto] = "/api/api.proto"
		files[_tplTypeAPIGenerate] = "/api/generate.go"
		tpls[_tplTypeGRPCServer] = _tplGRPCServer
		tpls[_tplTypeGRPCToml] = _tplGRPCToml
		tpls[_tplTypeService] = _tplGPRCService
	} else {
		tpls[_tplTypeService] = _tplService
		tpls[_tplTypeMain] = delgrpc(_tplMain)
	}
	if err = os.MkdirAll(p.Path, 0755); err != nil {
		return
	}
	for t, v := range files {
		i := strings.LastIndex(v, "/")
		if i > 0 {
			dir := v[:i]
			if err = os.MkdirAll(p.Path+dir, 0755); err != nil {
				return
			}
		}
		if err = write(p.Path+v, tpls[t]); err != nil {
			return
		}
	}
	if p.WithGRPC {
		if err = genpb(); err != nil {
			return
		}
	}
	err = updateProwAction()
	return
}

func genpb() error {
	cmd := exec.Command("go", "generate", p.Path+"/api/generate.go")
	return cmd.Run()
}

func delgrpc(tpl string) string {
	var buf bytes.Buffer
	lines := strings.Split(tpl, "\n")
	for _, l := range lines {
		if strings.Contains(l, "grpc") {
			continue
		}
		if strings.Contains(l, "warden") {
			continue
		}
		buf.WriteString(l)
		buf.WriteString("\n")
	}
	return buf.String()
}

func write(name, tpl string) (err error) {
	data, err := parse(tpl)
	if err != nil {
		return
	}
	return ioutil.WriteFile(name, data, 0644)
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
