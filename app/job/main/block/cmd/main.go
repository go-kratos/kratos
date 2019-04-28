package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/block/conf"
	"go-common/app/job/main/block/http"
	"go-common/library/log"
	manager "go-common/library/queue/databus/report"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() err(%+v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	http.Init()
	// manager log init
	manager.InitManager(conf.Conf.ManagerLog)
	log.Info("block-job start")

	signalHandler()
}

func signalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("block get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("block-job exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
