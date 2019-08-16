package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "protc"
	app.Usage = "protobuf生成工具"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "bm",
			Usage:       "whether to use BM for generation",
			Destination: &withBM,
		},
		cli.BoolFlag{
			Name:        "grpc",
			Usage:       "whether to use gRPC for generation",
			Destination: &withGRPC,
		},
		cli.BoolFlag{
			Name:        "swagger",
			Usage:       "whether to use swagger for generation",
			Destination: &withSwagger,
		},
		cli.BoolFlag{
			Name:        "ecode",
			Usage:       "whether to use ecode for generation",
			Destination: &withEcode,
		},
	}
	app.Action = func(c *cli.Context) error {
		return protocAction(c)
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
