package track

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() {})
		f(d)
	}
}

func Test_trackvideo(t *testing.T) {
	Convey("根据filename + aid获取视频track信息", t, WithDao(func(d *Dao) {
		aid := int64(101010)
		filename := "i123123123123"
		_, err := d.VideoTrack(context.TODO(), filename, aid)
		So(err, ShouldBeNil)
	}))
}

func Test_ArchiveTrack(t *testing.T) {
	Convey("", t, WithDao(func(d *Dao) {
		aid := int64(101010)
		_, err := d.ArchiveTrack(context.TODO(), aid)
		So(err, ShouldBeNil)
	}))
}
