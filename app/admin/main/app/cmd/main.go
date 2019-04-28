package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/app/conf"
	"go-common/app/admin/main/app/http"
	"go-common/library/log"
	"go-common/library/net/trace"
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
	log.Info("app-admin start")
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// service init
	http.Init(conf.Conf)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("app-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("app-admin exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
