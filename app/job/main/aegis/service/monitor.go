package service

import (
	"context"
	"errors"
	"fmt"
	moniMdl "go-common/app/job/main/aegis/model/monitor"
	accApi "go-common/app/service/main/account/api"
	upApi "go-common/app/service/main/up/api/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// monitorArchive 稿件业务监控
// 注意：oa有可能是nil，使用前必须判断！！！
func (s *Service) monitorArchive(act string, oa, na *moniMdl.BinlogArchive) (errs []error) {
	var (
		c      = context.TODO()
		logs   []string
		err    error
		errs2  []error
		nAddit *moniMdl.ArchiveAddit
	)
	defer func() {
		logStr := strings.Join(logs, "\n")
		if x := recover(); x != nil {
			log.Error("s.monitorArchive() unknown panic(%v)", x)
		} else if len(errs) > 0 {
			log.Error("s.monitorArchive(\n act: %s \n oa: %+v \n na: %+v) \n logStr:\n %v \n error:%v", act, oa, na, logStr, errs)
		} else {
			log.Info("s.monitorArchive(\n act: %s \n oa: %+v \n na: %+v) \n logStr:\n %v", act, oa, na, logStr)
		}
	}()
	if (na.Attr>>moniMdl.ArchiveBitPGC)&int64(1) == 1 {
		logs = append(logs, "忽略PGC稿件")
		return
	}
	if na == nil {
		err = errors.New("new msg nil")
		errs = append(errs, err)
		logs = append(logs, "databus数据异常，new msg nil")
		return
	}
	na.IsSpecTID = moniMdl.SpecialTypeIDs[na.TypeID]
	if nAddit, err = s.moniDao.ArchiveAttr(c, na.ID); err != nil {
		logs = append(logs, fmt.Sprintf("warn:稿件Addit获取失败！aid:%d error:%v", na.ID, err))
		if err != ecode.NothingFound {
			errs = append(errs, err)
			return
		}
		err = nil
	} else {
		na.Addit = nAddit
	}
	errs2 = s.monitorHandle(moniMdl.BusArc, *na, na.ID)
	errs = append(errs, errs2...)
	return
}

// monitorUpDelArc 监控UP主删除稿件。监控指定的UP主。1、在高能联盟的up主特殊用户组；2、粉丝数超过50w
func (s *Service) monitorUpDelArc(id int64, obj interface{}) (satisfy bool, logs []string, err error) {
	var (
		c       = context.TODO()
		pReply  *accApi.ProfileStatReply
		upReply *upApi.HighAllyUpsReply
		a       *moniMdl.BinlogArchive
	)
	logs = append(logs, "s.monitorUpDelArc() begin")
	if obj == nil {
		logs = append(logs, "\t obj是nil")
		err = errors.New("obj is nil")
		return
	}
	switch obj.(type) {
	case *moniMdl.BinlogArchive:
		a = obj.(*moniMdl.BinlogArchive)
	case moniMdl.BinlogArchive:
		ac := obj.(moniMdl.BinlogArchive)
		a = &ac
	default:
		logs = append(logs, fmt.Sprintf("\t 未知类型:%+v", obj))
		err = errors.New("unknown interface type")
		return
	}
	logs = append(logs, fmt.Sprintf("\t archive:%+v", a))
	if a == nil {
		err = errors.New("archive is nil")
		return
	}
	if a.State != moniMdl.ArchiveStateDel {
		logs = append(logs, "\t 非删除，忽略")
		return
	}
	if a.Copyright != moniMdl.ArchiveOriginal {
		logs = append(logs, "\t 非自制，忽略")
		return
	}
	if err = s.moniDao.AddToDelArc(c, a); err != nil {
		logs = append(logs, fmt.Sprintf("\t 添加删除信息到redis失败。error:%v", err))
		err = nil
	}
	if id == moniMdl.RuleHighUpDelArc {
		if upReply, err = s.up.GetHighAllyUps(c, &upApi.HighAllyUpsReq{Mids: []int64{a.MID}}); err != nil {
			logs = append(logs, fmt.Sprintf("\t 获取UP主高能信息失败。error:%v", err))
			log.Error("\t s.monitorUpDelArc() s.up.GetHighAllyUps() error:%v", err)
		}
		logs = append(logs, fmt.Sprintf("\t 用户信息：%+v", upReply))
		if upReply != nil {
			if _, ok := upReply.Lists[a.MID]; ok {
				logs = append(logs, "\t UP主属于高能联盟")
				satisfy = true
				return
			}
		}
	} else if id == moniMdl.RuleFamUpDelArc {
		if pReply, err = s.acc.ProfileWithStat3(c, &accApi.MidReq{Mid: a.MID}); err != nil {
			logs = append(logs, fmt.Sprintf("\t 获取UP主信息失败。error:%v", err))
			log.Error("\t s.monitorUpDelArc() s.acc.ProfileWithStat3() error:%v", err)
		}
		logs = append(logs, fmt.Sprintf("\t 用户信息：%+v", pReply))
		if pReply != nil && pReply.Follower >= 500000 {
			logs = append(logs, "\t UP主属于大UP主")
			satisfy = true
			return
		}
	}
	return
}

// monitorVideo 视频监控
func (s *Service) monitorVideo(act string, ov, nv *moniMdl.BinlogVideo) (errs []error) {
	errs = s.monitorHandle(moniMdl.BusVideo, nv, nv.ID)
	return
}

// monitorHandle 处理监控数据
func (s *Service) monitorHandle(bid int64, nObj interface{}, oid int64) (errs []error) {
	var (
		c            = context.TODO()
		rules        []*moniMdl.Rule
		logs, logs2  []string
		err          error
		errs2        []error
		oKeys, nKeys []string
	)
	defer func() {
		logStr := strings.Join(logs, "\n")
		if x := recover(); x != nil {
			log.Error("s.monitorHandle() unknown panic(%v)", x)
		} else if len(errs) > 0 {
			log.Error("s.monitorHandle(\n na: %+v) \n logStr:\n %v \n error:%v", nObj, logStr, errs)
		} else {
			log.Info("s.monitorHandle(\n na: %+v) \n logStr:\n %v", nObj, logStr)
		}
	}()
	if nObj == nil {
		err = errors.New("new msg nil")
		errs = append(errs, err)
		logs = append(logs, "databus数据异常，new msg nil")
		return
	}
	if rules, err = s.moniDao.RulesByBid(c, bid); err != nil {
		logs = append(logs, "获取监控配置失败")
		return
	}
	if len(rules) == 0 {
		logs = append(logs, "监控配置不存在")
		return
	}
	for _, rule := range rules {
		var allSatisfy = true
		//如果是监控UP主大量删稿，则特殊处理
		if rule.ID == moniMdl.RuleHighUpDelArc || rule.ID == moniMdl.RuleFamUpDelArc {
			if allSatisfy, logs2, err = s.monitorUpDelArc(rule.ID, nObj); err != nil {
				errs = append(errs, err)
			}
			logs = append(logs, logs2...)
		} else {
			for field, cdt := range rule.RuleConf.MoniCdt {
				var (
					val     int64
					satisfy bool
				)
				if val, err = s.reflectIntVal(nObj, field, 0); err != nil {
					errs = append(errs, err)
					logs = append(logs, fmt.Sprintf("没有找到字段%s", field))
				}
				if satisfy, err = s.monitorCompSatisfy(cdt.Comp, val); err != nil {
					allSatisfy = false
					break
				}
				if !satisfy {
					allSatisfy = false
					break
				}
			}
		}
		if allSatisfy { //如果满足所有条件，则移入监控
			nKeys = append(nKeys, fmt.Sprintf(moniMdl.RedisPrefix, rule.ID))
		} else { //如果有条件不满足，则移出监控
			oKeys = append(oKeys, fmt.Sprintf(moniMdl.RedisPrefix, rule.ID))
		}
	}

	logs = append(logs, fmt.Sprintf("%d移出keys：%v", oid, oKeys))
	logs = append(logs, fmt.Sprintf("%d移入keys：%v", oid, nKeys))
	logs2, errs2 = s.monitorSave(oKeys, nKeys, oid)
	logs = append(logs, logs2...)
	if len(errs2) != 0 {
		errs = append(errs, errs2...)
	}
	return
}

// monitorSave 保存结果
func (s *Service) monitorSave(oKeys, nKeys []string, oid int64) (logs []string, errs []error) {
	var (
		c     = context.TODO()
		logs2 []string
		err   error
	)
	defer func() {
		logStr := strings.Join(logs, "\n")
		if x := recover(); x != nil {
			log.Error("s.monitorSave() unknown panic(%v)", x)
		} else if len(errs) != 0 {
			log.Error("s.monitorSave(\n oKeys: %v \n nKeys: %v \n oid: %d) \n logStr:\n %v \n error:%v", oKeys, nKeys, oid, logStr, errs)
		} else {
			log.Info("s.monitorSave(\n oKeys: %v \n nKeys: %v \n oid: %d) \n logStr:\n %v", oKeys, nKeys, oid, logStr)
		}
	}()
	//从旧key中移出
	logs2, err = s.moniDao.RemFromSet(c, oKeys, oid)
	logs = append(logs, logs2...)
	if err != nil {
		errs = append(errs, err)
	}
	//清除过期的旧key
	logs2, err = s.moniDao.ClearExpireSet(c, oKeys)
	logs = append(logs, logs2...)
	if err != nil {
		errs = append(errs, err)
	}
	//移入到新key
	logs2, err = s.moniDao.AddToSet(c, nKeys, oid)
	logs = append(logs, logs2...)
	if err != nil {
		errs = append(errs, err)
	}
	//清除过期的移入key
	logs2, err = s.moniDao.ClearExpireSet(c, nKeys)
	logs = append(logs, logs2...)
	if err != nil {
		errs = append(errs, err)
	}

	return
}

// monitorNotify 监控报警
func (s *Service) monitorNotify() {
	defer func() {
		log.Warn("monitorNotify exited.")
	}()
	var (
		c        = context.TODO()
		err      error
		rules    []*moniMdl.Rule
		stats    *moniMdl.Stats
		min, max int64
	)
	for {
		log.Info("s.monitorNotify() begin")
		if rules, err = s.moniDao.ValidRules(c); err != nil {
			log.Error("s.monitorNotify() rules:%+v error:%v", rules, err)
			time.Sleep(1 * time.Minute)
			continue
		}
		for _, rule := range rules {
			if rule.ID == moniMdl.RuleHighUpDelArc || rule.ID == moniMdl.RuleFamUpDelArc {
				// 删稿监控特殊处理
				s.wg.Add(1)
				go s.moniUpDelArcNotify(rule)
				continue
			}
			if min, max, err = s.monitorNotifyTime(rule.RuleConf); err != nil {
				log.Error("s.monitorNotify() s.monitorNotifyTime() rule:%+v error:%v", rule, err)
				continue
			}
			if stats, err = s.moniDao.MoniRuleStats(c, rule.ID, min, max); err != nil {
				log.Error("s.monitorNotify() s.moniDao.MoniRuleStats(%d,%d,%d) error:%v", rule.ID, min, max, err)
				continue
			}
			notify := s.moniSatisfyNotify(rule.RuleConf, stats)
			if notify {
				title := fmt.Sprintf("%s监控(aegis)", rule.RuleConf.Name)
				body := fmt.Sprintf("当前滞留时间为%s超过阀值，滞留量为%d，整体量为%d \n报警时间：%s", secondsFormat(stats.MaxTime), stats.MoniCount, stats.TotalCount, time.Now().Format("2006-01-02 15:04:05"))
				url := ""
				switch rule.BID {
				case moniMdl.BusVideo:
					url = fmt.Sprintf("http://manager.bilibili.co/#!/video/list?monitor_list=%d_%d_%d", rule.Type, rule.BID, rule.ID)
				case moniMdl.BusArc:
					url = fmt.Sprintf("http://manager.bilibili.co/#!/archive_utils/all?monitor_list=%d_%d_%d", rule.Type, rule.BID, rule.ID)
				}
				body += fmt.Sprintf("\n跳转链接：<a href=\"%s\">点击跳转</a> %s", url, url)
				if err = s.monitorSendNotify(c, rule.RuleConf.Notify.Way, rule.RuleConf.Notify.Member, title, body); err != nil {
					log.Error("s.monitorNotify() s.monitorSendNotify(%d,%v,%s,%s) error:%v", rule.RuleConf.Notify.Way, rule.RuleConf.Notify.Member, title, body, err)
				}
			}
		}
		time.Sleep(30 * time.Minute)
	}
}

// moniSatisfyNotify 检查监控是否满足报警，目前只有数量+时长的条件
func (s *Service) moniSatisfyNotify(conf *moniMdl.RuleConf, stats *moniMdl.Stats) (notify bool) {
	notify = true
	if _, ok := conf.MoniCdt["time"]; ok {
		threshold := conf.NotifyCdt["time"].Value
		comp := conf.NotifyCdt["time"].Comp
		switch comp {
		case moniMdl.CompGT:
			if int64(stats.MaxTime) < threshold {
				notify = false
			}
		case moniMdl.CompLT:
			if int64(stats.MaxTime) > threshold {
				notify = false
			}
		}
	}
	if _, ok := conf.NotifyCdt["count"]; ok {
		threshold := conf.NotifyCdt["count"].Value
		comp := conf.NotifyCdt["count"].Comp
		switch comp {
		case moniMdl.CompGT:
			if int64(stats.MoniCount) < threshold {
				notify = false
			}
		case moniMdl.CompLT:
			if int64(stats.MoniCount) > threshold {
				notify = false
			}
		}
	}
	return
}

// moniUpDelArcNotify 特殊处理UP主删稿的逻辑
func (s *Service) moniUpDelArcNotify(rule *moniMdl.Rule) (err error) {
	var (
		c          = context.TODO()
		min, max   int64
		oidMap     map[int64]int
		oids, mids []int64
		infos      map[int64]*moniMdl.DelArcInfo
		delMap     map[int64][]*moniMdl.DelArcInfo
		accStats   map[int64]*accApi.ProfileStatReply
		threshold  int
	)
	defer func() {
		s.wg.Done()
	}()
	if _, ok := rule.RuleConf.NotifyCdt["count"]; !ok {
		err = errors.New("notify count config error")
		log.Error("s.moniUpDelArcNotify(%+v) 没有count监控配置", rule)
		return
	}
	threshold = int(rule.RuleConf.NotifyCdt["count"].Value)
	delMap = make(map[int64][]*moniMdl.DelArcInfo)
	min, max, err = s.monitorNotifyTime(rule.RuleConf)
	if err != nil {
		log.Error("s.monitorNotifyTime() rule:%+v error:%v", rule, err)
		return
	}
	if oidMap, err = s.moniDao.MoniRuleOids(c, rule.ID, min, max); err != nil {
		log.Error("s.moniDao.MoniRuleOids() rule:%+v error:%v", rule, err)
		return
	}
	for oid := range oidMap {
		oids = append(oids, oid)
	}
	if infos, err = s.moniDao.ArcDelInfos(c, oids); err != nil {
		log.Error("s.moniUpDelArcNotify() s.moniDao.ArcDelInfos(%v) error(%v)", oids, err)
		return
	}
	for _, info := range infos {
		delMap[info.MID] = append(delMap[info.MID], info)
		mids = append(mids, info.MID)
	}
	if accStats, err = s.multiAccounts(c, mids); err != nil {
		log.Error("s.moniUpDelArcNotify() s.multiAccounts(%v) error(%v)", mids, err)
		accStats = make(map[int64]*accApi.ProfileStatReply)
	}

	for mid, ins := range delMap {
		if _, ok := accStats[mid]; !ok {
			log.Error("s.monitorNotify() account ")
			accStats[mid] = &accApi.ProfileStatReply{
				Profile: &accApi.Profile{
					Name: "nil",
					Mid:  mid,
				},
			}
		}
		if len(ins) >= threshold {
			var (
				title, content string
			)
			for _, v := range ins {
				if title == "" {
					title = fmt.Sprintf("【异常删稿报警】“%s” 24内已删除%d个自制稿件 ", accStats[mid].Profile.Name, len(ins))
				}
				if content == "" {
					content = fmt.Sprintf("监控规则：%d——%s——%s<br />报警时间：%s<br /><br />", rule.ID, rule.Name, rule.RuleConf.Name, time.Now().Format("2006-01-02 15:04:05"))
					content += fmt.Sprintf("<b>UP主昵称:%s；mid: %d；当前粉丝数:%d； 24内已删除:%d；</b><br /><br />", accStats[mid].Profile.Name, accStats[mid].Profile.Mid, accStats[mid].Follower, len(ins))
					content += "<table border=\"1\" style=\"border-collapse: collapse;\"><tr><th>标题</th><th>av号</th><th>删除时间</th></tr>"
				}
				content += fmt.Sprintf("<tr><td style=\"padding: 5px 10px;\"> %s </td><td style=\"padding: 5px 10px;\"> %d </td><td style=\"padding: 5px 10px;\"> %s </td></tr>", v.Title, v.AID, v.Time)
			}
			content += "</table>"
			if err = s.monitorSendNotify(c, rule.RuleConf.Notify.Way, rule.RuleConf.Notify.Member, title, content); err != nil {
				log.Error("s.moniUpDelArcNotify(%d) s.monitorSendNotify(%d,%v,%s,%s) error:%v", rule.ID, rule.RuleConf.Notify.Way, rule.RuleConf.Notify.Member, title, content, err)
			}
		}
	}
	return
}

// monitorNotifyTime 计算监控报警的score区间
func (s *Service) monitorNotifyTime(conf *moniMdl.RuleConf) (tFrom, tTo int64, err error) {
	now := time.Now().Unix()
	if _, ok := conf.NotifyCdt["time"]; !ok {
		err = errors.New("配置的 NotifyCdt 中不存在 time")
		return
	}
	timeCdt := conf.NotifyCdt["time"].Value
	compCdt := conf.NotifyCdt["time"].Comp
	switch compCdt {
	case moniMdl.CompGT:
		tFrom = 0
		tTo = now - timeCdt
	case moniMdl.CompLT:
		tFrom = now - timeCdt
		tTo = now
	default:
		err = errors.New("配置的 NotifyCdt 中 comparison 不合法: " + compCdt)
		return
	}
	return
}

// reflectIntVal 反射Int值。支持多级查询，比如Addit.MissionID。
func (s *Service) reflectIntVal(obj interface{}, field string, dep int) (val int64, err error) {
	if dep > 10 {
		err = fmt.Errorf("too deep:%d", dep)
		return
	}
	if obj == nil {
		err = errors.New("s.reflectIntVal() obj is invalid memory address or nil pointer dereference")
		return
	}
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		dep += 1
		return s.reflectIntVal(reflect.ValueOf(obj).Elem().Interface(), field, dep)
	}
	if strings.Contains(field, ".") {
		fs := strings.Split(field, ".")
		for i, v := range fs {
			f := strings.Join(fs[i+1:], ".")
			fv := reflect.ValueOf(obj).FieldByName(v)
			if !fv.IsValid() {
				err = fmt.Errorf("s.reflectIntVal() field not found. field:%s obj: %+v", field, obj)
				return
			}
			if fv.IsNil() {
				err = fmt.Errorf("s.reflectIntVal() field is nil. field:%s obj: %+v", field, obj)
				return
			}
			dep += 1
			return s.reflectIntVal(fv.Elem().Interface(), f, dep)
		}
	} else {
		fv := reflect.ValueOf(obj).FieldByName(field)
		if !fv.IsValid() {
			err = fmt.Errorf("s.reflectIntVal() field not found. field:%s obj: %+v", field, obj)
			return
		}
		val = fv.Int()
	}
	return
}

