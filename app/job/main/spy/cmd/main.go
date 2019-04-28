package main

import (
	"flag"
	"os"
	"time"

	"go-common/app/job/main/spy/conf"
	"go-common/app/job/main/spy/http"
	"go-common/app/job/main/spy/service"
	"go-common/library/log"
	"go-common/library/os/signal"
	"go-common/library/syscall"
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

	// service init
	s = service.New(conf.Conf)
	http.Init(conf.Conf, s)

	log.Info("spy-job start")
	signalHandler()
}

func signalHandler() {
	var (
		err error
		ch  = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("get a signal %s, stop the spy-job process", si.String())
			if err = s.Close(); err != nil {
				log.Error("close spy-job error(%v)", err)
			}
			s.Wait()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
