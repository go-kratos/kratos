package service

import (
	"context"

	"go-common/app/interface/main/ugcpay/model"
	"go-common/library/log"
)

// IncomeAssetOverview 获得收入总览数据
func (s *Service) IncomeAssetOverview(ctx context.Context, mid int64) (inc *model.IncomeAssetOverview, err error) {
	inc, err = s.dao.IncomeAssetOverview(ctx, mid)
	return
}

// IncomeAssetList 获得稿件维度的收入数据
func (s *Service) IncomeAssetList(ctx context.Context, mid int64, ver int64, ps, pn int64) (inc *model.IncomeAssetMonthly, page *model.Page, err error) {
	if inc, err = s.dao.IncomeUserAssetList(ctx, mid, ver, ps, pn); err != nil {
		return
	}

	// 获得稿件标题信息
	var (
		aids     = make([]int64, 0)
		titleMap = make(map[int64]string)
	)
	for _, l := range inc.List {
		aids = append(aids, l.OID)
	}
	if titleMap, err = s.dao.ArchiveTitles(ctx, aids); err != nil {
		log.Error("s.dao.ArchiveTitles aids: %+v, err: %+v", aids, err)
		err = nil
	}
	for _, l := range inc.List {
		l.Title = titleMap[l.OID]
	}

	page = &model.Page{
		Size:  inc.Page.Size,
		Num:   inc.Page.Num,
		Total: inc.Page.Total,
	}
	return
}
