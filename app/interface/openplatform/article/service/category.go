package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/interface/openplatform/article/conf"
	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
)

func (s *Service) loadCategoriesproc() {
	for {
		time.Sleep(time.Minute * 10)
		s.loadCategories()
	}
}

func (s *Service) loadCategories() {
	for {
		c, err := s.dao.Categories(context.TODO())
		if err != nil || len(c) == 0 {
			dao.PromError("service:获取分类")
			time.Sleep(time.Second)
			continue
		}
		s.primaryCategories = transformPrimaryCategory(c)
		removeCategoryBanner(c)
		s.categoriesMap = c
		s.Categories = transformCategory(c)
		s.categoriesReverseMap = transformReverseCategory(c)
		s.categoryParents = transformCategoryParents(c)
		return
	}
}

func removeCategoryBanner(cs map[int64]*artmdl.Category) {
	for _, c := range cs {
		c.BannerURL = ""
	}
}

// transformPrimaryCategory 生成一级分区列表
func transformPrimaryCategory(cs map[int64]*artmdl.Category) (res []*artmdl.Category) {
	res = make([]*artmdl.Category, 0, len(cs))
	for _, c := range cs {
		if c.ParentID == 0 {
			nc := new(artmdl.Category)
			*nc = *c
			res = append(res, nc)
		}
	}
	sort.Sort(artmdl.Categories(res))
	if len(res) > 4 {
		res[4].Name = conf.Conf.Article.AppCategoryName
		res[4].BannerURL = conf.Conf.Article.AppCategoryURL
	}
	return
}

// 生成一级分类含children的数组
func transformCategory(cs map[int64]*artmdl.Category) (res artmdl.Categories) {
	newCs := make(map[int64]*artmdl.Category)
	for k, c := range cs {
		n := *c
		newCs[k] = &n
	}
	m := make(map[int64][]*artmdl.Category)
	for _, c := range newCs {
		m[c.ParentID] = append(m[c.ParentID], c)
	}
	for _, x := range m {
		sort.Sort(artmdl.Categories(x))
	}
	res = m[0]
	for id, cate := range newCs {
		if len(m[id]) != 0 {
			cate.Children = m[id]
		}
	}
	return
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

// 生成每个分区的父分区列表(含本分区)
func transformCategoryParents(cs map[int64]*artmdl.Category) (res map[int64][]*artmdl.Category) {
	res = make(map[int64][]*artmdl.Category)
	for _, c := range cs {
		categories := []*artmdl.Category{c}
		parentID := c.ParentID
		for parentID != 0 {
			if cs[parentID] != nil {
				categories = append(categories, cs[parentID])
				parentID = cs[parentID].ParentID
				continue
			}
			break
		}
		//reverse categories
		for i := len(categories)/2 - 1; i >= 0; i-- {
			opp := len(categories) - 1 - i
			categories[i], categories[opp] = categories[opp], categories[i]
		}
		res[c.ID] = categories
	}
	return
}
