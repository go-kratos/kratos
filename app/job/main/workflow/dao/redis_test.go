package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSetList(t *testing.T) {
	convey.Convey("SetList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = "test_list"
			ids = []int64{1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetList(c, key, ids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoExistKey(t *testing.T) {
	convey.Convey("ExistKey", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			exist, err := d.ExistKey(c, key)
			ctx.Convey("Then err should be nil.exist should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(exist, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetString(t *testing.T) {
	convey.Convey("SetString", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			key = ""
			val = ""
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetString(c, key, val)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetCrash(t *testing.T) {
	convey.Convey("SetCrash", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetCrash(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoIsCrash(t *testing.T) {
	convey.Convey("IsCrash", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			exist, err := d.IsCrash(c)
			ctx.Convey("Then err should be nil.exist should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(exist, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUperInfoCache(t *testing.T) {
	convey.Convey("UpInfoCache", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{0}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			groups, err := d.UperInfoCache(c, mids)
			ctx.Convey("Then err should be nil.groups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(groups, convey.ShouldNotBeNil)
				fmt.Println("len groups", len(groups), groups)
			})
		})
	})
}

func TestDaoSetWeightSortedSet(t *testing.T) {
	convey.Convey("SetWeightSortedSet", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			newWeigt = map[int64]int64{1: 200, 2: 1, 3: 10}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.SetWeightSortedSet(c, 1, newWeigt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoKeysSingleExpire(t *testing.T) {
	convey.Convey("KeysSingleExpire", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.SingleExpire(c, 1)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
