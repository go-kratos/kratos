package usersuit

import (
	"context"

	usmdl "go-common/app/service/main/usersuit/model"
)

// PointFlag .
func (s *Service) PointFlag(c context.Context, arg *usmdl.ArgMID) (res *usmdl.PointFlag, err error) {
	return s.usRPC.PointFlag(c, arg)
}
