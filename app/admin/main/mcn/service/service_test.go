package service

import (
	"flag"
	"os"

	"go-common/app/admin/main/mcn/conf"
)

var (
	s *Service
)

func init() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.mcn-admin")
		flag.Set("conf_token", "BVWgBtBvS2pkTBbmxAl0frX6KRA14d5P")
		flag.Set("tree_id", "6813")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/mcn-admin-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
}
