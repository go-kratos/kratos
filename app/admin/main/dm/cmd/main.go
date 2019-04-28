package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/admin/main/dm/conf"
	"go-common/app/admin/main/dm/http"
	"go-common/app/admin/main/dm/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
	manager "go-common/library/queue/databus/report"
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
	// manager log init
	manager.InitManager(conf.Conf.ManagerLog)
	// service init
	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	log.Info("dm-admin start")
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("dm-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("dm-admin exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
