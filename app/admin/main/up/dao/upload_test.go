package dao

import (
	"context"
	"flag"
	"io"
	"os"
	"testing"

	"go-common/app/admin/main/up/conf"

	"github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.up-admin")
		flag.Set("conf_token", "930697bb7def4df0713ef8080596b863")
		flag.Set("tree_id", "36438")
		flag.Set("conf_version", "1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../cmd/up-admin.toml")
	}
	if os.Getenv("UT_LOCAL_TEST") != "" {
		flag.Set("conf", "../cmd/up-admin.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestDaoUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = ""
			fileType = ""
			expire   = int64(0)
			body     io.Reader
			bfsConf  = &conf.Bfs{}
		)
		//defer gock.OffAll()
		//httpMock("PUT", bfsConf.Addr+uploadurl).Reply(200)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			location, err := Upload(c, fileName, fileType, expire, body, bfsConf)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
				ctx.So(location, convey.ShouldEqual, "")
			})
		})
	})
}

func TestDaoAuthorize(t *testing.T) {
	convey.Convey("Authorize", t, func(ctx convey.C) {
		var (
			key    = ""
			secret = ""
			method = ""
			bucket = ""
			file   = ""
			expire = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			authorization := Authorize(key, secret, method, bucket, file, expire)
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(authorization, convey.ShouldNotBeNil)
			})
		})
	})
}
