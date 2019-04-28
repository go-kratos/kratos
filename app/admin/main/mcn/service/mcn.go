package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"go-common/app/admin/main/mcn/dao/up"
	"go-common/app/admin/main/mcn/model"
	accgrpc "go-common/app/service/main/account/api"
	memgrpc "go-common/app/service/main/member/api"
	blkmdl "go-common/app/service/main/member/model/block"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// McnSignEntry .
func (s *Service) McnSignEntry(c context.Context, arg *model.MCNSignEntryReq) error {
	var (
		err           error
		count, lastID int64
		stime, etime  xtime.Time
		blockInfo     *memgrpc.BlockInfoReply
		tx            *xsql.Tx
		now           = time.Now()
		date          = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	)
	if tx, err = s.dao.BeginTran(c); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if blockInfo, err = s.memGRPC.BlockInfo(c, &memgrpc.MemberMidReq{Mid: arg.MCNMID}); err != nil {
		return err
	}
	if blockInfo.BlockStatus > int32(blkmdl.BlockStatusFalse) {
		return ecode.MCNSignIsBlocked
	}
	if count, err = s.dao.McnSignNoOKState(c, arg.MCNMID); err != nil {
		return err
	}
	if count > 0 {
		return ecode.MCNSignNoOkState
	}
	if stime, etime, err = arg.ParseTime(); err != nil {
		return err
	}
	if etime.Time().Before(date) || etime.Time().Equal(date) {
		return ecode.MCNSignEtimeNLEQNowTime
	}
	if count, err = s.dao.McnSignCountUQTime(c, arg.MCNMID, stime, etime); err != nil {
		return err
	}
	if count > 0 {
		return ecode.MCNSignCycleNotUQErr
	}
	arg.AttrPermitSet()
	if lastID, err = s.dao.TxAddMcnSignEntry(tx, arg.MCNMID, arg.BeginDate, arg.EndDate, arg.Permission); err != nil {
		return err
	}
	for _, v := range arg.SignPayInfo {
		if _, err = s.dao.TxAddMcnSignPay(tx, arg.MCNMID, lastID, v.PayValue, v.DueDate, ""); err != nil {
			return err
		}
	}
	s.worker.Add(func() {
		index := []interface{}{int8(model.MCNSignStateNoApply), lastID}
		content := map[string]interface{}{
			"sign_id":    lastID,
			"mcn_mid":    arg.MCNMID,
			"state":      int8(model.MCNSignStateNoApply),
			"begin_date": arg.BeginDate,
			"end_date":   arg.EndDate,
			"permission": arg.Permission,
		}
		s.AddAuditLog(context.Background(), model.MCNLogBizID, int8(model.MCNSignActionEntry), model.MCNSignActionEntry.String(), arg.UID, arg.UserName, []int64{arg.MCNMID}, index, content)
	})
	s.worker.Add(func() {
		index := []interface{}{arg.MCNMID}
		content := map[string]interface{}{
			"sign_id":       lastID,
			"mcn_mid":       arg.MCNMID,
			"sign_pay_info": arg.SignPayInfo,
		}
		s.AddAuditLog(context.Background(), model.MCNPayDateLogBizID, int8(model.MCNSignCycleActionAdd), model.MCNSignCycleActionAdd.String(), arg.UID, arg.UserName, []int64{lastID}, index, content)
	})
	s.worker.Add(func() {
		s.dao.DelMcnSignCache(context.Background(), arg.MCNMID)
	})
	return nil
}

