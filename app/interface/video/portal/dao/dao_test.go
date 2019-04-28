package dao

import (
	"flag"
	"go-common/app/interface/video/portal/conf"
	"os"
	"testing"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "video.live.portal")
		flag.Set("conf_token", "b112fc03806e30a1d20bae5da907a711")
		flag.Set("tree_id", "60575")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "prod")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/portal.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}
