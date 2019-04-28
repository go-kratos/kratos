package service

import (
	"context"
	"fmt"
	"sort"

	"go-common/app/service/main/account/model"
	"go-common/app/service/main/assist/model/assist"
	"go-common/app/service/main/assist/model/notify"
	"go-common/library/log"
)

// AddAssist add assist.
func (s *Service) AddAssist(c context.Context, mid, assistMid int64) (err error) {
	if err = s.checkFollow(c, mid, assistMid); err != nil {
		log.Error("s.checkFollow(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if err = s.checkIdentify(c, assistMid); err != nil {
		log.Error("s.checkIdentify err: (%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if err = s.checkBanned(c, assistMid); err != nil {
		log.Error("s.checkBanned err: (%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if err = s.checkIsAssist(c, mid, assistMid); err != nil {
		log.Error("s.checkIsAssist err: (%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if err = s.checkTotalLimit(c, mid); err != nil {
		log.Error("s.limitDailyCntAddAllAss err: (%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if err = s.checkSameLimit(c, mid, assistMid); err != nil {
		log.Error("s.limitDailyCntAddSameAss err: (%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if err = s.checkMaxAssistCnt(c, mid); err != nil {
		log.Error("s.limitAssistCnt err: (%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if _, err = s.ass.AddAssist(c, mid, assistMid); err != nil {
		log.Error("s.ass.AddAssist(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	var card *model.Card
	if card, err = s.acc.Card(c, mid, ""); err != nil {
		log.Error("s.acc.Card(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	_ = s.SendSysNotify(c, notify.AddAssNotifyAct, mid, assistMid, card.Name, "")
	// set cache async
	s.asyncCache(func() {
		_ = s.ass.DelCacheAss(context.TODO(), mid)
		_ = s.ass.DelAssUpAllCache(context.TODO(), assistMid)
		_ = s.ass.IncrAssCnt(context.TODO(), mid, assistMid)
	})
	return
}

// DelAssist delete assist.
func (s *Service) DelAssist(c context.Context, mid, assistMid int64) (err error) {
	if err = s.checkIsNotAssist(c, mid, assistMid); err != nil {
		log.Error("s.checkIsNotAssist err: (%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	if _, err = s.ass.DelAssist(c, mid, assistMid); err != nil {
		log.Error("s.ass.DelAssist(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	// del cache
	s.asyncCache(func() {
		_ = s.ass.DelCacheAss(context.TODO(), mid)
		_ = s.ass.DelAssUpAllCache(context.TODO(), assistMid)
	})
	var card *model.Card
	if card, err = s.acc.Card(c, mid, ""); err != nil {
		log.Error("s.acc.Card(%d,%d) error(%v)", mid, assistMid, err)
		return
	}
	log.Info("s.SendSysNotify (%d,%d,%s)", mid, assistMid, card.Name)
	_ = s.SendSysNotify(c, notify.DelAssNotifyAct, mid, assistMid, card.Name, "")
	return
}

// Assists get assists
func (s *Service) Assists(c context.Context, mid int64) (as []*assist.Assist, err error) {
	if as, err = s.ass.Assists(c, mid); err != nil {
		log.Error("s.ass.Assists(%d) error(%v)", mid, err)
		return
	}
	totalm, err := s.ass.LogCount(c, mid)
	if err != nil {
		log.Error("s.ass.LogCount(%d) error(%v)", mid, err)
		return
	}
	if len(totalm) == 0 {
		return
	}
	for _, a := range as {
		a.Total = totalm[a.AssistMid]
	}
	return
}

// AssistsMidsTotal get multi mids assists group by mid
func (s *Service) AssistsMidsTotal(c context.Context, mid int64, assmids []int64) (totalm map[int64]map[int8]map[int8]int, err error) {
	if totalm, err = s.ass.AssistsMidsTotal(c, mid, assmids); err != nil {
		log.Error("s.ass.AssistsMidsTotal(%d) error(%v)", mid, err)
		return
	}
	return
}

// Assist get assist allow default value is 1
func (s *Service) Assist(c context.Context, mid, assistMid, tp int64) (ar *assist.AssistRes, err error) {
	assistMids := []int64{}
	ar = &assist.AssistRes{Assist: 0, Allow: 0, Count: 0}
	if assistMids, err = s.ass.GetCacheAss(c, mid); err != nil {
		return
	}
	if assistMids == nil {
		assistMids = []int64{} // for set [] cache
		// no cache, and get from db
		var assists []*assist.Assist
		if assists, err = s.ass.Assists(c, mid); err != nil {
			return
		}
		for _, a := range assists {
			assistMids = append(assistMids, a.AssistMid)
		}
		// set cache async
		s.asyncCache(func() {
			_ = s.ass.SetCacheAss(context.TODO(), mid, assistMids)
		})
	}
	// check up has assist
	if len(assistMids) == 0 {
		return
	}
	// check is assist
	isAss := false
	for _, v := range assistMids {
		if v == assistMid {
			isAss = true
			break
		}
	}
	if !isAss {
		return
	}
	ar.Assist = 1
	// get daily count
	var cnt int64
	if cnt, err = s.ass.DailyLogCount(c, mid, assistMid, tp); err != nil {
		return
	}
	if cnt < s.c.MaxTypeCnt {
		ar.Allow = 1
	}
	ar.Count = cnt
	return
}

// AssistIDs get assists list for others.
func (s *Service) AssistIDs(c context.Context, mid int64) (assistMids []int64, err error) {
	if assistMids, err = s.ass.GetCacheAss(c, mid); err != nil {
		// cache failed, return
		err = nil
		return
	}
	// no cache, and get from db
	if assistMids == nil {
		assistMids = []int64{} // for set [] cache
		var assists []*assist.Assist
		if assists, err = s.ass.Assists(c, mid); err != nil {
			return
		}
		for _, a := range assists {
			assistMids = append(assistMids, a.AssistMid)
		}
		// set cache async
		s.asyncCache(func() {
			_ = s.ass.SetCacheAss(context.TODO(), mid, assistMids)
		})
	}
	return
}

// SendSysNotify fn
func (s *Service) SendSysNotify(c context.Context, action int, mid, assistMid int64, upUname, assistUname string) (err error) {
	if action == notify.AddAssNotifyAct {
		err = s.msg.Send(c, notify.Mc, notify.AddTitle, fmt.Sprintf(notify.AddContent, upUname, mid, upUname, assistMid), assistMid)
	}
	if action == notify.DelAssNotifyAct {
		err = s.msg.Send(c, notify.Mc, notify.DelTitle, fmt.Sprintf(notify.DelContent, upUname, mid), assistMid)
	}
	if action == notify.DelAssNotifyFollowerAct {
		err = s.msg.Send(c, notify.Mc, notify.DelFollowerTitle, fmt.Sprintf(notify.DelFollowerContent, assistUname, assistMid), mid)
	}
	log.Info("action(%d), mid(%d), assistMid(%d), upUname(%s), assistUname(%s)", action, mid, assistMid, upUname, assistUname)
	return
}

// Exit delete assist from follower.
func (s *Service) Exit(c context.Context, mid, assistMid int64) (err error) {
	if _, err = s.ass.DelAssist(c, mid, assistMid); err != nil {
		log.Error("s.ass.DelAssist(%d,%d) error(%v)", mid, assistMid, err)
	}
	log.Info("s.DelAssForRLFollower (%d,%d,%s)", mid, assistMid)
	// del cache
	s.asyncCache(func() {
		_ = s.ass.DelCacheAss(context.TODO(), mid)
		_ = s.ass.DelAssUpAllCache(context.TODO(), assistMid)
	})
	var (
		card = &model.Card{}
	)
	if card, err = s.acc.Card(c, assistMid, ""); err != nil {
		return
	}
	if assistUname := card.Name; len(assistUname) != 0 {
		log.Info("s.SendSysNotify (%d,%d,%s)", mid, assistMid, assistUname)
		_ = s.SendSysNotify(c, notify.DelAssNotifyFollowerAct, mid, assistMid, "", assistUname)
	}
	return
}

// AssistUps get ups who already sign me as assist.
func (s *Service) AssistUps(c context.Context, assistMid, pn, ps int64) (assistUpsPager *assist.AssistUpsPager, err error) {
	var (
		total  int64
		upMids = []int64{}
		assUps = []*assist.AssistUp{}
		ups    = make(map[int64]*assist.Up, ps)
	)
	assistUpsPager = &assist.AssistUpsPager{
		Data: []*assist.AssistUp{},
		Pager: assist.Pager{
			Pn:    pn,
			Ps:    ps,
			Total: total,
		},
	}
	if upMids, ups, total, err = s.ass.AssUpCacheWithScore(c, assistMid, (pn-1)*ps, pn*ps); err != nil {
		// cache failed, return
		log.Error("s.ass.AssUpCacheWithScore, (upMids %v), (ps,pn: %d, %d) err:%v", upMids, ps, pn, err)
		err = nil
		return
	}
	// no cache, and get from db
	if total == 0 {
		upMids, ups, total, err = s.ass.Ups(c, assistMid, pn, ps)
		if err != nil {
			log.Error("s.ass.Ups assistMid:(%d) error(%v)", assistMid, err)
			return
		}
		s.asyncCache(func() {
			_ = s.ass.AddAssUpAllCache(context.TODO(), assistMid, ups)
		})
	}
	if len(upMids) > 0 {
		var cards map[int64]*model.Card
		if cards, err = s.acc.Cards(c, upMids); err != nil {
			log.Error("s.ass.Cards(%d) error(%v)", ups, err)
			return
		}
		for _, card := range cards {
			if _, ok := ups[card.Mid]; ok {
				assUps = append(assUps, &assist.AssistUp{
					Mid:            card.Mid,
					Name:           card.Name,
					Sign:           card.Sign,
					Avatar:         card.Face,
					OfficialVerify: card.Official,
					CTime:          ups[card.Mid].CTime,
					Vip:            card.Vip,
				})
			}
		}
		sort.Sort(assist.SortUpsByCtime(assUps))
	}
	assistUpsPager.Pager.Total = total
	assistUpsPager.Data = assUps
	return
}
