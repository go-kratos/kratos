package service

import (
	"context"
	"time"

	"go-common/app/admin/openplatform/sug/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_setup  int8 = 1
	_remove int8 = 0
	_none   int8 = -1
)

// SourceSearch search season or items.
func (s *Service) SourceSearch(c context.Context, params *model.SourceSearch) (interface{}, error) {
	switch params.Type {
	case model.TypeSeason:
		return s.dao.SeasonList(c, params)
	case model.TypeItems:
		return s.dao.ItemList(c, params)
	}
	return nil, ecode.SugSearchTypeErr
}

// Search search season and items match list.
func (s *Service) Search(c context.Context, params *model.Search) (list []model.SugList, err error) {
	list, err = s.dao.SearchV2(c, params)
	if err != nil {
		err = ecode.SugEsSearchErr
	}
	if listLen := len(list); listLen > 0 {
		for i := 0; i < listLen-1; i++ {
			for j := 0; j < listLen-i-1; j++ {
				if list[j].Score < list[j+1].Score {
					list[j], list[j+1] = list[j+1], list[j]
				}
			}
		}
	}
	return
}

// MatchOperate match operate.
func (s *Service) MatchOperate(c context.Context, params *model.MatchOperate) (result interface{}, err error) {
	item, err := s.dao.GetItem(c, params.ItemsID)
	if err != nil {
		err = ecode.SugItemNone
		return
	}
	season, err := s.dao.GetSeason(c, params.SeasonID)
	if err != nil || season.ID == 0 {
		err = ecode.SugSeasonNone
		return
	}
	switch params.OpType {
	case _setup:
		if _, err = s.AddMatch(c, item, season); err != nil {
			err = ecode.SugOpErr
			return
		}
	case _remove:
		if err = s.DelMatch(c, item, season); err != nil {
			err = ecode.SugOpErr
			return
		}
	default:
		err = ecode.SugOpTypeErr
	}
	return
}

// DbOperate db operate.
func (s *Service) DbOperate(c context.Context, op int8, item model.Item, season model.Season, location string) (err error) {
	matchType, err := s.dao.GetMatchType(c, season.ID, item.ItemsID)
	if err != nil {
		return err
	}
	if matchType == _none {
		if _, err = s.dao.InsertMatch(c, season, item, op, time.Now().Unix(), location); err != nil {
			return err
		}
	}
	if _, err = s.dao.UpdateMatch(c, season, item, op, location); err != nil {
		return err
	}
	return
}

// AddMatch add match.
func (s *Service) AddMatch(c context.Context, item model.Item, season model.Season) (location string, err error) {
	if location, _ = s.dao.CreateItemPNG(item); location == "" {
		return
	}
	if err = s.DbOperate(c, _setup, item, season, location); err != nil {
		return
	}
	if _, err := s.dao.SetItem(c, &item, location); err != nil {
		log.Error("set item(%v) error(%v)", item, err)
	}
	if err := s.dao.SetSug(c, season.ID, item.ItemsID, time.Now().Unix()); err != nil {
		log.Error("zAdd season(%v) error(%v)", season, err)
	}
	return
}

// DelMatch del match.
func (s *Service) DelMatch(c context.Context, item model.Item, season model.Season) (err error) {
	if err = s.DbOperate(c, _remove, item, season, ""); err != nil {
		return
	}
	if err = s.dao.DelSug(c, season.ID, item.ItemsID); err != nil {
		log.Error("zRem item(%d) error(%v)", item, err)
	}
	return
}
