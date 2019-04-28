package bws

import (
	"context"
	"time"

	bwsmdl "go-common/app/interface/main/activity/model/bws"
	suitmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// Award  achievement award
func (s *Service) Award(c context.Context, loginMid int64, p *bwsmdl.ParamAward) (err error) {
	var (
		userAchieves []*bwsmdl.UserAchieve
		userAward    int64 = -1
	)
	if _, ok := s.awardMids[loginMid]; !ok {
		err = ecode.ActivityNotAwardAdmin
		return
	}
	if p.Key == "" {
		if p.Key, err = s.midToKey(c, p.Mid); err != nil {
			return
		}
	}
	if userAchieves, err = s.dao.UserAchieves(c, p.Bid, p.Key); err != nil {
		err = ecode.ActivityAchieveFail
		return
	}
	if len(userAchieves) == 0 {
		err = ecode.ActivityNoAchieve
		return
	}
	for _, v := range userAchieves {
		if v.Aid == p.Aid {
			userAward = v.Award
			break
		}
	}
	if userAward == -1 {
		err = ecode.ActivityNoAchieve
		return
	} else if userAward == _noAward {
		err = ecode.ActivityNoAward
		return
	} else if userAward == _awardAlready {
		err = ecode.ActivityAwardAlready
		return
	}
	if err = s.dao.Award(c, p.Key, p.Aid); err != nil {
		log.Error("s.dao.Award key(%s)  error(%v)", p.Key, err)
	}
	s.dao.DelCacheUserAchieves(c, p.Bid, p.Key)
	return
}

// Achievements achievements list
func (s *Service) Achievements(c context.Context, p *bwsmdl.ParamID) (rs *bwsmdl.Achievements, err error) {
	var mapCnt map[int64]int64
	if rs, err = s.dao.Achievements(c, p.Bid); err != nil || rs == nil || len(rs.Achievements) == 0 {
		log.Error("s.dao.Achievements error(%v)", err)
		err = ecode.ActivityAchieveFail
		return
	}
	if mapCnt, err = s.countAchieves(c, p.Bid, p.Day); err != nil || len(mapCnt) == 0 {
		err = nil
		return
	}
	for _, achieve := range rs.Achievements {
		achieve.UserCount = mapCnt[achieve.ID]
	}
	return
}

func (s *Service) countAchieves(c context.Context, bid int64, day string) (rs map[int64]int64, err error) {
	var countAchieves []*bwsmdl.CountAchieves
	if day == "" {
		day = today()
	}
	if countAchieves, err = s.dao.AchieveCounts(c, bid, day); err != nil {
		log.Error("s.dao.RawCountAchieves error(%v)", err)
		return
	}
	rs = make(map[int64]int64, len(countAchieves))
	for _, countAchieve := range countAchieves {
		rs[countAchieve.Aid] = countAchieve.Count
	}
	return
}

// Achievement Achievement
func (s *Service) Achievement(c context.Context, p *bwsmdl.ParamID) (rs *bwsmdl.Achievement, err error) {
	var (
		achieves *bwsmdl.Achievements
	)
	if achieves, err = s.dao.Achievements(c, p.Bid); err != nil || achieves == nil || len(achieves.Achievements) == 0 {
		log.Error("s.dao.Achievements error(%v)", err)
		err = ecode.ActivityAchieveFail
		return
	}
	for _, Achievement := range achieves.Achievements {
		if Achievement.ID == p.ID {
			rs = Achievement
			break
		}
	}
	if rs == nil {
		err = ecode.ActivityIDNotExists
	}
	return
}

func (s *Service) userAchieves(c context.Context, bid int64, key string) (res []*bwsmdl.UserAchieveDetail, err error) {
	var (
		usAchieves []*bwsmdl.UserAchieve
		achieves   *bwsmdl.Achievements
	)
	if usAchieves, err = s.dao.UserAchieves(c, bid, key); err != nil {
		err = ecode.ActivityUserAchieveFail
		return
	}
	if len(usAchieves) == 0 {
		return
	}
	if achieves, err = s.dao.Achievements(c, bid); err != nil || achieves == nil || len(achieves.Achievements) == 0 {
		err = ecode.ActivityAchieveFail
		return
	}
	achievesMap := make(map[int64]*bwsmdl.Achievement, len(achieves.Achievements))
	for _, v := range achieves.Achievements {
		achievesMap[v.ID] = v
	}
	for _, v := range usAchieves {
		detail := &bwsmdl.UserAchieveDetail{UserAchieve: v}
		if achieve, ok := achievesMap[v.Aid]; ok {
			detail.Name = achieve.Name
			detail.Icon = achieve.Icon
			detail.Dic = achieve.Dic
			detail.LockType = achieve.LockType
			detail.Unlock = achieve.Unlock
			detail.Bid = achieve.Bid
			detail.IconBig = achieve.IconBig
			detail.IconActive = achieve.IconActive
			detail.IconActiveBig = achieve.IconActiveBig
			detail.SuitID = achieve.SuitID
		}
		res = append(res, detail)
	}
	return
}

func (s *Service) addAchieve(c context.Context, mid int64, achieve *bwsmdl.Achievement, key string) (err error) {
	var uaID int64
	if uaID, err = s.dao.AddUserAchieve(c, achieve.Bid, achieve.ID, achieve.Award, key); err != nil {
		err = ecode.ActivityAddAchieveFail
		return
	}
	if err = s.dao.AppendUserAchievesCache(c, achieve.Bid, key, &bwsmdl.UserAchieve{ID: uaID, Aid: achieve.ID, Award: achieve.Award, Ctime: xtime.Time(time.Now().Unix())}); err != nil {
		return
	}
	s.cache.Do(c, func(c context.Context) {
		s.dao.IncrCacheAchieveCounts(c, achieve.Bid, achieve.ID, today())
		var (
			keyID int64
			e     error
		)
		if mid == 0 {
			if mid, keyID, e = s.keyToMid(c, key); e != nil || mid == 0 {
				log.Warn("Lottery keyID(%d) key(%s) error(%v)", keyID, key, e)
			}
		}
		if mid > 0 {
			if achieve.SuitID > 0 {
				arg := &suitmdl.ArgGrantByMids{Mids: []int64{mid}, Pid: achieve.SuitID, Expire: s.c.Rule.BwsSuitExpire}
				if e := s.suitRPC.GrantByMids(c, arg); e != nil {
					log.Error("addAchieve s.suit.GrantByMids(%d,%d) error(%v)", mid, achieve.SuitID, e)
				}
				log.Warn("Suit mid(%d) suitID(%d)", mid, achieve.SuitID)
			}
			if _, ok := s.lotteryAids[achieve.ID]; ok {
				s.dao.AddLotteryMidCache(c, achieve.ID, mid)
			}
		}
	})
	return
}
