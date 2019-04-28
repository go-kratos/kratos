package service

import (
	"context"
	"time"

	adminmodel "go-common/app/admin/main/mcn/model"
	"go-common/app/interface/main/mcn/dao/cache"
	"go-common/app/interface/main/mcn/dao/global"
	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	accgrpc "go-common/app/service/main/account/api"
	memgrpc "go-common/app/service/main/member/api"
	memmdl "go-common/app/service/main/member/model"
	"go-common/app/service/main/member/model/block"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"go-common/app/interface/main/mcn/conf"
	"go-common/app/interface/main/mcn/dao/mcndao"

	"strings"

	"github.com/jinzhu/gorm"
)

// getMcnWithState
// if state is nil, state is not checked
func (s *Service) getMcnWithState(c context.Context, mcnmid int64, state ...model.MCNSignState) (mcnSign *mcnmodel.McnSign, err error) {
	mcnSign, err = s.mcndao.McnSign(c, mcnmid)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}
	if mcnSign == nil {
		err = ecode.NothingFound
		return
	}
	var ok = false
	if state == nil {
		ok = true
	} else {
		for _, s := range state {
			if mcnSign.State == s {
				ok = true
				break
			}
		}
	}

	if !ok {
		log.Info("mcnmid=%d, mcn is in %d, should in (%v)", mcnmid, mcnSign.State, state)
		err = ecode.MCNNotAllowed
		return
	}
	return
}

func (s *Service) checkPermission(c context.Context, mcnMid, upMid int64, permissions ...adminmodel.AttrBasePermit) (res bool) {
	var permLen = len(permissions)
	if permLen == 0 {
		return
	} else if permLen == 1 {
		// 基础权限直接放过
		if permissions[0] == adminmodel.AttrBasePermitBit {
			return true
		}
	}

	mcnSign, err := s.getMcnWithState(c, mcnMid, model.MCNSignStateOnSign)
	if err != nil {
		log.Error("get mcn fail, mcnmid=%d, err=%v", mcnMid, err)
		return
	}

	permForUp, err := s.mcndao.UpPermission(c, mcnSign.ID, upMid)
	if err != nil || permForUp == nil {
		log.Error("get up permission fail, signID=%d, upmid=%d, err=%v or up not found", mcnSign.ID, upMid, err)
		return
	}

	// 比较mcn与up的权限
	var wantPermission uint32
	for _, v := range permissions {
		wantPermission = wantPermission | (1 << v)
	}
	var resultPermission = wantPermission & mcnSign.Permission & permForUp.Permission
	if resultPermission != wantPermission {
		log.Warn("mcn doesnt have permission, mcn perm=0x%x, up perm=0x%x, want=0x%x, lack=0x%x", mcnSign.Permission, permForUp.Permission, wantPermission, resultPermission^wantPermission)
		return
	}
	res = true
	return
}

//McnGetState mcn state
func (s *Service) McnGetState(c context.Context, arg *mcnmodel.GetStateReq) (res *mcnmodel.McnGetStateReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}
	res = new(mcnmodel.McnGetStateReply)
	res.State = int8(mcnSign.State)
	if mcnSign.State == model.MCNSignStateOnReject {
		res.RejectReason = mcnSign.RejectReason
	}
	log.Info("mcn_state=%d, mcn_id=%d", res.State, arg.McnMid)
	return
}

//McnExist .
func (s *Service) McnExist(c context.Context, arg *mcnmodel.GetStateReq) (res *mcnmodel.McnExistReply, err error) {
	res = new(mcnmodel.McnExistReply)
	_, err = s.getMcnWithState(c, arg.McnMid)
	if err == ecode.NothingFound {
		res.Exist = 0
		return
	} else if err != nil {
		log.Error("error get state, err=%s", err)
		return
	}

	res.Exist = 1
	return
}

// McnBaseInfo .
func (s *Service) McnBaseInfo(c context.Context, arg *mcnmodel.GetStateReq) (res *mcnmodel.McnBaseInfoReply, err error) {
	res = new(mcnmodel.McnBaseInfoReply)
	mcnSign, err := s.getMcnWithState(c, arg.McnMid)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}
	res.CopyFromMcnInfo(mcnSign)
	return
}

