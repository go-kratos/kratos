package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/admin/main/videoup/conf"
	"go-common/app/admin/main/videoup/http"
	"go-common/app/admin/main/videoup/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/queue/databus/report"
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
	defer log.Close()
	log.Info("videoup-admin start")
	report.InitManager(conf.Conf.ManagerReport)
	ecode.Init(conf.Conf.Ecode)
	// service init
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("videoup-admin get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			svr.Close()
			log.Info("videoup-admin exit")
			time.Sleep(1 * time.Second)
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
