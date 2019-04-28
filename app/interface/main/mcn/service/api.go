package service

import (
	"context"

	"go-common/app/interface/main/mcn/dao/mcndao"
	"go-common/app/interface/main/mcn/model/datamodel"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/log"
)

// McnGetRankArchiveLikesAPI get rank archive likes
func (s *Service) McnGetRankArchiveLikesAPI(c context.Context, arg *mcnmodel.McnGetRankAPIReq) (res *mcnmodel.McnGetRankUpFansReply, err error) {
	res, err = s.getRankResultAPI(c, arg, s.mcndao.GetRankArchiveLikes)
	return
}

func (s *Service) getRankResultAPI(c context.Context, arg *mcnmodel.McnGetRankAPIReq, rankFunc mcndao.RankFunc) (res *mcnmodel.McnGetRankUpFansReply, err error) {
	v, err := rankFunc(arg.SignID)
	if err != nil || v == nil {
		log.Error("get rank fail, sign id=%d, err=%s", arg.SignID, err)
		return
	}

	res = new(mcnmodel.McnGetRankUpFansReply)
	res.Result = v.GetList(arg.Tid, arg.DataType)
	res.TypeList = v.GetTypeList(arg.DataType)
	return
}

// GetMcnSummaryAPI .
func (s *Service) GetMcnSummaryAPI(c context.Context, arg *mcnmodel.McnGetDataSummaryReq) (res *mcnmodel.McnGetDataSummaryReply, err error) {
	return s.datadao.GetMcnSummaryCache(c, arg.SignID, datamodel.GetLastDay())
}

// GetIndexIncAPI .
func (s *Service) GetIndexIncAPI(c context.Context, arg *mcnmodel.McnGetIndexIncReq) (res *mcnmodel.McnGetIndexIncReply, err error) {
	return s.datadao.GetIndexIncCache(c, arg.SignID, datamodel.GetLastDay(), arg.Type)
}

// GetIndexSourceAPI .
func (s *Service) GetIndexSourceAPI(c context.Context, arg *mcnmodel.McnGetIndexSourceReq) (res *mcnmodel.McnGetIndexSourceReply, err error) {
	return s.datadao.GetIndexSourceCache(c, arg.SignID, datamodel.GetLastDay(), arg.Type)
}

// GetPlaySourceAPI .
func (s *Service) GetPlaySourceAPI(c context.Context, arg *mcnmodel.McnGetPlaySourceReq) (res *mcnmodel.McnGetPlaySourceReply, err error) {
	return s.datadao.GetPlaySourceCache(c, arg.SignID, datamodel.GetLastDay())
}

// GetMcnFansAPI .
func (s *Service) GetMcnFansAPI(c context.Context, arg *mcnmodel.McnGetMcnFansReq) (res *mcnmodel.McnGetMcnFansReply, err error) {
	return s.datadao.GetMcnFansCache(c, arg.SignID, datamodel.GetLastDay())
}

// GetMcnFansIncAPI .
func (s *Service) GetMcnFansIncAPI(c context.Context, arg *mcnmodel.McnGetMcnFansIncReq) (res *mcnmodel.McnGetMcnFansIncReply, err error) {
	return s.datadao.GetMcnFansIncCache(c, arg.SignID, datamodel.GetLastDay())
}

// GetMcnFansDecAPI .
func (s *Service) GetMcnFansDecAPI(c context.Context, arg *mcnmodel.McnGetMcnFansDecReq) (res *mcnmodel.McnGetMcnFansDecReply, err error) {
	return s.datadao.GetMcnFansDecCache(c, arg.SignID, datamodel.GetLastDay())
}

// GetMcnFansAttentionWayAPI .
func (s *Service) GetMcnFansAttentionWayAPI(c context.Context, arg *mcnmodel.McnGetMcnFansAttentionWayReq) (res *mcnmodel.McnGetMcnFansAttentionWayReply, err error) {
	return s.datadao.GetMcnFansAttentionWayCache(c, arg.SignID, datamodel.GetLastDay())
}

// GetFansBaseFansAttrAPI .
func (s *Service) GetFansBaseFansAttrAPI(c context.Context, arg *mcnmodel.McnGetBaseFansAttrReq) (res *mcnmodel.McnGetBaseFansAttrReply, err error) {
	return s.datadao.GetFansBaseFansAttrCache(c, arg.SignID, datamodel.GetLastWeek(), arg.UserType)
}

// GetFansAreaAPI .
func (s *Service) GetFansAreaAPI(c context.Context, arg *mcnmodel.McnGetFansAreaReq) (res *mcnmodel.McnGetFansAreaReply, err error) {
	return s.datadao.GetFansAreaCache(c, arg.SignID, datamodel.GetLastWeek(), arg.UserType)
}

// GetFansTypeAPI .
func (s *Service) GetFansTypeAPI(c context.Context, arg *mcnmodel.McnGetFansTypeReq) (res *mcnmodel.McnGetFansTypeReply, err error) {
	return s.datadao.GetFansTypeCache(c, arg.SignID, datamodel.GetLastWeek(), arg.UserType)
}

// GetFansTagAPI .
func (s *Service) GetFansTagAPI(c context.Context, arg *mcnmodel.McnGetFansTagReq) (res *mcnmodel.McnGetFansTagReply, err error) {
	return s.datadao.GetFansTagCache(c, arg.SignID, datamodel.GetLastWeek(), arg.UserType)
}
