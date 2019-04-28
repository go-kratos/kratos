package pgc

import (
	"context"
	"flag"
	"os"
	"strings"

	"go-common/app/interface/main/tv/conf"

	"gopkg.in/h2non/gock.v1"
)

var (
	d   *Dao
	ctx = context.Background()
)

func init() {
	// dir, _ := filepath.Abs("../../cmd/tv-interface.toml")
	// flag.Set("conf", dir)
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.web-svr.tv-interface")
		flag.Set("conf_token", "07c1826c1f39df02a1411cdd6f455879")
		flag.Set("tree_id", "15326")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}
