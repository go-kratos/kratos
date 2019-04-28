package medal

import (
	"context"
	"encoding/json"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/medal"
	"go-common/library/ecode"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bouk/monkey"
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

func TestRank(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		mid  int64
		data = make([]*medal.FansRank, 0)
		res  = &struct {
			Code int               `json:"code"`
			Data []*medal.FansRank `json:"data"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.fansRankURI).Reply(-502)
		data, err = d.Rank(c, mid)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 0
		res.Data = append(res.Data, &medal.FansRank{})
		js, _ := json.Marshal(res)
		httpMock("Do", d.fansRankURI).Reply(200).JSON(string(js))
		data, err = d.Rank(c, mid)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
}

func TestRename(t *testing.T) {
	var (
		c            = context.TODO()
		err          error
		mid          int64
		name, ck, ak string
		res          = &struct {
			Code int `json:"code"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.fansRankURI).Reply(-502)
		err = d.Rename(c, mid, name, ak, ck)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 0
		ak = "iamak"
		js, _ := json.Marshal(res)
		httpMock("Do", d.fansRankURI).Reply(200).JSON(string(js))
		err = d.Rename(c, mid, name, ak, ck)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestRecentFans(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		mid  = int64(2089809)
		data = make([]*medal.RecentFans, 0)
		res  = &struct {
			Code int                 `json:"code"`
			Data []*medal.RecentFans `json:"data"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.recentFansURI).Reply(-502)
		data, err = d.RecentFans(c, mid)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 0
		res.Data = append(res.Data, &medal.RecentFans{})
		js, _ := json.Marshal(res)
		httpMock("Do", d.recentFansURI).Reply(200).JSON(string(js))
		data, err = d.RecentFans(c, mid)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		req := &http.Request{
			Header: make(map[string][]string),
		}
		req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
		monkeyNewGetRequest(req, ecode.CreativeFansMedalErr)
		data, err = d.RecentFans(c, mid)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
}

func monkeyNewGetRequest(req *http.Request, err error) {
	monkey.Patch(http.NewRequest, func(_, _ string, _ io.Reader) (*http.Request, error) {
		return req, err
	})
}

func TestMedal(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		mid  = int64(2089809)
		data = &medal.Medal{}
		res  = &struct {
			Code int          `json:"code"`
			Data *medal.Medal `json:"data"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.getMedalURI).Reply(-502)
		data, err = d.Medal(c, mid)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 0
		res.Data = &medal.Medal{
			UID: "2089809",
		}
		js, _ := json.Marshal(res)
		httpMock("Do", d.getMedalURI).Reply(200).JSON(string(js))
		data, err = d.Medal(c, mid)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
}

func TestStatus(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		mid  = int64(2089809)
		data *medal.Status
		res  = &struct {
			Code int           `json:"code"`
			Data *medal.Status `json:"data"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.checkStatusURI).Reply(-502)
		data, err = d.Status(c, mid)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 20032
		js, _ := json.Marshal(res)
		httpMock("Do", d.checkStatusURI).Reply(200).JSON(string(js))
		data, err = d.Status(c, mid)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(data, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestOpenMedal(t *testing.T) {
	var (
		c    = context.TODO()
		err  error
		mid  int64
		name string
		res  = &struct {
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.openMedalURI).Reply(-502)
		err = d.OpenMedal(c, mid, name)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 0
		res.Msg = "ok"
		js, _ := json.Marshal(res)
		httpMock("Do", d.openMedalURI).Reply(200).JSON(string(js))
		err = d.OpenMedal(c, mid, name)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		res.Code = 0
		res.Msg = "ok"
		req := &http.Request{
			URL:    &url.URL{},
			Header: make(map[string][]string),
		}
		req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		monkeyNewGetRequest(req, nil)
		err = d.OpenMedal(c, mid, name)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestCheckMedal(t *testing.T) {
	var (
		c     = context.TODO()
		err   error
		mid   int64
		valid int
		name  string
		res   = &struct {
			Code int `json:"code"`
			Data struct {
				Enable bool `json:"enable"`
			} `json:"data"`
		}{}
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("Do", d.checkMedalURI).Reply(-502)
		valid, err = d.CheckMedal(c, mid, name)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(valid, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		defer gock.OffAll()
		res.Code = 0
		js, _ := json.Marshal(res)
		httpMock("Do", d.checkMedalURI).Reply(200).JSON(string(js))
		valid, err = d.CheckMedal(c, mid, name)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(valid, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
	convey.Convey("3", t, func(ctx convey.C) {
		res.Code = 0
		req := &http.Request{
			URL:    &url.URL{},
			Header: make(map[string][]string),
		}
		req.Header.Set("X-BiliLive-UID", strconv.FormatInt(mid, 10))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		monkeyNewGetRequest(req, ecode.CreativeFansMedalErr)
		valid, err = d.CheckMedal(c, mid, name)
		ctx.Convey("3", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(valid, convey.ShouldBeGreaterThanOrEqualTo, 0)
		})
	})
}
