package service

import (
	"context"
	"errors"
	"go-common/app/admin/main/aegis/model/monitor"
	"go-common/library/log"
	"time"
)

// MonitorBuzResult 获取业务的监控结果
func (s *Service) MonitorBuzResult(c context.Context, bid int64) (res []*monitor.RuleResultData, err error) {
	var (
		rules    []*monitor.Rule
		uids     []int64
		uNames   map[int64]string
		statsMap map[int64]*monitor.Stats
		min, max int64
	)
	statsMap = make(map[int64]*monitor.Stats)
	if rules, err = s.mysql.MoniBizRules(c, bid); err != nil {
		log.Error("s.MonitorResult(%d) error:%v", bid, err)
		return
	}
	for _, rule := range rules {
		uids = append(uids, rule.UID)

		if min, max, err = s.monitorNotifyTime(rule.RuleConf); err != nil {
			log.Error("s.MonitorBuzResult(%d) s.monitorNotifyTime(%+v) error:%v", bid, rule.RuleConf, err)
			continue
		}
		if statsMap[rule.ID], err = s.redis.MoniRuleStats(c, rule.ID, min, max); err != nil {
			log.Error("s.redis.MoniRuleStats(%d,%+v) error:%v", rule.ID, rule.RuleConf, err)
			err = nil
			statsMap[rule.ID] = &monitor.Stats{}
		}
	}
	if uNames, err = s.http.GetUnames(c, uids); err != nil {
		log.Error("s.MonitorResult(%d) s.http.ManagerUNames(%v) error:%v", bid, uids, err)
		err = nil
	}

	for _, rule := range rules {
		var (
			uName string
			stats *monitor.Stats
		)
		if _, ok := uNames[rule.UID]; ok {
			uName = uNames[rule.UID]
		}
		if _, ok := statsMap[rule.ID]; ok {
			stats = statsMap[rule.ID]
		}
		data := &monitor.RuleResultData{
			Rule: rule,
			User: &monitor.User{
				ID:       rule.UID,
				UserName: uName,
				NickName: uName,
			},
			Stats: stats,
		}
		res = append(res, data)
	}
	return
}

// MonitorResultOids 获取
func (s *Service) MonitorResultOids(c context.Context, rid int64) (res map[int64]int, err error) {
	var (
		min, max int64
		rule     *monitor.Rule
	)
	if rule, err = s.mysql.MoniRule(c, rid); err != nil {
		log.Error("s.MonitorResultOids(%d) error:%v", rid, err)
		return
	}
	if min, max, err = s.monitorNotifyTime(rule.RuleConf); err != nil {
		log.Error("s.MonitorResultOids(%d) s.monitorNotifyTime() error:%v", rid, err)
		return
	}
	return s.redis.MoniRuleOids(c, rid, min, max)
}

// monitorNotifyTime 计算监控报警的score区间
func (s *Service) monitorNotifyTime(conf *monitor.RuleConf) (tFrom, tTo int64, err error) {
	now := time.Now().Unix()
	if _, ok := conf.NotifyCdt["time"]; !ok {
		err = errors.New("配置的 NotifyCdt 中不存在 time")
		return
	}
	timeCdt := conf.NotifyCdt["time"].Value
	compCdt := conf.NotifyCdt["time"].Comp
	switch compCdt {
	case monitor.CompGT:
		tFrom = 0
		tTo = now - timeCdt
	case monitor.CompLT:
		tFrom = now - timeCdt
		tTo = now
	default:
		err = errors.New("配置的 NotifyCdt 中 comparison 不合法: " + compCdt)
		return
	}
	return
}
