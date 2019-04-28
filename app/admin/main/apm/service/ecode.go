package service

import (
	"context"
	"go-common/app/admin/main/apm/model/ecode"
	"go-common/library/log"
)

// GetCodes ...
func (s *Service) GetCodes(c context.Context, Interval1, Interval2 string) (data []*codes.Codes, err error) {
	data, err = s.dao.GetCodes(c, Interval1, Interval2)
	if err != nil {
		log.Error("service GetCodes error(%v)", err)
	}
	return
}
