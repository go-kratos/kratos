package like

import (
	"context"
	"strconv"

	"go-common/app/interface/main/activity/model/bnj"
	"go-common/app/interface/main/activity/model/like"
	suitmdl "go-common/app/service/main/usersuit/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Reward get bnj preview reward.
func (s *Service) Reward(c context.Context, mid int64, step int) (err error) {
	reward, ok := s.reward[step]
	if !ok {
		err = ecode.RequestErr
		return
	}
	var (
		likeActs  map[int64]int
		likeScore map[int64]int64
		check     bool
	)
	actID := s.c.Bnj2019.ActID
	subID := s.c.Bnj2019.SubID
	if likeScore, err = s.dao.LikeActLidCounts(c, []int64{subID}); err != nil {
		log.Error("Reward s.dao.LikeActLidCounts(subID:%d) error(%+v)", subID, err)
		return
	}
	if score, ok := likeScore[subID]; !ok || score < reward.Condition {
		err = ecode.ActivityBnjSubLow
		return
	}
	if likeActs, err = s.dao.LikeActs(c, actID, mid, []int64{subID}); err != nil {
		log.Error("Reward s.dao.LikeActs(subID:%d,actID:%d,mid:%d) error(%+v)", subID, actID, mid, err)
		return
	}
	if isSub, ok := likeActs[subID]; !ok || isSub == 0 {
		if _, err = s.LikeAct(c, &like.ParamAddLikeAct{Sid: actID, Lid: subID, Score: 1}, mid); err != nil {
			return
		}
	}
	// check has reward
	if check, err = s.bnjDao.CacheHasReward(c, mid, subID, step); err != nil || !check {
		log.Error("Reward s.dao.CacheHasReward(mid:%d,subID:%d,step:%d) error(%v) check(%v)", mid, subID, step, err, check)
		err = ecode.ActivityBnjHasReward
		return
	}
	switch reward.RewardType {
	case bnj.RewardTypePendant:
		rewardID, e := strconv.ParseInt(reward.RewardID, 10, 64)
		if e != nil {
			err = ecode.ActivityRewardConfErr
			return
		}
		err = s.suit.GrantByMids(c, &suitmdl.ArgGrantByMids{Mids: []int64{mid}, Pid: rewardID, Expire: reward.Expire})
	case bnj.RewardTypeCoupon:
		err = s.bnjDao.GrantCoupon(c, mid, reward.RewardID)
		// TODO check err code
	}
	if err != nil {
		log.Error("Reward (%+v) error(%+v)", reward, err)
		err = ecode.ActivityBnjRewardFail
		if e := s.bnjDao.DelCacheHasReward(c, mid, subID, step); e != nil {
			log.Error("s.dao.DelCacheHasReward(%+v) error(%v)", reward, e)
		}
	}
	return
}
