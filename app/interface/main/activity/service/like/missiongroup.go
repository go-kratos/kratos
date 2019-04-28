package like

import (
	"context"
	"net"
	"time"

	ldao "go-common/app/interface/main/activity/dao/like"
	l "go-common/app/interface/main/activity/model/like"
	accapi "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
)

// MissionLike launch buff activity .
func (s *Service) MissionLike(c context.Context, sid, mid int64) (lid int64, err error) {
	var (
		subject      *l.SubjectItem
		now          = time.Now().Unix()
		missionGroup int64
		group        *l.MissionGroup
	)
	if subject, err = s.dao.ActSubject(c, sid); err != nil {
		log.Error("s.dao.ActSubject(%d) error(%+v)", sid, err)
		return
	}
	if subject.ID == 0 || subject.Type != l.MISSIONGROUP {
		err = ecode.ActivityNotExist
		return
	}
	if subject.Stime.Time().Unix() > now {
		err = ecode.ActivityNotStart
		return
	}
	if subject.Etime.Time().Unix() < now {
		err = ecode.ActivityOverEnd
		return
	}
	if missionGroup, err = s.dao.LikeMissionBuff(c, sid, mid); err != nil {
		log.Error("s.dao.LikeMissionBuff(%d,%d) error(%+v)", sid, mid, err)
		return
	}
	if missionGroup > 0 {
		err = ecode.ActivityHasMissionGroup
		return
	}
	group = &l.MissionGroup{
		Sid:   sid,
		Mid:   mid,
		State: ldao.MissionStateInit,
	}
	if lid, err = s.dao.MissionGroupAdd(c, group); err != nil {
		log.Error("s.dao.MissionGroupAdd(%d,%d) error(%+v)", sid, mid, err)
		return
	}
	s.dao.AddCacheLikeMissionBuff(c, sid, lid, mid)
	return
}

// MissionInfo .
func (s *Service) MissionInfo(c context.Context, sid, lid, mid int64) (res *l.MissionInfo, err error) {
	var hasBufErr, hasHelpErr error
	res = &l.MissionInfo{}
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() error {
		res.HasBuff, hasBufErr = s.dao.LikeMissionBuff(errCtx, sid, mid)
		if hasBufErr != nil {
			log.Error("s.dao.LikeMissionBuff(%d,%d) error(%+v)", sid, mid, hasBufErr)
		}
		return nil
	})
	eg.Go(func() error {
		res.HasHelp, hasHelpErr = s.dao.ActMission(errCtx, sid, lid, mid)
		if hasHelpErr != nil {
			log.Error("s.dao.ActMission(%d,%d,%d) error(%+v)", sid, lid, mid, hasHelpErr)
		}
		return nil
	})
	eg.Wait()
	return
}

// MissionUser .
func (s *Service) MissionUser(c context.Context, sid, lid int64) (res *l.MissionFriends, err error) {
	var (
		groups map[int64]*l.MissionGroup
		group  *l.MissionGroup
		member *accapi.InfoReply
	)
	if groups, err = s.dao.MissionGroupItems(c, []int64{lid}); err != nil {
		log.Error("s.dao.MissionGroupItems(%v) error(%v)", lid, err)
		return
	}
	if _, ok := groups[lid]; !ok {
		err = ecode.ActivityNotExist
		return
	}
	group = groups[lid]
	if group.ID == 0 || group.Sid != sid {
		err = ecode.ActivityNotExist
		return
	}
	if member, err = s.accClient.Info3(c, &accapi.MidReq{Mid: group.Mid}); err != nil {
		log.Error(" s.acc.Info3(c,&accmdl.ArgMids{Mid:%d}) error(%v)", group.Mid, err)
		return
	}
	res = &l.MissionFriends{
		Name: member.Info.Name,
		Face: member.Info.Face,
		Mid:  member.Info.Mid,
	}
	return
}

