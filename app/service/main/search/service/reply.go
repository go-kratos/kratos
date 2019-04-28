package service

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/dao"
	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
)

// ReplyRecord gets reply record.
func (s *Service) ReplyRecord(c context.Context, sp *model.ReplyRecordParams) (res *model.SearchResult, err error) {
	if res, err = s.dao.ReplyRecord(c, sp); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索replyrecord失败", sp.Bsp.AppID), "s.dao.ReplyRecord(%v,%d,%d) error(%v) ", sp, sp.Bsp.Pn, sp.Bsp.Ps, err)
		err = ecode.SearchReplyRecordFailed
	}
	return
}
