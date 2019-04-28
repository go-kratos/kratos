package ugc

import (
	"fmt"
	"testing"

	"go-common/app/job/main/tv/model/ugc"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_Manual(t *testing.T) {
	Convey("TestDao_Manual", t, WithDao(func(d *Dao) {
		res, err := d.Manual(ctx)
		for _, v := range res {
			fmt.Println(v)
		}
		So(err, ShouldBeNil)
		So(len(res), ShouldBeGreaterThan, 0)
	}))
}

func TestDao_Ppmnl(t *testing.T) {
	Convey("TestDao_Ppmnl", t, WithDao(func(d *Dao) {
		err := d.Ppmnl(ctx, 10099763)
		So(err, ShouldBeNil)
	}))
}

func TestDao_UpdateArc(t *testing.T) {
	Convey("TestDao_UpdateArc", t, WithDao(func(d *Dao) {
		err := d.UpdateArc(ctx, &ugc.ArchDatabus{
			Aid:       10099763,
			Mid:       452156,
			Videos:    1,
			TypeID:    174,
			Title:     "test",
			Cover:     "testPic",
			Content:   "testDesc",
			Duration:  300,
			Copyright: 1,
			PubTime:   "2018-06-05",
			State:     5,
		})
		So(err, ShouldBeNil)
	}))
}
