package service

import (
	"context"
	"time"

	"go-common/library/log"
	"go-common/library/stat/prom"
)

const _businessArchive = 3
const _delLen = 1000

func (s *Service) shouldDelete() bool {
	now := time.Now()
	return now.Hour() >= s.c.Job.DeleteStartHour && now.Hour() < s.c.Job.DeleteEndHour
}

func (s *Service) deleteproc() {
	for {
		now := time.Now()
		if !s.shouldDelete() {
			time.Sleep(time.Minute)
			continue
		}
		if ok, err := s.dao.DelLock(context.Background()); err != nil {
			time.Sleep(time.Second)
			continue
		} else if !ok {
			log.Info("not get lock wait.")
			time.Sleep(time.Hour * 6)
			continue
		}
		log.Info("start clean db")
		bs, err := s.dao.Businesses(context.Background())
		if err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, b := range bs {
			if b.TTL <= 0 {
				continue
			}
			endTime := time.Unix(now.Unix()-b.TTL, 0)
			startTime, err := s.dao.EarlyHistory(context.Background(), b.ID)
			if err != nil {
				continue
			}
			log.Info("start clean business %s start:%v end: %v", b.Name, startTime, endTime)
			var count int64
			for startTime.Before(endTime) {
				if !s.shouldDelete() {
					log.Info("%s not delete time.", b.Name)
					break
				}
				partTime := startTime.Add(time.Duration(s.c.Job.DeleteStep))
				rows, err := s.dao.DeleteHistories(context.Background(), b.ID, startTime, partTime)
				prom.BusinessInfoCount.Add("del-"+b.Name, rows)
				if err != nil {
					time.Sleep(time.Second)
					continue
				}
				count += rows
				// 删除完这个时间段的数据后再删除下个时间段
				if rows == 0 {
					startTime = partTime
				}
			}
			log.Info("end clean business %s, rows: %v", b.Name, count)
		}
		log.Info("end clean db")
		time.Sleep(time.Hour * 6)
	}
}
