package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/mcn/model"
	accgrpc "go-common/app/service/main/account/api"
	memgrpc "go-common/app/service/main/member/api"
	blkmdl "go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// RecommendAdd .
func (s *Service) RecommendAdd(c context.Context, arg *model.RecommendUpReq) error {
	var (
		err       error
		ok        bool
		bindMids  []int64
		blockInfo *memgrpc.BlockInfoReply
		bi        *model.UpBaseInfo
		pi        *model.UpPlayInfo
		mbi       map[int64]*model.UpBaseInfo
		mpi       map[int64]*model.UpPlayInfo
		rpp       *model.McnUpRecommendPool
	)
	if blockInfo, err = s.memGRPC.BlockInfo(c, &memgrpc.MemberMidReq{Mid: arg.UpMid}); err != nil {
		return err
	}
	if blockInfo.BlockStatus > int32(blkmdl.BlockStatusFalse) {
		return fmt.Errorf("添加到推荐池的mid为(%d)的up主已经被封禁", arg.UpMid)
	}
	if bindMids, err = s.dao.McnUpBindMids(c, []int64{arg.UpMid}); err != nil {
		return err
	}
	if len(bindMids) > 0 {
		return fmt.Errorf("添加到推荐池的mid为(%s)的up主已经被绑定", xstr.JoinInts(bindMids))
	}
	if rpp, err = s.dao.McnUpRecommendMid(c, arg.UpMid); err != nil {
		return err
	}
	if rpp != nil && rpp.State != model.MCNUPRecommendStateDel {
		return ecode.MCNRecommendUpInPool
	}
	if mbi, err = s.dao.UpBaseInfoMap(c, []int64{arg.UpMid}); err != nil {
		return err
	}
	if mpi, err = s.dao.UpPlayInfoMap(c, []int64{arg.UpMid}); err != nil {
		return err
	}
	rp := &model.McnUpRecommendPool{UpMid: arg.UpMid}
	if bi, ok = mbi[arg.UpMid]; ok {
		rp.FansCount = bi.FansCount
		rp.ActiveTid = bi.ActiveTid
	}
	if pi, ok = mpi[arg.UpMid]; ok {
		rp.ArchiveCount = pi.ArticleCount
		rp.PlayCountAccumulate = pi.PlayCountAccumulate
		rp.PlayCountAverage = pi.PlayCountAverage
	}
	if _, err = s.dao.AddMcnUpRecommend(c, rp); err != nil {
		return err
	}
	s.worker.Add(func() {
		index := []interface{}{int8(model.MCNUPRecommendStateOff), arg.UpMid}
		content := map[string]interface{}{
			"up_mid":                arg.UpMid,
			"fans_count":            rp.FansCount,
			"archive_count":         rp.ArchiveCount,
			"play_count_accumulate": rp.PlayCountAccumulate,
			"play_count_average":    rp.PlayCountAverage,
			"active_tid":            rp.ActiveTid,
			"source":                model.MCNUPRecommendStateManual,
		}
		s.AddAuditLog(context.Background(), model.MCNRecommendLogBizID, int8(model.MCNUPRecommendActionAdd), model.MCNUPRecommendActionAdd.String(), arg.UID, arg.UserName, []int64{arg.UpMid}, index, content)
	})
	return nil
}

// RecommendOP .
func (s *Service) RecommendOP(c context.Context, arg *model.RecommendStateOpReq) error {
	var (
		err                                   error
		blockInfosReply                       *memgrpc.BlockBatchInfoReply
		blockMids, bindMids, banMids, recMids []int64
		mrp                                   map[int64]*model.McnUpRecommendPool
	)
	if len(arg.UpMids) == 0 {
		return ecode.MCNRecommendUpMidsIsEmpty
	}
	if arg.Action == model.MCNUPRecommendActionOn || arg.Action == model.MCNUPRecommendActionRestore {
		if blockInfosReply, err = s.memGRPC.BlockBatchInfo(c, &memgrpc.MemberMidsReq{Mids: arg.UpMids}); err != nil {
			return err
		}
		for _, v := range blockInfosReply.BlockInfos {
			if v.BlockStatus > int32(blkmdl.BlockStatusFalse) {
				blockMids = append(blockMids, v.MID)
			}
		}
		if len(blockMids) > 0 {
			return fmt.Errorf("推荐的mid为(%s)的up主已经被封禁", xstr.JoinInts(blockMids))
		}
		if bindMids, err = s.dao.McnUpBindMids(c, arg.UpMids); err != nil {
			return err
		}
		if len(bindMids) > 0 {
			return fmt.Errorf("推荐的mid为(%s)的up主已经被绑定", xstr.JoinInts(bindMids))
		}
	}
	state := arg.Action.GetState()
	if state == model.MCNUPRecommendStateUnKnown {
		return ecode.MCNRecommendUpStateFlowErr
	}
	if mrp, err = s.dao.McnUpRecommendMids(c, arg.UpMids); err != nil {
		return err
	}
	for _, upMids := range arg.UpMids {
		if rp, ok := mrp[upMids]; ok {
			if rp.State == model.MCNUPRecommendStateBan && arg.Action == model.MCNUPRecommendActionOn {
				banMids = append(banMids, upMids)
			}
			if rp.State == model.MCNUPRecommendStateOn && arg.Action == model.MCNUPRecommendActionRestore {
				recMids = append(recMids, upMids)
			}
		}
	}
	if len(banMids) > 0 {
		return fmt.Errorf("推荐的mid为(%s)的up主已被禁止推荐,不能推荐", xstr.JoinInts(banMids))
	}
	if len(recMids) > 0 {
		return fmt.Errorf("推荐的mid为(%s)的up主已被推荐,不需要恢复", xstr.JoinInts(recMids))
	}
	if _, err = s.dao.UpMcnUpsRecommendOP(c, arg.UpMids, state); err != nil {
		return err
	}
	for _, mid := range arg.UpMids {
		s.worker.Add(func() {
			index := []interface{}{int8(state), mid}
			content := map[string]interface{}{
				"up_mid": mid,
			}
			s.AddAuditLog(context.Background(), model.MCNRecommendLogBizID, int8(state), arg.Action.String(), arg.UID, arg.UserName, []int64{mid}, index, content)
		})
	}
	return nil
}

// RecommendList .
func (s *Service) RecommendList(c context.Context, arg *model.MCNUPRecommendReq) (res *model.McnUpRecommendListReply, err error) {
	var (
		mids, tids []int64
		tpNames    map[int64]string
		accsReply  *accgrpc.InfosReply
	)
	res = new(model.McnUpRecommendListReply)
	res.Page = arg.Page
	if res.TotalCount, err = s.dao.McnUpRecommendTotal(c, arg); err != nil {
		return
	}
	if res.TotalCount <= 0 {
		return
	}
	if res.List, err = s.dao.McnUpRecommends(c, arg); err != nil {
		return
	}
	if len(res.List) <= 0 {
		return
	}
	for _, v := range res.List {
		mids = append(mids, v.UpMid)
		tids = append(tids, int64(v.ActiveTid))
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) err(%v)", mids, err)
		err = nil
	}
	tpNames = s.videoup.GetTidName(tids)
	infos := accsReply.Infos
	for _, v := range res.List {
		if info, ok := infos[v.UpMid]; ok {
			v.UpName = info.Name
		}
		if tyName, ok := tpNames[int64(v.ActiveTid)]; ok {
			v.TpName = tyName
		} else {
			v.TpName = model.DefaultTyName
		}
	}
	return
}
