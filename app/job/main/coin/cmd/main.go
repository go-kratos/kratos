package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/coin/conf"
	"go-common/app/job/main/coin/http"
	"go-common/app/job/main/coin/service"
	"go-common/library/log"
)

var (
	s *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	log.Info("archive-service_consumer start")
	s = service.New(conf.Conf)
	http.Init(conf.Conf, s)
	signalHandler()
}

func signalHandler() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			s.Close()
			s.Wait()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
