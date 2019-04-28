package like

import (
	"context"
	"testing"

	"fmt"

	"github.com/smartystreets/goconvey/convey"
)

func TestLikeVoteLog(t *testing.T) {
	convey.Convey("VoteLog", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			sid   = int64(0)
			aid   = int64(0)
			mid   = int64(0)
			stage = int64(0)
			vote  = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.VoteLog(c, sid, aid, mid, stage, vote)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeNewestSubject(t *testing.T) {
	convey.Convey("NewestSubject", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			typeIDs = []int64{10256}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.NewestSubject(c, typeIDs)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				fmt.Printf("%+v", res)
			})
		})
	})
}

func TestLikeRawActSubject(t *testing.T) {
	convey.Convey("RawActSubject", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.RawActSubject(c, id)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSubjectListMoreSid(t *testing.T) {
	convey.Convey("SubjectListMoreSid", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			minSid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SubjectListMoreSid(c, minSid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestLikeSubjectMaxID(t *testing.T) {
	convey.Convey("SubjectMaxID", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.SubjectMaxID(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
