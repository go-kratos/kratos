package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/relation-cache/model"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_relationFidTable     = "user_relation_fid_"
	_relationMidTable     = "user_relation_mid_"
	_relationStatTable    = "user_relation_stat_"
	_relationTagUserTable = "user_relation_tag_user_"
)

func (s *Service) relationBinLogproc(ctx context.Context) {
	for msg := range s.relationBinLog.Messages() {
		if err := s.handleRelationBinLog(ctx, msg); err != nil {
			log.Error("Failed to handle relation binlog: %s: %+v", BeautifyMessage(msg), err)
			continue
		}
		log.Info("Succeed to handle relation binlog: %s", BeautifyMessage(msg))
	}
}

func (s *Service) handleRelationBinLog(ctx context.Context, msg *databus.Message) error {
	defer func() {
		if err := msg.Commit(); err != nil {
			log.Error("Failed to commit message: %+v", BeautifyMessage(msg))
			return
		}
	}()

	mu := &model.Message{}
	if err := json.Unmarshal(msg.Value, mu); err != nil {
		return errors.WithStack(err)
	}

	switch {
	case strings.HasPrefix(mu.Table, _relationStatTable):
		if err := s.stat(ctx, mu.Action, mu.New, mu.Old); err != nil {
			return err
		}
	case strings.HasPrefix(mu.Table, _relationMidTable):
		if err := s.relationMid(ctx, mu.Action, mu.New, mu.Old); err != nil {
			return err
		}
	case strings.HasPrefix(mu.Table, _relationFidTable):
		if err := s.relationFid(ctx, mu.Action, mu.New, mu.Old); err != nil {
			return err
		}
	case strings.HasPrefix(mu.Table, _relationTagUserTable):
		if err := s.tagUser(ctx, mu.New); err != nil {
			return err
		}
	}

	return nil
}

// stat
func (s *Service) stat(ctx context.Context, action string, nwMsg []byte, oldMsg []byte) (err error) {
	ms := &model.Stat{}
	if err = json.Unmarshal(nwMsg, ms); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", nwMsg, err)
		return
	}
	mo := &model.Stat{}
	if len(oldMsg) > 0 {
		if err = json.Unmarshal(oldMsg, mo); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", oldMsg, err)
			err = nil
		}
	}
	return s.dao.DelStatCache(ctx, ms.Mid)
}

// relationMid
func (s *Service) relationMid(ctx context.Context, action string, nwMsg []byte, oldMsg []byte) error {
	mr := &model.Relation{}
	if err := json.Unmarshal(nwMsg, mr); err != nil {
		return errors.WithStack(err)
	}

	f := &relation.Following{
		Mid:       mr.Fid,
		Attribute: mr.Attribute,
		MTime:     xtime.Time(time.Now().Unix()),
	}
	if err := s.upFollowingCache(ctx, mr.Mid, f); err != nil {
		return err
	}

	return s.dao.DelTagsCache(ctx, mr.Mid)
}

// relationFid
func (s *Service) relationFid(ctx context.Context, action string, nwMsg []byte, oldMsg []byte) error {
	var or *model.Relation
	mr := &model.Relation{}
	if err := json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%v) error(%v)", nwMsg, err)
		return err
	}

	if len(oldMsg) > 0 {
		or = new(model.Relation)
		if err := json.Unmarshal(oldMsg, or); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", oldMsg, err)
		}
	}

	return s.dao.DelFollowerCache(ctx, mr.Fid)
}

func (s *Service) tagUser(ctx context.Context, newMsg []byte) (err error) {
	var tags struct {
		Fid int64 `json:"fid"`
		Mid int64 `json:"mid"`
	}
	if err = json.Unmarshal(newMsg, &tags); err != nil {
		log.Error("json.Unmarshal err(%v)", err)
		return
	}
	return s.dao.DelTagsCache(ctx, tags.Mid)
}

func (s *Service) upFollowingCache(c context.Context, mid int64, f *relation.Following) (err error) {
	if f.Attribute == 0 {
		s.dao.DelFollowing(c, mid, f)
	} else {
		if err = s.dao.AddFollowingCache(c, mid, f); err != nil {
			return
		}
	}
	if err = s.dao.DelFollowingCache(c, mid); err != nil {
		return
	}
	return s.dao.DelTagCountCache(c, mid)
}
