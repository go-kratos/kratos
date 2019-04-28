package service

import (
	"context"
	"time"

	"go-common/library/log"
)

func (s *Service) checkReadStatus() {
	for {
		var c = context.TODO()
		readSet, err := s.dao.ReadPingSet(c)
		if err != nil || len(readSet) == 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		now := time.Now().Unix()
		for _, read := range readSet {
			last, err := s.dao.ReadPing(c, read.Buvid, read.Aid)
			if err != nil {
				time.Sleep(time.Second)
				continue
			}
			if now-last < 30 {
				continue
			}
			if err = s.dao.DelReadPingSet(c, read); err != nil {
				time.Sleep(time.Second)
				continue
			}
			if last == 0 {
				log.Error("阅读心跳没取到:buvid(%s) aid(%d)", read.Buvid, read.Aid)
				continue
			}
			read.EndTime = last
			log.Info("上传用户阅读记录至数据中心: %+v", read)
			s.ReadInfoc(read.Aid, read.Mid, read.Buvid, read.IP, read.EndTime-read.StartTime, read.From)
		}
		return
	}
}

func (s *Service) checkReadStatusProc() {
	for {
		s.checkReadStatus()
		time.Sleep(time.Minute)
	}
}
