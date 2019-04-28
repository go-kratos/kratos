package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/aegis/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

func (s *Service) initCache() {
	s.newactiveBizFlow = make(map[string]struct{})
	s.syncConfigCache(context.Background())
	s.syncConsumerCache(context.Background())
	s.oldactiveBizFlow = s.newactiveBizFlow
}

func (s *Service) cacheProc() {
	for {
		s.syncTaskCache()
		time.Sleep(3 * time.Minute)
		s.syncConfigCache(context.Background())
		s.syncWeightWatch(context.Background())
		s.syncConsumerCache(context.Background())
	}
}

func (s *Service) syncTaskCache() {
	var (
		tasks  []*model.Task
		lastid = int64(0)
		err    error
	)

	// 1.停滞任务,10分钟未变化，检查是否遗漏
	lastid = 0
	for {
		mtime := time.Now().Add(-10 * time.Minute)
		if tasks, lastid, err = s.dao.QueryTask(context.Background(), model.TaskStateDispatch, mtime, lastid, 1000); err != nil || len(tasks) == 0 {
			break
		}
		for _, task := range tasks {
			log.Info("检测到遗漏 停滞任务(%+v)", task)
			s.dao.SetTask(context.Background(), task)
			s.dao.PushPersonalTask(context.Background(), task.BusinessID, task.FlowID, task.UID, task.ID)
		}
		time.Sleep(time.Second)
	}

	// 2.延迟任务,半小时未变化，检查是否遗漏
	lastid = 0
	for {
		mtime := time.Now().Add(-30 * time.Minute)
		if tasks, lastid, err = s.dao.QueryTask(context.Background(), model.TaskStateDelay, mtime, lastid, 1000); err != nil || len(tasks) == 0 {
			break
		}
		for _, task := range tasks {
			log.Info("检测到遗漏 延迟任务(%+v)", task)
			s.dao.SetTask(context.Background(), task)
			s.dao.PushDelayTask(context.Background(), task.BusinessID, task.FlowID, task.UID, task.ID)
		}
		time.Sleep(time.Second)
	}

	// 3.实时任务,1小时未变化，检查是否遗漏
	lastid = 0
	for {
		mtime := time.Now().Add(-60 * time.Minute)
		if tasks, lastid, err = s.dao.QueryTask(context.Background(), model.TaskStateInit, mtime, lastid, 1000); err != nil || len(tasks) == 0 {
			break
		}
		for _, task := range tasks {
			log.Info("检测到遗漏 实时任务(%+v)", task)
			s.dao.SetTask(context.Background(), task)
		}

		s.dao.PushPublicTask(context.Background(), tasks...)
		time.Sleep(time.Second)
	}
}

func (s *Service) syncConfigCache(c context.Context) (err error) {
	s.oldactiveBizFlow = s.newactiveBizFlow

	configs, err := s.dao.TaskActiveConfigs(c)
	if err != nil {
		return
	}

	rangeWCCache := make(map[int64]map[string]*model.RangeWeightConfig)
	equalWCCache := make(map[string][]*model.EqualWeightConfig)
	assignCache := make(map[string][]*model.AssignConfig)
	activeBizFlow := make(map[string]struct{})

	for _, item := range configs {
		key := fmt.Sprintf("%d-%d", item.BusinessID, item.FlowID)
		activeBizFlow[key] = struct{}{}
		switch item.ConfType {
		case model.TaskConfigAssign:
			assign := new(struct {
				Mids string `json:"mids"`
				Uids string `json:"uids"`
			})
			if err = json.Unmarshal([]byte(item.ConfJSON), assign); err != nil {
				log.Error("json.Unmarshal error(%v)", err)
				continue
			}
			ac := &model.AssignConfig{}
			if item.UID > 0 {
				ac.Admin = item.UID
			} else {
				ac.Admin = 399
			}
			assign.Mids = strings.TrimSpace(assign.Mids)
			assign.Uids = strings.TrimSpace(assign.Uids)
			if ac.Mids, err = xstr.SplitInts(assign.Mids); err != nil {
				log.Error("xstr.SplitInts error(%v)", err)
				continue
			}
			if ac.Uids, err = xstr.SplitInts(assign.Uids); err != nil {
				log.Error("xstr.SplitInts error(%v)", err)
				continue
			}

			if _, ok := assignCache[key]; ok {
				assignCache[key] = append(assignCache[key], ac)
			} else {
				assignCache[key] = []*model.AssignConfig{ac}
			}
		case model.TaskConfigRangeWeight:
			wcitem := &model.RangeWeightConfig{}
			if err = json.Unmarshal([]byte(item.ConfJSON), wcitem); err != nil {
				log.Error("json.Unmarshal error(%v)", err)
				continue
			}

			if _, ok := rangeWCCache[item.BusinessID]; ok {
				rangeWCCache[item.BusinessID][wcitem.Name] = wcitem
			} else {
				rangeWCCache[item.BusinessID] = map[string]*model.RangeWeightConfig{
					wcitem.Name: wcitem,
				}
			}

		case model.TaskConfigEqualWeight:
			ewcitem := &model.EqualWeightConfig{}
			if err = json.Unmarshal([]byte(item.ConfJSON), ewcitem); err != nil {
				log.Error("json.Unmarshal error(%v)", err)
				continue
			}

			ewcitem.Uname = item.Uname
			ewcitem.Description = item.Description
			ewcitem.IDs = strings.TrimSpace(ewcitem.IDs)
			if _, ok := equalWCCache[key]; ok {
				equalWCCache[key] = append(equalWCCache[key], ewcitem)
			} else {
				equalWCCache[key] = []*model.EqualWeightConfig{ewcitem}
			}
		}
	}
	s.rangeWeightCfg = rangeWCCache
	s.equalWeightCfg = equalWCCache
	s.assignConfig = assignCache
	s.newactiveBizFlow = activeBizFlow

	return
}

func (s *Service) syncWeightWatch(c context.Context) {
	for key := range s.oldactiveBizFlow {
		if _, ok := s.newactiveBizFlow[key]; !ok {
			if wm, ok := s.wmHash[key]; ok {
				wm.close = true
				log.Info("关闭权重计算器 bizid(%d) flowid(%d)", wm.businessID, wm.flowID)
				delete(s.wmHash, key)
			}
		}
	}

	for key := range s.newactiveBizFlow {
		if _, ok := s.oldactiveBizFlow[key]; !ok {
			bizid, _ := parseKey(key)
			s.wmHash[key] = NewWeightManager(s, s.getWeightOpt(bizid), key)
		}
	}
}

func (s *Service) getWeightOpt(bizid int) *model.WeightOPT {
	for _, item := range s.c.BizCfg.WeightOpt {
		if item.BusinessID == int64(bizid) {
			return item
		}
	}
	return nil
}

func (s *Service) syncConsumerCache(c context.Context) (err error) {
	s.ccMux.Lock()
	defer s.ccMux.Unlock()
	consumerCache, err := s.dao.TaskActiveConsumer(c)
	if err != nil {
		return
	}
	s.consumerCache = consumerCache
	return
}

// getWeightCache .
func (s *Service) getWeightCache(c context.Context, businessid, flowid int64) (rwc map[string]*model.RangeWeightConfig, ewc []*model.EqualWeightConfig) {
	key := fmt.Sprintf("%d-%d", businessid, flowid)
	rwc = s.rangeWeightCfg[businessid]
	ewc = s.equalWeightCfg[key]
	return
}
