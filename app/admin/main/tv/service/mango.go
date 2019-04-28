package service

import (
	"fmt"

	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

// recomExist checks whether the recom exist or not
func (s *Service) recomExist(id int64) (exist bool) {
	var recom model.MangoRecom
	if err := s.DB.Where("id = ?", id).Where("deleted = 0").Find(&recom).Error; err != nil {
		log.Error("[recomExist] ID %d, Err %v", id, err)
		return
	}
	if recom.ID != 0 {
		return true
	}
	return
}

func (s *Service) resExist(rid int64, rtype int) (exist bool) {
	var recom model.MangoRecom
	if err := s.DB.Where("rid = ?", rid).Where("rtype = ?", rtype).Where("deleted = 0").Find(&recom).Error; err != nil {
		log.Error("[resExist] rid %d, Err %v", rid, err)
		return
	}
	if recom.ID != 0 {
		return true
	}
	return
}

// MangoList picks the mango recom list data
func (s *Service) MangoList(c *bm.Context) (data *model.MangoListResp, err error) {
	var (
		recoms   []*model.MangoRecom
		msg      = s.c.Cfg.MangoErr
		invalids []string
		mre      *model.MRecomMC
	)
	data = &model.MangoListResp{List: make([]*model.MangoRecom, 0)}
	if err = s.DB.Where("deleted = ?", 0).Order("`rorder` ASC").Find(&recoms).Error; err != nil {
		log.Error("[MangoList] DB query fail(%v)", err)
		return
	}
	for _, v := range recoms { // check whether the archive or the season is still valid, otherwise we delete it and remind the user
		if v.Rtype == _TypePGC {
			if ok, _ := s.snValid(v.RID); ok {
				data.List = append(data.List, v)
			} else {
				invalids = append(invalids, fmt.Sprintf("p%d", v.RID))
				s.MangoDel(c, v.ID)
			}
		} else if v.Rtype == _TypeUGC {
			if ok, _ := s.arcValid(v.RID); ok {
				data.List = append(data.List, v)
			} else {
				invalids = append(invalids, fmt.Sprintf("u%d", v.RID))
				s.MangoDel(c, v.ID)
			}
		} else {
			log.Error("MangoList ID %d, Rid %d, Type %d, TypeError", v.ID, v.RID, v.Rtype)
			invalids = append(invalids, fmt.Sprintf("%d", v.RID))
		}
	}
	if len(invalids) > 0 {
		data.Message = msg + joinStr(invalids)
	}
	if mre, err = s.dao.GetMRecom(c); err != nil {
		log.Error("MangoList GetMRecom Err %v", err)
		err = nil
		return
	}
	data.Pubtime = mre.Pubtime.Time().Format("2006-01-02 15:04:05")
	return
}

// joinStr joins strings
func joinStr(src []string) (res string) {
	for k, v := range src {
		if k == len(src)-1 {
			res = res + v
		} else {
			res = res + v + ","
		}
	}
	return
}

// MangoAdd adds the mango recom data
func (s *Service) MangoAdd(c *bm.Context, rtype int, rids []int64) (data *model.MangoAdd, err error) {
	data = &model.MangoAdd{
		Succ:     make([]int64, 0),
		Invalids: make([]int64, 0),
	}
	var (
		succRecoms []*model.MangoRecom
		newRids    []int64
	)
	for _, v := range rids { // 检查是否存在
		if ok := s.resExist(v, rtype); ok {
			data.Invalids = append(data.Invalids, v)
			continue
		}
		newRids = append(newRids, v)
	}
	if rtype == _TypePGC { // 检查对应的pgc和ugc是否存在并有效
		for _, v := range newRids {
			if ok, sn := s.snValid(v); !ok {
				data.Invalids = append(data.Invalids, v)
			} else {
				succRecoms = append(succRecoms, sn.ToMango())
			}
		}
	} else if rtype == _TypeUGC {
		for _, v := range newRids {
			if ok, arc := s.arcValid(v); !ok {
				data.Invalids = append(data.Invalids, v)
			} else {
				var pid int32
				if _, pid, err = s.arcPName(arc.TypeID); err != nil || pid == 0 {
					log.Warn("MangoAdd Aid %d, TypeID %d, Err %v", v, arc.TypeID, err)
					data.Invalids = append(data.Invalids, v)
					continue
				}
				succRecoms = append(succRecoms, arc.ToMango(int(pid)))
			}
		}
	} else {
		err = ecode.TvDangbeiWrongType
		return
	}
	if len(succRecoms) > 0 { // 选取最大顺序，在之后递增
		tx := s.DB.Begin()
		maxOrder := s.dao.MaxOrder(c)
		for _, v := range succRecoms {
			maxOrder = maxOrder + 1
			v.Rorder = maxOrder
			if err = tx.Create(v).Error; err != nil {
				log.Error("MangoAdd Create Rid %d, Recom %v, Err %v", v.RID, v, err)
				tx.Rollback()
				return
			}
		}
		tx.Commit() // add succ
		for _, v := range succRecoms {
			data.Succ = append(data.Succ, v.RID)
		}
	}
	return
}

// MangoDel deletes the mango resource
func (s *Service) MangoDel(c *bm.Context, id int64) (err error) {
	if !s.recomExist(id) {
		return ecode.NothingFound
	}
	return s.dao.DelMRecom(c, id)
}

// MangoEdit edits the mango resource
func (s *Service) MangoEdit(c *bm.Context, req *model.ReqMangoEdit) (err error) {
	if !s.recomExist(req.ID) {
		return ecode.NothingFound
	}
	if err = s.DB.Model(&model.MangoRecom{}).Where("id = ?", req.ID).Update(map[string]interface{}{
		"title":     req.Title,
		"cover":     req.Cover,
		"content":   req.Content,
		"staff":     req.Staff,
		"jid":       req.JID,
		"playcount": req.Playcount,
	}).Error; err != nil {
		log.Error("MangoDel ID %d, Mango %v, Err %v", req.ID, req, err)
	}
	return
}

// MangoPub publish the latest order of ids
func (s *Service) MangoPub(c *bm.Context, ids []int64) (err error) {
	var order = 0
	for _, v := range ids {
		if !s.recomExist(v) {
			return errors.Wrap(ecode.Int(404), fmt.Sprintf("ID: %d", v))
		}
	}
	tx := s.DB.Begin()
	for _, v := range ids {
		order = order + 1
		if err = tx.Model(&model.MangoRecom{}).Where("id = ?", v).Update(map[string]int{"rorder": order}).Error; err != nil {
			log.Error("MangoPub ID %d, Err %v", v, err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	err = s.dao.MangoRecom(c, ids)
	return
}
