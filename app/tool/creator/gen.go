package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/otokaze/mock/mockgen"
	"github.com/otokaze/mock/mockgen/model"
)

func genTest(parses []*parse) (err error) {
	for _, p := range parses {
		switch {
		case strings.HasSuffix(p.Path, "_mock.go") ||
			strings.HasSuffix(p.Path, ".intf.go"):
			continue
		case strings.HasSuffix(p.Path, "dao.go") ||
			strings.HasSuffix(p.Path, "service.go"):
			err = p.genTestMain()
		default:
			err = p.genUTTest()
		}
		if err != nil {
			break
		}
	}
	return
}

func (p *parse) genUTTest() (err error) {
	var (
		buffer bytes.Buffer
		impts  = strings.Join([]string{
			`"context"`,
			`"testing"`,
			`"github.com/smartystreets/goconvey/convey"`,
		}, "\n\t")
		content []byte
	)
	filename := strings.Replace(p.Path, ".go", "_test.go", -1)
	if _, err = os.Stat(filename); (_func == "" && err == nil) ||
		(err != nil && os.IsExist(err)) {
		err = nil
		return
	}
	for _, impt := range p.Imports {
		impts += "\n\t" + impt
	}
	if _func == "" {
		buffer.WriteString(fmt.Sprintf(tpPackage, p.Package))
		buffer.WriteString(fmt.Sprintf(tpImport, impts))
	}
	for _, parseFunc := range p.Funcs {
		if _func != "" && _func != parseFunc.Name {
			continue
		}
		var (
			methodK string
			tpVars  string
			vars    []string
			val     []string
			notice  = "Then "
			reset   string
		)
		if parseFunc.Method != nil {
			switch {
			case strings.Contains(p.Path, "/library"):
				methodK = ""
			case strings.Contains(p.Path, "/service"):
				methodK = "s."
			case strings.Contains(p.Path, "/dao"):
				methodK = "d."
			}
			// if IsService(p.Package) {
			// 	methodK = "s."
			// } else {
			// 	methodK = "d."
			// }
		}
		tpTestFuncs := fmt.Sprintf(tpTestFunc, strings.Title(p.Package), parseFunc.Name, parseFunc.Name, "%s", "%s", "%s")
		tpTestFuncBeCall := methodK + parseFunc.Name + "(%s)\n\t\t\tconvCtx.Convey(\"%s\", func(convCtx convey.C) {"
		if parseFunc.Result == nil {
			tpTestFuncBeCall = fmt.Sprintf(tpTestFuncBeCall, "%s", "No return values")
			tpTestFuncs = fmt.Sprintf(tpTestFuncs, "%s", tpTestFuncBeCall, "%s")
		}
		for k, res := range parseFunc.Result {
			if res.K == "" {
				res.K = fmt.Sprintf("p%d", k+1)
			}
			var so string
			if res.V == "error" {
				res.K = "err"
				so = fmt.Sprintf("\tconvCtx.So(%s, convey.ShouldBeNil)", res.K)
				notice += "err should be nil."
			} else {
				so = fmt.Sprintf("\tconvCtx.So(%s, convey.ShouldNotBeNil)", res.K)
				val = append(val, res.K)
			}
			if len(parseFunc.Result) <= k+1 {
				if len(val) != 0 {
					notice += strings.Join(val, ",") + " should not be nil."
				}
				tpTestFuncBeCall = fmt.Sprintf(tpTestFuncBeCall, "%s", notice)
				res.K += " := " + tpTestFuncBeCall
			} else {
				res.K += ", %s"
			}
			tpTestFuncs = fmt.Sprintf(tpTestFuncs, "%s", res.K+"\n\t\t\t%s", "%s")
			tpTestFuncs = fmt.Sprintf(tpTestFuncs, "%s", "%s", so, "%s")
		}
		if parseFunc.Params == nil {
			tpTestFuncs = fmt.Sprintf(tpTestFuncs, "%s", "", "%s")
		}
		for k, pType := range parseFunc.Params {
			if pType.K == "" {
				pType.K = fmt.Sprintf("a%d", k+1)
			}
			var (
				init   string
				params = pType.K
			)
			switch {
			case strings.HasPrefix(pType.V, "context"):
				init = params + " = context.Background()"
			case strings.HasPrefix(pType.V, "[]byte"):
				init = params + " = " + pType.V + "(\"\")"
			case strings.HasPrefix(pType.V, "[]"):
				init = params + " = " + pType.V + "{}"
			case strings.HasPrefix(pType.V, "int") ||
				strings.HasPrefix(pType.V, "uint") ||
				strings.HasPrefix(pType.V, "float") ||
				strings.HasPrefix(pType.V, "double"):
				init = params + " = " + pType.V + "(0)"
			case strings.HasPrefix(pType.V, "string"):
				init = params + " = \"\""
			case strings.Contains(pType.V, "*xsql.Tx"):
				// init = params + ",_ = d.BeginTran(c)"
				init = params + ",_ = " + methodK + "BeginTran(c)"
				reset += "\n\t" + params + ".Commit()"
			case strings.HasPrefix(pType.V, "*"):
				init = params + " = " + strings.Replace(pType.V, "*", "&", -1) + "{}"
			case strings.Contains(pType.V, "chan"):
				init = params + " = " + pType.V
			case pType.V == "time.Time":
				init = params + " = time.Now()"
			case strings.Contains(pType.V, "chan"):
				init = params + " = " + pType.V
			default:
				init = params + " " + pType.V
			}
			vars = append(vars, "\t\t"+init)
			if len(parseFunc.Params) > k+1 {
				params += ", %s"
			}
			tpTestFuncs = fmt.Sprintf(tpTestFuncs, "%s", params, "%s")
		}
		if len(vars) > 0 {
			tpVars = fmt.Sprintf(tpVar, strings.Join(vars, "\n\t"))
		}
		tpTestFuncs = fmt.Sprintf(tpTestFuncs, tpVars, "%s")
		if reset != "" {
			tpTestResets := fmt.Sprintf(tpTestReset, reset)
			tpTestFuncs = fmt.Sprintf(tpTestFuncs, tpTestResets)
		} else {
			tpTestFuncs = fmt.Sprintf(tpTestFuncs, "")
		}
		buffer.WriteString(tpTestFuncs)
	}
	var (
		file *os.File
		flag = os.O_RDWR | os.O_CREATE | os.O_APPEND
	)
	if file, err = os.OpenFile(filename, flag, 0644); err != nil {
		return
	}
	if _func == "" {
		content, _ = GoImport(filename, buffer.Bytes())
	} else {
		content = buffer.Bytes()
	}
	if _, err = file.Write(content); err != nil {
		return
	}
	if err = file.Close(); err != nil {
		return
	}
	return
}

