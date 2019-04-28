package service

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go-common/app/interface/main/growup/model"

	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

const (
	ceilings = 50000
	rule     = 1
	detail   = 3
	question = 5
	answer   = 6
)

// GetAwardUpStatus get award up status
func (s *Service) GetAwardUpStatus(c context.Context, awardID, mid int64) (status *model.AwardUpStatus, err error) {
	joined, err := s.isJoined(c, mid, awardID)
	if err != nil {
		return
	}

	fans, err := s.getUpFans(c, mid)
	if err != nil {
		return
	}

	status = &model.AwardUpStatus{
		Joined:    joined,
		Qualified: checkQualification(fans),
	}
	return
}

// GetWinningRecord get winning record
func (s *Service) GetWinningRecord(c context.Context, mid int64) (rs []*model.WinningRecord, err error) {
	awardIDs, err := s.dao.GetAwardJoinRecord(c, mid)
	if err != nil {
		return
	}

	if len(awardIDs) == 0 {
		return
	}

	as := make([]int64, 0)
	for awardID := range awardIDs {
		as = append(as, awardID)
	}

	awards, err := s.dao.JoinedSpecialAwards(c, as)
	if err != nil {
		return
	}

	sort.Slice(awards, func(i, j int) bool {
		return awards[i].CycleStart < awards[j].CycleStart
	})

	// am map[award_id]prize_id
	am, err := s.dao.AwardIDsByWinner(c, mid)
	if err != nil {
		return
	}

	now := time.Now().Unix()
	rs = make([]*model.WinningRecord, 0)
	for _, award := range awards {
		// in selection
		if now > int64(award.CycleEnd) && award.OpenStatus == 1 {
			rs = append(rs, &model.WinningRecord{
				AwardID:   award.AwardID,
				AwardName: award.AwardName,
				State:     2,
			})
			continue
		}
		// finished
		if award.OpenStatus == 2 {
			wr := &model.WinningRecord{
				AwardID:   award.AwardID,
				AwardName: award.AwardName,
			}
			if prizeID, ok := am[award.AwardID]; ok {
				wr.PrizeID = prizeID
				wr.State = 1
			}
			rs = append(rs, wr)
		}
	}
	return
}

// GetWinningPoster get prize winning poster
func (s *Service) GetWinningPoster(c context.Context, mid int64, awardID, prizeID int64) (poster *model.Poster, err error) {
	accs, err := s.dao.AccountInfos(c, []int64{mid})
	if err != nil {
		return
	}
	if len(accs) <= 0 {
		return
	}

	award, err := s.dao.GetAwardSchedule(c, awardID)
	if err != nil {
		return
	}

	// am map[award_id]division_name
	names, err := s.dao.DivisionName(c, mid)
	if err != nil {
		return
	}

	bonus, err := s.dao.AwardBonus(c, awardID, prizeID)
	if err != nil {
		return
	}

	poster = &model.Poster{
		AwardName: award.AwardName,
		Nickname:  accs[mid].Nickname,
		Face:      accs[mid].Face,
		PrizeName: fmt.Sprintf("最佳%s新秀奖", names[awardID]),
		Date:      award.CycleEnd.Time().Format("2006-01"),
		Bonus:     bonus,
	}
	return
}

// JoinAward sign up award
func (s *Service) JoinAward(c context.Context, mid int64, awardID int64) (err error) {
	joined, err := s.isJoined(c, mid, awardID)
	if err != nil {
		return
	}

	if joined {
		err = ecode.GrowupSpecialAwardJoined
		return
	}

	fans, err := s.getUpFans(c, mid)
	if err != nil {
		return
	}

	if !checkQualification(fans) {
		err = ecode.GrowupSpecialAwardUnqualified
		return
	}
	_, err = s.dao.AddToAwardRecord(c, mid, awardID)
	return
}

// if fans count >= ceilings, no qualification
func checkQualification(fans int64) bool {
	return fans < ceilings
}

// if joined special award
func (s *Service) isJoined(c context.Context, mid, awardID int64) (joined bool, err error) {
	count, err := s.dao.JoinedCount(c, mid, awardID)
	if err != nil {
		return
	}
	joined = count != 0
	return
}

// AwardList award_id: award_name
func (s *Service) AwardList(c context.Context) (as []*model.SimpleSpecialAward, err error) {
	as, err = s.dao.PastAwards(c)
	if err != nil {
		return
	}

	sort.Slice(as, func(i, j int) bool {
		return as[i].CycleStart < as[j].CycleStart
	})
	return
}

// Winners get winners by award id
func (s *Service) Winners(c context.Context, awardID int64) (as []*model.Account, err error) {
	mids, err := s.dao.GetWinners(c, awardID)
	if err != nil {
		return
	}
	infos, err := s.dao.AccountInfos(c, mids)
	if err != nil {
		return
	}
	for _, mid := range mids {
		a := &model.Account{Mid: mid}
		as = append(as, a)
		if info, ok := infos[mid]; ok {
			a.Name = info.Nickname
			a.Face = info.Face
		}
	}
	return
}

