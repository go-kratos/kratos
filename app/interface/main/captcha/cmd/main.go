package main

import (
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/interface/main/captcha/conf"
	"go-common/app/interface/main/captcha/http"
	"go-common/app/interface/main/captcha/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	log.Init(conf.Conf.XLog)
	defer log.Close()
	log.Info("captcha-service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	svr := service.New(conf.Conf)
	ecode.Init(conf.Conf.Ecode)
	http.Init(conf.Conf, svr)
	// signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("captcha-service get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			// svr.Close()
			log.Info("captcha-service exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
