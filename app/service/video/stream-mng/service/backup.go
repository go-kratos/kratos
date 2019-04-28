package service

import (
	"context"
	"go-common/app/service/video/stream-mng/model"
)

// CreateBackupStream 创建备用流
func (s *Service) CreateBackupStream(ctx context.Context, bs *model.BackupStream) (*model.BackupStream, error) {
	res, err := s.dao.CreateBackupStream(ctx, bs)
	if err == nil {

		// 更新redis
		s.dao.UpdateStreamStatusCache(ctx, &model.StreamStatus{
			RoomID:          bs.RoomID,
			StreamName:      bs.StreamName,
			DefaultUpStream: bs.DefaultVendor,
			DefaultChange:   true,
			Origin:          bs.OriginUpstream,
			OriginChange:    true,
			Forward:         bs.Streaming,
			ForwardChange:   true,
			Key:             bs.Key,
			Add:             true,
		})
	}

	return res, err
}
