package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-common/app/job/main/passport-game-data/conf"
	"go-common/app/job/main/passport-game-data/http"
	"go-common/app/job/main/passport-game-data/service"
	"go-common/library/log"
	"go-common/library/net/trace"
)

var (
	mode int

	compareMidListFile string
	diffLogFile        string
	diffParseResFile   string
)

const (
	_modeNormal       = 0
	_modeCompareOnly  = 1
	_modeParseDiffLog = 2
	_modeInitCloud    = 3
)

func init() {
	flag.IntVar(&mode, "mode", _modeNormal, "mode for starting this job, 0 for normal, 1 for compare only, 2 for parse diff log")

	flag.StringVar(&compareMidListFile, "compare_mid_list_file", "/tmp/mids.txt", "compare mid list file path")

	flag.StringVar(&diffLogFile, "diff_log_file", "/tmp/diff.txt", "diff log file path")

	flag.StringVar(&diffParseResFile, "diff_parse_res_file", "/tmp/diff_parse_res.txt", "diff parse result file path")
}

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	trace.Init(conf.Conf.Tracer)
	defer trace.Close()

	var srv *service.Service

	switch mode {
	case _modeCompareOnly:
		s := service.NewCompareOnly(conf.Conf)
		if err := s.CompareFromMidListFile(context.TODO(), compareMidListFile); err != nil {
			log.Error("service.CompareFromMidListFile(%s) error(%v)", compareMidListFile, err)
		}
	case _modeParseDiffLog:
		log.Info("parse diff log from %s", diffLogFile)
		if err := service.ParseDiffLog(diffLogFile, diffParseResFile); err != nil {
			log.Error("service.ParseDiffLog(%s) error(%v)", diffLogFile, err)
		}
		time.Sleep(time.Second * 2)
		return
	case _modeInitCloud:
		s := service.NewInitCloud(conf.Conf)
		s.InitCloud(context.TODO())
		time.Sleep(time.Second * 2)
	case _modeNormal:
		// service init
		srv = service.New(conf.Conf)
		http.Init(conf.Conf, srv)
		// signal handler
		log.Info("passport-game-data-job start")
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Info("passport-game-data-job get a signal %s", s.String())
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			if srv != nil {
				srv.Close()
			}
			log.Info("passport-game-data-job exit")
			return
		case syscall.SIGHUP:
			// TODO reload
		default:
			return
		}
	}
}
