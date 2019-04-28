package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsert(t *testing.T) {
	convey.Convey("Insert", t, func(ctx convey.C) {
		var (
			id          = uint64(0)
			content     = ""
			forceUpdate bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Insert(id, content, forceUpdate)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDaoFlush(t *testing.T) {
	convey.Convey("Flush", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Flush()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestDaoRemove(t *testing.T) {
	convey.Convey("Remove", t, func(ctx convey.C) {
		var (
			id          = uint64(0)
			forceUpdate bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Remove(id, forceUpdate)
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}
