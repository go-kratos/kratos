package dao

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"

	"go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/app/service/live/dao-anchor/model"
)

func TestDaoNormalizeRoomIDs(t *testing.T) {
	var (
		c           = context.TODO()
		inputIDs    = []int64{5910, 5901, 63, 53, 5010, 115, 666}
		turnedFlags = []bool{false, false, true, true, false, true, true}
	)

	convey.Convey("When normalize a given list of ids", t, func(ctx convey.C) {
		normalized, err := d.dbNormalizeRoomIDs(c, inputIDs)
		ctx.Convey("Then short-id is turned into room-id while room-id keeps untouched", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(normalized), convey.ShouldEqual, len(inputIDs))
			for i, turned := range turnedFlags {
				if turned {
					ctx.So(normalized[i] > inputIDs[i], convey.ShouldBeTrue)
				} else {
					ctx.So(normalized[i] == inputIDs[i], convey.ShouldBeTrue)
				}
			}

			// ctx.So(d.shortIDMapping.caches[0].Len(), convey.ShouldEqual, len(inputIDs))
		})

		ctx.Convey("Then order of result will be preserved", func(ctx convey.C) {
			roomIds := []int64{5901, 57796, 5010, 1011}
			results, err := d.dbNormalizeRoomIDs(c, roomIds)
			ctx.So(err, convey.ShouldBeNil)
			for i := range results {
				ctx.So(roomIds[i], convey.ShouldEqual, results[i])
			}
		})
	})
}

func TestDaoFetchAreas(t *testing.T) {
	var (
		c = context.TODO()
	)

	convey.Convey("When given a valid main area id", t, func(ctx convey.C) {
		req := &v1.FetchAreasReq{
			AreaId: 3,
		}
		resp, err := d.fetchAreas(c, req)
		ctx.Convey("Then we will get a list of its subarea's info along with its info", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(resp.Info.AreaId, convey.ShouldEqual, req.AreaId)
			ctx.So(resp.Info.AreaName, convey.ShouldNotBeEmpty)
			ctx.So(resp.Areas, convey.ShouldNotBeEmpty)
		})
	})

	convey.Convey("When given a non-existed main area id", t, func(ctx convey.C) {
		req := &v1.FetchAreasReq{
			AreaId: 999,
		}
		_, err := d.fetchAreas(c, req)
		ctx.Convey("Then we will get nothing", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoFetchAnchorInfo(t *testing.T) {
	var (
		c = context.TODO()
	)

	convey.Convey("When to fetch anchor info for a given uid", t, func(ctx convey.C) {
		uid := []int64{2}
		const RoomID = 1024
		resp := make(map[int64]*v1.RoomData)
		err := d.dbFetchAnchorInfo(c, uid, resp, false)
		ctx.Convey("Then we will get room data for the anchor", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(resp, convey.ShouldNotBeEmpty)

			data, ok := resp[RoomID]

			ctx.So(ok, convey.ShouldBeTrue)

			alv := data.AnchorLevel
			ctx.So(alv.MaxLevel, convey.ShouldEqual, model.MaxAnchorLevel)
			ctx.So(alv.Level, convey.ShouldEqual, 1)
			ctx.So(alv.Color, convey.ShouldEqual, 0)
			ctx.So(alv.Left, convey.ShouldEqual, 0)
			ctx.So(alv.Right, convey.ShouldEqual, 49)
		})
	})
}
