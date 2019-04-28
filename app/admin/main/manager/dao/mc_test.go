package dao

import (
	"context"
	"testing"

	"go-common/library/net/http/blademaster/middleware/permit"

	"github.com/smartystreets/goconvey/convey"
)

func TestSession(t *testing.T) {
	convey.Convey("Session", t, func() {
		sid := "1234567890"
		_, err := d.Session(context.Background(), sid)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestSetSession(t *testing.T) {
	convey.Convey("SetSession", t, func() {
		p := &permit.Session{
			Sid: "1234567890",
		}
		err := d.SetSession(context.Background(), p)
		convey.So(err, convey.ShouldBeNil)
	})
}
