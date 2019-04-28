package service

import (
	"context"
	"go-common/app/service/video/stream-mng/common"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// ChangeSrc 切换cdn
func (s *Service) ChangeSrc(c context.Context, rid int64, toOrigin int64, source string, operateName string, reason string) error {
	_, origin, err := s.dao.OriginUpStreamInfo(c, rid)

	if err != nil {
		return err
	}

	// 更新正式流数据库
	toSrc := common.BitwiseMapSrc[toOrigin]
	err = s.dao.UpdateOfficialStreamStatus(c, rid, toSrc)
	if err != nil {
		return err
	}

	s.dao.UpdateStreamStatusCache(c, &model.StreamStatus{
		RoomID:          rid,
		DefaultChange:   true,
		DefaultUpStream: toOrigin,
	})

	go func(ctx context.Context, rid int64, origin int64, toOrigin int64, reason, operateName, source string) {
		// 更新main-stream
		if err := s.dao.ChangeDefaultVendor(ctx, rid, toOrigin); err != nil {
			log.Infov(ctx, log.KV("change_main_stream_default_err", err.Error()))
		}

		// 更新redis被切记录
		s.dao.UpdateChangeSrcCache(ctx, rid, origin)

		// 记录日志
		s.RecordChangeLog(ctx, &model.StreamChangeLog{
			RoomID:      rid,
			FromOrigin:  origin,
			ToOrigin:    toOrigin,
			Reason:      reason,
			OperateName: operateName,
			Source:      source,
		})
	}(metadata.WithContext(c), rid, origin, toOrigin, reason, operateName, source)

	return nil
}
