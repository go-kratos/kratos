package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/search/dao"
	"go-common/app/admin/main/search/model"
	"go-common/library/ecode"
)

// ArchiveCheck gets archive check.
func (s *Service) ArchiveCheck(c context.Context, sp *model.ArchiveCheckParams) (res *model.SearchResult, err error) {
	if res, err = s.dao.ArchiveCheck(c, sp); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索archivecheck失败", sp.Bsp.AppID), "s.dao.SearchArchiveCheck(%v) error(%v) ", sp, err)
		err = ecode.SearchArchiveCheckFailed
	}
	return
}

// Video gets video relation.
func (s *Service) Video(c context.Context, sp *model.VideoParams) (res *model.SearchResult, err error) {
	if res, err = s.dao.Video(c, sp); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索video失败", sp.Bsp.AppID), "s.dao.Video(%v) error(%v) ", sp, err)
		err = ecode.SearchVideoFailed
	}
	return
}

// TaskQa .
func (s *Service) TaskQa(c context.Context, sp *model.TaskQa) (res *model.SearchResult, err error) {
	if res, err = s.dao.TaskQa(c, sp); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索TaskQa失败", sp.Bsp.AppID), "s.dao.TaskQa(%v) error(%v) ", sp, err)
		err = ecode.SearchVideoFailed
	}
	return
}

// ArchiveCommerce .
func (s *Service) ArchiveCommerce(c context.Context, sp *model.ArchiveCommerce) (res *model.SearchResult, err error) {
	if res, err = s.dao.ArchiveCommerce(c, sp); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索ArchiveCommerce失败", sp.Bsp.AppID), "s.dao.TaskQa(%v) error(%v) ", sp, err)
		err = ecode.SearchVideoFailed
	}
	return
}
