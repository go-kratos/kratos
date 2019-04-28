package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

// ugcLabels refreshes ugc labels
func (s *Service) ugcLabels() (err error) {
	var firstTps = s.firstTps()
	if err = s.ugcTpLabel(firstTps); err != nil {
		log.Error("ugcTpLabel Err %v", err)
		return
	}
	if err = s.ugcPubLabel(firstTps); err != nil {
		log.Error("ugcPubLabel Tps %v, Err %v", firstTps, err)
	}
	return
}

func (s *Service) ugcPubLabel(firstTps []int32) (err error) {
	var (
		exist bool
		cfg   = s.c.Cfg.RefLabel
	)
	for _, v := range firstTps { // check and create pub_time label
		pubtCore := &model.LabelCore{
			CatType:   model.UgcLabel,
			Category:  v,
			Param:     model.ParamUgctime,
			Value:     cfg.AllValue,
			Valid:     1,
			ParamName: cfg.UgcTime,
			Name:      cfg.AllName,
		}
		if exist, err = s.labelExist(pubtCore); err != nil {
			return
		}
		if !exist {
			if err = s.DB.Create(&model.LabelDB{LabelCore: *pubtCore}).Error; err != nil {
				log.Error("ugcLabels Pubtime All, Create Err %v", err)
				return
			}
		}
		time.Sleep(20 * time.Millisecond)
	}
	return
}

func (s *Service) firstTps() (tps []int32) {
	for _, v := range s.ArcTypes {
		if v.Pid == 0 {
			tps = append(tps, v.ID)
		}
	}
	return
}

func (s *Service) ugcTpLabel(firstTps []int32) (err error) {
	var (
		exist bool
		cfg   = s.c.Cfg.RefLabel
	)
	for _, v := range firstTps {
		extCore := &model.LabelCore{
			CatType:   model.UgcLabel,
			Param:     model.ParamTypeid,
			Valid:     1,
			Category:  v,
			Value:     cfg.AllValue,
			ParamName: cfg.UgcType,
			Name:      cfg.AllName,
		}
		if exist, err = s.labelExist(extCore); err != nil {
			return
		}
		if !exist {
			if err = s.DB.Create(&model.LabelDB{LabelCore: *extCore}).Error; err != nil {
				log.Error("ugcLabels ArcTypeID %d, Create ALL Label Err %v", v, err)
				return
			}
		}
	}
	for _, v := range s.ArcTypes {
		if v.Pid == 0 { // if first level type, we check the "ALL" label
			continue
		}
		extCore := &model.LabelCore{
			CatType:  model.UgcLabel,
			Param:    model.ParamTypeid,
			Valid:    1,
			Category: v.Pid,
			Value:    fmt.Sprintf("%d", v.ID),
		}
		if exist, err = s.labelExist(extCore); err != nil {
			return
		}
		if !exist {
			label := model.LabelDB{}
			label.FromArcTp(v, s.c.Cfg.RefLabel.UgcType)
			if err = s.DB.Create(&label).Error; err != nil {
				log.Error("ugcLabels ArcTypeID %d, Create Err %v", v.ID, err)
				return
			}
			time.Sleep(20 * time.Millisecond)
		}
	}
	return
}

// labelsExist distinguishes whether the ids are existing
func (s *Service) labelsExist(ids []int64) (err error) {
	var (
		labels []*model.LabelDB
		idmap  = make(map[int64]int)
	)
	if err = s.DB.Where(fmt.Sprintf("id IN (%s)", xstr.JoinInts(ids))).Where("deleted = 0").Find(&labels).Error; err != nil {
		log.Error("labelsExist Ids %v, Err %v", ids, err)
		return
	}
	if len(labels) >= len(ids) {
		return
	}
	for _, v := range labels {
		idmap[v.ID] = 1
	}
	for _, v := range ids {
		if _, ok := idmap[v]; !ok {
			log.Warn("labelsExist ids %v, not exist %d", ids, v)
			return ecode.RequestErr
		}
	}
	return
}

