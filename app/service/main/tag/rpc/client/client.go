package tag

import (
	"context"

	"go-common/app/service/main/tag/model"
	"go-common/library/net/rpc"
)

const (
	_infoByID            = "RPC.InfoByID"
	_infoByName          = "RPC.InfoByName"
	_checkName           = "RPC.CheckName"
	_count               = "RPC.Count"
	_counts              = "RPC.Counts"
	_infoByIDs           = "RPC.InfoByIDs"
	_infoByNames         = "RPC.InfoByNames"
	_subTags             = "RPC.SubTags"
	_resTags             = "RPC.ResTags"
	_resTagLog           = "RPC.ResTagLog"
	_addCustomSubTag     = "RPC.AddCustomSubTag"
	_resCustomSubTag     = "RPC.CustomSubTag"
	_resAction           = "RPC.ResAction"
	_resActionMap        = "RPC.ResActionMap"
	_addSub              = "RPC.AddSub"
	_cancelSub           = "RPC.CancelSub"
	_like                = "RPC.Like"
	_hate                = "RPC.Hate"
	_hideTag             = "RPC.HideTag"
	_createTag           = "RPC.CreateTag"
	_createTags          = "RPC.CreateTags"
	_platformUpBind      = "RPC.PlatformUpBind"
	_platformAdminBind   = "RPC.PlatformAdminBind"
	_platformUserBind    = "RPC.PlatformUserBind"
	_reportAction        = "RPC.ReportAction"
	_limitResource       = "RPC.LimitResource"
	_tagGroup            = "RPC.TagGroup"
	_whiteUser           = "RPC.WhiteUser"
	_resOidsByTid        = "RPC.ResOidsByTid"
	_recommandTag        = "RPC.RecommandTag"
	_rankingHot          = "RPC.RankingHot"
	_rankingRegion       = "RPC.RankingRegion"
	_rankingBangumi      = "RPC.RankingBangumi"
	_hots                = "RPC.Hots"
	_hotMap              = "RPC.HotMap"
	_prids               = "RPC.Prids"
	_rids                = "RPC.Rids"
	_defaultUpBind       = "RPC.DefaultUpBind"
	_defaultAdminBind    = "RPC.DefaultAdminBind"
	_addCustomSubChannel = "RPC.AddCustomSubChannel"
	_resCustomSubChannel = "RPC.CustomSubChannel"
)

var _noRes = &struct{}{}

const (
	_appid = "community.service.tag"
)

// Service .
type Service struct {
	client *rpc.Client2
}

// New .
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// InfoByID .
func (s *Service) InfoByID(c context.Context, arg *model.ArgID) (res *model.Tag, err error) {
	res = new(model.Tag)
	err = s.client.Call(c, _infoByID, arg, res)
	return
}

// InfoByName .
func (s *Service) InfoByName(c context.Context, arg *model.ArgName) (res *model.Tag, err error) {
	res = new(model.Tag)
	err = s.client.Call(c, _infoByName, arg, res)
	return
}

// CheckName .
func (s *Service) CheckName(c context.Context, arg *model.ArgCheckName) (res *model.Tag, err error) {
	res = new(model.Tag)
	err = s.client.Call(c, _checkName, arg, res)
	return
}

// Count .
func (s *Service) Count(c context.Context, arg *model.ArgID) (res *model.Count, err error) {
	res = new(model.Count)
	err = s.client.Call(c, _count, arg, res)
	return
}

// Counts .
func (s *Service) Counts(c context.Context, arg *model.ArgIDs) (res map[int64]*model.Count, err error) {
	err = s.client.Call(c, _counts, arg, &res)
	return
}

// InfoByIDs .
func (s *Service) InfoByIDs(c context.Context, arg *model.ArgIDs) (res []*model.Tag, err error) {
	err = s.client.Call(c, _infoByIDs, arg, &res)
	return
}

// InfoByNames .
func (s *Service) InfoByNames(c context.Context, arg *model.ArgNames) (res []*model.Tag, err error) {
	err = s.client.Call(c, _infoByNames, arg, &res)
	return
}

// ResTags .
func (s *Service) ResTags(c context.Context, arg *model.ArgResTags) (res []*model.Resource, err error) {
	err = s.client.Call(c, _resTags, arg, &res)
	return
}

// ResTagLog .
func (s *Service) ResTagLog(c context.Context, arg *model.ArgResTagLog) (res []*model.ResourceLog, err error) {
	err = s.client.Call(c, _resTagLog, arg, &res)
	return
}

// AddSub .
func (s *Service) AddSub(c context.Context, arg *model.ArgAddSub) (err error) {
	err = s.client.Call(c, _addSub, arg, _noRes)
	return
}

// CancelSub .
func (s *Service) CancelSub(c context.Context, arg *model.ArgCancelSub) (err error) {
	err = s.client.Call(c, _cancelSub, arg, _noRes)
	return
}

// SubTags .
func (s *Service) SubTags(c context.Context, arg *model.ArgSub) (res *model.ResSub, err error) {
	res = new(model.ResSub)
	err = s.client.Call(c, _subTags, arg, res)
	return
}

// AddCustomSubTag .
func (s *Service) AddCustomSubTag(c context.Context, arg *model.ArgCustomSub) (err error) {
	err = s.client.Call(c, _addCustomSubTag, arg, _noRes)
	return
}

// CustomSubTag .
func (s *Service) CustomSubTag(c context.Context, arg *model.ArgSub) (res *model.ResSubSort, err error) {
	res = new(model.ResSubSort)
	err = s.client.Call(c, _resCustomSubTag, arg, res)
	return
}

