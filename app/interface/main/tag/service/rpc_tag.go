package service

import (
	"context"
	"fmt"

	"go-common/app/interface/main/tag/model"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

func (s *Service) tag(c context.Context, tid, mid int64) (res *rpcModel.Tag, err error) {
	arg := &rpcModel.ArgID{
		ID:     tid,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.InfoByID(c, arg); err != nil {
		log.Error("s.tagRPC.InfoByID()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) tags(c context.Context, tids []int64, mid int64) (res []*rpcModel.Tag, err error) {
	arg := &rpcModel.ArgIDs{
		IDs:    tids,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.InfoByIDs(c, arg); err != nil {
		log.Error("s.tagRPC.InfoByIDs() ArgIDs:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) tagName(c context.Context, mid int64, name string) (res *rpcModel.Tag, err error) {
	arg := &rpcModel.ArgName{
		Name:   name,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.InfoByName(c, arg); err != nil {
		log.Error("s.tagRPC.InfoByName()ArgName:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) tagNames(c context.Context, mid int64, names []string) (res []*rpcModel.Tag, err error) {
	arg := &rpcModel.ArgNames{
		Names:  names,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.InfoByNames(c, arg); err != nil {
		log.Error("s.tagRPC.InfoByNames() ArgName:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) count(c context.Context, tid, mid int64) (res *rpcModel.Count, err error) {
	arg := &rpcModel.ArgID{
		ID:     tid,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.Count(c, arg); err != nil {
		log.Error("s.tagRPC.Count() ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) countMap(c context.Context, tids []int64, mid int64) (res map[int64]*rpcModel.Count, err error) {
	arg := &rpcModel.ArgIDs{
		IDs:    tids,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.Counts(c, arg); err != nil {
		log.Error("s.tagRPC.Counts() ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) subTag(c context.Context, mid int64, pn, ps, order int) (res *rpcModel.ResSub, err error) {
	arg := &rpcModel.ArgSub{
		Mid:    mid,
		Pn:     pn,
		Ps:     ps,
		Order:  order,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.SubTags(c, arg); err != nil {
		log.Error("s.tagRPC.SubTags()ArgSub:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) resTags(c context.Context, mid, oid int64, typ int32) (res []*rpcModel.Resource, err error) {
	arg := &rpcModel.ArgResTags{
		Oid:    oid,
		Type:   typ,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.ResTags(c, arg); err != nil {
		log.Error("s.tagRPC.ResTags()ArgResTags:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) resTagMap(c context.Context, mid, oid int64, typ int32) (rem map[int64]*rpcModel.Resource, err error) {
	arg := &rpcModel.ArgResTags{
		Oid:    oid,
		Type:   typ,
		Mid:    mid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	var res []*rpcModel.Resource
	if res, err = s.tagRPC.ResTags(c, arg); err != nil {
		log.Error("s.tagRPC.ResTags()ArgResTags:%+v, error(%v)", arg, err)
		return
	}
	rem = make(map[int64]*rpcModel.Resource)
	for _, r := range res {
		rem[r.Tid] = r
	}
	return
}

func (s *Service) resTagLog(c context.Context, mid, oid int64, typ int32, pn, ps int) (res []*rpcModel.ResourceLog, err error) {
	arg := &rpcModel.ArgResTagLog{
		Oid:    oid,
		Type:   typ,
		Mid:    mid,
		Pn:     pn,
		Ps:     ps,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.ResTagLog(c, arg); err != nil {
		log.Error("s.tagRPC.ResTagLog()ArgResTagLog:%+v, error(%v)", arg, err)
	}
	return
}
func (s *Service) resTagLogMap(c context.Context, mid, oid int64, typ int32, pn, ps int) (rem map[string]*rpcModel.ResourceLog, err error) {
	arg := &rpcModel.ArgResTagLog{
		Oid:    oid,
		Type:   typ,
		Mid:    mid,
		Pn:     pn,
		Ps:     ps,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	var res []*rpcModel.ResourceLog
	if res, err = s.tagRPC.ResTagLog(c, arg); err != nil {
		log.Error("s.tagRPC.ResTagLog()ArgResTagLog:%+v, error(%v)", arg, err)
	}
	rem = make(map[string]*rpcModel.ResourceLog, len(res))
	for _, v := range res {
		k := fmt.Sprintf("%d_%d_%d_%d_%d", v.Oid, v.Type, v.Tid, v.Mid, v.Action)
		rem[k] = v
	}
	return
}

func (s *Service) upCustomSubTags(c context.Context, mid int64, tp int, tids []int64) (err error) {
	arg := &rpcModel.ArgCustomSub{
		Mid:    mid,
		Type:   tp,
		Tids:   tids,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.AddCustomSubTag(c, arg); err != nil {
		log.Error("s.tagRPC.AddCustomSubTag()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) customSubTag(c context.Context, mid int64, tp, pn, ps, order int) (res *rpcModel.ResSubSort, err error) {
	arg := &rpcModel.ArgSub{
		Mid:    mid,
		Type:   tp,
		Pn:     pn,
		Ps:     ps,
		Order:  order,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.CustomSubTag(c, arg); err != nil {
		log.Error("s.tagRPC.CustomSubTag()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

// resTagsService rpc tag service user tag .
func (s *Service) resTagsService(c context.Context, aid, mid int64, typ int32) (reTs []*model.Tag, tids []int64, err error) {
	var res []*rpcModel.Resource
	res, err = s.resTags(c, mid, aid, typ)
	if err != nil {
		return
	}
	for _, v := range res {
		tids = append(tids, v.Tag.ID)
		reTs = append(reTs, &model.Tag{
			ID:      v.Tag.ID,
			Name:    v.Tag.Name,
			Cover:   v.Tag.Cover,
			Content: v.Tag.Content,
			Type:    int8(v.Tag.Type),
			State:   int8(v.Tag.State),
			CTime:   v.Tag.CTime,
			MTime:   v.Tag.MTime,
			// TODO Count
			IsAtten:   int8(v.Tag.Attention),
			Role:      int8(v.Role),
			Likes:     int64(v.Like),
			Hates:     int64(v.Hate),
			Attribute: int8(v.Attr),
		})
	}
	return
}

// arcTagsService rpc tag service user tag .
func (s *Service) arcTagsService(c context.Context, aid, mid int64, typ int32) (ats []*model.ArcTag, err error) {
	var res []*rpcModel.Resource
	res, err = s.resTags(c, mid, aid, typ)
	if err != nil {
		return
	}
	for _, v := range res {
		ats = append(ats, &model.ArcTag{
			Aid:       v.Oid,
			Mid:       v.Mid,
			Tid:       v.Tid,
			Likes:     int64(v.Like),
			Hates:     int64(v.Hate),
			Attribute: int8(v.Attr),
			Role:      int8(v.Role),
			State:     int8(v.State),
			CTime:     v.CTime,
			MTime:     v.MTime,
		})
	}
	return
}

func (s *Service) resAction(c context.Context, mid, oid int64, typ int32, tids []int64) (res map[int64]int32, err error) {
	arg := &rpcModel.ArgResActions{
		Mid:    mid,
		Oid:    oid,
		Type:   typ,
		Tids:   tids,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.ResActionMap(c, arg); err != nil {
		log.Error("s.tagRPC.Count() ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) addSub(c context.Context, mid int64, tids []int64) (err error) {
	arg := &rpcModel.ArgAddSub{
		Mid:    mid,
		Tids:   tids,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.AddSub(c, arg); err != nil {
		log.Error("s.tagRPC.AddSub()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) cancelSub(c context.Context, mid int64, tid int64) (err error) {
	arg := &rpcModel.ArgCancelSub{
		Mid:    mid,
		Tid:    tid,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.CancelSub(c, arg); err != nil {
		log.Error("s.tagRPC.CancelSub()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) likeService(c context.Context, mid, tid, oid int64, typ int32) (err error) {
	arg := &rpcModel.ArgResAction{
		Mid:    mid,
		Tid:    tid,
		Oid:    oid,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.Like(c, arg); err != nil {
		log.Error("s.tagRPC.Like()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) hateService(c context.Context, mid, tid, oid int64, typ int32) (err error) {
	arg := &rpcModel.ArgResAction{
		Mid:    mid,
		Tid:    tid,
		Oid:    oid,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.Hate(c, arg); err != nil {
		log.Error("s.tagRPC.Hate()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) createTag(c context.Context, tag *rpcModel.Tag) (err error) {
	arg := &rpcModel.ArgCreate{
		Tag:    tag,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.CreateTag(c, arg); err != nil {
		log.Error("s.tagRPC.CreateTag()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) createTags(c context.Context, tags []*rpcModel.Tag) (err error) {
	arg := &rpcModel.ArgCreate{
		Tags:   tags,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.CreateTags(c, arg); err != nil {
		log.Error("s.tagRPC.CreateTag()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

//
func (s *Service) platformUpBind(c context.Context, oid int64, mid int64, tids []int64, typ int32) (err error) {
	arg := &rpcModel.ArgUPBind{
		Mid:    mid,
		Tids:   tids,
		Oid:    oid,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.PlatformUpBind(c, arg); err != nil {
		log.Error("s.tagRPC.PlatformUpBind()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

//
func (s *Service) platformAdminBind(c context.Context, oid int64, mid int64, tids []int64, typ int32) (err error) {
	arg := &rpcModel.ArgUPBind{
		Mid:    mid,
		Tids:   tids,
		Oid:    oid,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.PlatformAdminBind(c, arg); err != nil {
		log.Error("s.tagRPC.PlatformAdminBind()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

// bindResTag .
func (s *Service) platformUserBind(c context.Context, oid int64, mid int64, tid int64, typ, role, action int32) (err error) {
	arg := &rpcModel.ArgUserBind{
		Oid:    oid,
		Mid:    mid,
		Tid:    tid,
		Type:   typ,
		Role:   role,
		Action: action,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.PlatformUserBind(c, arg); err != nil {
		log.Error("s.tagRPC.PlatformUserBind()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

// reportAction .
func (s *Service) reportAction(c context.Context, oid, logID, mid int64, typ, partID, reason, score int32) (err error) {
	arg := &rpcModel.ArgReportAction{
		Oid:    oid,
		Mid:    mid,
		LogID:  logID,
		Type:   typ,
		PartID: partID,
		Score:  score,
		Reason: reason,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.ReportAction(c, arg); err != nil {
		log.Error("s.tagRPC.ReportAction()ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) limitResource(c context.Context) (res []*rpcModel.ResourceLimit, err error) {
	if res, err = s.tagRPC.LimitResource(c); err != nil {
		log.Error("s.tagRPC.LimitResource() error(%v)", err)
	}
	return
}

func (s *Service) whiteUserService(c context.Context) (res map[int64]struct{}, err error) {
	if res, err = s.tagRPC.WhiteUser(c); err != nil {
		log.Error("s.tagRPC.WhiteUser() error(%v)", err)
	}
	return
}

func (s *Service) tagGroup(c context.Context) (res []*rpcModel.Synonym, err error) {
	if res, err = s.tagRPC.TagGroup(c); err != nil {
		log.Error("s.tagRPC.TagGroup() error(%v)", err)
	}
	return
}

func (s *Service) resOidsByTid(c context.Context, tid int64, typ int32) (res []int64, err error) {
	arg := &rpcModel.ArgRes{
		Tid:    tid,
		Limit:  1000,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if res, err = s.tagRPC.ResOidsByTid(c, arg); err != nil {
		log.Error("s.tagRPC.ResOidsByTid() error(%v)", err)
	}
	return
}

func (s *Service) recommandTagService(c context.Context) (res map[int64]map[string][]*rpcModel.UploadTag, err error) {
	if res, err = s.tagRPC.RecommandTag(c); err != nil {
		log.Error("s.tagRPC.RecommandTag() error(%v)", err)
	}
	return
}

func (s *Service) hots(c context.Context, rid, typ int64) (res []*model.HotTag, err error) {
	arg := &rpcModel.ArgHots{
		Rid:    rid,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	var hoTags []*rpcModel.HotTag
	if hoTags, err = s.tagRPC.Hots(c, arg); err != nil {
		log.Error("s.tagRPC.Hots()Hots:%+v, error(%v)", arg, err)
		return
	}
	for _, v := range hoTags {
		ht := &model.HotTag{
			Rid:       v.Rid,
			Tid:       v.Tid,
			Tname:     v.Tname,
			HighLight: v.HighLight,
			IsAtten:   v.IsAtten,
		}
		res = append(res, ht)
	}
	return
}

func (s *Service) hotMap(c context.Context) (res map[int16][]int64, err error) {
	if res, err = s.tagRPC.HotMap(c); err != nil {
		log.Error("s.tagRPC.HotMap() error(%v)", err)
	}
	return
}

func (s *Service) prids(c context.Context) (res []int64, err error) {
	if res, err = s.tagRPC.Prids(c); err != nil {
		log.Error("s.tagRPC.Prids() error(%v)", err)
	}
	return
}

func (s *Service) ridsService(c context.Context) (rids []int64, pridM map[int64]int64, ridMap map[int64][]int64, err error) {
	if pridM, err = s.tagRPC.Rids(c); err != nil {
		log.Error("s.tagRPC.Rids()  error(%v)", err)
		return
	}
	ridMap = make(map[int64][]int64)
	tmpMap := make(map[int64]struct{})
	for k, v := range pridM {
		ridMap[v] = append(ridMap[v], k)
		if _, ok := tmpMap[k]; !ok {
			rids = append(rids, k)
		}
		if _, ok := tmpMap[v]; !ok {
			rids = append(rids, v)
		}
	}
	return
}

func (s *Service) defaultUpBind(c context.Context, oid, mid int64, tids []int64, typ int32) (err error) {
	arg := &rpcModel.ArgDefaultBind{
		Mid:    mid,
		Tids:   tids,
		Oid:    oid,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.DefaultUpBind(c, arg); err != nil {
		log.Error("s.tagRPC.DefaultUpBind() ArgID:%+v, error(%v)", arg, err)
	}
	return
}

func (s *Service) defaultAdminBind(c context.Context, oid, mid int64, tids []int64, typ int32) (err error) {
	arg := &rpcModel.ArgDefaultBind{
		Mid:    mid,
		Tids:   tids,
		Oid:    oid,
		Type:   typ,
		RealIP: metadata.String(c, metadata.RemoteIP),
	}
	if err = s.tagRPC.DefaultAdminBind(c, arg); err != nil {
		log.Error("s.tagRPC.DefaultAdminBind() ArgID:%+v, error(%v)", arg, err)
	}
	return
}
