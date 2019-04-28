package service

import (
	"context"
	"time"

	"go-common/app/job/main/tag/model"
)

func (s *Service) businessCacheproc() {
	for {
		time.Sleep(time.Minute * 10)
		s.businessCaches()
	}
}

// businessCache .
func (s *Service) businessCaches() (err error) {
	business, err := s.dao.Business(context.TODO(), model.BusinessStateNormal)
	if err != nil {
		return
	}
	s.businessCache.Store(business)
	return
}