func (s *Service) labelExist(req *model.LabelCore) (exist bool, err error) {
	var (
		label = model.LabelDB{}
		db    = s.DB.Model(label).Where("deleted = 0")
	)
	if req.ID != 0 {
		db = db.Where("id = ?", req.ID)
	}
	if req.Category != 0 {
		db = db.Where("category = ?", req.Category)
	}
	if req.CatType != 0 {
		db = db.Where("cat_type = ?", req.CatType)
	}
	if req.Param != "" {
		db = db.Where("param = ?", req.Param)
	}
	if req.Value != "" {
		db = db.Where("value = ?", req.Value)
	}
	if req.Name != "" {
		db = db.Where("name = ?", req.Name)
	}
	if err = db.First(&label).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
			return
		}
		log.Error("labelExist V %v, Err %v", req, err)
		return
	}
	if label.ID > 0 {
		exist = true
	}
	return
}

func (s *Service) pgcLabels() (err error) {
	var (
		result *model.PgcCond
		exist  bool
	)
	for _, cat := range s.c.Cfg.SupportCat.PGCTypes {
		if result, err = s.dao.PgcCond(context.Background(), cat); err != nil {
			log.Error("PgcCond Cat %d, Err %v", cat, err)
			return
		}
		for _, cond := range result.Filter {
			if len(cond.Value) == 0 {
				continue
			}
			for _, v := range cond.Value {
				if exist, err = s.labelExist(&model.LabelCore{
					CatType:  model.PgcLabel,
					Category: cat,
					Param:    cond.ID,
					Value:    v.ID,
				}); err != nil {
					return
				}
				if exist {
					continue
				}
				label := model.LabelDB{}
				label.FromPgcCond(v, cond, cat)
				if err = s.DB.Create(&label).Error; err != nil {
					log.Error("pgcLabels Param %s, Cond %v Create Err %v", cond.ID, v, err)
					return
				}
				time.Sleep(20 * time.Millisecond)
			}
		}
	}
	return
}

// AddUgcTm adds time label for ugc
func (s *Service) AddUgcTm(tm *model.UgcTime) (err error) {
	var (
		exist bool
		timeV = tm.TimeV()
	)
	if exist, err = s.labelExist(&model.LabelCore{
		Category: tm.Category,
		Param:    model.ParamUgctime,
		Name:     tm.Name,
	}); err != nil {
		return
	}
	if exist {
		err = ecode.TvLabelExist
		return
	}
	label := model.LabelDB{}
	label.FromUgcTime(tm, s.c.Cfg.RefLabel.UgcTime)
	if err = s.DB.Create(&label).Error; err != nil {
		log.Error("ugcTimeLabel Time %s, Create Err %v", timeV, err)
	}
	return
}

// EditUgcTm edits a time label by name
func (s *Service) EditUgcTm(tm *model.EditUgcTime) (err error) {
	var exist bool
	if exist, err = s.labelExist(&model.LabelCore{
		ID: tm.ID,
	}); err != nil {
		return
	}
	if !exist {
		err = ecode.NothingFound
		return
	}
	if err = s.DB.Model(&model.LabelDB{}).Where("id = ?", tm.ID).Update(map[string]string{
		"name":  tm.Name,
		"value": tm.TimeV(),
	}).Error; err != nil {
		log.Error("ugcTimeLabel LabelID %d, Update Err %v", tm.ID, err)
	}
	return
}

// ActLabels act on labels
func (s *Service) ActLabels(ids []int64, act int) (err error) {
	if err = s.labelsExist(ids); err != nil {
		return
	}
	if err = s.DB.Model(&model.LabelDB{}).Where(fmt.Sprintf("id IN (%s)", xstr.JoinInts(ids))).Update(map[string]int{
		"valid": act,
	}).Error; err != nil {
		log.Error("ActLabels LabelID %s, Update Err %v", ids, err)
	}
	return
}

// DelLabels deletes labels
func (s *Service) DelLabels(ids []int64) (err error) {
	var exist bool
	for _, v := range ids {
		if exist, err = s.labelExist(&model.LabelCore{
			ID:      v,
			CatType: model.UgcLabel,
			Param:   model.ParamUgctime,
		}); err != nil {
			return
		}
		if !exist {
			log.Warn("DelLabels IDs %v, ID not exist %d", ids, v)
			return ecode.RequestErr
		}
	}
	if err = s.DB.Exec(fmt.Sprintf("UPDATE tv_label SET deleted = 1 WHERE id IN (%s)", xstr.JoinInts(ids))).Error; err != nil {
		log.Error("DelLabels LabelID %s, Update Err %v", ids, err)
	}
	return
}

