package service

import (
	"context"
	"fmt"

	pb "go-common/app/service/main/thumbup/api"
	"go-common/app/service/main/thumbup/dao"
	"go-common/app/service/main/thumbup/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// 用户点赞列表数据量少的话 全取比较好 点赞数量多的话 同时一次性取的少的话 zscore比较好
var _rangeLimit = 20 //经验值

// HasLike .
func (s *Service) HasLike(c context.Context, business string, mid int64, messageIDs []int64) (res map[int64]int8, all map[int64]*pb.UserLikeState, err error) {
	if len(messageIDs) == 0 {
		return
	}
	var businessID int64
	if businessID, err = s.CheckBusiness(business); err != nil {
		return
	}
	if mid == 0 {
		return
	}
	var (
		likes        = make(map[int64]*pb.UserLikeState)
		dislikes     = make(map[int64]*pb.UserLikeState)
		messageIDMap = make(map[int64]bool)
	)
	for _, id := range messageIDs {
		messageIDMap[id] = true
	}
	group := errgroup.Group{}
	limit := s.dao.BusinessIDMap[businessID].UserLikesLimit
	if s.dao.BusinessIDMap[businessID].EnableUserLikeList() {
		group.Go(func() (err error) {
			exist, _ := s.dao.ExpireUserLikesCache(c, mid, businessID, model.StateLike)
			if exist && len(messageIDs) <= _rangeLimit {
				likes, err = s.dao.UserLikeExists(c, mid, businessID, messageIDs, model.StateLike)
				return
			}
			var items []*model.ItemLikeRecord
			if items, err = s.dao.UserLikeList(c, mid, businessID, model.StateLike, 0, limit); err != nil {
				return err
			}
			s.dbus.Do(c, func(c context.Context) {
				s.dao.PubUserMsg(c, business, mid, model.StateLike)
			})
			for _, item := range items {
				if messageIDMap[item.MessageID] {
					likes[item.MessageID] = &pb.UserLikeState{
						Mid:   mid,
						Time:  item.Time,
						State: pb.State_STATE_LIKE,
					}
				}
			}
			return
		})
	}
	if s.dao.BusinessIDMap[businessID].EnableUserDislikeList() {
		group.Go(func() (err error) {
			exist, _ := s.dao.ExpireUserLikesCache(c, mid, businessID, model.StateDislike)
			if exist && len(messageIDs) <= _rangeLimit {
				dislikes, err = s.dao.UserLikeExists(c, mid, businessID, messageIDs, model.StateDislike)
				return
			}
			var items []*model.ItemLikeRecord
			if items, err = s.dao.UserLikeList(c, mid, businessID, model.StateDislike, 0, limit); err != nil {
				return err
			}
			s.dbus.Do(c, func(c context.Context) {
				s.dao.PubUserMsg(c, business, mid, model.StateDislike)
			})
			for _, item := range items {
				if messageIDMap[item.MessageID] {
					dislikes[item.MessageID] = &pb.UserLikeState{
						Mid:   mid,
						Time:  item.Time,
						State: pb.State_STATE_DISLIKE,
					}
				}
			}
			return
		})
	}
	err = group.Wait()
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("cache.Exists(%v) err: %+v", messageIDs, err)), log.KV("mid", mid), log.KV("business", business))
		dao.PromError("has_like:Exists")
	}
	res = make(map[int64]int8)
	all = make(map[int64]*pb.UserLikeState)
	for id, v := range likes {
		if v != nil {
			res[id] = model.StateLike
			all[id] = v
		}
	}
	for id, v := range dislikes {
		if v != nil {
			res[id] = model.StateDislike
			all[id] = v
		}
	}
	if len(messageIDs) == 1 {
		log.Info("has_like business: %s mid: %v message_id: %v, resp: %+v, err: %v", business, mid, messageIDs[0], res, err)
	}
	return
}
