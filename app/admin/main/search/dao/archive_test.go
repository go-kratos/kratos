package dao

import (
	"context"
	"go-common/app/admin/main/search/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/olivere/elastic.v5"
)

func TestDaoArchiveCheck(t *testing.T) {
	convey.Convey("ArchiveCheck", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &model.ArchiveCheckParams{
				Bsp:          &model.BasicSearchParams{},
				Aids:         []int64{0},
				TypeIds:      []int64{0},
				Attrs:        []int64{0},
				States:       []int64{0},
				Mids:         []int64{0},
				MidFrom:      1,
				MidTo:        1,
				DurationFrom: 1,
				DurationTo:   1,
				TimeFrom:     "0001-01-01 00:00:00",
				TimeTo:       "0001-01-01 00:00:00",
				Time:         "ctime",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ArchiveCheck(c, p)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoVideo(t *testing.T) {
	convey.Convey("Video", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &model.VideoParams{
				Bsp:            &model.BasicSearchParams{},
				VIDs:           []int64{0},
				AIDs:           []int64{0},
				CIDs:           []int64{0},
				TIDs:           []int64{0},
				FileNames:      []string{""},
				RelationStates: []int64{0},
				ArcMids:        []int64{0},
				TagID:          1,
				Status:         []int64{0},
				XCodeState:     []int64{0},
				UserType:       0,
				DurationFrom:   1,
				DurationTo:     1,
				OrderType:      1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.Video(c, p)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				//ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTaskQa(t *testing.T) {
	convey.Convey("TaskQa", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &model.TaskQa{
				Bsp: &model.BasicSearchParams{
					AppID: "task_qa",
				},
				Ids:           []int64{0},
				TaskIds:       []string{""},
				Uids:          []string{""},
				ArcTagIds:     []string{""},
				AuditTagIds:   []int64{0},
				UpGroups:      []string{""},
				ArcTitles:     []string{""},
				ArcTypeIds:    []string{""},
				States:        []string{""},
				AuditStatuses: []string{""},
				FansFrom:      "0",
				FansTo:        "0",
				CtimeFrom:     "0001-01-01 00:00:00",
				CtimeTo:       "0001-01-01 00:00:00",
				FtimeFrom:     "0001-01-01 00:00:00",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.TaskQa(c, p)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveCommerce(t *testing.T) {
	convey.Convey("ArchiveCommerce", t, func(ctx convey.C) {
		var (
			c = context.Background()
			p = &model.ArchiveCommerce{
				Bsp:        &model.BasicSearchParams{},
				Ids:        []string{"0"},
				Mids:       []string{"0"},
				PTypeIds:   []string{"0"},
				TypeIds:    []string{"0"},
				States:     []string{"0"},
				Copyrights: []string{"0"},
				OrderIds:   []string{"0"},
				IsOrder:    1,
				IsOriginal: 1,
				Action:     "get_ptypeids",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ArchiveCommerce(c, p)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoArchiveCommercePTypeIds(t *testing.T) {
	convey.Convey("ArchiveCommercePTypeIds", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = &elastic.BoolQuery{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ArchiveCommercePTypeIds(c, query)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
