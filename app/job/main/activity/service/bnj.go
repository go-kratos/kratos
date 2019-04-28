package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"go-common/app/admin/main/activity/model"
	"go-common/app/job/main/activity/model/bnj"
	"go-common/app/service/main/account/api"
	"go-common/library/log"
)

const (
	_preScore      = 20000
	_stepOne       = 1
	_stepTwo       = 2
	_stepThree     = 3
	_stepFour      = 4
	_stepFlagValue = "1"
)

var bnjSteps = []int{_stepOne, _stepTwo, _stepThree, _stepFour}

func (s *Service) bnjproc() {
	defer s.waiter.Done()
	var (
		c      = context.Background()
		lastTs int64
	)
	for {
		if s.closed {
			return
		}
		if s.bnjTimeFinish == 1 {
			log.Warn("bnjproc bnjTimeFinish")
			return
		}
		msg, ok := <-s.bnjSub.Messages()
		if !ok {
			log.Info("bnjproc: databus consumer exit!")
			return
		}
		msg.Commit()
		m := &bnj.ResetMsg{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("bnjproc json.Unmarshal(%s) error(%+v)", msg.Value, err)
			continue
		}
		if m.Mid <= 0 {
			continue
		}
		// broadcast max every 1s
		if m.Ts-lastTs < 1 {
			continue
		}
		lastTs = m.Ts
		atomic.StoreInt64(&s.bnjLessSecond, s.bnjMaxSecond)
		// default name
		pushMsg := &bnj.Push{Second: s.bnjLessSecond, Name: ""}
		if info, err := s.accClient.Info3(c, &api.MidReq{Mid: m.Mid}); err != nil || info == nil {
			log.Error("bnjproc s.accClient.Info3(%d) error(%v)", m.Mid, err)
		} else {
			var name []rune
			runes := []rune(info.Info.Name)
			nameLen := len(runes)
			if nameLen == 2 {
				name = append(runes[0:1], []rune("*")...)
			} else if nameLen > 2 {
				for i, v := range runes {
					if i == 0 {
						name = append(name, v)
					} else if i == nameLen-1 {
						name = append(name, runes[nameLen-1:]...)
					} else if i == 1 {
						name = append(name, []rune("*")...)
					}
				}
			} else {
				name = runes
			}
			pushMsg.Name = string(name)
			if pushStr, err := json.Marshal(pushMsg); err != nil {
				log.Error("bnjproc json.Marshal(%+v) error(%v)", pushMsg, err)
			} else {
				atomic.StoreInt64(&s.bnjLessSecond, s.bnjMaxSecond)
				log.Info("bnjproc mid(%d) reset lessTime(%d) maxTime(%d)", m.Mid, s.bnjLessSecond, s.bnjMaxSecond)
				if err := s.retryPushAll(context.Background(), string(pushStr), _retryTimes); err != nil {
					log.Error("bnjproc s.bnj.PushAll(%s) error(%v)", string(pushStr), err)
				}
			}
		}
		log.Info("bnjproc key:%s partition:%d offset:%d", msg.Key, msg.Partition, msg.Offset)
	}
}

func bnjWxFlagKey(lid int64, step int) string {
	return fmt.Sprintf("bnj_wx_%d_%d", lid, step)
}

func bnjMsgFlagKey(lid int64, step int) string {
	return fmt.Sprintf("bnj_msg_%d_%d", lid, step)
}

