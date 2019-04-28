package feed

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Audit(t *testing.T) {
	Convey("should get audit", t, func() {
		_, err := s.Audit(context.Background(), "", 1, 2)
		So(err, ShouldBeNil)
	})
}
