package main

import (
	"flag"
	"os"
	"time"

	"go-common/app/admin/main/spy/conf"
	"go-common/app/admin/main/spy/http"
	"go-common/app/admin/main/spy/service"
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
	log.Init(conf.Conf.Log)
	defer log.Close()
	// service init
	http.Init(conf.Conf)
	log.Info("spy-admin start")
	signalHandler()
}

func signalHandler() {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGSTOP)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			log.Info("get a signal %s, stop the spy-admin process", si.String())
			s.Wait()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
