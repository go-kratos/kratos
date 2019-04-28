package service

import (
	"context"
	"errors"
	"fmt"
	"go-common/app/service/video/stream-mng/model"
	"go-common/library/log"
)

// 这是一个过渡接口，用于搬运现有数据到新表中
func (s *Service) syncMainStream(ctx context.Context, roomID int64, streamName string) error {
	if roomID <= 0 && streamName == "" {
		return errors.New("invalid params")
	}

	var err error
	exists, err := s.dao.GetMainStreamFromDB(ctx, roomID, streamName)
	if err != nil && err.Error() != "sql: no rows in result set" {
		log.Errorv(ctx, log.KV("log", fmt.Sprintf("sync_stream_data_error = %v", err)))
		return err
	}
	if exists != nil && (exists.RoomID == roomID || exists.StreamName == streamName) {
		return nil
	}

	var full *model.StreamFullInfo
	if roomID > 0 && streamName == "" {
		full, err = s.GetStreamInfo(ctx, roomID, "")
	} else if roomID <= 0 && streamName != "" {
		full, err = s.GetStreamInfo(ctx, 0, streamName)
	}

	if err != nil {
		return err
	}
	if full == nil {
		return errors.New("unknow response")
	}

	for _, ss := range full.List {
		if ss.Type == 1 {
			ms := &model.MainStream{
				RoomID:        full.RoomID,
				StreamName:    ss.StreamName,
				Key:           ss.Key,
				DefaultVendor: ss.DefaultUpStream,
				Status:        1,
			}

			if ms.DefaultVendor == 0 {
				ms.DefaultVendor = 1
			}

			_, err := s.dao.CreateNewStream(ctx, ms)
			if err != nil {
				log.Errorv(ctx, log.KV("log", fmt.Sprintf("sync_stream_data_error = %v", err)))
			}
			break
		}
	}

	return nil
}
