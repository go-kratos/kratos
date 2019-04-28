package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var tplFlag = false
var grpcFlag = false

// 使用app/tool/warden/protoc.sh 生成grpc .pb.go
var protocShRunned = false

var syncLiveDoc = false

func main() {
	app := cli.NewApp()
	app.Name = "bmgen"
	app.Usage = "根据proto文件生成bm框架或者grpc代码: \n" +
		"用法1，在项目的任何一个位置运行bmgen，会自动找到proto文件\n" +
		"用法2，bmgen proto文件(文件必须在一个项目的api中（项目包含cmd和api目录）)\n"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:   "update",
			Usage:  "更新工具本身",
			Action: actionUpdate,
		},
		{
			Name:   "clean-live-doc",
			Usage:  "清除live-doc已经失效的分支的文档",
			Action: actionCleanLiveDoc,
		},
	}
	app.Action = actionGenerate
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "interface, i",
			Value: "",
			Usage: "generate for the interface project `RELATIVE-PATH`, eg. live/live-demo",
		},
		cli.BoolFlag{
			Name:        "grpc, g",
			Destination: &grpcFlag,
			Usage:       "废弃，会自动检测是否需要grpc",
		},
		cli.BoolFlag{
			Name:        "tpl, t",
			Destination: &tplFlag,
			Usage:       "是否生成service模板代码",
		},
		cli.BoolFlag{
			Name:        "live-doc, l",
			Destination: &syncLiveDoc,
			Usage:       "同步该项目的文档到live-doc https://git.bilibili.co/live-dev/live-doc/tree/doc/proto/$branch/$package/$filename.$Service.md",
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
		ioutil.WriteFile("/tmp/bmgentime", []byte(strconv.FormatInt(current, 10)), 0777)
		err = actionUpdate(ctx)
		if err != nil {
			return
		}
	}
	return
}

func actionCleanLiveDoc(ctx *cli.Context) (err error) {
	fmt.Println("暂未实现，敬请期待")
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
			err = runCmd("brew install protobuf")
			if err != nil {
				return
			}
		} else {
			fmt.Println("找不到protoc，请于 https://github.com/protocolbuffers/protobuf/releases 下载里面的protoc-$VERSION-$OS.zip 安装, 注意把文件拷贝到正确的未知")
			return errors.New("找不到protoc")
		}
	}
	err = runCmd("which protoc-gen-gogofast || go get github.com/gogo/protobuf/protoc-gen-gogofast")
	return
}

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
	iPath := ctx.String("i")
	if iPath != "" {
		iPath = goPath + "/src/go-common/app/interface/" + iPath
		if !fileExist(iPath) {
			return cli.NewExitError(fmt.Sprintf("interface project not found: "+iPath), 1)
		}
		pbs := filesWithSuffix(iPath+"/api", ".pb", ".proto")
		if len(pbs) == 0 {
			return cli.NewExitError(fmt.Sprintf("no pbs found in path: "+iPath+"/api"), 1)
		}
		filesToGenerate = pbs
		fmt.Printf(".pb files found %v\n", pbs)
	} else {
		if f == "" {
			// if is is empty, look up project that contains current dir
			abs, _ := filepath.Abs(".")
			proj := lookupProjPath(abs)
			if proj == "" {
				return cli.NewExitError("current dir is not in any project : "+abs, 1)
			}
			if proj != "" {
				pbs := filesWithSuffix(proj+"/api", ".pb", ".proto")
				if len(pbs) == 0 {
					return cli.NewExitError(fmt.Sprintf("no pbs found in path: "+proj+"/api"), 1)
				}
				filesToGenerate = pbs
				fmt.Printf(".pb files found %v\n", pbs)
			}
		}
	}

	for _, p := range filesToGenerate {
		if !fileExist(p) {
			return cli.NewExitError(fmt.Sprintf("file not exist: "+p), 1)
		}
		generateForFile(p, goPath)
	}
	if syncLiveDoc {
		err = actionSyncLiveDoc(ctx)
	}
	return
}

