package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/usersuit/conf"
	"go-common/app/admin/main/usersuit/http"
	"go-common/app/admin/main/usersuit/service"
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
	// service init
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	// service init
	ecode.Init(conf.Conf.Ecode)
	// signal handler
	log.Info("usersuit-admin start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("usersuit-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svr.Close()
			log.Info("usersuit-admin exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
