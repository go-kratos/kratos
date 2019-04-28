package service

import (
	"context"
	"fmt"
	"go-common/app/job/main/videoup-report/model/archive"
	"go-common/app/job/main/videoup-report/model/monitor"
	"go-common/library/log"
	"time"
)

// hdlMonitorArc deal with archive stay stats
func (s *Service) hdlMonitorArc(nw, old *archive.Archive) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.hdlMonitorArc panic(%v)", pErr)
		}
	}()
	var (
		oKey, nKey string
		kFormat    = monitor.RedisPrefix + monitor.SuffixArc
		addit      *archive.Addit
	)
	log.Info("hdlMonitorArc (%v,%v)", nw, old)
	if addit, err = s.arc.Addit(context.TODO(), nw.ID); err != nil {
		log.Error("s.hdlMonitorArc() s.arc.Addit(%d) error(%v)", nw.ID, err)
		return
	}
	//去掉PGC稿件
	if addit != nil && (addit.UpFrom == archive.UpFromPGC || addit.UpFrom == archive.UpFromPGCSecret || addit.UpFrom == archive.UpFromCoopera) {
		return
	}
	if old != nil {
		if nw.Round == old.Round && nw.State == old.State {
			return
		}
		oKey = fmt.Sprintf(kFormat, monitor.BusArc, old.Round, old.State)
	}
	//忽略Round 99 state不为-6,-1,0,1的数据
	if nw.Round != archive.RoundEnd && (nw.State == archive.StateForbidFixed || nw.State == archive.StateForbidWait || nw.State == archive.StateOpen || nw.State == archive.StateOrange) {
		nKey = fmt.Sprintf(kFormat, monitor.BusArc, nw.Round, nw.State)
	}
	//回查忽略活动稿件
	if addit != nil && addit.MissionID > 0 && (nw.Round == archive.RoundReviewFirst || nw.Round == archive.RoundReviewFirstWaitTrigger || nw.Round == archive.RoundReviewSecond || nw.Round == archive.RoundTriggerClick) {
		nKey = ""
	}
	log.Info("hdlMonitorArc () s.monitorSave(%s,%s,%d)", oKey, nKey, nw.ID)
	err = s.monitorSave(oKey, nKey, nw.ID)
	return
}

// hdlMonitorVideo 视频审核监控
func (s *Service) hdlMonitorVideo(nv, ov *archive.Video) (err error) {
	var (
		oKey, nKey string
		kFormat    = monitor.RedisPrefix + monitor.SuffixVideo
	)
	log.Info("hdlMonitorVideo (%v,%v)", nv, ov)
	if ov != nil {
		if nv.Status == ov.Status {
			return
		}
		oKey = fmt.Sprintf(kFormat, monitor.BusVideo, ov.Status)
	}
	if nv.Status == archive.VideoStatusSubmit || nv.Status == archive.VideoStatusWait {
		nKey = fmt.Sprintf(kFormat, monitor.BusVideo, nv.Status)
	}
	log.Info("hdlMonitorVideo () s.monitorSave(%s,%s,%d) filename(%s)", oKey, nKey, nv.ID, nv.Filename)
	err = s.monitorSave(oKey, nKey, nv.ID)
	return
}

func (s *Service) monitorSave(oKey, nKey string, oid int64) (err error) {
	defer func() {
		if pErr := recover(); pErr != nil {
			log.Error("s.monitorSave panic(%v)", pErr)
		}
	}()
	var c = context.TODO()
	if oKey != "" {
		if err = s.redis.RemMonitorStats(c, oKey, oid); err != nil {
			log.Error("s.monitorSave() RemMonitorStats(%s,%d) error(%v)", oKey, oid, err)
		}
		s.redis.ClearMonitorStats(c, oKey)
	}
	if nKey != "" {
		if err = s.redis.AddMonitorStats(c, nKey, oid); err != nil {
			log.Error("s.monitorSave() AddMonitorStats(%s,%d) error(%v)", nKey, oid, err)
			return
		}
		s.redis.ClearMonitorStats(c, nKey)
	}
	if err != nil {
		log.Error("s.monitorSave(%s,%s,%d) error(%v)", oKey, nKey, oid, err)
	}
	return
}

//monitorNotifyEmail 发送监控通知
func (s *Service) monitorNotify() {
	var (
		c    = context.TODO()
		data []*monitor.RuleResultData
		err  error
	)
	defer func() {
		if err := recover(); err != nil {
			log.Error("monitorNotifyEmail() panic(%v)", err)
		}
	}()
	// 从admin获取报警数据
	if data, err = s.dataDao.MonitorNotify(c); err != nil {
		log.Error("s.dataDao.MonitorNotify() error(%v)", err)
		return
	}
	for _, v := range data {
		if v.Rule.State != monitor.RuleStateOK {
			log.Error("monitorNotify() ignore rule(%d) state(%d)", v.Rule.ID, v.Rule.State)
			continue
		}
		subject := fmt.Sprintf("%s监控", v.Rule.Name)
		body := fmt.Sprintf("当前滞留时间为%s超过阀值，滞留量为%d，整体量为%d \n%s", secondsFormat(v.Stats.MaxTime), v.Stats.MoniCount, v.Stats.TotalCount, time.Now().Format("2006-01-02 15:04:05"))
		url := ""
		switch v.Rule.Business {
		case monitor.BusVideo:
			url = fmt.Sprintf("http://manager.bilibili.co/#!/video/list?monitor_list=%d_%d_%d", v.Rule.Type, v.Rule.Business, v.Rule.ID)
		case monitor.BusArc:
			url = fmt.Sprintf("http://manager.bilibili.co/#!/archive_utils/all?monitor_list=%d_%d_%d", v.Rule.Type, v.Rule.Business, v.Rule.ID)
		}
		body += fmt.Sprintf("\n跳转链接：%s", url)
		if v.Rule.RuleConf.Notify.Way == monitor.NotifyTypeEmail {
			tpl := s.email.MonitorNotifyTemplate(subject, body, v.Rule.RuleConf.Notify.Member)
			log.Info("monitorNotify() email template(%v)", *tpl)
			s.email.PushToRedis(c, tpl)
		} else {
			log.Error("monitorNotify() unknown notify rule(%d) type(%s)", v.Rule.ID, v.Rule.RuleConf.Notify.Way)
			continue
		}
	}
}
