package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/search/dao"
	"go-common/library/ecode"
)

// Update update some indices.
func (s *Service) Update(c context.Context, esName string, bulkData []dao.BulkItem) (err error) {
	if err = s.dao.UpdateBulk(c, esName, bulkData); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 更新失败", esName), "s.dao.updateBulk error(%v) ", err)
		err = ecode.SearchUpdateIndexFailed
	}
	return
}

// MapUpdate map update.
func (s *Service) MapUpdate(c context.Context, p []dao.BulkMapItem) (err error) {
	err = s.dao.UpdateMapBulk(c, "ssd_archive", p)
	return
}

// Index .
func (s *Service) Index(c context.Context, esName string, bulkData []dao.BulkItem) (err error) {
	if err = s.dao.BulkIndex(c, esName, bulkData); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 写入失败", esName), "s.dao.BulkIndex error(%v) ", err)
		err = ecode.SearchUpdateIndexFailed
	}
	return
}
