package main

import (
	"flag"
	"log"
	"os"

	"go-common/app/service/main/dapper/agent"
	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/pkg/util"
	xlog "go-common/library/log"
)

var debug bool

func init() {
	flag.BoolVar(&debug, "debug", false, "debug model decode and print span on stdout")
}

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	cfg, err := conf.LoadAgentConfig()
	if err != nil {
		log.Fatalf("local agent config error: %s", err)
	}
	xlog.Init(cfg.Log)
	defer xlog.Close()
	ag, err := agent.New(cfg, debug)
	if err != nil {
		log.Fatalf("new agent service error: %s", err)
	}
	util.HandlerReload(func(s os.Signal) {
		xlog.Warn("receive signal %s, dapper agent reload config", s)
		cfg, err := conf.LoadAgentConfig()
		if err != nil {
			xlog.Error("load config error: %s, reload config fail!", err)
			return
		}
		if err := ag.Reload(cfg); err != nil {
			xlog.Error("reload config error: %s", err)
		}
	})
	util.HandlerExit(func(s os.Signal) int {
		if err := ag.Close(); err != nil {
			xlog.Error("close agent error: %s", err)
			return 1
		}
		return 0
	})
}
