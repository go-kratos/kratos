package dao

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDaoCountBlocked(t *testing.T) {
	Convey("should return err be nil and num>=0", t, func() {
		num, err := d.CountBlocked(context.TODO(), 27515415, time.Now())
		So(err, ShouldBeNil)
		So(num, ShouldBeGreaterThanOrEqualTo, 0)
	})
}

func TestDaoBlockedInfoID(t *testing.T) {
	Convey("should return err be nil and num>=0", t, func() {
		id, err := d.BlockedInfoID(context.TODO(), 27515415)
		So(err, ShouldBeNil)
		So(id, ShouldBeGreaterThanOrEqualTo, 0)
	})
}
