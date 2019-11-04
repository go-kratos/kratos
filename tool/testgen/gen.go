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
			`. "github.com/smartystreets/goconvey/convey"`,
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
		impts += "\n\t\"" + impt.V + "\""
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
		if method := ConvertMethod(p.Path); method != "" {
			methodK = method + "."
		}
		tpTestFuncs := fmt.Sprintf(tpTestFunc, strings.Title(p.Package), parseFunc.Name, "", parseFunc.Name, "%s", "%s", "%s")
		tpTestFuncBeCall := methodK + parseFunc.Name + "(%s)\n\t\t\tConvey(\"%s\", func() {"
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
				so = fmt.Sprintf("\tSo(%s, ShouldBeNil)", res.K)
				notice += "err should be nil."
			} else {
				so = fmt.Sprintf("\tSo(%s, ShouldNotBeNil)", res.K)
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
		new                bool
		buffer             bytes.Buffer
		impts              string
		vars, mainFunc     string
		content            []byte
		instance, confFunc string
		tomlPath           = "**PUT PATH TO YOUR CONFIG FILES HERE**"
		filename           = strings.Replace(p.Path, ".go", "_test.go", -1)
	)
	if p.Imports["paladin"] != nil {
		new = true
	}
	// if _intfMode {
	// 	imptsList = append(imptsList, `"github.com/golang/mock/gomock"`)
	// 	for _, field := range p.Structs {
	// 		var hit bool
	// 		pkgName := strings.Split(field.V, ".")[0]
	// 		interfaceName := strings.Split(field.V, ".")[1]
	// 		if p.Imports[pkgName] != nil {
	// 			if hit, err = checkInterfaceMock(strings.Split(field.V, ".")[1], p.Imports[pkgName].V); err != nil {
	// 				return
	// 			}
	// 		}
	// 		if hit {
	// 			imptsList = append(imptsList, "mock"+p.Imports[pkgName].K+" \""+p.Imports[pkgName].V+"/mock\"")
	// 			pkgName = "mock" + strings.Title(pkgName)
	// 			interfaceName = "Mock" + interfaceName
	// 			varsList = append(varsList, "mock"+strings.Title(field.K)+" *"+pkgName+"."+interfaceName)
	// 			mockStmt += "\tmock" + strings.Title(field.K) + " = " + pkgName + ".New" + interfaceName + "(mockCtrl)\n"
	// 			newStmt += "\t\t" + field.K + ":\tmock" + strings.Title(field.K) + ",\n"
	// 		} else {
	// 			pkgName = subString(field.V, "*", ".")
	// 			if p.Imports[pkgName] != nil && pkgName != "conf" {
	// 				imptsList = append(imptsList, p.Imports[pkgName].K+" \""+p.Imports[pkgName].V+"\"")
	// 			}
	// 			switch {
	// 			case strings.HasPrefix(field.V, "*conf."):
	// 				newStmt += "\t\t" + field.K + ":\tconf.Conf,\n"
	// 			case strings.HasPrefix(field.V, "*"):
	// 				newStmt += "\t\t" + field.K + ":\t" + strings.Replace(field.V, "*", "&", -1) + "{},\n"
	// 			default:
	// 				newStmt += "\t\t" + field.K + ":\t" + field.V + ",\n"
	// 			}
	// 		}
	// 	}
	// 	mockStmt = fmt.Sprintf(_tpTestServiceMainMockStmt, mockStmt)
	// 	newStmt = fmt.Sprintf(_tpTestServiceMainNewStmt, newStmt)
	// }
	if instance = ConvertMethod(p.Path); instance == "s" {
		vars = strings.Join([]string{"s *Service"}, "\n\t")
		mainFunc = tpTestServiceMain
	} else {
		vars = strings.Join([]string{"d *Dao"}, "\n\t")
		mainFunc = tpTestDaoMain
	}
	if new {
		impts = strings.Join([]string{`"os"`, `"flag"`, `"testing"`, p.Imports["paladin"].V}, "\n\t")
		confFunc = fmt.Sprintf(tpTestMainNew, instance+" = New()")
	} else {
		impts = strings.Join(append([]string{`"os"`, `"flag"`, `"testing"`}), "\n\t")
		confFunc = fmt.Sprintf(tpTestMainOld, instance+" = New(conf.Conf)")
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		buffer.WriteString(fmt.Sprintf(tpPackage, p.Package))
		buffer.WriteString(fmt.Sprintf(tpImport, impts))
		buffer.WriteString(fmt.Sprintf(tpVar, vars))
		buffer.WriteString(fmt.Sprintf(mainFunc, tomlPath, confFunc))
		content, _ = GoImport(filename, buffer.Bytes())
		ioutil.WriteFile(filename, content, 0644)
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
		if strings.Contains(parse.Path, "monkey.go") ||
			strings.Contains(parse.Path, "/mock/") {
			continue
		}
		var (
			path                      = strings.Split(filepath.Dir(parse.Path), "/")
			pack                      = ConvertHump(path[len(path)-1])
			refer                     = path[len(path)-1]
			mockVar, mockType, srcDir string
		)
		for i := len(path) - 1; i > len(path)-4; i-- {
			if path[i] == "dao" || path[i] == "service" {
				srcDir = strings.Join(path[:i+1], "/")
				break
			}
			pack = ConvertHump(path[i-1]) + pack
		}
		if mockVar = ConvertMethod(parse.Path); mockType == "d" {
			mockType = "*" + refer + ".Dao"
		} else {
			mockType = "*" + refer + ".Service"
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
