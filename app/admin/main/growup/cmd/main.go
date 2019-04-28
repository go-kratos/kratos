package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/growup/conf"
	"go-common/app/admin/main/growup/http"
	"go-common/app/admin/main/growup/service"
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
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	log.SetFormat("[%D %T] [%L] [%S] %M")
	log.Info("growup-admin start")
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("growup-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("growup-admin exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
