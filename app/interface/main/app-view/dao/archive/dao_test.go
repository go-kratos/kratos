package archive

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/net/rpc"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

var (
	d *Dao
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-view")
		flag.Set("conf_token", "3a4CNLBhdFbRQPs7B4QftGvXHtJo92xw")
		flag.Set("tree_id", "4575")
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

func TestViewRPC(t *testing.T) {
	convey.Convey("TestViewRPC", t, func(ctx convey.C) {
		var addr = "172.22.36.185:6089" // new
		// addr = "172.22.38.5:6089"       // old
		client := rpc.Dial(addr, xtime.Duration(100*time.Millisecond), nil)
		var (
			view  *archive.View3
			views map[int64]*archive.View3
			err   error
		)
		if err = client.Call(context.TODO(), "RPC.View3", &archive.ArgAid{Aid: 10111165}, &view); err != nil {
			ctx.Println(err)
			return
		}
		ctx.Println(view)
		if err = client.Call(context.TODO(), "RPC.Views3", &archive.ArgAids{Aids: []int64{10111165}}, &views); err != nil {
			ctx.Println(err)
			return
		}
		ctx.Println(views)
	})
}

func TestPing(t *testing.T) {
	var (
		c = context.TODO()
	)
	convey.Convey("Ping", t, func(ctx convey.C) {
		d.Ping(c)
	})
}

func TestArchive3(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(0)
	)
	convey.Convey("Archive3", t, func(ctx convey.C) {
		_, err := d.Archive3(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchives(t *testing.T) {
	var (
		c    = context.TODO()
		aids = []int64{1}
	)
	convey.Convey("Archives", t, func(ctx convey.C) {
		_, err := d.Archives(c, aids)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestShot(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		cid = int64(2)
	)
	convey.Convey("Shot", t, func(ctx convey.C) {
		d.Shot(c, aid, cid)
	})
}

func TestUpCount2(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
	)
	convey.Convey("UpCount2", t, func(ctx convey.C) {
		_, err := d.UpCount2(c, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpArcs3(t *testing.T) {
	var (
		c   = context.TODO()
		mid = int64(1)
		pn  = int(1)
		ps  = int(20)
	)
	convey.Convey("UpArcs3", t, func(ctx convey.C) {
		d.UpArcs3(c, mid, pn, ps)
	})
}

func TestProgress(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(1)
		mid = int64(1)
	)
	convey.Convey("Progress", t, func(ctx convey.C) {
		_, err := d.Progress(c, aid, mid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchive(t *testing.T) {
	var (
		c   = context.TODO()
		aid = int64(-1)
	)
	convey.Convey("Archive", t, func(ctx convey.C) {
		_, err := d.Archive(c, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}
