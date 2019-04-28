package service

import (
	"context"
	"math"
	"time"

	"go-common/app/job/main/figure-timer/conf"
	"go-common/app/job/main/figure-timer/model"
	"go-common/library/log"
)

const (
	bizTypePosLawful = iota
	bizTypeNegLawful
	// bizTypePosWide
	// bizTypeNegWide
	bizTypePosFriendly
	bizTypeNegFriendly
	bizTypePosCreativity
	// bizTypeNegCreativity
	bizTypePosBounty
	// bizTypeNegBounty
)

// HandleFigure handle all figure score for a mid
func (s *Service) HandleFigure(c context.Context, mid int64, weekVer int64) (err error) {
	var (
		figure             = &model.Figure{}
		userInfo           *model.UserInfo
		actionCounter      *model.ActionCounter
		records            []*model.FigureRecord
		newRecord          *model.FigureRecord
		weekVerRecordsFrom = time.Unix(weekVer, 0).AddDate(0, 0, -7*52).Unix()
		weekVerRecordsTo   = time.Unix(weekVer, 0).AddDate(0, 0, -7).Unix() + 1
	)
	//1. get user_info from hbase
	if userInfo, err = s.dao.UserInfo(c, mid, weekVer); err != nil {
		return
	}
	//2. get action_counter from hbase
	if actionCounter, err = s.dao.ActionCounter(c, mid, weekVer); err != nil {
		return
	}
	//3. get figure_records from hbase
	if records, err = s.dao.CalcRecords(c, mid, weekVerRecordsFrom, weekVerRecordsTo); err != nil {
		return
	}
	//4. calc figure
	figure, newRecord = s.CalcFigure(c, userInfo, []*model.ActionCounter{actionCounter}, records, weekVer)
	log.Info("User figure [%+v]", figure)
	log.Info("User newRecord [%+v]", newRecord)
	//5. save to db
	time.Sleep(time.Millisecond)
	if figure.ID, err = s.dao.UpsertFigure(c, figure); err != nil {
		log.Error("%+v", err)
	}
	//6. save new record
	if err = s.dao.PutCalcRecord(c, newRecord, weekVer); err != nil {
		log.Error("%+v", err)
	}
	//7. remove existed redis
	if err = s.dao.RemoveCache(c, mid); err != nil {
		log.Error("%+v", err)
	}
	rank.AddScore(figure.Score)
	return
}

