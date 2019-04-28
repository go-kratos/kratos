package main

import (
	"flag"
	"go-common/library/net/trace"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/live/push-search/conf"
	"go-common/app/job/live/push-search/http"
	"go-common/app/job/live/push-search/service/migrate"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("push-search-service start")
	if len(os.Args[1:]) > 0 && os.Args[1:][0] == "migrate" {
		roomId := os.Args[1:][1]
		isTest := os.Args[1:][2]
		ms := migrate.NewMigrateS(conf.Conf)
		go ms.Migrate(roomId, isTest)
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
		for {
			m := <-c
			switch m {
			case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
				ms.Close()
				log.Info("push-search-service-migrate exit")
				time.Sleep(time.Second)
				return
			case syscall.SIGHUP:
			default:
				return
			}
		}
	}
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)
	http.Init(conf.Conf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			http.Srv.Close()
			log.Info("push-search-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
