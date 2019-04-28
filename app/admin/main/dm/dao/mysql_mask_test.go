package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMaskUps(t *testing.T) {
	Convey("test mask ups", t, func() {
		res, total, err := testDao.MaskUps(context.Background(), 1, 50)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		So(total, ShouldNotBeNil)
		t.Log(res, total)
	})
}

func TestMaskUpOpen(t *testing.T) {
	Convey("test mask up open", t, func() {
		res, err := testDao.MaskUpOpen(context.Background(), 27515266, 1, "")
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
		t.Log(res)
	})
}
