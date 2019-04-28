package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/relation/conf"
	"go-common/app/job/main/relation/http"
	"go-common/app/job/main/relation/service"
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
	log.Info("go-common/app/job/main/relation start")
	srv = service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	signalHandler()
}

func signalHandler() {
	var (
		err error
		ch  = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			if err = srv.Close(); err != nil {
				log.Error("srv close consumer error(%v)", err)
			}
			srv.Wait()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
