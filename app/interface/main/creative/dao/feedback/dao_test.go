package feedback

import (
	"context"
	"encoding/json"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/feedback"
	"go-common/library/ecode"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative")
		flag.Set("conf_token", "96b6a6c10bb311e894c14a552f48fef8")
		flag.Set("tree_id", "2305")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.SetTransport(gock.DefaultTransport)
	return r
}

func TestCloseSession(t *testing.T) {
	var (
		sessionID int64
		ip        string
		c         = context.TODO()
		err       error
		res       struct {
			Code int `json:"code"`
		}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Post", d.closeURI).Reply(-502)
		err = d.CloseSession(c, sessionID, ip)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res.Code = 20010
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Post", d.closeURI).Reply(200).JSON(string(js))
		err = d.CloseSession(c, sessionID, ip)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		res.Code = 0
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Post", d.closeURI).Reply(200).JSON(string(js))
		err = d.CloseSession(c, sessionID, ip)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestAddFeedback(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		res struct {
			Code int `json:"code"`
		}
		mid, tagID, sessionID                           int64
		qq, content, aid, browser, imgURL, platform, ip string
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Post", d.addURI).Reply(-502)
		err = d.AddFeedback(c, mid, tagID, sessionID, qq, content, aid, browser, imgURL, platform, ip)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res.Code = 20010
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Post", d.addURI).Reply(200).JSON(string(js))
		err = d.AddFeedback(c, mid, tagID, sessionID, qq, content, aid, browser, imgURL, platform, ip)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		res.Code = 0
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Post", d.addURI).Reply(200).JSON(string(js))
		err = d.AddFeedback(c, mid, tagID, sessionID, qq, content, aid, browser, imgURL, platform, ip)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDetail(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		res struct {
			Code int               `json:"code"`
			Data []*feedback.Reply `json:"data"`
		}
		mid, sessionID int64
		ip             string
		data           []*feedback.Reply
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Get", d.detailURI).Reply(-502)
		data, err = d.Detail(c, mid, sessionID, ip)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res.Code = 20010
		res.Data = make([]*feedback.Reply, 0)
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Get", d.detailURI).Reply(200).JSON(string(js))
		data, err = d.Detail(c, mid, sessionID, ip)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		res.Code = 0
		res.Data = make([]*feedback.Reply, 0)
		res.Data = append(res.Data, &feedback.Reply{})
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Get", d.detailURI).Reply(200).JSON(string(js))
		data, err = d.Detail(c, mid, sessionID, ip)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(data, convey.ShouldNotBeNil)
		})
	})
}

func TestFeedbacks(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		res struct {
			Code  int                  `json:"code"`
			Data  []*feedback.Feedback `json:"data"`
			Count int64                `json:"total"`
		}
		mid, ps, pn, cnt  int64
		tagID             = int64(1)
		end, platform, ip string
		start             = "1090"
		state             = "open"
		data              []*feedback.Feedback
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Get", d.listURI).Reply(-502)
		data, cnt, err = d.Feedbacks(c, mid, ps, pn, tagID, state, start, end, platform, ip)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
			ctx.So(data, convey.ShouldBeNil)
			ctx.So(cnt, convey.ShouldEqual, 0)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res.Code = 20010
		res.Data = make([]*feedback.Feedback, 0)
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Get", d.listURI).Reply(200).JSON(string(js))
		data, cnt, err = d.Feedbacks(c, mid, ps, pn, tagID, state, start, end, platform, ip)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.CreativeFeedbackErr)
			ctx.So(data, convey.ShouldBeNil)
			ctx.So(cnt, convey.ShouldEqual, 0)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		res.Code = 0
		res.Data = make([]*feedback.Feedback, 0)
		res.Data = append(res.Data, &feedback.Feedback{})
		res.Count = 100
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("Get", d.listURI).Reply(200).JSON(string(js))
		data, cnt, err = d.Feedbacks(c, mid, ps, pn, tagID, state, start, end, platform, ip)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(data, convey.ShouldNotBeNil)
			ctx.So(cnt, convey.ShouldNotEqual, 0)
		})
	})
}
