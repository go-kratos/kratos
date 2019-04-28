package watermark

import (
	"context"
	"testing"

	"github.com/davecgh/go-spew/spew"

	"go-common/app/interface/main/creative/model/watermark"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_WaterMark(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		res       *watermark.Watermark
		MID       = int64(27515256)
		localHost = "127.0.0.1"
	)
	wp := &watermark.WatermarkParam{
		MID:   MID,
		State: 0,
		Ty:    1,
		Pos:   1,
		IP:    localHost,
	}
	Convey("WaterMark", t, WithService(func(s *Service) {
		res, err = s.WaterMark(c, MID)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
	Convey("WaterMarkSet", t, WithService(func(s *Service) {
		res, err = s.WaterMarkSet(c, wp)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}

func Test_WaterMarkSet(t *testing.T) {
	var (
		c         = context.TODO()
		err       error
		res       *watermark.Watermark
		MID       = int64(27515256)
		localHost = "127.0.0.1"
	)
	wp := &watermark.WatermarkParam{
		MID:   MID,
		State: 0,
		Ty:    1,
		Pos:   1,
		IP:    localHost,
	}
	Convey("WaterMarkSet", t, WithService(func(s *Service) {
		res, err = s.WaterMarkSet(c, wp)
		spew.Dump(res, err)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
	}))
}
