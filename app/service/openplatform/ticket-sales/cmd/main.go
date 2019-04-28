package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/openplatform/ticket-sales/conf"
	"go-common/app/service/openplatform/ticket-sales/server/grpc"
	"go-common/app/service/openplatform/ticket-sales/server/http"
	"go-common/app/service/openplatform/ticket-sales/service"
	"go-common/app/service/openplatform/ticket-sales/service/mis"
	"go-common/library/conf/paladin"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("ticket-sales.toml", conf.Conf); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	log.Info("ticket-sales start")
	ecode.Init(nil)
	// service init
	srv := service.New(conf.Conf)
	misSrv := mis.New(srv.Get())
	grpc.New(srv, misSrv)
	http.Init(conf.Conf, srv)
	log.Info("ready to serv")
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("ticket-sales get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			// rpcSvr.Close()
			time.Sleep(time.Second * 2)
			log.Info("ticket-sales exit")
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
