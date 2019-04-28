package service

import (
	"context"
	"time"

	"go-common/library/conf/env"
	"go-common/library/log"
)

// taskAddMonitor add monitor task
func (s *Service) taskAddMonitor() {
	if err := s.AddMonitor(context.TODO()); err != nil {
		log.Error("task.taskAddMonitor error(%v)", err)
	}
}

// taskAddCache add cache task
func (s *Service) taskAddCache() {
	if err := s.RanksCache(context.Background()); err != nil {
		log.Error("task.RanksCache error(%v)", err)
	}
	if err := s.AppsCache(context.Background()); err != nil {
		log.Error("task.AppsCache error(%v)", err)
	}
}

// taskRankWechatReport send rank report to wechat group task
func (s *Service) taskRankWechatReport() {
	if env.DeployEnv != env.DeployEnvProd || time.Now().Weekday() == time.Sunday || time.Now().Weekday() == time.Saturday {
		return
	}
	if err := s.RankWechatReport(context.TODO()); err != nil {
		log.Error("task.taskRankWechatReport error(%v)", err)
	}
}

// taskWeeklyWechatReport send Weekly report and reset redis every Friday 19:00
func (s *Service) taskWeeklyWechatReport() {
	if env.DeployEnv == env.DeployEnvProd && time.Now().Weekday() == time.Friday {
		if err := s.SummaryWechatReport(context.TODO()); err != nil {
			log.Error("task.taskWeeklyWechatReport error(%v)", err)
		}
		if err := s.dao.SetAppCovCache(context.TODO()); err != nil {
			log.Error("task.taskWeeklyWechatReport error(%v)", err)
		}
	}
}