// AwardDetail get award detail include schedule & resource
func (s *Service) AwardDetail(c context.Context, awardID int64) (data map[string]interface{}, err error) {
	schedule, err := s.dao.GetAwardSchedule(c, awardID)
	if err != nil {
		return
	}
	// rs map[resource_type]map[index]content
	rs, err := s.dao.GetResources(c, awardID)
	if err != nil {
		return
	}

	qas := make([]*model.QA, len(rs[question]))
	for i := 0; i < len(rs[question]); i++ {
		qas[i] = &model.QA{}
	}

	res := map[string]interface{}{
		"qa":     qas,
		"rule":   "",
		"detail": "",
	}

	for rt, cs := range rs {
		if rt == rule {
			res["rule"] = cs[1]
		}
		if rt == detail {
			res["detail"] = cs[1]
		}
		if rt == question || rt == answer {
			for index, content := range cs {
				if index < 0 || index > len(qas) {
					continue
				}
				qa := qas[index-1]
				if rt == question {
					qa.Question = content
				}
				if rt == answer {
					qa.Answer = content
				}
			}
		}
	}

	data = map[string]interface{}{
		"schedule": schedule,
		"resource": res,
	}
	return
}

// SpecialAwardInfo special award info
func (s *Service) SpecialAwardInfo(c context.Context, mid int64) (data map[string]interface{}, err error) {
	var (
		nowTime   = xtime.Time(time.Now().Unix())
		winRecord []string
		upStates  []*model.UpAwardState
	)

	awards, nowAward, nextAward, err := s.getRecentSpecialAward(c, nowTime)
	if err != nil {
		log.Error("s.getRecentSpecialAward error(%v)", err)
		return
	}

	if mid > 0 {
		winRecord, err = s.getAwardWinRecord(c, mid, awards)
		if err != nil {
			log.Error("s.getAwardWinRecord error(%v)", err)
			return
		}
		upStates, err = s.getUpAwardState(c, mid, awards)
		if err != nil {
			log.Error("s.getUpAwardState error(%v)", err)
			return
		}
	}
	data = map[string]interface{}{
		"win_record": winRecord,
		"now":        nowAward,
		"next":       nextAward,
		"up_states":  upStates,
	}
	return
}

// get now and next special award
func (s *Service) getRecentSpecialAward(c context.Context, nowTime xtime.Time) (awards []*model.SpecialAward, nowAward, nextAward *model.SpecialAward, err error) {
	awards, err = s.dao.GetSpecialAwards(c)
	if err != nil {
		log.Error("s.dao.GetSpecialAwards error(%v)", err)
		return
	}
	sort.Slice(awards, func(i, j int) bool {
		return awards[i].CycleStart < awards[j].CycleStart
	})
	for i := 0; i < len(awards); i++ {
		if awards[i].CycleStart <= nowTime && awards[i].CycleEnd >= nowTime {
			nowAward = awards[i]
		} else if awards[i].CycleStart > nowTime {
			nextAward = awards[i]
			break
		}
	}
	if nowAward != nil {
		nowAward.Duration = int64(nowAward.CycleEnd - nowTime)
		nowAward.Divisions, err = s.dao.GetSpecialAwardDivision(c, nowAward.AwardID)
		if err != nil {
			return
		}
	}
	if nextAward != nil {
		nextAward.Duration = int64(nextAward.CycleStart - nowTime)
		nextAward.Divisions, err = s.dao.GetSpecialAwardDivision(c, nextAward.AwardID)
		if err != nil {
			return
		}
	}
	return
}

func (s *Service) getAwardWinRecord(c context.Context, mid int64, awards []*model.SpecialAward) (awardNames []string, err error) {
	awardIDs, err := s.dao.GetAwardWinRecord(c, mid)
	if err != nil {
		log.Error("s.dao.GetAwardWinRecord error(%v)", err)
		return
	}
	awardNames = make([]string, 0)
	for i := len(awards) - 1; i >= 0; i-- {
		if awardIDs[awards[i].AwardID] {
			awardNames = append(awardNames, awards[i].AwardName)
		}
	}
	return
}

func (s *Service) getUpAwardState(c context.Context, mid int64, awards []*model.SpecialAward) (upStates []*model.UpAwardState, err error) {
	upStates = make([]*model.UpAwardState, 0)
	now := xtime.Time(time.Now().Unix())
	awardIDs, err := s.dao.GetAwardJoinRecord(c, mid)
	if err != nil {
		log.Error("s.dao.GetAwardJoinRecord error(%v)", err)
		return
	}
	winIDs, err := s.dao.GetAwardWinRecord(c, mid)
	if err != nil {
		log.Error("s.dao.GetAwardWinRecord error(%v)", err)
		return
	}
	for i := 0; i < len(awards); i++ {
		upState := &model.UpAwardState{AwardName: awards[i].AwardName}
		date := awards[i].AnnounceDate.Time()
		doubleCreativeStart := time.Date(date.Year(), date.Month()+1, 15, 0, 0, 0, 0, time.Local)
		doubleCreativeEnd := doubleCreativeStart.AddDate(0, 0, 14)
		if now > awards[i].CycleEnd && now <= awards[i].AnnounceDate && awardIDs[awards[i].AwardID] { // 评选中
			upState.State = 1
		} else if now > awards[i].AnnounceDate && now < xtime.Time(doubleCreativeStart.Unix()) && winIDs[awards[i].AwardID] { // 双倍即将开始
			upState.State = 2
		} else if now >= xtime.Time(doubleCreativeStart.Unix()) && now <= xtime.Time(doubleCreativeEnd.Unix()) && winIDs[awards[i].AwardID] { // 双倍中
			upState.State = 3
		}
		if upState.State != 0 {
			upStates = append(upStates, upState)
		}
	}
	return
}
