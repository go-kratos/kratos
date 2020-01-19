package main

import (
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

var appHelpTemplate = `{{if .Usage}}{{.Usage}}{{end}}

USAGE:
   kratos new {{if .UsageText}}{{.UsageText}}{{else}}{{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}{{end}}{{if .Version}}{{if not .HideVersion}}

VERSION:
   {{.Version}}{{end}}{{end}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if len .Authors}}

AUTHOR{{with $length := len .Authors}}{{if ne 1 $length}}S{{end}}{{end}}:
   {{range $index, $author := .Authors}}{{if $index}}
   {{end}}{{$author}}{{end}}{{end}}{{if .VisibleCommands}}

OPTIONS:
   {{range $index, $option := .VisibleFlags}}{{if $index}}
   {{end}}{{$option}}{{end}}{{end}}{{if .Copyright}}

COPYRIGHT:
   {{.Copyright}}{{end}}
`

func main() {
	app := cli.NewApp()
	app.Name = ""
	app.Usage = "kratos 新项目创建工具"
	app.UsageText = "项目名 [options]"
	app.HideVersion = true
	app.CustomAppHelpTemplate = appHelpTemplate
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "d",
			Value:       "",
			Usage:       "指定项目所在目录",
			Destination: &p.path,
		},
		&cli.BoolFlag{
			Name:        "http",
			Usage:       "只使用http 不使用grpc",
			Destination: &p.onlyHTTP,
		},
		&cli.BoolFlag{
			Name:        "grpc",
			Usage:       "只使用grpc 不使用http",
			Destination: &p.onlyGRPC,
		},
		&cli.BoolFlag{
			Name:        "proto",
			Usage:       "废弃参数 无作用",
			Destination: &p.none,
			Hidden:      true,
		},
	}
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		app.Run([]string{"-h"})
		return
	}
	p.Name = os.Args[1]
	app.Action = runNew
	args := append([]string{os.Args[0]}, os.Args[2:]...)
	err := app.Run(args)
	if err != nil {
		panic(err)
	}
}
