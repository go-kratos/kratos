package service

import (
	"context"
	"strconv"

	"go-common/app/job/openplatform/open-sug/model"
	"go-common/library/log"
)

func (s *Service) buildData() (err error) {
	var (
		itemList   []*model.Item
		scoreSlice []model.Score
	)
	if itemList, err = s.dao.FetchItem(context.TODO()); err != nil {
		log.Error("pull project list error (%v)", err)
		return
	}
	for _, item := range itemList {
		if scoreSlice, err = s.dao.SeasonData(context.TODO(), item); err != nil {
			log.Error("es item match season fail")
			continue
		}
		if len(scoreSlice) > 0 {
			if _, err = s.dao.SetItem(context.TODO(), item); err != nil {
				log.Error("set item redis error(%v)", err)
				continue
			}
			for _, score := range scoreSlice {
				score.Score = score.Score * s.BuildScore(item, score.SeasonID)
				if rowNum, _ := s.dao.InsertMatch(context.TODO(), item, score); rowNum > 0 {
					if _, err = s.dao.SetSug(context.TODO(), score.SeasonID, item.ID, score.Score); err != nil {
						log.Error("redis dao.SetSug error(%v)", err)
						continue
					}
				}
			}
		}
	}
	return
}

// BuildScore ...
func (s *Service) BuildScore(item *model.Item, seasonID string) (normalizationScore float64) {
	normalizationScore = 0.5*(s.Normalization(item.SalesCount, s.dao.ItemSalesMax[seasonID], s.dao.ItemSalesMin[seasonID])) + 0.2*(s.Normalization(item.CommentCount, s.dao.ItemCommentMax[seasonID], s.dao.ItemCommentMin[seasonID])) + 0.3*(s.Normalization(item.WishCount, s.dao.ItemWishMax[seasonID], s.dao.ItemCommentMin[seasonID]))
	return
}

// Normalization ...
func (s *Service) Normalization(x int, max int, min int) (normalizationScore float64) {
	if x == 0 {
		return 0
	}
	if float64(max-min) == 0 {
		return 1
	}
	return float64(x-min) / float64(max-min)
}

// Filter ...
func (s *Service) Filter() {
	c := context.TODO()
	if bindItems, err := s.dao.GetBind(c); err == nil {
		for _, buildItem := range bindItems {
			s.dao.GetItem(c, buildItem)
			log.Info("add filter data items_id(%d) season_id(%d)", buildItem.Item.ID, buildItem.SeasonID)
			if _, err = s.dao.SetItem(c, buildItem.Item); err != nil {
				log.Error("set item redis error(%v)", err)
				continue
			}
			log.Info("add sug res items_id(%d) season_id(%d)", buildItem.Item.ID, buildItem.SeasonID)
			if _, err = s.dao.SetSug(c, strconv.FormatInt(buildItem.SeasonID, 10), buildItem.Item.ID, float64(buildItem.Score)); err != nil {
				log.Error("redis dao.SetSug error(%v)", err)
				continue
			}
			s.dao.UpdatePic(c, buildItem)
		}
	}
}
