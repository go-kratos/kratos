package faq

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/faq"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d             *Dao
	errConnClosed = errors.New("redigo: connection closed")
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

func TestFaqDetail(t *testing.T) {
	var (
		qTypeID         string
		keyFlag, pn, ps int
		c               = context.TODO()
		err             error
		res             struct {
			Code  string        `json:"retCode"`
			Data  []*faq.Detail `json:"items"`
			Total int           `json:"totalCount"`
		}
		data  []*faq.Detail
		total int
	)
	convey.Convey("1", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.searchURL).Reply(-502)
		data, total, err = d.Detail(c, qTypeID, keyFlag, pn, ps)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.HelpDetailError)
			ctx.So(total, convey.ShouldBeZeroValue)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res.Code = "000000"
		res.Data = make([]*faq.Detail, 0)
		res.Total = 100
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("GET", d.searchURL).Reply(200).JSON(string(js))
		data, total, err = d.Detail(c, qTypeID, keyFlag, pn, ps)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(data, convey.ShouldNotBeNil)
			ctx.So(total, convey.ShouldEqual, 100)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		res.Code = "0000001"
		res.Data = make([]*faq.Detail, 0)
		res.Total = 100
		js, _ := json.Marshal(res)
		defer gock.OffAll()
		httpMock("GET", d.searchURL).Reply(200).JSON(string(js))
		data, total, err = d.Detail(c, qTypeID, keyFlag, pn, ps)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(data, convey.ShouldBeNil)
			ctx.So(err, convey.ShouldEqual, ecode.HelpDetailError)
			ctx.So(total, convey.ShouldEqual, 0)
		})
	})
}

func TestFaqDetailCache(t *testing.T) {
	var (
		qTypeID         = "faq_id"
		keyFlag, pn, ps = int(1), int(1), int(10)
		c               = context.TODO()
		err             error
		data            []*faq.Detail
		total           int
	)
	convey.Convey("1", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.redis), "Get", func(_ *redis.Pool, _ context.Context) redis.Conn {
			return redis.MockWith(errConnClosed)
		})
		defer connGuard.Unpatch()
		data, total, err = d.DetailCache(c, qTypeID, keyFlag, pn, ps)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(total, convey.ShouldBeZeroValue)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		connGuard := monkey.Patch(redis.Values, func(_ interface{}, _ error) ([]interface{}, error) {
			detail := &faq.Detail{}
			data, _ := json.Marshal(detail)
			res := []interface{}{
				data,
			}
			return res, nil
		})
		defer connGuard.Unpatch()
		data, total, err = d.DetailCache(c, qTypeID, keyFlag, pn, ps)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(total, convey.ShouldBeZeroValue)
			ctx.So(data, convey.ShouldBeNil)
		})
	})
}

func TestFaqSetDetailCache(t *testing.T) {
	var (
		qTypeID         = "faq_id"
		keyFlag, pn, ps = int(1), int(1), int(10)
		c               = context.TODO()
		err             error
		data            = make([]*faq.Detail, 0)
		total           int
	)
	convey.Convey("1", t, func(ctx convey.C) {
		connGuard := monkey.PatchInstanceMethod(reflect.TypeOf(d.redis), "Get", func(_ *redis.Pool, _ context.Context) redis.Conn {
			return redis.MockWith(errConnClosed)
		})
		defer connGuard.Unpatch()
		err = d.SetDetailCache(c, qTypeID, keyFlag, pn, ps, total, data)
		ctx.Convey("1", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("2", t, func(ctx convey.C) {
		data = append(data, &faq.Detail{})
		err = d.SetDetailCache(c, qTypeID, keyFlag, pn, ps, total, data)
		ctx.Convey("2", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
	// convey.Convey("2", t, func(ctx convey.C) {
	// 	connGuard := monkey.Patch(redis.Values, func(_ interface{}, _ error) ([]interface{}, error) {
	// 		detail := &faq.Detail{}
	// 		data, _ := json.Marshal(detail)
	// 		res := []interface{}{
	// 			data,
	// 		}
	// 		return res, nil
	// 	})
	// 	defer connGuard.Unpatch()
	// 	data, total, err = d.DetailCache(c, qTypeID, keyFlag, pn, ps)
	// 	ctx.Convey("2", func(ctx convey.C) {
	// 		ctx.So(err, convey.ShouldNotBeNil)
	// 		ctx.So(total, convey.ShouldBeZeroValue)
	// 		ctx.So(data, convey.ShouldBeNil)
	// 	})
	// })
}
