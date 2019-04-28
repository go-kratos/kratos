package dao

import (
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/answer/conf"
)

var (
	d   *Dao
	mid int64 = 7593623
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account-law.answer")
		flag.Set("conf_appid", "main.account-law.answer")
		flag.Set("conf_token", "ba3ee255695e8d7b46782268ddc9c8a3")
		flag.Set("tree_id", "25260")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_env", "10")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/answer-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
