package feed

import (
	"context"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_DislikeCancel(t *testing.T) {
	Convey("should get DislikeCancel", t, func() {
		err := s.DislikeCancel(context.Background(), 1, 2, "", "", 3, 9, 4, 8, 5, 6, "", time.Now())
		So(err, ShouldBeNil)
	})
}
