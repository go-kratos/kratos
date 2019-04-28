package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/admin/main/macross/conf"
	"go-common/app/admin/main/macross/http"
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
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("macross start")
	// service init
	http.Init(conf.Conf)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("macross get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("macross exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
