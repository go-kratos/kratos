package archive

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_videoOperAttrsCtimes(t *testing.T) {
	Convey("获取vid的审核属性记录，按照ctime排序", t, WithDao(func(d *Dao) {
		vid := int64(8943315)
		attrs, ctimes, err := d.VideoOperAttrsCtimes(context.TODO(), vid)
		So(err, ShouldBeNil)
		So(len(attrs), ShouldEqual, len(ctimes))
	}))
}

func Test_ArchiveOper(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		m, err := d.ArchiveOper(context.Background(), 10116994)
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	}))
}

func Test_VideoOper(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		m, err := d.VideoOper(context.Background(), 10116994)
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	}))
}

func Test_PassedOper(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		m, err := d.PassedOper(context.Background(), 10116994)
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	}))
}

func Test_VideoOperAttrsCtimes(t *testing.T) {
	Convey("test archive", t, WithDao(func(d *Dao) {
		_, _, err := d.VideoOperAttrsCtimes(context.Background(), 10116994)
		So(err, ShouldBeNil)
	}))
}

func Test_UpVideoOper(t *testing.T) {
	var c = context.Background()
	Convey("UpVideoOper", t, WithDao(func(d *Dao) {
		_, err := d.UpVideoOper(c, 0, 0)
		So(err, ShouldBeNil)
	}))
}
