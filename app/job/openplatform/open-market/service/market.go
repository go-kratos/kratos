package service

import (
	"context"

	"go-common/app/job/openplatform/open-market/model"
	"go-common/library/log"
)

//fetch market data ,from db to es
func (s *Service) marketProc() (err error) {
	var projectList []*model.Project
	if projectList, err = s.dao.FetchProject(context.TODO()); err != nil {
		log.Error("pull project list error (%v)", err)
		return
	}
	for _, project := range projectList {
		var orderData, wishData, favoriteData, commentData, pvData, uvData map[int32]int64
		if orderData, err = s.dao.OrderData(context.TODO(), project.ID, project.StartTime); err != nil {
			log.Error("fetch project_order_data [%d] error (%v)", project.ID, err)
		}
		if commentData, err = s.dao.CommentData(context.TODO(), project.ID, project.StartTime); err != nil {
			log.Error("fetch project_comment_data [%d] error (%v)", project.ID, err)
		}
		if wishData, err = s.dao.WishData(context.TODO(), project.ID, project.StartTime); err != nil {
			log.Error("fetch project_wish_data [%d] error (%v)", project.ID, err)
		}
		if favoriteData, err = s.dao.FavoriteData(context.TODO(), project.ID, project.StartTime); err != nil {
			log.Error("fetch project_favorite_data [%d] error (%v)", project.ID, err)
		}
		if pvData, uvData, err = s.dao.QueryPUVCount(context.TODO(), project.ID); err != nil {
			log.Error("fetch project_puv_data [%d] error (%v)", project.ID, err)
		}
		project.PV = pvData
		project.UV = uvData
		project.SaleInfo = orderData
		project.CommentInfo = commentData
		project.WishInfo = wishData
		project.FavoriteInfo = favoriteData
		if err = s.dao.SaveData(context.TODO(), project); err != nil {
			log.Error("save [%d]data to es  error (%v)", project.ID, err)
		}
	}
	return
}
