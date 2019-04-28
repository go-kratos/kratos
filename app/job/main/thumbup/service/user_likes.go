package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/thumbup/model"
	xmdl "go-common/app/service/main/thumbup/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func newUserLikeMsg(msg *databus.Message) (res interface{}, err error) {
	userLikesMsg := new(xmdl.UserMsg)
	if err = json.Unmarshal(msg.Value, &userLikesMsg); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
		return
	}
	log.Info("get user like msg: %+v", userLikesMsg)
	res = userLikesMsg
	return
}

func userLikesSplit(msg *databus.Message, data interface{}) int {
	um, ok := data.(*xmdl.UserMsg)
	if !ok {
		log.Error("user like msg err: mid: 0 %s", msg.Value)
		return 0
	}
	return int(um.Mid)
}

func (s *Service) userLikesDo(ms []interface{}) {
	for _, m := range ms {
		um, ok := m.(*xmdl.UserMsg)
		if !ok {
			continue
		}
		var (
			err        error
			businessID int64
			ctx        = context.Background()
		)
		if businessID, err = s.checkBusiness(um.Business); err != nil {
			log.Warn("userlikes: checkBusiness(%s) err:%+v", um.Business, err)
			continue
		}
		var exist bool
		if exist, _ = s.dao.ExpireUserLikesCache(ctx, um.Mid, businessID, um.State); exist {
			log.Warn("userlikes: ExpireUserLikesCache(%+v) exist ignore", um, err)
			continue
		}
		for i := 0; i < _retryTimes; i++ {
			if err = s.addUserLikesCache(context.Background(), um); err == nil {
				break
			}
		}
		if err != nil {
			log.Error("userLikes fail params(%+v) err: %+v", m, err)
		} else {
			log.Info("userLikes success params(%+v)", m)
		}
	}
}

// addUserLikesCache .
func (s *Service) addUserLikesCache(c context.Context, p *xmdl.UserMsg) (err error) {
	var businessID int64
	if businessID, err = s.checkBusiness(p.Business); err != nil {
		log.Error("s.checkBusiness business(%s) error(%v)", p.Business, err)
		return
	}
	var items []*model.ItemLikeRecord
	var limit = s.businessIDMap[businessID].UserLikesLimit
	if items, err = s.dao.UserLikes(c, p.Mid, businessID, p.State, limit); err != nil {
		log.Error("s.dao.UserLikes mid(%d) businessID(%d)(%d) type(%d) error(%v)", p.Mid, businessID, p.State, err)
		return
	}
	err = s.dao.AddUserLikesCache(c, p.Mid, businessID, items, p.State, limit)
	return
}

func (s *Service) addUserlikeRecord(c context.Context, mid, businessID int64, state int8, item *model.ItemLikeRecord) (err error) {
	var exist bool
	if exist, err = s.dao.ExpireUserLikesCache(c, mid, businessID, state); (err != nil) || !exist {
		return
	}
	limit := s.businessIDMap[businessID].UserLikesLimit
	err = s.dao.AppendCacheUserLikeList(c, mid, item, businessID, state, limit)
	return
}
