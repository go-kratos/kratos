package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
)

type GetHandler struct {
}

func (handler *GetHandler) NeedCheckUid() bool {
	return true
}

func (handler *GetHandler) NeedTransactionMutex() bool {
	return false
}
func (handler *GetHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	platform, _ := params[0].(string)
	withMetal, _ := params[1].(int)
	if !model.IsValidPlatform(platform) {
		err = ecode.RequestErr
		return
	}
	r, err := ws.s.dao.GetMelonseedByCache(ws.c, uid)
	if err != nil {
		err = ecode.ServerErr
		return
	}
	if withMetal == 1 {
		metal, _ := ws.s.dao.GetMetal(ws.c, uid)
		v = model.GetMelonseedWithMetalResp(platform, r, metal)

	} else {
		v = model.GetMelonseedResp(platform, r)
	}
	return

}

func (s *Service) Get(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	handler := GetHandler{}
	return s.execByHandler(&handler, c, basicParam, uid, params...)
}
