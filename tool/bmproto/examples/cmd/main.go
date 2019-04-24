package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bilibili/kratos/tool/bmproto/examples/internal/service"
	"github.com/bilibili/kratos/tool/bmproto/examples/internal/server/http"
	"github.com/bilibili/kratos/pkg/log"
)

func main() {
	flag.Parse()
	log.Init(&log.Config{
		Dir:"/tmp/",
	})
	svc := service.NewGreeterService()
	httpSvc := http.New(svc)
	defer log.Close()
	log.Info("service start")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			ctx, cancel := context.WithTimeout(context.Background(), 35*time.Second)
			defer cancel()
			httpSvc.Shutdown(ctx)
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
