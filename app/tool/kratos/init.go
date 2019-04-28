package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey"
	"github.com/urfave/cli"
)

var (
	// 允许建立项目的部门
	depts = map[string]bool{
		"main":         true,
		"live":         true,
		"openplatform": true,
		"search":       true,
		"ep":           true,
		"bbq":          true,
		"video":        true,
		"bplus":        true,
		"ops":          true,
	}
	// 允许建立的项目类型
	types = map[string]bool{
		"interface": true,
		"admin":     true,
		"job":       true,
		"service":   true,
	}
)

const (
	_textModeFastInit    = "一键初始化项目"
	_textModeInteraction = "自定义项目参数"
	_textYes             = "是"
	_textNo              = "否"
)

func runInit(ctx *cli.Context) (err error) {
	if ctx.NumFlags() == 0 {
		if err = interact(); err != nil {
			return
		}
	}
	if ok := check(); !ok {
		return nil
	}
	if err = create(); err != nil {
		println("项目初始化失败: ", err.Error())
		return nil
	}
	fmt.Printf(`项目初始化成功！
注意：请先创建rider、服务树节点、在配置中心创建uat环境配置文件，否则提交mr后无法运行单元测试！
相关帮助信息见 http://info.bilibili.co/pages/viewpage.action?pageId=7567510
`)
	return nil
}

func initPwd() (ok bool) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	ps := strings.Split(pwd, string(os.PathSeparator))
	plen := len(ps)
	if plen < 3 {
		// 至少要有三个目录层级：部门、项目类型、项目名
		return
	}
	name := ps[plen-1]
	dept := ps[plen-2]
	typ := ps[plen-3]
	if !depts[dept] {
		return
	}
	if !types[typ] {
		return
	}
	if name == "" {
		return
	}
	p.Name = name
	p.Department = dept
	p.Type = typ
	p.Path = pwd
	return true
}

func check() (ok bool) {
	root, err := goPath()
	if err != nil || root == "" {
		log.Printf("can not read GOPATH, use ~/go as default GOPATH")
		root = path.Join(os.Getenv("HOME"), "go")
	}
	if !validate() {
		return
	}
	p.Path = fmt.Sprintf("%s/src/go-common/app/%s/%s/%s", strings.TrimRight(root, "/"), p.Type, p.Department, p.Name)
	return true
}

func goPath() (string, error) {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	if len(gopaths) == 1 {
		return gopaths[0], nil
	}
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abspwd, err := filepath.Abs(pwd)
	if err != nil {
		return "", err
	}
	for _, gp := range gopaths {
		absgp, err := filepath.Abs(gp)
		if err != nil {
			return "", err
		}
		if strings.HasPrefix(abspwd, absgp) {
			return absgp, nil
		}
	}
	return "", fmt.Errorf("can't found current gopath")
}

func interact() (err error) {
	qs1 := &survey.Select{
		Message: "你想怎么玩？",
		Options: []string{_textModeFastInit, _textModeInteraction},
	}
	var ans1 string
	if err = survey.AskOne(qs1, &ans1, nil); err != nil {
		return
	}
	switch ans1 {
	case _textModeFastInit:
		if ok := initPwd(); !ok {
			println("Notice: Not in project directory. Skipped fast init.")
		}
		return
	case _textModeInteraction:
		// go on
	default:
		return
	}
	var ds, ts []string
	for d := range depts {
		ds = append(ds, d)
	}
	for t := range types {
		ts = append(ts, t)
	}
	qs := []*survey.Question{
		{
			Name: "department",
			Prompt: &survey.Select{
				Message: "请选择选择部门：",
				Options: ds,
				Default: "main",
			},
		},
		{
			Name: "type",
			Prompt: &survey.Select{
				Message: "请选择项目类型：",
				Options: ts,
			},
		},
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "请输入项目名称：",
			},
			Validate: survey.Required,
		},
		{
			Name: "owner",
			Prompt: &survey.Input{
				Message: "请输入项目负责人：",
			},
		},
		{
			Name: "useGRPC",
			Prompt: &survey.Select{
				Message: "是否使用 gRPC ？",
				Options: []string{_textYes, _textNo},
				Default: _textNo,
			},
		},
	}
	ans := struct {
		Department string
		Type       string
		Name       string
		Owner      string
		UseGRPC    string
	}{}
	if err = survey.Ask(qs, &ans); err != nil {
		return
	}
	p.Name = ans.Name
	p.Department = ans.Department
	p.Type = ans.Type
	p.Owner = ans.Owner
	if ans.UseGRPC == _textYes {
		p.WithGRPC = true
	}
	return
}
