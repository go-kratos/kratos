package upper

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"

	"go-common/app/job/main/tv/conf"
	ugcMdl "go-common/app/job/main/tv/model/ugc"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx = context.TODO()
	d   *Dao
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

func TestDao_LoadUpMeta(t *testing.T) {
	Convey("TestDao_LoadUpMeta", t, WithDao(func(d *Dao) {
		res, err := d.LoadUpMeta(ctx, 88895270)
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func TestDao_CountUP(t *testing.T) {
	Convey("TestDao_CountUP", t, WithDao(func(d *Dao) {
		res, err := d.CountUP(ctx)
		So(err, ShouldBeNil)
		fmt.Println(res)
	}))
}

func TestDao_PickUppers(t *testing.T) {
	Convey("TestDao_PickUppers", t, WithDao(func(d *Dao) {
		res, myLast, err := d.PickUppers(ctx, 0, 50)
		So(err, ShouldBeNil)
		So(myLast, ShouldBeGreaterThan, 0)
		So(len(res), ShouldBeGreaterThan, 0)
		fmt.Println(myLast)
	}))
}

func TestDao_RefreshUp(t *testing.T) {
	Convey("TestDao_RefreshUp", t, WithDao(func(d *Dao) {
		res, _, _ := d.PickUppers(ctx, 0, 50)
		if len(res) > 0 {
			err := d.RefreshUp(ctx, &ugcMdl.ReqSetUp{
				MID:    res[0],
				Value:  "Face666",
				UpType: _upFace,
			})
			So(err, ShouldBeNil)
		}
	}))
}

func TestDao_TosyncUps(t *testing.T) {
	Convey("TestDao_TosyncUps", t, WithDao(func(d *Dao) {
		mids, err := d.TosyncUps(ctx)
		So(err, ShouldBeNil)
		fmt.Println(mids)
	}))
}

func TestDao_FinsyncUps(t *testing.T) {
	Convey("TestDao_FinsyncUps", t, WithDao(func(d *Dao) {
		err := d.FinsyncUps(ctx, 88895270)
		So(err, ShouldBeNil)
	}))
}

func TestDao_IsSame(t *testing.T) {
	Convey("TestDao_IsSame", t, WithDao(func(d *Dao) {
		var u = &ugcMdl.Upper{
			OriFace: "http://i0.hdslb.com/bfs/face/589568b0b0228d5a38b2e6f14b83bcb174637140.jpg",
		}
		var faceSame, _ = u.IsSame("", "http://i1.hdslb.com/bfs/face/589568b0b0228d5a38b2e6f14b83bcb174637140.jpg")
		So(faceSame, ShouldBeTrue)
		faceSame, _ = u.IsSame("", "http://i0.hdslb.com/bfs/face/xxx.jpg")
		So(faceSame, ShouldBeFalse)
	}))
}
