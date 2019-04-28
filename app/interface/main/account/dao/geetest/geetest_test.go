package geetest

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/account/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.account.account-interface")
		flag.Set("conf_token", "967eef77ad40b478234f11b0d489d6d6")
		flag.Set("tree_id", "3815")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/account-interface-example.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestGeetestPreProcess(t *testing.T) {
	convey.Convey("PreProcess", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			mid        = int64(11)
			geeType    = ""
			gc         = d.c.Geetest.PC
			newCaptcha = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			challenge, err := d.PreProcess(c, mid, geeType, gc, newCaptcha)
			ctx.Convey("Then err should be nil.challenge should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(challenge, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGeetestValidate(t *testing.T) {
	convey.Convey("Validate", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			challenge  = "ac7e67bb9b8e3e567186159e63df2153cs"
			seccode    = "7bda121544bd035032b2b5e390bce5d7"
			clientType = "web"
			captchaID  = "0926f406a4a5433a45faf33ae9a721c3"
			mid        = int64(11)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Validate(c, challenge, seccode, clientType, captchaID, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
