package service

import (
	"context"
	"sort"
	"time"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/app/job/openplatform/article/dao"
	"go-common/app/job/openplatform/article/model"
	"go-common/library/log"
)

var (
	_recommendCategory = int64(0)
)

// UpdateSort update sort
func (s *Service) UpdateSort(c context.Context) (err error) {
	ids := []int64{_recommendCategory}
	for id := range s.categoriesMap {
		ids = append(ids, id)
	}
	var (
		cache   bool
		allArts map[int64]map[int][][2]int64
	)
	for _, categoryID := range ids {
		for _, field := range artmdl.SortFields {
			if (categoryID == _recommendCategory) && (field != artmdl.FieldNew) {
				continue
			}
			if cache, err = s.dao.ExpireSortCache(c, categoryID, field); err != nil {
				dao.PromError("recommends:更新文章排序")
				return
			}
			// cache不存在才更新
			if cache {
				continue
			}
			dao.PromInfo("sort:初始化排序缓存")
			if allArts == nil {
				allArts, _ = s.loadSortArts(context.TODO())
			}
			if field == artmdl.FieldNew {
				s.updateNewArts(c, categoryID)
			} else {
				s.updateStatSort(c, allArts, categoryID, field)
			}
		}
	}
	return
}

func (s *Service) updateNewArts(c context.Context, categoryID int64) (err error) {
	var arts [][2]int64
	if categoryID == _recommendCategory {
		if arts, err = s.dao.NewestArtIDs(c, s.c.Job.MaxNewArtsNum); err == nil {
			log.Info("s.updateNewArts() len: %v", len(arts))
			// 不异步
			err = s.dao.AddSortCaches(c, _recommendCategory, artmdl.FieldNew, arts, s.c.Job.MaxNewArtsNum)
		}
		return
	}
	ids := []int64{}
	cs := s.categoriesReverseMap[categoryID]
	if len(cs) == 0 {
		// 二级分区
		ids = append(ids, categoryID)
	} else {
		// 一级分区聚合子分区
		for _, s := range cs {
			ids = append(ids, s.ID)
		}
	}
	if arts, err = s.dao.NewestArtIDByCategory(c, ids, s.c.Job.MaxNewArtsNum); err == nil {
		// 不异步
		err = s.dao.AddSortCaches(c, categoryID, artmdl.FieldNew, arts, s.c.Job.MaxNewArtsNum)
	}
	return
}

func (s *Service) updateStatSort(c context.Context, allArts map[int64]map[int][][2]int64, categoryID int64, field int) (err error) {
	if allArts == nil {
		return
	}
	if categoryID == _recommendCategory {
		//推荐下没有排序
		return
	}
	arts := trimArts(allArts[categoryID][field], int(s.c.Job.MaxSortArtsNum))
	err = s.dao.AddSortCaches(c, categoryID, field, arts, s.c.Job.MaxSortArtsNum)
	return
}

func trimArts(arts [][2]int64, max int) (res [][2]int64) {
	sort.Slice(arts, func(i, j int) bool {
		return arts[i][1] > arts[j][1]
	})
	if len(arts) > max {
		return arts[:max]
	}
	return arts
}

func (s *Service) loadSortArts(c context.Context) (res map[int64]map[int][][2]int64, err error) {
	// init category and field
	res = make(map[int64]map[int][][2]int64)
	for id := range s.categoriesMap {
		res[id] = make(map[int][][2]int64)
	}
	var (
		arts      []*model.SearchArticle
		limitTime = time.Now().Unix() - s.sortLimitTime
	)
	if arts, err = s.dao.SearchArts(c, limitTime); err != nil {
		dao.PromError("sort:初始化排序计数失败")
		return
	}
	for _, art := range arts {
		if artmdl.NoDistributeAttr(art.Attributes) || artmdl.NoRegionAttr(art.Attributes) {
			continue
		}
		if res[art.CategoryID] == nil {
			res[art.CategoryID] = make(map[int][][2]int64)
		}
		res[art.CategoryID][artmdl.FieldFav] = append(res[art.CategoryID][artmdl.FieldFav], [2]int64{art.ID, art.StatsFavorite})
		res[art.CategoryID][artmdl.FieldLike] = append(res[art.CategoryID][artmdl.FieldLike], [2]int64{art.ID, art.StatsLikes})
		res[art.CategoryID][artmdl.FieldReply] = append(res[art.CategoryID][artmdl.FieldReply], [2]int64{art.ID, art.StatsReply})
		res[art.CategoryID][artmdl.FieldView] = append(res[art.CategoryID][artmdl.FieldView], [2]int64{art.ID, art.StatsView})
		var parentID int64
		if id, ok := s.categoriesMap[art.CategoryID]; ok {
			parentID = id.ParentID
		} else {
			continue
		}
		if res[parentID] == nil {
			res[parentID] = make(map[int][][2]int64)
		}
		res[parentID][artmdl.FieldFav] = append(res[parentID][artmdl.FieldFav], [2]int64{art.ID, art.StatsFavorite})
		res[parentID][artmdl.FieldLike] = append(res[parentID][artmdl.FieldLike], [2]int64{art.ID, art.StatsLikes})
		res[parentID][artmdl.FieldReply] = append(res[parentID][artmdl.FieldReply], [2]int64{art.ID, art.StatsReply})
		res[parentID][artmdl.FieldView] = append(res[parentID][artmdl.FieldView], [2]int64{art.ID, art.StatsView})
	}
	return
}

func (s *Service) loadCategoriesproc() {
	for {
		time.Sleep(time.Minute * 10)
		s.loadCategories()
	}
}

func (s *Service) loadCategories() {
	for {
		c, err := s.articleRPC.CategoriesMap(context.TODO(), &artmdl.ArgIP{})
		if err != nil || len(c) == 0 {
			dao.PromError("service:获取分类")
			log.Error("s.articleRPC.CategoriesMap err %v", err)
			time.Sleep(time.Second)
			continue
		}
		s.categoriesMap = c
		s.categoriesReverseMap = transformReverseCategory(c)
		return
	}
}

// 生成某个分类下的所有子分类
func transformReverseCategory(cs map[int64]*artmdl.Category) (res map[int64][]*artmdl.Category) {
	res = make(map[int64][]*artmdl.Category)
	for _, c := range cs {
		n := c
		old := c
		for (n != nil) && (n.ParentID != 0) {
			res[n.ParentID] = append(res[n.ParentID], old)
			n = cs[n.ParentID]
		}
	}
	return
}
