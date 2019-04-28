package service

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Service(t *testing.T) {
	c := context.TODO()
	var err error
	s := New(nil)
	Convey("test", t, func() {
		s.Ping(c)
		So(err, ShouldBeNil)
	})
}
