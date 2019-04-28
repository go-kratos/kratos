package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Ping(t *testing.T) {
	var err error
	var c = context.Background()
	Convey("test archive video", t, func() {
		if err = s.Ping(c); err != nil {
			t.Fail()
		}
		So(err, ShouldBeNil)
	})
}
