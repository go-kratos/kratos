package vip

import (
	"context"

	"go-common/app/interface/main/account/model"
	vipmod "go-common/app/service/main/vip/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// CodeOpen .
func (s *Service) CodeOpen(c context.Context, mid int64, code, token, verify string) (codeInfo *model.ResourceCode, err error) {
	var (
		codeResp *model.ResourceCodeResq
	)
	if codeResp, err = s.vipDao.CodeOpen(c, mid, code, token, verify); err != nil {
		err = errors.WithStack(err)
		return
	}
	codeInfo = codeResp.Data
	return
}

//CodeOpeneds sel code opened
func (s *Service) CodeOpeneds(c context.Context, arg *model.CodeInfoReq, ip string) (resp []*vipmod.CodeInfoResp, err error) {

	if err = s.checkIP(arg.Appkey, ip); err != nil {
		err = errors.WithStack(err)
		return
	}
	if resp, err = s.vipDao.CodeOpeneds(c, arg, ip); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func (s *Service) checkIP(appkey, ip string) (err error) {
	var (
		strings []string
		ok      bool
	)
	if strings, ok = s.c.Vipproperty.CodeOpenwhiteIPMap[appkey]; !ok {
		log.Error("checkIP s.c.Vipproperty.CodeOpenwhiteIPMap empty(%s)", appkey)
		err = ecode.VipWhiteIPListErr
		return
	}
	for _, v := range strings {
		if v == ip {
			return
		}
	}
	log.Error("checkIP fail(%s, %s)", appkey, ip)
	err = ecode.VipWhiteIPListErr
	return
}

// CodeVerify .
func (s *Service) CodeVerify(c context.Context) (token *model.Token, err error) {
	if token, err = s.vipDao.CodeVerify(c); err != nil {
		err = errors.WithStack(err)
	}
	return
}

//UseBatch use resource batch
func (s *Service) UseBatch(c context.Context, arg *vipmod.ArgUseBatch) (err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if err = s.checkIP(arg.Appkey, ip); err != nil {
		err = errors.WithStack(err)
		return
	}
	if err = s.vipRPC.ResourceBatchOpenVip(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}
