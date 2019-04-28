package bws

import (
	"context"
	"sort"
	"time"

	bwsmdl "go-common/app/interface/main/activity/model/bws"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func (s *Service) points(c context.Context, bid int64) (rs map[string][]*bwsmdl.Point, err error) {
	var (
		points                 *bwsmdl.Points
		dp, game, clockin, egg []*bwsmdl.Point
	)
	if points, err = s.dao.Points(c, bid); err != nil || points == nil || len(points.Points) == 0 {
		log.Error("s.dao.Points error(%v)", err)
		err = ecode.ActivityPointFail
		return
	}
	for _, point := range points.Points {
		switch point.LockType {
		case _dpType:
			dp = append(dp, point)
		case _gameType:
			game = append(game, point)
		case _clockinType:
			clockin = append(clockin, point)
		case _eggType:
			egg = append(egg, point)
		}
	}
	rs = make(map[string][]*bwsmdl.Point, 4)
	if len(dp) == 0 {
		rs[_dp] = _emptPoints
	} else {
		rs[_dp] = dp
	}
	if len(game) == 0 {
		rs[_game] = _emptPoints
	} else {
		rs[_game] = game
	}
	if len(clockin) == 0 {
		rs[_clockin] = _emptPoints
	} else {
		rs[_clockin] = clockin
	}
	if len(egg) == 0 {
		rs[_egg] = _emptPoints
	} else {
		rs[_egg] = egg
	}
	return
}

// Points points list
func (s *Service) Points(c context.Context, p *bwsmdl.ParamPoints) (rs map[string][]*bwsmdl.Point, err error) {
	var points map[string][]*bwsmdl.Point
	if points, err = s.points(c, p.Bid); err != nil {
		return
	}
	rs = make(map[string][]*bwsmdl.Point)
	switch p.Tp {
	case _allType:
		rs = points
	case _dpType:
		rs[_dp] = points[_dp]
	case _gameType:
		rs[_game] = points[_game]
	case _clockinType:
		rs[_clockin] = points[_clockin]
	case _eggType:
		rs[_egg] = points[_egg]
	}
	return
}

// Point point
func (s *Service) Point(c context.Context, p *bwsmdl.ParamID) (rs *bwsmdl.Point, err error) {
	var (
		points *bwsmdl.Points
	)
	if points, err = s.dao.Points(c, p.Bid); err != nil || points == nil || len(points.Points) == 0 {
		log.Error("s.dao.Points error(%v)", err)
		err = ecode.ActivityPointFail
		return
	}
	for _, point := range points.Points {
		if point.ID == p.ID {
			rs = point
			break
		}
	}
	if rs == nil {
		err = ecode.ActivityIDNotExists
	}
	return
}

// Unlock unlock point.
func (s *Service) Unlock(c context.Context, owner int64, arg *bwsmdl.ParamUnlock) (err error) {
	var (
		point         *bwsmdl.Point
		userPoints    []*bwsmdl.UserPointDetail
		userAchieves  []*bwsmdl.UserAchieveDetail
		achieves      *bwsmdl.Achievements
		unLockCnt, hp int64
		addAchieve    *bwsmdl.Achievement
		lockAchieves  []*bwsmdl.Achievement
	)
	if arg.Key == "" {
		if arg.Key, err = s.midToKey(c, arg.Mid); err != nil {
			return
		}
	}
	if point, err = s.Point(c, &bwsmdl.ParamID{ID: arg.Pid, Bid: arg.Bid}); err != nil {
		return
	}
	if point.Ower != owner && !s.isAdmin(owner) {
		err = ecode.ActivityNotOwner
		return
	}
	if point.LockType == _gameType {
		if arg.GameResult != bwsmdl.GameResWin && arg.GameResult != bwsmdl.GameResFail {
			err = ecode.ActivityGameResult
			return
		}
	}
	if userPoints, err = s.userPoints(c, arg.Bid, arg.Key); err != nil {
		return
	}
	userPidMap := make(map[int64]int64, len(userPoints))
	for _, v := range userPoints {
		if point.LockType != _gameType && v.Pid == point.ID {
			err = ecode.ActivityHasUnlock
			return
		}
		if _, ok := userPidMap[v.Pid]; !ok && v.LockType == point.LockType {
			if v.LockType == _gameType {
				if v.Points == v.Unlocked {
					unLockCnt++
					userPidMap[v.Pid] = v.Pid
				}
			} else {
				unLockCnt++
				userPidMap[v.Pid] = v.Pid
			}
		}
		hp += v.Points
	}
	lockPoint := point.Unlocked
	if point.LockType == _gameType && arg.GameResult == bwsmdl.GameResFail {
		lockPoint = point.LoseUnlocked
	}
	if hp+lockPoint < 0 {
		err = ecode.ActivityLackHp
		return
	}
	if userAchieves, err = s.userAchieves(c, arg.Bid, arg.Key); err != nil {
		return
	}
	if err = s.addUserPoint(c, arg.Bid, arg.Pid, lockPoint, arg.Key); err != nil {
		return
	}
	if achieves, err = s.dao.Achievements(c, arg.Bid); err != nil || len(achieves.Achievements) == 0 {
		log.Error("s.dao.Achievements error(%v)", err)
		err = ecode.ActivityAchieveFail
		return
	}
	for _, v := range achieves.Achievements {
		if point.LockType == v.LockType {
			lockAchieves = append(lockAchieves, v)
		}
	}
	if len(lockAchieves) > 0 {
		sort.Slice(lockAchieves, func(i, j int) bool { return lockAchieves[i].Unlock > lockAchieves[j].Unlock })
		if point.LockType == _gameType {
			if arg.GameResult == bwsmdl.GameResWin {
				unLockCnt++
			}
		} else {
			unLockCnt++
		}
		for _, ach := range lockAchieves {
			if unLockCnt >= ach.Unlock {
				addAchieve = ach
				break
			}
		}
	}
	if addAchieve != nil {
		for _, v := range userAchieves {
			if v.Aid == addAchieve.ID {
				return
			}
		}
		s.addAchieve(c, arg.Mid, addAchieve, arg.Key)
	}
	return
}

func (s *Service) userPoints(c context.Context, bid int64, key string) (res []*bwsmdl.UserPointDetail, err error) {
	var (
		usPoints []*bwsmdl.UserPoint
		points   *bwsmdl.Points
	)
	if usPoints, err = s.dao.UserPoints(c, bid, key); err != nil {
		err = ecode.ActivityUserPointFail
		return
	}
	if len(usPoints) == 0 {
		return
	}
	if points, err = s.dao.Points(c, bid); err != nil || points == nil || len(points.Points) == 0 {
		log.Error("s.dao.Points error(%v)", err)
		err = ecode.ActivityPointFail
		return
	}
	pointsMap := make(map[int64]*bwsmdl.Point, len(points.Points))
	for _, v := range points.Points {
		pointsMap[v.ID] = v
	}
	for _, v := range usPoints {
		detail := &bwsmdl.UserPointDetail{UserPoint: v}
		if point, ok := pointsMap[v.Pid]; ok {
			detail.Name = point.Name
			detail.Icon = point.Icon
			detail.Fid = point.Fid
			detail.Image = point.Image
			detail.Unlocked = point.Unlocked
			detail.LockType = point.LockType
			detail.Dic = point.Dic
			detail.Rule = point.Rule
			detail.Bid = point.Bid
		}
		res = append(res, detail)
	}
	return
}

func (s *Service) addUserPoint(c context.Context, bid, pid, points int64, key string) (err error) {
	var usPtID int64
	if usPtID, err = s.dao.AddUserPoint(c, bid, pid, points, key); err != nil {
		err = ecode.ActivityUnlockFail
		return
	}
	err = s.dao.AppendUserPointsCache(c, bid, key, &bwsmdl.UserPoint{ID: usPtID, Pid: pid, Points: points, Ctime: xtime.Time(time.Now().Unix())})
	return
}
