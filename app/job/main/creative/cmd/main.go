package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/creative/conf"
	"go-common/app/job/main/creative/http"
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
	log.Info("creative-jobstart")
	// service init
	http.Init(conf.Conf)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("creative-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			http.Svc.Close()
			log.Info("creative-jobexit")
			time.Sleep(1 * time.Second)
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
