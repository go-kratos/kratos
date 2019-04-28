package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_LikeSync(t *testing.T) {
	Convey("work", t, WithDao(func(d *Dao) {
		err := d.LikeSync(context.Background(), 1, 20)
		So(err, ShouldBeNil)
	}))
}
