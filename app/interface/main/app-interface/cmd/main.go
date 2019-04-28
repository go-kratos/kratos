package main

import (
	"flag"
	"go-common/app/interface/main/app-interface/conf"
	"go-common/app/interface/main/app-interface/http"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
	"os"
	"os/signal"
	"syscall"
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
	log.Info("app-interface start")
	// init trace
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// ecode init
	ecode.Init(conf.Conf.Ecode)
	// service init
	http.Init(conf.Conf)
	// init pprof conf.Conf.Perf
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("app-interface get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("app-interface exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
