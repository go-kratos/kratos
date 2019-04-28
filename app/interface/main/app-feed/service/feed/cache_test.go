package feed

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_indexCache(t *testing.T) {
	Convey("should get indexCache", t, func() {
		_, err := s.indexCache(context.Background(), 1, 2)
		So(err, ShouldBeNil)
	})
}
