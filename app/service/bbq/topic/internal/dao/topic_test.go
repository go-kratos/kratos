package dao

import (
	"context"
	"fmt"
	"go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRawTopicInfo(t *testing.T) {
	convey.Convey("RawTopicInfo", t, func(convCtx convey.C) {
		var (
			ctx      = context.Background()
			topicIDs = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.RawTopicInfo(ctx, topicIDs)
			log.Infow(ctx, "topics", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCacheTopicInfo(t *testing.T) {
	convey.Convey("CacheTopicInfo", t, func(convCtx convey.C) {
		var (
			ctx      = context.Background()
			topicIDs = []int64{1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := d.CacheTopicInfo(ctx, topicIDs)
			log.Infow(ctx, "res", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAddCacheTopicInfo(t *testing.T) {
	convey.Convey("AddCacheTopicInfo", t, func(convCtx convey.C) {
		var (
			ctx        = context.Background()
			topicInfos map[int64]*api.TopicInfo
		)
		topicInfos = make(map[int64]*api.TopicInfo)
		topicInfos[1] = &api.TopicInfo{TopicId: 1, Name: "Test", State: 0, Desc: "test for tester"}
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.AddCacheTopicInfo(ctx, topicInfos)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelCacheTopicInfo(t *testing.T) {
	convey.Convey("DelCacheTopicInfo", t, func(convCtx convey.C) {
		var (
			ctx     = context.Background()
			topicID = int64(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			d.DelCacheTopicInfo(ctx, topicID)
			res, _ := d.CacheTopicInfo(ctx, []int64{topicID})
			topicInfo := res[topicID]
			convCtx.Convey("No return values", func(convCtx convey.C) {
				convCtx.So(topicInfo, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoInsertTopics(t *testing.T) {
	convey.Convey("InsertTopics", t, func(convCtx convey.C) {
		var (
			ctx    = context.Background()
			topics map[string]*api.TopicInfo
		)
		topics = make(map[string]*api.TopicInfo)
		//topicName := fmt.Sprintf("test_%d", rand.Int()%10000000)
		topicName := "Test"
		topics[topicName] = &api.TopicInfo{Name: topicName, Score: float64(rand.Int()%10000) / float64(10000)}
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			_, err := d.InsertTopics(ctx, topics)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
		log.Infow(ctx, "log", "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx:"+topicName)

		topicName = fmt.Sprintf("test_%d", rand.Int()%10000000)
		topics[topicName] = &api.TopicInfo{Name: topicName, Score: float64(rand.Int()%10000) / float64(10000)}
		topicName = fmt.Sprintf("test_%d", rand.Int()%10000000)
		topics[topicName] = &api.TopicInfo{Name: topicName, Score: float64(rand.Int()%10000) / float64(10000)}
		convCtx.Convey("multi insert", func(convCtx convey.C) {
			_, err := d.InsertTopics(ctx, topics)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})

		longTopics := make(map[string]*api.TopicInfo)
		longName := "test_toolonggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggggg"
		longTopics[longName] = &api.TopicInfo{Name: longName}
		convCtx.Convey("error case", func(convCtx convey.C) {
			_, duplicateErr := d.InsertTopics(ctx, longTopics)
			convCtx.Convey("Then err should not be nil.", func(convCtx convey.C) {
				convCtx.So(duplicateErr, convey.ShouldAlmostEqual, ecode.TopicNameLenErr)
			})
		})

	})
}

func TestDaoTopicID(t *testing.T) {
	convey.Convey("TopicID", t, func(convCtx convey.C) {
		var (
			ctx   = context.Background()
			names = []string{"Test"}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			topics, err := d.TopicID(ctx, names)
			log.Infow(ctx, "names", names, "topics", topics)
			convCtx.Convey("Then err should be nil.topics should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(topics, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateTopic(t *testing.T) {
	convey.Convey("TopicID", t, func(convCtx convey.C) {
		convCtx.Convey("update desc", func(convCtx convey.C) {
			var (
				ctx     = context.Background()
				topicID = int64(1)
				field   = "desc"
				value   = "update_desc"
			)
			origin, _ := d.TopicInfo(ctx, []int64{topicID})
			originTopic := origin[topicID]
			err := d.UpdateTopic(ctx, topicID, field, value)
			curr, _ := d.TopicInfo(ctx, []int64{topicID})
			currTopic := curr[topicID]
			d.UpdateTopic(ctx, topicID, field, originTopic.Desc)
			convCtx.Convey("Then err should be nil.topics should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(currTopic.Desc, convey.ShouldEqual, value)
			})
		})
		convCtx.Convey("update state", func(convCtx convey.C) {
			var (
				ctx     = context.Background()
				topicID = int64(1)
				field   = "state"
				value   = model.TopicStateUnavailable
			)
			err := d.UpdateTopic(ctx, topicID, field, value)
			curr, _ := d.TopicInfo(ctx, []int64{topicID})
			currTopic := curr[topicID]
			d.UpdateTopic(ctx, topicID, field, model.TopicStateAvailable)
			convCtx.Convey("Then err should be nil.topics should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(currTopic.State, convey.ShouldEqual, value)
			})
		})
	})
}

func TestDaoListUnAvailableTopics(t *testing.T) {
	convey.Convey("ListUnAvailableTopics", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			page = int32(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			list, hasMore, err := d.ListUnAvailableTopics(ctx, page, model.CmsTopicSize)
			log.Infow(ctx, "list", list, "has_more", hasMore)
			convCtx.Convey("Then err should be nil.list,hasMore should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListRankTopics(t *testing.T) {
	convey.Convey("ListRankTopics", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			page = int32(1)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			list, hasMore, err := d.ListRankTopics(ctx, page, model.DiscoveryTopicSize)
			log.Infow(ctx, "list", list, "topics", hasMore)
			convCtx.Convey("Then err should be nil.list,hasMore should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(hasMore, convey.ShouldBeTrue)
				convCtx.So(list, convey.ShouldNotBeNil)
			})
		})
		convCtx.Convey("stick test", func(convCtx convey.C) {
			originStickList, _ := d.GetStickTopic(ctx)
			d.setStickTopic(ctx, []int64{111111110, 111111111, 111111112})
			list, hasMore, err := d.ListRankTopics(ctx, 1, model.DiscoveryTopicSize)
			log.Infow(ctx, "list", list, "has_more", hasMore)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(list[0], convey.ShouldEqual, 111111110)
			convCtx.So(list[1], convey.ShouldEqual, 111111111)
			convCtx.So(list[2], convey.ShouldEqual, 111111112)
			convCtx.So(len(list), convey.ShouldEqual, 3+model.DiscoveryTopicSize)
			// 恢复原来的置顶话题
			if len(originStickList) > 0 {
				d.setStickTopic(ctx, originStickList)
			}
		})
	})
}

func TestDaogetStickTopic(t *testing.T) {
	convey.Convey("GetStickTopic", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			list, err := d.GetStickTopic(ctx)
			convCtx.Convey("Then err should be nil.list should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(list, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoStickTopic(t *testing.T) {
	convey.Convey("StickTopic", t, func(convCtx convey.C) {
		var (
			ctx       = context.Background()
			opTopicID = int64(1)
			op        = int64(1)
		)
		originStickList, _ := d.GetStickTopic(ctx)

		convCtx.Convey("common stick operate", func(convCtx convey.C) {
			err := d.StickTopic(ctx, opTopicID, op)
			newStickList, _ := d.GetStickTopic(ctx)
			log.V(1).Infow(ctx, "new_stick_list", newStickList)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(newStickList[0], convey.ShouldEqual, 1)
		})

		convCtx.Convey("common cancel stick operate", func(convCtx convey.C) {
			err := d.StickTopic(ctx, opTopicID, 0)
			newCancelStickList, _ := d.GetStickTopic(ctx)
			log.V(1).Infow(ctx, "new_cancel_stick_list", newCancelStickList)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(newCancelStickList[0], convey.ShouldNotEqual, 1)
		})

		convCtx.Convey("stick num test", func(convCtx convey.C) {
			for i := 1; i < model.MaxStickTopicNum+3; i++ {
				d.StickTopic(ctx, int64(i), 1)
			}
			list, _ := d.GetStickTopic(ctx)
			log.V(1).Infow(ctx, "list", list)
			convCtx.So(len(list), convey.ShouldEqual, model.MaxStickTopicNum)
		})

		// 恢复原来的置顶话题
		if len(originStickList) > 0 {
			d.setStickTopic(ctx, originStickList)
		}
	})
}
