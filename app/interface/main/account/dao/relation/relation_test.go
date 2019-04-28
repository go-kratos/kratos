package relation

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"

	"go-common/app/interface/main/account/conf"
	bm "go-common/library/net/http/blademaster"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
)

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}
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
	d.httpClient.SetTransport(gock.DefaultTransport)

	m.Run()
	os.Exit(0)
}

func TestRelationpaltform(t *testing.T) {
	var (
		device = &bm.Device{}
	)
	convey.Convey("paltform", t, func(ctx convey.C) {
		p1 := paltform(device)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestRelationbuvid(t *testing.T) {
	var (
		device = &bm.Device{}
	)
	convey.Convey("buvid", t, func(ctx convey.C) {
		p1 := buvid(device)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestRelationRecommend(t *testing.T) {
	var (
		ctx         = context.Background()
		mid         = int64(0)
		serviceArea = ""
		mainTids    = ""
		subTids     = ""
		device      = &bm.Device{}
		pagesize    = int64(0)
		ip          = ""
	)
	convey.Convey("Recommend", t, func(c convey.C) {
		httpMock("GET", d.recommendURL).Reply(200).JSON(`{"code":0,"message":"0","data":[]}`)
		defer gock.OffAll()

		p1, err := d.Recommend(ctx, mid, serviceArea, mainTids, subTids, device, pagesize, ip)
		c.Convey("Then err should be nil.p1 should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestRelationTagSuggestRecommend(t *testing.T) {
	var (
		ctx       = context.Background()
		mid       = int64(0)
		contextID = ""
		tagname   = ""
		device    = &bm.Device{}
		pagesize  = int64(0)
		ip        = ""
	)
	convey.Convey("TagSuggestRecommend", t, func(c convey.C) {
		httpMock("GET", d.recommendURL).Reply(200).JSON(`{"code":0,"message":"0","data":[]}`)
		defer gock.OffAll()
		p1, err := d.TagSuggestRecommend(ctx, mid, contextID, tagname, device, pagesize, ip)
		c.Convey("Then err should be nil.p1 should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}
