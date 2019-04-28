package service

import (
	"go-common/library/log"
)

//TransToCheckBack ..
func (s *Service) TransToCheckBack() {
	log.Info("deliveryNewVdieoToCms begin")
	s.dao.TransToCheckBack()
}

//TransToReview ...
func (s *Service) TransToReview() {
	log.Info("TransToReview begin")
	s.dao.TransToReview()
}