func (s *Service) initBnjSecond() {
	for {
		time.Sleep(time.Second)
		if time.Now().Unix() < s.c.Bnj2019.StartTime.Unix() {
			continue
		}
		break
	}
	if value, err := s.retryCacheTimeFinish(context.Background(), _retryTimes); err != nil {
		log.Error("initBnjSecond s.dao.retryCacheTimeFinish error(%v)", err)
		return
	} else if value > 0 {
		log.Warn("initBnjSecond time finish")
		atomic.StoreInt64(&s.bnjTimeFinish, value)
		return
	}
	// init step flag
	for _, v := range s.c.Bnj2019.Time {
		if v.Step > 0 {
			if value, err := s.retryRsGet(context.Background(), bnjMsgFlagKey(s.c.Bnj2019.LID, v.Step), _retryTimes); err != nil {
				log.Error("initBnjSecond msg s.dao.retryRsGet error(%v)")
			} else if value != "" {
				log.Info("initBnjSecond msg bnjMsgFlagMap[step:%d]", v.Step)
				s.bnjMsgFlagMap[v.Step] = 1
			}
			if value, err := s.retryRsGet(context.Background(), bnjWxFlagKey(s.c.Bnj2019.LID, v.Step), _retryTimes); err != nil {
				log.Error("initBnjSecond wx s.dao.retryRsGet error(%v)")
			} else if value != "" {
				log.Info("initBnjSecond wx bnjWxMsgFlagMap[step:%d]", v.Step)
				s.bnjWxMsgFlagMap[v.Step] = 1
			}
		}
	}
	scores, err := s.retryBatchLikeActSum(context.Background(), []int64{s.c.Bnj2019.LID}, _retryTimes)
	if err != nil {
		// TODO need to restart
		log.Error("initBnjSecond failed s.dao.BatchLikeActSum(%d) error(%v)", s.c.Bnj2019.LID, err)
		return
	}
	if score, ok := scores[s.c.Bnj2019.LID]; ok {
		for i, v := range s.c.Bnj2019.Time {
			if score >= v.Score {
				atomic.StoreInt64(&s.bnjMaxSecond, v.Second)
				atomic.StoreInt64(&s.bnjLessSecond, v.Second)
				break
			}
			if i == len(s.c.Bnj2019.Time)-1 {
				atomic.StoreInt64(&s.bnjMaxSecond, v.Second)
				atomic.StoreInt64(&s.bnjLessSecond, v.Second)
			}
		}
	} else {
		// max second
		atomic.StoreInt64(&s.bnjMaxSecond, s.c.Bnj2019.Time[len(s.c.Bnj2019.Time)-1].Second)
		atomic.StoreInt64(&s.bnjLessSecond, s.c.Bnj2019.Time[len(s.c.Bnj2019.Time)-1].Second)
	}
	if lessSecond, err := s.bnj.CacheLessTime(context.Background()); err != nil {
		log.Error("initBnjSecond s.dao.CacheLessTime error(%v)", err)
	} else if lessSecond > 0 {
		atomic.StoreInt64(&s.bnjLessSecond, lessSecond)
	}
	log.Warn("initBnjSecond maxSecond(%d) lessSecond(%d)", s.bnjMaxSecond, s.bnjLessSecond)
	go s.maxSecondproc()
	go s.lessSecondproc()
	s.waiter.Add(1)
	go s.bnjproc()
}

