package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_DelImage(t *testing.T) {
	Convey("DelImage", t, func() {
		Convey("when not exist, err should not nil", func() {
			// FIXME : 该UT不可重入
			err := d.DelImage(context.Background(), "130d0e4aa718fcab4abd3ad756ef57017e42bcf4.png", d.c.FaceBFS)
			So(err, ShouldNotBeNil)
		})
	})
}
