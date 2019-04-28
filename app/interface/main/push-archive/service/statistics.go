package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/interface/main/push-archive/dao"
	"go-common/app/interface/main/push-archive/model"
	"go-common/library/log"
)

// 切割需要存储的统计数据
func (s *Service) formStatisticsProc(aid int64, group string, included []int64, excluded *[]int64) {
	params := model.NewBatchParam(map[string]interface{}{
		"aid":         aid,
		"group":       group,
		"type":        model.StatisticsPush,
		"createdTime": time.Now(),
	}, nil)
	dao.Batch(&included, 1000, 2, params, s.formPushStatistic)

	params.Params["type"] = model.StatisticsUnpush
	dao.Batch(excluded, 1000, 2, params, s.formPushStatistic)
}

// 组建统计对象
func (s *Service) formPushStatistic(fans *[]int64, params map[string]interface{}) (err error) {
	ln := len(*fans)
	b, err := json.Marshal(*fans)
	if err != nil {
		log.Error("formStatistic json.Marshal error(%v) fans(%v) params(%v)", err, fans, params)
		return
	}
	ps := &model.PushStatistic{
		Aid:         params["aid"].(int64),
		Group:       params["group"].(string),
		Type:        params["type"].(int),
		Mids:        string(b),
		MidsCounter: ln,
		CTime:       params["createdTime"].(time.Time),
	}
	if err = s.dao.AddStatisticsCache(context.TODO(), ps); err != nil {
		log.Error("formPushStatistic s.dao.AddStatisticsCache error(%v), pushstatistic(%v)", err, ps)
		return
	}
	return
}

// 每日定时清除推送的统计数据,只保留最近几天的数据
func (s *Service) clearStatisticsProc() {
	for {
		// 到指定时间
		clearTime, err := s.getTodayTime(s.c.ArcPush.PushStatisticsClearTime)
		if err != nil {
			log.Error("clearStatisticsProc getTodayTime(%s) error(%v)", s.c.ArcPush.PushStatisticsClearTime, err)
			continue
		}
		dur := clearTime.Unix() - time.Now().Unix()
		if dur < 0 || dur > 60 {
			time.Sleep(time.Second * 50)
			continue
		}
		// 需要删除数据的最大时间
		time.Sleep(time.Second * 60)
		deadline, err := s.getDeadline()
		if err != nil {
			log.Error("clearStatisticsProc getDeadline error(%v)", err)
			continue
		}
		log.Info("start to clear statistics before deadline(%s)", deadline.Format("2006-01-02 15:04:05"))
		var min, max, mid int64
		for i := 0; i < 3; i++ {
			if min == 0 && max == 0 {
				min, max, err = s.dao.GetStatisticsIDRange(context.TODO(), deadline)
				mid = min
			}
			for err == nil && mid < max {
				min = mid
				mid = min + 5000
				if mid > max {
					mid = max
				}
				_, err = s.dao.DelStatisticsByID(context.TODO(), min, mid)
				time.Sleep(time.Second)
			}
			if err == nil {
				log.Info("success end clear statistics before deadline(%s)", deadline.Format("2006-01-02 15:04:05"))
				break
			}
		}
		if err != nil {
			s.dao.WechatMessage(fmt.Sprintf("clearStatisticsProc: push-archive failed to clear expired(%s) push_statistics, error(%v)", deadline.Format("2006-01-02 15:04:05"), err))
		}
	}
}

// 获取需要删除数据的 最大时间
func (s *Service) getDeadline() (deadline time.Time, err error) {
	dd := "00:00:00"
	today, err := s.getTodayTime(dd)
	if err != nil {
		log.Error("clearStatisticsProc getTodayTime(%s) error(%v)", dd, err)
		return
	}

	deadline = today.AddDate(0, 0, -1*s.c.ArcPush.PushStatisticsKeepDays+1)
	return
}

// 统计数据落库
func (s *Service) saveStatisticsProc() {
	defer s.wg.Done()
	for {
		select {
		case _, ok := <-s.CloseCh:
			if !ok {
				log.Info("CloseCh is closed, close the saveStatisticsProc")
				return
			}
		default:
		}

		ps, err := s.dao.GetStatisticsCache(context.TODO())
		if err != nil {
			log.Error("saveStatisticsProc s.dao.GetStatisticsCache error(%v)", err)
			time.Sleep(time.Millisecond * 100)
			continue
		}
		if ps == nil {
			time.Sleep(time.Millisecond * 100)
			continue
		}

		if _, err := s.dao.SetStatistics(context.TODO(), ps); err != nil {
			log.Error("saveStatisticsProc s.dao.SetStatistics error(%v) pushstatistic(%v)", err, ps)
			s.dao.AddStatisticsCache(context.TODO(), ps)
		}
		time.Sleep(time.Millisecond * 100)
	}
}
