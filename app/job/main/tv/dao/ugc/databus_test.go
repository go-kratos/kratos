package ugc

import (
	"fmt"
	"testing"

	"go-common/app/service/main/archive/api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_UpInList(t *testing.T) {
	Convey("TestDao_UpInList", t, WithDao(func(d *Dao) {
		res, err := d.UpInList(ctx, 27515615)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
		fmt.Println(res)
		res2, err2 := d.UpInList(ctx, 100997637777)
		So(err2, ShouldBeNil)
		So(res2, ShouldEqual, 0)
	}))
}

func TestDao_PickVideos(t *testing.T) {
	Convey("TestDao_PickVideos", t, WithDao(func(d *Dao) {
		res, err := d.PickVideos(ctx, 10099763)
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
		for _, v := range res {
			fmt.Println(v)
		}
	}))
}

func TestDao_InsertVideos(t *testing.T) {
	Convey("TestDao_InsertVideos", t, WithDao(func(d *Dao) {
		tx, err := d.BeginTran(ctx)
		So(err, ShouldBeNil)
		err = d.TxAddVideos(tx, []*api.Page{
			{
				Cid:      10126229,
				Part:     "test",
				Duration: 2333,
				Desc:     "test",
				Page:     999,
			},
		}, 10098693)
		tx.Commit()
		So(err, ShouldBeNil)
	}))
}
