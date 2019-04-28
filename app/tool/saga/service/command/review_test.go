package command

import (
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCommandReviewedOwner(t *testing.T) {
	convey.Convey("ownerReviewed", t, func(ctx convey.C) {
		var (
			owners        = []string{"a", "b"}
			reviewedUsers = []string{"a"}
			username      = "d"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			isowner, reviewed := reviewedOwner(owners, reviewedUsers, username)
			fmt.Println(isowner, reviewed)
			ctx.Convey("Then isowner,reviewed should not be nil.", func(ctx convey.C) {
				ctx.So(reviewed, convey.ShouldBeTrue)
				ctx.So(isowner, convey.ShouldBeFalse)
			})
		})
	})
}

func TestCommandReviewedNum(t *testing.T) {
	convey.Convey("reviewedNum", t, func(ctx convey.C) {
		var (
			reviewers     = []string{"a", "b"}
			reviewedUsers = []string{"c"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num := reviewedNum(reviewers, reviewedUsers)
			ctx.Convey("Then num should not be nil.", func(ctx convey.C) {
				ctx.So(num, convey.ShouldEqual, 0)
			})
		})
	})
	convey.Convey("reviewedNum", t, func(ctx convey.C) {
		var (
			reviewers     = []string{"a", "b"}
			reviewedUsers = []string{"a"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			num := reviewedNum(reviewers, reviewedUsers)
			ctx.Convey("Then num should not be nil.", func(ctx convey.C) {
				ctx.So(num, convey.ShouldEqual, 1)
			})
		})
	})
}
