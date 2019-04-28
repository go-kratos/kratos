package archive

import (
	"context"
	"go-common/app/interface/main/creative/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

// VideoPoints 获取视频最新生效的高能看点
func (s *Service) VideoPoints(c context.Context, aid, cid int64) (vp *archive.ViewPointRow, err error) {
	vp, err = s.arc.ViewPoint(c, aid, cid)
	return
}

// WebViewPoints 创作中心高能看点接口
func (s *Service) WebViewPoints(c context.Context, aid, cid, mid int64, ip string) (vps []*archive.ViewPointRow, err error) {
	a, err := s.arc.Archive(c, aid, ip)
	if err != nil {
		return
	}
	if a.Author.Mid != mid {
		err = ecode.AccessDenied
		log.Error("WebViewPoints() wrong auth. Login mid(%d) archive mid(%d)", mid, a.Author.Mid)
		return
	}
	if vps, err = s.arc.RawViewPoints(c, aid, cid, 3); err != nil {
		log.Error("s.arc.RawViewPoints(%d,%d) error(%v)", aid, cid, err)
		return
	}
	return
}
