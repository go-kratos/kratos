package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/vip/conf"
	"go-common/app/job/main/vip/http"
	"go-common/app/job/main/vip/service"
	"go-common/library/log"
)

var (
	s *service.Service
)

func main() {

	flag.Parse()
	// init conf,log,trace,stat,perf.
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	s = service.New(conf.Conf)
	http.Init(conf.Conf, s)
	// rpcSvr := rpc.New(conf.Conf, svr)
	// signal handler
	log.Info("vip-job start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		signal := <-c
		log.Info("vip-job get a signal %s", signal.String())
		switch signal {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Close()
			time.Sleep(time.Second * 2)
			log.Info("vip-job exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