// monitorCompSatisfy 验证值是否满足监控表达式
func (s *Service) monitorCompSatisfy(com string, val int64) (is bool, err error) {
	var (
		v    int64
		vals []int64
	)
	//暂时支持!=、>、<、=、in
	if strings.Contains(com, "!=") {
		//"!=10"
		if v, err = strconv.ParseInt(strings.Replace(com, "!=", "", -1), 10, 64); err != nil {
			return
		}
		is = v != val
	} else if strings.Contains(com, ">=") {
		//">=10"
		if v, err = strconv.ParseInt(strings.Replace(com, ">=", "", -1), 10, 64); err != nil {
			return
		}
		is = val >= v
	} else if strings.Contains(com, "<=") {
		//"<=10"
		if v, err = strconv.ParseInt(strings.Replace(com, "<=", "", -1), 10, 64); err != nil {
			return
		}
		is = val <= v
	} else if strings.Contains(com, "=") {
		//"=10"
		if v, err = strconv.ParseInt(strings.Replace(com, "=", "", -1), 10, 64); err != nil {
			return
		}
		is = v == val
	} else if strings.Contains(com, "!in") {
		//"in(1,2,3)"
		com = strings.Replace(com, "!in(", "", -1)
		com = strings.Replace(com, ")", "", -1)
		if vals, err = xstr.SplitInts(com); err != nil {
			return
		}
		is = true
		for _, v := range vals {
			if val == v {
				is = false
				break
			}
		}
	} else if strings.Contains(com, "in") {
		//"in(1,2,3)"
		com = strings.Replace(com, "in(", "", -1)
		com = strings.Replace(com, ")", "", -1)
		if vals, err = xstr.SplitInts(com); err != nil {
			return
		}
		for _, v := range vals {
			if val == v {
				is = true
				break
			}
		}
	} else if strings.Contains(com, ">") {
		//">10"
		if v, err = strconv.ParseInt(strings.Replace(com, ">", "", -1), 10, 64); err != nil {
			return
		}
		is = val > v
	} else if strings.Contains(com, "<") {
		//"<10"
		if v, err = strconv.ParseInt(strings.Replace(com, "<", "", -1), 10, 64); err != nil {
			return
		}
		is = val < v
	} else {
		err = errors.New("unknown comparison")
	}
	return
}