// McnSignList .
func (s *Service) McnSignList(c context.Context, arg *model.MCNSignStateReq) (res *model.MCNSignListReply, err error) {
	var (
		count         int64
		signIDs, mids []int64
		accsReply     *accgrpc.InfosReply
		accInfos      map[int64]*accgrpc.Info
		mcns          []*model.MCNSignInfoReply
		sm            map[int64][]*model.SignPayInfoReply
	)
	res = new(model.MCNSignListReply)
	res.Page = arg.Page
	if count, err = s.dao.McnSignTotal(c, arg); err != nil {
		return res, err
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if mcns, err = s.dao.McnSigns(c, arg); err != nil {
		return res, err
	}
	for _, v := range mcns {
		v.AttrPermitVal()
		signIDs = append(signIDs, v.SignID)
		mids = append(mids, v.McnMid)
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) error(%+v)", &accgrpc.MidsReq{Mids: mids}, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	if sm, err = s.dao.McnSignPayMap(c, signIDs); err != nil {
		return res, err
	}
	for _, v := range mcns {
		if info, ok := accInfos[v.McnMid]; ok {
			v.McnName = info.Name
		}
		v.ContractLink = model.BuildBfsURL(v.ContractLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
		v.CompanyLicenseLink = model.BuildBfsURL(v.CompanyLicenseLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
		if sp, ok := sm[v.SignID]; ok {
			v.SignPayInfo = sp
		}
		sort.Slice(v.SignPayInfo, func(i int, j int) bool {
			return v.SignPayInfo[i].DueDate < v.SignPayInfo[j].DueDate
		})
	}

	res.List = mcns
	return res, nil
}

// McnSignOP .
func (s *Service) McnSignOP(c context.Context, arg *model.MCNSignStateOpReq) error {
	var (
		err      error
		ok       bool
		accReply *accgrpc.InfoReply
		m        *model.MCNSignInfoReply
		payInfo  []*model.SignPayInfoReply
		sm       map[int64][]*model.SignPayInfoReply
		now      = xtime.Time(time.Now().Unix())
		state    model.MCNSignState
	)
	if !arg.Action.NotRightAction() {
		return ecode.MCNSignUnknownReviewErr
	}
	if m, err = s.dao.McnSign(c, arg.SignID); err != nil {
		return err
	}
	if m == nil {
		return ecode.MCNCSignUnknownInfoErr
	}
	state = arg.Action.GetState(m.State)
	if state == model.MCNSignStateUnKnown {
		log.Warn("mcn_sign action(%d) old state(%d) to new err state(-1)", arg.Action, m.State)
		return ecode.MCNSignStateFlowErr
	}
	if m.State.IsOnReviewState(arg.Action) {
		return ecode.MCNSignOnlyReviewOpErr
	}
	if now < m.BeginDate && arg.Action == model.MCNSignActionPass {
		state = model.MCNSignStateOnPreOpen
	}
	if arg.Action != model.MCNSignActionReject {
		now = 0
		arg.RejectReason = ""
	}
	if _, err = s.dao.UpMcnSignOP(c, arg.SignID, int8(state), now, arg.RejectReason); err != nil {
		return err
	}
	if accReply, err = s.accGRPC.Info3(c, &accgrpc.MidReq{Mid: m.McnMid}); err != nil {
		log.Error("s.accGRPC.Info3(%+v) error(%+v)", &accgrpc.MidReq{Mid: m.McnMid}, err)
		err = nil
	}
	if sm, err = s.dao.McnSignPayMap(c, []int64{arg.SignID}); err != nil {
		log.Error("s.dao.McnSignPayMap(%+v) error(%+v)", []int64{arg.SignID}, err)
		err = nil
	}
	if payInfo, ok = sm[arg.SignID]; !ok {
		payInfo = nil
	}
	s.worker.Add(func() {
		var name string
		if accReply.Info != nil {
			name = accReply.Info.Name
		}
		index := []interface{}{int8(state), arg.SignID, time.Now().Unix(), name}
		content := map[string]interface{}{
			"sign_id":       arg.SignID,
			"mcn_mid":       m.McnMid,
			"mcn_name":      name,
			"begin_date":    m.BeginDate,
			"end_date":      m.EndDate,
			"reject_reason": arg.RejectReason,
			"reject_time":   now,
			"contract_link": m.ContractLink,
			"sign_pay_info": payInfo,
			"permission":    m.Permission,
		}
		s.AddAuditLog(context.Background(), model.MCNLogBizID, int8(arg.Action), arg.Action.String(), arg.UID, arg.UserName, []int64{m.McnMid}, index, content)
	})
	s.worker.Add(func() {
		s.sendMsg(context.Background(), &model.ArgMsg{MSGType: arg.Action.GetmsgType(state), MIDs: []int64{m.McnMid}, Reason: arg.RejectReason})
	})
	s.worker.Add(func() {
		s.dao.DelMcnSignCache(context.Background(), m.McnMid)
	})
	return nil
}

// McnUPReviewList .
func (s *Service) McnUPReviewList(c context.Context, arg *model.MCNUPStateReq) (res *model.MCNUPReviewListReply, err error) {
	var (
		count     int64
		mids      []int64
		ups       []*model.MCNUPInfoReply
		mbi       map[int64]*model.UpBaseInfo
		accsReply *accgrpc.InfosReply
		accInfos  map[int64]*accgrpc.Info
	)
	res = new(model.MCNUPReviewListReply)
	res.Page = arg.Page
	if count, err = s.dao.McnUpTotal(c, arg); err != nil {
		return res, err
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if ups, err = s.dao.McnUps(c, arg); err != nil {
		return res, err
	}
	for _, v := range ups {
		v.AttrPermitVal()
		mids = append(mids, v.McnMid)
		mids = append(mids, v.UpMid)
	}
	if mbi, err = s.dao.UpBaseInfoMap(c, mids); err != nil {
		return res, err
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) error(%+v)", &accgrpc.MidsReq{Mids: mids}, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	for _, v := range ups {
		if up, ok := mbi[v.McnMid]; ok {
			v.FansCount = up.FansCount
			v.ActiveTid = up.ActiveTid
		}
		if info, ok := accInfos[v.McnMid]; ok {
			v.McnName = info.Name
		}
		if info, ok := accInfos[v.UpMid]; ok {
			v.UpName = info.Name
		}
		v.UpAuthLink = model.BuildBfsURL(v.UpAuthLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
		v.ContractLink = model.BuildBfsURL(v.ContractLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
	}
	res.List = ups
	return res, nil
}

// McnUPOP .
func (s *Service) McnUPOP(c context.Context, arg *model.MCNUPStateOpReq) error {
	var (
		err             error
		upName, mcnName string
		accsReply       *accgrpc.InfosReply
		accInfos        map[int64]*accgrpc.Info
		up              *model.MCNUPInfoReply
		m               *model.MCNSignInfoReply
		now             = time.Now()
		xnow            = xtime.Time(now.Unix())
		date            = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		state           model.MCNUPState
	)
	state = arg.Action.GetState()
	if state == model.MCNUPStateUnKnown {
		log.Warn("mcn_up action(%d) to new err state(-1)", arg.Action)
		return ecode.MCNSignStateFlowErr
	}
	if !arg.Action.NotRightAction() {
		return ecode.MCNUpUnknownReviewErr
	}
	if arg.Action.NoRejectState() {
		xnow = 0
		arg.RejectReason = ""
	}
	if up, err = s.dao.McnUp(c, arg.SignUpID); err != nil {
		return err
	}
	if up == nil {
		return ecode.MCNCUpUnknownInfoErr
	}
	if m, err = s.dao.McnSign(c, up.SignID); err != nil {
		return err
	}
	if m.State != model.MCNSignStateOnSign && arg.Action.NoRejectState() {
		return ecode.MCNUpPassOnEffectSign
	}
	if up.State.IsOnReviewState() {
		return ecode.MCNUpOnlyReviewOpErr
	}
	if up.BeginDate.Time().After(date) && state == model.MCNUPStateOnSign {
		state = model.MCNUPStateOnPreOpen
	}
	if _, err = s.dao.UpMcnUpOP(c, arg.SignUpID, int8(state), xnow, arg.RejectReason); err != nil {
		return err
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: []int64{up.McnMid, up.UpMid}}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) error(%+v)", &accgrpc.MidsReq{Mids: []int64{up.McnMid, up.UpMid}}, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	if info, ok := accInfos[up.McnMid]; ok {
		mcnName = info.Name
	}
	if info, ok := accInfos[up.UpMid]; ok {
		upName = info.Name
	}
	if arg.Action == model.MCNUPActionPass {
		s.dao.UpMcnUpsRecommendOP(c, []int64{up.UpMid}, model.MCNUPRecommendStateDel)
	}
	s.worker.Add(func() {
		index := []interface{}{int8(state), arg.SignUpID, up.McnMid, mcnName, upName}
		content := map[string]interface{}{
			"sign_up_id":    arg.SignUpID,
			"sign_id":       up.SignID,
			"mcn_mid":       up.McnMid,
			"up_mid":        up.UpMid,
			"begin_date":    up.BeginDate,
			"end_date":      up.EndDate,
			"contract_link": up.ContractLink,
			"up_auth_link":  up.UpAuthLink,
			"reject_reason": arg.RejectReason,
			"reject_time":   now,
			"permission":    up.Permission,
		}
		s.AddAuditLog(context.Background(), model.MCNLogBizID, int8(arg.Action), arg.Action.String(), arg.UID, arg.UserName, []int64{up.UpMid}, index, content)
	})
	s.worker.Add(func() {
		argMsg := &model.ArgMsg{
			MSGType:     arg.Action.GetmsgType(true),
			MIDs:        []int64{up.McnMid},
			Reason:      arg.RejectReason,
			CompanyName: m.CompanyName,
			McnName:     mcnName,
			McnMid:      up.McnMid,
			UpName:      upName,
			UpMid:       up.UpMid,
		}
		s.sendMsg(context.Background(), argMsg)
		argMsg.MSGType = arg.Action.GetmsgType(false)
		argMsg.MIDs = []int64{up.UpMid}
		s.sendMsg(context.Background(), argMsg)
	})
	return nil
}

// McnPermitOP .
func (s *Service) McnPermitOP(c context.Context, arg *model.MCNSignPermissionReq) (err error) {
	var (
		open, closed []string
		m            *model.MCNSignInfoReply
	)
	if m, err = s.dao.McnSign(c, arg.SignID); err != nil {
		return err
	}
	if m == nil {
		return ecode.MCNCSignUnknownInfoErr
	}
	m.AttrPermitVal()
	arg.AttrPermitSet()
	open, closed = s.getPermitOpenOrClosed(arg.Permission, m.Permission)
	if len(open) == 0 && len(closed) == 0 {
		return
	}
	if _, err = s.dao.UpMCNPermission(c, arg.SignID, arg.Permission); err != nil {
		return
	}
	s.worker.Add(func() {
		index := []interface{}{arg.SignID}
		content := map[string]interface{}{
			"sign_id":    arg.SignID,
			"mcn_mid":    m.McnMid,
			"state":      m.State,
			"begin_date": m.BeginDate,
			"end_date":   m.EndDate,
			"permission": arg.Permission,
		}
		s.AddAuditLog(context.Background(), model.MCNLogBizID, int8(model.MCNSignActionPermit), model.MCNSignActionPermit.String(), arg.UID, arg.UserName, []int64{m.McnMid}, index, content)
	})
	s.worker.Add(func() {
		argMsg := &model.ArgMsg{
			MIDs: []int64{m.McnMid},
		}
		if len(open) > 0 {
			argMsg.MSGType = model.McnPermissionOpen
			argMsg.Permission = strings.Join(open, "、")
			s.sendMsg(context.Background(), argMsg)
		}
		if len(closed) > 0 {
			argMsg.MSGType = model.McnPermissionClosed
			argMsg.Permission = strings.Join(closed, "、")
			s.sendMsg(context.Background(), argMsg)
		}
	})
	s.worker.Add(func() {
		s.dao.DelMcnSignCache(context.Background(), m.McnMid)
	})
	return
}

func (s *Service) getPermitOpenOrClosed(a, b uint32) (open, closed []string) {
	for permit := range model.PermitMap {
		var c, d = model.AttrVal(a, uint(permit)), model.AttrVal(b, uint(permit))
		if c == d {
			continue
		}
		if c > d {
			open = append(open, permit.String())
		} else {
			closed = append(closed, permit.String())
		}
	}
	return
}

func (s *Service) getUpPermitString(permission uint32) (ps []string) {
	for permit := range model.PermitMap {
		var p = model.AttrVal(permission, uint(permit))
		if p <= 0 {
			continue
		}
		ps = append(ps, permit.String())
	}
	return
}

// McnUPPermitList .
func (s *Service) McnUPPermitList(c context.Context, arg *model.MCNUPPermitStateReq) (res *model.McnUpPermitApplyListReply, err error) {
	var (
		count              int64
		upMids, mids, tids []int64
		accsReply          *accgrpc.InfosReply
		accInfos           map[int64]*accgrpc.Info
		mbi                map[int64]*model.UpBaseInfo
		ms                 []*model.McnUpPermissionApply
	)
	res = new(model.McnUpPermitApplyListReply)
	res.Page = arg.Page
	if count, err = s.dao.McnUpPermitTotal(c, arg); err != nil {
		return
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if ms, err = s.dao.McnUpPermits(c, arg); err != nil {
		return
	}
	for _, m := range ms {
		upMids = append(upMids, m.UpMid)
		mids = append(mids, m.UpMid)
		mids = append(mids, m.McnMid)
	}
	if mbi, err = s.dao.UpBaseInfoMap(c, upMids); err != nil {
		return
	}
	for _, v := range mbi {
		tids = append(tids, int64(v.ActiveTid))
	}
	for _, m := range ms {
		if bi, ok := mbi[m.UpMid]; ok {
			m.ActiveTID = bi.ActiveTid
			m.FansCount = bi.FansCount
		}
	}
	if len(mids) > 0 {
		if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
			log.Error("s.accGRPC.Infos3(%+v) err(%v)", mids, err)
			err = nil
		} else {
			accInfos = accsReply.Infos
		}
	}
	tpNames := s.videoup.GetTidName(tids)
	for _, v := range ms {
		v.AttrPermitVal()
		if info, ok := accInfos[v.McnMid]; ok {
			v.McnName = info.Name
		}
		if info, ok := accInfos[v.UpMid]; ok {
			v.UpName = info.Name
		}
		if tyName, ok := tpNames[int64(v.ActiveTID)]; ok {
			v.TypeName = tyName
		} else {
			v.TypeName = model.DefaultTyName
		}
	}
	res.List = ms
	return
}

// McnUPPermitOP .
func (s *Service) McnUPPermitOP(c context.Context, arg *model.MCNUPPermitOPReq) (err error) {
	var (
		tx       *xsql.Tx
		ps       []string
		accReply *accgrpc.InfoReply
		now      = xtime.Time(time.Now().Unix())
		m        *model.McnUpPermissionApply
	)
	if !arg.Action.NotRightAction() {
		return ecode.MCNSignUnknownReviewErr
	}
	if m, err = s.dao.McnUpPermit(c, arg.ID); err != nil {
		return err
	}
	if m == nil {
		return ecode.MCNUpAbnormalDataErr
	}
	state := arg.Action.GetState()
	if state == model.MCNUPPermissionStateUnKnown {
		log.Warn("mcn-up permit action(%d) old state(%d) to new err state(-1)", arg.Action, m.State)
		return ecode.MCNUpPermitStateFlowErr
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	if arg.Action != model.MCNUPPermissionActionFail {
		now = 0
		arg.RejectReason = ""
	}
	m.State = state
	m.RejectReason = arg.RejectReason
	m.RejectTime = now
	m.AdminName = arg.UserName
	m.AdminID = arg.UID
	if _, err = s.dao.TxUpPermitApplyOP(tx, m); err != nil {
		return err
	}
	if state == model.MCNUPPermissionStatePass {
		if _, err = s.dao.TxMcnUpPermitOP(tx, m.SignID, m.McnMid, m.UpMid, m.NewPermission, m.UpAuthLink); err != nil {
			return err
		}
		ps = s.getUpPermitString(m.NewPermission)
	}
	if accReply, err = s.accGRPC.Info3(c, &accgrpc.MidReq{Mid: m.UpMid}); err != nil {
		log.Error("s.accGRPC.Info3(%+v) error(%+v)", &accgrpc.MidReq{Mid: m.McnMid}, err)
		err = nil
	}
	var name string
	if accReply.Info != nil {
		name = accReply.Info.Name
	}
	s.worker.Add(func() {
		argMsg := &model.ArgMsg{
			MIDs: []int64{m.McnMid},
		}
		if state == model.MCNUPPermissionStatePass {
			argMsg.UpName = name
			argMsg.MSGType = model.McnOperAgreeChangePermit
			argMsg.Permission = strings.Join(ps, "、")
		} else {
			argMsg.UpName = name
			argMsg.MSGType = model.McnOperNotAgreeChangePermit
			argMsg.Reason = m.RejectReason
		}
		s.sendMsg(context.Background(), argMsg)
	})
	s.worker.Add(func() {
		s.dao.DelMcnUpperCache(context.Background(), m.SignID, m.UpMid)
	})
	return
}

// MCNList .
func (s *Service) MCNList(c context.Context, arg *model.MCNListReq) (res *model.MCNListReply, err error) {
	var (
		count         int64
		signIDs, mids []int64
		mcns          []*model.MCNListOne
		accsReply     *accgrpc.InfosReply
		accInfos      map[int64]*accgrpc.Info
		payInfos      map[int64][]*model.SignPayInfoReply
	)
	res = new(model.MCNListReply)
	res.Page = arg.Page
	count, err = s.dao.MCNListTotal(c, arg)
	if err != nil {
		return
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if mcns, signIDs, mids, err = s.dao.MCNList(c, arg); err != nil {
		return
	}
	if len(signIDs) <= 0 {
		return
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) err(%v)", mids, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	if payInfos, err = s.dao.MCNPayInfos(c, signIDs); err != nil {
		return
	}
	for k, v := range mcns {
		v.AttrPermitVal()
		mcns[k].PayInfos = payInfos[v.ID]
		if info, ok := accInfos[v.MCNMID]; ok {
			v.MCNName = info.Name
		}
	}
	res.List = mcns
	return
}

// MCNPayEdit .
func (s *Service) MCNPayEdit(c context.Context, arg *model.MCNPayEditReq) error {
	if _, err := s.dao.UpMCNPay(c, arg); err != nil {
		return err
	}
	s.worker.Add(func() {
		index := []interface{}{arg.ID, arg.MCNMID}
		content := map[string]interface{}{
			"id":        arg.ID,
			"mcn_mid":   arg.MCNMID,
			"sign_id":   arg.SignID,
			"due_date":  arg.DueDate,
			"pay_value": arg.PayValue,
		}
		s.AddAuditLog(context.Background(), model.MCNPayDateLogBizID, int8(model.MCNSignCycleActionUp), model.MCNSignCycleActionUp.String(), arg.UID, arg.UserName, []int64{arg.SignID}, index, content)
	})
	s.worker.Add(func() {
		s.dao.DelMcnSignCache(context.Background(), arg.MCNMID)
	})
	return nil
}

// MCNPayStateEdit .
func (s *Service) MCNPayStateEdit(c context.Context, arg *model.MCNPayStateEditReq) error {
	var (
		err     error
		PayInfo *model.SignPayInfoReply
	)
	if _, err = s.dao.UpMCNPayState(c, arg); err != nil {
		return err
	}
	if PayInfo, err = s.dao.MCNPayInfo(c, arg); err != nil {
		log.Error("s.dao.MCNPayInfo(%+v) err(%v)", arg, err)
		err = nil
	}
	s.worker.Add(func() {
		index := []interface{}{arg.ID, arg.MCNMID}
		content := map[string]interface{}{
			"id":        arg.ID,
			"mcn_mid":   arg.MCNMID,
			"sign_id":   arg.SignID,
			"due_date":  PayInfo.DueDate,
			"pay_value": PayInfo.PayValue,
			"state":     arg.State,
		}
		s.AddAuditLog(context.Background(), model.MCNPayDateLogBizID, int8(model.MCNSignCycleActionUp), model.MCNSignCycleActionUp.String(), arg.UID, arg.UserName, []int64{arg.SignID}, index, content)
	})
	s.worker.Add(func() {
		s.dao.DelMcnSignCache(context.Background(), arg.MCNMID)
	})
	return err
}

// MCNStateEdit .
func (s *Service) MCNStateEdit(c context.Context, arg *model.MCNStateEditReq) error {
	if arg.Action != model.MCNSignActionBlock && arg.Action != model.MCNSignActionClear && arg.Action != model.MCNSignActionRestore {
		return ecode.RequestErr
	}
	var (
		err error
		ms  *model.MCNSignInfoReply
	)
	if ms, err = s.dao.McnSign(c, arg.ID); err != nil {
		return err
	}
	arg.State = arg.Action.GetState(ms.State)
	if arg.State == model.MCNSignStateUnKnown {
		log.Warn("mcn_sign action(%d) old state(%d) to new err state(-1)", arg.Action, ms.State)
		return ecode.MCNSignStateFlowErr
	}
	if _, err = s.dao.UpMCNState(c, arg); err != nil {
		return err
	}
	s.worker.Add(func() {
		index := []interface{}{arg.State, arg.ID}
		content := map[string]interface{}{
			"sign_id":              arg.ID,
			"mcn_mid":              ms.McnMid,
			"begin_date":           ms.BeginDate,
			"end_date":             ms.EndDate,
			"contract_link":        ms.ContractLink,
			"company_name":         ms.CompanyName,
			"company_license_id":   ms.CompanyLicenseID,
			"company_license_link": ms.CompanyLicenseLink,
			"contact_title":        ms.ContactTitle,
			"contact_idcard":       ms.ContactIdcard,
			"contact_phone":        ms.ContactPhone,
			"contact_name":         ms.ContactName,
			"state":                arg.State,
			"permission":           ms.Permission,
		}
		s.AddAuditLog(context.Background(), model.MCNLogBizID, int8(arg.Action), arg.Action.String(), arg.UID, arg.UserName, []int64{arg.MCNMID}, index, content)
	})
	s.worker.Add(func() {
		s.dao.DelMcnSignCache(context.Background(), arg.MCNMID)
	})
	s.worker.Add(func() {
		s.sendMsg(context.Background(), &model.ArgMsg{MSGType: arg.Action.GetmsgType(ms.State), MIDs: []int64{arg.MCNMID}})
	})
	return err
}

// MCNRenewal .
func (s *Service) MCNRenewal(c context.Context, arg *model.MCNRenewalReq) (err error) {
	var (
		signID       int64
		tx           *xsql.Tx
		stime, etime time.Time
		m            *model.MCNSignInfoReply
		ms           = &model.MCNSign{}
		ups          = make([]*model.MCNUP, 0)
	)
	if stime, err = time.ParseInLocation(model.TimeFormatDay, arg.BeginDate, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) error(%+v)", arg.BeginDate, err)
		return err
	}
	if etime, err = time.ParseInLocation(model.TimeFormatDay, arg.EndDate, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) error(%+v)", arg.EndDate, err)
		return err
	}
	if m, err = s.dao.McnSignByMCNMID(c, arg.MCNMID); err != nil {
		return err
	}
	if arg.ID != m.SignID {
		return ecode.MCNRenewalAlreadyErr
	}
	if stime.Before(m.EndDate.Time()) || etime.Before(stime) {
		return ecode.MCNRenewalDateErr
	}
	if m.State.IsRenewalState() && m.EndDate.Time().AddDate(0, 0, 1).After(time.Now()) && m.EndDate.Time().AddDate(0, 0, -30).Before(time.Now()) {
		return ecode.MCNRenewalNotInRangeErr
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()
	arg.AttrPermitSet()
	ms.MCNMID = m.McnMid
	ms.CompanyName = m.CompanyName
	ms.CompanyLicenseID = m.CompanyLicenseID
	ms.CompanyLicenseLink = m.CompanyLicenseLink
	ms.ContractLink = arg.ContractLink
	ms.ContactName = m.ContactName
	ms.ContactTitle = m.ContactTitle
	ms.ContactIdcard = m.ContactIdcard
	ms.ContactPhone = m.ContactPhone
	ms.BeginDate = xtime.Time(stime.Unix())
	ms.EndDate = xtime.Time(etime.Unix())
	ms.Permission = arg.Permission
	if signID, err = s.dao.TxAddMCNRenewal(tx, ms); err != nil {
		return err
	}
	if len(arg.SignPayInfo) > 0 {
		if err = s.dao.TxAddMCNPays(tx, signID, m.McnMid, arg.SignPayInfo); err != nil {
			return err
		}
	}
	if ups, err = s.dao.TxMCNRenewalUPs(tx, arg.ID, m.McnMid); err != nil {
		return err
	}
	if len(ups) > 0 {
		if err = s.dao.TxAddMCNUPs(tx, signID, m.McnMid, ups); err != nil {
			return err
		}
	}
	s.worker.Add(func() {
		index := []interface{}{int8(model.MCNSignStateOnSign), signID}
		content := map[string]interface{}{
			"sign_id":              signID,
			"mcn_mid":              m.McnMid,
			"begin_date":           ms.BeginDate,
			"end_date":             ms.EndDate,
			"contract_link":        m.ContractLink,
			"company_name":         m.CompanyName,
			"company_license_id":   m.CompanyLicenseID,
			"company_license_link": m.CompanyLicenseLink,
			"contact_title":        m.ContactTitle,
			"contact_idcard":       m.ContactIdcard,
			"contact_phone":        m.ContactPhone,
			"permission":           m.Permission,
		}
		s.AddAuditLog(context.Background(), model.MCNLogBizID, int8(model.MCNSignStateOnSign), model.MCNSignStateOnSign.String(), arg.UID, arg.UserName, []int64{arg.MCNMID}, index, content)
	})
	s.worker.Add(func() {
		index := []interface{}{arg.MCNMID}
		content := map[string]interface{}{
			"sign_id":       signID,
			"mcn_mid":       arg.MCNMID,
			"sign_pay_info": arg.SignPayInfo,
		}
		s.AddAuditLog(context.Background(), model.MCNPayDateLogBizID, int8(model.MCNSignCycleActionAdd), model.MCNSignCycleActionAdd.String(), arg.UID, arg.UserName, []int64{signID}, index, content)
	})
	s.worker.Add(func() {
		s.dao.DelMcnSignCache(context.Background(), arg.MCNMID)
	})
	s.worker.Add(func() {
		s.sendMsg(context.Background(), &model.ArgMsg{MSGType: model.McnRenewcontract, MIDs: []int64{arg.MCNMID}})
	})
	return
}

// MCNInfo .
func (s *Service) MCNInfo(c context.Context, arg *model.MCNInfoReq) (res *model.MCNInfoReply, err error) {
	var accReply *accgrpc.InfoReply
	if res, err = s.dao.MCNInfo(c, arg); err != nil || res == nil {
		return
	}
	if accReply, err = s.accGRPC.Info3(c, &accgrpc.MidReq{Mid: res.MCNMID}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) err(%v)", arg, err)
		err = nil
	}
	if accReply.Info != nil {
		res.MCNName = accReply.Info.Name
	}
	res.CompanyLicenseLink = model.BuildBfsURL(res.CompanyLicenseLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
	res.ContractLink = model.BuildBfsURL(res.ContractLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
	return
}

// MCNUPList .
func (s *Service) MCNUPList(c context.Context, arg *model.MCNUPListReq) (res *model.MCNUPListReply, err error) {
	var (
		count        int64
		upMIDs, tids []int64
		tpNames      map[int64]string
		mcn          *model.MCNSignInfoReply
		mcnUPs       []*model.MCNUPInfoReply
		accsReply    *accgrpc.InfosReply
		accInfos     map[int64]*accgrpc.Info
	)
	res = new(model.MCNUPListReply)
	res.Page = arg.Page
	if count, err = s.dao.MCNUPListTotal(c, arg); err != nil {
		return
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if mcnUPs, err = s.dao.MCNUPList(c, arg); err != nil {
		return
	}
	if mcn, err = s.dao.McnSign(c, arg.SignID); err != nil {
		log.Error("s.dao.McnSign(%d) err(%v)", arg.SignID, err)
		err = nil
	}
	for _, v := range mcnUPs {
		if mcn != nil {
			v.Permission = v.Permission & mcn.Permission
		}
		v.AttrPermitVal()
		upMIDs = append(upMIDs, v.UpMid)
		tids = append(tids, int64(v.ActiveTid))
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: upMIDs}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) err(%v)", upMIDs, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	tpNames = s.videoup.GetTidName(tids)
	for _, v := range mcnUPs {
		if info, ok := accInfos[v.UpMid]; ok {
			v.UpName = info.Name
		}
		if tyName, ok := tpNames[int64(v.ActiveTid)]; ok {
			v.TpName = tyName
		} else {
			v.TpName = model.DefaultTyName
		}
		v.UpAuthLink = model.BuildBfsURL(v.UpAuthLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
		v.ContractLink = model.BuildBfsURL(v.ContractLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
	}
	res.List = mcnUPs
	return
}

// MCNUPStateEdit .
func (s *Service) MCNUPStateEdit(c context.Context, arg *model.MCNUPStateEditReq) error {
	if arg.Action != model.MCNUPActionFreeze && arg.Action != model.MCNUPActionRelease && arg.Action != model.MCNUPActionRestore {
		return ecode.RequestErr
	}
	var (
		err             error
		up              *model.MCNUPInfoReply
		m               *model.MCNSignInfoReply
		upName, mcnName string
		accsReply       *accgrpc.InfosReply
		accInfos        map[int64]*accgrpc.Info
	)
	arg.State = arg.Action.GetState()
	if arg.State == model.MCNUPStateUnKnown {
		log.Warn("mcn_up action(%d) to new err state(-1)", arg.Action)
		return ecode.MCNSignStateFlowErr
	}
	if _, err = s.dao.UpMCNUPState(c, arg); err != nil {
		return err
	}
	if up, err = s.dao.McnUp(c, arg.ID); err != nil {
		log.Error("s.dao.McnUp(%v) err(%v)", arg.ID, err)
	}
	if m, err = s.dao.McnSign(c, up.SignID); err != nil {
		log.Error("s.dao.McnSign error(%+v)", err)
		err = nil
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: []int64{up.McnMid, up.UpMid}}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) error(%+v)", &accgrpc.MidsReq{Mids: []int64{up.McnMid, up.UpMid}}, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	if info, ok := accInfos[up.McnMid]; ok {
		mcnName = info.Name
	}
	if info, ok := accInfos[up.UpMid]; ok {
		upName = info.Name
	}
	s.worker.Add(func() {
		index := []interface{}{int8(arg.State), arg.ID, up.McnMid, up.SignID}
		content := map[string]interface{}{
			"sign_up_id":    arg.ID,
			"sign_id":       up.SignID,
			"mcn_mid":       up.McnMid,
			"up_mid":        up.UpMid,
			"begin_date":    up.BeginDate,
			"end_date":      up.EndDate,
			"contract_link": up.ContractLink,
			"up_auth_link":  up.UpAuthLink,
			"state":         arg.State,
			"reject_reason": up.RejectReason,
			"reject_time":   up.RejectTime,
		}
		s.AddAuditLog(context.Background(), model.MCNLogBizID, int8(arg.Action), arg.Action.String(), arg.UID, arg.UserName, []int64{arg.UPMID}, index, content)
	})
	s.worker.Add(func() {
		argMsg := &model.ArgMsg{
			MSGType:     arg.Action.GetmsgType(true),
			MIDs:        []int64{up.McnMid},
			CompanyName: m.CompanyName,
			McnName:     mcnName,
			McnMid:      up.McnMid,
			UpName:      upName,
			UpMid:       up.UpMid,
		}
		s.sendMsg(context.Background(), argMsg)
		argMsg.MSGType = arg.Action.GetmsgType(false)
		argMsg.MIDs = []int64{up.UpMid}
		s.sendMsg(context.Background(), argMsg)
	})
	return nil
}

// MCNCheatList .
func (s *Service) MCNCheatList(c context.Context, arg *model.MCNCheatListReq) (res *model.MCNCheatListReply, err error) {
	var (
		count     int64
		mids      []int64
		mcnCheats []*model.MCNCheatReply
		accsReply *accgrpc.InfosReply
		accInfos  map[int64]*accgrpc.Info
	)
	res = new(model.MCNCheatListReply)
	res.Page = arg.Page
	if count, err = s.dao.MCNCheatListTotal(c, arg); err != nil {
		return
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if mcnCheats, mids, err = s.dao.MCNCheatList(c, arg); err != nil {
		return
	}
	if len(mcnCheats) <= 0 || len(mids) <= 0 {
		return
	}
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) err(%v)", mids, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	for _, v := range mcnCheats {
		if info, ok := accInfos[v.UpMID]; ok {
			v.UpName = info.Name
		}
		if info, ok := accInfos[v.MCNMID]; ok {
			v.MCNName = info.Name
		}
	}
	res.List = mcnCheats
	return
}

// MCNCheatUPList .
func (s *Service) MCNCheatUPList(c context.Context, arg *model.MCNCheatUPListReq) (res *model.MCNCheatUPListReply, err error) {
	var (
		count        int64
		cheatUPInfos []*model.MCNCheatUPReply
	)
	res = new(model.MCNCheatUPListReply)
	res.Page = arg.Page
	if count, err = s.dao.MCNCheatUPListTotal(c, arg); err != nil {
		return
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if cheatUPInfos, err = s.dao.MCNCheatUPList(c, arg); err != nil {
		return
	}
	var mids []int64
	for _, v := range cheatUPInfos {
		mids = append(mids, v.MCNMID)
	}
	mids = up.SliceUnique(mids)
	var (
		accInfos  map[int64]*accgrpc.Info
		accsReply *accgrpc.InfosReply
	)
	if accsReply, err = s.accGRPC.Infos3(c, &accgrpc.MidsReq{Mids: mids}); err != nil {
		log.Error("s.accGRPC.Infos3(%+v) err(%v)", mids, err)
		err = nil
	} else {
		accInfos = accsReply.Infos
	}
	for _, v := range cheatUPInfos {
		if info, ok := accInfos[v.MCNMID]; ok {
			v.MCNName = info.Name
		}
	}
	res.List = cheatUPInfos
	return
}

// MCNImportUPInfo .
func (s *Service) MCNImportUPInfo(c context.Context, arg *model.MCNImportUPInfoReq) (res *model.MCNImportUPInfoReply, err error) {
	var profileRepley *accgrpc.ProfileReply
	if res, err = s.dao.MCNImportUPInfo(c, arg); err != nil {
		return
	}
	res.UpMID = arg.UPMID
	if profileRepley, err = s.accGRPC.Profile3(c, &accgrpc.MidReq{Mid: arg.UPMID}); err != nil {
		log.Error("s.accGRPC.Profile3(%d) err(%v)", arg.UPMID, err)
		err = nil
	}
	if profileRepley.Profile != nil {
		res.UpName = profileRepley.Profile.Name
		res.JoinTime = profileRepley.Profile.JoinTime
	}
	return
}

// MCNImportUPRewardSign .
func (s *Service) MCNImportUPRewardSign(c context.Context, arg *model.MCNImportUPRewardSignReq) (err error) {
	_, err = s.dao.UpMCNImportUPRewardSign(c, arg)
	return
}

// MCNIncreaseList .
func (s *Service) MCNIncreaseList(c context.Context, arg *model.MCNIncreaseListReq) (res *model.MCNIncreaseListReply, err error) {
	var (
		count         int64
		increaseDatas []*model.MCNIncreaseReply
	)
	res = new(model.MCNIncreaseListReply)
	res.Page = arg.Page
	if count, err = s.dao.MCNIncreaseListTotal(c, arg); err != nil {
		return
	}
	if count <= 0 {
		return
	}
	res.TotalCount = int(count)
	if increaseDatas, err = s.dao.MCNIncreaseList(c, arg); err != nil {
		return
	}
	res.List = increaseDatas
	return
}
