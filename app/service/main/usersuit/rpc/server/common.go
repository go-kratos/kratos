package server

import (
	"go-common/app/service/main/usersuit/model"
	"go-common/library/net/rpc/context"
)

// PointFlag obtain new pendant noify.
func (r *RPC) PointFlag(c context.Context, arg *model.ArgMID, res *model.PointFlag) (err error) {
	var pf *model.PointFlag
	if pf, err = r.s.PointFlag(c, arg); err == nil && pf != nil {
		*res = *pf
	}
	return
}
