package gorpc

import (
	rpcmodel "go-common/app/service/main/member/model/block"
	"go-common/library/net/rpc/context"
)

// BlockInfo is
func (r *RPC) BlockInfo(c context.Context, arg *rpcmodel.RPCArgInfo, res *rpcmodel.RPCResInfo) (err error) {
	var (
		blockInfos []*rpcmodel.BlockInfo
	)
	if blockInfos, err = r.block.Infos(c, []int64{arg.MID}); err != nil {
		return
	}
	if len(blockInfos) < 1 {
		res.Parse(r.block.DefaultUser(arg.MID))
	}
	res.Parse(blockInfos[0])
	return
}

// BlockBatchInfo is
func (r *RPC) BlockBatchInfo(c context.Context, arg *rpcmodel.RPCArgBatchInfo, res *[]*rpcmodel.RPCResInfo) (err error) {
	var (
		blockInfos []*rpcmodel.BlockInfo
	)
	if blockInfos, err = r.block.Infos(c, arg.MIDs); err != nil {
		return
	}
	for _, info := range blockInfos {
		r := &rpcmodel.RPCResInfo{}
		r.Parse(info)
		*res = append(*res, r)
	}
	return
}
