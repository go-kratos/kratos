package main

import (
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

const (
	toolDoc = "https://github.com/bilibili/kratos/blob/master/doc/wiki-cn/kratos-tool.md"
)

type Tool struct {
	Name         string    `json:"name"`
	Alias        string    `json:"alias"`
	BuildTime    time.Time `json:"build_time"`
	Install      string    `json:"install"`
	Requirements []string  `json:"requirements"`
	Dir          string    `json:"dir"`
	Summary      string    `json:"summary"`
	Platform     []string  `json:"platform"`
	Author       string    `json:"author"`
	URL          string    `json:"url"`
	Hidden       bool      `json:"hidden"`
	requires     []*Tool
}

func toolAction(c *cli.Context) (err error) {
	if c.NArg() == 0 {
		sort.Slice(toolIndexs, func(i, j int) bool { return toolIndexs[i].BuildTime.After(toolIndexs[j].BuildTime) })
		for _, t := range toolIndexs {
			if t.Hidden {
				continue
			}
			updateTime := t.BuildTime.Format("2006/01/02")
			fmt.Printf("%s%s: %s Author(%s) [%s]\n", color.HiMagentaString(t.Name), getNotice(t), color.HiCyanString(t.Summary), t.Author, updateTime)
		}
		fmt.Println("\n安装工具: kratos tool install demo")
		fmt.Println("执行工具: kratos tool demo")
		fmt.Println("安装全部工具: kratos tool install all")
		fmt.Println("全部升级: kratos tool upgrade all")
		fmt.Println("\n详细文档：", toolDoc)
		return
	}
	commond := c.Args().First()
	switch commond {
	case "upgrade":
		upgradeAll()
		return
	case "install":
		name := c.Args().Get(1)
		if name == "all" {
			installAll()
		} else {
			install(name)
		}
		return
	case "check_install":
		if e := checkInstall(c.Args().Get(1)); e != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", e))
		}
		return
	}
	if e := installAndRun(commond, c.Args().Slice()[1:]); e != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("%v\n", e))
	}
	return
}

func installAndRun(name string, args []string) (err error) {
	for _, t := range toolList() {
		if name == t.Name {
			if !t.installed() || t.needUpdated() {
				t.install()
			}
			pwd, _ := os.Getwd()
			err = runTool(t.Name, pwd, t.toolPath(), args)
			return
		}
	}
	return fmt.Errorf("找不到%s", name)
}

func checkInstall(name string) (err error) {
	for _, t := range toolList() {
		if name == t.Name {
			if !t.installed() || t.needUpdated() {
				t.install()
			}
			return
		}
	}
	return fmt.Errorf("找不到%s", name)
}

func upgradeAction(c *cli.Context) error {
	install("kratos")
	return nil
}

func install(name string) {
	if name == "" {
		fmt.Fprintf(os.Stderr, color.HiRedString("请填写要安装的工具名称\n"))
		return
	}
	for _, t := range toolList() {
		if name == t.Name {
			t.install()
			return
		}
	}
	fmt.Fprintf(os.Stderr, color.HiRedString("安装失败 找不到 %s\n", name))
	return
}

func installAll() {
	for _, t := range toolList() {
		if t.Install != "" {
			t.install()
		}
	}
}

func upgradeAll() {
	for _, t := range toolList() {
		if t.needUpdated() {
			t.install()
		}
	}
}

func toolList() (tools []*Tool) {
	return toolIndexs
}

func getNotice(t *Tool) (notice string) {
	if !t.supportOS() || t.Install == "" {
		return
	}
	notice = color.HiGreenString("(未安装)")
	if t.installed() {
		notice = color.HiBlueString("(已安装)")
		if t.needUpdated() {
			notice = color.RedString("(有更新)")
		}
	}
	return
}

func (t Tool) needUpdated() bool {
	for _, r := range t.requires {
		if r.needUpdated() {
			return true
		}
	}
	if !t.supportOS() || t.Install == "" {
		return false
	}
	if f, err := os.Stat(t.toolPath()); err == nil {
		if t.BuildTime.After(f.ModTime()) {
			return true
		}
	}
	return false
}

func (t Tool) toolPath() string {
	name := t.Alias
	if name == "" {
		name = t.Name
	}
	gobin := Getenv("GOBIN")
	if gobin != "" {
		return filepath.Join(gobin, name)
	}
	return filepath.Join(gopath(), "bin", name)
}

func (t Tool) installed() bool {
	_, err := os.Stat(t.toolPath())
	return err == nil
}

func (t Tool) supportOS() bool {
	for _, p := range t.Platform {
		if strings.ToLower(p) == runtime.GOOS {
			return true
		}
	}
	return false
}

func (t Tool) install() {
	if t.Install == "" {
		fmt.Fprintf(os.Stderr, color.RedString("%s: 自动安装失败详情请查看文档：%s\n", t.Name, toolDoc))
		return
	}
	fmt.Println(t.Install)
	cmds := strings.Split(t.Install, " ")
	if len(cmds) > 0 {
		if err := runTool(t.Name, path.Dir(t.toolPath()), cmds[0], cmds[1:]); err == nil {
			color.Green("%s: 安装成功!", t.Name)
		}
	}
}

func (t Tool) updated() bool {
	if !t.supportOS() || t.Install == "" {
		return false
	}
	if f, err := os.Stat(t.toolPath()); err == nil {
		if t.BuildTime.After(f.ModTime()) {
			return true
		}
	}
	return false
}

func gopath() (gp string) {
	gopaths := strings.Split(Getenv("GOPATH"), string(filepath.ListSeparator))

	if len(gopaths) == 1 && gopaths[0] != "" {
		return gopaths[0]
	}
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	abspwd, err := filepath.Abs(pwd)
	if err != nil {
		return
	}
	for _, gopath := range gopaths {
		if gopath == "" {
			continue
		}
		absgp, err := filepath.Abs(gopath)
		if err != nil {
			return
		}
		if strings.HasPrefix(abspwd, absgp) {
			return absgp
		}
	}
	return build.Default.GOPATH
}

func runTool(name, dir, cmd string, args []string) (err error) {
	toolCmd := &exec.Cmd{
		Path:   cmd,
		Args:   append([]string{cmd}, args...),
		Dir:    dir,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	}
	if filepath.Base(cmd) == cmd {
		var lp string
		if lp, err = exec.LookPath(cmd); err == nil {
			toolCmd.Path = lp
		}
	}
	if err = toolCmd.Run(); err != nil {
		if e, ok := err.(*exec.ExitError); !ok || !e.Exited() {
			fmt.Fprintf(os.Stderr, "运行 %s 出错: %v\n", name, err)
		}
	}
	return
}
