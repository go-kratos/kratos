package service

import (
	"flag"

	"go-common/app/job/main/passport-game-data/conf"
	"sync"
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
	// service init
	s = New(conf.Conf)
}
