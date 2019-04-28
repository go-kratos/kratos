package service

import (
	"context"

	dm2Mdl "go-common/app/interface/main/dm2/model"
	"go-common/app/interface/main/dm2/model/oplog"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/assist/model/assist"
	"go-common/library/ecode"
	"go-common/library/log"
)

// assist
func (s *Service) assist(c context.Context, mid int64, aid int64) (upID int64, isUp bool, err error) {
	var ares *assist.AssistRes
	arc, err := s.acvSvc.Archive3(c, &arcMdl.ArgAid2{Aid: aid})
	if err != nil {
		log.Error("s.acvSvc.Archive3(%d) error(%v)", aid, err)
		return
	}
	upID = arc.Author.Mid
	isUp = upID == mid
	if isUp {
		return
	}
	arg := &assist.ArgAssist{Mid: upID, AssistMid: mid, Type: assist.TypeDm}
	if ares, err = s.astSvc.Assist(c, arg); err != nil {
		log.Error("s.astSvc.Assist(%v) error(%v)", arg, err)
		return
	}
	if ares.Assist == 0 {
		err = ecode.DMAssistNo
		return
	}
	if ares.Allow < 1 {
		err = ecode.DMAssistLimit
	}
	return
}

// AssistBanned up主屏蔽
func (s *Service) AssistBanned(c context.Context, mid, cid int64, dmids []int64) (err error) {
	arg := &dm2Mdl.ArgBanUsers{
		Mid:   mid,
		Oid:   cid,
		DMIDs: dmids,
	}
	if err = s.dmRPC.BanUsers(c, arg); err != nil {
		log.Error("dmRPC.BanUsers(%+v) error(%v)", arg, err)
	}
	return
}

// AssistUptBanned 更新up主屏蔽
func (s *Service) AssistUptBanned(c context.Context, mid int64, hash string, active int8) (err error) {
	arg := &dm2Mdl.ArgEditUpFilters{
		Mid:     mid,
		Type:    dm2Mdl.FilterTypeID,
		Active:  active,
		Filters: []string{hash},
	}
	if _, err = s.dmRPC.EditUpFilters(c, arg); err != nil {
		log.Error("dmRPC.EditUpFilters(%+v) error(%v)", arg, err)
	}
	return
}

// AssistDelBanned2 批量撤销up主屏蔽
func (s *Service) AssistDelBanned2(c context.Context, mid, aid int64, hashes []string) (err error) {
	arg := &dm2Mdl.ArgCancelBanUsers{
		Mid:     mid,
		Aid:     aid,
		Filters: hashes,
	}
	if err = s.dmRPC.CancelBanUsers(c, arg); err != nil {
		log.Error("dmRPC.CancelBanUsers(%+v) error(%v)", arg, err)
	}
	return
}

// AssistBannedUsers 获取up主屏蔽列表
func (s *Service) AssistBannedUsers(c context.Context, mid, aid int64) (hashes []string, err error) {
	upID, _, err := s.assist(c, mid, aid)
	if err != nil {
		if err == ecode.DMAssistLimit {
			err = nil
		} else {
			log.Error("s.assist(%d,%d) error(%v)", mid, aid, err)
			return
		}
	}
	arg := &dm2Mdl.ArgUpFilters{Mid: upID}
	res, err := s.dmRPC.UpFilters(c, arg)
	if err != nil {
		log.Error("dmRPC.UpFilters(%+v) error(%v)", arg, err)
		return
	}
	for _, v := range res {
		if v.Type == dm2Mdl.FilterTypeID {
			hashes = append(hashes, v.Filter)
		}
	}
	return
}

// AssistDeleteDM  assist delete dm.
func (s *Service) AssistDeleteDM(c context.Context, mid, oid int64, dmids []int64) (err error) {
	arg := &dm2Mdl.ArgEditDMState{
		Type:         dm2Mdl.SubTypeVideo,
		Oid:          oid,
		Mid:          mid,
		State:        dm2Mdl.StateDelete, // must be this value
		Dmids:        dmids,
		Source:       oplog.SourcePlayer,
		OperatorType: oplog.OperatorMember,
	}
	if err = s.dmRPC.EditDMState(c, arg); err != nil {
		log.Error("dmRPC.EditDMState(%v) error(%v)", arg, err)
	}
	return
}