//McnApply .
func (s *Service) McnApply(c context.Context, arg *mcnmodel.McnApplyReq) (res *mcnmodel.CommonReply, err error) {

	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnReject, model.MCNSignStateNoApply)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	var sign mcnmodel.McnSign
	if err = s.uniqueChecker.CheckIsUniqe(arg); err != nil {
		log.Info("check unique fail, err=%s, arg=%v", err, arg)
		return
	}

	arg.CopyTo(&sign)
	sign.ID = mcnSign.ID
	sign.State = model.MCNSignStateOnReview
	var db = s.mcndao.GetMcnDB()
	if err = db.Table(sign.TableName()).Where("id=?", sign.ID).Updates(map[string]interface{}{
		"company_name":         sign.CompanyName,
		"company_license_id":   sign.CompanyLicenseID,
		"contact_name":         sign.ContactName,
		"contact_title":        sign.ContactTitle,
		"contact_idcard":       sign.ContactIdcard,
		"contact_phone":        sign.ContactPhone,
		"company_license_link": sign.CompanyLicenseLink,
		"contract_link":        sign.ContractLink,
		"state":                sign.State,
	}).Error; err != nil {
		log.Error("save mcn fail, mcn mid=%d, row id=%d", sign.McnMid, sign.ID)
		err = ecode.ServerErr
		return
	}
	s.mcndao.DelCacheMcnSign(c, arg.McnMid)
	s.worker.Add(func() {
		s.loadMcnUniqueCache()
	})
	return
}

