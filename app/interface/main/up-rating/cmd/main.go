package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/up-rating/conf"
	"go-common/app/interface/main/up-rating/http"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	// init conf,log,trace,stat,perf
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	http.Init(conf.Conf)
	log.Info("up-rating-interface start")
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("up-rating-interface get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("up-rating-interface exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