// MissionLikeAct help to mission group .
func (s *Service) MissionLikeAct(c context.Context, p *l.ParamMissionLikeAct, mid int64) (data *l.MissionLikeAct, err error) {
	var (
		subject                                         *l.SubjectItem
		groups                                          map[int64]*l.MissionGroup
		group                                           *l.MissionGroup
		memberRly                                       *accapi.ProfileReply
		now                                             = time.Now().Unix()
		ActMissionID, likeLimit, missionActCount, mLid  int64
		score                                           = int64(1)
		missionActList                                  *l.ActMissionGroup
		lottery                                         *l.Lottery
		subErr, groupErr, missionErr, caculErr, incrErr error
	)
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() error {
		subject, subErr = s.dao.ActSubject(errCtx, p.Sid)
		return subErr
	})
	eg.Go(func() error {
		groups, groupErr = s.dao.MissionGroupItems(errCtx, []int64{p.Lid})
		return groupErr
	})
	if err = eg.Wait(); err != nil {
		log.Error("MissionLikeAct:eg.Wait error(%v)", err)
		return
	}
	if _, ok := groups[p.Lid]; ok {
		group = groups[p.Lid]
	} else {
		err = ecode.ActivityNotExist
		return
	}
	if subject.ID == 0 || subject.Type != l.MISSIONGROUP || group.ID == 0 || group.Sid != p.Sid {
		err = ecode.ActivityNotExist
		return
	}
	if group.Mid == mid {
		err = ecode.ActivityMGNotYourself
		return
	}
	if memberRly, err = s.accClient.Profile3(c, &accapi.MidReq{Mid: mid}); err != nil {
		log.Error(" s.acc.Profile3(c,&accmdl.ArgMid{Mid:%d}) error(%v)", mid, err)
		return
	}
	if err = s.judgeUser(c, subject, memberRly.Profile); err != nil {
		return
	}
	if subject.Lstime.Time().Unix() >= now {
		err = ecode.ActivityMissionNotStart
		return
	}
	if subject.Letime.Time().Unix() <= now {
		err = ecode.ActivityMissionHasEnd
		return
	}
	if ActMissionID, err = s.dao.ActMission(c, p.Sid, p.Lid, mid); err != nil {
		log.Error("s.dao.ActMission(%v) error(%+v)", p, err)
		return
	}
	if ActMissionID > 0 {
		err = ecode.ActivityHasMission
		return
	}
	if subject.LikeLimit > 0 {
		if likeLimit, err = s.dao.MissionLikeLimit(c, p.Sid, mid); err != nil {
			log.Error("s.dao.ActMission(%v) error(%+v)", p, err)
			return
		}
		if likeLimit >= subject.LikeLimit {
			err = ecode.ActivityOverMissionLimit
			return
		}
	}
	if missionActCount, err = s.dao.SetMissionTop(c, p.Sid, p.Lid, score, now); err != nil {
		log.Error("s.dao.SetMissionTop(%v) error(%+v)", p, err)
		return
	}
	missionActList = &l.ActMissionGroup{
		Lid:    p.Lid,
		Sid:    p.Sid,
		Mid:    mid,
		Action: score,
		IPv6:   make([]byte, 0),
	}
	if IPv6 := net.ParseIP(metadata.String(c, metadata.RemoteIP)); IPv6 != nil {
		missionActList.IPv6 = IPv6
	}
	if mLid, err = s.dao.AddActMission(c, missionActList); err != nil {
		log.Error("s.dao.AddActMission(%v) error(%+v)", p, err)
		return
	}
	egT, errCtxT := errgroup.WithContext(c)
	egT.Go(func() error {
		missionErr = s.dao.AddCacheActMission(errCtxT, p.Sid, mLid, p.Lid, mid)
		return missionErr
	})
	egT.Go(func() error {
		caculErr = s.CalculateAchievement(errCtxT, p.Sid, group.Mid, missionActCount)
		return caculErr
	})
	egT.Go(func() error {
		_, incrErr = s.dao.InrcMissionLikeLimit(errCtxT, p.Sid, mid, score)
		return incrErr
	})
	if err = egT.Wait(); err != nil {
		log.Error("MissionLikeAct:eg.Wait add cache error(%v)", err)
		return
	}
	if lottery, err = s.dao.LotteryIndex(c, s.c.Rule.LotteryActID, int64(0), int64(0), mid); err != nil {
		log.Error("s.dao.LotteryIndex(%d) mid(%d) error(%+v)", s.c.Rule.LotteryActID, mid, err)
		return
	}
	data = &l.MissionLikeAct{
		Mlid:    mLid,
		Lottery: lottery,
	}
	return
}

