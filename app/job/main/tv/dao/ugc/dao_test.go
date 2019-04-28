package ugc

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/job/main/tv/conf"
	"go-common/app/job/main/tv/model/ugc"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx      = context.TODO()
	d        *Dao
	queryMid = "SELECT mid FROM ugc_archive WHERE result=1 AND valid=1 AND deleted=0 LIMIT 1"
	queryAid = "SELECT aid FROM ugc_archive WHERE result=1 AND valid=1 AND deleted=0 LIMIT 1"
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		dir, _ := filepath.Abs("../../cmd/tv-job-test.toml")
		flag.Set("conf", dir)
		conf.Init()
		if d == nil {
			d = New(conf.Conf)
		}
		f(d)
	}
}

func TestDao_UpArcs(t *testing.T) {
	Convey("TestDao_UpArcs", t, WithDao(func(d *Dao) {
		var mid int64
		d.DB.QueryRow(ctx, queryMid).Scan(&mid)
		if mid == 0 {
			return
		}
		res, err := d.UpArcs(ctx, mid)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestDao_SetArcCMS(t *testing.T) {
	Convey("TestDao_SetArcCMS", t, WithDao(func(d *Dao) {
		var aid int64
		d.DB.QueryRow(ctx, queryAid).Scan(&aid)
		if aid != 0 {
			err := d.SetArcCMS(ctx, &ugc.ArcCMS{
				AID:   aid,
				Title: "test",
			})
			So(err, ShouldBeNil)
			fmt.Println(aid)
		}
	}))
}

func TestDao_CountArcs(t *testing.T) {
	Convey("TestDao_CountArcs", t, WithDao(func(d *Dao) {
		res, err := d.CountUpArcs(ctx, 452156)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
		fmt.Println(res)
	}))
}
