package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/openplatform/open-market/conf"
	"go-common/app/job/openplatform/open-market/http"
	"go-common/app/job/openplatform/open-market/service"
	"go-common/library/log"
)

var (
	svr *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	// service init
	svr = service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	log.Info("open-market-job start")
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("open-market-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("open-market-job exit")
			if err := svr.Close(); err != nil {
				log.Error("srv close consumer error(%v)", err)
			}
			time.Sleep(2 * time.Second)
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