func generateForFile(f string, goPath string) {
	absPath, _ := filepath.Abs(f)
	base := filepath.Base(f)
	projPath := lookupProjPath(absPath)
	fileFolder := filepath.Dir(absPath)
	var relativePath string
	if projPath != "" {
		relativePath = absPath[len(projPath)+1:]
	}

	var cmd string
	if strings.Index(relativePath, "api/liverpc") == 0 {
		// need not generate for liverpc
		fmt.Printf("skip for liverpc \n")
		return
	}
	genTpl := 0
	if tplFlag {
		genTpl = 1
	}
	if !strings.Contains(relativePath, "api/http") {
		//非http 生成grpc和http的代码
		isInsideGoCommon := strings.Contains(fileFolder, "go-common")
		if !protocShRunned && isInsideGoCommon {
			//protoc.sh 只能在大仓库中使用
			protocShRunned = true
			cmd = fmt.Sprintf(`cd "%s" && "%s/src/go-common/app/tool/warden/protoc.sh"`, fileFolder, goPath)
			runCmd(cmd)
		}
		genGrpcArg := "--gogofast_out=plugins=grpc:."
		if isInsideGoCommon {
			// go-common中已经用上述命令生成过了
			genGrpcArg = ""
		}
		cmd = fmt.Sprintf(`cd "%s" && protoc --bm_out=tpl=%d:. `+
			`%s -I. -I%s/src -I"%s/src/go-common" -I"%s/src/go-common/vendor" "%s"`,
			fileFolder, genTpl, genGrpcArg, goPath, goPath, goPath, base)
	} else {
		// 只生成http的代码
		var pbOutArg string
		if strings.LastIndex(f, ".pb") == len(f)-3 {
			// ends with .pb
			log.Printf("\n\033[0;33m======！！！！WARNING！！！！========\n" +
				".pb文件生成代码的功能已经不再维护，请尽快迁移到.proto, 详情：\nhttp://info.bilibili.co/pages/viewpage.action?pageId=11864735#proto文件格式-.pb迁移到.proto\n" +
				"======！！！！WARNING！！！！========\033[0m\n")
			pbOutArg = "--gogofasterg_out=json_emitdefault=1,suffix=.pbg.go:."
		} else {
			pbOutArg = "--gogofast_out=."
		}
		cmd = fmt.Sprintf(`cd "%s" && protoc --bm_out=tpl=%d:. `+
			pbOutArg+` -I. -I"%s/src" -I"%s/src/go-common" -I"%s/src/go-common/vendor" "%s"`,
			fileFolder, genTpl, goPath, goPath, goPath, base)
	}

	err := runCmd(cmd)
	if err == nil {
		fmt.Println("ok")
	}
}

// files all files with suffix
func filesWithSuffix(dir string, suffixes ...string) (result []string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			for _, suffix := range suffixes {
				if strings.HasSuffix(info.Name(), suffix) {
					result = append(result, path)
					break
				}
			}
		}
		return nil
	})
	return
}

