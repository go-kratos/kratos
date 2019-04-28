package service

import (
	"context"
	"encoding/json"

	"go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

func (s *Service) award10kCondition(ctx context.Context, mid int64, achieve model.AchieveFlag) error {
	audit, err := s.Audit(ctx, mid, metadata.String(ctx, metadata.RemoteIP))
	if err != nil {
		return err
	}
	log.Info("award 10k audit info by mid: %d and audit: %+v", mid, audit)

	if audit.Rank < 10000 {
		return ecode.RelAwardInsufficientRank
	}
	if !audit.BindTel {
		return ecode.RelAwardPhoneRequired
	}
	if audit.Blocked {
		return ecode.RelAwardIsBlocked
	}

	if s.dao.HasReachAchieve(ctx, mid, achieve) {
		return nil
	}

	stat, err := s.Stat(ctx, mid)
	if err != nil {
		return err
	}
	if stat.Follower < 10000 {
		return ecode.RelAwardInsufficientFollower
	}
	return nil
}

// AchieveGet is
func (s *Service) AchieveGet(ctx context.Context, arg *model.ArgAchieveGet) (*model.AchieveGetReply, error) {
	if arg.Award != "10k" {
		return nil, ecode.RequestErr
	}

	if err := s.award10kCondition(ctx, arg.Mid, model.FollowerAchieve10k); err != nil {
		return nil, err
	}

	achieve := &model.Achieve{
		Award: arg.Award,
		Mid:   arg.Mid,
	}
	js, err := json.Marshal(achieve)
	if err != nil {
		return nil, errors.Wrapf(err, "achieve: %s", achieve)
	}
	token, err := encrypt([]byte(s.c.Relation.AchieveKey), string(js))
	if err != nil {
		log.Error("Failed to encrypt achieve with key: %s, text: %s: %+v", s.c.Relation.AchieveKey, string(js), errors.WithStack(err))
		return nil, ecode.RelAwardIsBlocked
	}
	return &model.AchieveGetReply{AwardToken: token}, nil
}

// Achieve is
func (s *Service) Achieve(ctx context.Context, arg *model.ArgAchieve) (*model.Achieve, error) {
	js, err := decrypt([]byte(s.c.Relation.AchieveKey), arg.AwardToken)
	if err != nil {
		log.Error("Failed to decrypt achieve with key: %s, token: %s: %+v", s.c.Relation.AchieveKey, arg.AwardToken, errors.WithStack(err))
		return nil, ecode.RelAwardInfoFailed
	}

	achieve := new(model.Achieve)
	if err := json.Unmarshal([]byte(js), achieve); err != nil {
		log.Error("Failed to parse achieve data: %s: %+v", js, err)
		return nil, ecode.RelAwardInfoFailed
	}

	return achieve, nil
}
