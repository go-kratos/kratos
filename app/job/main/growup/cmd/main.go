package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/growup/conf"
	"go-common/app/job/main/growup/http"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	log.Info("growup-job start")
	http.Init(conf.Conf)
	signalHandler()
}

func signalHandler() {
	var ch = make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		log.Info("growup-job get a signal %s", si.String())
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			log.Info("growup-job [verison: %s] exit")
			http.Close()
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
