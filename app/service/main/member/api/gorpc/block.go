package gorpc

import (
	"context"

	rpcmodel "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
)

const (
	_blockInfo      = "RPC.BlockInfo"
	_blockBatchInfo = "RPC.BlockBatchInfo"
)

// BlockInfo is
func (s *Service) BlockInfo(c context.Context, arg *rpcmodel.RPCArgInfo) (res *rpcmodel.RPCResInfo, err error) {
	res = new(rpcmodel.RPCResInfo)
	err = s.client.Call(c, _blockInfo, arg, res)
	return
}

// BlockBatchInfo len(mids) <= 50
func (s *Service) BlockBatchInfo(c context.Context, arg *rpcmodel.RPCArgBatchInfo) (res []*rpcmodel.RPCResInfo, err error) {
	if len(arg.MIDs) == 0 {
		return
	}
	if len(arg.MIDs) > 50 {
		err = ecode.RequestErr
		return
	}
	err = s.client.Call(c, _blockBatchInfo, arg, &res)
	return
}
