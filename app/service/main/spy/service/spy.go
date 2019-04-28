package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/dao"
	"go-common/app/service/main/spy/model"
	spy "go-common/app/service/main/spy/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/ip"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// HandleEvent handle spy-event.
func (s *Service) HandleEvent(c context.Context, eventMsg *model.EventMessage) (err error) {
	var (
		factor *model.Factor
		mid    = eventMsg.ActiveMid
	)
	ui, err := s.UserInfo(c, mid, ip.InternalIP())
	if err != nil {
		log.Error("s.UserInfo(%d) err(%v)", mid, err)
		return
	}
	// if blocked already , return
	// if ui.State == model.StateBlock {
	// 	return
	// }
	// get factor by servieName , eventName , riskLevel
	if factor, err = s.factor(c, eventMsg.Service, eventMsg.Event, eventMsg.RiskLevel); err != nil {
		log.Error("s.factor(%s, %s, %d) error(%v)", eventMsg.Service, eventMsg.Event, eventMsg.RiskLevel, err)
		return
	}
	s.promEventScore.Incr("change")
	// mark eventScore decrease
	if factor.FactorVal < 1.0 {
		s.promEventScore.Incr("decrease")
	}
	// calc all score
	ui.EventScore = s.calcFactorScore(ui.EventScore, factor)
	ui.Score = s.calcScore(ui.BaseScore, ui.EventScore)
	if err = s.updateEventInfo(c, ui, factor, eventMsg); err != nil {
		log.Error("s.updateEventInfo(%v, %v, %v) error(%v)", ui, factor, eventMsg, err)
		return
	}
	s.updateInfoCache(c, ui)
	// check punishment.
	_, err = s.BlockFilter(c, ui)
	if err != nil {
		log.Error("s.BlockFilter(%v) error(%v)", ui, err)
		return
	}
	return
}

