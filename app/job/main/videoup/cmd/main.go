package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/library/log"

	"go-common/app/job/main/videoup/conf"
	"go-common/app/job/main/videoup/http"
	"go-common/app/job/main/videoup/service"
)

var (
	srv *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		log.Info("conf.Init() error(%v)", err)
		panic(err)
	}
	// init log
	log.Init(conf.Conf.XLog)
	defer log.Close()

	srv = service.New(conf.Conf)
	http.Init(conf.Conf, srv)
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
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
