package service

import (
	"context"

	"go-common/app/interface/main/mcn/model/mcnmodel"
)

//CmdReloadRankCache reload cache
func (s *Service) CmdReloadRankCache(c context.Context, arg *mcnmodel.CmdReloadRank) (res *mcnmodel.CommonReply, err error) {
	err = s.mcndao.ReloadRank(arg.SignID)
	return
}
