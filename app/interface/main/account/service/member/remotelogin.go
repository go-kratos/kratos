package member

import (
	"context"

	secmodel "go-common/app/service/main/secure/model"
	"go-common/library/log"
)

// Status query user's remote login status.
func (s *Service) Status(c context.Context, mid int64, uuid string) (res *secmodel.Msg, err error) {
	arg := &secmodel.ArgSecure{Mid: mid, UUID: uuid}
	if res, err = s.secureRPC.Status(c, arg); err != nil {
		log.Error("s.secureRPC.Status(mid:%d) error (%v)", mid, err)
		return
	}
	return
}
