package service

import (
	"context"

	"go-common/app/service/main/spy/conf"
	"go-common/app/service/main/spy/dao"
	"go-common/app/service/main/spy/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	// RestoreBaseScoreEvent represents a event triggered when user resets base score.
	RestoreBaseScoreEvent = "restore_base_score"
)

const (
	_defaultService   = ""
	_defaultRiskLevel = 1

	_telIsBound       = "_telIsBound"
	_telIsNotBound    = "_telIsNotBound"
	_telIsLowRisk     = "_telIsLowRisk"
	_telIsMediumRisk  = "_telIsMediumRisk"
	_telIsHighRisk    = "_telIsHighRisk"
	_telIsUnknownRisk = "_telIsUnknownRisk"
	_mailIsBound      = "_mailIsBound"
	_mailIsNotBound   = "_mailIsNotBound"
	_idenIsAuthed     = "_idenIsAuthed"
	_idenIsNotAuthed  = "_idenIsNotAuthed"
)

const (
	// IterFail represents iteration failure.
	IterFail IterState = iota
	// IterSuccess represents iteration success.
	IterSuccess
)

// Rule represents a rule.
type Rule = string

// Rules represents an array of rules.
type Rules = []Rule

// IterState represents the state of iteration.
type IterState = int

// RuleFunc represents a judgement function.
type RuleFunc func(info *JudgementInfo) (result bool, err error)

var ruleToFuncMaps = map[string]RuleFunc{
	_telIsBound:       telIsBound,
	_telIsNotBound:    telIsNotBound,
	_telIsLowRisk:     telIsLowRisk,
	_telIsMediumRisk:  telIsMediumRisk,
	_telIsHighRisk:    telIsHighRisk,
	_telIsUnknownRisk: telIsUnknownRisk,
	_mailIsBound:      mailIsBound,
	_mailIsNotBound:   mailIsNotBound,
	_idenIsAuthed:     idenIsAuthed,
	_idenIsNotAuthed:  idenIsNotAuthed,
}

// FactoryMeta includes details to find a factor.
type FactoryMeta struct {
	serviceName string
	eventName   string
	riskLevel   int8
}

// JudgementInfo includes all needed information to judge factor.
type JudgementInfo struct {
	mid         int64
	ip          string
	ctx         context.Context
	s           *Service
	auditInfo   *model.AuditInfo
	profileInfo *model.ProfileInfo
	telRiskInfo *model.TelRiskInfo
}

// RuleMap represents a map-like struct which contains rules and factoryMeta.
type RuleMap struct {
	rules       Rules
	factoryMeta FactoryMeta
}

// NewJudgementInfo returns a new judgement info.
func NewJudgementInfo(c context.Context, s *Service, mid int64, ip string) *JudgementInfo {
	return &JudgementInfo{
		s:   s,
		mid: mid,
		ctx: c,
		ip:  ip,
	}
}

func (ji *JudgementInfo) getAuditInfo() (auditInfo *model.AuditInfo, err error) {
	if ji.auditInfo != nil {
		return ji.auditInfo, nil
	}
	auditInfo, err = ji.s.dao.AuditInfo(ji.ctx, ji.mid, ji.ip)
	if err != nil {
		log.Error("s.dao.AuditInfo(%d, %s) error(%v)", ji.mid, ji.ip, err)
		return
	}
	ji.auditInfo = auditInfo
	return
}

func (ji *JudgementInfo) getProfileInfo() (profileInfo *model.ProfileInfo, err error) {
	if ji.profileInfo != nil {
		return ji.profileInfo, nil
	}
	profileInfo, err = ji.s.dao.ProfileInfo(ji.ctx, ji.mid, ji.ip)
	if err != nil {
		log.Error("s.dao.ProfileInfo(%d, %s) error(%v)", ji.mid, ji.ip, err)
		return
	}
	ji.profileInfo = profileInfo
	return ji.profileInfo, nil
}