// CalculateAchievement .
func (s *Service) CalculateAchievement(c context.Context, sid, mid int64, missionCount int64) (err error) {
	var (
		achieves *l.Achievements
		avState  int64
	)
	if achieves, err = s.dao.ActLikeAchieves(c, sid); err != nil {
		log.Error("s.dao.ActLikeAchieves(%d,%d) error(%v)", sid, mid, err)
		return
	}
	if len(achieves.Achievements) > 0 {
		for _, v := range achieves.Achievements {
			if v.Unlock == missionCount {
				if v.Award == ldao.HaveAward {
					avState = ldao.AwardNotChange
				} else {
					avState = ldao.AwardNoGet
				}
				if _, err = s.dao.AddUserAchievment(c, &l.ActLikeUserAchievement{Aid: v.ID, Sid: sid, Mid: mid, Award: avState}); err != nil {
					log.Error("s.dao.AddUserAchievment(%d,%d,%v) error(%+v)", sid, mid, v, err)
					return
				}
				break
			}
		}
	}
	return
}

// MissionRank get user rank .
func (s *Service) MissionRank(c context.Context, sid, mid int64) (data *l.MissionRank, err error) {
	data = &l.MissionRank{Rank: -1}
	if data.Lid, err = s.dao.LikeMissionBuff(c, sid, mid); err != nil {
		log.Error("s.dao.LikeMissionBuff(%d,%d) error(%+v)", sid, mid, err)
		return
	}
	if data.Lid > 0 {
		if data.Score, err = s.dao.MissionLidScore(c, sid, data.Lid); err != nil {
			log.Error("s.dao.MissionLidScore(%d,%d) error(%+v)", sid, data.Lid, err)
			return
		}
		if data.Rank, err = s.dao.MissionLidRank(c, sid, data.Lid); err != nil {
			log.Error("s.dao.MissionLidRank(%d,%d) error(%+v)", sid, data.Lid, err)
			return
		}
		if data.Rank >= 0 {
			data.Rank = data.Rank + 1
		}
	}
	return
}

// MissionTops get the top list .
func (s *Service) MissionTops(c context.Context, sid int64, num int) (data []*l.MissionFriends, err error) {
	var (
		lids       []int64
		lidsList   map[int64]*l.MissionGroup
		mids       []int64
		membersRly *accapi.InfosReply
	)
	if lids, err = s.dao.MissionScoreList(c, sid, 0, num-1); err != nil {
		log.Error("s.dao.MissionScoreList(%d) error(%+v)", sid, err)
		return
	}
	if len(lids) > 0 {
		if lidsList, err = s.dao.MissionGroupItems(c, lids); err != nil {
			log.Error("s.dao.MissionGroupItems(%v) error(%v)", lids, err)
			return
		}
		mids = make([]int64, 0, len(lidsList))
		for _, v := range lidsList {
			if v.ID > 0 {
				mids = append(mids, v.Mid)
			}
		}
		if len(mids) > 0 {
			if membersRly, err = s.accClient.Infos3(c, &accapi.MidsReq{Mids: mids}); err != nil {
				log.Error("s.acc.Infos3(%v) error(%v)", mids, err)
				return
			}
		}
		data = make([]*l.MissionFriends, 0, len(lids))
		for _, v := range lids {
			if _, ok := lidsList[v]; ok {
				n := &l.MissionFriends{Mid: lidsList[v].Mid}
				if membersRly != nil {
					if val, y := membersRly.Infos[lidsList[v].Mid]; y {
						n.Name = val.Name
						n.Face = val.Face
					}
				}
				data = append(data, n)
			}
		}
	}
	return
}

