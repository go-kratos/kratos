package dao

import (
	"context"
	"go-common/app/interface/main/credit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaocaseObtainMIDKey(t *testing.T) {
	convey.Convey("caseObtainMIDKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := caseObtainMIDKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaocaseVoteCIDMIDKey(t *testing.T) {
	convey.Convey("caseVoteCIDMIDKey", t, func(convCtx convey.C) {
		var (
			mid = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := caseVoteCIDMIDKey(mid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCaseObtainMID(t *testing.T) {
	convey.Convey("CaseObtainMID", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			isToday bool
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			cases, err := d.CaseObtainMID(c, mid, isToday)
			convCtx.Convey("Then err should be nil.cases should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(cases, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIsExpiredObtainMID(t *testing.T) {
	convey.Convey("IsExpiredObtainMID", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			isToday bool
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			ok, err := d.IsExpiredObtainMID(c, mid, isToday)
			convCtx.Convey("Then err should be nil.ok should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(ok, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGrantCases(t *testing.T) {
	convey.Convey("GrantCases", t, func(convCtx convey.C) {
		var (
			c = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			mcases, err := d.GrantCases(c)
			convCtx.Convey("Then err should be nil.mcases should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(mcases, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLoadVoteCaseMID(t *testing.T) {
	convey.Convey("LoadVoteCaseMID", t, func(convCtx convey.C) {
		var (
			c       = context.Background()
			mid     = int64(0)
			mcases  map[int64]*model.SimCase
			isToday bool
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.LoadVoteCaseMID(c, mid, mcases, isToday)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetVoteCaseMID(t *testing.T) {
	convey.Convey("SetVoteCaseMID", t, func(convCtx convey.C) {
		var (
			c   = context.Background()
			mid = int64(0)
			sc  = &model.SimCase{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetVoteCaseMID(c, mid, sc)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