// AddCustomSubChannel .
func (s *Service) AddCustomSubChannel(c context.Context, arg *model.ArgCustomSub) (err error) {
	err = s.client.Call(c, _addCustomSubChannel, arg, _noRes)
	return
}

// CustomSubChannel .
func (s *Service) CustomSubChannel(c context.Context, arg *model.ArgCustomChannel) (res *model.ResSubSort, err error) {
	res = new(model.ResSubSort)
	err = s.client.Call(c, _resCustomSubChannel, arg, res)
	return
}

// Like .
func (s *Service) Like(c context.Context, arg *model.ArgResAction) (err error) {
	err = s.client.Call(c, _like, arg, _noRes)
	return
}

// Hate .
func (s *Service) Hate(c context.Context, arg *model.ArgResAction) (err error) {
	err = s.client.Call(c, _hate, arg, _noRes)
	return
}

// ResAction .
func (s *Service) ResAction(c context.Context, arg *model.ArgResAction) (res int32, err error) {
	err = s.client.Call(c, _resAction, arg, &res)
	return
}

// ResActionMap .
func (s *Service) ResActionMap(c context.Context, arg *model.ArgResActions) (res map[int64]int32, err error) {
	err = s.client.Call(c, _resActionMap, arg, &res)
	return
}

// HideTag .
func (s *Service) HideTag(c context.Context, arg *model.ArgHide) (err error) {
	err = s.client.Call(c, _hideTag, arg, _noRes)
	return
}

// CreateTag .
func (s *Service) CreateTag(c context.Context, arg *model.ArgCreate) (err error) {
	err = s.client.Call(c, _createTag, arg, _noRes)
	return
}

// CreateTags .
func (s *Service) CreateTags(c context.Context, arg *model.ArgCreate) (err error) {
	err = s.client.Call(c, _createTags, arg, _noRes)
	return
}

// PlatformUpBind .
func (s *Service) PlatformUpBind(c context.Context, arg *model.ArgUPBind) (err error) {
	err = s.client.Call(c, _platformUpBind, arg, _noRes)
	return
}

// PlatformAdminBind .
func (s *Service) PlatformAdminBind(c context.Context, arg *model.ArgUPBind) (err error) {
	err = s.client.Call(c, _platformAdminBind, arg, _noRes)
	return
}

// PlatformUserBind .
func (s *Service) PlatformUserBind(c context.Context, arg *model.ArgUserBind) (err error) {
	err = s.client.Call(c, _platformUserBind, arg, _noRes)
	return
}

// ReportAction .
func (s *Service) ReportAction(c context.Context, arg *model.ArgReportAction) (err error) {
	err = s.client.Call(c, _reportAction, arg, _noRes)
	return
}

// LimitResource .
func (s *Service) LimitResource(c context.Context) (res []*model.ResourceLimit, err error) {
	err = s.client.Call(c, _limitResource, _noRes, &res)
	return
}

// WhiteUser .
func (s *Service) WhiteUser(c context.Context) (res map[int64]struct{}, err error) {
	err = s.client.Call(c, _whiteUser, _noRes, &res)
	return
}

// TagGroup .
func (s *Service) TagGroup(c context.Context) (res []*model.Synonym, err error) {
	err = s.client.Call(c, _tagGroup, _noRes, &res)
	return
}

// ResOidsByTid .
func (s *Service) ResOidsByTid(c context.Context, arg *model.ArgRes) (res []int64, err error) {
	err = s.client.Call(c, _resOidsByTid, arg, &res)
	return
}

// RecommandTag .
func (s *Service) RecommandTag(c context.Context) (res map[int64]map[string][]*model.UploadTag, err error) {
	err = s.client.Call(c, _recommandTag, _noRes, &res)
	return
}

// RankingHot .
func (s *Service) RankingHot(c context.Context) (res []*model.Tag, err error) {
	err = s.client.Call(c, _rankingHot, _noRes, &res)
	return
}

// RankingRegion .
func (s *Service) RankingRegion(c context.Context, arg *model.ArgRankingRegion) (res []*model.RankingRegion, err error) {
	err = s.client.Call(c, _rankingRegion, arg, &res)
	return
}

// RankingBangumi .
func (s *Service) RankingBangumi(c context.Context) (res *model.ResBangumi, err error) {
	res = new(model.ResBangumi)
	err = s.client.Call(c, _rankingBangumi, _noRes, res)
	return
}

// Hots .
func (s *Service) Hots(c context.Context, arg *model.ArgHots) (res []*model.HotTag, err error) {
	err = s.client.Call(c, _hots, arg, &res)
	return
}

// HotMap .
func (s *Service) HotMap(c context.Context) (res map[int16][]int64, err error) {
	err = s.client.Call(c, _hotMap, _noRes, &res)
	return
}

// Prids .
func (s *Service) Prids(c context.Context) (res []int64, err error) {
	err = s.client.Call(c, _prids, _noRes, &res)
	return
}

// Rids .
func (s *Service) Rids(c context.Context) (res map[int64]int64, err error) {
	err = s.client.Call(c, _rids, _noRes, &res)
	return
}

// DefaultUpBind .
func (s *Service) DefaultUpBind(c context.Context, arg *model.ArgDefaultBind) (err error) {
	err = s.client.Call(c, _defaultUpBind, arg, _noRes)
	return
}

// DefaultAdminBind .
func (s *Service) DefaultAdminBind(c context.Context, arg *model.ArgDefaultBind) (err error) {
	err = s.client.Call(c, _defaultAdminBind, arg, _noRes)
	return
}
