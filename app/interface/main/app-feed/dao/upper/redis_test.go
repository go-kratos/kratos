package upper

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UnreadCountCache(t *testing.T) {
	Convey("should get UnreadCountCache", t, func() {
		_, err := d.UnreadCountCache(context.Background(), 1)
		So(err, ShouldBeNil)
	})
}
