package service

import (
	"fmt"
	"strconv"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

//ModulesAdd is used for add modules
func (s *Service) ModulesAdd(v *model.Modules) (err error) {
	var (
		order uint8
		mod   *model.Modules
	)
	if mod, err = s.isModulesExists(v.PageID, v.Title); err != nil {
		return
	}
	if mod != nil {
		err = fmt.Errorf("当前模块下，标题已存在")
		return
	}
	if order, err = s.getOrder(v.PageID); err != nil {
		return
	}
	//在已存在顺序上加一
	v.Order = order + 1
	if err = s.DB.Model(&model.Modules{}).Create(v).Error; err != nil {
		return
	}
	return
}

//isModulesExists is use for checking is module exists
func (s *Service) isModulesExists(pageID string, title string) (v *model.Modules, err error) {
	v = &model.Modules{}
	w := map[string]interface{}{
		"deleted": model.ModulesNotDelete,
		"page_id": pageID,
		"title":   title,
	}
	if err = s.DB.Where(w).First(v).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return
	}
	return
}

//getOrder is used for getting existed model order
func (s *Service) getOrder(pageID string) (order uint8, err error) {
	var v model.Modules
	if err = s.DB.Where("deleted = 0").Where("page_id = ?", pageID).Order("`order` DESC").First(&v).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		log.Error("getOrder Err %v", err)
		return
	}
	order = v.Order
	return
}

//ModulesList is used for get module list
func (s *Service) ModulesList(pageID string) (v []*model.Modules, err error) {
	selectStr := []string{
		"id",
		"title",
		"page_id",
		"source",
		"type",
		"flexible",
		"icon",
		"capacity",
		"more",
		"`order`",
		"moretype",
		"morepage",
		"valid",
		"src_type",
	}
	w := map[string]interface{}{
		"deleted": model.ModulesNotDelete,
		"page_id": pageID,
	}
	if err = s.DB.Where(w).Select(selectStr).Order("`order` ASC").Find(&v).Error; err != nil {
		return
	}
	for i := range v {
		attr := v[i]
		pid, _ := strconv.Atoi(attr.PageID)
		switch pid {
		case model.PageMain:
			attr.PageID = "主页"
		case model.PageJP:
			attr.PageID = "番剧"
		case model.PageMovie:
			attr.PageID = "电影"
		case model.PageDocumentary:
			attr.PageID = "纪录片"
		case model.PageCN:
			attr.PageID = "国创"
		case model.PageSoapopera:
			attr.PageID = "电视剧"
		}
		t, _ := strconv.Atoi(attr.Type)
		switch t {
		case model.TypeSevenFocus:
			attr.Type = "首页七格焦点图"
		case model.TypeFiveFocus:
			attr.Type = "5格焦点"
		case model.TypeSixFocus:
			attr.Type = "6格焦点"
		case model.TypeVertListFirst:
			attr.Type = "竖图1列表"
		case model.TypeVertListSecond:
			attr.Type = "竖图2列表"
		case model.TypeHorizList:
			attr.Type = "横图列表"
		case model.TypeZhuiFan:
			attr.Type = "追番模块"
		}
	}
	return
}

//ModulesEditGet is used for get module with module id
func (s *Service) ModulesEditGet(id uint64) (v *model.Modules, err error) {
	selectStr := []string{
		"id",
		"title",
		"page_id",
		"type",
		"source",
		"flexible",
		"icon",
		"capacity",
		"more",
		"`order`",
		"moretype",
		"morepage",
		"src_type",
	}
	w := map[string]interface{}{
		"id":      id,
		"deleted": model.ModulesNotDelete,
	}
	v = &model.Modules{}
	if err = s.DB.Where(w).Select(selectStr).First(v).Error; err != nil {
		return
	}
	return
}

