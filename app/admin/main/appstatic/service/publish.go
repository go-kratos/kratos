package service

import (
	"context"

	"go-common/app/admin/main/appstatic/model"
	"go-common/library/log"
)

// GendDiff picks the already generated diff packages
func (s *Service) GendDiff(resID int) (generated map[int64]int64, err error) {
	generated = make(map[int64]int64)
	genVers := []*model.Ver{}
	if err = s.DB.Where("file_type IN (1,2)"). // 1=diff pkg, 2=diff pkg calculation in progress
							Where("is_deleted = 0").Where("resource_id = ?", resID).Select("id, from_ver").Find(&genVers).Error; err != nil {
		log.Error("generatedDiff Error %v", err)
		return
	}
	for _, v := range genVers {
		generated[v.FromVer] = v.ID
	}
	return
}

// Publish returns the second trigger result
func (s *Service) Publish(ctx context.Context, resID int) (data *model.PubResp, err error) {
	var (
		prodVers, testVers []int64 // the history versions that we should generate for
		currRes            *model.Resource
		generated          map[int64]int64
		prodMore, testMore []int64
	)
	// pick history versions to calculate diff
	if prodVers, testVers, currRes, err = s.pickDiff(resID); err != nil {
		return
	}
	// pick already generated diff packages
	if generated, err = s.GendDiff(resID); err != nil {
		return
	}
	// filter already generated
	for _, v := range prodVers {
		if _, ok := generated[v]; !ok {
			prodMore = append(prodMore, v)
		}
	}
	for _, v := range testVers {
		if _, ok := generated[v]; !ok {
			testMore = append(testMore, v)
		}
	}
	// put diff packages in our DB
	if err = s.putDiff(resID, mergeSlice(prodMore, testMore), currRes); err != nil {
		return
	}
	data = &model.PubResp{
		CurrVer:  currRes.Version,
		DiffProd: prodMore,
		DiffTest: testMore,
	}
	// add the publish resID into to push list
	err = s.newPush(ctx, resID)
	return
}
