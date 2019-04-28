package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_VideoByID(t *testing.T) {
	Convey("VideoByID", t, WithDao(func(d *Dao) {
		_, err := d.VideoByID(context.Background(), 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_VideosByAid(t *testing.T) {
	Convey("VideosByAid", t, WithDao(func(d *Dao) {
		_, err := d.VideosByAid(context.Background(), 1)
		So(err, ShouldBeNil)
	}))
}

func TestDao_VideoByIDs(t *testing.T) {
	Convey("VideoByIDs", t, WithDao(func(d *Dao) {
		_, err := d.VideoByIDs(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}
func Test_VideoStateMap(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, err := d.VideoStateMap(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}

func Test_VideoAidMap(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, err := d.VideoAidMap(context.Background(), []int64{1, 2})
		So(err, ShouldBeNil)
	}))
}
