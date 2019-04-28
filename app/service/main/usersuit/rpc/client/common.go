package client

import (
	"context"

	"go-common/app/service/main/usersuit/model"
)

const (
	_pointFlag = "RPC.PointFlag"
)

// PointFlag obtain new pendant noify.
func (s *Service2) PointFlag(c context.Context, arg *model.ArgMID) (res *model.PointFlag, err error) {
	res = new(model.PointFlag)
	err = s.client.Call(c, _pointFlag, arg, res)
	return
}
