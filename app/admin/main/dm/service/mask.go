package service

import (
	"context"

	"go-common/app/admin/main/dm/model"
	accountApi "go-common/app/service/main/account/api"
	"go-common/app/service/main/archive/api"
	arcMdl "go-common/app/service/main/archive/model/archive"

	"go-common/library/ecode"
	"go-common/library/log"
)

// MaskState get mask state
func (s *Service) MaskState(c context.Context, tp int32, oid int64) (open, mobile, web int32, err error) {
	sub, err := s.dao.Subject(c, tp, oid)
	if err != nil {
		return
	}
	if sub == nil {
		err = ecode.ArchiveNotExist
		return
	}
	return sub.AttrVal(model.AttrSubMaskOpen), sub.AttrVal(model.AttrSubMblMaskReady), sub.AttrVal(model.AttrSubWebMaskReady), err
}

// UpdateMaskState update mask state
func (s *Service) UpdateMaskState(c context.Context, tp int32, oid int64, plat int8, state int32) (err error) {
	var (
		archive3 *api.Arc
		err1     error
		duration int64
		typeID   int32
	)
	sub, err := s.dao.Subject(c, tp, oid)
	if err != nil {
		return
	}
	if sub == nil {
		err = ecode.ArchiveNotExist
		return
	}
	if archive3, err1 = s.arcRPC.Archive3(c, &arcMdl.ArgAid2{Aid: sub.Pid}); err1 == nil && archive3 != nil {
		duration = archive3.Duration
		typeID = archive3.TypeID
	}
	if state == model.AttrYes {
		switch plat {
		case model.MaskPlatMbl:
			if sub.AttrVal(model.AttrSubMblMaskReady) == model.AttrNo {
				if err = s.dao.GenerateMask(c, oid, sub.Mid, plat, 0, sub.Pid, duration, typeID); err != nil {
					return
				}
			}
		case model.MaskPlatWeb:
			if sub.AttrVal(model.AttrSubWebMaskReady) == model.AttrNo {
				if err = s.dao.GenerateMask(c, oid, sub.Mid, plat, 0, sub.Pid, duration, typeID); err != nil {
					return
				}
			}
		default:
			if sub.AttrVal(model.AttrSubMblMaskReady) == model.AttrNo || sub.AttrVal(model.AttrSubWebMaskReady) == model.AttrNo {
				if err = s.dao.GenerateMask(c, oid, sub.Mid, plat, 0, sub.Pid, duration, typeID); err != nil {
					return
				}
			}
		}
	}
	sub.AttrSet(state, model.AttrSubMaskOpen)
	_, err = s.dao.UpSubjectAttr(c, tp, oid, sub.Attr)
	return
}

// GenerateMask generate mask
func (s *Service) GenerateMask(c context.Context, tp int32, oid int64, plat int8) (err error) {
	var (
		archive3 *api.Arc
		err1     error
		duration int64
		typeID   int32
	)
	sub, err := s.dao.Subject(c, tp, oid)
	if err != nil {
		return
	}
	if sub == nil {
		err = ecode.ArchiveNotExist
		return
	}
	if archive3, err1 = s.arcRPC.Archive3(c, &arcMdl.ArgAid2{Aid: sub.Pid}); err1 == nil && archive3 != nil {
		duration = archive3.Duration
		typeID = archive3.TypeID
	}
	err = s.dao.GenerateMask(c, oid, sub.Mid, plat, 1, sub.Pid, duration, typeID)
	return
}

// MaskUps get mask up infos.
func (s *Service) MaskUps(c context.Context, pn, ps int64) (res *model.MaskUpRes, err error) {

	MaskUps, total, err := s.dao.MaskUps(c, pn, ps)
	if err != nil {
		return
	}
	mids := make([]int64, 0, len(MaskUps))
	for _, up := range MaskUps {
		mids = append(mids, up.Mid)
	}
	arg := &accountApi.MidsReq{Mids: mids}
	uInfos, err := s.accountRPC.Infos3(c, arg)
	if err != nil {
		log.Error("s.accRPC.Infos3(%v) error(%v)", mids, err)
		return
	}
	for _, up := range MaskUps {
		if info, ok := uInfos.GetInfos()[up.Mid]; ok {
			up.Name = info.Name
		}
	}
	res = &model.MaskUpRes{
		Result: MaskUps,
		Page: &model.PageInfo{
			Num:   pn,
			Size:  ps,
			Total: total,
		},
	}
	return
}

// MaskUpOpen add mask up
func (s *Service) MaskUpOpen(c context.Context, mids []int64, state int32, comment string) (err error) {
	midMap := make(map[int64]struct{})
	ids := make([]int64, 0, len(mids))
	for _, mid := range mids {
		if _, ok := midMap[mid]; !ok {
			midMap[mid] = struct{}{}
			ids = append(ids, mid)
		}
	}
	// 验证mids
	arg := &accountApi.MidsReq{Mids: ids}
	uInfos, err := s.accountRPC.Infos3(c, arg)
	if err != nil {
		log.Error("s.accRPC.Infos3(%v) error(%v)", ids, err)
		return
	}
	if len(uInfos.GetInfos()) < len(ids) {
		err = ecode.AccountInexistence
		log.Error("s.MaskUpOpen length diff(%d,%d)", len(ids), len(uInfos.GetInfos()))
		return
	}
	for _, id := range ids {
		if _, err = s.dao.MaskUpOpen(c, id, state, comment); err != nil {
			return
		}
	}
	return
}