// monitorSendNotify 发送监控通知
func (s *Service) monitorSendNotify(c context.Context, way int8, members []string, title, content string) (err error) {
	switch way {
	case moniMdl.NotifyTypeEmail:
		log.Info("s.monitorSendNotify() begin. way:%d members:%v title:%s content:%s", way, members, title, content)
		if err = s.email.MonitorEmailAsync(c, members, title, content); err != nil {
			log.Error("s.email.SendMonitorNotify(%v,%s,%s) error:%v", members, title, content, err)
			return
		}
	default:
		err = errors.New("unknown notify way")
		log.Error("s.monitorSendNotify(%d,%v,%s,%s) unknown notify way", way, members, title, content)
		return
	}
	return
}

// multiAccounts 批量获取用户数据
func (s *Service) multiAccounts(c context.Context, mids []int64) (res map[int64]*accApi.ProfileStatReply, err error) {
	var (
		mark map[int64]bool
	)
	res = make(map[int64]*accApi.ProfileStatReply)
	mark = make(map[int64]bool)
	if len(mids) == 0 {
		return
	}
	for _, v := range mids {
		if mark[v] {
			continue
		}
		mark[v] = true
		var r *accApi.ProfileStatReply
		if r, err = s.acc.ProfileWithStat3(c, &accApi.MidReq{Mid: v}); err != nil {
			log.Error("s.multiAccounts() s.acc.ProfileWithStat3(%d) error(%v)", v, err)
			continue
		}
		res[v] = r
	}
	return
}

// monitorEmailProc 发送监控邮件任务
func (s *Service) monitorEmailProc() {
	defer s.wg.Done()
	for {
		s.email.MonitorEmailProc()
		time.Sleep(200 * time.Millisecond)
	}
}

// secondsFormat 将秒转成 时:分:秒
func secondsFormat(sec int) (str string) {
	if sec < 0 {
		return "--:--:--"
	}
	if sec == 0 {
		return "00:00:00"
	}
	h := math.Floor(float64(sec) / 3600)
	m := math.Floor((float64(sec) - 3600*h) / 60)
	se := sec % 60
	return fmt.Sprintf("%02d:%02d:%02d", int64(h), int64(m), se)
}
