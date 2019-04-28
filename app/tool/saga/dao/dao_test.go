package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/tool/saga/conf"
)

var (
	d   *Dao
	ctx = context.Background()
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "test.ep.saga")
		flag.Set("conf_token", "1d10044ab80b21d37ea67445f0475237")
		flag.Set("tree_id", "56733")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "fat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/saga-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New()
	defer d.Close()
	os.Exit(m.Run())
}