//McnBindUpApply .
func (s *Service) McnBindUpApply(c context.Context, arg *mcnmodel.McnBindUpApplyReq) (res *mcnmodel.McnBindUpApplyReply, err error) {
	if arg.BeginDate > arg.EndDate {
		err = ecode.MCNUpBindUpSTimeLtETime
		return
	}
	// 查询mcn状态
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	// 0.检查是否封禁
	var blockArg = memgrpc.MemberMidReq{Mid: arg.UpMid, RemoteIP: metadata.String(c, metadata.RemoteIP)}
	var blockInfo, e = global.GetMemGRPC().BlockInfo(c, &blockArg)
	if e == nil {
		if blockInfo.BlockStatus > int32(block.BlockStatusFalse) {
			log.Info("up is blocked, mid=%d, blockstatus=%d", arg.UpMid, blockInfo.BlockStatus)
			err = ecode.MCNUpBindUpIsBlocked
			return
		}
	} else {
		log.Error("get block info error, err=%s", e)
	}

	// 1.检查是否是蓝V用户
	var (
		memberInfo *memgrpc.MemberInfoReply
		memberArg  = memgrpc.MemberMidReq{Mid: arg.UpMid, RemoteIP: metadata.String(c, metadata.RemoteIP)}
	)
	if memberInfo, err = global.GetMemGRPC().Member(c, &memberArg); err != nil {
		log.Error("get member info error, err=%s", err)
	} else {
		if memberInfo.OfficialInfo != nil &&
			(memberInfo.OfficialInfo.Role == memmdl.OfficialRoleBusiness ||
				memberInfo.OfficialInfo.Role == memmdl.OfficialRoleGov ||
				memberInfo.OfficialInfo.Role == memmdl.OfficialRoleMedia ||
				memberInfo.OfficialInfo.Role == memmdl.OfficialRoleOther) {
			err = ecode.MCNUpBindUpIsBlueUser
			return
		}
	}

	// 2.查询当前up状态
	upList, err := s.mcndao.GetUpBind("up_mid=?", arg.UpMid)
	if err != nil {
		log.Error("get up bind fail, err=%s", err)
		err = ecode.ServerErr
		return
	}

	// 3.可以申请绑定的up主才能绑定
	var mcnUp *mcnmodel.McnUp
	for _, v := range upList {
		if !v.IsBindable() {
			log.Info("up is in state(%d), cannot be bind again. id=%d, upmid=%d, signid=%d, mcnSign=%d", v.State, v.ID, v.UpMid, v.SignID, v.McnMid)
			err = ecode.MCNUpCannotBind
			return
		}
		if v.IsBeingBindedWithMcn(mcnSign) {
			log.Info("up is being binded with mcnSign, state=%d, id=%d, signid=%d, mcnSign=%d, mcn_mid=%d", v.State, v.State, v.ID, v.SignID, v.McnMid)
			err = ecode.MCNUpBindUpAlreadyInProgress
			return
		}
		if v.SignID == mcnSign.ID {
			mcnUp = v
		}
	}

	if arg.UpType == 1 {
		// 站外up主需要满足条件：
		// 1.粉丝数≤100  或  2. 投稿数＜2及90天内未投稿 （1，2并列关系，满足其一即可申请）
		baseInfoMap, e := s.mcndao.GetUpBaseInfo("article_count_accumulate, activity, fans_count, mid", []int64{arg.UpMid})
		if e == nil {
			var upInfo, ok = baseInfoMap[arg.UpMid]
			if ok && upInfo != nil {
				//upInfo.Activity 1高，2中，3低，4流失
				//高=30天内有投稿
				//中=31~90天内有投稿
				//低=91~180天内有投稿
				//流失=180内以上未投稿
				if !(upInfo.FansCount <= 100 || (upInfo.ArticleCountAccumulate < 2 && upInfo.Activity > 2)) {
					err = ecode.MCNUpOutSiteIsNotQualified
					log.Error("outsite cannot bind, up fans count(%d) > 100", upInfo.FansCount)
					return
				}
			} else {
				log.Warn("up info is not found in up base info, up=%d", arg.UpMid)
			}
		}
	}

	// 站外信息是否OK
	if !arg.IsSiteInfoOk() {
		err = ecode.MCNUpBindInvalidURL
		log.Warn("arg error, up is out site up, but site url is not valid, arg=%v", arg)
		return
	}

	// 只能设置mcn自己有的权限，如果要设置其他权限，返回错误。
	// 只有mcn有的权限，才可以申请up主的权限
	_, err = mcnShouldContainUpPermission(mcnSign.Permission, arg.GetAttrPermitVal())
	if err != nil {
		return
	}

	// 3.绑定Up主，如果已有记录，则更新记录
	bindup, affectedRow, err := s.mcndao.BindUp(mcnUp, mcnSign, arg)
	if err != nil {
		log.Error("fail to bind up, mcnmid=%d, upmid=%d err=%s", arg.McnMid, arg.UpMid, err)
		return
	}
	res = new(mcnmodel.McnBindUpApplyReply)
	res.BindID = bindup.ID
	// 4.发送站内信息
	if arg.UpMid != arg.McnMid {
		var nickname = global.GetName(c, arg.McnMid)
		var msg = adminmodel.ArgMsg{
			MSGType:     adminmodel.McnUpBindAuthApply,
			MIDs:        []int64{arg.UpMid},
			McnName:     nickname,
			McnMid:      arg.McnMid,
			CompanyName: mcnSign.CompanyName,
			SignUpID:    bindup.ID,
		}
		s.sendMsg(&msg)
	}
	log.Info("bind up apply success, mcn=%d, upmid=%d, rowaffected=%d", arg.McnMid, arg.UpMid, affectedRow)
	return
}

