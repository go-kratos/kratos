package dao

import (
	"context"
	"net/http"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"gopkg.in/h2non/gock.v1"
)

func TestDaoParseUTFiles(t *testing.T) {
	convey.Convey("ParseUTFiles", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			url = "http://uat-i0.hdslb.com/bfs/test/03af5d0707e4c6ec41a3eb797c8e9897f124515c.html"
		)
		gock.Off()
		d.client.SetTransport(http.DefaultClient.Transport)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			pkgs, err := d.ParseUTFiles(c, url)
			//t.Logf("the pkg[0] is PKG: %s, Coverage: %v, HTMLURL: %s", pkgs[0].PKG, pkgs[0].Coverage, pkgs[0].HTMLURL)
			ctx.Convey("Then err should be nil.pkgs should not be nil.", func(ctx convey.C) {
				ctx.So(pkgs, convey.ShouldNotBeNil)
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoGitCommitInfo(t *testing.T) {
	convey.Convey("GitCommitInfo", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			commitID = "ae1377033a11ca85a19bca365af32a5b0ebea31c"
		)
		gock.Off()
		d.client.SetTransport(http.DefaultClient.Transport)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			commit, err := d.GitLabCommits(c, commitID)
			ctx.Convey("Then err should be nil. gitcommit should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(commit, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSendWechatToUsers(t *testing.T) {
	convey.Convey("SendWechatToUsers", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			users = []string{"zhaobingqing"}
			msg   = "给你我的小心心❤"
		)
		d.client.SetTransport(gock.DefaultTransport)
		ctx.Convey("When everything gose postive", func(ctx convey.C) {
			httpMock("POST", _sagaWechatURL+"/message/send").Reply(200).JSON(`{"code":0,"message":"0"}`)
			err := d.SendWechatToUsers(c, users, msg)
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When http status != 200", func(ctx convey.C) {
			httpMock("POST", _sagaWechatURL+"/message/send").Reply(404)
			err := d.SendWechatToUsers(c, users, msg)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			httpMock("POST", _sagaWechatURL+"/message/send").Reply(200).JSON(`{"code":-401,"message":"0"}`)
			err := d.SendWechatToUsers(c, users, msg)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			gock.Off()
			d.client.SetTransport(http.DefaultClient.Transport)
		})
	})
}

func TestDaoSendWechatToGroup(t *testing.T) {
	convey.Convey("SendWechatToGroup", t, func(ctx convey.C) {
		var (
			msg = "信息测试~@002157"
			c   = context.Background()
		)
		d.client.SetTransport(gock.DefaultTransport)
		ctx.Convey("When everything gose postive", func(ctx convey.C) {
			httpMock("POST", _sagaWechatURL+"/appchat/send").Reply(200).JSON(`{"code":0,"message":"0"}`)
			err := d.SendWechatToGroup(c, d.c.WeChat.ChatID, msg)
			ctx.Convey("Then err should be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When http status != 200", func(ctx convey.C) {
			httpMock("POST", _sagaWechatURL+"/appchat/send").Reply(404)
			err := d.SendWechatToGroup(c, d.c.WeChat.ChatID, msg)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When http response code != 0", func(ctx convey.C) {
			httpMock("POST", _sagaWechatURL+"/appchat/send").Reply(200).JSON(`{"code":-401,"message":"0"}`)
			err := d.SendWechatToGroup(c, d.c.WeChat.ChatID, msg)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Reset(func() {
			gock.Off()
			d.client.SetTransport(http.DefaultClient.Transport)
		})
	})
}

func TestDaoGetCoverage(t *testing.T) {
	convey.Convey("GetCoverage", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			commitID = "8d2f1b49661c7089e2b595eafff326033a138c23"
			pkg      = "go-common/app/admin/main/apm"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cov, err := d.GetCoverage(c, commitID, pkg)
			ctx.Convey("Then err should be nil.cov should not be nil.", func(ctx convey.C) {
				// t.Logf("the cov of %s is %.2f", pkg, cov)
				ctx.So(cov, convey.ShouldNotBeNil)
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		ctx.Convey("When the pkg is empty string", func(ctx convey.C) {
			pkg = ""
			_, err := d.GetCoverage(c, commitID, pkg)
			ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
