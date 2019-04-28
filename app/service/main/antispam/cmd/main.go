package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/main/antispam/conf"
	"go-common/app/service/main/antispam/http"
	rpc "go-common/app/service/main/antispam/rpc/server"
	"go-common/app/service/main/antispam/service"
	ecode "go-common/library/ecode/tip"

	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	flag.Parse()
	if err := conf.Init(conf.ConfPath); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	ecode.Init(conf.Conf.Ecode)
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	defer log.Close()
	log.Info("antispam start")
	svr := service.New(conf.Conf)
	rpcSvr := rpc.New(conf.Conf, svr)
	http.Init(conf.Conf, svr)
	// init signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("get a signal %s, stop the consume process", si.String())
			rpcSvr.Close()
			time.Sleep(time.Second * 2)
			svr.Close()
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
