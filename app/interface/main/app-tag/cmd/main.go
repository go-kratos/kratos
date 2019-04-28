package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/app-tag/conf"
	"go-common/app/interface/main/app-tag/http"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init ecode
	ecode.Init(nil)
	// init log
	log.Init(conf.Conf.XLog)
	defer log.Close()
	log.Info("app-tag start")
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
		log.Info("app-tag get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("app-tag exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