// MissionFriendsList .
func (s *Service) MissionFriendsList(c context.Context, p *l.ParamMissionFriends, mid int64) (data []*l.MissionFriends, err error) {
	var (
		groups      map[int64]*l.MissionGroup
		ActList     []*l.ActMissionGroup
		ActMissions *l.ActMissionGroups
		mids        []int64
		membersRly  *accapi.InfosReply
		score       int64
		actLen      int
	)
	if groups, err = s.dao.MissionGroupItems(c, []int64{p.Lid}); err != nil {
		log.Error("s.dao.MissionGroupItems(%d) error(%v)", p.Lid, err)
		return
	}
	if _, ok := groups[p.Lid]; !ok {
		err = ecode.ActivityNotExist
		return
	}
	if groups[p.Lid].ID == 0 || groups[p.Lid].Mid != mid || groups[p.Lid].Sid != p.Sid {
		err = ecode.ActivityNotExist
		return
	}
	if ActMissions, err = s.dao.ActMissionFriends(c, p.Sid, p.Lid); err != nil {
		log.Error("s.dao.ActMissionFriends(%v) error(%+v)", p, err)
		return
	}
	ActList = ActMissions.ActMissionGroups
	actLen = len(ActList)
	score, _ = s.dao.MissionLidScore(c, p.Sid, p.Lid)
	if int64(actLen) < score && actLen < p.Size {
		// need to update cache
		s.dao.DelCacheActMissionFriends(c, p.Sid, p.Lid)
		if ActMissions, err = s.dao.ActMissionFriends(c, p.Sid, p.Lid); err != nil {
			log.Error("s.dao.ActMissionFriends(%v) error(%+v)", p, err)
			return
		}
		ActList = ActMissions.ActMissionGroups
		actLen = len(ActList)
	}
	if actLen > p.Size {
		ActList = ActList[:p.Size]
		actLen = p.Size
	}
	mids = make([]int64, 0, actLen)
	for _, v := range ActList {
		mids = append(mids, v.Mid)
	}
	if len(mids) > 0 {
		if membersRly, err = s.accClient.Infos3(c, &accapi.MidsReq{Mids: mids}); err != nil {
			log.Error("s.acc.Infos3(%v) error(%v)", mids, err)
			return
		}
	}
	data = make([]*l.MissionFriends, 0, len(ActList))
	for _, v := range ActList {
		n := &l.MissionFriends{Mid: v.Mid}
		if membersRly != nil {
			if val, y := membersRly.Infos[v.Mid]; y {
				n.Name = val.Name
				n.Face = val.Face
			}
		}
		data = append(data, n)
	}
	return
}

// MissionAward .
func (s *Service) MissionAward(c context.Context, sid, mid int64) (data []*l.MissionAward, err error) {
	var (
		achieves     *l.Achievements
		userAchieves []*l.ActLikeUserAchievement
		userAchMap   map[int64]*l.ActLikeUserAchievement
		achErr       error
		userErr      error
	)
	eg, errCtx := errgroup.WithContext(c)
	eg.Go(func() error {
		achieves, achErr = s.dao.ActLikeAchieves(errCtx, sid)
		return achErr
	})
	eg.Go(func() error {
		userAchieves, userErr = s.dao.UserAchievement(c, sid, mid)
		return userErr
	})
	if err = eg.Wait(); err != nil {
		log.Error("MissionAward:eg.Wait() error(%+v)", err)
		return
	}
	userAchMap = make(map[int64]*l.ActLikeUserAchievement, len(userAchieves))
	for _, v := range userAchieves {
		userAchMap[v.Aid] = v
	}
	for _, val := range achieves.Achievements {
		n := &l.MissionAward{Name: val.Name, Image: val.Image}
		if v, ok := userAchMap[val.ID]; ok {
			n.ID = v.ID
			n.Award = v.Award
		}
		data = append(data, n)
	}
	return
}

// MissionAchieve .
func (s *Service) MissionAchieve(c context.Context, sid, id, mid int64) (res int64, err error) {
	var (
		useActAchieve *l.ActLikeUserAchievement
		award         int64
	)
	if useActAchieve, err = s.dao.ActUserAchieve(c, id); err != nil {
		log.Error("s.dao.ActUserAchieve(%d) error(%+v)", id, err)
		return
	}
	if useActAchieve.ID == 0 || useActAchieve.Mid != mid || useActAchieve.Sid != sid {
		err = ecode.ActivityNotAward
		return
	}
	if award, err = s.dao.CacheActUserAward(c, id); err != nil {
		log.Info("s.dao.CacheActUserAward(%d) error(%v)", id, err)
	}
	if award > 0 {
		err = ecode.ActivityHasAward
		return
	}
	if res, err = s.dao.ActUserAchieveChange(c, id, ldao.AwardHasChange); err != nil {
		log.Error("s.dao.ActUserAchieveChange(%d) error(%+v)", id, err)
		return
	}
	s.dao.AddCacheActUserAward(c, id, id)
	return
}
