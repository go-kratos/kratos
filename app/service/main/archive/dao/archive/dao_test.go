package archive

import (
	"context"
	"flag"
	"go-common/app/service/main/archive/conf"
	"os"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.archive-service")
		flag.Set("conf_token", "Y2LJhIsHx87nJaOBSxuG5TeZoLdBFlrE")
		flag.Set("tree_id", "2302")
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

func TestArchiveArchive3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(-1)
	)
	convey.Convey("Archive3", t, func(ctx convey.C) {
		_, err := d.Archive3(c, aid)
		ctx.Convey("Then err should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveArchives3(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{1009799822}
	)
	convey.Convey("Archives3", t, func(ctx convey.C) {
		am, err := d.Archives3(c, aids)
		ctx.Convey("Then err should be nil.am should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(am, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveVideos3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10097272)
	)
	convey.Convey("Videos3", t, func(ctx convey.C) {
		vs, err := d.Videos3(c, aid)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vs, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveVideosByAids3(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{10097272, 10098500}
	)
	convey.Convey("VideosByAids3", t, func(ctx convey.C) {
		vs, err := d.VideosByAids3(c, aids)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(vs, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveVideo3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		cid = int64(2)
	)
	convey.Convey("Video3", t, func(ctx convey.C) {
		v, err := d.Video3(c, aid, cid)
		ctx.Convey("Then err should not be nil.v should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(v, convey.ShouldBeNil)
		})
	})
}

func TestArchiveDescription(t *testing.T) {
	var (
		c    = context.TODO()
		aid  = int64(1)
		desc = "我是个大描述我是个大描述"
	)
	convey.Convey("Description", t, func(ctx convey.C) {
		err := d.addDescCache(c, aid, desc)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
		desc, err := d.descCache(context.TODO(), aid)
		ctx.Convey("Then err should be nil.desc should not be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(desc, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveUpVideo3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		cid = int64(2)
	)
	convey.Convey("UpVideo3", t, func(ctx convey.C) {
		v, err := d.Video3(c, aid, cid)
		ctx.Convey("Then err should not be nil.v should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(v, convey.ShouldBeNil)
		})
	})
}

func TestArchiveUpperCache(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(10098500)
	)
	convey.Convey("UpperCache", t, func(ctx convey.C) {
		err := d.UpArchiveCache(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func Test_Ping(t *testing.T) {
	var c = context.TODO()
	convey.Convey("Ping", t, func(ctx convey.C) {
		err := d.Ping(c)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
