package dao

import (
	"context"
	"encoding/json"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/log"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoInsertTopicVideo(t *testing.T) {
	convey.Convey("InsertTopicVideo", t, func(convCtx convey.C) {
		var (
			ctx      = context.Background()
			svid     = rand.Int63() % 1000000
			topicIDs = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			rowsAffected, err := d.InsertTopicVideo(ctx, svid, topicIDs)
			convCtx.Convey("Then err should be nil.rowsAffected should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(rowsAffected, convey.ShouldEqual, 1)
			})
		})
	})
}

func TestDaoUpdateVideoScore(t *testing.T) {
	convey.Convey("UpdateVideoScore", t, func(convCtx convey.C) {
		var (
			ctx   = context.Background()
			svid  = int64(1)
			score = float64(1.0)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateVideoScore(ctx, svid, score)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoUpdateVideoState(t *testing.T) {
	convey.Convey("UpdateVideoState", t, func(convCtx convey.C) {
		var (
			ctx   = context.Background()
			svid  = int64(1)
			state = int32(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.UpdateVideoState(ctx, svid, state)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoListTopicVideos(t *testing.T) {
	convey.Convey("ListTopicVideos", t, func(convCtx convey.C) {
		var (
			ctx     = context.Background()
			topicID = int64(1)
		)
		res, hasMore, err := d.ListTopicVideos(ctx, topicID, "", "", model.TopicVideoSize)
		log.V(1).Infow(ctx, "res", res)
		convCtx.Convey("Then err should be nil.res,hasMore should not be nil.", func(convCtx convey.C) {
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(hasMore, convey.ShouldNotBeNil)
			convCtx.So(res, convey.ShouldNotBeNil)
		})

		originStickList, _ := d.GetStickTopicVideo(ctx, topicID)
		newStickList := []int64{1, 2, 3, 4, 5, 6}
		d.SetStickTopicVideo(ctx, topicID, newStickList)
		var data []byte

		convCtx.Convey("cursor_in_rank && direction_next", func(convCtx convey.C) {
			data, _ = json.Marshal(model.CursorValue{Offset: 0, StickRank: 2})
			res, hasMore, err = d.ListTopicVideos(ctx, topicID, "", string(data), model.TopicVideoSize)
			log.V(1).Infow(ctx, "res", res, "has_more", hasMore)

			convCtx.So(res[0].Svid, convey.ShouldEqual, 3)
			convCtx.So(res[1].Svid, convey.ShouldEqual, 4)
			convCtx.So(res[2].Svid, convey.ShouldEqual, 5)
			convCtx.So(res[3].Svid, convey.ShouldEqual, 6)
			// 检验cursor值是否符合要求
			var unmarshalCursor model.CursorValue
			json.Unmarshal([]byte(res[3].CursorValue), &unmarshalCursor)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(unmarshalCursor.StickRank, convey.ShouldEqual, 6)
			convCtx.So(unmarshalCursor.Offset, convey.ShouldEqual, 0)
			json.Unmarshal([]byte(res[4].CursorValue), &unmarshalCursor)
			convCtx.So(unmarshalCursor.StickRank, convey.ShouldEqual, 0)
			convCtx.So(unmarshalCursor.Offset, convey.ShouldEqual, 1)

			convCtx.So(len(res), convey.ShouldEqual, model.TopicVideoSize+4)
			convCtx.So(hasMore, convey.ShouldBeTrue)
			convCtx.So(res, convey.ShouldNotBeNil)
		})

		convCtx.Convey("cursor_in_rank && direction_prev", func(convCtx convey.C) {
			data, _ = json.Marshal(model.CursorValue{Offset: 0, StickRank: 4})
			res, hasMore, err = d.ListTopicVideos(ctx, topicID, string(data), "", model.TopicVideoSize)
			log.V(1).Infow(ctx, "res", res, "has_more", hasMore)

			convCtx.So(res[0].Svid, convey.ShouldEqual, 3)
			convCtx.So(res[1].Svid, convey.ShouldEqual, 2)
			convCtx.So(res[2].Svid, convey.ShouldEqual, 1)
			convCtx.So(len(res), convey.ShouldEqual, 3)
			convCtx.So(hasMore, convey.ShouldBeFalse)

			// 边缘情况，选择了rank=1的视频
			data, _ = json.Marshal(model.CursorValue{Offset: 0, StickRank: 1})
			res, hasMore, err = d.ListTopicVideos(ctx, topicID, string(data), "", model.TopicVideoSize)
			log.V(1).Infow(ctx, "res", res, "has_more", hasMore)

			convCtx.So(len(res), convey.ShouldEqual, 0)
			convCtx.So(hasMore, convey.ShouldBeFalse)
		})

		convCtx.Convey("direction_prev", func(convCtx convey.C) {
			data, _ = json.Marshal(model.CursorValue{Offset: 1, StickRank: 0})
			res, hasMore, err = d.ListTopicVideos(ctx, topicID, string(data), "", model.TopicVideoSize)
			log.V(1).Infow(ctx, "res", res, "has_more", hasMore)

			convCtx.So(len(res), convey.ShouldEqual, len(newStickList))
			convCtx.So(hasMore, convey.ShouldBeFalse)

			data, _ = json.Marshal(model.CursorValue{Offset: 5, StickRank: 0})
			res, hasMore, err = d.ListTopicVideos(ctx, topicID, string(data), "", model.TopicVideoSize)
			log.V(1).Infow(ctx, "res", res, "has_more", hasMore)
			// 检验cursor值是否符合要求
			var unmarshalCursor model.CursorValue
			json.Unmarshal([]byte(res[0].CursorValue), &unmarshalCursor)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(unmarshalCursor.StickRank, convey.ShouldEqual, 0)
			convCtx.So(unmarshalCursor.Offset, convey.ShouldEqual, 4)
			json.Unmarshal([]byte(res[1].CursorValue), &unmarshalCursor)
			convCtx.So(unmarshalCursor.StickRank, convey.ShouldEqual, 0)
			convCtx.So(unmarshalCursor.Offset, convey.ShouldEqual, 3)
		})

		convCtx.Convey("direction_next", func(convCtx convey.C) {
			data, _ = json.Marshal(model.CursorValue{Offset: 1, StickRank: 0})
			res, hasMore, err = d.ListTopicVideos(ctx, topicID, "", string(data), model.TopicVideoSize)
			log.V(1).Infow(ctx, "res", res, "has_more", hasMore)

			convCtx.So(len(res), convey.ShouldEqual, int(model.TopicVideoSize))
			convCtx.So(hasMore, convey.ShouldBeTrue)

			data, _ = json.Marshal(model.CursorValue{Offset: 5, StickRank: 0})
			res, hasMore, err = d.ListTopicVideos(ctx, topicID, "", string(data), model.TopicVideoSize)
			log.V(1).Infow(ctx, "res", res, "has_more", hasMore)
			// 检验cursor值是否符合要求
			var unmarshalCursor model.CursorValue
			json.Unmarshal([]byte(res[0].CursorValue), &unmarshalCursor)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(unmarshalCursor.StickRank, convey.ShouldEqual, 0)
			convCtx.So(unmarshalCursor.Offset, convey.ShouldEqual, 6)
			json.Unmarshal([]byte(res[1].CursorValue), &unmarshalCursor)
			convCtx.So(unmarshalCursor.StickRank, convey.ShouldEqual, 0)
			convCtx.So(unmarshalCursor.Offset, convey.ShouldEqual, 7)

		})

		// 恢复原来的置顶话题
		if len(originStickList) > 0 {
			d.SetStickTopicVideo(ctx, topicID, originStickList)
		}

	})
}

func TestDaogetStickTopicVideo(t *testing.T) {
	convey.Convey("GetStickTopicVideo", t, func(convCtx convey.C) {
		var (
			ctx     = context.Background()
			topicID = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			list, err := d.GetStickTopicVideo(ctx, topicID)
			log.V(1).Infow(ctx, "list", list)
			convCtx.Convey("Then err should be nil.list should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStickTopicVideo(t *testing.T) {
	convey.Convey("StickTopicVideo", t, func(convCtx convey.C) {
		var (
			ctx       = context.Background()
			opTopicID = int64(1)
			opSvid    = int64(1)
			op        = int64(1)
		)

		originStickList, _ := d.GetStickTopicVideo(ctx, opTopicID)

		convCtx.Convey("common stick operate", func(convCtx convey.C) {
			err := d.StickTopicVideo(ctx, opTopicID, opSvid, op)
			newStickList, _ := d.GetStickTopicVideo(ctx, opTopicID)
			log.V(1).Infow(ctx, "new_stick_list", newStickList)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(newStickList[0], convey.ShouldEqual, 1)
		})

		convCtx.Convey("common cancel stick operate", func(convCtx convey.C) {
			err := d.StickTopicVideo(ctx, opTopicID, opSvid, 0)
			newCancelStickList, _ := d.GetStickTopicVideo(ctx, opTopicID)
			log.V(1).Infow(ctx, "new_cancel_stick_list", newCancelStickList)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(newCancelStickList[0], convey.ShouldNotEqual, 1)
		})

		convCtx.Convey("stick num test", func(convCtx convey.C) {
			for i := 1; i < model.MaxStickTopicVideoNum+3; i++ {
				d.StickTopicVideo(ctx, opTopicID, int64(i), 1)
			}
			list, _ := d.GetStickTopicVideo(ctx, opTopicID)
			log.V(1).Infow(ctx, "list", list)
			convCtx.So(len(list), convey.ShouldEqual, model.MaxStickTopicVideoNum)
		})

		// 恢复原来的置顶话题
		if len(originStickList) > 0 {
			d.SetStickTopicVideo(ctx, opTopicID, originStickList)
		}

	})
}

func TestDaoGetVideoTopic(t *testing.T) {
	convey.Convey("GetVideoTopic", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			svid = int64(1547635456050324977)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.GetVideoTopic(ctx, svid)
			data, _ := json.Marshal(res)
			log.V(1).Infow(ctx, "res", string(data))
			convCtx.Convey("Then err should be nil.list should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
