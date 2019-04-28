package service

import (
	"context"
	"sort"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// CheckChannelCategory check channel category.
func (s *Service) CheckChannelCategory(name string) (string, error) {
	var rb []byte
	bytes := []byte(strings.Trim(name, " "))
	for _, b := range bytes {
		if b < 0x20 || b == 0x7f {
			return "", ecode.RequestErr
		}
		rb = append(rb, b)
	}
	if len(rb) == 0 {
		return "", ecode.RequestErr
	}
	return string(rb), nil
}

// ChannelCategories channel categories.
func (s *Service) ChannelCategories(c context.Context) (res []*model.ChannelCategory, err error) {
	categories, err := s.dao.ChannelCategories(c)
	if err != nil {
		return
	}
	res = make([]*model.ChannelCategory, 0, len(categories))
	for _, category := range categories {
		if category.State == model.StateNormal {
			category.INTShield = category.AttrVal(model.ChannelCategoryAttrINT)
			res = append(res, category)
		}
	}
	sort.Sort(model.ChannelCategorySort(res))
	return
}

// AddChannelCategory add channel category.
func (s *Service) AddChannelCategory(c context.Context, name string, shieldINT int32) (err error) {
	category, err := s.dao.ChannelCategoryByName(c, name)
	if err != nil {
		return
	}
	if category != nil {
		switch category.State {
		case model.StateDel:
			err = ecode.ChannelTypeDeled
		case model.StateNormal:
			err = ecode.ChannelTypeExist
		default:
			err = ecode.TagOperateFail
		}
		return
	}
	order, err := s.dao.CountChannelCategory(c)
	if err != nil {
		return
	}
	category = &model.ChannelCategory{
		Name:  name,
		Order: order,
		State: model.StateNormal,
	}
	category.AttrSet(model.ChannelCategoryAttrINT, shieldINT)
	_, err = s.dao.InsertChannelCategory(c, category)
	return
}

// DeleteChannelCategory delete channel category.
func (s *Service) DeleteChannelCategory(c context.Context, id int64, name string) (err error) {
	var (
		category *model.ChannelCategory
		channels []*model.Channel
	)
	if name != "" {
		if category, err = s.dao.ChannelCategoryByName(c, name); err != nil {
			return
		}
		if category == nil {
			return ecode.ChanTypeNotExist
		}
		id = category.ID
	}
	if channels, _, err = s.dao.ChannelsByType(c, id); err != nil {
		return
	}
	for _, channel := range channels {
		if channel.State > model.ChanStateStop || channel.State == model.ChanStateOffline {
			return ecode.ChannelCanNotDel
		}
	}
	_, err = s.dao.StateChannelCategory(c, id, model.StateDel)
	return
}

// SortChannelCategory sort channel category.
func (s *Service) SortChannelCategory(c context.Context, ids []int64) (err error) {
	var (
		sortMap        = make(map[int64]int32, len(ids))
		allCategories  []*model.ChannelCategory
		sortCategories = make([]*model.ChannelCategory, 0, len(ids))
	)
	if allCategories, err = s.dao.ChannelCategories(c); err != nil {
		return
	}
	for offset, id := range ids {
		sortMap[id] = int32(offset)
	}
	for _, category := range allCategories {
		k, ok := sortMap[category.ID]
		if !ok {
			continue
		}
		category.Order = k
		sortCategories = append(sortCategories, category)
	}
	if len(sortCategories) == 0 {
		return
	}
	_, err = s.dao.UpdateChannelCategories(c, sortCategories)
	return
}

// CategoryShieldINT .
func (s *Service) CategoryShieldINT(c context.Context, id int64, state int32, uname string) (err error) {
	var (
		category *model.ChannelCategory
		channels []*model.Channel
	)
	if category, err = s.dao.ChannelCategory(c, id); err != nil {
		return
	}
	if category == nil {
		return ecode.ChanTypeNotExist
	}
	if category.AttrVal(model.ChannelCategoryAttrINT) == state {
		return ecode.CategoryNoChange
	}
	category.AttrSet(model.ChannelCategoryAttrINT, state)
	if state == model.CategoryStateShieldINT {
		if channels, _, err = s.dao.ChannelsByType(c, id); err != nil {
			return
		}
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.CategoryShieldINT BeginTran(%d,%v) error(%v)", id, uname, err)
		return
	}
	if _, err = s.dao.TxUpChannelCategoryAttr(tx, category); err != nil {
		tx.Rollback()
		return
	}
	if state == model.CategoryStateShieldINT {
		for _, channel := range channels {
			if channel.AttrVal(model.ChannelAttrINT) == model.ChannelStateShieldINT {
				continue
			}
			channel.AttrSet(model.ChannelAttrINT, model.ChannelStateShieldINT)
			channel.Operator = uname
			if _, err = s.dao.TxUpChannelAttr(tx, channel); err != nil {
				tx.Rollback()
				return
			}
		}
	}
	return tx.Commit()
}
