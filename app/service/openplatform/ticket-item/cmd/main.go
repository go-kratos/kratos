package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/openplatform/ticket-item/conf"
	rpc "go-common/app/service/openplatform/ticket-item/server/grpc"
	"go-common/app/service/openplatform/ticket-item/server/http"
	"go-common/app/service/openplatform/ticket-item/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	trace.Init(nil)
	defer trace.Close()
	log.Info("ticket-item start")
	svr := service.New(conf.Conf)
	// http Service init
	http.Init(conf.Conf, svr)
	// grpc service init
	gsvr := rpc.New(conf.Conf)
	// init pprof conf.Conf.Perf
	// signal handler
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		log.Info("ticket-item get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT:
			svr.Close()
			gsvr.Shutdown(ctx)
			log.Info("ticket-item exit")
			time.Sleep(time.Second) // 休眠10s
			return
		case syscall.SIGHUP:
		// TODO reload
		default:
			return
		}
	}
}
