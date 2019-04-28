package service

import (
	"flag"
	"os"
	"testing"

	"go-common/app/job/main/app/conf"
)

var (
	s *Service
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-job")
		flag.Set("conf_token", "613aae0ddd1cc47a79920d6115cea472")
		flag.Set("tree_id", "2861")
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
	s = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestArchiveArchive3(t *testing.T) {
	// var (
	// 	c   = context.TODO()
	// 	aid = int64(-1)
	// )
	// convey.Convey("Archive3", t, func(ctx convey.C) {
	// 	_, err := d.Archive3(c, aid)
	// 	ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
	// 		ctx.So(err, convey.ShouldNotBeNil)
	// 	})
	// })
}
