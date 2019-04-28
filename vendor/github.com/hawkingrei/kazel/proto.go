package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

type ProtoInfo struct {
	src         []string
	importPath  string
	packageName string
	imports     []string
	isGogo      bool
	hasServices bool
}

var protoRe = buildProtoRegexp()

const (
	importSubexpIndex    = 1
	packageSubexpIndex   = 2
	goPackageSubexpIndex = 3
	serviceSubexpIndex   = 4
	goCommonTimeIndex    = 5
)

func protoFileInfo(goPrefix, basepath string, protosrc []string) ProtoInfo {
	//info := fileNameInfo(dir, rel, name)
	var info ProtoInfo
	info.src = protosrc
	for _, srcpath := range info.src {
		content, err := ioutil.ReadFile(filepath.Join(basepath, srcpath))
		if err != nil {
			log.Printf("%s: error reading proto file: %v", srcpath, err)
			return info
		}

		for _, match := range protoRe.FindAllSubmatch(content, -1) {
			switch {
			case match[importSubexpIndex] != nil:
				imp := unquoteProtoString(match[importSubexpIndex])
				info.imports = append(info.imports, imp)

			case match[packageSubexpIndex] != nil:
				pkg := string(match[packageSubexpIndex])
				if info.packageName == "" {
					info.packageName = strings.Replace(pkg, ".", "_", -1)
				}

			case match[goPackageSubexpIndex] != nil:
				gopkg := unquoteProtoString(match[goPackageSubexpIndex])
				// If there's no / in the package option, then it's just a
				// simple package name, not a full import path.
				if strings.LastIndexByte(gopkg, '/') == -1 {
					info.packageName = gopkg
				} else {
					if i := strings.LastIndexByte(gopkg, ';'); i != -1 {
						info.importPath = gopkg[:i]
						info.packageName = gopkg[i+1:]
					} else {
						info.importPath = gopkg
						info.packageName = path.Base(gopkg)
					}
				}

			case match[serviceSubexpIndex] != nil:
				info.hasServices = true
			default:
				// Comment matched. Nothing to extract.
			}
		}
		rtime := regexp.MustCompile("go-common/library/time.Time")
		if rtime.FindAllSubmatchIndex(content, -1) != nil {
			info.imports = append(info.imports, "//library/time:go_default_library")
		}
		sort.Strings(info.imports)

		if info.packageName == "" {
			stem := strings.TrimSuffix(filepath.Base(srcpath), ".proto")
			fs := strings.FieldsFunc(stem, func(r rune) bool {
				return !(unicode.IsLetter(r) || unicode.IsNumber(r) || r == '_')
			})
			info.packageName = strings.Join(fs, "_")
		}
	}
	if len(info.imports) > 1 {
		info.imports = unique(info.imports)
	}
	for _, v := range info.imports {
		if strings.Contains(v, "gogo") {
			info.isGogo = true
		}
	}
	info.importPath = filepath.Join(goPrefix, basepath)
	return info
}

// Based on https://developers.google.com/protocol-buffers/docs/reference/proto3-spec
func buildProtoRegexp() *regexp.Regexp {
	hexEscape := `\\[xX][0-9a-fA-f]{2}`
	octEscape := `\\[0-7]{3}`
	charEscape := `\\[abfnrtv'"\\]`
	charValue := strings.Join([]string{hexEscape, octEscape, charEscape, "[^\x00\\'\\\"\\\\]"}, "|")
	strLit := `'(?:` + charValue + `|")*'|"(?:` + charValue + `|')*"`
	ident := `[A-Za-z][A-Za-z0-9_]*`
	fullIdent := ident + `(?:\.` + ident + `)*`
	importStmt := `\bimport\s*(?:public|weak)?\s*(?P<import>` + strLit + `)\s*;`
	packageStmt := `\bpackage\s*(?P<package>` + fullIdent + `)\s*;`
	goPackageStmt := `\boption\s*go_package\s*=\s*(?P<go_package>` + strLit + `)\s*;`
	serviceStmt := `(?P<service>service)`
	comment := `//[^\n]*`
	protoReSrc := strings.Join([]string{importStmt, packageStmt, goPackageStmt, serviceStmt, comment}, "|")
	return regexp.MustCompile(protoReSrc)
}

func unquoteProtoString(q []byte) string {
	// Adjust quotes so that Unquote is happy. We need a double quoted string
	// without unescaped double quote characters inside.
	noQuotes := bytes.Split(q[1:len(q)-1], []byte{'"'})
	if len(noQuotes) != 1 {
		for i := 0; i < len(noQuotes)-1; i++ {
			if len(noQuotes[i]) == 0 || noQuotes[i][len(noQuotes[i])-1] != '\\' {
				noQuotes[i] = append(noQuotes[i], '\\')
			}
		}
		q = append([]byte{'"'}, bytes.Join(noQuotes, []byte{'"'})...)
		q = append(q, '"')
	}
	if q[0] == '\'' {
		q[0] = '"'
		q[len(q)-1] = '"'
	}

	s, err := strconv.Unquote(string(q))
	if err != nil {
		log.Panicf("unquoting string literal %s from proto: %v", q, err)
	}
	return s
}

func unique(intSlice []string) []string {
	keys := make(map[string]interface{})
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = nil
			list = append(list, entry)
		}
	}
	return list
}
