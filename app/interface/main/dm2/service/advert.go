package service

import (
	"context"

	"go-common/app/interface/main/dm2/model"
)

// DMAdvert dm advert.
func (s *Service) DMAdvert(c context.Context, arg *model.ADReq) (res *model.ADResp, err error) {
	ad, err := s.dao.DMAdvert(c, arg.Aid, arg.Oid, arg.Mid, arg.Build, arg.Buvid, arg.MobiApp, arg.ADExtra)
	if err != nil || ad == nil {
		return
	}
	res = ad.Convert(arg.ClientIP)
	return
}
