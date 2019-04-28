package service

import (
	"context"
	"strings"
	"unicode"

	"go-common/app/admin/main/filter/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
)

// AreaGroupList 获得当前所有area分组信息
func (s *Service) AreaGroupList(ctx context.Context, ps, pn int) (total int, list []*model.AreaGroup, err error) {
	if total, err = s.dao.AreaGroupTotal(ctx); err != nil {
		return
	}
	if list, err = s.dao.AreaGroupList(ctx, pn, ps); err != nil {
		return
	}
	return
}

func areaGroupNameCheck(name string) bool {
	if strings.TrimSpace(name) == "" {
		return false
	}
	if len(strings.Split(name, "")) > 10 {
		return false
	}
	return true
}

// AddAreaGroup 增加一个area分组
func (s *Service) AddAreaGroup(ctx context.Context, groupName string, adid int, adName string) (err error) {
	if !areaGroupNameCheck(groupName) {
		err = ecode.FilterInvalidAreaGroupName
		return
	}
	var (
		areaGroup *model.AreaGroup
	)
	if areaGroup, err = s.dao.AreaGroupByName(ctx, groupName); err != nil {
		return
	}
	if areaGroup != nil {
		err = ecode.FilterDuplicateAreaGroup
		return
	}
	areaGroup = &model.AreaGroup{
		Name: groupName,
	}
	var (
		groupID int64
		log     = &model.AreaGroupLog{
			AdID:   adid,
			AdName: adName,
			State:  model.LogStateAdd,
		}
		tx *xsql.Tx
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	if groupID, err = s.dao.TxInsertAreaGroup(ctx, tx, areaGroup); err != nil {
		tx.Rollback()
		return
	}
	if err = s.dao.TxInsertAreaGroupLog(ctx, tx, groupID, log); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// AreaList 获得当前所有area信息
func (s *Service) AreaList(ctx context.Context, groupID int, pn, ps int) (total int, list []*model.Area, err error) {
	if total, err = s.dao.AreaTotal(ctx, groupID); err != nil {
		return
	}
	if list, err = s.dao.AreaList(ctx, groupID, pn, ps); err != nil {
		return
	}
	return
}

func areaNameCheck(name string, showName string) error {
	if strings.TrimSpace(name) == "" {
		return ecode.FilterInvalidAreaName
	}
	if strings.TrimSpace(showName) == "" {
		return ecode.FilterInvalidAreaShowName
	}
	for _, s := range name {
		if !unicode.Is(unicode.Scripts["Latin"], s) && string(s) != "_" {
			return ecode.FilterInvalidAreaName
		}
	}
	if len(strings.Split(showName, "")) >= 10 {
		return ecode.FilterInvalidAreaShowName
	}
	return nil
}

// AddArea 增加area
func (s *Service) AddArea(ctx context.Context, groupID int, areaName string, areaShowName string, commonFlag bool, adID int, adName string) (err error) {
	if err = areaNameCheck(areaName, areaShowName); err != nil {
		return
	}
	var (
		area      *model.Area
		areaGroup *model.AreaGroup
	)
	if areaGroup, err = s.dao.AreaGroup(ctx, groupID); err != nil {
		return
	}
	if areaGroup == nil {
		err = ecode.FilterAreaGroupNotFound
		return
	}
	if area, err = s.dao.AreaByName(ctx, areaName); err != nil {
		return
	}
	if area != nil {
		err = ecode.FilterDuplicateArea
		return
	}
	area = &model.Area{
		GroupID:    groupID,
		Name:       areaName,
		ShowName:   areaShowName,
		CommonFlag: commonFlag,
	}
	var (
		tx     *xsql.Tx
		areaID int64
		log    = &model.AreaLog{
			AdID:    adID,
			AdName:  adName,
			Comment: "",
			State:   model.LogStateAdd,
		}
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	if areaID, err = s.dao.TxInsertArea(ctx, tx, area); err != nil {
		tx.Rollback()
		return
	}
	if err = s.dao.TxInsertAreaLog(ctx, tx, areaID, log); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// EditArea 编辑area
func (s *Service) EditArea(ctx context.Context, areaID int, commonFlag bool, adID int, adName string, adComment string) (err error) {
	var (
		area *model.Area
	)
	if area, err = s.dao.Area(ctx, areaID); err != nil {
		return
	}
	if area == nil {
		err = ecode.FilterInvalidArea
		return
	}
	area = &model.Area{
		ID:         areaID,
		GroupID:    area.GroupID,
		Name:       area.Name,
		ShowName:   area.ShowName,
		CommonFlag: commonFlag,
	}
	var (
		tx  *xsql.Tx
		log = &model.AreaLog{
			AdID:    adID,
			AdName:  adName,
			Comment: adComment,
			State:   model.LogStateEdit,
		}
	)
	if tx, err = s.dao.BeginTran(ctx); err != nil {
		return
	}
	if err = s.dao.TxUpdateArea(ctx, tx, area); err != nil {
		tx.Rollback()
		return
	}
	if err = s.dao.TxInsertAreaLog(ctx, tx, int64(area.ID), log); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// AreaLog 获得area的所有日志信息
func (s *Service) AreaLog(ctx context.Context, areaID int) (list []*model.AreaLog, err error) {
	if list, err = s.dao.AreaLog(ctx, areaID); err != nil {
		return
	}
	return
}
