package upcrmservice

import (
	"context"
	"go-common/app/admin/main/up/model"
	"time"
)

//CommandRefreshUpRank refresh up rank
func (s *Service) CommandRefreshUpRank(con context.Context, arg *model.CommandCommonArg) (result model.CommandCommonResult, err error) {
	s.RefreshCache(time.Now())
	return
}
