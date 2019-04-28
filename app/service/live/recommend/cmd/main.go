package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/service/live/recommend/internal/conf"
	"go-common/app/service/live/recommend/internal/dao"
	"go-common/app/service/live/recommend/internal/server/grpc"
	"go-common/app/service/live/recommend/internal/server/http"
	"go-common/app/service/live/recommend/internal/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

var runJob = false
var jobFilePath = ""
var jobWorkerNum = 1000
var jobOffset = -1

func main() {
	flag.BoolVar(&runJob, "runJob", false, "跑redis脚本")
	flag.StringVar(&jobFilePath, "jobFile", "", "推荐文件地址")
	flag.IntVar(&jobOffset, "jobOffset", -1, "操作偏移，即从n行开始跑, 默认会从上一次中断的地方开始")
	flag.IntVar(&jobWorkerNum, "jobWorkerNum", 1000, "worker数量")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("recommend-service start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)

	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	grpc.Init(conf.Conf)
	go dao.StartRefreshJob()
	go dao.StartRoomFeatureJob(conf.Conf)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			log.Info("recommend-service exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
