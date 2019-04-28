package service

import (
	"context"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) checkResTag(c context.Context, oid, tid int64, op, typ int32) (rtm map[int64]*model.Resource, err error) {
	if rtm, err = s.resTags(c, oid, typ); err != nil {
		return
	}
	if len(rtm) == 0 {
		if op == model.SpamActionDel {
			err = ecode.TagArcTagNotExist
		}
		return
	}
	if op == model.SpamActionAdd && len(rtm) >= s.conf.Tag.ResTagMaxNum {
		err = ecode.TagArcTagMaxNum
		return
	}
	_, ok := rtm[tid]
	if op == model.SpamActionAdd && ok {
		err = ecode.TagArcTagExist
	} else if op == model.SpamActionDel && !ok {
		err = ecode.TagArcTagNotExist
	}
	return
}

func (s *Service) checkUserAdd(c context.Context, oid, tid, author, mid int64, typ int32, ip string) (err error) {
	if s.noLimit(author, mid) {
		return
	}
	if typ == int32(model.ResTypeArchive) {
		if attr, ok := s.limitRes[oid]; ok {
			if (attr>>model.LimitAttrAllowAdd)&0x1 == 1 { // restag属性 不可绑定tag
				err = ecode.TagArcCannotAddTag
				return
			}
		}
	}
	// TODO 审核添加限制
	//var id int
	//if id, err = s.dao.ReportStatus(c, oid, tid, typ); err != nil {
	//	return
	//}
	//if id == 0 || id == 5 || id == 4 {
	//	err = ecode.TagAddNotRptPassed
	//	return
	//}
	return
}

func (s *Service) checkUserDel(c context.Context, oid, tid, author, mid int64, typ int32, ip string) (err error) {
	if s.noLimit(author, mid) {
		return
	}
	// 稿件属性check
	if typ == model.ResTypeArchive {
		if attr, ok := s.limitRes[oid]; ok {
			if (attr>>model.LimitAttrAllowDel)&0x1 == 1 { // restag属性 不可删除tag
				err = ecode.TagArcCannotAddTag
				return
			}
		}
	}

	// TODO 审核添加限制
	//var id int
	//if id, err = s.dao.ReportStatus(c, oid, tid, typ); err != nil {
	//	return
	//}
	//if id == 0 || id == 5 || id == 4 {
	//	err = ecode.TagDelNotRptPassed
	//	return
	//}
	// 官方活动tag不可删
	var tag *model.Tag
	if tag, err = s.tag(c, tid); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	if tag.Type == model.TypeOfficailActivity {
		err = ecode.TagIsOfficailTag
		return
	}
	// up主tag不可删
	var rt *model.Resource
	if rt, err = s.resource(c, oid, tid, int32(typ)); err != nil {
		return
	}
	if rt.Role == model.ResRoleUpper && mid != author {
		err = ecode.TagUpTagCannotDel
		return
	}
	if rt.Locked() {
		err = ecode.TagArcTagisLocked
	}
	return
}

func (s *Service) noLimit(author, mid int64) (ok bool) {
	if author == mid {
		return true
	}
	_, ok = s.whiteUser[mid]
	return
}

func (s *Service) checkResType(c context.Context, oid int64, typ int32, ip string) (res *api.Arc, err error) {
	switch typ {
	case model.ResTypeArchive:
		arg := &archive.ArgAid2{Aid: oid, RealIP: ip}
		if res, err = s.arcRPC.Archive3(c, arg); err != nil {
			if ecode.Cause(err).Code() == ecode.NothingFound.Code() {
				return nil, ecode.NothingFound
			}
			return nil, err
		}
		if res == nil || res.Aid == 0 {
			err = ecode.ArchiveNotExist
			return
		}
		if !res.IsNormal() {
			err = ecode.TagOperateFail
		}
		return
	case model.ResTypeArticle:

	default:
		log.Warn("this is resource is not found oid:%d,type:%d", oid, typ)
		return nil, ecode.NothingFound
	}
	return
}