func (p *parse) genTestMain() (err error) {
	var (
		buffer                  bytes.Buffer
		left, right             int
		vars, mainFunc, fileAbs string
	)
	if fileAbs, err = filepath.Abs(p.Path); err != nil {
		return
	}
	var pathSplit []string
	if pathSplit = strings.Split(fileAbs, "/"); len(pathSplit) == 1 {
		pathSplit = strings.Split(fileAbs, "\\")
	}

	for index, value := range pathSplit {
		if value == "go-common" {
			left = index
		}
		if value == "main" ||
			value == "ep" ||
			value == "openplatform" ||
			value == "live" ||
			value == "video" ||
			value == "tool" ||
			value == "bbq" ||
			value == "library" {
			right = index
			break
		}
	}
	var (
		tomlPath string
		filePath = filepath.Dir(fileAbs)
		cmdFiles []os.FileInfo
		cmdPath  = strings.Join(pathSplit[:right+2], "/") + "/cmd"
	)
	if cmdFiles, err = ioutil.ReadDir(cmdPath); err != nil {
		return
	}
	for _, f := range cmdFiles {
		if !strings.HasSuffix(f.Name(), ".toml") {
			continue
		}
		if tomlPath, err = filepath.Rel(filePath, cmdPath); err != nil {
			return
		}
		if tomlPath != "" {
			tomlPath = tomlPath + "/" + f.Name()
			break
		}
	}
	var (
		filename = strings.Replace(p.Path, ".go", "_test.go", -1)
		flag     = pathSplit[right+2 : right+3][0]
		conf     = strings.Join(pathSplit[left:right+2], "/") + "/conf"
		impts    = strings.Join([]string{`"os"`, `"flag"`, `"testing"`, "\"" + conf + "\""}, "\n\t")
	)
	if IsService(flag) {
		vars = strings.Join([]string{"s *Service"}, "\n\t")
		mainFunc = tpTestServiceMain
	} else {
		vars = strings.Join([]string{"d *Dao"}, "\n\t")
		mainFunc = tpTestDaoMain
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		buffer.WriteString(fmt.Sprintf(tpPackage, p.Package))
		buffer.WriteString(fmt.Sprintf(tpImport, impts))
		buffer.WriteString(fmt.Sprintf(tpVar, vars))
		buffer.WriteString(fmt.Sprintf(mainFunc, tomlPath))
		ioutil.WriteFile(filename, buffer.Bytes(), 0644)
	}
	return
}

