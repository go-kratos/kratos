package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/thumbup/model"
	xmdl "go-common/app/service/main/thumbup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

func newLikeMsg(msg *databus.Message) (res interface{}, err error) {
	likeMsg := new(xmdl.LikeMsg)
	if err = json.Unmarshal(msg.Value, &likeMsg); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
		return
	}
	log.Info("get like event msg: %+v", likeMsg)
	res = likeMsg
	return
}

func likeSplit(msg *databus.Message, data interface{}) int {
	lm, ok := data.(*xmdl.LikeMsg)
	if !ok {
		log.Error("get like event msg: not ok %s", msg.Value)
		return 0
	}
	return int(lm.Mid)
}

func (s *Service) likeDo(ms []interface{}) {
	for _, m := range ms {
		lm, ok := m.(*xmdl.LikeMsg)
		if !ok {
			log.Error("get like event msg: not ok %+v", m)
			continue
		}
		var (
			oldState   int8
			err        error
			businessID int64
			ctx        = context.Background()
		)
		if businessID, err = s.checkBusinessOrigin(lm.Business, lm.OriginID); err != nil {
			log.Warn("like event: business(%v, %v) err: +v", lm.Business, lm.OriginID)
			continue
		}
		if oldState, err = s.dao.LikeState(ctx, lm.Mid, businessID, lm.OriginID, lm.MessageID); err != nil {
			log.Warn("like event: likeState(%+v) err: +v", lm)
			time.Sleep(time.Millisecond * 50)
			continue
		}
		var newState, likeType int8
		if newState, likeType, err = s.checkState(ctx, oldState, lm); err != nil {
			log.Warn("repeat like mid(%d) likeType(%d) oldState(%d) newState(%d) bid(%d) oid(%d) messageID(%d)",
				lm.Mid, likeType, oldState, newState, businessID, lm.OriginID, lm.MessageID)
			continue
		}
		var stat model.Stats
		if stat, err = s.dao.UpdateLikeState(ctx, lm.Mid, businessID, lm.OriginID, lm.MessageID, newState, lm.LikeTime); err != nil {
			log.Warn("like event: UpdateLikeState(%+v) err: +v", lm)
			time.Sleep(time.Millisecond * 50)
			continue
		}
		// 聚合数据
		key := fmt.Sprintf("%d-%d-%d", businessID, lm.OriginID, lm.MessageID)
		s.merge.Add(ctx, key, lm)
		stat = calculateCount(stat, likeType)
		s.updateCache(ctx, lm.Mid, businessID, lm.OriginID, lm.MessageID, likeType, lm.LikeTime, &stat)
		s.dao.PubStatDatabus(ctx, lm.Business, lm.Mid, &stat, lm.UpMid)
		log.Info("like event: like success params(%+v)", m)

		// 拜年祭
		target := s.mergeTarget(lm.Business, lm.MessageID)
		if target <= 0 {
			continue
		}
		if stat, err = s.dao.Stat(ctx, businessID, 0, target); err != nil {
			continue
		}
		lm.MessageID = target
		key = fmt.Sprintf("%d-%d-%d", businessID, 0, target)
		s.merge.Add(ctx, key, lm)
		stat = calculateCount(stat, likeType)
		s.updateStatCache(ctx, businessID, 0, &stat)
		s.dao.PubStatDatabus(ctx, lm.Business, lm.Mid, &stat, 0)
		log.Info("like success params(%+v)", m)
	}
}

func (s *Service) countsSplit(key string) int {
	messageIDStr := strings.Split(key, "-")[2]
	messageID, _ := strconv.Atoi(messageIDStr)
	return messageID % s.c.Merge.Worker
}

