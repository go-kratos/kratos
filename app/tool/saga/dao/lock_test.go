package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLock(t *testing.T) {
	var (
		key string
		ok  bool
		err error
	)
	Convey("TEST Lock", t, func() {
		ok, err = d.TryLock(ctx, key, "test", 1)
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		ok, err = d.TryLock(ctx, key, "test", 1)
		So(err, ShouldBeNil)
		So(ok, ShouldBeFalse)

		err = d.UnLock(ctx, key)
		So(err, ShouldBeNil)

		ok, err = d.TryLock(ctx, key, "test", 1)
		So(err, ShouldBeNil)
		So(ok, ShouldBeTrue)

		ok, err = d.TryLock(ctx, key, "test", 1)
		So(err, ShouldBeNil)
		So(ok, ShouldBeFalse)
	})
}
