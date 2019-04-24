package main

import (
	"bufio"
	"fmt"
	"kratos/tool/bmproto/pkg/project"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

type bmgenConf struct {
	Tpl          bool   `toml:"tpl"`
	ExplicitHTTP bool   `toml:"explicit_http"`
	AppID        string `toml:"app_id"`
}

var fileConf = &bmgenConf{}
var flagVerbose bool
var flagExplicitHTTP bool
var flagTpl bool
var flagAppID string

// 使用 kratos/tool/warden/protoc.sh 生成 grpc .pb.go
// key 为目录，表示这个目录有没有运行过 protoc.sh
var protocShRunnedMap = map[string]struct{}{}

func main() {
	app := cli.NewApp()
	app.Name = "bmgen"
	app.Usage = "根据proto文件生成bm框架或者grpc代码: \n" +
		"用法1，在项目的任何一个位置运行bmgen，会自动找到proto文件（项目包含cmd和api目录）\n"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:   "update",
			Usage:  "更新工具本身",
			Action: actionUpdate,
		},
	}
	app.Action = actionGenerate
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "vv",
			Destination: &flagVerbose,
			Usage:       "打印详细信息",
		},
		cli.BoolFlag{
			Name:        "tpl, t",
			Destination: &flagTpl,
			Usage:       "是否生成service模板代码",
		},
		cli.BoolFlag{
			Name:        "explicit_http, e",
			Destination: &flagExplicitHTTP,
			Usage:       "只有显示指定的http option的method才生成相应http代码",
		},
		cli.StringFlag{
			Name:        "app_id",
			Destination: &flagAppID,
			Usage:       "服务树名字, 例如live.live.room-service, 同步到bapi时需要使用",
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func autoUpdate(ctx *cli.Context) (err error) {
	tfile, e := ioutil.ReadFile("/tmp/bmgentime")
	if e != nil {
		tfile = []byte{'0'}
	}
	ts, _ := strconv.ParseInt(string(tfile), 0, 64)
	current := time.Now().Unix()
	if (current - int64(ts)) > 3600*12 { // 12 hours old
		_ = ioutil.WriteFile("/tmp/bmgentime", []byte(strconv.FormatInt(current, 10)), 0777)
		err = actionUpdate(ctx)
		if err != nil {
			return
		}
	}
	return
}

func installDependencies() (err error) {
	// 检查protoc
	//
	e := runCmd("which protoc")
	if e != nil {
		var uname string
		uname, err = runCmdRet("uname")
		if err != nil {
			return
		}
		if uname == "Darwin" {
			infof("Installing protobuf using homebrew...")
			err = runCmd("brew install protobuf")
			if err != nil {
				return
			}
		} else {
			errorf("找不到protoc，请于 https://github.com/protocolbuffers/protobuf/releases 下载里面的protoc-$VERSION-$OS.zip 安装, 注意把文件拷贝到正确的未知")
			return errors.New("找不到protoc")
		}
	}
	err = runCmd("which protoc-gen-gogofast || go get github.com/gogo/protobuf/protoc-gen-gogofast")
	return
}

const projectNotFoundErr = "当前目录没有在任何项目内，一个项目应当至少包含 /api 目录"

// actionGenerate invokes protoc to generate files
func actionGenerate(ctx *cli.Context) (err error) {
	if err = autoUpdate(ctx); err != nil {
		return
	}
	err = installDependencies()
	if err != nil {
		return
	}
	f := ctx.Args().Get(0)

	goPath := initGopath()
	if !fileExist(goPath) {
		return cli.NewExitError(fmt.Sprintf("GOPATH not exist: "+goPath), 1)
	}

	filesToGenerate := []string{f}
	var proj string
	if f == "" {
		// if is is empty, look up project that contains current dir
		abs, _ := filepath.Abs(".")
		proj = project.LookupProjPath(abs)
		if proj == "" {
			return cli.NewExitError(projectNotFoundErr, 1)
		}
		if proj != "" {
			pbs := filesWithSuffix(proj+"/api", ".pb", ".proto")
			if len(pbs) == 0 {
				return cli.NewExitError(fmt.Sprintf("no pbs found in path: "+proj+"/api"), 1)
			}
			filesToGenerate = pbs
			logf("proto files found %v", pbs)
		}
	} else {
		abs, _ := filepath.Abs(f)
		proj = project.LookupProjPath(abs)
	}

	initProjConf(proj)

	for _, p := range filesToGenerate {
		if !fileExist(p) {
			return cli.NewExitError(fmt.Sprintf("file not exist: "+p), 1)
		}
		err = generateForFile(p, goPath)
		if err != nil {
			return
		}
	}
	if err == nil {
		infof("Done.")
	}
	return
}

func initProjConf(projPath string) {
	if projPath == "" {
		return
	}
	var confPath = projPath + "/bmgen.toml"
	if !fileExist(confPath) {
		return
	}
	infof("使用配置: " + confPath)
	_, err := toml.DecodeFile(confPath, &fileConf)
	if err != nil {
		errorf("Decode 配置错误:  %+v", err)
		return
	}
	if fileConf.Tpl {
		flagTpl = true
		infof("flag from file: tpl=true")
	}
	if fileConf.ExplicitHTTP {
		flagExplicitHTTP = true
		infof("flag from file: explicit_http=true")
	}
	if flagAppID == "" && fileConf.AppID != "" {
		flagAppID = fileConf.AppID
		infof("flag from file: app_id=" + flagAppID)
	}
}

func generateForFile(f string, goPath string) (err error) {
	absPath, _ := filepath.Abs(f)
	projPath := project.LookupProjPath(absPath)
	fileFolder := filepath.Dir(absPath)
	var relativePath string
	if projPath != "" {
		relativePath = absPath[len(projPath)+1:]
	}
	var goSrcPath = goPath + "/src"
	var relativeInGoPath = absPath[len(goPath+"/src")+1:]

	var cmd string
	if strings.Index(relativePath, "api/liverpc") == 0 {
		// need not generate for liverpc
		logf("skip for liverpc")
		return
	}
	explicitHTTPArg := 0
	if flagExplicitHTTP {
		explicitHTTPArg = 1
	}

	var (
		tplOutArg     string
		bmOutArg      string
	)
	if flagTpl {
		tplOutArg = "--tpl_out=."
	}
	bmOutArg = fmt.Sprintf("--bm_out=explicit_http=%d:.", explicitHTTPArg)
	if !strings.Contains(relativePath, "api/http") {
		//非http 生成grpc和http的代码
		isInsideGoCommon := strings.Contains(fileFolder, "go-common")
		if _, protocShRunned := protocShRunnedMap[fileFolder]; !protocShRunned && isInsideGoCommon {
			//protoc.sh 只能在大仓库中使用
			protocShRunnedMap[fileFolder] = struct{}{}
			cmd = fmt.Sprintf(`cd "%s" && "%s/src/kratos/tool/warden/protoc.sh"`, fileFolder, goPath)
			err = runCmd(cmd)
			if err != nil {
				return
			}
		}
		genGrpcArg := "--gogofast_out=plugins=grpc:."
		if isInsideGoCommon {
			genGrpcArg = ""
		}

		cmd = fmt.Sprintf(`cd "%s" && protoc %s %s %s`+
			` -I%s -I"%s/go-common" -I"%s/go-common/vendor" "%s"`,
			goSrcPath, bmOutArg, tplOutArg, genGrpcArg,
			goSrcPath, goSrcPath, goSrcPath, relativeInGoPath)
	} else {
		// 只生成http的代码
		var pbOutArg string
		if strings.LastIndex(f, ".pb") == len(f)-3 {
			// ends with .pb

			errorf("\n======！！！！WARNING！！！！========\n" +
				".pb文件生成代码的功能已经不再维护，请尽快迁移到.proto\n" +
				"======！！！！WARNING！！！！========")
			e := runCmd("which protoc-gen-gogofasterg")
			if e != nil {
				infof("Installing protoc-gen-gogofasterg")
				err = runCmd("go get github.com/liues1992/gogoprotobuf/protoc-gen-gogofasterg")
				if err != nil {
					return
				}
			}

			pbOutArg = "--gogofasterg_out=json_emitdefault=1,suffix=.pbg.go:."
		} else {
			pbOutArg = "--gogofast_out=."
		}
		cmd = fmt.Sprintf(`cd "%s" && protoc %s %s %s`+
			` -I"%s" -I"%s/go-common" -I"%s/go-common/vendor" "%s"`,
			goSrcPath, bmOutArg, tplOutArg, pbOutArg,
			goSrcPath, goSrcPath, goSrcPath, relativeInGoPath)
	}

	err = runCmd(cmd)
	if err != nil {
		return
	}
	return
}

// file is a go source file
// return empty string if not found
func findGoPkgForFile(file string) string {
	fileHandle, _ := os.Open(file)
	fileScanner := bufio.NewScanner(fileHandle)
	for fileScanner.Scan() {
		var text = fileScanner.Text()
		var reg = regexp.MustCompile(`package (.+)`)
		matches := reg.FindStringSubmatch(text)
		if len(matches) >= 2 {
			pkg := matches[1]
			return pkg
		}
	}
	// package has to be found in a go file
	panic("package not found in a go file: " + file)
}

// actionUpdate update the tools its self
func actionUpdate(ctx *cli.Context) (err error) {
	logf("Updating bmgen.....")
	goPath := initGopath()
	goCommonPath := goPath + "/src/go-common"
	if !fileExist(goCommonPath) {
		return cli.NewExitError("go-common not exist : "+goCommonPath, 1)
	}
	cmd := fmt.Sprintf(`go install "kratos/tool/bmproto/..."`)
	if err = runCmd(cmd); err != nil {
		err = cli.NewExitError(err.Error(), 1)
		return
	}
	logf("Updated!")
	return
}
