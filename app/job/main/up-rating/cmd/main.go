package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/up-rating/conf"
	"go-common/app/job/main/up-rating/http"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)

	http.Init(conf.Conf)
	signalHandler()
}

func signalHandler() {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		log.Info("up-rating-job get a signal %s", si.String())
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			http.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
