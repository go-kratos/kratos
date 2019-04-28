package dao

import (
	"context"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"go-common/app/job/main/push/conf"

	. "github.com/smartystreets/goconvey/convey"
)

var d *Dao

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") == "uat" {
		flag.Set("app_id", "main.web-svr.push-job")
		flag.Set("conf_token", "4de43ccf842485eea314fd8a48f1ee84")
		flag.Set("tree_id", "5220")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		dir, _ := filepath.Abs("../cmd/push-job-test.toml")
		flag.Set("conf", dir)
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func Test_Wechat(t *testing.T) {
	Convey("test wechat message", t, func() {
		err := d.SendWechat("test send wechat message")
		So(err, ShouldBeNil)
	})
}

func Test_DpDownloadFile(t *testing.T) {
	Convey("data platform download file", t, func() {
		_, err := d.DpDownloadFile(context.Background(), "https://raw.githubusercontent.com/Bilibili/discovery/master/README.md")
		So(err, ShouldBeNil)
	})
}

func Test_DpSubmitQuery(t *testing.T) {
	Convey("data platform submit query", t, func() {
		url, err := d.DpSubmitQuery(context.Background(), "select device_token from basic.dws_push_buvid where log_date='20180707'")
		So(err, ShouldNotBeNil)
		t.Logf("url(%v)", url)
	})
}
