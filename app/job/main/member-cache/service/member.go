package service

import (
	"context"
	"encoding/json"
	"strings"

	"go-common/app/job/main/member-cache/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

// consts
const (
	_MemberBaseInfo      = "user_base_"
	_MemberMoral         = "user_moral"
	_memberExp           = "user_exp_"
	_memberRealnameApply = "realname_apply"
	_memberRealnameInfo  = "realname_info"
)

// consts
const (
	expMulti = 100
	level1   = 1
	level2   = 200
	level3   = 1500
	level4   = 4500
	level5   = 10800
	level6   = 28800
)

func isExpAndLevelChange(mu *model.Binlog) (bool, bool) {
	if mu.Action == "insert" {
		return true, true
	}
	if len(mu.Old) <= 0 || len(mu.New) <= 0 {
		return false, false
	}
	old := &model.ExpMessage{}
	new := &model.ExpMessage{}
	if err := json.Unmarshal(mu.New, new); err != nil {
		return false, false
	}
	if err := json.Unmarshal(mu.Old, old); err != nil {
		return false, false
	}
	expChange := false
	levelChange := false
	if old.Exp != new.Exp {
		expChange = true
	}
	if level(old.Exp) != level(new.Exp) {
		levelChange = true
	}
	return expChange, levelChange
}

func level(exp int64) int8 {
	exp = exp / expMulti
	switch {
	case exp < level1:
		return 0
	case exp < level2:
		return 1
	case exp < level3:
		return 2
	case exp < level4:
		return 3
	case exp < level5:
		return 4
	case exp < level6:
		return 5
	default:
		return 6
	}
}

func resolveBaseAct(old json.RawMessage, new json.RawMessage) string {
	if old == nil || new == nil {
		return model.ActUpdateFace
	}
	ob := &model.MemberBase{}
	if err := json.Unmarshal(old, ob); err != nil {
		log.Error("Failed to parse data: %s", string(old))
		return model.ActUpdateFace
	}
	nb := &model.MemberBase{}
	if err := json.Unmarshal(new, nb); err != nil {
		log.Error("Failed to parse data: %s", string(new))
		return model.ActUpdateFace
	}
	if ob.Name != nb.Name {
		return model.ActUpdateUname
	}
	if ob.Face != nb.Face {
		return model.ActUpdateFace
	}
	return model.ActUpdateFace
}

func (s *Service) handleMemberBinLog(ctx context.Context, msg *databus.Message) error {
	defer func() {
		if err := msg.Commit(); err != nil {
			log.Error("Failed to commit message: %+v", BeautifyMessage(msg))
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

	switch {
	case strings.HasPrefix(mu.Table, _MemberBaseInfo):
		if err := s.dao.DelBaseInfoCache(ctx, mmid.Mid); err != nil {
			return err
		}
		s.NotifyPurgeCache(ctx, mmid.Mid, resolveBaseAct(mu.Old, mu.New))
	case mu.Table == _MemberMoral:
		if err := s.dao.DelMoralCache(ctx, mmid.Mid); err != nil {
			return err
		}
		s.NotifyPurgeCache(ctx, mmid.Mid, model.ActUpdateMoral)
	case mu.Table == _memberRealnameInfo || mu.Table == _memberRealnameApply:
		if err := s.dao.DeleteRealnameCache(ctx, mmid.Mid); err != nil {
			return err
		}
		s.NotifyPurgeCache(ctx, mmid.Mid, model.ActUpdateRealname)
	case strings.HasPrefix(mu.Table, _memberExp):
		mexp := &model.NewExp{}
		if err := json.Unmarshal(mu.New, mexp); err != nil {
			return errors.WithStack(err)
		}
		if err := s.dao.SetExpCache(ctx, mexp.Mid, mexp.Exp); err != nil {
			return err
		}
		expChange, levelChange := isExpAndLevelChange(mu)
		if expChange {
			s.NotifyPurgeCache(ctx, mmid.Mid, model.ActUpdateExp)
		}
		if levelChange {
			s.NotifyPurgeCache(ctx, mmid.Mid, model.ActUpdateLevel)
		}
	default:
		if mmid.Mid <= 0 {
			log.Info("Invalid message: %+v", BeautifyMessage(msg))
			return nil
		}
		s.deleteAllCache(ctx, mmid.Mid)
		s.NotifyPurgeCache(ctx, mmid.Mid, model.ActUpdateByAdmin)
	}

	return nil
}

func (s *Service) deleteAllCache(ctx context.Context, mid int64) error {
	if err := s.dao.DelBaseInfoCache(ctx, mid); err != nil {
		log.Error("Failed to delete cache: %+v", err)
	}
	if err := s.dao.DelMoralCache(ctx, mid); err != nil {
		log.Error("Failed to delete cache: %+v", err)
	}
	if err := s.dao.DeleteRealnameCache(ctx, mid); err != nil {
		log.Error("Failed to delete cache: %+v", err)
	}
	if err := s.dao.DelExpCache(ctx, mid); err != nil {
		log.Error("Failed to delete cache: %+v", err)
	}
	return nil
}

func (s *Service) memberBinLogproc(ctx context.Context) {
	for msg := range s.memberBinLog.Messages() {
		if err := s.handleMemberBinLog(ctx, msg); err != nil {
			log.Error("Failed to handle member binlog: %s: %+v", BeautifyMessage(msg), err)
			continue
		}
		log.Info("Succeed to handle member binlog: %s", BeautifyMessage(msg))
	}
}
