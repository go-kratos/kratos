package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/main/bns/agent"
	"go-common/app/service/main/bns/conf"
	"go-common/library/log"

	_ "go-common/app/service/main/bns/agent/backend/discovery"
)

var confPath string

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}

	cfg, err := conf.LoadConfig(confPath)
	if err != nil {
		panic(fmt.Sprintf("loadconfig from: %s error: %s", confPath, err))
	}

	log.Init(cfg.Log)
	defer log.Close()

	ag, err := agent.New(cfg)
	if err != nil {
		panic(err)
	}

	if err := ag.Start(); err != nil {
		panic(err)
	}

	log.Info("bns start ...")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for s := range ch {
		log.Info("bns get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("bns exit ...")
			return
		case syscall.SIGHUP:
			log.Warn("reload is not support yet!")
		default:
			os.Exit(1)
		}
	}
}
