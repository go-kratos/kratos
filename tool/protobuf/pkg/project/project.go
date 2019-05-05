package project

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bilibili/kratos/tool/protobuf/pkg/utils"
	"github.com/pkg/errors"
	"github.com/siddontang/go/ioutil2"
)

// if proto file is inside a project (that has a /api directory)
// this present a project info
// 必须假设proto文件的路径就是相对work-dir的路径，否则无法找到proto文件以及对应的project
type ProjectInfo struct {
	// AbsolutePath of the project
	AbsolutePath string
	// ImportPath of project
	ImportPath string
	// dir name of project
	Name string
	// parent dir of project, maybe empty
	Department string
	// grandma dir of project , maybe empty
	Typ            string
	HasInternalPkg bool
	// 从工作目录(working directory)到project目录的相对路径 比如a/b .. ../a
	// 作用是什么？
	// 假设目录结构是
	// -project
	//   - api/api.proto
	//   - internal/service
	// 我想在 internal/service下生成一个文件 service.go
	// work-dir 为 project/api
	// proto 生成命令为 protoc --xx_out=. api.proto
	// 那么在 protoc plugin 中的文件输出路径就得是 ../internal/service/ => {pathRefToProj}internal/service
	//
	PathRefToProj string
}

func NewProjInfo(file string, modDirName string, modImportPath string) (projInfo *ProjectInfo, err error) {
	projInfo = &ProjectInfo{}
	wd, err := os.Getwd()
	if err != nil {
		panic("cannot get working directory")
	}
	protoAbs := wd + "/" + file
	protoAbs, _ = filepath.Abs(protoAbs)

	if !ioutil2.FileExists(protoAbs) {
		return nil, errors.Errorf("Cannot find proto file in current dir %s, file: %s ", wd, file)
	}
	//appIndex := strings.Index(wd, modDirName)
	//if appIndex == -1 {
	//	err = errors.New("not in " + modDirName)
	//	return nil, err
	//}

	projPath := LookupProjPath(protoAbs)

	if projPath == "" {
		err = errors.New("not in project")
		return nil, err
	}
	rel, _ := filepath.Rel(wd, projPath)
	projInfo.PathRefToProj = rel
	projInfo.AbsolutePath = projPath
	if ioutil2.FileExists(projPath + "/internal") {
		projInfo.HasInternalPkg = true
	}

	i := strings.Index(projInfo.AbsolutePath, modDirName)
	if i == -1 {
		err = errors.Errorf("project is not inside module, project=%s, module=%s", projPath, modDirName)
	}
	relativePath := projInfo.AbsolutePath[i+len(modDirName):]
	projInfo.ImportPath = modImportPath + relativePath
	projInfo.Name = filepath.Base(projPath)
	if p := filepath.Dir(projPath); p != "/" {
		projInfo.Department = filepath.Base(p)
		if p = filepath.Dir(p); p != "/" {
			projInfo.Typ = filepath.Base(p)
		}
	}
	return projInfo, nil
}

// LookupProjPath get project path by proto absolute path
// assume that proto is in the project's api directory
func LookupProjPath(protoAbs string) (result string) {
	f := func(protoAbs string, dirs []string) string {
		lastIndex := len(protoAbs)
		curPath := protoAbs

		for lastIndex > 0 {
			found := true
			for _, d := range dirs {
				if !utils.IsDir(curPath + "/" + d) {
					found = false
					break
				}
			}
			if found {
				return curPath
			}
			lastIndex = strings.LastIndex(curPath, string(os.PathSeparator))
			curPath = protoAbs[:lastIndex]
		}
		result = ""
		return result
	}

	firstStep := f(protoAbs, []string{"cmd", "api"})
	if firstStep != "" {
		return firstStep
	}
	return f(protoAbs, []string{"api"})
}