// 扫描md文档并且同步到live-doc
func actionSyncLiveDoc(ctx *cli.Context) (err error) {
	//mdFiles := filesWithSuffix(
	// if is is empty, look up project that contains current dir
	abs, _ := filepath.Abs(".")
	proj := lookupProjPath(abs)
	if proj == "" {
		err = cli.NewExitError("current dir is not in any project : "+abs, 1)
		return
	}

	var branch string
	branch, err = runCmdRet("git rev-parse --abbrev-ref HEAD")
	if err != nil {
		return
	}
	branch = url.QueryEscape(branch)

	err = runCmd("[ -d ~/.brpc ] || mkdir ~/.brpc")
	if err != nil {
		return
	}
	var liveDocUrl = "git@git.bilibili.co:live-dev/live-doc.git"
	err = runCmd("cd ~/.brpc && [ -d ~/.brpc/live-doc ] || git clone " + liveDocUrl)
	if err != nil {
		return
	}
	err = runCmd("cd ~/.brpc/live-doc && git checkout doc && git pull origin doc")
	if err != nil {
		return
	}

	mdFiles := filesWithSuffix(proj+"/api", ".md")

	var u *user.User
	u, err = user.Current()
	if err != nil {
		return err
	}
	var liveDocDir = u.HomeDir + "/.brpc/live-doc"

	for _, file := range mdFiles {
		fileHandle, _ := os.Open(file)
		fileScanner := bufio.NewScanner(fileHandle)
		for fileScanner.Scan() {
			var text = fileScanner.Text()
			var reg = regexp.MustCompile(`package=([\w\.]+)`)
			matches := reg.FindStringSubmatch(text)
			if len(matches) >= 2 {
				pkg := matches[1]
				relativeDir := strings.Replace(pkg, ".", "/", -1)
				var destDir = liveDocDir + "/protodoc/" + branch + "/" + relativeDir
				runCmd("mkdir -p " + destDir)
				err = runCmd(fmt.Sprintf("cp %s %s", file, destDir))
				if err != nil {
					return
				}
			} else {
				fmt.Println("package not found", file)
			}
			break
		}
		fileHandle.Close()
	}
	err = runCmd(`cd "` + liveDocDir + `" && git add -A`)
	if err != nil {
		return
	}
	var gitStatus string
	gitStatus, _ = runCmdRet(`cd "` + liveDocDir + `" && git status --porcelain`)
	if gitStatus != "" {
		err = runCmd(`cd "` + liveDocDir + `" && git commit -m 'update doc from proto' && git push origin doc`)
		if err != nil {
			return
		}
	}
	fmt.Printf("文档已经生成至 %s%s\n", "https://git.bilibili.co/live-dev/live-doc/tree/doc/protodoc/", url.QueryEscape(branch))
	return
}

// actionUpdate update the tools its self
func actionUpdate(ctx *cli.Context) (err error) {
	log.Print("Updating bmgen.....")
	goPath := initGopath()
	goCommonPath := goPath + "/src/go-common"
	if !fileExist(goCommonPath) {
		return cli.NewExitError("go-common not exist : "+goCommonPath, 1)
	}
	cmd := fmt.Sprintf(`go install "go-common/app/tool/bmproto/..."`)
	if err = runCmd(cmd); err != nil {
		err = cli.NewExitError(err.Error(), 1)
		return
	}
	log.Print("Updated!")
	return
}

// runCmd runs the cmd & print output (both stdout & stderr)
func runCmd(cmd string) (err error) {
	fmt.Printf("CMD: %s \n", cmd)
	out, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	fmt.Print(string(out))
	return
}

func runCmdRet(cmd string) (out string, err error) {
	fmt.Printf("CMD: %s \n", cmd)
	outBytes, err := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
	out = strings.Trim(string(outBytes), "\n\r\t ")
	return
}

// lookupProjPath get project path by proto absolute path
// assume that proto is in the project's model directory
func lookupProjPath(protoAbs string) (result string) {
	lastIndex := len(protoAbs)
	curPath := protoAbs

	for lastIndex > 0 {
		if fileExist(curPath+"/cmd") && fileExist(curPath+"/api") {
			result = curPath
			return
		}
		lastIndex = strings.LastIndex(curPath, string(os.PathSeparator))
		curPath = protoAbs[:lastIndex]
	}
	result = ""
	return
}

func fileExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil
}

func initGopath() string {
	root, err := goPath()
	if err != nil || root == "" {
		log.Printf("can not read GOPATH, use ~/go as default GOPATH")
		root = path.Join(os.Getenv("HOME"), "go")
	}
	return root
}

func goPath() (string, error) {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	if len(gopaths) == 1 {
		return gopaths[0], nil
	}
	for _, gp := range gopaths {
		absgp, err := filepath.Abs(gp)
		if err != nil {
			return "", err
		}
		if fileExist(absgp + "/src/go-common") {
			return absgp, nil
		}
	}
	return "", fmt.Errorf("can't found current gopath")
}
