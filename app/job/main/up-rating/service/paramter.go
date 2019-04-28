package service

import (
	"context"

	"go-common/app/job/main/up-rating/model"
)

func (s *Service) getAllParamter(c context.Context) (rp *model.RatingParameter, err error) {
	paramters, err := s.dao.GetAllParamter(c)
	if err != nil {
		return
	}
	rp = &model.RatingParameter{
		WDP:      paramters["wdp"],
		WDC:      paramters["wdc"],
		WDV:      paramters["wdv"],
		WMDV:     paramters["wmdv"],
		WCS:      paramters["wcs"],
		WCSR:     paramters["wcsr"],
		WMAAFans: paramters["wmaafans"],
		WMAHFans: paramters["wmahfans"],
		WIS:      paramters["wis"],
		WISR:     paramters["wisr"],
		HBASE:    paramters["hbase"],
		HR:       paramters["hr"],
		HV:       paramters["hv"],
		HVM:      paramters["hvm"],
		HL:       paramters["hl"],
		HLM:      paramters["hlm"],
	}
	return
}
