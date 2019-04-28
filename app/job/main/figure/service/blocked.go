package service

import (
	"context"
)

var (
	rageMap = map[int16]string{
		1: "S",
		2: "A",
		3: "B",
		4: "C",
	}
)

// BlockedKPIInfo .
func (s *Service) BlockedKPIInfo(c context.Context, mid int64, rage int16) (err error) {
	if len(rageMap[rage]) != 0 {
		s.figureDao.BlockedRage(c, mid, rage)
	} else {
		s.figureDao.BlockedRage(c, mid, int16(0))
	}
	return
}
