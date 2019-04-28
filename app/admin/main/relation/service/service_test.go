package service

import (
	"flag"
	"go-common/app/admin/main/relation/conf"
)

var s *Service

func init() {
	flag.Parse()

	if err := conf.Init(); err != nil {
		panic(err)
	}

	s = New(conf.Conf)
}
