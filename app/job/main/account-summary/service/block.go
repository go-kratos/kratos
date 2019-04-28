package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/account-summary/model"
	"go-common/library/log"
)

func (s *Service) blockBinLogproc(ctx context.Context) {
	for msg := range s.BlockBinLog.Messages() {
		blog := &model.CanalBinLog{}
		if err := json.Unmarshal(msg.Value, blog); err != nil {
			log.Error("Failed to unmarshal canal bin log: %+v, value: %s: %+v", msg, string(msg.Value), err)
			msg.Commit()
			continue
		}

		log.Info("Handling message key: %s, value: %s", msg.Key, string(msg.Value))
		s.blockBinLogHandle(ctx, blog)
		msg.Commit()
	}
}

func (s *Service) blockBinLogHandle(ctx context.Context, blog *model.CanalBinLog) {
	if len(blog.New) == 0 {
		log.Error("Failed to sync to hbase with empty new field: %+v", blog)
		return
	}

	switch blog.Table {
	case "block_user":
		midl := &model.MidBinLog{}
		if err := json.Unmarshal(blog.New, midl); err != nil {
			log.Error("Failed to unmarsha new data: %s: %+v", string(blog.New), err)
			return
		}
		// FIXME: 一段时间后改用 syncBlock
		if err := s.SyncOne(ctx, midl.Mid); err != nil {
			log.Error("Failed to sync block with mid: %d: %+v", midl.Mid, err)
			return
		}
	default:
		log.Warn("Unable to hanlde binlog: %+v, old: %s, new: %s", blog, string(blog.Old), string(blog.New))
	}
}