// CalcFigure calc figure.
func (s *Service) CalcFigure(c context.Context, userInfo *model.UserInfo, actionCounters []*model.ActionCounter, records []*model.FigureRecord, weekVer int64) (figure *model.Figure, newRecord *model.FigureRecord) {
	figure = &model.Figure{Mid: userInfo.Mid, Ver: weekVer}
	newRecord = &model.FigureRecord{Mid: userInfo.Mid, Version: time.Unix(weekVer, 0)}
	var (
		posx                 float64
		negx                 float64
		newPosx, newNegx     float64
		k1, k2, k3, k4, k5   float64 = s.c.Property.Calc.K1, s.c.Property.Calc.K2, s.c.Property.Calc.K3, s.c.Property.Calc.K4, s.c.Property.Calc.K5
		posOffset, negOffset int64
		lawfulBase           = s.c.Property.Calc.InitLawfulScore
		lawfulPosMax         = s.c.Property.Calc.LawfulPosMax
		lawfulNegMax         = s.c.Property.Calc.LawfulNegMax
		lawfulPosK           = s.c.Property.Calc.LawfulPosK
		lawfulNegK1          = s.c.Property.Calc.LawfulNegK1
		lawfulNegK2          = s.c.Property.Calc.LawfulNegK2
		lawfulPosL           = s.c.Property.Calc.LawfulPosL
		lawfulNegL           = s.c.Property.Calc.LawfulNegL
		lawfulPosC3          = s.c.Property.Calc.LawfulPosC3
		lawfulNegC1          = s.c.Property.Calc.LawfulNegC1
		lawfulPosQ3          = s.c.Property.Calc.LawfulPosQ3
		lawfulNegQ1          = s.c.Property.Calc.LawfulNegQ1
		wideBase             = s.c.Property.Calc.InitWideScore
		widePosMax           = s.c.Property.Calc.WidePosMax
		widePosK             = s.c.Property.Calc.WidePosK
		wideC1               = s.c.Property.Calc.WideC1 //有播放的活跃天数
		wideQ1               = s.c.Property.Calc.WideQ1
		wideC2               = s.c.Property.Calc.WideC2 // 账号累计经验值
		wideQ2               = s.c.Property.Calc.WideQ2
		friendlyBase         = s.c.Property.Calc.InitFriendlyScore
		friendlyPosMax       = s.c.Property.Calc.FriendlyPosMax
		friendlyNegMax       = s.c.Property.Calc.FriendlyNegMax
		friendlyPosK         = s.c.Property.Calc.FriendlyPosK
		friendlyNegK         = s.c.Property.Calc.FriendlyNegK
		friendlyPosL         = s.c.Property.Calc.FriendlyPosL
		friendlyNegL         = s.c.Property.Calc.FriendlyNegL
		bountyBase           = s.c.Property.Calc.InitBountyScore
		bountyMax            = s.c.Property.Calc.BountyMax
		bountyPosL           = s.c.Property.Calc.BountyPosL
		bountyK              = s.c.Property.Calc.BountyK
		bountyQ1             = s.c.Property.Calc.BountyQ1
		bountyC1             = s.c.Property.Calc.BountyC1
		creativityBase       = s.c.Property.Calc.InitCreativityScore
		creativityPosMax     = s.c.Property.Calc.CreativityPosMax
		creativityPosK       = s.c.Property.Calc.CreativityPosK
		creativityPosL1      = s.c.Property.Calc.CreativityPosL1
	)
	//1. lawful
	newPosx, posx = s.calcActionX(lawfulPosL, bizTypePosLawful, actionCounters, records, weekVer)
	posx += lawfulPosQ3 * lawfulPosC3 * float64(userInfo.DisciplineCommittee)
	posOffset = calcOffset(lawfulPosMax, lawfulPosK, posx)

	spyScore := userInfo.SpyScore
	negx = lawfulNegQ1 * lawfulNegC1 * float64(80-float64(spyScore))
	if negx < 0 {
		negx = 0.0
	}
	negOffset = calcOffset(lawfulNegMax, lawfulNegK1, negx)
	newNegx, negx = s.calcActionX(lawfulNegL, bizTypeNegLawful, actionCounters, records, weekVer)
	negOffset += calcOffset(lawfulNegMax, lawfulNegK2, negx)
	figure.LawfulScore = int32(lawfulBase + posOffset - negOffset)
	if figure.LawfulScore < 0 {
		figure.LawfulScore = 0
	}
	newRecord.XPosLawful, newRecord.XNegLawful = int64(newPosx), int64(newNegx)

	//2. wide
	newPosx, newNegx = 0, 0
	posx = wideQ1*wideC1*float64(userInfo.ArchiveViews) + wideQ2*wideC2*float64(userInfo.Exp)
	posOffset = calcOffset(widePosMax, widePosK, posx)
	negx = 0
	negOffset = 0
	figure.WideScore = int32(wideBase + posOffset - negOffset)
	if figure.WideScore < 0 {
		figure.WideScore = 0
	}
	newRecord.XPosWide, newRecord.XNegWide = int64(newPosx), int64(newNegx)

	//3. friendly
	newPosx, newNegx = 0, 0
	newPosx, posx = s.calcActionX(friendlyPosL, bizTypePosFriendly, actionCounters, records, weekVer)
	posOffset = calcOffset(friendlyPosMax, friendlyPosK, posx)
	newNegx, negx = s.calcActionX(friendlyNegL, bizTypeNegFriendly, actionCounters, records, weekVer)
	negOffset = calcOffset(friendlyNegMax, friendlyNegK, negx)
	figure.FriendlyScore = int32(friendlyBase + posOffset - negOffset)
	if figure.FriendlyScore < 0 {
		figure.FriendlyScore = 0
	}
	newRecord.XPosFriendly, newRecord.XNegFriendly = int64(newPosx), int64(newNegx)

	//4. bounty
	newPosx, newNegx = 0, 0
	var bountyVIP float64
	if userInfo.VIPStatus > 0 {
		bountyVIP = 1
	}
	newPosx, posx = s.calcActionX(bountyPosL, bizTypePosBounty, actionCounters, records, weekVer)
	posx += bountyQ1 * bountyC1 * bountyVIP
	posOffset = calcOffset(bountyMax, bountyK, posx)
	negx = 0
	negOffset = 0
	figure.BountyScore = int32(bountyBase + posOffset - negOffset)
	if figure.BountyScore < 0 {
		figure.BountyScore = 0
	}
	newRecord.XPosBounty, newRecord.XNegBounty = int64(newPosx), int64(newNegx)

	//5. creativity
	newPosx, newNegx = 0, 0
	newPosx, posx = s.calcActionX(creativityPosL1, bizTypePosCreativity, actionCounters, records, weekVer)
	posOffset = calcOffset(creativityPosMax, creativityPosK, posx)
	negx = 0
	negOffset = 0
	figure.CreativityScore = int32(creativityBase + posOffset - negOffset)
	if figure.CreativityScore < 0 {
		figure.CreativityScore = 0
	}
	newRecord.XPosCreativity, newRecord.XNegCreativity = int64(newPosx), int64(newNegx)

	//6. calc score
	figure.Score = int32(math.Floor(k1*float64(figure.LawfulScore) + k2*float64(figure.WideScore) + k3*float64(figure.FriendlyScore) + k4*float64(figure.CreativityScore) + k5*float64(figure.BountyScore)))
	return
}

