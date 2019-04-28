package service

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"time"
)

// CreateOfficalStream 创建正式流
func (s *Service) CreateOfficalStream(c context.Context, streamName string, key string, rid int64) bool {
	succ := false
	infos := []*model.OfficialStream{}
	// 创建所有cdn的数据 ，为防止其中一条插入失败， 使用事务操作
	for k, v := range common.CdnMapSrc {
		var upRank int64

		if k == common.QNName {
			upRank = 1
		}

		ts := time.Now()

		stream := model.OfficialStream{
			RoomID:              rid,
			Name:                streamName,
			Key:                 key,
			Src:                 v,
			UpRank:              upRank,
			DownRank:            0,
			Status:              0,
			LastStatusUpdatedAt: ts,
			UpdateAt:            ts,
			CreateAt:            ts,
		}

		infos = append(infos, &stream)
	}

	err := s.dao.CreateOfficialStream(c, infos)

	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("stream create faild err= %v, room_id= %d", err, rid)))
	} else {
		succ = true
		go func(ctx context.Context, roomID int64) {
			s.syncMainStream(ctx, roomID, "")
		}(metadata.WithContext(c), rid)
	}

	return succ
}