func (s *Service) maxSecondproc() {
	ctx := context.Background()
	for {
		if s.closed {
			return
		}
		time.Sleep(time.Second)
		if scores, err := s.dao.BatchLikeActSum(context.Background(), []int64{s.c.Bnj2019.LID}); err != nil {
			log.Error("maxSecondproc s.dao.BatchLikeActSum(%d) error(%v)", s.c.Bnj2019.LID, err)
		} else {
			if score, ok := scores[s.c.Bnj2019.LID]; ok {
				for _, v := range s.c.Bnj2019.Time {
					if score >= v.Score {
						atomic.StoreInt64(&s.bnjMaxSecond, v.Second)
						if s.bnjLessSecond > s.bnjMaxSecond {
							atomic.StoreInt64(&s.bnjLessSecond, s.bnjMaxSecond)
						}
						msg := v.Msg
						mc := v.MsgMc
						msgTitle := v.MsgTitle
						step := v.Step
						if step > 0 && s.bnjMsgFlagMap[step] == 0 {
							if err = s.retryRsSet(ctx, bnjMsgFlagKey(s.c.Bnj2019.LID, step), _stepFlagValue, _retryTimes); err != nil {
								log.Error("s.retryRsSet(%d,%d) error(%v)", s.c.Bnj2019.LID, step, err)
								break
							}
							if msg != "" && msgTitle != "" && mc != "" {
								go s.sendMessageToSubs(ctx, s.c.Bnj2019.LID, mc, msgTitle, msg, _retryTimes)
							} else {
								log.Error("bnj msg conf step(%d) error", step)
								break
							}
							log.Info("bnj send msg step:%d finish", step)
							s.bnjMsgFlagMu.Lock()
							s.bnjMsgFlagMap[step] = 1
							s.bnjMsgFlagMu.Unlock()
						}
						break
					}
				}
				for _, v := range s.c.Bnj2019.Time {
					if score+_preScore >= v.Score {
						wxMsg := v.WxMsg
						step := v.Step
						if step > 0 && s.bnjWxMsgFlagMap[step] == 0 {
							if err = s.retryRsSet(ctx, bnjWxFlagKey(s.c.Bnj2019.LID, step), _stepFlagValue, _retryTimes); err != nil {
								log.Error("s.retryRsSet(%d,%d) error(%v)", s.c.Bnj2019.LID, step, err)
								break
							}
							if wxMsg != "" && s.c.Bnj2019.WxUser != "" {
								if err = s.retrySendWechat(ctx, s.c.Bnj2019.WxTitle, wxMsg, s.c.Bnj2019.WxUser, _retryTimes); err != nil {
									log.Error("s.retrySendWechat(%s,%s) error(%v)", s.c.Bnj2019.WxTitle, wxMsg, err)
									break
								}
							} else {
								log.Error("bnj wx msg conf step(%d) error", step)
								break
							}
							log.Info("bnj send wx step:%d finish", step)
							s.bnjWxMsgFlagMu.Lock()
							s.bnjWxMsgFlagMap[step] = 1
							s.bnjWxMsgFlagMu.Unlock()
						}
						break
					}
				}
			} else {
				log.Warn("maxSecondproc lid not found")
			}
		}
	}
}

func (s *Service) lessSecondproc() {
	for {
		if s.closed {
			return
		}
		time.Sleep(time.Second)
		atomic.AddInt64(&s.bnjLessSecond, -1)
		if s.c.Bnj2019.GameCancel != 0 {
			log.Warn("lessSecondproc bnj game cancel")
			atomic.StoreInt64(&s.bnjLessSecond, 0)
		}
		if s.bnjLessSecond <= 0 {
			if err := s.retryAddCacheTimeFinish(context.Background(), 1, _retryTimes); err != nil {
				log.Error("lessSecondproc s.bnj.AddCacheTimeFinish error(%v)", err)
				continue
			}
			log.Warn("lessSecondproc bnj time Finish")
			atomic.StoreInt64(&s.bnjTimeFinish, 1)
			pushMsg := &bnj.Push{Second: 0, Name: "", TimelinePic: s.c.Bnj2019.TimelinePic, H5TimelinePic: s.c.Bnj2019.H5TimelinePic}
			if pushStr, err := json.Marshal(pushMsg); err != nil {
				log.Error("lessSecondproc json.Marshal(%+v) error(%v)", pushMsg, err)
			} else {
				atomic.StoreInt64(&s.bnjLessSecond, s.bnjMaxSecond)
				if err := s.retryPushAll(context.Background(), string(pushStr), _retryTimes); err != nil {
					log.Error("lessSecondproc s.bnj.PushAll error(%v)", err)
				}
			}
			return
		}
		pushMsg := &bnj.Push{Second: s.bnjLessSecond, Name: ""}
		if pushStr, err := json.Marshal(pushMsg); err != nil {
			log.Error("lessSecondproc json.Marshal(%+v) error(%v)", pushMsg, err)
		} else {
			if err := s.retryPushAll(context.Background(), string(pushStr), 1); err != nil {
				log.Error("lessSecondproc s.bnj.PushAll error(%v)", err)
			}
		}
	}
}

