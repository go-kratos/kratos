package service

import (
	"context"
	"go-common/app/service/live/wallet/dao"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

type GetTidHandler struct {
}

func (handler *GetTidHandler) NeedCheckUid() bool {
	return false
}
func (handler *GetTidHandler) NeedTransactionMutex() bool {
	return false
}
func (handler *GetTidHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	serviceType, _ := params[0].(int32)
	if !model.IsValidServiceType(serviceType) {
		err = ecode.RequestErr
		return
	}

	callParams, _ := params[1].(string)
	if callParams == "" {
		err = ecode.RequestErr
		return
	}
	log.Info("getTid info : type:%d,callParams:%s", serviceType, callParams)
	tid := dao.GetTid(model.ServiceType(serviceType), callParams)

	v = model.GetTidResp(tid)

	return

}

func (s *Service) GetTid(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	handler := GetTidHandler{}
	return s.execByHandler(&handler, c, basicParam, uid, params...)
}
