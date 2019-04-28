package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDao_CheckChannelReview(t *testing.T) {
	convey.Convey("实时查询稿件所绑定频道，是否需要回查", t, WithDao(func(d *Dao) {
		in, ids, err := d.CheckChannelReview(context.TODO(), 1)
		t.Logf("in(%v) ids(%s) err(%v)", in, ids, err)
		convey.So(err, convey.ShouldBeNil)
	}))
}
