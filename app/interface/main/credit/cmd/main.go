package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/credit/conf"
	"go-common/app/interface/main/credit/http"
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
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	// service init
	ecode.Init(conf.Conf.Ecode)
	http.Init(conf.Conf)
	// signal handler
	log.Info("credit start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("credit get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			log.Info("credit exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