func (ji *JudgementInfo) getTelRiskInfo() (telRiskInfo *model.TelRiskInfo, err error) {
	var (
		telRiskLevel    int8
		event           *model.Event
		eventHistory    *model.UserEventHistory
		unicomGiftState int
	)
	if ji.telRiskInfo != nil {
		return ji.telRiskInfo, nil
	}
	telRiskLevel, err = ji.s.TelRiskLevel(ji.ctx, ji.mid, "")
	if err != nil {
		log.Error("s.dao.TelRiskLevel(%d) error:(%v)", ji.mid, err)
		return
	}
	if event, err = ji.s.dao.Event(ji.ctx, RestoreBaseScoreEvent); err != nil {
		return
	}
	if eventHistory, err = ji.s.dao.EventHistoryByMidAndEvent(ji.ctx, ji.mid, event.ID); err != nil {
		return
	}
	if unicomGiftState, err = ji.s.dao.UnicomGiftState(ji.ctx, ji.mid); err != nil {
		return
	}

	ji.telRiskInfo = &model.TelRiskInfo{
		TelRiskLevel:    telRiskLevel,
		RestoreHistory:  eventHistory,
		UnicomGiftState: unicomGiftState,
	}
	return ji.telRiskInfo, nil
}

func (s *Service) getRulesMap() []RuleMap {
	return []RuleMap{
		// 绑定手机状态：绑定手机，且天御认定为低风险	邮箱绑定状态：any	实名认证状态：通过认证无论人工还是芝麻
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelLowRiskAndIdenAuth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsBound, _telIsLowRisk, _idenIsAuthed,
			}},
		// 绑定手机状态：绑定手机，且天域认定为低风险	邮箱绑定状态：any	实名认证状态：无
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelLowRiskAndIdenUnauth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsBound, _telIsLowRisk, _idenIsNotAuthed,
			}},
		// 绑定手机状态：绑定手机，且天域认定为未知风险	邮箱绑定状态：any	实名认证状态：通过认证无论人工还是芝麻
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelUnknownRiskAndIdenAuth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsBound, _telIsUnknownRisk, _idenIsAuthed,
			}},
		// 绑定手机状态：绑定手机，且天域认定为中风险	邮箱绑定状态：any	实名认证状态：通过认证无论人工还是芝麻
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelMediumRiskAndIdenAuth,
				riskLevel:   _defaultRiskLevel,
			},

			rules: Rules{
				_telIsBound, _telIsMediumRisk, _idenIsAuthed,
			}},
		// 绑定手机状态：绑定手机，且天域认定为未知风险	邮箱绑定状态：any	实名认证状态：无
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelUnknownRiskAndIdenUnauth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsBound, _telIsUnknownRisk, _idenIsNotAuthed,
			}},
		// 绑定手机状态：绑定手机，且天域认定为中风险	邮箱绑定状态：any	实名认证状态：无
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelMediumRiskAndIdenUnauth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsBound, _telIsMediumRisk, _idenIsNotAuthed,
			}},
		// 绑定手机状态：无	绑定邮箱状态：有	实名认证状态：any
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindMailAndIdenUnknown,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsNotBound, _mailIsBound,
			}},
		// 绑定手机状态：绑定手机，且天域认定为高风险	邮箱绑定状态：any	实名认证状态：实名认证，无论人工还是芝麻
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelHighRiskAndIdenAuth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsBound, _telIsHighRisk, _idenIsAuthed,
			}},
		// 无绑定手机状态：绑定手机	绑定邮箱状态：无	实名认证状态：无
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindNothingV2,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsNotBound, _mailIsNotBound, _idenIsNotAuthed,
			}},
		// 无绑定手机状态：绑定手机	绑定邮箱状态：无	实名认证状态：有
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindNothingAndIdenAuth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsNotBound, _mailIsNotBound, _idenIsAuthed,
			}},
		// 绑定手机状态：绑定手机，且天域认定为高风险	邮箱绑定状态：any	实名认证状态：无
		{
			factoryMeta: FactoryMeta{
				serviceName: _defaultService,
				eventName:   conf.Conf.Property.Event.BindTelHighRiskAndIdenUnauth,
				riskLevel:   _defaultRiskLevel,
			},
			rules: Rules{
				_telIsBound, _telIsHighRisk, _idenIsNotAuthed,
			}},
	}
}

func (s *Service) getJudgementInfo(c context.Context, mid int64, ip string) (judgementInfo *JudgementInfo, err error) {
	return NewJudgementInfo(c, s, mid, ip), nil
}

func (s *Service) getRuleFunc(rule Rule) (ruleFunc RuleFunc, err error) {
	ruleFunc, ok := ruleToFuncMaps[rule]
	if !ok {
		err = ecode.SpyRuleNotExist
		return
	}
	return
}

