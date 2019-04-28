package service

import (
	"context"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
)

func (s *Service) isWhiteMid(mid int64) (uid int64, ok bool) {
	uid, ok = s.whiteMidsCache[mid]
	return
}

func (s *Service) bindBgmWithVideo(aid, mid int64, nvs []*archive.Video) (err error) {
	for _, v := range nvs {
		log.Info("bindBgmWithVideo aid(%d),mid(%d),v.sid(%d),v.cid(%d)", aid, mid, v.Sid, v.Cid)
		if v.Sid > 0 {
			s.bgm.Bind(context.TODO(), aid, v.Sid, v.Cid)
		}
	}
	return
}