//ModulesEditPost is used for update module value
func (s *Service) ModulesEditPost(id uint64, v *model.Modules) (err error) {
	var (
		mod = &model.Modules{}
	)
	if mod, err = s.isModulesExists(v.PageID, v.Title); err != nil {
		return
	}
	if mod != nil && mod.ID != id {
		err = fmt.Errorf("当前模块下，标题已存在")
		return
	}
	return s.DB.Model(&model.Modules{}).Where("id = ?", id).Update(v).Error
}

//GetModPub is used for get publish status from MC
func (s *Service) GetModPub(c *bm.Context, pageID string) (p model.ModPub, err error) {
	return s.dao.GetModPub(c, pageID)
}

//ModulesPublish is used for publish module or deleted modules
func (s *Service) ModulesPublish(c *bm.Context, pageID string, state uint8, ids []int, deletedIds []int) (err error) {
	if len(ids) > 30 {
		err = fmt.Errorf("模块发布不能超过30个")
		return
	}
	tx := s.DB.Begin()
	for k, v := range ids {
		up := map[string]interface{}{
			"order": k + 1,
			"valid": model.ModulesValid,
		}
		where := map[string]interface{}{
			"id":      v,
			"page_id": pageID,
		}
		if err = tx.Model(&model.Modules{}).Where(where).Update(up).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if len(deletedIds) > 0 {
		deletedUp := map[string]interface{}{
			"deleted": model.ModulesDelete,
		}
		if err = s.DB.Model(&model.Modules{}).Where("id in (?)", deletedIds).Where("page_id=?", pageID).
			Update(deletedUp).Error; err != nil {
			tx.Rollback()
			return
		}
	}
	if err = s.SetPublish(c, pageID, state); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

//SetPublish is used for set publish status
func (s *Service) SetPublish(c *bm.Context, pageID string, state uint8) (err error) {
	nowTime := time.Now()
	t := nowTime.Format("2006-01-02 15:04:05")
	p := model.ModPub{
		Time:  t,
		State: state,
	}
	return s.dao.SetModPub(c, pageID, p)
}

// TypeSupport distinguish whether the source is supported or not
func (s *Service) TypeSupport(srcType int, source int) bool {
	if srcType == _TypePGC {
		_, ok := s.supCatMap.PgcMap[int32(source)]
		return ok
	}
	if srcType == _TypeUGC {
		_, ok := s.supCatMap.UgcMap[int32(source)]
		return ok
	}
	return false
}

// loadCats, reload pgc & ugc support cats
func (s *Service) loadCats() {
	var (
		pgcTypes   = s.c.Cfg.SupportCat.PGCTypes
		ugcTypes   = s.c.Cfg.SupportCat.UGCTypes
		newCats    = []*model.ParentCat{}
		newCatsMap = &model.SupCats{
			UgcMap: make(map[int32]int),
			PgcMap: make(map[int32]int),
		}
	)
	// load supporting pgc types
	if len(pgcTypes) > 0 {
		for _, v := range pgcTypes {
			newCats = append(newCats, &model.ParentCat{
				ID:   v,
				Name: s.pgcCatToName(int(v)),
				Type: _TypePGC,
			})
			newCatsMap.PgcMap[v] = 1
		}
	}
	// load support ugc first level types and their children
	if len(ugcTypes) > 0 {
		for _, v := range ugcTypes {
			newCatsMap.UgcMap[v] = 1
			var ugcCat = &model.ParentCat{
				ID:   v,
				Type: _TypeUGC,
			}
			if tp, ok := s.ArcTypes[v]; ok {
				ugcCat.Name = tp.Name
			}
			for _, types := range s.ArcTypes { // gather children
				if types.Pid == v {
					ugcCat.Children = append(ugcCat.Children, &model.CommonCat{
						ID:   types.ID,
						Name: types.Name,
						PID:  types.Pid,
						Type: _TypeUGC,
					})
					newCatsMap.UgcMap[types.ID] = 1
				}
			}
			newCats = append(newCats, ugcCat)
		}
	}
	if len(newCats) > 0 {
		s.SupCats = newCats
		s.supCatMap = newCatsMap
	}
}
