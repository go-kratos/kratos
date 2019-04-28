package ugc

import (
	"context"

	arccli "go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

// viewCache distinguishes the archive's license status,
// if it's ok, we call it's RPC and save the result in MC cache to accelerate the view page loading
func (s *Service) viewCache(aid int64) {
	if aid == 0 {
		return
	}
	var (
		c   = context.Background()
		err error
		arg = &arccli.ViewRequest{Aid: aid}
		v   *arccli.ViewReply
	)
	if v, err = s.arcClient.View(c, arg); err != nil {
		if ecode.Cause(err) == ecode.NothingFound {
			log.Warn("s.arcRPC.View3(%v) error(%v)", arg, err)
			err = nil
			return
		}
		log.Error("s.arcRPC.View3(%v) error(%v)", arg, err)
		return
	}
	if err = s.arcDao.UpArcCache(c, v.Arc); err != nil {
		log.Error("viewCache %+v", err)
		return
	}
	if s.arcDao.UpViewCache(c, v); err != nil {
		log.Error("viewCache %+v", err)
	}
}
