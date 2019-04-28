package service

import (
	"context"
	"go-common/app/admin/main/videoup/model/archive"
	"go-common/app/admin/main/videoup/model/manager"
	"go-common/app/admin/main/videoup/model/monitor"
	"go-common/library/log"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// MonitorResult 获取监控业务的统计结果列表
func (s *Service) MonitorResult(c context.Context, p *monitor.RuleResultParams) (res []*monitor.RuleResultData, err error) {
	var (
		rules []*monitor.Rule
		mUser *manager.User
	)
	res = []*monitor.RuleResultData{}
	if rules, err = s.monitor.GetRules(c, p.Type, p.Business, true); err != nil {
		return
	}
	for _, v := range rules {
		var stats *monitor.Stats
		if stats, err = s.MonitorStats(c, v); err != nil {
			return
		}
		if mUser, err = s.mng.User(c, v.UID); err != nil {
			mUser = &manager.User{
				ID: v.UID,
			}
			err = nil
		}
		tmp := &monitor.RuleResultData{
			Stats: stats,
			User:  mUser,
			Rule:  v,
		}
		res = append(res, tmp)
	}
	return
}

// MonitorStats 根据business和rule获取统计结果
func (s *Service) MonitorStats(c context.Context, rule *monitor.Rule) (stats *monitor.Stats, err error) {
	var (
		qualiKeys []string //符合条件的统计redis key
	)
	qualiKeys, err = s.RuleQualifiedKeys(c, rule)
	//根据符合条件的redis key qualiKeys获取数据
	stats = &monitor.Stats{}
	for _, k := range qualiKeys {
		var res *monitor.Stats
		if res, err = s.monitor.StatsResult(c, k, rule.RuleConf); err != nil {
			return
		}
		stats.TotalCount += res.TotalCount
		stats.MoniCount += res.MoniCount
		if res.MaxTime > stats.MaxTime {
			stats.MaxTime = res.MaxTime
		}
	}
	return
}

// MonitorRuleUpdate 更新监控规则
func (s *Service) MonitorRuleUpdate(c context.Context, rule *monitor.Rule) (err error) {
	rule.CTime = time.Now().Format("2006-01-02 15:04:05")
	err = s.monitor.SetRule(c, rule)
	return
}

// RuleQualifiedKeys 获取监控业务中符合监控条件的Redis key
func (s *Service) RuleQualifiedKeys(c context.Context, rule *monitor.Rule) (qualiKeys []string, err error) {
	var (
		allKeys []string //当前业务的所有统计redis key
		//qualiKeys []string                //符合条件的统计redis key
		prefix  string                  //当前业务redis 可以前缀
		keyConf *monitor.KeyConf        //当前业务redis key 的字段配置信息
		moniCdt = rule.RuleConf.MoniCdt //当前业务中，需要统计的条件
		ok      bool
	)
	qualiKeys = []string{}
	if keyConf, ok = monitor.RedisKeyConf[rule.Business]; !ok {
		err = errors.New("Business Not Exists")
		return
	}
	prefix, allKeys, err = s.monitor.BusStatsKeys(c, rule.Business)
	kFields := keyConf.KFields //Redis key中的字段
	for _, fullK := range allKeys {
		k := strings.Replace(fullK, prefix, "", 1) //去掉redis key 中的前缀monitor_stats_{business}_
		ks := strings.Split(k, "_")                //redis key 的切割，格式：%d_%d
		if len(ks) != len(kFields) {               //KFields中的字段数必须与redis key中数量一致
			err = errors.New("KeyConf 中 KFields 的字段配置错误")
			return
		}
		//把满足条件的Redis key找出来
		qualified := true //当前key是否满足条件
		i := -1
		for kf := range kFields {
			i++
			var (
				fv int64
			)
			fv, err = strconv.ParseInt(ks[i], 10, 64) //Redis key 中第i位的值
			if err != nil {
				return
			}
			if _, ok = moniCdt[kf]; !ok {
				err = errors.New("配置的 moniCdt 中不存在: " + kf)
				return
			}
			switch moniCdt[kf].Comp {
			case monitor.CompE:
				if fv != moniCdt[kf].Value {
					qualified = false
				}
			case monitor.CompNE:
				if fv == moniCdt[kf].Value {
					qualified = false
				}
			case monitor.CompGET:
				if fv < moniCdt[kf].Value {
					qualified = false
				}
			case monitor.CompLET:
				if fv > moniCdt[kf].Value {
					qualified = false
				}
			case monitor.CompGT:
				if fv <= moniCdt[kf].Value {
					qualified = false
				}
			case monitor.CompLT:
				if fv >= moniCdt[kf].Value {
					qualified = false
				}
			default:
				err = errors.New("配置的 MoniCdt 中 comparison 不合法: " + moniCdt[kf].Comp)
				return
			}
		}
		if qualified {
			qualiKeys = append(qualiKeys, fullK)
		}
	}
	return
}

// MoniStayOids 获取监控范围内，滞留的oids
func (s *Service) MoniStayOids(c context.Context, tp, bid int8, id int64) (total int, oidMap map[int64]int, qualiKeys []string, err error) {
	var (
		rule *monitor.Rule
	)
	if rule, err = s.monitor.GetRule(c, tp, bid, id); err != nil {
		return
	}
	//查找符合条件的统计redis key
	if qualiKeys, err = s.RuleQualifiedKeys(c, rule); err != nil {
		return
	}
	oidMap, total, err = s.monitor.StayOids(c, rule, qualiKeys)
	log.Info("MoniStayOids(%d,%d,%d) oidMap(%v)", tp, bid, id, oidMap)
	return
}

// MonitorStayOids 获取监控范围内，滞留的oids
func (s *Service) MonitorStayOids(c context.Context, id int64) (oidMap map[int64]int, err error) {
	return s.data.MonitorOids(c, id)
}

// MonitorNotifyResult 获取达到了报警阀值的数据
func (s *Service) MonitorNotifyResult(c context.Context) (res []*monitor.RuleResultData, err error) {
	var (
		rules []*monitor.Rule
		stats *monitor.Stats
	)
	res = []*monitor.RuleResultData{}
	if rules, err = s.monitor.GetAllRules(c, false); err != nil {
		log.Error("MonitorNotifyCheck() error(%v)", err)
		return
	}
	for _, v := range rules {
		if v.Business == monitor.BusVideo {
			s.MonitorCheckVideoStatus(c, v.Type, v.ID)
		}
		if stats, err = s.MonitorStats(c, v); err != nil {
			log.Error("MonitorNotifyCheck() error(%v)", err)
			err = nil
			continue
		}
		notify := true
		//暂时只有time、count这两个报警条件
		if _, ok := v.RuleConf.NotifyCdt["time"]; ok {
			threshold := v.RuleConf.NotifyCdt["time"].Value
			comp := v.RuleConf.NotifyCdt["time"].Comp
			switch comp {
			case monitor.CompGT:
				if int64(stats.MaxTime) < threshold {
					notify = false
				}
			case monitor.CompLT:
				if int64(stats.MaxTime) > threshold {
					notify = false
				}
			}
		}
		if _, ok := v.RuleConf.NotifyCdt["count"]; ok {
			threshold := v.RuleConf.NotifyCdt["count"].Value
			comp := v.RuleConf.NotifyCdt["count"].Comp
			switch comp {
			case monitor.CompGT:
				if int64(stats.MoniCount) < threshold {
					notify = false
				}
			case monitor.CompLT:
				if int64(stats.MoniCount) > threshold {
					notify = false
				}
			}
		}
		if notify {
			tmp := &monitor.RuleResultData{
				Stats: stats,
				Rule:  v,
			}
			res = append(res, tmp)
		}
	}
	return
}

// MonitorCheckVideoStatus 检查视频的稿件状态，如果是-100则剔除SortedSet的数据
func (s *Service) MonitorCheckVideoStatus(c context.Context, tp int8, id int64) (err error) {
	var (
		vidMap     map[int64]int
		vidAidMap  map[int64]int64
		arcStates  map[int64]int
		vids, aids []int64
		bid        = monitor.BusVideo
		keys       []string
	)
	if _, vidMap, keys, err = s.MoniStayOids(c, tp, bid, id); err != nil {
		log.Error("s.MoniStayOids(%d,%d,%d) error(%v)", tp, bid, id, err)
		return
	}
	for vid := range vidMap {
		vids = append(vids, vid)
	}
	if vidAidMap, err = s.arc.VideoAidMap(c, vids); err != nil {
		log.Error("s.VideoAidMap(%d) error(%v)", vids, err)
		return
	}
	for _, aid := range vidAidMap {
		aids = append(aids, aid)
	}
	if arcStates, err = s.arc.ArcStateMap(c, aids); err != nil {
		log.Error("s.ArcStateMap(%d) error(%v)", aids, err)
		return
	}
	for vid, aid := range vidAidMap {
		if _, ok := arcStates[aid]; !ok {
			continue
		}
		if arcStates[aid] != int(archive.StateForbidUpDelete) && arcStates[aid] != int(archive.StateForbidLock) && arcStates[aid] != int(archive.StateForbidRecycle) {
			continue
		}
		for _, k := range keys {
			if err = s.monitor.RemMonitorStats(c, k, vid); err != nil {
				log.Error("s.monitor.RemMonitorStats(%s,%d) error(%v)", k, vid, err)
			}
		}
	}
	return
}
