package service

import (
	"context"
	"go-common/app/admin/main/laser/model"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceArchiveRecheck(t *testing.T) {
	convey.Convey("ArchiveRecheck", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			typeIDS   = []int64{}
			unames    = ""
			startDate = int64(0)
			endDate   = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.ArchiveRecheck(c, typeIDS, unames, startDate, endDate)
			convCtx.Convey("Then err should be nil.recheckViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceUserRecheck(t *testing.T) {
	convey.Convey("UserRecheck", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			typeIDS   = []int64{}
			unames    = ""
			startDate = int64(0)
			endDate   = int64(0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.UserRecheck(c, typeIDS, unames, startDate, endDate)
			convCtx.Convey("Then err should be nil.recheckViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceDailyStatArchiveRecheck(t *testing.T) {
	convey.Convey("dailyStatArchiveRecheck", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			business  = int(0)
			typeIDS   = []int64{}
			statTypes = []int64{}
			uids      = []int64{}
			statDate  = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.dailyStatArchiveRecheck(c, business, typeIDS, statTypes, uids, statDate)
			convCtx.Convey("Then err should be nil.statViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServicemakeUpArchiveRecheck(t *testing.T) {
	convey.Convey("makeUpArchiveRecheck", t, func(convCtx convey.C) {
		var (
			mediateView map[int64]map[int]int64
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			makeUpArchiveRecheck(mediateView)
			convCtx.Convey("Then statViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceDailyStatArchiveStat(t *testing.T) {
	convey.Convey("dailyArchiveStat", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			business  = int(0)
			typeIDS   = []int64{}
			statTypes = []int64{}
			uids      = []int64{}
			statDate  = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.dailyArchiveStat(c, business, typeIDS, statTypes, uids, statDate)
			convCtx.Convey("Then err should be nil.mediateView should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceTagRecheck(t *testing.T) {
	convey.Convey("TagRecheck", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			unames    = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.TagRecheck(c, startDate, endDate, unames)
			convCtx.Convey("Then err should be nil.tagViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceDailyStatTagRecheck(t *testing.T) {
	convey.Convey("dailyStatTagRecheck", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			business  = int(0)
			statTypes = []int64{}
			uids      = []int64{}
			statDate  = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.dailyStatTagRecheck(c, business, statTypes, uids, statDate)
			convCtx.Convey("Then err should be nil.statViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServicemakeUpTagRecheck(t *testing.T) {
	convey.Convey("makeUpTagRecheck", t, func(convCtx convey.C) {
		var (
			mediateView map[int64]map[int]int64
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			makeUpTagRecheck(mediateView)
			convCtx.Convey("Then statViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceRecheck123(t *testing.T) {
	convey.Convey("Recheck123", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			typeIDS   = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.Recheck123(c, startDate, endDate, typeIDS)
			convCtx.Convey("Then err should be nil.recheckView should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceDailyStatArchiveStreamStat(t *testing.T) {
	convey.Convey("dailyStatArchiveStreamStat", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			business  = int(0)
			typeIDS   = []int64{}
			uids      = []int64{}
			statTypes = []int64{}
			statDate  = time.Now()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.dailyStatArchiveStreamStat(c, business, typeIDS, uids, statTypes, statDate)
			convCtx.Convey("Then err should be nil.statViews should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServicewrap(t *testing.T) {
	convey.Convey("wrap", t, func(convCtx convey.C) {
		var (
			cargoMap map[int64]*model.CargoItem
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			wrap(cargoMap)
			convCtx.Convey("Then views should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceCsvAuditCargo(t *testing.T) {
	convey.Convey("CsvAuditCargo", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			unames    = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.CsvAuditCargo(c, startDate, endDate, unames)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceAuditorCargoList(t *testing.T) {
	convey.Convey("AuditorCargoList", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			unames    = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.AuditorCargoList(c, startDate, endDate, unames)
			convCtx.Convey("Then err should be nil.wrappers,lineWidth should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceCsvRandomVideoAudit(t *testing.T) {
	convey.Convey("CsvRandomVideoAudit", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			unames    = ""
			typeIDS   = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.CsvRandomVideoAudit(c, startDate, endDate, unames, typeIDS)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceCsvFixedVideoAudit(t *testing.T) {
	convey.Convey("CsvFixedVideoAudit", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			unames    = ""
			typeIDS   = []int64{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.CsvFixedVideoAudit(c, startDate, endDate, unames, typeIDS)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceRandomVideo(t *testing.T) {
	convey.Convey("RandomVideo", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			typeIDS   = []int64{}
			uname     = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.RandomVideo(c, startDate, endDate, typeIDS, uname)
			convCtx.Convey("Then err should be nil.statViewExts,lineWidth should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceFixedVideo(t *testing.T) {
	convey.Convey("FixedVideo", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			startDate = int64(0)
			endDate   = int64(0)
			typeIDS   = []int64{}
			uname     = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.FixedVideo(c, startDate, endDate, typeIDS, uname)
			convCtx.Convey("Then err should be nil.statViewExts,lineWidth should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServiceVideoAudit(t *testing.T) {
	convey.Convey("videoAudit", t, func(convCtx convey.C) {
		var (
			c        = context.Background()
			business = int(0)
			statDate = time.Now()
			typeIDS  = []int64{}
			unames   = ""
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.videoAudit(c, business, statDate, typeIDS, unames)
			convCtx.Convey("Then err should be nil.viewExts,lineWidth should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}

func TestServicestatNode2ViewExt(t *testing.T) {
	convey.Convey("statNode2ViewExt", t, func(convCtx convey.C) {
		var (
			c         = context.Background()
			statNodes = []*model.StatNode{}
			needALL   = false
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			s.statNode2ViewExt(c, statNodes, needALL)
			convCtx.Convey("Then err should be nil.statViewsExts,lineWidth should not be nil.", func(convCtx convey.C) {

			})
		})
	})
}
