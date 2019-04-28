package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/upload/conf"
	"go-common/app/job/main/upload/http"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("upload start")
	// service init
	http.Init(conf.Conf)
	// init pprof conf.Conf.Perf
	// report
	report.InitManager(nil)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("upload get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("upload exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
