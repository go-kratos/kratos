package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/reply-feed/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSlotsMapping(t *testing.T) {
	convey.Convey("SlotsMapping", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			slotsMap, err := d.SlotsMapping(context.Background())
			ctx.Convey("Then err should be nil.slotsMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(slotsMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSlotsStatManager(t *testing.T) {
	convey.Convey("SlotsStatManager", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			s, err := d.SlotsStatManager(context.Background())
			ctx.Convey("Then err should be nil.s should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(s, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoIdleSlot(t *testing.T) {
	convey.Convey("IdleSlot", t, func(ctx convey.C) {
		var (
			count = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			slots, err := d.IdleSlots(context.Background(), count)
			ctx.Convey("Then err should be nil.slots should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(slots, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoCountIdleSlot(t *testing.T) {
	convey.Convey("CountIdleSlot", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.CountIdleSlot(context.Background())
			ctx.Convey("Then err should be nil. count should greater than 0.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestDaoModifyState(t *testing.T) {
	convey.Convey("ModifyState", t, func(ctx convey.C) {
		var (
			name  = ""
			state = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.ModifyState(context.Background(), name, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateSlotsStat(t *testing.T) {
	convey.Convey("UpdateSlotsStat", t, func(ctx convey.C) {
		var (
			name      = ""
			algorithm = ""
			weight    = ""
			slots     = []int64{-1}
			state     = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateSlotsStat(context.Background(), name, algorithm, weight, slots, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestSlotsStatByName(t *testing.T) {
	convey.Convey("SlotsStatByName", t, func(ctx convey.C) {
		var (
			name = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			slots, algorithm, weight, err := d.SlotsStatByName(context.Background(), name)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(slots, convey.ShouldBeNil)
				ctx.So(algorithm, convey.ShouldEqual, "")
				ctx.So(weight, convey.ShouldEqual, "")
			})
		})
	})
}
func TestDaoUpdateWeight(t *testing.T) {
	convey.Convey("UpdateWeight", t, func(ctx convey.C) {
		var (
			name   = ""
			weight = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpdateWeight(context.Background(), name, weight)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpsertStatistics(t *testing.T) {
	convey.Convey("UpsertStatistics", t, func(ctx convey.C) {
		var (
			name = ""
			date = int(0)
			hour = int(0)
			s    = &model.StatisticsStat{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpsertStatistics(context.Background(), name, date, hour, s)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoStatisticsByDate(t *testing.T) {
	convey.Convey("StatisticsByDate", t, func(ctx convey.C) {
		var (
			begin = int64(0)
			end   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			stats, err := d.StatisticsByDate(context.Background(), begin, end)
			ctx.Convey("Then err should be nil.stats should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(stats, convey.ShouldNotBeNil)
			})
		})
	})
}
