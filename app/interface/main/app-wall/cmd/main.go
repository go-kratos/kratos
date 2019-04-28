package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/app-wall/conf"
	"go-common/app/interface/main/app-wall/http"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
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
	log.Info("app-wall start")
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// ecode init
	ecode.Init(conf.Conf.Ecode)
	// report init
	report.InitUser(conf.Conf.Report)
	// service init
	http.Init(conf.Conf)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("app-wall get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("app-wall exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
