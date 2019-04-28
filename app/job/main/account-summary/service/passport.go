package service

import (
	"context"
	"encoding/json"
	"strings"

	"go-common/app/job/main/account-summary/model"
	"go-common/library/log"
)

func (s *Service) passportBinLogproc(ctx context.Context) {
	for msg := range s.PassportBinLog.Messages() {
		blog := &model.CanalBinLog{}
		if err := json.Unmarshal(msg.Value, blog); err != nil {
			log.Error("Failed to unmarshal canal bin log: %+v, value: %s: %+v", msg, string(msg.Value), err)
			msg.Commit()
			continue
		}

		log.Info("Handling message key: %s, value: %s", msg.Key, string(msg.Value))
		s.passportBinLogHandle(ctx, blog)
		msg.Commit()
	}
}

func (s *Service) passportBinLogHandle(ctx context.Context, blog *model.CanalBinLog) {
	if len(blog.New) == 0 {
		log.Error("Failed to sync to hbase with empty new field: %+v", blog)
		return
	}

	switch {
	case strings.HasPrefix(blog.Table, "aso_account"):
		midl := &model.MidBinLog{}
		if err := json.Unmarshal(blog.New, midl); err != nil {
			log.Error("Failed to unmarsha new data: %s: %+v", string(blog.New), err)
			return
		}
		if err := s.syncPassportSummary(ctx, midl.Mid); err != nil {
			log.Error("Failed to sync passport summary with mid: %d: %+v", midl.Mid, err)
			return
		}
	default:
		log.Warn("Unable to hanlde binlog: %+v, old: %s, new: %s", blog, string(blog.Old), string(blog.New))
	}
}
