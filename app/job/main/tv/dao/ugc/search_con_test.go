package ugc

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_UgcCont(t *testing.T) {
	Convey("TestDao_UgcCont", t, WithDao(func(d *Dao) {
		var aid int
		d.DB.QueryRow(ctx, queryAid).Scan(&aid)
		if aid == 0 {
			fmt.Println("No ready audit Data")
			return
		}
		if res, maxID, err := d.UgcCont(ctx, aid, 50); err != nil {
			fmt.Println(err)
			So(len(res), ShouldEqual, 0)
			So(maxID, ShouldBeZeroValue)
		} else {
			So(len(res), ShouldBeGreaterThan, 0)
			for _, v := range res {
				fmt.Println(v)
			}
			So(maxID, ShouldBeGreaterThan, 0)
		}
	}))
}

func TestDao_UgcContCount(t *testing.T) {
	Convey("TestDao_UgcContCount", t, WithDao(func(d *Dao) {
		var aid int
		d.DB.QueryRow(ctx, queryAid).Scan(&aid)
		if aid != 0 {
			res, _ := d.UgcCnt(ctx)
			So(res, ShouldBeGreaterThan, 0)
			fmt.Println(res)
		}
	}))
}
