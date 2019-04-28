package service

import (
	"go-common/app/job/main/identify/conf"
	"sync"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}
