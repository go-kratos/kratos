package client

import (
	"context"

	"go-common/app/interface/main/dm2/model"
)

const _mask = "RPC.Mask"

// Mask get mask
func (s *Service) Mask(c context.Context, arg *model.ArgMask) (res *model.Mask, err error) {
	err = s.client.Call(c, _mask, arg, &res)
	return
}
