package grpc

import (
	"context"

	"go-common/app/service/main/member/api"
)

// BlockInfo 查询封禁信息
func (s *MemberServer) BlockInfo(ctx context.Context, req *api.MemberMidReq) (*api.BlockInfoReply, error) {
	res, err := s.blockSvr.Infos(ctx, []int64{req.Mid})
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return api.FromBlockInfo(s.blockSvr.DefaultUser(req.Mid)), nil
	}
	blockInfoReply := api.FromBlockInfo(res[0])
	return blockInfoReply, nil
}

// BlockBatchInfo 批量查询封禁信息
func (s *MemberServer) BlockBatchInfo(ctx context.Context, req *api.MemberMidsReq) (*api.BlockBatchInfoReply, error) {
	res, err := s.blockSvr.Infos(ctx, req.Mids)
	if err != nil {
		return nil, err
	}
	blockInfos := make([]*api.BlockInfoReply, 0, len(res))
	for i := range res {
		blockInfos = append(blockInfos, api.FromBlockInfo(res[i]))
	}
	blockInfosReply := &api.BlockBatchInfoReply{
		BlockInfos: blockInfos,
	}
	return blockInfosReply, nil
}

// BlockBatchDetail 批量查询封禁信息
func (s *MemberServer) BlockBatchDetail(ctx context.Context, req *api.MemberMidsReq) (reply *api.BlockBatchDetailReply, err error) {
	res, err := s.blockSvr.UserDetails(ctx, req.Mids)
	if err != nil {
		return nil, err
	}
	userDetails := make(map[int64]*api.BlockDetailReply, len(res))
	for mid := range res {
		userDetails[mid] = api.FromBlockUserDetail(res[mid])
	}
	blockUserDetailsReply := &api.BlockBatchDetailReply{
		BlockDetails: userDetails,
	}
	return blockUserDetailsReply, nil
}
