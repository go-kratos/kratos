package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/dm/conf"
	"go-common/app/interface/main/dm/http"
	"go-common/app/interface/main/dm/service"

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
	log.Init(conf.Conf.Xlog)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer func() {
		log.Close()
		// wait for a while to guarantee that all log messages are written
		time.Sleep(10 * time.Millisecond)
	}()
	// ecode init
	ecode.Init(conf.Conf.Ecode)
	// service init
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	log.Info("dm start")
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("dm get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("dm exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
