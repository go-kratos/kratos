package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/job/main/member/conf"
	"go-common/app/job/main/member/http"
	"go-common/app/job/main/member/service"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

var (
	svr *service.Service
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	report.InitUser(conf.Conf.UserReport)
	report.InitManager(conf.Conf.ManagerReport)
	log.Info("member-job start")
	svr = service.New(conf.Conf)
	http.Init(conf.Conf, svr)
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
			if err = svr.Close(); err != nil {
				log.Error("srv close consumer error(%v)", err)
			}
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
