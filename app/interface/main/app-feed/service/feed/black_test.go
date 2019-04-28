package feed

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_BlackList(t *testing.T) {
	Convey("should get BlackList", t, func() {
		_, err := s.BlackList(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}
