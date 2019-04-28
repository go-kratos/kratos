package cms

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmsMixedFilter(t *testing.T) {
	var (
		ctx  = context.Background()
		sids = []int64{}
		aids = []int64{}
	)
	convey.Convey("MixedFilter", t, func(c convey.C) {
		okSids, okAids := d.MixedFilter(ctx, sids, aids)
		c.Convey("Then okSids,okAids should not be nil.", func(c convey.C) {
			c.So(okAids, convey.ShouldNotBeNil)
			c.So(okSids, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsaidsFilter(t *testing.T) {
	var (
		ctx  = context.Background()
		aids = []int64{}
	)
	convey.Convey("aidsFilter", t, func(c convey.C) {
		okAids := d.aidsFilter(ctx, aids)
		c.Convey("Then okAids should not be nil.", func(c convey.C) {
			c.So(okAids, convey.ShouldNotBeNil)
		})
	})
}

func TestCmssidsFilter(t *testing.T) {
	var (
		ctx  = context.Background()
		sids = []int64{}
	)
	convey.Convey("sidsFilter", t, func(c convey.C) {
		okSids := d.sidsFilter(ctx, sids)
		c.Convey("Then okSids should not be nil.", func(c convey.C) {
			c.So(okSids, convey.ShouldNotBeNil)
		})
	})
}
