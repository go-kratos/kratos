package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/push/dao"
	pb "go-common/app/service/main/push/api/grpc/v1"
	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/log"
)

const (
	_retryCallback    = 5
	_delCallbackLimit = 5000
)

func (s *Service) callbackproc() {
	defer s.waiter.Done()
	var err error
	for {
		msg, ok := <-s.callbackCh
		if !ok {
			log.Warn("s.callbackproc() closed")
			return
		}
		for _, v := range msg {
			if v == nil {
				continue
			}
			arg := &pb.AddCallbackRequest{
				Task:     v.Task,
				APP:      v.APP,
				Platform: int32(v.Platform),
				Mid:      v.Mid,
				Pid:      int32(v.Pid),
				Token:    v.Token,
				Buvid:    v.Buvid,
				Click:    int32(v.Click),
			}
			if v.Extra != nil {
				arg.Extra = &pb.CallbackExtra{Status: int32(v.Extra.Status), Channel: int32(v.Extra.Channel)}
			}
			for i := 0; i < _retryCallback; i++ {
				if _, err = s.pushRPC.AddCallback(context.Background(), arg); err == nil {
					break
				}
				time.Sleep(20 * time.Millisecond)
			}
			if err != nil {
				log.Error("s.pushRPC.AddCallback(%+v) error(%v)", arg, err)
				dao.PromError("report:新增callback")
				continue
			}
			log.Info("add callback success task(%s) token(%s)", v.Task, v.Token)
			time.Sleep(time.Millisecond)
		}
	}
}

// consumeCallback consumes callback.
func (s *Service) consumeCallback() {
	defer s.waiter.Done()
	for {
		msg, ok := <-s.callbackSub.Messages()
		if !ok {
			log.Info("databus: push-job callback consumer exit!")
			close(s.callbackCh)
			return
		}
		s.callbackCnt++
		msg.Commit()
		var cbs []*pushmdl.Callback
		if err := json.Unmarshal(msg.Value, &cbs); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			dao.PromError("service:解析databus中callback消息")
			continue
		}
		log.Info("consumeCallback key(%s) partition(%d) offset(%d) msg(%v)", msg.Key, msg.Partition, msg.Offset, string(msg.Value))
		s.callbackCh <- cbs
	}
}

func (s *Service) delCallbacksproc() {
	for {
		now := time.Now()
		// 每天4点时删除七天前的callback数据
		if now.Hour() == 4 {
			var (
				err     error
				deleted int64
				b       = now.Add(time.Duration(-s.c.Job.DelCallbackInterval*24) * time.Hour)
				loc, _  = time.LoadLocation("Local")
				t       = time.Date(b.Year(), b.Month(), b.Day(), 23, 59, 59, 0, loc)
			)
			for {
				if deleted, err = s.dao.DelCallbacks(context.TODO(), t, _delCallbackLimit); err != nil {
					log.Error("s.delCallbacks(%v) error(%v)", t, err)
					s.dao.SendWechat("DB操作失败:push-job删除callback数据错误")
					time.Sleep(time.Second)
					continue
				}
				if deleted < _delCallbackLimit {
					break
				}
				time.Sleep(time.Second)
			}
			log.Info("delCallbacksproc success date(%v)", t)
			time.Sleep(time.Hour)
		}
		time.Sleep(time.Minute)
	}
}
