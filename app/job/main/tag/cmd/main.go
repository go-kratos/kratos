package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/tag/conf"
	"go-common/app/job/main/tag/http"
	"go-common/app/job/main/tag/service"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("tag-job start")
	svr := service.New(conf.Conf)
	http.Init(conf.Conf, svr)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-ch
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("tag-job stopped, get a signal %s, ", s.String())
			if err := svr.Close(); err != nil {
				log.Error("close consumer error(%v)", err)
			}
			svr.Wait()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
