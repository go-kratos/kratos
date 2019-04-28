package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/reply-feed/model"
	"go-common/library/log"
)

// func (s *Service) eventproc() {
// 	defer s.waiter.Done()
// 	msgs := s.eventConsumer.Messages()
// 	ctx := context.Background()
// 	for {
// 		msg, ok := <-msgs
// 		if !ok {
// 			log.Warn("databus consumer channel has been closed.")
// 			return
// 		}
// 		if msg.Topic != s.c.Databus.Event.Topic {
// 			log.Warn("wrong topic actual (%s) expect (%s)", msg.Topic, s.c.Databus.Stats.Topic)
// 			continue
// 		}
// 		value := &model.EventMsg{}
// 		if err := json.Unmarshal(msg.Value, value); err != nil {
// 			log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
// 			continue
// 		}
// 		switch value.Action {
// 		case model.DatabusActionReIdx:
// 			s.setReplySetBatch(ctx, value.Oid, value.Tp)
// 			s.upsertZSet(ctx, value.Oid, value.Tp)
// 		default:
// 			continue
// 		}
// 		msg.Commit()
// 		log.Info("consumer topic:%s, partitionId:%d, offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
// 	}
// }

func (s *Service) statsproc() {
	defer s.waiter.Done()
	msgs := s.statsConsumer.Messages()
	for {
		msg, ok := <-msgs
		if !ok {
			log.Warn("databus consumer channel has been closed.")
			return
		}
		if msg.Topic != s.c.Databus.Stats.Topic {
			log.Warn("wrong topic actual (%s) expect (%s)", msg.Topic, s.c.Databus.Stats.Topic)
			continue
		}
		value := &model.StatsMsg{}
		if err := json.Unmarshal(msg.Value, value); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", msg.Value, err)
			continue
		}
		// 脏数据
		if value.Reply == nil || value.Subject == nil || (value.Action == model.DatabusActionReport && value.Report == nil) {
			log.Error("illegal message (%v)", value)
			continue
		}
		ctx := context.Background()
		// 针对评论列表的流程
		s.replyListFlow(ctx, value)
		// 针对统计数据的流程
		s.statisticsFlow(ctx, value)
		msg.Commit()
		log.Info("consumer topic:%s, partitionId:%d, offset:%d, Key:%s, Value:%s", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
	}
}

func (s *Service) statisticsFlow(ctx context.Context, value *model.StatsMsg) {
	var (
		reply      = value.Reply
		oid        = value.Subject.Oid
		tp         = value.Subject.Type
		rpID       = reply.RpID
		isHotReply bool
		name       string
		err        error
	)
	s.statisticsLock.RLock()
	name = s.statisticsStats[value.Sharding()].Name
	s.statisticsLock.RUnlock()
	if value.HotCondition() {
		if !reply.IsRoot() {
			rpID = reply.Root
		}
		if name == model.DefaultAlgorithm {
			if isHotReply, err = s.dao.IsOriginHot(ctx, oid, rpID, tp); err != nil {
				return
			}
		} else {
			if isHotReply, err = s.isHot(ctx, name, oid, rpID, tp); err != nil {
				return
			}
		}
	}
	s.addUV(ctx, value, isHotReply)
	switch value.Action {
	case model.DatabusActionLike:
		if isHotReply {
			s.statisticsStats[value.Sharding()].HotLike++
		}
		s.statisticsStats[value.Sharding()].TotalLike++
	case model.DatabusActionHate:
		if isHotReply {
			s.statisticsStats[value.Sharding()].HotHate++
		}
		s.statisticsStats[value.Sharding()].TotalHate++
	case model.DatabusActionCancelLike:
		if isHotReply && s.statisticsStats[value.Sharding()].HotLike > 0 {
			s.statisticsStats[value.Sharding()].HotLike--
		}
		if s.statisticsStats[value.Sharding()].TotalLike > 0 {
			s.statisticsStats[value.Sharding()].TotalLike--
		}
	case model.DatabusActionCancelHate:
		if isHotReply && s.statisticsStats[value.Sharding()].HotHate > 0 {
			s.statisticsStats[value.Sharding()].HotHate--
		}
		if s.statisticsStats[value.Sharding()].TotalHate > 0 {
			s.statisticsStats[value.Sharding()].TotalHate--
		}
	case model.DatabusActionReport:
		if isHotReply {
			s.statisticsStats[value.Sharding()].HotReport++
		}
		s.statisticsStats[value.Sharding()].TotalReport++
	case model.DatabusActionReply:
		if reply.IsRoot() {
			s.statisticsStats[value.Sharding()].TotalRootReply++
		} else {
			if isHotReply {
				s.statisticsStats[value.Sharding()].HotChildReply++
			}
			s.statisticsStats[value.Sharding()].TotalChildReply++
		}
	}
}

func (s *Service) replyListFlow(ctx context.Context, value *model.StatsMsg) {
	var (
		subject     = value.Subject
		reply       = value.Reply
		oid         = subject.Oid
		tp          = subject.Type
		stat        *model.ReplyStat
		reportCount int
		err         error
	)
	if value.Report == nil {
		reportCount = 0
	} else {
		reportCount = value.Report.Count
	}
	// if root reply get stat, else get root reply stat
	if reply.IsRoot() {
		stat = &model.ReplyStat{
			RpID:        reply.RpID,
			Reply:       reply.RCount,
			Like:        reply.Like,
			Hate:        reply.Hate,
			Report:      reportCount,
			SubjectTime: subject.CTime,
			ReplyTime:   reply.CTime,
		}
	} else {
		if stat, err = s.GetStatByID(ctx, oid, tp, reply.Root); err != nil || stat == nil {
			return
		}
	}
	if reply.IsRoot() {
		switch value.Action {
		case model.DatabusActionTop, model.DatabusActionDel, model.DatabusActionRptDel:
			s.remReply(ctx, oid, tp, stat.RpID)
		case model.DatabusActionUnTop, model.DatabusActionRecover, model.DatabusActionReply:
			s.addReplySet(ctx, oid, tp, stat.RpID)
		case model.DatabusActionLike, model.DatabusActionCancelLike, model.DatabusActionCancelHate, model.DatabusActionHate, model.DatabusActionReport:
			s.updateStat(ctx, stat.RpID, stat)
		default:
			return
		}
	} else {
		switch value.Action {
		case model.DatabusActionReply, model.DatabusActionRecover:
			stat.Reply++
		case model.DatabusActionDel, model.DatabusActionRptDel:
			if stat.Reply > 0 {
				stat.Reply--
			}
		default:
			return
		}
		s.updateStat(ctx, stat.RpID, stat)
	}
	s.upsertZSet(ctx, oid, tp)
}
