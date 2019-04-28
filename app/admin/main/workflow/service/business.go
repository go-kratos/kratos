package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"go-common/app/admin/main/workflow/model"
	"go-common/app/admin/main/workflow/model/param"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_lenShortCut = 8
)

// 0: 有效 1: 无效 2: 流转 3: 众裁 4: 批量有效 5: 批量无效 6: 批量流转 7: 批量移交众裁
var btnName = map[uint]string{
	0: "有效",
	1: "无效",
	2: "流转",
	3: "众裁",
	4: "批量有效",
	5: "批量无效",
	6: "批量流转",
	7: "批量移交众裁",
}

var btnInitShortCut = map[uint]string{
	0: "A",
	1: "S",
	2: "D",
	3: "F",
	4: "J",
	5: "K",
	6: "H",
	7: "L",
}

// ListMeta will list business meta inforamtion from DAO
func (s *Service) ListMeta(c context.Context, itemType string) (metaList []*model.Meta, err error) {
	allMetas := s.dao.AllMetas(c)
	if err != nil {
		log.Error("Failed to fetch all metas from database: %v", err)
		return
	}

	metaList = make([]*model.Meta, 0, len(allMetas))
	for _, m := range allMetas {
		if itemType != "" && m.ItemType != itemType {
			continue
		}
		//sort rounds field asc by id
		sort.Sort(model.RoundSlice(m.Rounds))
		metaList = append(metaList, m)
	}
	//sort meta asc by business
	sort.Sort(model.MetaSlice(metaList))

	return
}

// ListBusAttr list business attr info
func (s *Service) ListBusAttr(ctx context.Context) (busAttr []*model.BusinessAttr, err error) {
	busAttr = make([]*model.BusinessAttr, 0)
	if err = s.dao.ORM.Table("workflow_business_attr").Find(&busAttr).Error; err != nil {
		return
	}
	return
}

// ListBusAttrV3 .
func (s *Service) ListBusAttrV3(ctx context.Context) (busAttr []*model.BusinessAttr, err error) {
	log.Info("start ListBusAttrV3")
	t := time.Now()
	busAttr = make([]*model.BusinessAttr, 0)
	if err = s.dao.ORM.Table("workflow_business_attr").Find(&busAttr).Error; err != nil {
		log.Info("query table workflow_business_attr error(%v)", err)
		return
	}

	for _, attr := range busAttr {
		btnShortCut := strings.Split(attr.ButtonKey, ",")
		if len(btnShortCut) != _lenShortCut {
			log.Warn("button short cut length not 8, load initial value")
			btnShortCut = []string{"A", "S", "D", "F", "J", "K", "H", "L"}
		}
		for i := uint(0); i < _lenShortCut; i++ {
			state := false
			mask := uint8(1 << i)
			if attr.Button&mask > 0 {
				state = true
			}
			if btnShortCut[i] == "" {
				btnShortCut[i] = btnInitShortCut[i]
			}
			attr.Buttons = append(attr.Buttons, &model.Button{
				Index: int(i),
				Name:  btnName[i],
				State: state,
				Key:   btnShortCut[i],
			})
		}
	}
	log.Info("end ListBusAttrV3 time(%v)", time.Since(t).String())
	return
}

// AddOrUpdateBusAttr add or update business attr info
func (s *Service) AddOrUpdateBusAttr(ctx context.Context, abap *param.AddBusAttrParam) (err error) {
	busAttr := &model.BusinessAttr{
		ID: abap.ID,
	}
	attr := map[string]interface{}{
		"bid":           abap.Bid,
		"name":          abap.Name,
		"deal_type":     abap.DealType,
		"expire_time":   abap.ExpireTime,
		"assign_type":   abap.AssignType,
		"assign_max":    abap.AssignMax,
		"group_type":    abap.GroupType,
		"business_name": abap.BusinessName,
	}

	if err = s.dao.ORM.Table("workflow_business_attr").
		Where("id=?", abap.ID).
		Assign(attr).FirstOrCreate(busAttr).Error; err != nil {
		log.Error("Failed to create business_attr(%+v): %v", busAttr, err)
		return
	}
	s.loadBusAttrs()
	return
}

// SetSwitch .
func (s *Service) SetSwitch(ctx context.Context, bs *param.BusAttrButtonSwitch) (err error) {
	attr := new(model.BusinessAttr)
	if err = s.dao.ORM.Table("workflow_business_attr").Where("bid = ?", bs.Bid).Find(attr).Error; err != nil {
		log.Error("Failed to find business_attr where bid = %d : %v", bs.Bid, err)
		return
	}
	oldBut := attr.Button
	mask := uint8(^(1 << bs.Index))
	oldBut = oldBut & mask
	newBut := oldBut + (bs.Switch << bs.Index)
	if err = s.dao.ORM.Table("workflow_business_attr").Where("bid = ?", bs.Bid).Update("button", newBut).Error; err != nil {
		log.Error("Failed to update business_attr button field where bid = %d, button = %d : %v", bs.Bid, newBut, err)
	}
	return
}

// SetShortCut .
func (s *Service) SetShortCut(ctx context.Context, sc *param.BusAttrButtonShortCut) (err error) {
	attr := new(model.BusinessAttr)
	if err = s.dao.ORM.Table("workflow_business_attr").Where("bid = ?", sc.Bid).Find(attr).Error; err != nil {
		log.Error("Failed to find business_attr where bid = %d : %v", sc.Bid, err)
		return
	}

	oldShortCut := strings.Split(attr.ButtonKey, ",")
	if len(oldShortCut) != 8 {
		oldShortCut = make([]string, 8)
	}
	oldShortCut[sc.Index] = sc.ShortCut

	newShortCut := strings.Join(oldShortCut, ",")
	if err = s.dao.ORM.Table("workflow_business_attr").Where("bid = ?", sc.Bid).Update("button_key", newShortCut).Error; err != nil {
		log.Error("Failed to update business_attr button_key field where bid = %d, button_key = %d : %v", sc.Bid, newShortCut, err)
	}
	return
}

// ManagerTag .
func (s *Service) ManagerTag(ctx context.Context) (map[int8]map[int64]*model.TagMeta, error) {
	return s.tagListCache, nil
}

// UserBlockInfo .
// http://info.bilibili.co/pages/viewpage.action?pageId=5417571
// http://info.bilibili.co/pages/viewpage.action?pageId=7559616
func (s *Service) UserBlockInfo(ctx context.Context, bi *param.BlockInfo) (resp model.BlockInfoResp, err error) {
	var sum int64
	if sum, err = s.dao.BlockNum(ctx, bi.Mid); err != nil {
		log.Error("s.dao.BlockNum(%d) error(%v)", bi.Mid, err)
		return
	}

	if resp, err = s.dao.BlockInfo(ctx, bi.Mid); err != nil {
		log.Error("s.dao.BlockInfo(%d) error(%v)", bi.Mid, err)
		return
	}
	resp.Data.BlockedSum = sum
	return
}

// SourceList .
func (s *Service) SourceList(ctx context.Context, src *param.Source) (data map[string]interface{}, err error) {
	// check if has external uri
	if _, ok := s.callbackCache[src.Bid]; !ok {
		err = ecode.WkfBusinessCallbackConfigNotFound
		return
	}
	uri := ""
	if uri = s.callbackCache[src.Bid].SourceAPI; uri == "" {
		return
	}

	return s.dao.SourceInfo(ctx, uri)
}
