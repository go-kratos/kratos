package service

import (
	"context"

	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/task"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// On 上线
func (s *Service) On(c context.Context, opt *common.BaseOptions) (err error) {
	if err = s.mc.ConsumerOn(c, opt); err != nil {
		log.Error("s.mc.ConsumerOn mc错误(%v)", err)
		err = nil
	}
	if err = s.mysql.ConsumerOn(c, opt); err != nil {
		return
	}
	go s.sendTaskConsumerLog(context.TODO(), "on", opt)
	return
}

// Off 离线
func (s *Service) Off(c context.Context, opt *common.BaseOptions) (err error) {
	if err = s.mc.ConsumerOff(c, opt); err != nil {
		log.Error("s.mc.ConsumerOff mc错误(%v)", err)
		err = nil
	}
	if err = s.mysql.ConsumerOff(c, opt); err != nil {
		return
	}
	go s.sendTaskConsumerLog(context.TODO(), "off", opt)
	return
}

// KickOut 踢出
func (s *Service) KickOut(c context.Context, opt *common.BaseOptions, kickuid int64) (err error) {
	opt.UID = kickuid
	unames, _ := s.http.GetUnames(c, []int64{kickuid})
	opt.Uname = unames[kickuid]
	return s.Off(c, opt)
}

// Watcher 监控管理
func (s *Service) Watcher(c context.Context, bizid, flowid int64, role int8) (watchers []*task.WatchItem, err error) {
	// 从数据库拿出24小时内有变化的或者依然在线的
	var wis []*task.WatchItem
	if wis, err = s.mysql.ConsumerStat(c, bizid, flowid); err != nil {
		return
	}

	bopt := &common.BaseOptions{
		BusinessID: bizid,
		FlowID:     flowid,
	}

	var (
		onuids, offuids []int64
		inxmap          = make(map[int64]int)
	)
	for _, item := range wis {
		bopt.UID = item.UID
		bopt.Uname = ""
		uname, urole, err := s.GetRole(c, bopt)
		// 组员列表展示 组员 + 角色获取失败的 + 非任务角色(例如平台管理员)
		if err != nil || role != urole || len(uname) == 0 {
			if role == task.TaskRoleLeader || urole == task.TaskRoleLeader {
				continue
			}
		}

		isOn, _ := s.mc.IsConsumerOn(c, bopt)
		item.Role = urole
		item.IsOnLine = isOn
		item.Uname = uname
		if item.IsOnLine {
			onuids = append(onuids, item.UID)
			item.LastOn = item.Mtime.Format("2006-01-02 15:04:05")
		} else {
			offuids = append(offuids, item.UID)
			item.LastOff = item.Mtime.Format("2006-01-02 15:04:05")
		}
		inxmap[item.UID] = len(watchers)
		watchers = append(watchers, item)
	}

	// 补充laston 或者 lastoff
	wg, ctx := errgroup.WithContext(c)
	if len(onuids) > 0 {
		wg.Go(func() error {
			at, err := s.searchConsumerLog(ctx, bizid, flowid, []string{"off", "kickout"}, onuids, len(onuids))
			if err == nil {
				for uid, ctime := range at {
					watchers[inxmap[uid]].LastOff = ctime
				}
			}
			return err
		})
	}
	if len(offuids) > 0 {
		wg.Go(func() error {
			at, err := s.searchConsumerLog(ctx, bizid, flowid, []string{"on"}, offuids, len(offuids))
			if err == nil {
				for uid, ctime := range at {
					watchers[inxmap[uid]].LastOn = ctime
				}
			}
			return err
		})
	}

	wg.Go(func() error {
		trans := func(c context.Context, ids []int64) (map[int64][]interface{}, error) {
			return s.MemberStats(c, bizid, flowid, ids)
		}
		return s.mulIDtoName(ctx, watchers, trans, "UID", "Count", "CompleteRate", "PassRate", "AvgUT")
	})
	wg.Wait()

	return
}

// IsOn .
func (s *Service) IsOn(c context.Context, opt *common.BaseOptions) bool {
	on, err := s.mc.IsConsumerOn(c, opt)
	if err != nil {
		log.Error("s.mc.ConsumerOff mc错误(%v)", err)
		if on, err = s.mysql.IsConsumerOn(c, opt); err != nil {
			log.Error("s.mysql.ConsumerOff mysql错误(%v)", err)
		}
	}
	return on
}
