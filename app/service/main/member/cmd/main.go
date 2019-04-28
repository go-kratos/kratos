package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/member/conf"
	"go-common/app/service/main/member/server/gorpc"
	"go-common/app/service/main/member/server/grpc"
	"go-common/app/service/main/member/server/http"
	"go-common/app/service/main/member/service"
	"go-common/library/log"
	"go-common/library/net/trace"
	"go-common/library/queue/databus/report"
)

func main() {
	flag.Parse()
	// init conf,log,trace,stat,perf.
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	report.InitUser(conf.Conf.ReportUser)
	report.InitManager(conf.Conf.ReportManager)
	svr := service.New(conf.Conf)
	rpcSvr := gorpc.New(conf.Conf, svr)
	ws := grpc.New(conf.Conf.WardenServer, svr)
	http.Init(conf.Conf, svr)
	// signal handler
	log.Info("member-service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("member-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ws.Shutdown(context.Background())
			rpcSvr.Close()
			time.Sleep(time.Second * 2)
			log.Info("member-service exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