func genInterface(parses []*parse) (err error) {
	var (
		parse *parse
		pkg   = make(map[string]string)
	)
	for _, parse = range parses {
		if strings.Contains(parse.Path, ".intf.go") {
			continue
		}
		dirPath := filepath.Dir(parse.Path)
		for _, parseFunc := range parse.Funcs {
			if (parseFunc.Method == nil) ||
				!(parseFunc.Name[0] >= 'A' && parseFunc.Name[0] <= 'Z') {
				continue
			}
			var (
				params  string
				results string
			)
			for k, param := range parseFunc.Params {
				params += param.K + " " + param.P + param.V
				if len(parseFunc.Params) > k+1 {
					params += ", "
				}
			}
			for k, res := range parseFunc.Result {
				results += res.K + " " + res.P + res.V
				if len(parseFunc.Result) > k+1 {
					results += ", "
				}
			}
			if len(results) != 0 {
				results = "(" + results + ")"
			}
			pkg[dirPath] += "\t" + fmt.Sprintf(tpIntfcFunc, parseFunc.Name, params, results)
		}
	}
	for k, v := range pkg {
		var buffer bytes.Buffer
		pathSplit := strings.Split(k, "/")
		filename := k + "/" + pathSplit[len(pathSplit)-1] + ".intf.go"
		if _, exist := os.Stat(filename); os.IsExist(exist) {
			continue
		}
		buffer.WriteString(fmt.Sprintf(tpPackage, pathSplit[len(pathSplit)-1]))
		buffer.WriteString(fmt.Sprintf(tpInterface, strings.Title(pathSplit[len(pathSplit)-1]), v))
		content, _ := GoImport(filename, buffer.Bytes())
		err = ioutil.WriteFile(filename, content, 0644)
	}
	return
}

func genMock(files ...string) (err error) {
	for _, file := range files {
		var pkg *model.Package
		if pkg, err = mockgen.ParseFile(file); err != nil {
			return
		}
		if len(pkg.Interfaces) == 0 {
			continue
		}
		var mockDir = pkg.SrcDir + "/mock"
		if _, err = os.Stat(mockDir); os.IsNotExist(err) {
			err = nil
			os.Mkdir(mockDir, 0744)
		}
		var mockPath = mockDir + "/" + pkg.Name + "_mock.go"
		if _, exist := os.Stat(mockPath); os.IsExist(exist) {
			continue
		}
		var g = &mockgen.Generator{Filename: file}
		if err = g.Generate(pkg, "mock", mockPath); err != nil {
			return
		}
		if err = ioutil.WriteFile(mockPath, g.Output(), 0644); err != nil {
			return
		}
	}
	return
}

func genMonkey(parses []*parse) (err error) {
	var (
		pkg = make(map[string]string)
	)
	for _, parse := range parses {
		if strings.Contains(parse.Path, "monkey.go") || strings.Contains(parse.Path, "/mock/") {
			continue
		}
		var (
			mockVar, mockType, srcDir, flag string
			path                            = strings.Split(filepath.Dir(parse.Path), "/")
			pack                            = ConvertHump(path[len(path)-1])
			refer                           = path[len(path)-1]
		)
		for i := len(path) - 1; i > len(path)-4; i-- {
			if path[i] == "dao" || path[i] == "service" {
				srcDir = strings.Join(path[:i+1], "/")
				flag = path[i]
				break
			}
			pack = ConvertHump(path[i-1]) + pack
		}
		if IsService(flag) {
			mockVar = "s"
			mockType = "*" + refer + ".Service"
		} else {
			mockVar = "d"
			mockType = "*" + refer + ".Dao"
		}
		for _, parseFunc := range parse.Funcs {
			if (parseFunc.Method == nil) || (parseFunc.Result == nil) ||
				!(parseFunc.Name[0] >= 'A' && parseFunc.Name[0] <= 'Z') {
				continue
			}
			var (
				funcParams, funcResults, mockKey, mockValue, funcName string
			)
			funcName = pack + parseFunc.Name
			for k, param := range parseFunc.Params {
				funcParams += "_ " + param.V
				if len(parseFunc.Params) > k+1 {
					funcParams += ", "
				}
			}
			for k, res := range parseFunc.Result {
				if res.K == "" {
					if res.V == "error" {
						res.K = "err"
					} else {
						res.K = fmt.Sprintf("p%d", k+1)
					}
				}
				mockKey += res.K
				mockValue += res.V
				funcResults += res.K + " " + res.P + res.V
				if len(parseFunc.Result) > k+1 {
					mockKey += ", "
					mockValue += ", "
					funcResults += ", "
				}
			}
			pkg[srcDir+"."+refer] += fmt.Sprintf(tpMonkeyFunc, funcName, funcName, mockVar, mockType, funcResults, mockVar, parseFunc.Name, mockType, funcParams, mockValue, mockKey)
		}
	}
	for path, content := range pkg {
		var (
			buffer   bytes.Buffer
			dir      = strings.Split(path, ".")
			mockDir  = dir[0] + "/mock"
			filename = mockDir + "/monkey_" + dir[1] + ".go"
		)
		if _, err = os.Stat(mockDir); os.IsNotExist(err) {
			err = nil
			os.Mkdir(mockDir, 0744)
		}
		if _, err := os.Stat(filename); os.IsExist(err) {
			continue
		}
		buffer.WriteString(fmt.Sprintf(tpPackage, "mock"))
		buffer.WriteString(content)
		content, _ := GoImport(filename, buffer.Bytes())
		ioutil.WriteFile(filename, content, 0644)
	}
	return
}