//McnUpConfirm .
func (s *Service) McnUpConfirm(c context.Context, arg *mcnmodel.McnUpConfirmReq) (res *mcnmodel.CommonReply, err error) {

	// 1.查询当前up状态
	upList, err := s.mcndao.GetUpBind("id=? and up_mid=? and state=?", arg.BindID, arg.UpMid, model.MCNUPStateNoAuthorize)
	if err != nil {
		log.Error("get up bind fail, err=%s", err)
		err = ecode.ServerErr
		return
	}

	// 不存在
	if len(upList) == 0 {
		log.Info("bind id not found, id=%d", arg.BindID)
		err = ecode.MCNNotAllowed
		return
	}

	var upBind = upList[0]
	// 查询mcn状态
	mcnSign, err := s.getMcnWithState(c, upBind.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.MCNStateInvalid {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if mcnSign.ID != upBind.SignID {
		log.Error("bind id not same with mcn signid, bind id=%d, signid=%d", upBind.ID, upBind.SignID)
		err = ecode.MCNUpBindInvalid
		return
	}

	var state = model.MCNUPStateOnRefuse
	if arg.Choice {
		state = model.MCNUPStateOnReview
	}
	// 更新状态
	err = s.mcndao.UpConfirm(arg, state)
	if err != nil {
		log.Error("fail to update up bind, bind id=%d, upmid=%d, err=%s", arg.BindID, arg.UpMid, err)
		err = ecode.ServerErr
		return
	}
	// 同意
	if arg.Choice {
		var mcnName = global.GetName(c, mcnSign.McnMid)
		var msg = adminmodel.ArgMsg{
			MSGType:     adminmodel.McnUpBindAuthReview,
			MIDs:        []int64{arg.UpMid},
			McnName:     mcnName,
			McnMid:      mcnSign.McnMid,
			CompanyName: mcnSign.CompanyName,
		}
		s.sendMsg(&msg)
	} else {
		var upName = global.GetName(c, arg.UpMid)
		var msg = adminmodel.ArgMsg{
			MSGType: adminmodel.McnUpBindAuthApplyRefuse,
			MIDs:    []int64{mcnSign.McnMid},
			UpMid:   arg.UpMid,
			UpName:  upName,
		}
		s.sendMsg(&msg)
	}
	log.Info("up bind change, bind id=%d, upmid=%d, isaccept=%t", arg.BindID, arg.UpMid, arg.Choice)
	return
}

//McnUpGetBind .
func (s *Service) McnUpGetBind(c context.Context, arg *mcnmodel.McnUpGetBindReq) (res *mcnmodel.McnGetBindReply, err error) {
	res, err = s.mcndao.GetBindInfo(arg)
	if err != nil {
		log.Error("fail to get bind info, err=%s", err)
		return
	}

	accInfo, err := global.GetInfo(c, int64(res.McnMid))
	if err == nil && accInfo != nil {
		res.McnName = accInfo.Name
	}

	res.Finish()
	res.UpAuthLink = model.BuildBfsURL(res.UpAuthLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
	return
}

//McnDataSummary .
func (s *Service) McnDataSummary(c context.Context, arg *mcnmodel.McnGetDataSummaryReq) (res *mcnmodel.McnGetDataSummaryReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	var today = time.Now().Add(-12 * time.Hour)
	res, err = s.datadao.GetMcnSummaryCache(c, mcnSign.ID, today)
	if err != nil {
		log.Error("fail to get mcn data, sign id=%d, mcnmid=%d, err=%s", mcnSign.ID, mcnSign.McnMid, err)
		return
	}
	// today is not found, try yesterday
	if res == nil {
		res, err = s.mcndao.McnDataSummary(c, mcnSign.ID, today.AddDate(0, 0, -1))
		if err != nil {
			log.Error("fail to get mcn data, sign id=%d, mcnmid=%d, err=%s", mcnSign.ID, mcnSign.McnMid, err)
			return
		}
	}

	if res == nil {
		log.Error("fail to get mcn data, res = nil, sign id=%d", mcnSign.ID)
		res = new(mcnmodel.McnGetDataSummaryReply)
	}
	return
}

//McnDataUpList .
func (s *Service) McnDataUpList(c context.Context, arg *mcnmodel.McnGetUpListReq) (res *mcnmodel.McnGetUpListReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	generateDate, err := s.mcndao.GetDataUpLatestDate(mcnmodel.DataTypeAccumulate, mcnSign.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			log.Warn("no data list found for mcn=%d, sign id=%d", mcnSign.McnMid, mcnSign.ID)
			var now = time.Now()
			generateDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		} else {
			log.Error("fail to get latest generate up date, err=%s", err)
			return
		}
	}

	// 获取数据
	upData, err := s.mcndao.GetAllUpData(int64(mcnSign.ID), arg.UpMid, generateDate)
	// 在正式数据出来之前，临时使用
	//upData, err := s.mcndao.GetAllUpDataTemp(int64(mcnSign.ID), arg.UpMid, time.Now())
	var mids []int64
	for _, v := range upData {
		mids = append(mids, v.UpMid)
		v.Permission = v.Permission & mcnSign.Permission
	}
	var infosReply *accgrpc.InfosReply
	var midtidmap map[int64]int64
	var accInfos map[int64]*accgrpc.Info
	if len(mids) > 0 {
		var e error
		infosReply, e = global.GetAccGRPC().Infos3(c, &accgrpc.MidsReq{Mids: mids})
		if e != nil {
			log.Warn("fail to get info, err=%s", e)
		} else {
			accInfos = infosReply.Infos
		}

		midtidmap, e = s.mcndao.GetActiveTid(mids)
		if e != nil {
			log.Warn("fail to get activit, err=%s", e)
		}
	}
	res = new(mcnmodel.McnGetUpListReply)
	for _, v := range upData {
		var info, ok = accInfos[v.UpMid]
		if ok {
			v.Name = info.Name
		}
		if v.State != int8(model.MCNUPStateOnSign) {
			// MCNUPStateOnSign 与 MCNUPStateOnPreOpen 状态下 不隐藏时间
			v.HideData(!(v.State == int8(model.MCNUPStateOnSign) ||
				v.State == int8(model.MCNUPStateOnPreOpen)))
		}

		tid, ok := midtidmap[v.UpMid]
		if ok {
			v.TidName = cache.GetTidName(tid)
			v.ActiveTid = int16(tid)
		}
		res.Result = append(res.Result, v)
	}
	res.Finish()
	return
}

//McnGetOldInfo .
func (s *Service) McnGetOldInfo(c context.Context, arg *mcnmodel.McnGetMcnOldInfoReq) (res *mcnmodel.McnGetMcnOldInfoReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateNoApply)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	info, err := s.mcndao.GetMcnOldInfo(mcnSign.McnMid)
	if err != nil {
		if err == ecode.NothingFound {
			err = nil
		} else {
			log.Error("fail get mcn old info err=%s", err)
			return
		}
	}

	res = new(mcnmodel.McnGetMcnOldInfoReply)
	res.Copy(info)
	return
}

func getUpPermitString(permission uint32) (ps []string) {
	for permit := range adminmodel.PermitMap {
		var p = adminmodel.AttrVal(permission, uint(permit))
		if p <= 0 {
			continue
		}
		ps = append(ps, permit.String())
	}
	return
}

// 检查up主权限，是否是mcn的子集
// mcnPermission mcn自己的permission
// upPermission up的permission
// return finalPermission = mcnPermission &upPermission
func mcnShouldContainUpPermission(mcnPermission, upPermission uint32) (finalPermission uint32, err error) {
	// 3.只能设置mcn自己有的权限，如果要设置其他权限，返回错误。
	// 只有mcn有的权限，才可以申请up主的权限
	finalPermission = mcnPermission & upPermission
	if finalPermission != upPermission {
		log.Error("mcn has no permission to change, mcn=0x%x, wantup=0x%x, notallowd=0x%x", mcnPermission, upPermission, finalPermission^upPermission)
		err = ecode.MCNChangePermissionLackPermission
		return
	}
	return
}

//McnChangePermit change up's permission
func (s *Service) McnChangePermit(c context.Context, arg *mcnmodel.McnChangePermitReq) (res *mcnmodel.McnChangePermitReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	// 1.检查Up主关系，只有“已签约”,“待开启”状态的Up主可以修改
	// 2.查询当前up状态
	upList, err := s.mcndao.GetUpBind("up_mid=? and sign_id=? and state in (?)", arg.UpMid, mcnSign.ID, mcndao.UpSignedStates)
	if err != nil {
		log.Error("get up bind fail, err=%v", err)
		err = ecode.ServerErr
		return
	}

	if len(upList) == 0 {
		log.Error("up is not in signed state with mcn, up_mid=%d and sign_id=%d , need state in (%v)", arg.UpMid, mcnSign.ID, mcndao.UpSignedStates)
		err = ecode.MCNUpSignStateInvalid
		return
	}
	var oldUp = upList[0]
	var newPermission = arg.GetAttrPermitVal()
	if oldUp.Permission == uint32(newPermission) {
		log.Error("permission not changed, arg=%+v", arg)
		err = ecode.MCNChangePermissionSamePermission
		return
	}
	// 2.只能设置mcn自己有的权限，如果要设置其他权限，返回错误。
	// 只有mcn有的权限，才可以申请up主的权限
	maskedPermission, err := mcnShouldContainUpPermission(mcnSign.Permission, newPermission)
	if err != nil {
		return
	}

	// 如果是自己，则直接进行修改
	if arg.UpMid == mcnSign.McnMid {
		var _, e = s.mcndao.UpdateBindUp(map[string]interface{}{
			"permission": maskedPermission,
		}, "up_mid=? and sign_id=?", arg.UpMid, mcnSign.ID)
		if e != nil {
			err = e
			log.Error("fail to change up permission, err=%v, arg=%v", err, arg)
			return
		}
		return
	}

	// 3.检查是否有对应up主的修改请求,如果有就拒绝这次修改
	existedApply, _ := s.mcndao.GetUpPermissionApply("id", "sign_id=? and up_mid=? and state in (?)", mcnSign.ID, arg.UpMid, mcndao.UpPermissionApplyCannotApplyStates)
	if len(existedApply) > 0 {
		log.Error("apply already exist, id=%d, sign id=%d, mid=%d", existedApply[0].ID, existedApply[0].SignID, existedApply[0].UpMid)
		err = ecode.MCNChangePermissionAlreadyInProgress
		return
	}

	// 真的去增加permission
	var permissionApply = mcnmodel.McnUpPermissionApply{
		SignID:        mcnSign.ID,
		McnMid:        mcnSign.McnMid,
		UpMid:         arg.UpMid,
		NewPermission: maskedPermission,
		OldPermission: oldUp.Permission,
		UpAuthLink:    arg.UpAuthLink,
	}

	var db = s.mcndao.GetMcnDB()
	err = db.Create(&permissionApply).Error
	if err != nil {
		log.Error("create permission apply fail, err=%v, arg=%+v", err, arg)
		return
	}
	// 返回bind_id
	res = &mcnmodel.McnChangePermitReply{BindID: permissionApply.ID}
	// 4.发送站内信息
	if arg.UpMid != arg.McnMid {
		var nickname = global.GetName(c, arg.McnMid)
		var msg = adminmodel.ArgMsg{
			MSGType:     adminmodel.McnApplyUpChangePermit,
			MIDs:        []int64{arg.UpMid},
			McnName:     nickname,
			McnMid:      arg.McnMid,
			CompanyName: mcnSign.CompanyName,
			SignUpID:    permissionApply.ID,
			Permission:  strings.Join(getUpPermitString(maskedPermission), "、"),
		}
		s.sendMsg(&msg)
	}
	return
}

//McnPermitApplyGetBind get permit apply bind
func (s *Service) McnPermitApplyGetBind(c context.Context, arg *mcnmodel.McnUpGetBindReq) (res *mcnmodel.McnGetBindReply, err error) {
	res, err = s.mcndao.GetUpPermissionBindInfo(arg)
	if err != nil {
		log.Error("fail to get bind info, err=%s", err)
		return
	}

	accInfo, err := global.GetInfo(c, int64(res.McnMid))
	if err == nil && accInfo != nil {
		res.McnName = accInfo.Name
	}

	res.Finish()
	res.UpAuthLink = model.BuildBfsURL(res.UpAuthLink, s.c.BFS.Key, s.c.BFS.Secret, s.c.BFS.Bucket, model.BfsEasyPath)
	return
}

//McnUpPermitApplyConfirm permit apply confirm
func (s *Service) McnUpPermitApplyConfirm(c context.Context, arg *mcnmodel.McnUpConfirmReq) (res *mcnmodel.CommonReply, err error) {

	// 1.查询当前up状态
	upList, err := s.mcndao.GetUpPermissionApply("*", "id=? and up_mid=? and state=?", arg.BindID, arg.UpMid, adminmodel.MCNUPPermissionStateNoAuthorize)
	if err != nil {
		log.Error("get up bind fail, err=%s", err)
		err = ecode.ServerErr
		return
	}

	// 不存在
	if len(upList) == 0 {
		log.Info("bind id not found, id=%d", arg.BindID)
		err = ecode.MCNNotAllowed
		return
	}

	var upBind = upList[0]
	// 查询mcn状态
	mcnSign, err := s.getMcnWithState(c, upBind.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.MCNStateInvalid {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if mcnSign.ID != upBind.SignID {
		log.Error("bind id not same with mcn signid, bind id=%d, signid=%d", upBind.ID, upBind.SignID)
		err = ecode.MCNUpBindInvalid
		return
	}

	var state = adminmodel.MCNUPPermissionStateOnRefuse
	if arg.Choice {
		state = adminmodel.MCNUPPermissionStateReview
	}
	// 更新状态
	err = s.mcndao.UpPermissionConfirm(arg, state)
	if err != nil {
		log.Error("fail to update up bind, bind id=%d, upmid=%d, err=%s", arg.BindID, arg.UpMid, err)
		err = ecode.ServerErr
		return
	}
	// 同意
	if arg.Choice {
		// do nothing.
	} else {
		var upName = global.GetName(c, arg.UpMid)
		var msg = adminmodel.ArgMsg{
			MSGType: adminmodel.McnUpNotAgreeChangePermit,
			MIDs:    []int64{mcnSign.McnMid},
			UpMid:   arg.UpMid,
			UpName:  upName,
		}
		s.sendMsg(&msg)
	}
	log.Info("up permission bind change, bind id=%d, upmid=%d, isaccept=%t", arg.BindID, arg.UpMid, arg.Choice)
	return
}

//McnPublicationPriceChange .
func (s *Service) McnPublicationPriceChange(c context.Context, arg *mcnmodel.McnPublicationPriceChangeReq) (res *mcnmodel.McnPublicationPriceChangeReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	// 1.检查上次刊例价修改时间,如果时间内不能修改返回错误
	publicationPriceCache, e := s.mcndao.CachePublicationPrice(c, mcnSign.ID, arg.UpMid)
	if e != nil {
		log.Warn("get modify time from cache fail, arg=%+v, err=%v", arg, err)
	}

	if publicationPriceCache == nil {
		// 初始化为0值
		publicationPriceCache = &mcnmodel.PublicationPriceCache{}
	}

	var lastModifyTime = publicationPriceCache.ModifyTime // 从缓存中获取
	var now = time.Now()
	if now.Before(lastModifyTime.Add(time.Duration(conf.Conf.Other.PublicationPriceChangeLimit))) {
		log.Error("publication change fail, last modify time=%s, timelimit=%+v, arg=%+v", lastModifyTime, conf.Conf.Other.PublicationPriceChangeLimit, arg)
		err = ecode.MCNPublicationFailTimeLimit
		return
	}

	// 2.检查Up主关系，只有“已签约”,“待开启”状态的Up主可以修改
	upList, err := s.mcndao.GetUpBind("up_mid=? and sign_id=? and state in (?)", arg.UpMid, mcnSign.ID, mcndao.UpSignedStates)
	if err != nil {
		log.Error("get up bind fail, err=%v", err)
		err = ecode.ServerErr
		return
	}

	if len(upList) == 0 {
		log.Error("up is not in signed state with mcn, up_mid=%d and sign_id=%d , need state in (%v)", arg.UpMid, mcnSign.ID, mcndao.UpSignedStates)
		err = ecode.MCNUpSignStateInvalid
		return
	}

	var up = upList[0]
	// 3.修改刊例价,更新上次修改时间
	var db = s.mcndao.GetMcnDB()
	err = db.Table(mcnmodel.TableNameMcnUp).Where("id=?", up.ID).Update("publication_price", arg.Price).Error
	if err != nil {
		log.Error("change publication price fail, err=%v, arg=%+v", err, arg)
		return
	}
	publicationPriceCache.ModifyTime = now
	e = s.mcndao.AddCachePublicationPrice(c, mcnSign.ID, publicationPriceCache, arg.UpMid)
	if e != nil {
		log.Warn("fail to add cache, arg=%+v, err=%v", arg, e)
	}
	return
}
