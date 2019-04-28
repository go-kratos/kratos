package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/admin/main/open/conf"
	"go-common/app/admin/main/open/http"
	"go-common/app/admin/main/open/service"
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
	defer func() {
		log.Close()
		// wait for a while to guarantee that all log messages are written
		time.Sleep(10 * time.Millisecond)
	}()
	//service init
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("open-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			svc.Close()
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
