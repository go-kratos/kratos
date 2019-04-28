package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/feedback/conf"
	"go-common/app/interface/main/feedback/http"
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
	// init log
	log.Init(conf.Conf.XLog)
	defer log.Close()
	log.Info("feedback start")
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
		log.Info("feedback get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("feedback exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