func (s *Service) retryBatchLikeActSum(c context.Context, lids []int64, retryCnt int) (res map[int64]int64, err error) {
	for i := 0; i < retryCnt; i++ {
		if res, err = s.dao.BatchLikeActSum(c, lids); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (s *Service) retryAddCacheTimeFinish(c context.Context, value int64, retryCnt int) (err error) {
	for i := 0; i < retryCnt; i++ {
		if err = s.bnj.AddCacheTimeFinish(c, value); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (s *Service) retryCacheTimeFinish(c context.Context, retryCnt int) (value int64, err error) {
	for i := 0; i < retryCnt; i++ {
		if value, err = s.bnj.CacheTimeFinish(c); err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return
}

func (s *Service) retryPushAll(c context.Context, msg string, retryCnt int) (err error) {
	for i := 0; i < retryCnt; i++ {
		if err = s.bnj.PushAll(c, msg); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (s *Service) cronInformationMessage() {
	log.Info("cronInformationMessage start cron")
	if s.c.Bnj2019.LID == 0 || s.c.Bnj2019.MidLimit == 0 {
		log.Error("cronInformationMessage conf error")
		return
	}
	var (
		c              = context.Background()
		title, msg, mc string
	)
	for _, v := range s.c.Bnj2019.Message {
		if time.Now().Unix() >= v.Start.Unix() {
			title = v.Title
			msg = v.Content
			mc = v.Mc
			break
		}
	}
	if title == "" || msg == "" || mc == "" {
		log.Error("cronInformationMessage message conf error")
		return
	}
	s.sendMessageToSubs(c, s.c.Bnj2019.LID, mc, title, msg, _retryTimes)
	log.Info("cronInformationMessage finish title(%s)", title)
}

func (s *Service) sendMessageToSubs(c context.Context, lid int64, mc, title, msg string, retryCnt int) {
	var minID int64
	log.Info("sendMessageToSubs mc:%s title:%s start", mc, title)
	for {
		time.Sleep(100 * time.Millisecond)
		actions, err := s.retryLikeActList(c, lid, minID, s.c.Bnj2019.MidLimit, retryCnt)
		if err != nil {
			log.Error("sendMessageToSubs s.dao.LikeActList(%d,%d,%d) error(%v)", lid, minID, s.c.Bnj2019.MidLimit, err)
			continue
		}
		if len(actions) == 0 {
			log.Info("sendMessageToSubs finish")
			break
		}
		var mids []int64
		for i, v := range actions {
			if v.Mid > 0 {
				mids = append(mids, v.Mid)
			}
			if i == len(actions)-1 {
				minID = v.ID
			}
		}
		if len(mids) == 0 {
			continue
		}
		if err = s.retrySendMessage(c, mids, mc, title, msg, _retryTimes); err != nil {
			log.Error("sendMessageToSubs s.dao.SendMessage(mids:%v) error(%v)", mids, err)
		}
	}
}

func (s *Service) retryLikeActList(c context.Context, lid, minID, limit int64, retryCnt int) (list []*model.LikeAction, err error) {
	for i := 0; i < retryCnt; i++ {
		if list, err = s.dao.LikeActList(c, lid, minID, limit); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (s *Service) retrySendMessage(c context.Context, mids []int64, mc, title, msg string, retryCnt int) (err error) {
	for i := 0; i < retryCnt; i++ {
		if err = s.bnj.SendMessage(c, mids, mc, title, msg); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (s *Service) retryRsSet(c context.Context, key, value string, retryCnt int) (err error) {
	for i := 0; i < retryCnt; i++ {
		if err = s.dao.RsSet(c, key, value); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (s *Service) retryRsGet(c context.Context, key string, retryCnt int) (value string, err error) {
	for i := 0; i < retryCnt; i++ {
		if value, err = s.dao.RsGet(c, key); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func (s *Service) retrySendWechat(c context.Context, title, msg, user string, retryCnt int) (err error) {
	for i := 0; i < retryCnt; i++ {
		if err = s.bnj.SendWechat(c, title, msg, user); err == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}
