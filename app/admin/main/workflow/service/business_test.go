package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestListMeta(t *testing.T) {
	convey.Convey("ListMeta", t, func() {
		ml, err := s.ListMeta(context.Background(), "challenge")
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(ml), convey.ShouldBeGreaterThanOrEqualTo, 1)
	})
}
