package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/http"
	"go-common/app/admin/ep/saga/server/grpc"
	"go-common/app/admin/ep/saga/service"
	"go-common/library/log"
)

const (
	_durationForClosingServer = 2 // second
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}

	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("saga-admin start")

	s := service.New()
	http.Init(s)
	grpcsvr, err := grpc.New(nil, s.Wechat())
	if err != nil {
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-c
		log.Info("saga-admin get a signal %s", si.String())
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			grpcsvr.Shutdown(context.Background())
			log.Info("saga-admin exit")
			s.Close()
			time.Sleep(_durationForClosingServer * time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
