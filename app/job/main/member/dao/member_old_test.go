package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoAccounts(t *testing.T) {
	convey.Convey("Accounts", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{1, 2}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			accs, errs, err := d.Accounts(c, mids)
			if err != nil {
				convCtx.Convey("When no rows int result", func(convCtx convey.C) {
					convCtx.So(err, convey.ShouldNotBeNil)
					convCtx.So(errs, convey.ShouldNotBeNil)
					convCtx.So(accs, convey.ShouldNotBeNil)
				})
			} else {
				convCtx.Convey("When have rows int result", func(convCtx convey.C) {
					convCtx.So(err, convey.ShouldBeNil)
					convCtx.So(errs, convey.ShouldNotBeNil)
					convCtx.So(accs, convey.ShouldNotBeNil)
				})
			}
		})
	})
}

func TestDaoName(t *testing.T) {
	convey.Convey("Name", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(1)
		)
		convCtx.Convey("Name search two situation:", func(convCtx convey.C) {
			name, err := d.Name(c, mid)
			if err != nil {
				convCtx.Convey("Name occur an error", func(convCtx convey.C) {
					convCtx.So(err, convey.ShouldNotBeNil)
					convCtx.So(name, convey.ShouldNotBeNil)
				})
			} else {
				convCtx.Convey("Name no err search", func(convCtx convey.C) {
					convCtx.So(err, convey.ShouldBeNil)
					convCtx.So(name, convey.ShouldNotBeNil)
				})
			}
		})
	})
}
