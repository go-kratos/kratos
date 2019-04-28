package service

import (
	"context"
	"go-common/app/job/bbq/video/dao"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"sync/atomic"
	"time"

	notice "go-common/app/service/bbq/notice-service/api/v1"
)

// 信号量，为了避免task执行超过周期，连续执行两个任务会出现问题
var i32 int32

// SysMsgTask 推送消息
func (s *Service) SysMsgTask() {
	ctx := context.Background()
	newValue := atomic.AddInt32(&i32, 1)
	defer atomic.AddInt32(&i32, -1)
	if newValue > 1 {
		log.Errorv(ctx, log.KV("log", "sysMsgTask pending"))
		return
	}

	res, err := s.dao.RawCheckTask(ctx, "checkSysMsg")
	if err != nil {
		log.Errorv(ctx, log.KV("log", "get last sysMsgTask id fail"))
		return
	}

	lastSysMsgID := res.LastCheck
	curSysMsgID := res.LastCheck
	list, err := s.dao.GetNewSysMsg(ctx, curSysMsgID)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "get new sysMsg fail"))
		return
	}
	if len(list) == 0 {
		log.Infov(ctx, log.KV("log", "no new sys msg to be sync to notice"))
		return
	}

	var mids []int64
	for _, item := range list {
		curSysMsgID = item.Id

		notice := notice.NoticeBase{
			Mid:        item.Receiver,
			ActionMid:  item.Sender,
			SvId:       0,
			NoticeType: 4,
			Text:       item.Text,
			JumpUrl:    item.JumpUrl,
			BizType:    dao.NoticeBizTypeSysMsg,
			BizId:      item.Id,
		}
		// 全量系统消息
		if item.Receiver == 0 {
			lastUserID := int64(0)
			if len(mids) == 0 {
				for {
					if userBases, err := s.dao.UsersByLast(ctx, lastUserID); err != nil {
						log.Errorv(ctx, log.KV("log", "sys msg task: get user base fail"))
						break
					} else {
						for _, userBase := range userBases {
							mids = append(mids, userBase.MID)
						}
						if len(userBases) > 0 {
							lastUserID = userBases[len(userBases)-1].ID
						} else {
							break
						}
					}

				}
			}
			midChan := make(chan int64, 20)
			go func() {
				for _, mid := range mids {
					midChan <- mid
				}
				close(midChan)
			}()
			startTime := time.Now()
			g := errgroup.Group{}
			for i := 0; i < 10; i++ {
				g.Go(func() error {
					subNotice := notice
					for mid := range midChan {
						subNotice.Mid = mid
						s.dao.CreateNotice(ctx, &subNotice)
					}
					return nil
				})
			}
			g.Wait()
			log.Info("total sys msg notice push: cost_time=%f, mid_len=%d", time.Since(startTime).Seconds(), len(mids))
		} else {
			s.dao.CreateNotice(ctx, &notice)
		}
	}

	if _, err := s.dao.UpdateTaskLastCheck(ctx, "checkSysMsg", curSysMsgID); err != nil {
		log.Errorv(ctx, log.KV("log", "update check_task mysql fail"))
		return
	}

	log.Infov(ctx, log.KV("log", "no new sys msg to be sync to notice"), log.KV("last_sys_id", lastSysMsgID), log.KV("cur_sys_id", curSysMsgID))
}
