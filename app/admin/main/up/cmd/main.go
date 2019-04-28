package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/admin/main/up/conf"
	"go-common/app/admin/main/up/http"
	"go-common/library/log"
	manager "go-common/library/queue/databus/report"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.XLog)
	defer log.Close()
	log.SetFormat("[%D %T] [%L] [%S] %M")
	log.Info("up-adminstart")
	// service init
	http.Init(conf.Conf)
	//logCli.Init(conf.Conf.LogCli)
	// manager log init
	manager.InitManager(nil)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("up-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			http.Svc.Close()
			log.Info("up-adminexit")
			time.Sleep(1 * time.Second)
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
