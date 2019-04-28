package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func Test_getVideoOperInfo(t *testing.T) {
	Init()

	convey.Convey("getVideoOperInfo", t, func() {
		list, err := s.getVideoOperInfo(context.TODO(), 8942606)
		for _, info := range list {
			t.Logf("info(%+v)", info)
		}
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(list), convey.ShouldEqual, 66)
	})
}
