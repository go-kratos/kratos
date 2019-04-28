package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/archive/conf"
	rpc "go-common/app/service/main/archive/server/gorpc"
	"go-common/app/service/main/archive/server/grpc"
	"go-common/app/service/main/archive/server/http"
	"go-common/app/service/main/archive/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	// init ecode
	ecode.Init(nil)
	// init log
	log.Init(conf.Conf.Xlog)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("archive-service start")
	// service init
	svr := service.New(conf.Conf)
	// statsd init
	rpcSvr := rpc.New(conf.Conf, svr)
	grpcSvr, err := grpc.New(nil, svr)
	if err != nil {
		panic(err)
	}
	http.Init(conf.Conf, svr)
	// init signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("archive-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			rpcSvr.Close()
			grpcSvr.Shutdown(context.TODO())
			time.Sleep(time.Second * 2)
			log.Info("archive-service exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
