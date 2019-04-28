package service

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestPlatformChallCount(t *testing.T) {
	convey.Convey("PlatformChallCount", t, func() {
		challCount, err := s.PlatformChallCount(context.Background(), 1, map[int8]int64{2: 11})
		convey.So(err, convey.ShouldBeNil)
		convey.So(challCount.TotalCount, convey.ShouldBeGreaterThanOrEqualTo, int32(0))
	})
}
