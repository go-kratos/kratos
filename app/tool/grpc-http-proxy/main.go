package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	xtime "go-common/library/time"
)

var bind = flag.String("bind", "localhost:10010", "bind addr example: localhost:10010")

func main() {
	flag.Parse()
	engineOuter := bm.DefaultServer(&bm.ServerConfig{
		Addr:    *bind,
		Timeout: xtime.Duration(time.Second),
	})
	// bm不支持 get post 绑定一个路径 为了获取body只能用post
	engineOuter.POST("/", handle)
	if err := engineOuter.Start(); err != nil {
		log.Error("engine.Start error(%v)", err)
		panic(err)
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Info("exit")
			return
		default:
			return
		}
	}
}
