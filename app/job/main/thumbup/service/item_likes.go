package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/thumbup/model"
	xmdl "go-common/app/service/main/thumbup/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func newItemLikeMsg(msg *databus.Message) (res interface{}, err error) {
	itemLikesMsg := new(xmdl.ItemMsg)
	if err = json.Unmarshal(msg.Value, &itemLikesMsg); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
		return
	}
	log.Info("get item like msg: %+v", itemLikesMsg)
	res = itemLikesMsg
	return
}

func itemLikesSplit(msg *databus.Message, data interface{}) int {
	im, ok := data.(*xmdl.ItemMsg)
	if !ok {
		log.Error("user item msg err: message_id: 0 %s", msg.Value)
		return 0
	}
	return int(im.MessageID)
}

func (s *Service) itemLikesDo(ms []interface{}) {
	for _, m := range ms {
		im, ok := m.(*xmdl.ItemMsg)
		if !ok {
			continue
		}
		var (
			err        error
			businessID int64
			ctx        = context.Background()
		)
		if businessID, err = s.checkBusinessOrigin(im.Business, im.OriginID); err != nil {
			continue
		}
		var exist bool
		if exist, _ = s.dao.ExpireItemLikesCache(ctx, im.MessageID, businessID, im.State); exist {
			continue
		}
		for i := 0; i < _retryTimes; i++ {
			if err = s.addItemLikesCache(ctx, im); err == nil {
				break
			}
		}
		log.Info("itemLikes success params(%+v)", m)
	}
}

// addCacheItemLikes .
func (s *Service) addItemLikesCache(c context.Context, p *xmdl.ItemMsg) (err error) {
	var businessID int64
	if businessID, err = s.checkBusinessOrigin(p.Business, p.OriginID); err != nil {
		log.Error("s.checkBusinessOrigin business(%s) originID(%s)", p.Business, p.OriginID, err)
		return
	}
	var items []*model.UserLikeRecord
	var limit = s.businessIDMap[businessID].MessageLikesLimit
	if items, err = s.dao.ItemLikes(c, businessID, p.OriginID, p.MessageID, p.State, limit); err != nil {
		log.Error("s.dao.ItemLikes businessID(%d) originID(%d) messageID(%d) type(%d) error(%v)", businessID, p.OriginID, p.MessageID, p.State, err)
		return
	}
	err = s.dao.AddItemLikesCache(c, businessID, p.MessageID, p.State, limit, items)
	return
}

func (s *Service) addItemlikeRecord(c context.Context, businessID, messageID int64, state int8, item *model.UserLikeRecord) (err error) {
	var exist bool
	if exist, err = s.dao.ExpireItemLikesCache(c, messageID, businessID, state); (err != nil) || !exist {
		return
	}
	limit := s.businessIDMap[businessID].MessageLikesLimit
	err = s.dao.AppendCacheItemLikeList(c, messageID, item, businessID, state, limit)
	return
}
