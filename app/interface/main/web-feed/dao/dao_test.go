package dao

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Dao(t *testing.T) {
	c := context.TODO()
	mid := int64(1)
	var err error
	d := New(nil)
	Convey("set cache", t, func() {
		PromError("")
		PromInfo("")
		err = d.SetFeedCache(c, mid, nil)
		_, err = d.FeedCache(c, mid)
		So(err, ShouldBeNil)
	})
}
