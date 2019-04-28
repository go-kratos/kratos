package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/live/recommend-job/internal/conf"
	"go-common/app/job/live/recommend-job/internal/server/http"
	"go-common/app/job/live/recommend-job/internal/service"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/net/trace"
)

var itemCFRun = false
var itemCFInputPath = ""
var itemCFWorkerNum = 1000

var userAreaRun = false
var uerAreaInputPath = ""

func main() {
	flag.BoolVar(&userAreaRun, "userarea.run", false, "用户分区缓存：是否手动跑")
	flag.StringVar(&uerAreaInputPath, "userarea.input", "", "用户分区缓存：输入文件地址")

	flag.BoolVar(&itemCFRun, "itemcf.run", false, "协同过滤推荐缓存到redis：是否手动跑")
	flag.StringVar(&itemCFInputPath, "itemcf.input", "", "协同过滤结果文件地址")
	flag.IntVar(&itemCFWorkerNum, "itemcf.workerNum", 1000, "worker数量")

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
	log.Info("recommend-job start")
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()
	ecode.Init(conf.Conf.Ecode)

	if itemCFRun {
		c := conf.Conf.ItemCFJob
		c.InputFile = itemCFInputPath
		c.WorkerNum = itemCFWorkerNum
		var job = service.ItemCFJob{
			Conf:       c,
			RedisConf:  conf.Conf.Redis,
			HadoopConf: conf.Conf.Hadoop,
		}
		job.Run()
		os.Exit(0)
	}

	if userAreaRun {
		c := conf.Conf.UserAreaJob
		c.InputFile = uerAreaInputPath
		c.WorkerNum = 1000
		var job = service.UserAreaJob{
			JobConf:    c,
			RedisConf:  conf.Conf.Redis,
			HadoopConf: conf.Conf.Hadoop,
		}
		job.Run()
		os.Exit(0)
	}

	svc := service.New(conf.Conf)
	http.Init(conf.Conf, svc)
	svc.RunCrontab()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			svc.Close()
			log.Info("recommend-job exit")
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}
