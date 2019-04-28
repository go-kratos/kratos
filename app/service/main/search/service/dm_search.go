package service

import (
	"context"
	"fmt"

	"go-common/app/service/main/search/dao"
	"go-common/app/service/main/search/model"
	"go-common/library/ecode"
)

func (s *Service) DmSearch(c context.Context, sp *model.DmSearchParams) (res *model.SearchResult, err error) {
	if res, err = s.dao.DmSearch(c, sp); err != nil {
		dao.PromError(fmt.Sprintf("es:%s 搜索dm_search失败", sp.Bsp.AppID), "s.dao.DmSearch(%v) error(%v)", sp, err)
		err = ecode.SearchDmFailed
	}
	return
}
