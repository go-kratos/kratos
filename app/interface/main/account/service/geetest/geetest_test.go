package geetest

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/account/conf"
	"go-common/app/interface/main/account/model"

	"github.com/smartystreets/goconvey/convey"
)

var (
	s *Service
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
		flag.Set("conf", "../cmd/account-interface-example.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func TestGeetestPreProcess(t *testing.T) {
	convey.Convey("PreProcess", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.GeeCaptchaRequest{
				MID: 1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.PreProcess(c, req)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
				ctx.Printf("%+v", res)
			})
		})
	})
}

func TestGeetestValidate(t *testing.T) {
	convey.Convey("Validate", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			req = &model.GeeCheckRequest{
				Challenge: "078348b20cda8680bd2a01ac79394c37",
				Validate:  "",
				Seccode:   "",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			stat := s.Validate(c, req)
			ctx.Convey("Then stat should not be nil.", func(ctx convey.C) {
				ctx.So(stat, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGeetestfailbackValidate(t *testing.T) {
	convey.Convey("failbackValidate", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			challenge = ""
			validate  = ""
			seccode   = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := s.failbackValidate(c, challenge, validate, seccode)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGeetestdecodeResponse(t *testing.T) {
	convey.Convey("decodeResponse", t, func(ctx convey.C) {
		var (
			challenge    = ""
			userresponse = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res := s.decodeResponse(challenge, userresponse)
			ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGeetestdecodeRandBase(t *testing.T) {
	convey.Convey("decodeRandBase", t, func(ctx convey.C) {
		var (
			challenge = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := s.decodeRandBase(challenge)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGeetestmd5Encode(t *testing.T) {
	convey.Convey("md5Encode", t, func(ctx convey.C) {
		var (
			values = []byte("")
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := s.md5Encode(values)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGeetestvalidateFailImage(t *testing.T) {
	convey.Convey("validateFailImage", t, func(ctx convey.C) {
		var (
			ans         = int(0)
			fullBgIndex = int(0)
			imgGrpIndex = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := s.validateFailImage(ans, fullBgIndex, imgGrpIndex)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
