package upper

import (
	"context"
	upMdl "go-common/app/interface/main/tv/model/upper"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpperupperMetaKey(t *testing.T) {
	var (
		MID = int64(0)
	)
	convey.Convey("upperMetaKey", t, func(c convey.C) {
		p1 := upperMetaKey(MID)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperLoadUpMeta(t *testing.T) {
	var (
		ctx = context.Background()
		mid = int64(0)
	)
	convey.Convey("LoadUpMeta", t, func(c convey.C) {
		upper, err := d.LoadUpMeta(ctx, mid)
		c.Convey("Then err should be nil.upper should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(upper, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperupMetaCache(t *testing.T) {
	var (
		ctx = context.Background()
		mid = int64(0)
	)
	convey.Convey("upMetaCache", t, func(c convey.C) {
		upper, err := d.upMetaCache(ctx, mid)
		c.Convey("Then err should be nil.upper should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(upper, convey.ShouldNotBeNil)
		})
	})
}

func TestUpperupMetaDB(t *testing.T) {
	var (
		ctx = context.Background()
		mid = int64(0)
	)
	convey.Convey("upMetaDB", t, func(c convey.C) {
		upper, err := d.upMetaDB(ctx, mid)
		c.Convey("Then err should be nil.upper should not be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
			c.So(upper, convey.ShouldNotBeNil)
		})
	})
}

func TestUppersetUpMetaCache(t *testing.T) {
	var (
		ctx   = context.Background()
		upper = &upMdl.Upper{}
	)
	convey.Convey("setUpMetaCache", t, func(c convey.C) {
		err := d.setUpMetaCache(ctx, upper)
		c.Convey("Then err should be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUpperaddUpMetaCache(t *testing.T) {
	var (
		upper = &upMdl.Upper{}
	)
	convey.Convey("addUpMetaCache", t, func(c convey.C) {
		d.addUpMetaCache(upper)
		c.Convey("No return values", func(c convey.C) {
		})
	})
}
