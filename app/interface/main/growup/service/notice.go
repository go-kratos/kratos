package service

import (
	"context"
	"fmt"

	"go-common/app/interface/main/growup/model"
)

// LatestNotice latest notice
func (s *Service) LatestNotice(c context.Context, platform int) (notice *model.Notice, err error) {
	return s.dao.LatestNotice(c, platform)
}

// GetNotices get notice
func (s *Service) GetNotices(c context.Context, typ int, platform int, offset, limit int) (notices []*model.Notice, total int64, err error) {
	typStr := ""
	if typ > 0 {
		typStr = fmt.Sprintf("type=%d AND", typ)
	}
	total, err = s.dao.NoticeCount(c, typStr, platform)
	if err != nil {
		return
	}
	notices, err = s.dao.Notices(c, typStr, platform, offset, limit)
	return
}
