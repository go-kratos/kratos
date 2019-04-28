package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMaskOneVideo(t *testing.T) {
	Convey("mask one video", t, func() {
		err := svr.maskOneVideo(context.TODO(), 8936701)
		So(err, ShouldBeNil)
		t.Logf("err:%v", err)
	})
}

func TestMaskOneArchive(t *testing.T) {
	Convey("mask one archive", t, func() {
		err := svr.maskOneArchive(context.TODO(), 10098039)
		So(err, ShouldBeNil)
		t.Logf("err:%v", err)
	})
}

func TestMaskOneCate(t *testing.T) {
	Convey("mask one cate", t, func() {
		err := svr.maskOneCate(context.TODO(), 185)
		So(err, ShouldBeNil)
		t.Logf("err:%v", err)
	})
}
