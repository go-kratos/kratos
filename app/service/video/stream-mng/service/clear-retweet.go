package service

import (
	"context"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// ClearStreamStatus 清理互推标记
func (s *Service) ClearStreamStatus(c context.Context, rid int64) error {
	// 主流状态更新 : 将up_rank=2的流up_rank 设置为0
	err := s.dao.UpdateOfficialUpRankStatus(c, rid, 2, 0)
	if err != nil {
		return err
	}
	// 查询options
	infos, err := s.dao.StreamFullInfo(c, rid, "")
	//没有查到结果
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("select_maskbyroomid_error = %v", err)))
		return err
	}
	var newmask int64
	var options int64
	if infos != nil && infos.List != nil {
		for _, v := range infos.List {
			if v.Type == 1 {
				newmask = v.Options &^ 4
				newmask = newmask &^ 8
			}
		}
	}

	// 清除缓存， 清除forward
	s.dao.UpdateStreamStatusCache(c, &model.StreamStatus{
		RoomID:        rid,
		Forward:       0,
		ForwardChange: true,
		Options:       newmask,
		OptionsChange: true,
	})

	// 同步数据
	go func(ctx context.Context, roomID int64, newoptions int64, options int64) {
		s.syncMainStream(ctx, roomID, "")
		s.dao.ClearMainStreaming(ctx, roomID, newmask, options)
	}(metadata.WithContext(c), rid, newmask, options)

	// 备用流的状态无需更新
	return nil
}
