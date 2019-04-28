package tomllint

import (
	"io"
	"regexp"
	"strconv"

	"github.com/BurntSushi/toml"

	"go-common/app/admin/main/config/pkg/lint"
)

var lineNumberRe *regexp.Regexp

const filetype = "toml"

type lintFn func(metadata toml.MetaData) []lint.LineErr

var lintFns []lintFn

type tomllint struct{}

// Lint toml file return lint.Error
func (tomllint) Lint(r io.Reader) lint.Error {
	var v interface{}
	var lintErr lint.Error
	metadata, err := toml.DecodeReader(r, &v)
	if err != nil {
		line := -1
		if match := lineNumberRe.FindStringSubmatch(err.Error()); len(match) == 2 {
			line, _ = strconv.Atoi(match[1])
		}
		lintErr = append(lintErr, lint.LineErr{Line: line, Message: err.Error()})
		return lintErr
	}
	for _, fn := range lintFns {
		if lineErrs := fn(metadata); lineErrs != nil {
			lintErr = append(lintErr, lineErrs...)
		}
	}
	if len(lintErr) == 0 {
		return nil
	}
	return lintErr
}

// not allowed defined kv that type is not Hash at top level
//func noTopKV(metadata toml.MetaData) []lint.LineErr {
//	var lineErrs []lint.LineErr
//	for _, keys := range metadata.Keys() {
//		if len(keys) != 1 {
//			continue
//		}
//		typeName := metadata.Type(keys...)
//		if typeName != "Hash" {
//			lineErrs = append(lineErrs, lint.LineErr{
//				Line:    -1,
//				Message: fmt.Sprintf("top level value must be Object, key: %s type is %s", keys[0], typeName),
//			})
//		}
//	}
//	return lineErrs
//}

// noApp not allowed app section exists
func noApp(metadata toml.MetaData) []lint.LineErr {
	if metadata.IsDefined("app") {
		return []lint.LineErr{{Line: -1, Message: "请删除无用 App 配置 see: http://git.bilibili.co/platform/go-common/issues/310 (゜-゜)つロ"}}
	}
	return nil
}

// noIdentify not allowed identify config
func noIdentify(metadata toml.MetaData) []lint.LineErr {
	if metadata.IsDefined("identify") {
		return []lint.LineErr{{Line: -1, Message: "请删除无用 Identify 配置 see: http://git.bilibili.co/platform/go-common/issues/310 (゜-゜)つロ"}}
	}
	return nil
}

// noCommon not allowed common config
func noCommon(metadata toml.MetaData) []lint.LineErr {
	count := 0
	commonKey := []string{"version", "user", "pid", "dir", "perf"}
	for _, key := range commonKey {
		if metadata.IsDefined(key) {
			count++
		}
	}
	if count > 0 {
		return []lint.LineErr{{Line: -1, Message: "请删除无用 Common 配置 see: http://git.bilibili.co/platform/go-common/issues/310 (゜-゜)つロ"}}
	}
	return nil
}

func init() {
	lint.RegisterLinter(filetype, tomllint{})
	lintFns = []lintFn{noApp, noIdentify, noCommon}
	lineNumberRe = regexp.MustCompile("^Near line (\\d+)")
}
