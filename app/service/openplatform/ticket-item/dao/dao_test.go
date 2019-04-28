package dao

import (
	"time"

	"go-common/app/service/openplatform/ticket-item/conf"
	"go-common/library/log"
)

func initConf() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
}

func startService() {
	initConf()
	d = New(conf.Conf)
	time.Sleep(time.Second * 2)
}
