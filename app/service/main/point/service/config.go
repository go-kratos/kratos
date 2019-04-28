package service

import (
	"context"
	"math"
	"strconv"

	"go-common/app/service/main/point/model"
	"go-common/library/log"
)

const (
	_defcursor  = int(^uint(0) >> 1)
	_defps      = 20
	_defpn      = 1
	_timeFormat = "2006-01-02 15:04:05"
)

// Config get point config.
func (s *Service) Config(c context.Context, changeType int, mid int64, bp float64) (point int64, err error) {
	var (
		rate  int64
		times int64
		ok    bool
	)
	if bp == 0 {
		return
	}
	if rate, ok = s.c.Property.PointGetRule[strconv.Itoa(changeType)]; !ok {
		return
	}
	point = int64(math.Ceil(float64(rate) * bp))
	if times, err = s.activityGiveTimes(c, mid, changeType, point); err != nil {
		log.Error("%+v", err)
		return
	}
	point += times * model.ActivityGivePoint
	return
}

//AllConfig all point config
func (s *Service) AllConfig(c context.Context) map[string]int64 {
	return s.c.Property.PointGetRule
}
