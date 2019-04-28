package feed

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UpperFeed(t *testing.T) {
	Convey("should get UpperFeed", t, func() {
		uf, _ := s.UpperFeed(context.Background(), 1, 2, 3, 4, 5, time.Now())
		So(uf, ShouldNotBeNil)
	})
}
