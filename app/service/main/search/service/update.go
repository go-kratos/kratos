package service

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/dao"
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

func (s *Service) UpdateMap(c context.Context, esName string, bulkData []dao.BulkMapItem) (err error) {
	if err = s.dao.UpdateMapBulk(c, esName, bulkData); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 更新失败", esName), "s.dao.updateBulk error(%v) ", err)
		err = ecode.SearchUpdateIndexFailed
	}
	return
}
