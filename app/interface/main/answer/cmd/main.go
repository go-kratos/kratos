package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/answer/conf"
	"go-common/app/interface/main/answer/http"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
	"go-common/library/text/translate/chinese"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	ecode.Init(conf.Conf.Ecode)
	chinese.Init()
	http.Init(conf.Conf)
	report.InitUser(conf.Conf.Report)
	// signal handler
	log.Info("answer-interface start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("answer get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			log.Info("answer-interface exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