func (s *Service) iterRules(info *JudgementInfo, rules Rules) (state IterState, err error) {
	var (
		counter  int
		ruleFunc RuleFunc
		result   bool
	)
	for _, rule := range rules {
		if ruleFunc, err = s.getRuleFunc(rule); err != nil {
			return
		}
		if result, err = ruleFunc(info); err != nil {
			return
		}
		if !result {
			log.Info("s.iterRules mid(%d) rules(%v) counter(%d) state(%d)", info.mid, rules, counter, IterFail)
			return IterFail, nil
		}
		counter++
	}
	log.Info("s.iterRules mid(%d) rules(%v) counter(%d) state(%d)", info.mid, rules, counter, IterSuccess)
	return IterSuccess, nil
}

func (s *Service) getBaseScoreFactor(info *JudgementInfo) (meta FactoryMeta, err error) {
	var (
		ruleMap RuleMap
		state   IterState
	)
	for _, ruleMap = range s.getRulesMap() {
		if state, err = s.iterRules(info, ruleMap.rules); err != nil {
			log.Error("s.iterRules(%v, %v) err(%v)", info, ruleMap.rules, err)
			return
		}
		if state == IterSuccess {
			return ruleMap.factoryMeta, nil
		}
	}
	err = ecode.SpyRulesNotMatch
	return
}

func telIsBound(info *JudgementInfo) (result bool, err error) {
	var (
		auditInfo *model.AuditInfo
	)
	if auditInfo, err = info.getAuditInfo(); err != nil {
		return
	}
	return auditInfo.BindTel, nil
}

func telIsNotBound(info *JudgementInfo) (result bool, err error) {
	var (
		auditInfo *model.AuditInfo
	)
	if auditInfo, err = info.getAuditInfo(); err != nil {
		return
	}
	return !auditInfo.BindTel, nil
}

func telIsLowRisk(info *JudgementInfo) (result bool, err error) {
	var (
		telRiskInfo *model.TelRiskInfo
	)
	if telRiskInfo, err = info.getTelRiskInfo(); err != nil {
		return
	}
	if telRiskInfo.RestoreHistory != nil || telRiskInfo.UnicomGiftState == 1 {
		return true, nil
	}
	return telRiskInfo.TelRiskLevel == dao.TelRiskLevelLow, nil
}

func telIsMediumRisk(info *JudgementInfo) (result bool, err error) {
	var (
		telRiskInfo *model.TelRiskInfo
	)
	if telRiskInfo, err = info.getTelRiskInfo(); err != nil {
		return
	}
	if telRiskInfo.RestoreHistory != nil || telRiskInfo.UnicomGiftState == 1 {
		return false, nil
	}
	return telRiskInfo.TelRiskLevel == dao.TelRiskLevelMedium, nil
}

func telIsHighRisk(info *JudgementInfo) (result bool, err error) {
	var (
		telRiskInfo *model.TelRiskInfo
	)
	if telRiskInfo, err = info.getTelRiskInfo(); err != nil {
		return
	}
	if telRiskInfo.RestoreHistory != nil || telRiskInfo.UnicomGiftState == 1 {
		return false, nil
	}
	return telRiskInfo.TelRiskLevel == dao.TelRiskLevelHigh, nil
}

func telIsUnknownRisk(info *JudgementInfo) (result bool, err error) {
	var (
		telRiskInfo *model.TelRiskInfo
	)
	if telRiskInfo, err = info.getTelRiskInfo(); err != nil {
		return
	}
	if telRiskInfo.RestoreHistory != nil || telRiskInfo.UnicomGiftState == 1 {
		return false, nil
	}
	return telRiskInfo.TelRiskLevel == dao.TelRiskLevelUnknown, nil
}

func mailIsBound(info *JudgementInfo) (result bool, err error) {
	var (
		auditInfo *model.AuditInfo
	)
	if auditInfo, err = info.getAuditInfo(); err != nil {
		return
	}
	return auditInfo.BindMail, nil
}

func mailIsNotBound(info *JudgementInfo) (result bool, err error) {
	var (
		auditInfo *model.AuditInfo
	)
	if auditInfo, err = info.getAuditInfo(); err != nil {
		return
	}
	return !auditInfo.BindMail, nil
}

func idenIsNotAuthed(info *JudgementInfo) (result bool, err error) {
	var (
		profileInfo *model.ProfileInfo
	)
	if profileInfo, err = info.getProfileInfo(); err != nil {
		return
	}
	return profileInfo.Identification == 0, nil
}

func idenIsAuthed(info *JudgementInfo) (result bool, err error) {
	var (
		profileInfo *model.ProfileInfo
	)
	if profileInfo, err = info.getProfileInfo(); err != nil {
		return
	}
	return profileInfo.Identification == 1, nil
}
