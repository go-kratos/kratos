package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/tv/conf"
	"go-common/app/interface/main/tv/http"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	// init conf,log,trace,stat,perf.
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	http.Init(conf.Conf)
	// signal handler
	log.Info("tv-interface start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("tv-interface get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("tv-interface exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