func (s *Service) updateEventInfo(c context.Context, ui *model.UserInfo, factor *model.Factor, eventMsg *model.EventMessage) (err error) {
	var (
		tx *sql.Tx
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	// update info.
	if err = s.dao.TxUpdateInfo(c, tx, ui); err != nil {
		log.Error("s.dao.UpsertInfo(%v) error(%v)", ui, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, factor, eventMsg); err != nil {
		log.Error("s.addHistory(%v, %v, %v) error(%v)", ui, factor, eventMsg, err)
		return
	}
	s.dao.PubScoreChange(c, ui.Mid, &model.ScoreChange{
		Mid:       ui.Mid,
		Score:     ui.Score,
		TS:        time.Now().Unix(),
		Reason:    eventMsg.Service,
		RiskLevel: eventMsg.RiskLevel,
	})
	return
}

// factor get facory by serviceName , eventName , riskLevel.
func (s *Service) factor(c context.Context, serviceName, eventName string, riskLevel int8) (factor *model.Factor, err error) {
	var (
		event *model.Event
	)
	if event, err = s.dao.Event(c, eventName); err != nil {
		log.Error("s.dao.Event(%s) error(%v)", eventName, err)
		return
	}
	if event == nil {
		log.Error("event(%s) not support", eventName)
		err = ecode.SpyEventNotExist
		return
	}
	if factor, err = s.dao.Factor(c, event.ServiceID, event.ID, riskLevel); err != nil {
		log.Error("s.dao.Factor(%d, %d, %d) error(%v)", event.ServiceID, event.ID, riskLevel, err)
		return
	}
	if factor == nil {
		log.Warn("factor(service=%d,event=%d,risklevel=%d) not support", event.ServiceID, event.ID, riskLevel)
		err = ecode.SpyFactorNotExist
		return
	}
	return
}

// UserInfo get UserInfo by mid , from cache or db or generate.
func (s *Service) UserInfo(c context.Context, mid int64, ip string) (ui *model.UserInfo, err error) {
	ui, err = s.userInfo(c, mid)
	if ui != nil {
		return
	}
	// user info not generated , so let's do it.
	if ui, err = s.initUserInfo(c, mid, ip); err != nil {
		log.Error("s.generateUserInfo(%d) error(%v)", mid, err)
		return
	}
	return
}

// UserInfoAsyn get UserInfo by mid , from cache or db or asyn generate.
func (s *Service) UserInfoAsyn(c context.Context, mid int64) (ui *model.UserInfo, err error) {
	ui, err = s.userInfo(c, mid)
	if ui != nil {
		return
	}
	ip := metadata.String(c, metadata.RemoteIP)
	s.infomission(func() {
		if _, err1 := s.initUserInfo(context.TODO(), mid, ip); err1 != nil {
			log.Error("s.generateUserInfo(%d) error(%v)", mid, err1)
			return
		}
	})
	return
}

// UserInfo get UserInfo by mid , from cache or db.
func (s *Service) userInfo(c context.Context, mid int64) (ui *model.UserInfo, err error) {
	var cacheFlag = true
	// get info from cache.
	if ui, err = s.dao.UserInfoCache(c, mid); err != nil {
		cacheFlag = false
		err = nil
		log.Error("s.dao.InfoCache(%d) error(%v)", mid, err)
	}
	// return if cache exist.
	if ui != nil {
		return
	}
	// get info from db.
	if ui, err = s.dao.UserInfo(c, mid); err != nil {
		log.Error("s.dao.Info(%d) error(%v)", mid, err)
		return
	}
	if ui != nil {
		// reload to cache
		if cacheFlag {
			s.updateInfoCache(c, ui)
		}
	}
	return
}

// Info returns user info.
func (s *Service) Info(c context.Context, mid int64) (ui *model.UserInfo, err error){
	ui, err = s.UserInfoAsyn(c, mid)
	if err != nil {
		return
	}
	if ui == nil {
		// asyn init user score , return def score.
		ui = &model.UserInfo{Mid: mid, Score: model.SpyInitScore}
	}
	return ui, nil
}

// initUserInfo .
func (s *Service) initUserInfo(c context.Context, mid int64, ip string) (ui *model.UserInfo, err error) {
	var (
		tx *sql.Tx
		// init default info
		fakeFactor = &model.Factor{
			NickName:  "initialization",
			ServiceID: -1,
			EventID:   conf.Conf.Property.Event.InitEventID,
			GroupID:   -1,
			RiskLevel: 1,
			FactorVal: 1.0,
		}
		fakeEventMsg = &model.EventMessage{
			IP:        "127.0.0.1",
			ActiveMid: mid,
			TargetMid: mid,
			Effect:    "initialization",
			RiskLevel: 1,
		}
	)
	// start generate logic.
	ui = &model.UserInfo{
		Mid:        mid,
		EventScore: conf.Conf.Property.Score.EventInit,
		State:      model.StateNormal,
	}
	if ui.BaseScore, err = s.getBaseScore(c, mid, ip); err != nil {
		log.Error("s.calcBaseScore(%d, %s) error(%v)", mid, ip, err)
		return
	}
	ui.Score = s.calcScore(ui.BaseScore, ui.EventScore)
	// add new info into db.
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if ui.ID, err = s.dao.TxAddInfo(c, tx, ui); err != nil {
		log.Error("s.userDao.NewInfo(%d, %v) error(%v)", mid, ui, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, fakeFactor, fakeEventMsg); err != nil {
		log.Error("s.addHistory(%v, %v, %v) error(%v)", ui, fakeFactor, fakeEventMsg, err)
		return
	}
	s.updateInfoCache(c, ui)
	s.dao.PubScoreChange(c, mid, &model.ScoreChange{
		Mid:   mid,
		Score: ui.Score,
		TS:    time.Now().Unix(),
	})
	return
}

func (s *Service) calcFactorScore(preScore int8, factor *model.Factor) (score int8) {
	return int8(float64(preScore) * factor.FactorVal)
}

// calc real score.
func (s *Service) calcScore(baseScore, eventScore int8) (score int8) {
	return int8(float64(baseScore) * float64(eventScore) / float64(_score))
}

// getBaseScore get base score.
func (s *Service) getBaseScore(c context.Context, mid int64, ip string) (score int8, err error) {
	var (
		judgementInfo *JudgementInfo
		factorMeta    FactoryMeta
		factor        *model.Factor
	)
	if judgementInfo, err = s.getJudgementInfo(c, mid, ip); err != nil {
		return
	}
	if factorMeta, err = s.getBaseScoreFactor(judgementInfo); err != nil {
		return
	}
	log.Info("getBaseScoreFactor mid(%d) meta(%v)", mid, factorMeta)
	if factor, err = s.factor(c, factorMeta.serviceName, factorMeta.eventName, factorMeta.riskLevel); err != nil {
		log.Error("s.factor(%s, %s) error(%v)", conf.Conf.Property.Event.ServiceName, conf.Conf.Property.Event.BindMailOnly, 1)
		return
	}
	score = s.calcFactorScore(conf.Conf.Property.Score.BaseInit, factor)
	return
}

// ReBuildPortrait reBuild user info.
func (s *Service) ReBuildPortrait(c context.Context, mid int64, reason string) (err error) {
	var (
		tx *sql.Tx
	)
	ui, err := s.UserInfo(c, mid, ip.InternalIP())
	if err != nil || ui == nil {
		log.Error("user portrait not fund, mid(%d, %v), err(%v)", mid, ui, err)
		return
	}
	s.promEventScore.Incr("change")
	s.resetScore(c, ui, ui.BaseScore, _score)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if err = s.dao.TxUpdateEventScoreReLive(c, tx, ui.Mid, ui.EventScore, ui.Score); err != nil {
		log.Error("s.dao.TxUpdateEventScoreReLive(%d, %d, %d), err(%v)", ui.Mid, ui.EventScore, ui.Score, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, &model.Factor{}, &model.EventMessage{Effect: reason}); err != nil {
		log.Error("s.addHistory(%+v), err(%v)", ui, err)
		return
	}
	s.updateInfoCache(c, ui)
	s.dao.PubScoreChange(c, mid, &model.ScoreChange{
		Mid:   mid,
		Score: ui.Score,
		TS:    time.Now().Unix(),
	})
	return
}

// UpdateUserScore update user score
func (s *Service) UpdateUserScore(c context.Context, mid int64, ip, effect string) (err error) {
	var (
		tx           *sql.Tx
		factor       = &model.Factor{}
		eventMsg     = &model.EventMessage{Effect: effect}
		preBaseScore int8
	)
	ui, err := s.UserInfo(c, mid, ip)
	if err != nil || ui == nil {
		log.Error("user portrait not found, mid(%d, %v), err(%v)", mid, ui, err)
		return
	}
	preBaseScore = ui.BaseScore
	if ui.BaseScore, err = s.getBaseScore(c, mid, ip); err != nil {
		log.Error("s.calcBaseScore(%d, %s) error(%v)", mid, ip, err)
		return
	}
	s.promBaseScore.Incr("change")
	if ui.BaseScore < preBaseScore {
		s.promBaseScore.Incr("decrease")
	}
	s.resetScore(c, ui, ui.BaseScore, ui.EventScore)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if err = s.dao.TxUpdateBaseScore(c, tx, ui); err != nil {
		log.Error("s.dao.UpdateBaseScore(%v), err(%v)", ui, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, factor, eventMsg); err != nil {
		log.Error("s.addHistory(%v, %v, %v) error(%v)", ui, factor, eventMsg, err)
		return
	}
	s.updateInfoCache(c, ui)
	s.dao.PubScoreChange(c, mid, &model.ScoreChange{
		Mid:   mid,
		Score: ui.Score,
		TS:    time.Now().Unix(),
	})
	return
}

// RefreshBaseScore refresh base score.
func (s *Service) RefreshBaseScore(c context.Context, arg *spy.ArgReset) (err error) {
	var (
		tx           *sql.Tx
		mid          = arg.Mid
		ip           = "ip"
		eventMsg     = &model.EventMessage{Effect: "手动更新基础分值"}
		preBaseScore int8
	)
	ui, err := s.UserInfo(c, mid, "ip")
	if err != nil || ui == nil {
		log.Error("user portrait not found, mid(%d, %v), err(%v)", mid, ui, err)
		return
	}
	preBaseScore = ui.BaseScore
	if ui.BaseScore, err = s.getBaseScore(c, mid, ip); err != nil {
		log.Error("s.calcBaseScore(%d, %s) error(%v)", mid, ip, err)
		return
	}
	s.promBaseScore.Incr("change")
	if ui.BaseScore < preBaseScore {
		s.promBaseScore.Incr("decrease")
	}
	s.resetScore(c, ui, ui.BaseScore, ui.EventScore)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if err = s.dao.TxUpdateBaseScore(c, tx, ui); err != nil {
		log.Error("s.dao.UpdateBaseScore(%v), err(%v)", ui, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, &model.Factor{}, eventMsg); err != nil {
		log.Error("s.addHistory(%v, %v) error(%v)", ui, eventMsg, err)
		return
	}
	s.updateInfoCache(c, ui)
	s.dao.PubScoreChange(c, mid, &model.ScoreChange{
		Mid:   mid,
		Score: ui.Score,
		TS:    time.Now().Unix(),
	})
	return
}

//UpdateBaseScore update base score.
func (s *Service) UpdateBaseScore(c context.Context, arg *spy.ArgReset) (err error) {
	var (
		tx       *sql.Tx
		event    *model.Event
		eventMsg = &model.EventMessage{}
	)
	ui, err := s.UserInfo(c, arg.Mid, "ip")
	if err != nil || ui == nil {
		log.Error("user portrait not fund, mid(%d, %v), err(%v)", arg.Mid, ui, err)
		return
	}
	s.resetScore(c, ui, _score, ui.EventScore)
	eventMsg.Effect = fmt.Sprintf("人工恢复基础得分(%s)", arg.Operator)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if err = s.dao.TxUpdateBaseScore(c, tx, ui); err != nil {
		log.Error("s.dao.UpdateBaseScore(%v), err(%v)", ui, err)
		return
	}
	if event, err = s.dao.Event(c, RestoreBaseScoreEvent); err != nil {
		log.Error("s.dao.Event(%s), err(%v)", RestoreBaseScoreEvent, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, &model.Factor{EventID: event.ID}, eventMsg); err != nil {
		log.Error("s.addHistory(%v), err(%v)", ui, err)
		return
	}
	s.dao.DelInfoCache(c, arg.Mid)
	s.dao.PubScoreChange(c, arg.Mid, &model.ScoreChange{
		Mid:   arg.Mid,
		Score: ui.Score,
		TS:    time.Now().Unix(),
	})
	return
}

//UpdateEventScore update event score.
func (s *Service) UpdateEventScore(c context.Context, arg *spy.ArgReset) (err error) {
	var (
		tx       *sql.Tx
		eventMsg = &model.EventMessage{}
	)
	ui, err := s.UserInfo(c, arg.Mid, "ip")
	if err != nil || ui == nil {
		log.Error("user portrait not fund, mid(%d, %v), err(%v)", arg.Mid, ui, err)
		return
	}
	s.resetScore(c, ui, ui.BaseScore, _score)
	eventMsg.Effect = fmt.Sprintf("人工恢复行为得分(%s)", arg.Operator)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if err = s.dao.TxUpdateEventScore(c, tx, ui.Mid, ui.EventScore, ui.Score); err != nil {
		log.Error("s.dao.UpdateEventScore(%v), err(%v)", ui, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, &model.Factor{}, eventMsg); err != nil {
		log.Error("s.addHistory(%v), err(%v)", ui, err)
		return
	}
	s.dao.DelInfoCache(c, arg.Mid)
	s.dao.PubScoreChange(c, arg.Mid, &model.ScoreChange{
		Mid:   arg.Mid,
		Score: ui.Score,
		TS:    time.Now().Unix(),
	})
	return
}

func (s *Service) txAddHistory(c context.Context, tx *sql.Tx, ui *model.UserInfo, factor *model.Factor, eventMsg *model.EventMessage) (err error) {
	// append event history
	var (
		ueh = &model.UserEventHistory{
			Mid:        ui.Mid,
			EventID:    factor.EventID,
			Score:      ui.Score,
			BaseScore:  ui.BaseScore,
			EventScore: ui.EventScore,
			Reason:     eventMsg.Effect,
			FactorVal:  factor.FactorVal,
		}
		remarkBytes []byte
	)
	if remarkBytes, err = json.Marshal(eventMsg); err != nil {
		log.Error("json.Marshal(%v) error(%v)", eventMsg, err)
		ueh.Remark = "{}"
	} else {
		ueh.Remark = string(remarkBytes)
	}
	if err = s.dao.TxAddEventHistory(c, tx, ueh); err != nil {
		log.Error("s.dao.AddEventHistory(%v) error(%v)", ueh, err)
		return
	}
	return
}

func (s *Service) resetScore(c context.Context, ui *model.UserInfo, baseScore, eventScore int8) (err error) {
	ui.BaseScore = baseScore
	ui.EventScore = eventScore
	ui.Score = int8(int(ui.EventScore) * int(ui.BaseScore) / _score)
	return
}

func (s *Service) updateInfoCache(c context.Context, ui *model.UserInfo) (err error) {
	s.mission(func() {
		if err = s.dao.SetUserInfoCache(context.TODO(), ui); err != nil {
			log.Error("s.dao.SetInfoCache(%v) error(%v)", ui, err)
		}
	})
	return
}

//ClearReliveTimes clear times.
func (s *Service) ClearReliveTimes(c context.Context, arg *spy.ArgReset) (err error) {
	var (
		tx       *sql.Tx
		eventMsg = &model.EventMessage{}
	)
	ui, err := s.UserInfo(c, arg.Mid, "ip")
	if err != nil || ui == nil {
		log.Error("user portrait not fund, mid(%d, %v), err(%v)", arg.Mid, ui, err)
		return
	}
	ui.ReliveTimes = 0
	eventMsg.Effect = fmt.Sprintf("清除封号记次(%s)", arg.Operator)
	if tx, err = s.dao.BeginTran(c); err != nil {
		log.Error("s.dao.BeginTran() err(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()
	if err = s.dao.TxClearReliveTimes(c, tx, ui); err != nil {
		log.Error("s.dao.ClearReliveTimes(%v), err(%v)", ui, err)
		return
	}
	if err = s.txAddHistory(c, tx, ui, &model.Factor{}, eventMsg); err != nil {
		log.Error("s.addHistory(%v), err(%v)", ui, err)
		return
	}
	s.dao.DelInfoCache(c, arg.Mid)
	return
}

// TelRiskLevel tel risk level.
func (s *Service) TelRiskLevel(c context.Context, mid int64, ip string) (riskLevel int8, err error) {
	var (
		tel   *model.TelInfo
		level int8
	)
	riskLevel = dao.TelRiskLevelUnknown
	if tel, err = s.dao.TelInfo(c, mid); err != nil {
		log.Error("s.dao.TelInfo error(%v)", err)
		return
	}
	if len(tel.Tel) == 0 {
		log.Warn("mid(%d) no tel info", mid)
		return
	}
	// white tel
	if telnum, theErr := strconv.ParseInt(tel.Tel, 10, 64); theErr == nil {
		for _, whiteTel := range s.c.Property.White.Tels {
			if telnum >= whiteTel.From && telnum <= whiteTel.To {
				log.Info("spy hit tel white from [%d] to [%d]", whiteTel.From, whiteTel.To)
				riskLevel = dao.TelRiskLevelLow
				return
			}
		}
	} else {
		log.Error("+v", errors.WithStack(theErr))
	}

	args := url.Values{}
	args.Set("accountType", fmt.Sprintf("%d", model.AccountType))
	args.Set("uid", fmt.Sprintf("%d", mid))
	args.Set("phoneNumber", tel.Tel)
	args.Set("registerTime", fmt.Sprintf("%d", tel.JoinTime))
	args.Set("registerIp", tel.JoinIP)
	if level, err = s.dao.RegisterProtection(c, args, ip); err != nil {
		log.Error("s.dao.RegisterProtection error(%v)", err)
		return
	}
	switch level {
	case model.Nomal:
		riskLevel = dao.TelRiskLevelLow
	case model.LevelOne:
		riskLevel = dao.TelRiskLevelMedium
	case model.LevelTwo:
		riskLevel = dao.TelRiskLevelMedium
	case model.LevelThree:
		riskLevel = dao.TelRiskLevelHigh
	case model.LevelFour:
		riskLevel = dao.TelRiskLevelHigh
	default:
		log.Error(" RegisterProtection not found level (%d)", level)
	}
	return
}
