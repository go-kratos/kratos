package gorm

import (
	"testing"

	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoDirectionByTransitionID(t *testing.T) {
	convey.Convey("DirectionByTransitionID", t, func(ctx convey.C) {
		_, err := d.DirectionByTransitionID(cntx, []int64{}, 1, true)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDirectionByID(t *testing.T) {
	convey.Convey("DirectionByID", t, func(ctx convey.C) {
		_, err := d.DirectionByID(cntx, -1)
		ctx.Convey("Then err should be 404", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.NothingFound)
		})
	})
}

func TestDaoDirectionList(t *testing.T) {
	var (
		pm = &net.ListDirectionParam{
			NetID:        1,
			Ps:           20,
			ID:           []int64{1},
			FlowID:       1,
			TransitionID: 1,
			Direction:    1,
		}
	)
	convey.Convey("DirectionList", t, func(ctx convey.C) {
		result, err := d.DirectionList(cntx, pm)
		ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(result, convey.ShouldNotBeNil)

			t.Logf("list result(%+v)", result)
			for _, dir := range result.Result {
				t.Logf("bylist dir(%+v)", dir)
			}
		})
	})
}

func TestDaoDirectionByUnique(t *testing.T) {
	convey.Convey("DirectionByUnique", t, func(ctx convey.C) {
		_, err := d.DirectionByUnique(cntx, 0, 0, 0, 0)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDirectionByFlowID(t *testing.T) {
	convey.Convey("DirectionByFlowID", t, func(ctx convey.C) {
		_, err := d.DirectionByFlowID(cntx, []int64{}, 1)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDirectionByNet(t *testing.T) {
	convey.Convey("DirectionByNet", t, func(ctx convey.C) {
		_, err := d.DirectionByNet(cntx, 0)
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoDirections(t *testing.T) {
	convey.Convey("Directions", t, func(ctx convey.C) {
		_, err := d.Directions(cntx, []int64{0})
		ctx.Convey("Then err should be nil", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
