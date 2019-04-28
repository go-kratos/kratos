package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-common/app/service/bbq/topic/api"
	"go-common/app/service/bbq/topic/internal/model"
	"go-common/library/log"
	"math/rand"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceUpdateVideoScore(t *testing.T) {
	convey.Convey("UpdateVideoScore", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			req = &api.UpdateVideoScoreReq{Svid: 1, Score: 0.33}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpdateVideoScore(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateVideoState(t *testing.T) {
	convey.Convey("UpdateVideoState", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			req = &api.UpdateVideoStateReq{Svid: 1, State: 0}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpdateVideoState(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceUpdateTopicDesc(t *testing.T) {
	convey.Convey("UpdateTopicDesc", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			desc = string("ehahahatest")
			req  = &api.TopicInfo{TopicId: 1, Desc: desc}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			origin, _ := s.dao.TopicInfo(ctx, []int64{1})
			originTopic := origin[1]
			res, err := s.UpdateTopicDesc(ctx, req)
			curr, _ := s.dao.TopicInfo(ctx, []int64{1})
			currTopic := curr[1]
			s.UpdateTopicDesc(ctx, originTopic)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(currTopic.Desc, convey.ShouldEqual, desc)
			})
		})
	})
}

func TestServiceUpdateTopicState(t *testing.T) {
	convey.Convey("UpdateTopicState", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			req = &api.TopicInfo{TopicId: 1, State: 1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.UpdateTopicState(ctx, req)
			curr, _ := s.dao.TopicInfo(ctx, []int64{1})
			currTopic := curr[1]
			req.State = 0
			s.UpdateTopicState(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(currTopic.State, convey.ShouldEqual, 1)
			})
		})
	})
}

func TestServiceListCmsTopics(t *testing.T) {
	convey.Convey("ListCmsTopics", t, func(convCtx convey.C) {
		convCtx.Convey("search name", func(convCtx convey.C) {
			var (
				ctx = context.Background()
				req = &api.ListCmsTopicsReq{Page: 1, Name: "Test"}
			)
			res, err := s.ListCmsTopics(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(res.HasMore, convey.ShouldBeFalse)
				convCtx.So(res.List[0].TopicId, convey.ShouldEqual, 1)
			})
		})
		convCtx.Convey("search topic id", func(convCtx convey.C) {
			var (
				ctx = context.Background()
				req = &api.ListCmsTopicsReq{TopicId: 1, Page: 1}
			)
			res, err := s.ListCmsTopics(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(res.HasMore, convey.ShouldBeFalse)
				convCtx.So(res.List[0].TopicId, convey.ShouldEqual, 1)
			})
		})
		convCtx.Convey("search state available", func(convCtx convey.C) {
			var (
				ctx = context.Background()
				req = &api.ListCmsTopicsReq{State: model.TopicStateAvailable, Page: 1}
			)
			res, err := s.ListCmsTopics(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(res.List[0].State, convey.ShouldEqual, model.TopicStateAvailable)
			})
		})
		convCtx.Convey("search state unavailable", func(convCtx convey.C) {
			var (
				ctx = context.Background()
				req = &api.ListCmsTopicsReq{State: model.TopicStateUnavailable, Page: 1}
			)
			res, err := s.ListCmsTopics(ctx, req)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
				convCtx.So(res.List[0].State, convey.ShouldEqual, model.TopicStateUnavailable)
			})
		})
	})
}

func TestServiceListDiscoveryTopics(t *testing.T) {
	convey.Convey("ListDiscoveryTopics", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			req = &api.ListDiscoveryTopicReq{Page: 1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.ListDiscoveryTopics(ctx, req)
			data, _ := json.Marshal(res)
			log.V(1).Infow(ctx, "res", string(data))
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceListTopicVideos(t *testing.T) {
	convey.Convey("ListTopicVideos", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			req = &api.TopicVideosReq{TopicId: 1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.ListTopicVideos(ctx, req)
			data, _ := json.Marshal(res)
			log.V(1).Infow(ctx, "res", string(data))
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

//
//func TestServiceStickTopic(t *testing.T) {
//	convey.Convey("StickTopic", t, func(convCtx convey.C) {
//		var (
//			ctx = context.Background()
//			in  = &api.StickTopicReq{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.StickTopic(ctx, in)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceStickTopicVideo(t *testing.T) {
//	convey.Convey("StickTopicVideo", t, func(convCtx convey.C) {
//		var (
//			ctx = context.Background()
//			in  = &api.StickTopicVideoReq{}
//		)
//		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
//			res, err := s.StickTopicVideo(ctx, in)
//			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
//				convCtx.So(err, convey.ShouldBeNil)
//				convCtx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestServiceregisterTopic(t *testing.T) {
	convey.Convey("registerTopic", t, func(convCtx convey.C) {
		var (
			ctx  = context.Background()
			svid = int64(rand.Int() % 1000000)
			list = []*api.TitleExtraItem{{Name: "Test"}, {Name: fmt.Sprintf("test_%d", rand.Int()%10000000)}}
		)
		log.V(1).Infow(ctx, "list", list)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.registerTopic(ctx, svid, list)
			log.V(1).Infow(ctx, "list", list, "res", res)
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicegetAvailableTopicInfo(t *testing.T) {
	convey.Convey("getAvailableTopicInfo", t, func(convCtx convey.C) {
		var (
			c    = context.Background()
			keys = []int64{2}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.getAvailableTopicInfo(c, keys)
			convCtx.So(err, convey.ShouldBeNil)
			convCtx.So(len(res), convey.ShouldEqual, 0)
		})
	})
}
func TestServiceVideoTopic(t *testing.T) {
	convey.Convey("VideoTopic", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			in  = &api.VideoTopicReq{Svid: 1547635456050324977}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.VideoTopic(ctx, in)
			data, _ := json.Marshal(res)
			log.V(1).Infow(ctx, "res", string(data))
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestServiceSetStickTopicVideo(t *testing.T) {
	convey.Convey("SetStickTopicVideo", t, func(convCtx convey.C) {
		var (
			ctx     = context.Background()
			topicID = int64(1)
			in      = &api.SetStickTopicVideoReq{TopicId: topicID, Svids: []int64{2, 3, 4, 5}}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			originList, _ := s.dao.GetStickTopicVideo(ctx, topicID)
			res, err := s.SetStickTopicVideo(ctx, in)
			currentList, _ := s.dao.GetStickTopicVideo(ctx, topicID)
			data, _ := json.Marshal(currentList)
			log.V(1).Infow(ctx, "res", string(data))
			s.SetStickTopicVideo(ctx, &api.SetStickTopicVideoReq{TopicId: topicID, Svids: originList})
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceListTopics(t *testing.T) {
	convey.Convey("ListTopics", t, func(convCtx convey.C) {
		var (
			ctx = context.Background()
			req = &api.ListTopicsReq{Page: 1}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			res, err := s.ListTopics(ctx, req)
			data, _ := json.Marshal(res)
			log.V(1).Infow(ctx, "res", string(data))
			convCtx.Convey("Then err should be nil.res should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
