package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSagaAuthKey(t *testing.T) {
	convey.Convey("sagaAuthKey", t, func(ctx convey.C) {
		var (
			projID = int(111)
			branch = "test"
			path   = "."
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := sagaAuthKey(projID, branch, path)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldEqual, "saga_auth_111_test_.")
			})
		})
	})
}

func TestDaoSetPathAuthH(t *testing.T) {
	convey.Convey("SetPathAuthH", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			projID    = int(111)
			branch    = "master"
			path      = "."
			owners    = []string{"zhanglin", "wuwei"}
			reviewers = []string{"tangyongqiang", "changhengyuan"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetPathAuthH(c, projID, branch, path, owners, reviewers)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoPathAuthH(t *testing.T) {
	convey.Convey("PathAuthH", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			projID   = int(111)
			branch   = "master"
			path     = "."
			owner    = []string{"zhanglin", "wuwei"}
			reviewer = []string{"tangyongqiang", "changhengyuan"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			owners, reviewers, err := d.PathAuthH(c, projID, branch, path)
			ctx.Convey("Then err should be nil.owners,reviewers should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fmt.Sprint(reviewers), convey.ShouldEqual, fmt.Sprint(reviewer))
				ctx.So(fmt.Sprint(owners), convey.ShouldEqual, fmt.Sprint(owner))
			})
		})
	})
}

func TestDaoDeletePathAuthH(t *testing.T) {
	convey.Convey("DeletePathAuthH", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			projID     = int(111)
			branch     = "master"
			path       = "."
			owners     = []string{"zhanglin", "wuwei"}
			reviewers  = []string{"tangyongqiang", "changhengyuan"}
			path2      = "test"
			owners2    = []string{"zhang", "wu"}
			reviewers2 = []string{"tang", "chang"}
		)

		ctx.Convey("When data not exist", func(ctx convey.C) {
			err := d.DeletePathAuthH(c, projID, branch, path)
			ctx.Convey("delete err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})

		ctx.Convey("When data exist", func(ctx convey.C) {
			err := d.SetPathAuthH(c, projID, branch, path, owners, reviewers)
			ctx.Convey("save path auth.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
			err = d.SetPathAuthH(c, projID, branch, path2, owners2, reviewers2)
			ctx.Convey("save path auth 2.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})

			os, rs, err := d.PathAuthH(c, projID, branch, path)
			ctx.Convey("get path auth", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fmt.Sprint(rs), convey.ShouldEqual, fmt.Sprint(reviewers))
				ctx.So(fmt.Sprint(os), convey.ShouldEqual, fmt.Sprint(owners))
			})
			os, rs, err = d.PathAuthH(c, projID, branch, path2)
			ctx.Convey("get path auth 2.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fmt.Sprint(rs), convey.ShouldEqual, fmt.Sprint(reviewers2))
				ctx.So(fmt.Sprint(os), convey.ShouldEqual, fmt.Sprint(owners2))
			})

			err = d.DeletePathAuthH(c, projID, branch, path)
			ctx.Convey("delete auth path.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})

			os, rs, err = d.PathAuthH(c, projID, branch, path)
			ctx.Convey("get again path auth.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fmt.Sprint(rs), convey.ShouldEqual, "[]")
				ctx.So(fmt.Sprint(os), convey.ShouldEqual, "[]")
				ctx.So(len(rs), convey.ShouldEqual, 0)
				ctx.So(len(os), convey.ShouldEqual, 0)
			})
			os, rs, err = d.PathAuthH(c, projID, branch, path2)
			ctx.Convey("get again path auth 2.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(fmt.Sprint(rs), convey.ShouldEqual, fmt.Sprint(reviewers2))
				ctx.So(fmt.Sprint(os), convey.ShouldEqual, fmt.Sprint(owners2))
			})

		})
	})
}
