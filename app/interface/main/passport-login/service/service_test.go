package service

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"

	"go-common/app/interface/main/passport-login/conf"
)

var (
	once sync.Once
	s    *Service
)

func startService() {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.passport.passport-login-interface")
		flag.Set("conf_token", "edc2d1ae5a49fc907eb173745e030264")
		flag.Set("tree_id", "62792")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/passport-login-interface.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	s = New(conf.Conf)
}

func TestNew(t *testing.T) {
	once.Do(startService)
}
