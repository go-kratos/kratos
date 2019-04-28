package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	apiv1 "go-common/app/service/main/dapper-query/api/v1"
	"go-common/app/service/main/dapper-query/conf"
	"go-common/app/service/main/dapper-query/service"
	"go-common/app/service/main/dapper-query/util"
	xlog "go-common/library/log"
	bm "go-common/library/net/http/blademaster"
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
	xlog.Init(nil)
	defer xlog.Close()
	xlog.Info("dapper-service start")

	// new dapper service
	srv, err := service.New(conf.Conf)
	if err != nil {
		log.Fatalf("new dapper service error: %s", err)
	}
	// init blademaster server
	engine := bm.NewServer(nil)
	engine.Use(bm.Recovery(), bm.Logger())
	engine.Ping(func(*bm.Context) {})
	engine.Inject("^/x/internal/dapper/ops-log", util.SessionIDMiddleware)
	apiv1.RegisterDapperQueryBMServer(engine, srv)
	if err := engine.Start(); err != nil {
		log.Fatalf("start bm server error: %s", err)
	}

	sch := make(chan os.Signal, 1)
	signal.Notify(sch, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	s := <-sch
	// program exit
	xlog.Info("dapper-service get a signal %s", s.String())
	xlog.Info("dapper-service exit")
}
