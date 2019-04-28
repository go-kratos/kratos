package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "kratos"
	app.Usage = "kratos tool"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "bazel build",
			Action:  bazelAction,
		},
		{
			Name:    "init",
			Aliases: []string{"i"},
			Usage:   "create new project",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "d",
					Value:       "",
					Usage:       "department name for create project",
					Destination: &p.Department,
				},
				cli.StringFlag{
					Name:        "t",
					Value:       "",
					Usage:       "project type name for create project",
					Destination: &p.Type,
				},
				cli.StringFlag{
					Name:        "n",
					Value:       "",
					Usage:       "project name for create project",
					Destination: &p.Name,
				},
				cli.StringFlag{
					Name:        "o",
					Value:       "",
					Usage:       "project owner for create project",
					Destination: &p.Owner,
				},
				cli.BoolFlag{
					Name:        "grpc",
					Usage:       "whether to use grpc for create project",
					Destination: &p.WithGRPC,
				},
			},
			Action: runInit,
		},
		{
			Name:    "update",
			Aliases: []string{"u"},
			Usage:   "update bazel building configure",
			Action:  updateAction,
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "kratos version",
			Action: func(c *cli.Context) error {
				fmt.Println(getVersion())
				return nil
			},
		},
		{
			Name:   "upgrade",
			Usage:  "kratos self-upgrade",
			Action: upgradeAction,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
