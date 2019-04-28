package goblin

import (
	"context"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"
)

// VerUpdate .
func (s *Service) VerUpdate(c context.Context, ver *model.VerUpdate) (result *model.HTTPData, errCode ecode.Codes, err error) {
	result, errCode, err = s.dao.VerUpdate(c, ver)
	return
}
