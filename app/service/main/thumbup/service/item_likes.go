package service

import (
	"context"

	"go-common/app/service/main/thumbup/model"
	"go-common/library/stat/prom"
	xtime "go-common/library/time"
)

// ItemLikes .
func (s *Service) ItemLikes(c context.Context, business string, originID, messageID int64, pn, ps int, mid int64) (res []*model.UserLikeRecord, err error) {
	return s.itemLikes(c, business, originID, messageID, pn, ps, model.StateLike, mid)
}

// ItemDislikes .
func (s *Service) ItemDislikes(c context.Context, business string, originID, messageID int64, pn, ps int, mid int64) (res []*model.UserLikeRecord, err error) {
	return s.itemLikes(c, business, originID, messageID, pn, ps, model.StateDislike, mid)
}

func (s *Service) itemLikes(c context.Context, business string, originID, messageID int64, pn, ps int, state int8, mid int64) (res []*model.UserLikeRecord, err error) {
	var businessID int64
	if businessID, err = s.CheckBusinessOrigin(business, originID); err != nil {
		return
	}
	if err = s.checkItemLikeType(businessID, state); err != nil {
		return
	}
	var (
		start = (pn - 1) * ps
		end   = start + ps - 1 // from cache, end-1
	)
	res, _ = s.dao.CacheItemLikeList(c, messageID, businessID, state, start, end)
	if len(res) != 0 {
		prom.CacheHit.Incr("itemLikeList")
	} else {
		s.dbus.Do(c, func(c context.Context) {
			s.dao.PubItemMsg(c, business, originID, messageID, state)
		})
		prom.CacheMiss.Incr("itemLikeList")
		res, err = s.dao.RawItemLikeList(c, messageID, businessID, originID, state, start, end)
		if err != nil {
			return
		}
	}
	for i, r := range res {
		if r.Mid == mid && len(res) > i+1 {
			res = res[i+1:]
			return
		}
	}
	return
}

// UpdateUpMids .
func (s *Service) UpdateUpMids(c context.Context, business string, data []*model.UpMidsReq) (err error) {
	var businessID int64
	if businessID, err = s.CheckBusiness(business); err != nil {
		return
	}
	for _, m := range data {
		if _, err = s.CheckBusinessOrigin(business, m.OriginID); err != nil {
			return
		}
	}
	prom.BusinessInfoCount.Add("update-mids-raw-"+business, int64(len(data)))
	rows, err := s.dao.UpdateUpMids(c, businessID, data)
	prom.BusinessInfoCount.Add("update-mids-affect-"+business, rows)
	return
}

// ItemHasLike item是否被mids点赞过 返回mid/点赞时间戳
func (s *Service) ItemHasLike(c context.Context, business string, originID, messageID int64, mids []int64) (res map[int64]*model.UserLikeRecord, err error) {
	var businessID int64
	var likes map[int64]int64
	if businessID, err = s.CheckBusinessOrigin(business, originID); err != nil {
		return
	}
	state := int8(model.StateLike)
	if err = s.checkItemLikeType(businessID, state); err != nil {
		return
	}
	var exist bool
	if exist, _ = s.dao.ExpireItemLikesCache(c, messageID, businessID, state); exist {
		if likes, err = s.dao.ItemLikeExists(c, messageID, businessID, mids, state); err != nil {
			exist = false
			err = nil
		}
	}
	if !exist {
		likes, err = s.dao.ItemHasLike(c, businessID, originID, messageID, mids, state)
		s.dbus.Do(c, func(c context.Context) {
			s.dao.PubItemMsg(c, business, originID, messageID, state)
		})
	}
	res = make(map[int64]*model.UserLikeRecord)
	for mid, t := range likes {
		res[mid] = &model.UserLikeRecord{
			Mid:  mid,
			Time: xtime.Time(t),
		}
	}
	return
}
