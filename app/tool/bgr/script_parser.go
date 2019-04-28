package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

var (
	_lints = make([]*lint, 0)
)

func walkScript(dir string) {
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			_log.Debugf("%+v", err)
			return err
		}
		if !strings.HasSuffix(info.Name(), ".bgl") {
			return nil
		}
		var (
			file    *os.File
			newErr  error
			scripts []*script
		)
		if file, newErr = os.Open(path); newErr != nil {
			newErr = errors.WithStack(newErr)
			return newErr
		}
		if scripts, newErr = fileToScript(file, path); newErr != nil {
			newErr = errors.WithStack(newErr)
			return newErr
		}
		for _, s := range scripts {
			registerLints(s)
		}
		return nil
	}
	if err := filepath.Walk(dir, fn); err != nil {
		panic(err)
	}
}

func fileToScript(file *os.File, path string) (scripts []*script, err error) {
	var (
		br        = bufio.NewReader(file)
		curScript *script
		line      []byte
		isPrefix  bool
	)
	for line, isPrefix, err = br.ReadLine(); err != io.EOF; line, isPrefix, err = br.ReadLine() {
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		if isPrefix {
			_log.Fatalf("parseScript file: %s/%s err: some line too long", path, file.Name())
		}

		strs := strings.Split(strings.TrimSpace(string(line)), " ")
		if len(strs) != 2 {
			continue
		}
		k, v := strs[0], strs[1]

		switch k {
		case "T":
			ts := strings.Split(strings.TrimSpace(v), ".")
			curScript = &script{
				dir: filepath.Dir(path),
				ts:  ts,
				l:   "e",
				d:   fmt.Sprintf("{%s : %s}", strings.Join(ts, "."), v),
			}
		case "V":
			if curScript != nil {
				curScript.v = v
			}
			scripts = append(scripts, curScript)
		case "L":
			if curScript != nil {
				curScript.l = v
			}
		case "D":
			if curScript != nil {
				curScript.d = v
			}
		}
	}
	err = nil
	return
}

func registerLints(script *script) {
	_lints = append(_lints, &lint{
		s:  script,
		fn: assembleLint(script),
	})
}

func assembleLint(script *script) func(curDir string, f *ast.File, node ast.Node) bool {
	var (
		reg *regexp.Regexp
		err error
	)
	_log.Debugf("assembleLint script: %+v", script)
	if reg, err = regexp.Compile(script.v); err != nil {
		_log.Fatalf("assembleLint script: %s, v compile error: %+v", script, err)
		return nil
	}
	return func(curDir string, f *ast.File, n ast.Node) bool {
		// if !strings.HasPrefix(curDir, script.dir) {
		// 	return true
		// }
		var (
			parse func(curDir string, f *ast.File, n ast.Node) (v string, hit bool)
			ok    bool
			k     = strings.Join(script.ts, ".")
		)
		if parse, ok = _parsers[k]; !ok {
			return true
		}
		content, hit := parse(curDir, f, n)
		if !hit {
			return true
		}
		return !reg.MatchString(content) // 返回是否是正常node（未命中lint）
	}
}
