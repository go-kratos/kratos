package service

import (
	"context"
	"fmt"
	"go-common/library/sync/errgroup"
	"strings"
	"time"

	"go-common/app/admin/main/videoup/model/archive"
	"go-common/library/log"
)

// HandsUp 签入
func (s *Service) HandsUp(c context.Context, uid int64, uname string) (err error) {
	if s.CheckOnline(c, uid) {
		log.Info("已经登入(%d)", uid)
		return
	}

	_, err = s.arc.TaskUserCheckIn(c, uid)
	if err != nil {
		log.Error("s.arc.TaskUserCheckIn(%d) error(%v)", uid, err)
		return
	}
	s.sendConsumerLog(c, &archive.ConsumerLog{
		UID:    uid,
		Uname:  uname,
		Action: archive.ActionHandsUP,
		Ctime:  time.Now().Format(archive.TimeFormatSec),
		Desc:   "checkin",
	})

	mapParas := map[string]interface{}{
		"action": archive.ActionHandsUP,
		"uid":    uid,
	}

	if _, err = s.arc.AddTaskHis(c, 0, archive.ActionHandsUP, 0, 0, uid, 0, 0, "checkin"); err != nil {
		log.Error("s.arc.AddTaskLog(%v) error(%v)", mapParas, uid)
		return
	}
	log.Info("用户签入(%d)", uid)
	return
}

// HandsOff 签出
func (s *Service) HandsOff(c context.Context, uid int64, fuid int64) (err error) {
	if fuid != 0 { //管理员强制踢出组员
		if !s.isLeader(c, uid) {
			return fmt.Errorf("只有组长能强制踢出")
		}
		log.Info("管理员%d踢出组员%d", uid, fuid)
		uid = fuid
	}

	err = s.checkOut(c, uid)
	if err != nil {
		log.Error("s.checkOut(%d) error(%v)", uid, err)
		return
	}
	s.sendConsumerLog(c, &archive.ConsumerLog{
		UID:    uid,
		Uname:  "",
		Action: archive.ActionHandsOFF,
		Ctime:  time.Now().Format(archive.TimeFormatSec),
		Desc:   "checkout",
	})
	s.Free(c, uid)
	return
}

// Online 用户列表
func (s *Service) Online(c context.Context) (cms []*archive.Consumers, err error) {
	cms, err = s.arc.Consumers(c)
	if err != nil {
		log.Error("s.arc.Consumers error(%v)", err)
		return
	}

	if len(cms) > 0 {
		var wg errgroup.Group
		wg.Go(func() error {
			if err := s.mulIDtoName(c, cms, s.mng.GetNameByUID, "UID", "UserName"); err != nil {
				log.Error("mulIDtoName s.mng.GetNameByUID error(%v)", err)
			}
			return nil
		})
		wg.Go(func() error {
			if err := s.mulIDtoName(c, cms, s.search.OutTime, "UID", "LastOut"); err != nil {
				log.Error("mulIDtoName s.search.OutTime error(%v)", err)
			}
			return nil
		})
		wg.Wait()
	}

	return
}

// InOutList 用户登入登出历史
func (s *Service) InOutList(c context.Context, unames string, bt, et string) (l []*archive.InQuit, err error) {
	uids := []int64{}
	if len(unames) > 0 {
		if res, err := s.mng.Uids(c, strings.Split(unames, ",")); err == nil {
			for _, uid := range res {
				uids = append(uids, uid)
			}
		}
	}
	// 前端参数是日期，搜索参数必须到秒
	if len(bt) > 0 && len(et) > 0 {
		bt = bt + " 00:00:00"
		et = et + " 23:59:59"
	}

	return s.search.InQuitList(c, uids, bt, et)
}

// CheckOnline 检查在线状态
func (s *Service) CheckOnline(c context.Context, uid int64) (on bool) {
	if s.arc.IsConsumerOn(c, uid) == 1 {
		on = true
	}
	return
}

// CheckGroup 检查用户组权限
func (s *Service) CheckGroup(c context.Context, uid int64) (role int8, err error) {
	role, err = s.mng.GetUserRole(c, uid)
	if err != nil || role == 0 {
		log.Error("非法用户(%d) error(%v)", uid, err)
		return
	}
	return
}

func (s *Service) checkOut(c context.Context, uid int64) (err error) {
	if s.arc.IsConsumerOn(c, uid) == 0 {
		log.Info("已经签出(%d)", uid)
		return
	}

	_, err = s.arc.TaskUserCheckOff(c, uid)
	if err != nil {
		log.Error("s.arc.TaskUserCheckOff(%d) error(%v)", uid, err)
		return
	}

	mapParas := map[string]interface{}{
		"action": archive.ActionHandsOFF,
		"uid":    uid,
	}

	if _, err = s.arc.AddTaskHis(c, 0, archive.ActionHandsOFF, 0, 0, uid, 0, 0, "checkOut"); err != nil {
		log.Error("s.arc.AddTaskLog(%v) error(%v)", mapParas, uid)
	}
	return
}
