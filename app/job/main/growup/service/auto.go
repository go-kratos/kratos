package service

import (
	"context"
	"fmt"
	"time"

	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_video         = 0
	_forbidDay     = 10
	_breachReason  = "非自制"
	_forbidReason  = "投递非本人自制稿件2次"
	_dismissReason = "投递非本人自制稿件3次及以上"
)

// AutoBreach auto breach
func (s *Service) AutoBreach(c context.Context, date string) (msg string, err error) {
	return s.autoAvBreach(c, date)
}

func (s *Service) autoAvBreach(c context.Context, date string) (msg string, err error) {
	avs, err := s.dao.GetAvBreachPre(c, _video, 1, date)
	if err != nil {
		log.Error("s.dao.GetAvBreachPre error(%v)", err)
		return
	}

	needBreach := make(map[int64]bool)
	for _, av := range avs {
		err = s.dao.DoAvBreach(c, av.MID, av.AvID, _video, _breachReason)
		if err != nil {
			log.Error("s.dao.DoAvBreach error(%v)", err)
			return
		}
		needBreach[av.AvID] = true
	}

	breach := make([]int64, 0)
	breachAvs, err := s.dao.GetAvBreachPre(c, _video, 2, date)
	if err != nil {
		log.Error("s.dao.GetAvBreachPre error(%v)", err)
		return
	}
	for _, av := range breachAvs {
		if needBreach[av.AvID] {
			breach = append(breach, av.AvID)
		}
	}
	msg = fmt.Sprintf("%s 自制转转载违规扣除稿件:%s", date, xstr.JoinInts(breach))
	return
}

// AutoPunish auto punish
func (s *Service) AutoPunish(c context.Context) (msg string, err error) {
	return s.autoUpPunish(c)
}

// 每周一检查并处罚 只检查非自制
func (s *Service) autoUpPunish(c context.Context) (msg string, err error) {
	now := time.Now()
	if now.Format(_layout) != getStartWeeklyDate(now).Format(_layout) {
		return
	}
	lastWeek := getStartWeeklyDate(getStartWeeklyDate(now).AddDate(0, 0, -1))
	avs, err := s.dao.GetAvBreach(c, "2018-10-15", now.AddDate(0, 0, -1).Format(_layout)) // 从2018-10-15上线开始计算扣除
	if err != nil {
		log.Error("s.dao.GetAvBreach error(%v)", err)
		return
	}

	blacks, err := s.listBlacklist(c, "")
	if err != nil {
		log.Error("s.listBlacklist error(%v)", err)
		return
	}
	mBlacks := make(map[int64]bool)
	for _, b := range blacks {
		mBlacks[b.AvID] = true
	}

	mids := make(map[int64]map[string]struct{})
	for _, av := range avs {
		if av.Reason != _breachReason {
			continue
		}
		if !mBlacks[av.AvID] {
			continue
		}
		date := getStartWeeklyDate(av.Date.Time()).Format(_layout)
		if _, ok := mids[av.MID]; !ok {
			mids[av.MID] = make(map[string]struct{})
		}
		mids[av.MID][date] = struct{}{}
	}

	forbidMIDs, dismissMIDs := make([]int64, 0), make([]int64, 0)
	for mid, times := range mids {
		// 当且仅当上周有发生过扣除才处罚up主，防止多次处罚
		if _, ok := times[lastWeek.Format(_layout)]; !ok {
			continue
		}
		if len(times) == 2 {
			err = s.autoForbid(c, mid)
			if err != nil {
				log.Error("s.autoForbid(%d) error(%v)", mid, err)
				return
			}
			forbidMIDs = append(forbidMIDs, mid)
		}
		if len(times) >= 3 {
			err = s.autoDismiss(c, mid)
			if err != nil {
				log.Error("s.autoDismiss(%d) error(%v)", mid, err)
				return
			}
			dismissMIDs = append(dismissMIDs, mid)
		}
	}
	msg = fmt.Sprintf("%s 自制转转载处罚up主: 封禁(%s), 清退(%s)", time.Now().Format(_layout), xstr.JoinInts(forbidMIDs), xstr.JoinInts(dismissMIDs))
	return
}

func (s *Service) autoForbid(c context.Context, mid int64) (err error) {
	accState, err := s.dao.GetUpStateByMID(c, mid)
	if err != nil {
		log.Error(" s.dao.GetUpStateByMID(%d) error(%v)", mid, err)
		return
	}
	if !(accState == 3 || accState == 7) {
		return
	}
	err = s.dao.DoUpForbid(c, mid, _forbidDay, _video, _forbidReason)
	if err != nil {
		log.Error("s.dao.DoUpForbid error(%v)", err)
	}
	return
}

func (s *Service) autoDismiss(c context.Context, mid int64) (err error) {
	accState, err := s.dao.GetUpStateByMID(c, mid)
	if err != nil {
		log.Error(" s.dao.GetUpStateByMID(%d) error(%v)", mid, err)
		return
	}
	if accState != 3 && accState != 7 {
		return
	}
	err = s.dao.DoUpDismiss(c, mid, _video, _dismissReason)
	if err != nil {
		log.Error("s.dao.DoUpDismiss error(%v)", err)
	}
	return
}

// AutoExamination auto examination
func (s *Service) AutoExamination(c context.Context) (msg string, err error) {
	return s.autoExamination(c)
}

func (s *Service) autoExamination(c context.Context) (msg string, err error) {
	ups, err := s.getAllUps(c, 2000)
	if err != nil {
		log.Error("s.getAllUps error(%v)", err)
		return
	}
	mids := make([]int64, 0, len(ups))
	for mid, up := range ups {
		if up.AccountState != 2 || up.IsDeleted == 1 {
			continue
		}
		if up.Fans < 10000 && up.TotalPlayCount < 500000 {
			continue
		}
		mids = append(mids, mid)
	}
	scores, err := s.dao.GetUpCreditScore(c, mids)
	if err != nil {
		log.Error("s.dao.GetUpCreditScore error(%v)", err)
		return
	}
	passMID := make([]int64, 0, len(mids))
	for _, mid := range mids {
		if score, ok := scores[mid]; ok {
			if score < 100 {
				continue
			}
		}
		passMID = append(passMID, mid)
	}
	err = s.doUpPass(c, passMID, _video)
	if err != nil {
		log.Error("s.doUpPass error(%v)", err)
		return
	}
	msg = fmt.Sprintf("%s 自动过审up主: %s", time.Now().Format(_layout), xstr.JoinInts(passMID))
	return
}

func (s *Service) doUpPass(c context.Context, mids []int64, ctype int) (err error) {
	start, end := 0, 100
	for {
		if start >= len(mids) {
			break
		}
		if end > len(mids) {
			end = len(mids)
		}
		err = s.dao.DoUpPass(c, mids[start:end], _video)
		if err != nil {
			log.Error("s.dao.DoUpPass error(%v)", err)
			return
		}
		start = end
		end += 100
	}
	return
}

func getStartWeeklyDate(date time.Time) time.Time {
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
	}
	return date
}
