package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/member-cache/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

func (s *Service) handleBlockBinLog(ctx context.Context, msg *databus.Message) error {
	defer func() {
		if err := msg.Commit(); err != nil {
			log.Error("Failed to commit message: %s: %+v", BeautifyMessage(msg), err)
			return
		}
	}()

	mu := &model.Binlog{}
	if err := json.Unmarshal(msg.Value, mu); err != nil {
		return errors.WithStack(err)
	}

	mmid := &model.NeastMid{}
	bs := mu.New
	if len(bs) <= 0 {
		bs = mu.Old
	}
	if err := json.Unmarshal(bs, mmid); err != nil {
		return errors.WithStack(err)
	}

	defer s.NotifyPurgeCache(ctx, mmid.Mid, model.ActBlockUser)
	return s.dao.DeleteUserBlockCache(ctx, mmid.Mid)
}

func (s *Service) blockBinLogproc(ctx context.Context) {
	for msg := range s.blockBinLog.Messages() {
		if err := s.handleBlockBinLog(ctx, msg); err != nil {
			log.Error("Failed to handle block binlog: %s: %+v", BeautifyMessage(msg), err)
			continue
		}
		log.Info("Succeed to handle block binlog: %s", BeautifyMessage(msg))
	}
}