func (s *Service) updateCountsDo(c context.Context, ch int, values map[string][]interface{}) {
	mItem := make(map[model.LikeItem]*model.LikeCounts)
	for _, vs := range values {
		for _, v := range vs {
			item := v.(*xmdl.LikeMsg)
			stat := calculateCount(model.Stats{}, item.Type)
			likesCount := stat.Likes
			dislikesCount := stat.Dislikes
			likeItem := model.LikeItem{
				Business:  item.Business,
				OriginID:  item.OriginID,
				MessageID: item.MessageID,
			}
			if mItem[likeItem] == nil {
				mItem[likeItem] = &model.LikeCounts{Like: likesCount, Dislike: dislikesCount, UpMid: item.UpMid}
			} else {
				mItem[likeItem].Like += likesCount
				mItem[likeItem].Dislike += dislikesCount
				if item.UpMid > 0 {
					mItem[likeItem].UpMid = item.UpMid
				}
			}
		}
	}
	for item, count := range mItem {
		for i := 0; i < _retryTimes; i++ {
			if err := s.dao.UpdateCounts(context.Background(), s.businessMap[item.Business].ID, item.OriginID, item.MessageID, count.Like, count.Dislike, count.UpMid); err == nil {
				break
			}
		}
	}
}

// checkBusiness .
func (s *Service) checkBusiness(business string) (id int64, err error) {
	b := s.businessMap[business]
	if b == nil {
		err = ecode.ThumbupBusinessBlankErr
		return
	}
	id = b.ID
	return
}

// checkBusinessOrigin .
func (s *Service) checkBusinessOrigin(business string, originID int64) (id int64, err error) {
	b := s.businessMap[business]
	if b == nil {
		err = ecode.ThumbupBusinessBlankErr
		return
	}
	if (b.EnableOriginID == 1 && originID == 0) || (b.EnableOriginID == 0 && originID != 0) {
		err = ecode.ThumbupOriginErr
		return
	}
	id = b.ID
	return
}

// updateCache .
func (s *Service) updateCache(c context.Context, mid, businessID, originID, messageID int64, likeType int8, likeTime time.Time, stat *model.Stats) {
	if stat != nil {
		s.updateStatCache(c, businessID, originID, stat)
	}
	business := s.businessIDMap[businessID]
	likeRecord := &model.ItemLikeRecord{MessageID: messageID, Time: xtime.Time(likeTime.Unix())}
	userRecord := &model.UserLikeRecord{Mid: mid, Time: xtime.Time(likeTime.Unix())}
	switch likeType {
	case model.TypeLike:
		if business.EnableUserLikeList() {
			s.addUserlikeRecord(c, mid, businessID, model.StateLike, likeRecord)
		}
		if business.EnableItemLikeList() {
			s.addItemlikeRecord(c, businessID, messageID, model.StateLike, userRecord)
		}
	case model.TypeCancelLike:
		if business.EnableUserLikeList() {
			s.dao.DelUserLikeCache(c, mid, businessID, messageID, model.StateLike)
		}
		if business.EnableItemLikeList() {
			s.dao.DelItemLikeCache(c, messageID, businessID, mid, model.StateLike)
		}
	case model.TypeDislike:
		if business.EnableUserDislikeList() {
			s.addUserlikeRecord(c, mid, businessID, model.StateDislike, likeRecord)
		}
		if business.EnableItemDislikeList() {
			s.addItemlikeRecord(c, businessID, messageID, model.StateDislike, userRecord)
		}
	case model.TypeCancelDislike:
		if business.EnableUserDislikeList() {
			s.dao.DelUserLikeCache(c, mid, businessID, messageID, model.StateDislike)
		}
		if business.EnableItemDislikeList() {
			s.dao.DelItemLikeCache(c, messageID, businessID, mid, model.StateDislike)
		}
	case model.TypeLikeReverse:
		if business.EnableUserLikeList() {
			s.dao.DelUserLikeCache(c, mid, businessID, messageID, model.StateLike)
		}
		if business.EnableItemLikeList() {
			s.dao.DelItemLikeCache(c, messageID, businessID, mid, model.StateLike)
		}
		if business.EnableUserDislikeList() {
			s.addUserlikeRecord(c, mid, businessID, model.StateDislike, likeRecord)
		}
		if business.EnableItemDislikeList() {
			s.addItemlikeRecord(c, businessID, messageID, model.StateDislike, userRecord)
		}
	case model.TypeDislikeReverse:
		if business.EnableUserDislikeList() {
			s.dao.DelUserLikeCache(c, mid, businessID, messageID, model.StateDislike)
		}
		if business.EnableItemDislikeList() {
			s.dao.DelItemLikeCache(c, messageID, businessID, mid, model.StateDislike)
		}
		if business.EnableUserLikeList() {
			s.addUserlikeRecord(c, mid, businessID, model.StateLike, likeRecord)
		}
		if business.EnableItemLikeList() {
			s.addItemlikeRecord(c, businessID, messageID, model.StateLike, userRecord)
		}
	}
}

