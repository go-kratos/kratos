package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/credit/conf"
	"go-common/app/job/main/credit/http"
	"go-common/app/job/main/credit/service"
	"go-common/library/log"
)

var (
	srv *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	log.Info("credit-job start")
	http.Init(conf.Conf)
	signalHandler()
}

func signalHandler() {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			srv.Close()
			time.Sleep(2 * time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