// x must >= 0
func calcOffset(max, k, x float64) (score int64) {
	return int64(math.Floor(max * (1 - math.Pow(math.E, -(k*x)))))
}

func (s *Service) calcActionX(L float64, bizType int8, actionCounters []*model.ActionCounter, records []*model.FigureRecord, weekVer int64) (newx, totalx float64) {
	var (
		day int64 = 24 * 3600
		t   float64

		lawfulPosC1 = s.c.Property.Calc.LawfulPosC1
		lawfulPosC2 = s.c.Property.Calc.LawfulPosC2
		lawfulPosQ1 = s.c.Property.Calc.LawfulPosQ1
		lawfulPosQ2 = s.c.Property.Calc.LawfulPosQ2
		lawfulNegC2 = s.c.Property.Calc.LawfulNegC2
		lawfulNegC3 = s.c.Property.Calc.LawfulNegC3
		lawfulNegQ2 = s.c.Property.Calc.LawfulNegQ2
		lawfulNegQ3 = s.c.Property.Calc.LawfulNegQ3

		friendlyPosQ1 = s.c.Property.Calc.FriendlyPosQ1
		friendlyPosC1 = s.c.Property.Calc.FriendlyPosC1
		friendlyPosQ2 = s.c.Property.Calc.FriendlyPosQ2
		friendlyPosC2 = s.c.Property.Calc.FriendlyPosC2
		friendlyPosQ3 = s.c.Property.Calc.FriendlyPosQ3
		friendlyPosC3 = s.c.Property.Calc.FriendlyPosC3
		friendlyNegQ1 = s.c.Property.Calc.FriendlyNegQ1
		friendlyNegC1 = s.c.Property.Calc.FriendlyNegC1
		friendlyNegQ2 = s.c.Property.Calc.FriendlyNegQ2
		friendlyNegC2 = s.c.Property.Calc.FriendlyNegC2
		friendlyNegQ3 = s.c.Property.Calc.FriendlyNegQ3
		friendlyNegC3 = s.c.Property.Calc.FriendlyNegC3
		friendlyNegQ4 = s.c.Property.Calc.FriendlyNegQ4
		friendlyNegC4 = s.c.Property.Calc.FriendlyNegC4

		creativityQ1 = s.c.Property.Calc.CreativityQ1
		creativityC1 = s.c.Property.Calc.CreativityC1

		bountyQ2 = s.c.Property.Calc.BountyQ2
		bountyC2 = s.c.Property.Calc.BountyC2
		bountyQ3 = s.c.Property.Calc.BountyQ3
		bountyC3 = s.c.Property.Calc.BountyC3
	)
	if L == 0.0 {
		return 0, 0
	}
	for _, ac := range actionCounters {
		t = float64(7 - (ac.Version.Unix()-weekVer)/day)
		if t <= 0 {
			continue
		}
		switch bizType {
		case bizTypePosLawful:
			if ac.ReportReplyPassed < 0 {
				ac.ReportReplyPassed = 0
			}
			newx += actionX(lawfulPosQ1, lawfulPosC1, float64(ac.ReportReplyPassed), t, L)
			if ac.ReportDanmakuPassed < 0 {
				ac.ReportDanmakuPassed = 0
			}
			newx += actionX(lawfulPosQ2, lawfulPosC2, float64(ac.ReportDanmakuPassed), t, L)
		case bizTypeNegLawful:
			if ac.PublishReplyDeleted < 0 {
				ac.PublishReplyDeleted = 0
			}
			newx += actionX(lawfulNegQ2, lawfulNegC2, float64(ac.PublishReplyDeleted), t, L)
			if ac.PublishDanmakuDeleted < 0 {
				ac.PublishDanmakuDeleted = 0
			}
			newx += actionX(lawfulNegQ3, lawfulNegC3, float64(ac.PublishDanmakuDeleted), t, L)
		case bizTypePosFriendly:
			newx += actionX(friendlyPosQ1, friendlyPosC1, float64(ac.CoinCount), t, L)
			newx += actionX(friendlyPosQ2, friendlyPosC2, float64(ac.ReplyCount), t, L)
			newx += actionX(friendlyPosQ3, friendlyPosC3, float64(ac.DanmakuCount), t, L)
		case bizTypeNegFriendly:
			newx += actionX(friendlyNegQ1, friendlyNegC1, float64(ac.CoinHighRisk), t, L)
			newx += actionX(friendlyNegQ2, friendlyNegC2, float64(ac.CoinLowRisk), t, L)
			if ac.PublishReplyDeleted < 0 {
				ac.PublishReplyDeleted = 0
			}
			newx += actionX(friendlyNegQ3, friendlyNegC3, float64(ac.PublishReplyDeleted), t, L)
			if ac.PublishDanmakuDeleted < 0 {
				ac.PublishDanmakuDeleted = 0
			}
			newx += actionX(friendlyNegQ4, friendlyNegC4, float64(ac.PublishDanmakuDeleted), t, L)
		case bizTypePosCreativity:
			replyLikeCount := float64(ac.ReplyLiked) - float64(ac.ReplyUnliked)
			if replyLikeCount < 0 {
				replyLikeCount = 0.0
			}
			newx += actionX(creativityQ1, creativityC1, replyLikeCount, t, L)
		case bizTypePosBounty:
			newx += actionX(bountyQ2, bountyC2, float64(ac.PayMoney), t, L)
			newx += actionX(bountyQ3, bountyC3, float64(ac.PayLiveMoney), t, L)
		}
	}
	totalx = newx
	for _, r := range records {
		t = float64((weekVer - r.Version.Unix()) / day)
		if t <= 0 {
			continue
		}
		switch bizType {
		case bizTypePosLawful:
			totalx += actionX(1, 1, float64(r.XPosLawful), t, L)
		case bizTypeNegLawful:
			totalx += actionX(1, 1, float64(r.XNegLawful), t, L)
		case bizTypePosFriendly:
			totalx += actionX(1, 1, float64(r.XPosFriendly), t, L)
		case bizTypeNegFriendly:
			totalx += actionX(1, 1, float64(r.XNegFriendly), t, L)
		case bizTypePosCreativity:
			totalx += actionX(1, 1, float64(r.XPosCreativity), t, L)
		case bizTypePosBounty:
			totalx += actionX(1, 1, float64(r.XPosBounty), t, L)
		}
	}
	return
}

