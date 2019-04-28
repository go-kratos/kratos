package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/dm2/conf"
	"go-common/library/log"
)

var (
	testDao *Dao
	c       = context.TODO()
)

func TestMain(m *testing.M) {
	flag.Set("app_id", "main.community.dm2")
	flag.Set("conf_token", "41bebe87fc9df75601583db2672f7c34")
	flag.Set("tree_id", "2291")
	flag.Set("conf_version", "docker-1")
	flag.Set("deploy_env", "uat")
	flag.Set("conf_host", "config.bilibili.co")
	flag.Set("conf_path", "/tmp")
	flag.Set("region", "sh")
	flag.Set("zone", "sh001")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Xlog)
	defer log.Close()
	testDao = New(conf.Conf)
	os.Exit(m.Run())
}
