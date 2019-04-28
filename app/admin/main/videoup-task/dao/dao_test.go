package dao

import (
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/admin/main/videoup-task/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.videoup-task-admin")
		flag.Set("conf_token", "7e2256b4bc8e4084cf848fe27a474a20")
		flag.Set("tree_id", "25095")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "dev")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/videoup-task-admin.toml")
	}

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	d.hclient.SetTransport(gock.DefaultTransport)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
