package service

import (
	"context"
	"flag"
	"testing"

	"go-common/app/admin/main/aegis/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	s    *Service
	cntx context.Context
)

func init() {
	flag.Set("conf", "../cmd/aegis-admin.toml")
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	cntx = context.Background()
}

/*
func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.aegis-admin")
		flag.Set("conf_token", "cad913269be022e1eb8c45a8d5408d78")
		flag.Set("tree_id", "60977")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/local.toml")
	}

	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	conf.Conf.AegisPub = &databus.Config{
		// Key:          "0PvKGhAqDvsK7zitmS8t",
		// Secret:       "0PvKGhAqDvsK7zitmS8u",
		// Group:        "databus_test_group",
		// Topic:        "databus_test_topic",
		Key:    "dbe67e6a4c36f877",
		Secret: "8c775ea242caa367ba5c876c04576571",
		Group:  "Test1-MainCommonArch-P",
		Topic:  "test1",
		Action: "pub",
		Name:   "databus",
		Proto:  "tcp",
		// Addr:         "172.16.33.158:6205",
		Addr:         "172.18.33.50:6205",
		Active:       10,
		Idle:         5,
		DialTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		IdleTimeout:  xtime.Duration(time.Minute),
	}
	s = New(conf.Conf)
	cntx = context.Background()
	m.Run()
	os.Exit(0)
}
*/

func TestServicePing(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := s.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestServiceClose(t *testing.T) {
	convey.Convey("Close", t, func(ctx convey.C) {
		//s.Close()
		//ctx.Convey("No return values", func(ctx convey.C) {
		//})
	})
}
