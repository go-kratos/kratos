package api

import (
	"go-common/app/service/main/member/model"
	block "go-common/app/service/main/member/model/block"
)

// FromBaseInfo convert from model.BaseInfo to v1.BaseInfoReply
func FromBaseInfo(model *model.BaseInfo) *BaseInfoReply {
	baseInfoReply := &BaseInfoReply{
		Mid:      model.Mid,
		Name:     model.Name,
		Sex:      model.Sex,
		Face:     model.Face,
		Sign:     model.Sign,
		Rank:     model.Rank,
		Birthday: model.Birthday,
	}
	return baseInfoReply
}

// FromLevelInfo convert from model.LevelInfo to v1.LevelInfoReply
func FromLevelInfo(model *model.LevelInfo) *LevelInfoReply {
	levelInfoReply := &LevelInfoReply{
		Cur:     model.Cur,
		Min:     model.Min,
		NowExp:  model.NowExp,
		NextExp: model.NextExp,
	}
	return levelInfoReply
}

// FromOfficialInfo convert from model.OfficialInfo to v1.OfficialInfoReply
func FromOfficialInfo(model *model.OfficialInfo) *OfficialInfoReply {
	officialInfoReply := &OfficialInfoReply{
		Role:  model.Role,
		Title: model.Title,
		Desc:  model.Desc,
	}
	return officialInfoReply
}

// FromMember convert from model.Member to v1.MemberInfoReply
func FromMember(res *model.Member) *MemberInfoReply {
	var baseInfo *BaseInfoReply
	var levelInfo *LevelInfoReply
	var officialInfo *OfficialInfoReply
	if res.BaseInfo != nil {
		baseInfo = FromBaseInfo(res.BaseInfo)
	}
	if res.LevelInfo != nil {
		levelInfo = FromLevelInfo(res.LevelInfo)
	}
	if res.OfficialInfo != nil {
		officialInfo = FromOfficialInfo(res.OfficialInfo)
	}
	memberInfoReply := &MemberInfoReply{
		BaseInfo:     baseInfo,
		LevelInfo:    levelInfo,
		OfficialInfo: officialInfo,
	}
	return memberInfoReply
}

// FromOfficialDoc convert from model.OfficalDoc to v1.OfficialDocInfoReply
func FromOfficialDoc(model *model.OfficialDoc) *OfficialDocInfoReply {
	officalDocInfoReply := &OfficialDocInfoReply{
		Mid:              model.Mid,
		Name:             model.Name,
		State:            int32(model.State),
		Role:             int8(model.Role),
		Title:            model.Title,
		Desc:             model.Desc,
		RejectReason:     model.RejectReason,
		Realname:         int8(model.Realname),
		Operator:         model.Operator,
		Telephone:        model.Telephone,
		Email:            model.Email,
		Address:          model.Address,
		Company:          model.Company,
		CreditCode:       model.CreditCode,
		Organization:     model.Organization,
		OrganizationType: model.OrganizationType,
		BusinessLicense:  model.BusinessLicense,
		BusinessScale:    model.BusinessScale,
		BusinessLevel:    model.BusinessLevel,
		BusinessAuth:     model.BusinessAuth,
		Supplement:       model.Supplement,
		Professional:     model.Professional,
		Identification:   model.Identification,
	}

	return officalDocInfoReply
}

// FromBlockInfo convert from model.BlockInfo to v1.OfficialDocInfoReply
func FromBlockInfo(model *block.BlockInfo) *BlockInfoReply {
	blockInfoReply := &BlockInfoReply{
		MID:         model.MID,
		BlockStatus: int32(model.BlockStatus),
		StartTime:   model.StartTime,
		EndTime:     model.EndTime,
	}
	return blockInfoReply
}

// FromBlockUserDetail convert from model.BlockUserDetail v1.OfficialDocInfoReply
func FromBlockUserDetail(model *block.BlockUserDetail) *BlockDetailReply {
	blockDetailReply := &BlockDetailReply{
		MID:        model.MID,
		BlockCount: model.BlockCount,
	}
	return blockDetailReply
}

// ToArgOfficialDoc convert from v1.officalDocReq to model.AragsOfficalDoc
func ToArgOfficialDoc(req *OfficialDocReq) *model.ArgOfficialDoc {
	argOfficialDoc := &model.ArgOfficialDoc{
		Mid:              req.Mid,
		Name:             req.Name,
		Role:             req.Role,
		Title:            req.Title,
		Desc:             req.Desc,
		Realname:         int8(req.Realname),
		Operator:         req.Operator,
		Telephone:        req.Telephone,
		Email:            req.Email,
		Address:          req.Address,
		Company:          req.Company,
		CreditCode:       req.CreditCode,
		Organization:     req.Organization,
		OrganizationType: req.OrganizationType,
		BusinessLicense:  req.BusinessLicense,
		BusinessScale:    req.BusinessScale,
		BusinessLevel:    req.BusinessLevel,
		BusinessAuth:     req.BusinessAuth,
		Supplement:       req.Supplement,
		Professional:     req.Professional,
		Identification:   req.Identification,
		SubmitSource:     req.SubmitSource,
	}
	return argOfficialDoc
}

// ToArgUpdateMoral convert from v1.UpdateMoralReq to model.ArgUpdateMoral
func ToArgUpdateMoral(req *UpdateMoralReq) *model.ArgUpdateMoral {
	updateMoral := &model.ArgUpdateMoral{
		Mid:        req.Mid,
		Delta:      req.Delta,
		Origin:     req.Origin,
		Reason:     req.Reason,
		ReasonType: req.ReasonType,
		Operator:   req.Operator,
		Remark:     req.Remark,
		Status:     req.Status,
		IsNotify:   req.IsNotify,
		IP:         req.Ip,
	}
	return updateMoral
}

// ToArgUpdateMorals convert from v1.UpdateMoralsReq to model.ArgUpdateMorals
func ToArgUpdateMorals(req *UpdateMoralsReq) *model.ArgUpdateMorals {
	updateMorals := &model.ArgUpdateMorals{
		Mids:       req.Mids,
		Delta:      req.Delta,
		Origin:     req.Origin,
		Reason:     req.Reason,
		ReasonType: req.ReasonType,
		Operator:   req.Operator,
		Remark:     req.Remark,
		Status:     req.Status,
		IsNotify:   req.IsNotify,
		IP:         req.Ip,
	}
	return updateMorals
}
