package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/thumbup/conf"
	"go-common/app/job/main/thumbup/server/http"
	"go-common/app/job/main/thumbup/service"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	srv := service.New(conf.Conf)
	http.Init(conf.Conf, srv)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			srv.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
