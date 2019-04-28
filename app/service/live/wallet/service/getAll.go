package service

import (
	"context"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
)

type GetAllHandler struct {
}

func (handler *GetAllHandler) NeedCheckUid() bool {
	return true
}

func (handler *GetAllHandler) NeedTransactionMutex() bool {
	return false
}
func (handler *GetAllHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	platform, _ := params[0].(string)
	withMetal, _ := params[1].(int)
	if !model.IsValidPlatform(platform) {
		err = ecode.RequestErr
		return
	}
	r, err := ws.s.dao.GetDetailByCache(ws.c, uid)
	if err != nil {
		err = ecode.ServerErr
		return
	}

	if withMetal == 1 {
		metal, _ := ws.s.dao.GetMetal(ws.c, uid)
		v = model.GetDetailWithMetalResp(platform, r, metal)
	} else {
		v = model.GetDetailResp(platform, r)
	}
	return

}

func (s *Service) GetAll(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	handler := GetAllHandler{}
	return s.execByHandler(&handler, c, basicParam, uid, params...)
}
