package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"text/template"
)

// project project config
type project struct {
	Name     string
	Owner    string
	Path     string
	WithGRPC bool
	Here     bool
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
	_tplTypeAPIGenerateWin
	_tplTypeGomod
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
		_tplTypeGomod:      "/go.mod",
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
		_tplTypeDao:            _tplDao,
		_tplTypeHTTPServer:     _tplHTTPServer,
		_tplTypeAPIProto:       _tplAPIProto,
		_tplTypeAPIGenerate:    _tplAPIGenerate,
		_tplTypeAPIGenerateWin: _tplAPIGenerateWin,
		_tplTypeMain:           _tplMain,
		_tplTypeChangeLog:      _tplChangeLog,
		_tplTypeContributors:   _tplContributors,
		_tplTypeReadme:         _tplReadme,
		_tplTypeMySQLToml:      _tplMySQLToml,
		_tplTypeMCToml:         _tplMCToml,
		_tplTypeRedisToml:      _tplRedisToml,
		_tplTypeAppToml:        _tplAppToml,
		_tplTypeHTTPToml:       _tplHTTPToml,
		_tplTypeModel:          _tplModel,
		_tplTypeGomod:          _tplGoMod,
	}
)

func create() (err error) {
	if p.WithGRPC {
		files[_tplTypeGRPCServer] = "/internal/server/grpc/server.go"
		files[_tplTypeAPIProto] = "/api/api.proto"
		files[_tplTypeAPIGenerate] = "/api/gen.go"
		files[_tplTypeAPIGenerateWin] = "/api/gen_windows.go"
		tpls[_tplTypeGRPCServer] = _tplGRPCServer
		tpls[_tplTypeGRPCToml] = _tplGRPCToml
		tpls[_tplTypeService] = _tplGPRCService
	} else {
		tpls[_tplTypeService] = _tplService
		tpls[_tplTypeMain] = delGRPC(_tplMain)
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
		if err = checkProtobuf(); err != nil {
			return
		}
		if err = genGRPC(); err != nil {
			return
		}
	}
	return
}

func checkProtobuf() error {
	if _, err := exec.LookPath("protoc"); err != nil {
		switch runtime.GOOS {
		case "darwin":
			// MacOS
			fmt.Println("brew install protobuf..")
			cmd := exec.Command("brew", "install", "protobuf")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err = cmd.Run(); err != nil {
				return err
			}
		default:
			return errors.New("生成 gRPC API 失败，您还没安装protobuf，请进行手动安装：https://github.com/protocolbuffers/protobuf/releases")
		}
	}
	return nil
}

func genGRPC() error {
	dir := path.Join(p.Path, "api")
	cmd := exec.Command("go", "generate", dir)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("generate gRPC error(%v)", err)
	}
	return nil
}

func delGRPC(tpl string) string {
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
