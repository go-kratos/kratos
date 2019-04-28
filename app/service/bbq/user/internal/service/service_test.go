package service

import (
	"flag"
	"go-common/app/service/bbq/user/internal/conf"
	"go-common/library/conf/paladin"
	"go-common/library/log"
	"os"
	"testing"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "")
		flag.Set("conf_token", "")
		flag.Set("tree_id", "")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/test.toml")
		flag.Set("deploy.env", "uat")
	}
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	if err := paladin.Watch("test.toml", conf.Conf); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	s = New(conf.Conf)
	os.Exit(m.Run())
}
