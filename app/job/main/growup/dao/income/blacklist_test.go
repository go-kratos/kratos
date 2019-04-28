package income

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Blacklist(t *testing.T) {
	Convey("Blacklist", t, func() {
		_, _, err := d.Blacklist(context.Background(), 0, 2000)
		So(err, ShouldBeNil)
	})
}