// DynamicLabels picks the defined pgc label types
func (s *Service) labelTypes() (tps []*model.TpLabel, err error) {
	var rows *sql.Rows
	//  select category, param , param_name from tv_label where deleted = 0 and cat_type = 1 group by category,param
	if rows, err = s.DB.Model(&model.LabelDB{}).Where("deleted = 0").
		Where("cat_type = ?", model.PgcLabel).Select("category, param, param_name").Group("category,param").Rows(); err != nil {
		log.Error("labelTypes rows Err %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var cont = &model.TpLabel{}
		if err = rows.Scan(&cont.Category, &cont.Param, &cont.ParamName); err != nil {
			log.Error("rows.Scan error(%v)", err)
			return
		}
		tps = append(tps, cont)
	}
	if err = rows.Err(); err != nil {
		log.Error("labelTypes rows Err %v", err)
	}
	return
}

func (s *Service) loadLabel() (err error) {
	var (
		newTps = make(map[int][]*model.TpLabel)
		labels []*model.TpLabel
	)
	if labels, err = s.labelTypes(); err != nil {
		log.Error("labelTypes err %v", err)
		time.Sleep(time.Duration(10 * time.Second))
		return
	}
	for _, v := range labels {
		if lbs, ok := newTps[v.Category]; ok {
			newTps[v.Category] = append(lbs, v)
		} else {
			newTps[v.Category] = append([]*model.TpLabel{}, v)
		}
	}
	if len(newTps) > 0 {
		s.labelTps = newTps
	}
	return
}

// LabelTp returns the category's label types
func (s *Service) LabelTp(category int) (lbs []*model.TpLabel, err error) {
	var ok bool
	if lbs, ok = s.labelTps[category]; !ok {
		err = ecode.RequestErr
	}
	return
}

// PickLabels picks ugc labels
func (s *Service) PickLabels(req *model.ReqLabel, catType int) (data []*model.LabelList, err error) {
	var (
		db = s.DB.Model(&model.LabelDB{}).
			Where("param = ?", req.Param).
			Where("category = ?", req.Category).
			Where("deleted = 0").
			Where("cat_type = ?", catType)
		labels []*model.LabelDB
	)
	if req.Title != "" {
		db = db.Where("name LIKE ?", "%"+req.Title+"%")
	}
	if req.ID != 0 {
		db = db.Where("id = ?", req.ID)
	}
	if err = db.Order("position ASC").Find(&labels).Error; err != nil {
		log.Error("PickLabels Req %v, Err %v", req, err)
	}
	for _, v := range labels {
		data = append(data, v.ToList())
	}
	return
}

// EditLabel edits a pgc label
func (s *Service) EditLabel(id int64, name string) (err error) {
	var exist bool
	if exist, err = s.labelExist(&model.LabelCore{
		ID: id,
	}); err != nil {
		return
	}
	if !exist {
		err = ecode.NothingFound
		return
	}
	if err = s.DB.Model(&model.LabelDB{}).Where("id = ?", id).Update(map[string]string{
		"name": name,
	}).Error; err != nil {
		log.Error("EditLabel LabelID %d, Update Err %v", id, err)
	}
	return
}

// PubLabel publish label's order
func (s *Service) PubLabel(ids []int64) (err error) {
	if len(ids) == 0 {
		return
	}
	var (
		labels   []*model.LabelDB
		position = 1
		tx       = s.DB.Begin()
		labelMap = make(map[int64]*model.LabelDB, len(ids))
	)
	if err = s.DB.Where(fmt.Sprintf("id IN (%s)", xstr.JoinInts(ids))).Find(&labels).Error; err != nil {
		log.Error("PubLabel Ids %v, Err %v", ids, err)
		return
	}
	if len(labels) == 0 {
		err = ecode.NothingFound
		return
	}
	for _, v := range labels {
		labelMap[v.ID] = v
	}
	for _, id := range ids {
		lbl, ok := labelMap[id]
		if !ok {
			err = ecode.RequestErr
			log.Warn("PubLabel Id %d Not found", id)
			return
		}
		if !lbl.SameType(labels[0]) {
			log.Error("PubLabel Id %d FirstLabel ID %d, Not same type", lbl.ID, labels[0].ID)
			err = ecode.RequestErr
			tx.Rollback()
			return
		}
		if err = tx.Model(&model.LabelDB{}).Where("id = ?", lbl.ID).Update(map[string]int{"position": position}).Error; err != nil {
			log.Error("PubLabel ID %d, Err %v", lbl.ID, err)
			tx.Rollback()
			return
		}
		position = position + 1
	}
	tx.Commit()
	return
}
