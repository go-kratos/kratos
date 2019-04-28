package task

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDao_GetWeightRedis(t *testing.T) {
	Convey("CheckChannelReview", t, WithDao(func(d *Dao) {
		m, err := d.GetWeightRedis(context.Background(), []int64{10098217})
		So(err, ShouldBeNil)
		So(m, ShouldNotBeNil)
	}))
}
