package main

import (
	"flag"
	"log"
	"os"

	"go-common/app/service/main/dapper/collector"
	"go-common/app/service/main/dapper/conf"
	"go-common/app/service/main/dapper/pkg/util"
	xlog "go-common/library/log"
)

func main() {
	if !flag.Parsed() {
		flag.Parse()
	}
	// load config file
	if err := conf.Init(); err != nil {
		log.Fatalf("init config error: %s", err)
	}
	// init xlog
	xlog.Init(conf.Conf.Log)
	defer xlog.Close()
	xlog.Info("dapper-service starting")

	// new collector service
	srv, err := collector.New(conf.Conf)
	if err != nil {
		log.Fatalf("new dapper service error: %s", err)
	}
	if err := srv.ListenAndStart(); err != nil {
		log.Fatalf("start dapper service error: %s", err)
	}
	//hsvr := http.New(srv)
	//if err := hsvr.Start(); err != nil {
	//	log.Fatalf("start dapper http server error: %s", err)
	//}
	util.HandlerExit(func(s os.Signal) int {
		xlog.Info("dapper-service get a signal %s", s.String())
		if err := srv.Close(); err != nil {
			xlog.Info("dapper-service exit, error: %s", err)
			return 1
		}
		xlog.Info("dapper-service exit")
		return 0
	})
}
