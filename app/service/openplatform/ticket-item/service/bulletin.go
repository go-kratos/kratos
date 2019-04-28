package service

import (
	"context"

	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/log"
)

// GetBulletins 获取公告信息 测试用
/**func (s *ItemService) GetBulletins(c context.Context, id *int64) (res []*item.BulletinInfo, err error) {
	res, err = s.dao.GetBulletins(c, *id)
	return
}**/

// BulletinInfo 添加新公告版本
func (s *ItemService) BulletinInfo(c context.Context, info *item.BulletinInfoRequest) (res *item.BulletinReply, err error) {
	var dbReturn bool
	var dbErr error
	if info.VerID == 0 {
		dbReturn, dbErr = s.dao.AddBulletin(c, info)
	} else {
		dbReturn, dbErr = s.dao.UpdateBulletin(c, info)
	}
	return &item.BulletinReply{Success: dbReturn}, dbErr
}

// BulletinCheck 审核公告
func (s *ItemService) BulletinCheck(c context.Context, info *item.BulletinCheckRequest) (res *item.BulletinReply, err error) {
	var dbReturn bool
	var dbErr error
	versionLogData := &model.VersionLog{
		VerID: info.VerID,
		Type:  1,
		Log:   info.Comment,
		Uname: info.Reviewer,
	}
	if info.OpType == 1 {
		dbReturn, dbErr = s.dao.PassBulletin(c, info.VerID)
		versionLogData.IsPass = 1
	} else {
		dbReturn, dbErr = s.dao.RejectVersion(c, info.VerID, model.VerTypeBulletin)
		versionLogData.IsPass = 0

	}
	// 添加审核记录
	if logErr := s.dao.AddVersionLog(c, versionLogData); logErr != nil {
		log.Error("新建审核记录失败")
		dbReturn = false
		dbErr = logErr
	}

	return &item.BulletinReply{Success: dbReturn}, dbErr
}

// BulletinState 添加新公告版本
func (s *ItemService) BulletinState(c context.Context, info *item.BulletinStateRequest) (res *item.BulletinReply, err error) {
	var dbReturn bool
	var dbErr error
	if info.OpType == 1 {
		// 上架
		dbReturn, dbErr = s.dao.PassBulletin(c, info.VerID)
	} else {
		// 下架
		var status int8
		if info.Source == 1 {
			status = -1 // 手动下架
		} else {
			status = -2 // 强制下架
		}
		dbReturn, dbErr = s.dao.UnpublishBulletin(c, info.VerID, status)

	}
	return &item.BulletinReply{Success: dbReturn}, dbErr
}
