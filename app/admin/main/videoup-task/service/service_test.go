package service

import (
	"flag"

	"go-common/app/admin/main/videoup-task/conf"
)

var s *Service

func Init() {
	if s != nil {
		return
	}

	flag.Set("conf", "../cmd/videoup-task-admin.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}

	s = New(conf.Conf)
}
