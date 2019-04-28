package service

import (
	"context"
	"testing"

	wkhmdl "go-common/app/interface/main/creative/model/weeklyhonor"
	upgrpc "go-common/app/service/main/up/api/v1"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
)

var (
	mtime        = xtime.Time(wkhmdl.LatestSunday().Unix())
	c            = context.Background()
	rawUpActives = []*upgrpc.UpActivity{
		{Mid: 1, Activity: 4}, {Mid: 2, Activity: 4}, {Mid: 3, Activity: 4}, {Mid: 4, Activity: 3}, {Mid: 5, Activity: 4}, {Mid: 6, Activity: 1}, {Mid: 7, Activity: 1}, {Mid: 24, Activity: 2},
	}
	mockHls = []*wkhmdl.HonorLog{
		// lev:SSR clicked
		{
			ID:    1,
			MID:   1,
			HID:   1,
			MTime: mtime,
		},
		// lev:SR & BlackList clicked
		{
			ID:    2,
			MID:   2,
			HID:   9,
			MTime: mtime,
		},
		// lev:R clicked
		{
			ID:    3,
			MID:   3,
			HID:   17,
			MTime: mtime,
		},
		// lev:A
		{
			ID:    4,
			MID:   4,
			HID:   26,
			MTime: mtime,
		},
		// lev:B clicked
		{
			ID:    5,
			MID:   5,
			HID:   35,
			MTime: mtime,
		},
		// lev:C clicked
		{
			ID:    6,
			MID:   6,
			HID:   46,
			MTime: mtime,
		},
		// lev:D
		{
			ID:    7,
			MID:   7,
			HID:   56,
			MTime: mtime,
		},
		// BlackList clicked
		{
			ID:  8,
			MID: 24,
			HID: 50,
		},
	}
	clickMap = map[int64]int32{
		1:  1,
		2:  1,
		3:  1,
		5:  1,
		6:  1,
		24: 1,
	}
)

func TestServiceSendMsg(t *testing.T) {
	convey.Convey("SendMsg", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			// mock
			s.honDao.MockUpActivesList(rawUpActives, 0, nil)
			s.honDao.MockUpCount(1, nil)
			s.honDao.MockLatestHonorLogs(mockHls, nil)
			s.honDao.MockClickCounts(clickMap, nil)
			s.honDao.MockSendNotify(nil)
			// test
			s.SendMsg()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServiceFlushHonor(t *testing.T) {
	convey.Convey("FlushHonor", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			// mock
			s.honDao.MockUpActivesList(rawUpActives, 0, nil)
			mockStat := wkhmdl.HonorStat{
				Play:      100,
				PlayLastW: 100,
			}
			s.honDao.MockHonorStat(&mockStat, nil)
			// test
			s.FlushHonor()
			ctx.Convey("No return values", func(ctx convey.C) {
			})
		})
	})
}

func TestServiceupsertHonor(t *testing.T) {
	convey.Convey("upsertHonor", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.upsertHonor(c, mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceTestSendMsg(t *testing.T) {
	convey.Convey("TestSendMsg", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mids = []int64{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.TestSendMsg(c, mids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServicefilterMids(t *testing.T) {
	convey.Convey("filterUnActiveMids", t, func(ctx convey.C) {
		var (
			filteredMids = []int64{1, 2, 3, 4, 6}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			s.honDao.MockUpCount(1, nil)
			s.honDao.MockLatestHonorLogs(mockHls, nil)
			s.honDao.MockClickCounts(clickMap, nil)
			mids, err := s.filterMids(c, rawUpActives)
			ctx.Convey("Then err should be nil.mids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				sunday := wkhmdl.LatestSunday()
				if isOddWeek(sunday) {
					filteredMids = append(filteredMids, 7)
				}
				ctx.So(mids, convey.ShouldResemble, filteredMids)
			})
		})
	})
}
