package goblin

import (
	"context"
	"fmt"

	"go-common/app/interface/main/tv/model"
	arcwar "go-common/app/service/main/archive/api"
	tvapi "go-common/app/service/main/tv/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_tvVipOk = 1
)

// UgcPlayurl returns the result of ugc play url
func (s *Service) UgcPlayurl(c context.Context, p *model.PlayURLReq, mid int64) (result map[string]interface{}, err error) {
	var (
		arc       *arcwar.Arc
		firstRes  map[string]interface{}
		firstResp *model.PlayURLResp
		tvRes     *tvapi.UserInfoReply
	)
	result = make(map[string]interface{})
	if firstRes, firstResp, err = s.dao.UgcPlayurl(c, p); err != nil {
		return
	}
	if _, ok := s.VipQns[p.Qn]; !ok { // if it doesn't request vip quality, let it go
		return firstRes, nil
	}
	if mid != 0 {
		if tvRes, err = s.tvCilent.UserInfo(c, &tvapi.UserInfoReq{Mid: mid}); err != nil && !ecode.EqualError(ecode.NothingFound, err) {
			log.Error("[playurl.UgcPlayurl] mid(%d) error(%s)", mid, err)
			return
		}
		if tvRes != nil && tvRes.Status == _tvVipOk {
			return firstRes, nil // if it's tv vip, let it go !
		}
		if arc, err = s.arcDao.Archive3(c, p.Avid); err != nil || arc == nil { // try author himself
			log.Warn("s.arcDao.Archive3 failed can not view Aid %d, Mid %", p.Avid, mid)
			return
		}
		if arc.Author.Mid == mid {
			return firstRes, nil // if it's upper himself, let it go
		}
	}
	// downgrade logic
	for _, qn := range firstResp.AcceptQuality {
		qnStr := fmt.Sprintf("%d", qn)
		if _, ok := s.VipQns[qnStr]; ok { // if vip
			continue
		}
		p.Qn = qnStr
		result, _, err = s.dao.UgcPlayurl(c, p)
		return
	}
	err = ecode.NothingFound // it doesn't have any other quality to allow downgrade
	log.Error("Allow Quality %v, Err %v", firstResp.AcceptQuality, err)
	return
}
