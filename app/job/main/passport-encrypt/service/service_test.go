package service

import (
	"flag"
	"sync"

	"go-common/app/job/main/passport-encrypt/conf"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}
