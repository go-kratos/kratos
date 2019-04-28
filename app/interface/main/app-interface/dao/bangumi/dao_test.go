package bangumi

import (
	"context"
	"flag"
	"os"
	"testing"

	"go-common/app/interface/main/app-interface/conf"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-interface")
		flag.Set("conf_token", "1mWvdEwZHmCYGoXJCVIdszBOPVdtpXb3")
		flag.Set("tree_id", "2688")
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
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

// go test -conf="../../app-interface-example.toml"  -v -test.run TestSeason
func TestSeason(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		aid = int64(1)
		ip  = ""
	)
	convey.Convey("Season", t, func(ctx convey.C) {
		_, err := d.Season(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

// go test -conf="../../app-interface-example.toml"  -v -test.run TestBPInfo
func TestBPInfo(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		aid = int64(1)
		ip  = ""
	)
	convey.Convey("BPInfo", t, func(ctx convey.C) {
		_, err := d.BPInfo(c, aid, mid, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

// go test -conf="../../app-interface-example.toml"  -v -test.run TestMovie
func TestMovie(t *testing.T) {
	var (
		c         = context.TODO()
		mid       = int64(1)
		aid       = int64(1)
		build     = 123
		mobileApp = "android"
		device    = "android"
		ip        = ""
	)
	convey.Convey("Movie", t, func(ctx convey.C) {
		_, err := d.Movie(c, aid, mid, build, mobileApp, device, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

// go test -conf="../../app-interface-example.toml"  -v -test.run TestConcern
func TestConcern(t *testing.T) {
	var (
		c    = context.TODO()
		mid  = int64(1)
		vmid = int64(1)
		pn   = 1
		ps   = 5
	)
	convey.Convey("Concern", t, func(ctx convey.C) {
		_, _, err := d.Concern(c, mid, vmid, pn, ps)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestFavDisplay(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("Concern", t, func(ctx convey.C) {
		_, _, err := d.FavDisplay(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			err = nil
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
