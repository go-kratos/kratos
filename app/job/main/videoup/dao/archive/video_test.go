package archive

import (
	"context"

	"go-common/app/job/main/videoup/model/archive"

	_ "github.com/go-sql-driver/mysql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_NewVideo(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.Video
	)
	Convey("NewVideo", t, WithDao(func(d *Dao) {
		sub, err = d.NewVideo(c, "2333")
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}

func Test_NewSumDuration(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int64
	)
	Convey("NewVideo", t, WithDao(func(d *Dao) {
		sub, err = d.NewSumDuration(c, 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeZeroValue)
	}))
}

func Test_NewVideoByAid(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.Video
	)
	Convey("NewVideoByAid", t, WithDao(func(d *Dao) {
		sub, err = d.NewVideoByAid(c, "filename", 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}

func Test_NewVideoCount(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int
	)
	Convey("NewVideoCount", t, WithDao(func(d *Dao) {
		sub, err = d.NewVideoCount(c, 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeZeroValue)
	}))
}

func Test_NewVideoCountCapable(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int
	)
	Convey("NewVideoCountCapable", t, WithDao(func(d *Dao) {
		sub, err = d.NewVideoCountCapable(c, 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeZeroValue)
	}))
}

// old TODO deprecated
func Test_Video(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.Video
	)
	Convey("Video", t, WithDao(func(d *Dao) {
		sub, err = d.NewVideo(c, "2333")
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}

func Test_SumDuration(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int64
	)
	Convey("Video", t, WithDao(func(d *Dao) {
		sub, err = d.SumDuration(c, 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeZeroValue)
	}))
}

func Test_VideoByAid(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub *archive.Video
	)
	Convey("VideoByAid", t, WithDao(func(d *Dao) {
		sub, err = d.VideoByAid(c, "filename", 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeNil)
	}))
}

func Test_VideoCount(t *testing.T) {
	var (
		c   = context.TODO()
		err error
	)
	Convey("VideoCount", t, WithDao(func(d *Dao) {
		_, err = d.VideoCount(c, 2333)
		So(err, ShouldBeNil)
	}))
}

func Test_VideoCountCapable(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub int
	)
	Convey("VideoCountCapable", t, WithDao(func(d *Dao) {
		sub, err = d.VideoCountCapable(c, 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeZeroValue)
	}))
}

func Test_Reason(t *testing.T) {
	var (
		c   = context.TODO()
		err error
		sub string
	)
	Convey("Reason", t, WithDao(func(d *Dao) {
		sub, err = d.Reason(c, 2333)
		So(err, ShouldBeNil)
		So(sub, ShouldBeEmpty)
	}))
}
