package service

import (
	"context"

	"go-common/app/interface/main/dm/model"
	"go-common/library/ecode"
)

func (s *Service) subject(c context.Context, tp int32, oid int64) (sub *model.Subject, err error) {
	if sub, err = s.dao.Subject(c, tp, oid); err != nil {
		return
	}
	if sub == nil {
		err = ecode.NothingFound
	}
	return
}
