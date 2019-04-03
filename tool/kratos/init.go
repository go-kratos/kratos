package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
	"gopkg.in/AlecAivazis/survey.v1"
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
	if !validate() {
		return nil
	}
	if err = create(); err != nil {
		fmt.Println("项目初始化失败: ", err.Error())
		return nil
	}
	fmt.Printf("项目[%s]初始化成功！\n", p.Path)
	return nil
}

func initPwd() (ok bool) {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}
	ps := strings.Split(pwd, string(os.PathSeparator))
	plen := len(ps)
	if plen < 1 {
		// 至少要有一个目录层级：项目名
		return
	}
	name := ps[plen-1]
	if name == "" {
		return
	}
	p.Name = name
	p.Path = pwd
	return true
}

func goPath() (gp string) {
	gopaths := strings.Split(os.Getenv("GOPATH"), ":")
	if len(gopaths) == 1 {
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
			fmt.Println("快速初始化失败！")
		}
		return
	case _textModeInteraction:
		// go on
	default:
		return
	}
	qs := []*survey.Question{
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
		{
			Name: "here",
			Prompt: &survey.Select{
				Message: "是否当前目录？默认为GOPATH下",
				Options: []string{_textYes, _textNo},
				Default: _textYes,
			},
		},
	}
	ans := struct {
		Name    string
		Owner   string
		UseGRPC string
		Here    string
	}{}
	if err = survey.Ask(qs, &ans); err != nil {
		return
	}
	p.Name = ans.Name
	p.Owner = ans.Owner
	if ans.UseGRPC == _textYes {
		p.WithGRPC = true
	}
	if ans.UseGRPC == _textYes {
		p.WithGRPC = true
	}
	if ans.Here == _textYes {
		p.Here = true
	}
	return
}
