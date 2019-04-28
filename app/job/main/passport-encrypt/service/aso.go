package service

import (
	"context"

	"go-common/app/job/main/passport-encrypt/model"
)

func (s *Service) saveEncryptAccount(c context.Context, account *model.EncryptAccount) (err error) {
	var affect int64
	if affect, err = s.d.AddAsoAccount(c, account); err != nil || affect == 0 {
		return
	}
	return
}

func (s *Service) updateEncryptAccount(c context.Context, account *model.EncryptAccount) (err error) {
	var affect int64
	if affect, err = s.d.UpdateAsoAccount(c, account); err != nil || affect == 0 {
		return
	}
	return
}

func (s *Service) delEncryptAccount(c context.Context, mid int64) (err error) {
	var affect int64
	if affect, err = s.d.DelAsoAccount(c, mid); err != nil || affect == 0 {
		return
	}
	return
}
