package service

import (
	"flag"
	"go-common/app/admin/main/laser/conf"
	"os"
)

var (
	s *Service
)

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "mobile.studio.laser-admin")
		flag.Set("conf_token", "25911b439f4636ce9083f91c4882dffa")
		flag.Set("tree_id", "19167")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/laser-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}
