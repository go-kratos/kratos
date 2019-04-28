package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/favorite/conf"
	"go-common/app/job/main/favorite/http"
	"go-common/app/job/main/favorite/service"
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
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("archive-service_consumer start")
	s = service.New(conf.Conf)
	http.Init(conf.Conf, s)
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
			if err = s.Close(); err != nil {
				log.Error("close consumer error(%v)", err)
			}
			s.Wait()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
