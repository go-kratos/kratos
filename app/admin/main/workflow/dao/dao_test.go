package dao

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/library/conf/paladin"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") == "uat" {
		flag.Set("app_id", "main.manager.workflow-admin")
		flag.Set("conf_token", "daebcaa3b1886d74e1b8f0361b34b04c")
		flag.Set("tree_id", "6812")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
		flag.Set("app_id", "main.manager.workflow-admin")
		flag.Set("deploy.env", "dev")
		flag.Set("conf_token", "1de17252107b89394b5e72e07bfbc8de")
		flag.Set("conf_path", "/tmp")
		flag.Set("conf_version", "docker-1")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("tree_id", "6812")
	}
	flag.Parse()
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	d = New()
	os.Exit(m.Run())
}

func TestPing(t *testing.T) {
	convey.Convey("Ping", t, func() {
		d.Ping(context.TODO())
	})
}
