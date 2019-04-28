package client

import (
	"context"
	"testing"

	"go-common/app/interface/main/dm2/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSubtitleGet(t *testing.T) {
	var (
		cid int64 = 401
		aid int64 = 4052462
		tp  int32 = 1
	)
	Convey("test mask list", t, func() {
		arg := &model.ArgSubtitleGet{Oid: cid, Aid: aid, Type: tp}
		res, err := svr.SubtitleGet(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("===============%+v", res)
	})
}

func TestSubtitleSujectSubmit(t *testing.T) {
	Convey("test mask list", t, func() {
		arg := &model.ArgSubtitleAllowSubmit{Aid: 4052462, AllowSubmit: true, Lan: "zh-CN"}
		err := svr.SubtitleSujectSubmit(context.TODO(), arg)
		So(err, ShouldBeNil)
	})
}

func TestSubtitleSubjectSubmitGet(t *testing.T) {
	Convey("test mask list", t, func() {
		arg := &model.ArgArchiveID{Aid: 4052462}
		res, err := svr.SubtitleSubjectSubmitGet(context.TODO(), arg)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Logf("===============%+v", res)
	})
}
