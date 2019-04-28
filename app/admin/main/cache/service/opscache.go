package service

import (
	"context"
	"time"
)

func (s *Service) loadOpsCache() {
	mcs, err := s.dao.OpsMemcaches(context.Background())
	if err == nil {
		s.opsMcs = mcs
	}
	rds, err := s.dao.OpsRediss(context.Background())
	if err == nil {
		s.opsRds = rds
	}
}

func (s *Service) loadOpsproc() {
	for {
		s.loadOpsCache()
		time.Sleep(time.Minute)
	}
}
