package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/interface/main/favorite/conf"
	"go-common/app/interface/main/favorite/http"
	"go-common/app/interface/main/favorite/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

var svc *service.Service

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
	log.Info("favorite start")
	ecode.Init(conf.Conf.Ecode)
	// service init
	svc = service.New(conf.Conf)
	http.Init(conf.Conf, svc)
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
			svc.Close()
			if err != nil {
				log.Error("svc.Close() error(%v)", err)
			}
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
