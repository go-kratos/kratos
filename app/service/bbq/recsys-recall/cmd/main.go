package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"go-common/app/service/bbq/recsys-recall/conf"
	"go-common/app/service/bbq/recsys-recall/server/grpc"
	"go-common/app/service/bbq/recsys-recall/server/http"
	"go-common/app/service/bbq/recsys-recall/service"
	"go-common/app/service/bbq/recsys-recall/service/index"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	// cpuf, err := os.Create("cpu_profile")
	// if err != nil {
	// 	panic(err)
	// }
	// pprof.StartCPUProfile(cpuf)
	// defer pprof.StopCPUProfile()

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	srv := service.New(conf.Conf)

	// 加载正排索引
	index.Init(conf.Conf)

	grpc.New(srv)
	http.Init(conf.Conf)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
