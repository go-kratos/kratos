package service

import (
	"context"

	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/datamodel"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/ecode"
	"go-common/library/log"
)

//GetMcnGetIndexInc .
func (s *Service) GetMcnGetIndexInc(c context.Context, arg *mcnmodel.McnGetIndexIncReq) (res *mcnmodel.McnGetIndexIncReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}
	if res, err = s.datadao.GetIndexIncCache(c, mcnSign.ID, datamodel.GetLastDay(), arg.Type); err != nil {
		log.Error("get data fail, err=%v", err)
		return
	}
	return
}

// GetMcnGetIndexSource .
func (s *Service) GetMcnGetIndexSource(c context.Context, arg *mcnmodel.McnGetIndexSourceReq) (res *mcnmodel.McnGetIndexSourceReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetIndexSourceCache(c, mcnSign.ID, datamodel.GetLastDay(), arg.Type); err != nil {
		log.Error("get data fail, err=%v", err)
		return
	}

	return
}

// GetPlaySource .
func (s *Service) GetPlaySource(c context.Context, arg *mcnmodel.McnGetPlaySourceReq) (res *mcnmodel.McnGetPlaySourceReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetPlaySourceCache(c, mcnSign.ID, datamodel.GetLastDay()); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetMcnFans .
func (s *Service) GetMcnFans(c context.Context, arg *mcnmodel.McnGetMcnFansReq) (res *mcnmodel.McnGetMcnFansReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetMcnFansCache(c, mcnSign.ID, datamodel.GetLastDay()); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetMcnFansInc .
func (s *Service) GetMcnFansInc(c context.Context, arg *mcnmodel.McnGetMcnFansIncReq) (res *mcnmodel.McnGetMcnFansIncReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetMcnFansIncCache(c, mcnSign.ID, datamodel.GetLastDay()); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetMcnFansDec .
func (s *Service) GetMcnFansDec(c context.Context, arg *mcnmodel.McnGetMcnFansDecReq) (res *mcnmodel.McnGetMcnFansDecReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetMcnFansDecCache(c, mcnSign.ID, datamodel.GetLastDay()); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetMcnFansAttentionWay .
func (s *Service) GetMcnFansAttentionWay(c context.Context, arg *mcnmodel.McnGetMcnFansAttentionWayReq) (res *mcnmodel.McnGetMcnFansAttentionWayReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetMcnFansAttentionWayCache(c, mcnSign.ID, datamodel.GetLastDay()); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetBaseFansAttrReq .
func (s *Service) GetBaseFansAttrReq(c context.Context, arg *mcnmodel.McnGetBaseFansAttrReq) (res *mcnmodel.McnGetBaseFansAttrReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetFansBaseFansAttrCache(c, mcnSign.ID, datamodel.GetLastWeek(), arg.UserType); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetFansArea .
func (s *Service) GetFansArea(c context.Context, arg *mcnmodel.McnGetFansAreaReq) (res *mcnmodel.McnGetFansAreaReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetFansAreaCache(c, mcnSign.ID, datamodel.GetLastWeek(), arg.UserType); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetFansType .
func (s *Service) GetFansType(c context.Context, arg *mcnmodel.McnGetFansTypeReq) (res *mcnmodel.McnGetFansTypeReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetFansTypeCache(c, mcnSign.ID, datamodel.GetLastWeek(), arg.UserType); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}

// GetFansTag .
func (s *Service) GetFansTag(c context.Context, arg *mcnmodel.McnGetFansTagReq) (res *mcnmodel.McnGetFansTagReply, err error) {
	mcnSign, err := s.getMcnWithState(c, arg.McnMid, model.MCNSignStateOnSign)
	if err != nil {
		if err != ecode.NothingFound {
			log.Error("error get state, err=%s", err)
		}
		return
	}

	if res, err = s.datadao.GetFansTagCache(c, mcnSign.ID, datamodel.GetLastWeek(), arg.UserType); err != nil {
		log.Error("get data fail, err=%v", err)
	}

	return
}
