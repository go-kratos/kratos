package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/library/conf/env"
	"go-common/library/net/metadata"

	"go-common/app/job/live/xroom-feed/internal/model"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
)

const (
	_whiteListType  = 1
	_condType       = 2
	_fileListType   = 3
	_areaType       = "area"
	_onlineType     = "online"
	_incomeType     = "income"
	_dmsType        = "dms"
	_liveDaysType   = "live_days"
	_hourRankType   = "hour_rank"
	_anchorCateType = "anchor_cate"
	_roomStatusType = "room_status"
	_condAnd        = "and"
	_condOr         = "or"

	_confTypeString = "string"
	_confTypeRange  = "range"
	_confTypeTop    = "top"
)

// JobTTl ...
type JobTTl struct {
	ttl int
}

func (s *Service) reloadConfFromDb() {
	t1 := time.NewTicker(time.Second * 5)
	defer t1.Stop()
	for {
		select {
		case <-t1.C:
			confList, err := s.dao.GetConfFromDb()
			if err != nil {
				log.Error("[reloadConfFromDb]getConfFromDb_error:%+v", err)
				continue
			}

			key, expire := s.getRecConfKey()
			listStr, err := json.Marshal(confList)
			if err == nil {
				s.dao.SetRecPoolCache(context.TODO(), key, string(listStr), expire)
			}
			s.ruleConf.Store(confList)
		}
	}
}

func (s *Service) loadConfFromDb() {
	confList, err := s.dao.GetConfFromDb()
	if err != nil {
		log.Error("[loadConfFromDb]getConfFromDb_error:%+v", err)
		return
	}
	s.ruleConf.Store(confList)
}

func (s *Service) parseCondConf(cond string) (condConf *model.RuleProtocol, err error) {
	condConf = new(model.RuleProtocol)
	err = json.Unmarshal([]byte(cond), &condConf)
	if err != nil {
		log.Error("[parseCondConf]Unmarshal err: %+v", err)
		return
	}

	return
}

func (s *Service) reloadRecList() {
	t1 := time.NewTicker(time.Second * 20)
	defer t1.Stop()
	for {
		select {
		case <-t1.C:
			s.genRecListJob()
		}
	}
}

// 对所有配置的规则取数据的主程序
func (s *Service) genRecListJob() {
	log.Info("[genRecListJob]genRecListJob start")
	ttl := 20
	ttl, err := s.ac.Get("JobTtl").Int()
	if err != nil {
		log.Error("[genRecListJob]getJobTtlConfFromSven_error:%+v", err)
	}
	confList := s.ruleConf.Load()
	// assert
	res, ok := confList.([]*model.RecPoolConf)
	if !ok {
		log.Error("[genRecListJob]conf assert error! %+v", confList)
		return
	}
	if len(res) <= 0 {
		log.Warn("[genRecListJob]confList_empty")
		return
	}

	//可配置时间结束子任务 内部要有阻塞操作
	cCtx := metadata.NewContext(context.TODO(), metadata.MD{metadata.Color: env.Color})
	ctx, cancel := context.WithTimeout(cCtx, time.Duration(ttl)*time.Second)
	wg := errgroup.Group{}

	// 每条规则之间无关联，可以不互相cancel，只看超时ctx (可在内部对每个操作设置)
	// 规则获取loop
	for _, rule := range res {
		// 白名单规则
		if rule.ConfType == _whiteListType || rule.ConfType == _fileListType {
			ruleId := rule.Id
			wg.Go(func() error {
				listStr, err := s.dao.GetWhiteList(ctx, ruleId)
				if err != nil {
					log.Error("[genRecListJob]GetWhiteList_err:%+v", err)
					return nil
				}

				listStrArr := strings.Split(listStr, ",")
				listIntArr := make([]int64, 0)
				for _, roomIdStr := range listStrArr {
					roomIdInt, _ := strconv.ParseInt(roomIdStr, 10, 64)
					listIntArr = append(listIntArr, roomIdInt)
				}

				if len(listIntArr) > 0 {
					ids := s.setRecInfoCache(ctx, listIntArr)
					key, expire := s.getRecPoolKey(ruleId)
					s.dao.SetRecPoolCache(ctx, key, xstr.JoinInts(ids), expire)
				}
				return nil
			})
		}

		// 条件规则
		if rule.ConfType == _condType {
			ruleId := rule.Id
			ruleStr := rule.Rules
			wg.Go(func() error {
				condConf, err := s.parseCondConf(ruleStr)
				if err != nil {
					log.Error("[genRecListJob]parseCondConf_err:%+v", err)
					return nil
				}
				if len(condConf.Condition) == 0 {
					log.Error("[genRecListJob]parseCondConf_empty")
					return nil
				}
				roomIds := s.genCondConfRoomList(ctx, condConf.Condition, condConf.Cond, ruleId)

				if len(roomIds) > 0 {
					ids := s.setRecInfoCache(ctx, roomIds)
					key, expire := s.getRecPoolKey(ruleId)
					s.dao.SetRecPoolCache(ctx, key, xstr.JoinInts(ids), expire)
				}
				return nil
			})
		}
	}

	// 超时检测
	select {
	case <-ctx.Done():
		{
			log.Info("[genRecListJob]job time out: %d", ttl)
			cancel()
		}
	}

	err = wg.Wait()
	if err != nil {
		log.Error("[genRecListJob]rule loop wait err:%+v", err)
	}
	log.Info("[genRecListJob]genRecListJob done")
}
