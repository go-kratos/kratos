package service

import (
	"context"

	"go-common/app/interface/main/spread/model"
)

// BangumiContent .
func (s *Service) BangumiContent(c context.Context, pn, ps int, typ int8, appkey string) (res model.BangumiResp, err error) {
	res, err = s.dao.BangumiContent(c, pn, ps, typ, appkey)
	return
}

// BangumiOff .
func (s *Service) BangumiOff(c context.Context, pn, ps int, typ int8, ts int64, appkey string) (res model.BangumiOffResp, err error) {
	res, err = s.dao.BangumiOff(c, pn, ps, typ, appkey, ts)
	return
}