// updateStateCache .
func (s *Service) updateStatCache(c context.Context, businessID, originID int64, stat *model.Stats) (err error) {
	if stat == nil {
		return
	}
	var ok bool
	if originID == 0 {
		err = s.dao.AddStatsCache(c, businessID, stat)
	} else {
		if ok, err = s.dao.ExpireHashStatsCache(c, businessID, originID); ok {
			err = s.dao.AddHashStatsCache(c, businessID, originID, stat)
		}
	}
	return
}

func calculateCount(stat model.Stats, typ int8) model.Stats {
	var likesCount, dislikeCount int64
	switch typ {
	case model.TypeLike:
		likesCount = 1
	case model.TypeCancelLike:
		likesCount = -1
	case model.TypeDislike:
		dislikeCount = 1
	case model.TypeCancelDislike:
		dislikeCount = -1
	case model.TypeLikeReverse:
		likesCount = -1
		dislikeCount = 1
	case model.TypeDislikeReverse:
		likesCount = 1
		dislikeCount = -1
	}
	stat.Likes += likesCount
	stat.Dislikes += dislikeCount
	return stat
}

// checkState .
func (s *Service) checkState(c context.Context, oldState int8, lm *xmdl.LikeMsg) (newState, likeType int8, err error) {
	likeType = lm.Type
	if oldState == model.StateBlank {
		switch lm.Type {
		case model.TypeLike:
			newState = model.StateLike
		case model.TypeCancelLike:
			err = ecode.ThumbupCancelLikeErr
		case model.TypeDislike:
			newState = model.StateDislike
		case model.TypeCancelDislike:
			err = ecode.ThumbupCancelDislikeErr
		}
	} else if oldState == model.StateLike {
		switch lm.Type {
		case model.TypeLike:
			err = ecode.ThumbupDupLikeErr
			limit := s.businessMap[lm.Business].UserLikesLimit
			bid := s.businessMap[lm.Business].ID
			likeRecord := &model.ItemLikeRecord{MessageID: lm.MessageID, Time: xtime.Time(lm.LikeTime.Unix())}
			if exists, err1 := s.dao.ExpireUserLikesCache(c, lm.Mid, bid, model.StateLike); err1 == nil && exists {
				s.dao.AppendCacheUserLikeList(c, lm.Mid, likeRecord, bid, model.StateLike, limit)
			}
		case model.TypeCancelLike:
			newState = model.StateBlank
		case model.TypeDislike:
			likeType = model.TypeLikeReverse
			newState = model.StateDislike
		case model.TypeCancelDislike:
			err = ecode.ThumbupCancelDislikeErr
		}
	} else if oldState == model.StateDislike {
		switch lm.Type {
		case model.TypeLike:
			likeType = model.TypeDislikeReverse
			newState = model.StateLike
		case model.TypeCancelLike:
			err = ecode.ThumbupCancelLikeErr
		case model.TypeDislike:
			err = ecode.ThumbupDupDislikeErr
			limit := s.businessMap[lm.Business].UserLikesLimit
			bid := s.businessMap[lm.Business].ID
			likeRecord := &model.ItemLikeRecord{MessageID: lm.MessageID, Time: xtime.Time(lm.LikeTime.Unix())}
			if exists, err1 := s.dao.ExpireUserLikesCache(c, lm.Mid, bid, model.StateDislike); err1 == nil && exists {
				s.dao.AppendCacheUserLikeList(c, lm.Mid, likeRecord, bid, model.StateDislike, limit)
			}
		case model.TypeCancelDislike:
			newState = model.StateBlank
		}
	} else {
		log.Warn("oldState abnormal mid:%d business:%v oid:%d messageID:%d oldState:%d", lm.Mid, lm.Business, lm.OriginID, lm.MessageID, oldState)
	}
	return
}

func (s *Service) mergeTarget(business string, aid int64) int64 {
	if s.statMerge != nil && s.statMerge.Business == business && s.statMerge.Sources[aid] {
		return s.statMerge.Target
	}
	return 0
}