func actionX(q, c, x, t, L float64) float64 {
	return q * c * x * math.Pow(math.E, -math.Pow((t/L), 2))
}

// InitFigure initialize user figure
func (s *Service) InitFigure(c context.Context, mid int64, ver string) (figure *model.Figure, err error) {
	figure = &model.Figure{}
	figure.Mid = mid
	figure.LawfulScore = int32(conf.Conf.Property.Calc.InitLawfulScore)
	figure.WideScore = int32(conf.Conf.Property.Calc.InitWideScore)
	figure.FriendlyScore = int32(conf.Conf.Property.Calc.InitFriendlyScore)
	figure.BountyScore = int32(conf.Conf.Property.Calc.InitBountyScore)
	figure.CreativityScore = int32(conf.Conf.Property.Calc.InitCreativityScore)
	figure.Ver = s.curVer

	if figure.ID, err = s.dao.UpsertFigure(c, figure); err != nil {
		return
	}
	if err = s.dao.SetFigureCache(c, figure); err != nil {
		log.Error("%+v", err)
	}
	return
}

// get ever monday start time ts.
func weekVersion(now time.Time) (ts int64) {
	var (
		wd int
		w  time.Weekday
	)
	w = now.Weekday()
	switch w {
	case time.Sunday:
		wd = 6
	default:
		wd = int(w) - 1
	}
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -wd).Unix()
}
