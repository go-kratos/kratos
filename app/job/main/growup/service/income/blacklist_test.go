package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Blacklist(t *testing.T) {
	Convey("Blacklist", t, func() {
		_, err := s.Blacklist(context.Background(), 10)
		So(err, ShouldBeNil)
	})
}
