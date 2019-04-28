package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/up-rating/model"
)

// InsertParameter insert parameter
func (s *Service) InsertParameter(c context.Context, name, remark string, value int) (err error) {
	_, err = s.dao.InsertParameter(c, fmt.Sprintf("('%s', %d, '%s')", name, value, remark))
	return
}

func (s *Service) getTypeScore(c context.Context, ctype int64) (score int64, err error) {
	params, err := s.getAllParameter(c)
	if err != nil {
		return
	}

	switch ctype {
	case 0:
		score = params.WCSR + params.HR + params.WISR
	case 1:
		score = params.WCSR
	case 2:
		score = params.WISR
	case 3:
		score = params.HR
	}
	return
}

func (s *Service) getAllParameter(c context.Context) (rp *model.RatingParameter, err error) {
	parameters, err := s.dao.GetAllParameter(c)
	if err != nil {
		return
	}
	rp = &model.RatingParameter{
		WDP:      parameters["wdp"],
		WDC:      parameters["wdc"],
		WDV:      parameters["wdv"],
		WMDV:     parameters["wmdv"],
		WCS:      parameters["wcs"],
		WCSR:     parameters["wcsr"],
		WMAAFans: parameters["wmaafans"],
		WMAHFans: parameters["wmahfans"],
		WIS:      parameters["wis"],
		WISR:     parameters["wisr"],
		HBASE:    parameters["hbase"],
		HR:       parameters["hr"],
		HV:       parameters["hv"],
		HVM:      parameters["hvm"],
		HL:       parameters["hl"],
		HLM:      parameters["hlm"],
	}
	return
}
