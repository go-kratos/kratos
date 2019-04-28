package service

import (
	"time"

	"go-common/app/interface/main/captcha/conf"
	xtime "go-common/library/time"
)

var (
	_defaultBusiness = &conf.Business{
		BusinessID: "default",
		LenStart:   4,
		LenEnd:     4,
		Width:      100,
		Length:     50,
		TTL:        xtime.Duration(300 * time.Second),
	}
)

// LookUp look up business services.
func (s *Service) LookUp(bid string) (business *conf.Business) {
	for _, b := range s.conf.Business {
		if b.BusinessID == bid {
			return b
		}
	}
	return _defaultBusiness
}
