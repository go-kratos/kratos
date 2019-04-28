package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/creative/model"
	"go-common/library/log"
)

const (
	_archive         = "archive"
	_insert          = "insert"
	_update          = "update"
	fromArchiveUp    = 0
	fromArchiveNewUp = 1

	_relationMidTable  = "user_relation_mid_"
	_relationStatTable = "user_relation_stat_"
)

//pub  up auth msg.
func (s *Service) pub(mid int64, from, isAuthor int) (err error) {
	c := context.TODO()
	msg := &model.Msg{
		MID:       mid,
		From:      from,
		IsAuthor:  isAuthor,
		TimeStamp: time.Now().Unix(),
	}
	if err = s.upPub.Send(c, strconv.FormatInt(mid, 10), msg); err != nil {
		log.Error("pub mid(%d) error(%v)", mid, err)
	}
	return
}

func (s *Service) arcNotifyCanalConsume() {
	var err error
	for msg := range s.arcNotifySub.Messages() {
		msg.Commit()
		s.arcNotifyMo++
		m := &model.CanalMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if m.Table == _archive {
			s.arcNotifyMsg(m)
			log.Info("arcNotifyCanalConsume key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
		}
	}

	s.wg.Done()
}

func (s *Service) arcNotifyMsg(m *model.CanalMsg) {
	var (
		err    error
		newArc = &model.Archive{}
		oldArc = &model.Archive{}
	)
	if err = json.Unmarshal(m.New, newArc); err != nil {
		log.Error("arcNotifyMsg newMsg json.Unmarshal(%s) error(%v)", string(m.New), err)
		return
	}
	mid := newArc.MID
	if mid <= 0 {
		log.Error("arcNotifyMsg mid (%d) error", mid)
		return
	}
	if m.Action == _insert && newArc.State >= 0 { //0->1
		s.pub(mid, fromArchiveUp, 1)
	} else if m.Action == _update {
		if err = json.Unmarshal(m.Old, oldArc); err != nil {
			log.Error("arcNotifyMsg oldMsg json.Unmarshal(%s) error(%v)", string(m.Old), err)
			return
		}
		if oldArc.State < 0 && (newArc.State >= 0 || newArc.State == -6) { //0->1
			s.pub(mid, fromArchiveUp, 1)
		}
		if (oldArc.State >= 0 || oldArc.State == -6) && newArc.State < 0 { //1->0
			cnt, err := s.arc.UpCount(context.Background(), mid)
			if err != nil {
				log.Error("arcNotifyMsg s.arc.UpCount(%d) error(%v)", mid, err)
				return
			}
			if cnt <= 0 {
				s.pub(mid, fromArchiveUp, 0)
			}
		}
	}
}

func (s *Service) arcCanalConsume() {
	var err error
	for msg := range s.arcSub.Messages() {
		msg.Commit()
		s.arcMo++
		m := &model.CanalMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if m.Table == _archive && m.Action == _insert {
			arc := &model.Archive{}
			if err = json.Unmarshal(m.New, arc); err != nil {
				log.Error("creative-job binglog newMsg json.Unmarshal(%s) error(%v)", m.New, err)
				continue
			}
			if arc.MID > 0 {
				s.pub(arc.MID, fromArchiveNewUp, 1)
			}
			log.Info("arcCanalConsume key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
		}
	}
	s.wg.Done()
}

func (s *Service) task() {
	for msg := range s.taskSub.Messages() {
		msg.Commit()
		log.Info("databus task key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
		s.taskSubQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- msg
	}
	for _, c := range s.taskSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) share() {
	for msg := range s.shareSub.Messages() {
		msg.Commit()
		m := &model.ShareMsg{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus share json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if m.TP != 3 {
			continue
		}
		mid, authorMID := m.MID, s.getMIDByAID(m.OID)
		if mid != authorMID { //不是该用户分享的稿件则不做任何处理
			log.Warn("s.arc.Archive mid(%d)|author mid(%d)", mid, authorMID)
			continue
		}
		log.Info("databus share key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
		s.shareSubQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- m
	}

	for _, c := range s.shareSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) relation() {
	var err error
	newF, advancedF := s.c.Task.NewFollower, s.c.Task.AdvancedFollower
	for msg := range s.relationSub.Messages() {
		msg.Commit()

		rl := &model.RelaMessage{}
		if err = json.Unmarshal(msg.Value, rl); err != nil {
			log.Error("databus relation json.Unmarshal (%v) error(%v)", msg.Value, err)
			continue
		}

		if !strings.HasPrefix(rl.Table, _relationStatTable) && !strings.HasPrefix(rl.Table, _relationMidTable) {
			continue
		}
		if strings.HasPrefix(rl.Table, _relationStatTable) {
			ost := &model.Stat{}
			st := &model.Stat{}

			if rl.Action == "update" {
				if err = json.Unmarshal(rl.Old, ost); err != nil {
					log.Error("relation old msg json.Unmarshal(%s) error(%v)", string(rl.Old), err)
					continue
				}
			}

			if err = json.Unmarshal(rl.New, st); err != nil {
				log.Error("relation new msg json.Unmarshal(%s) error(%v)", string(rl.New), err)
				continue
			}

			isFollower := false
			if ost.Follower < newF && st.Follower >= newF { //新手任务粉丝数限制
				isFollower = true
			}
			if ost.Follower < advancedF && st.Follower >= advancedF { //进阶任务粉丝数限制
				isFollower = true
			}

			if isFollower {
				log.Info("databus relation key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
				s.followerQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- st
			}
		}

		if strings.HasPrefix(rl.Table, _relationMidTable) {
			fl := &model.Relation{}
			if err = json.Unmarshal(rl.New, fl); err != nil {
				log.Error("json.Unmarshal(%s) error(%v)", string(rl.New), err)
				continue
			}
			if fl.FID != s.c.Task.BiliMID { //过滤关注哔哩哔哩创作中心
				continue
			}
			log.Info("databus relation key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
			s.relationQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- fl
		}
	}

	for _, c := range s.relationQueue {
		close(c)
	}

	for _, c := range s.followerQueue {
		close(c)
	}

	s.wg.Done()
}

func (s *Service) statView() {
	statView, statViewUp := s.c.Task.StatView, s.c.Task.StatViewUp
	for msg := range s.statViewSub.Messages() {
		msg.Commit()
		m := &model.StatView{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus statView json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if !strings.EqualFold(m.Type, "archive") {
			continue
		}
		if m.Count >= statView && m.Count <= statViewUp {
			s.statViewSubQueue[s.shardingQueueIndex(msg.Key, s.statViewQueueLen)] <- m
		}
	}

	for _, c := range s.statViewSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) statLike() {
	statLike, statLikeUp := s.c.Task.StatLike, s.c.Task.StatLikeUp
	for msg := range s.statLikeSub.Messages() {
		msg.Commit()
		m := &model.StatLike{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus statLike json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if !strings.EqualFold(m.Type, "archive") {
			continue
		}
		if m.Count >= statLike && m.Count <= statLikeUp {
			s.statLikeSubQueue[s.shardingQueueIndex(msg.Key, s.statLikeQueueLen)] <- m
		}
	}

	for _, c := range s.statLikeSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) statShare() {
	statShare, statShareUp := s.c.Task.StatShare, s.c.Task.StatShareUp
	for msg := range s.statShareSub.Messages() {
		msg.Commit()
		m := &model.StatShare{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus statShare json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if !strings.EqualFold(m.Type, "archive") {
			continue
		}
		if m.Count >= statShare && m.Count <= statShareUp {
			log.Info("databus statShare key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
			s.statShareSubQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- m
		}
	}

	for _, c := range s.statShareSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) statCoin() {
	statCoin, statCoinUp := s.c.Task.StatCoin, s.c.Task.StatCoinUp
	for msg := range s.statCoinSub.Messages() {
		msg.Commit()
		m := &model.StatCoin{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus statCoin json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if !strings.EqualFold(m.Type, "archive") {
			continue
		}
		if m.Count >= statCoin && m.Count <= statCoinUp {
			log.Info("databus statCoin key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
			s.statCoinSubQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- m
		}
	}

	for _, c := range s.statCoinSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) statFav() {
	statFav, statFavUp := s.c.Task.StatFav, s.c.Task.StatFavUp
	for msg := range s.statFavSub.Messages() {
		msg.Commit()
		m := &model.StatFav{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus statFav json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if !strings.EqualFold(m.Type, "archive") {
			continue
		}
		if m.Count >= statFav && m.Count <= statFavUp {
			log.Info("databus statFav key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
			s.statFavSubQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- m
		}
	}

	for _, c := range s.statFavSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) statReply() {
	statReply, statReplyUp := s.c.Task.StatReply, s.c.Task.StatReplyUp
	for msg := range s.statReplySub.Messages() {
		msg.Commit()
		m := &model.StatReply{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus statReply json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if !strings.EqualFold(m.Type, "archive") || m.Count < s.c.Task.StatReply {
			continue
		}
		if m.Count >= statReply && m.Count <= statReplyUp {
			log.Info("databus statReply key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
			s.statReplySubQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- m
		}
	}

	for _, c := range s.statReplySubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) statDM() {
	statDM, statDMUp := s.c.Task.StatDM, s.c.Task.StatDMUp
	for msg := range s.statDMSub.Messages() {
		msg.Commit()
		m := &model.StatDM{}
		if err := json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus statDM json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}

		if !strings.EqualFold(m.Type, "archive") || m.Count < s.c.Task.StatDM {
			continue
		}
		if m.Count >= statDM && m.Count <= statDMUp {
			log.Info("databus statDM key(%s) value(%s) partition(%d) offset(%d) commit", msg.Key, msg.Value, msg.Partition, msg.Offset)
			s.statDMSubQueue[s.shardingQueueIndex(msg.Key, s.databusQueueLen)] <- m
		}
	}

	for _, c := range s.statDMSubQueue {
		close(c)
	}
	s.wg.Done()
}

func (s *Service) newUp() {
	var err error
	for msg := range s.newUpSub.Messages() {
		msg.Commit()
		m := &model.CanalMsg{}
		if err = json.Unmarshal(msg.Value, m); err != nil {
			log.Error("databus newUp json.Unmarshal(%v) error(%v)", string(msg.Value), err)
			continue
		}
		if m.Table != _archive {
			continue
		}

		newArc := &model.Archive{}
		if err = json.Unmarshal(m.New, newArc); err != nil {
			log.Error("databus newUp newMsg json.Unmarshal(%s) error(%v)", string(m.New), err)
			continue
		}
		mid, aid := newArc.MID, newArc.AID
		if mid == 0 || aid == 0 {
			log.Error("databus newUp mid(%d) | aid(%d) error", mid, aid)
			continue
		}

		isUp := false
		if m.Action == "insert" && newArc.State >= 0 {
			isUp = true
		} else if m.Action == "update" {
			oldArc := &model.Archive{}
			if err = json.Unmarshal(m.Old, oldArc); err != nil {
				log.Error("newUp oldMsg json.Unmarshal(%s) error(%v)", string(m.Old), err)
				continue
			}
			if oldArc.State < 0 && (newArc.State >= 0 || newArc.State == -6) { //0->1
				isUp = true
			}
		}

		if isUp {
			av, err := s.arc.View(context.Background(), mid, aid) //获取投稿来源
			if err != nil {
				log.Error("newUp s.arc.View mid(%d) mid(%d) av(%+v) error(%v)", mid, aid, av, err)
			}

			if av != nil && av.Archive != nil && (av.Archive.UpFrom == 3 || av.Archive.UpFrom == 8 || av.Archive.UpFrom == 9) { // 3-App , 8-Android , 9-IOS
				log.Info("databus mobile mid(%d) aid(%d) av(%+v)", mid, aid, av)
				s.mobileUpQueue[s.shardingQueueIndex(strconv.FormatInt(mid, 10), s.databusQueueLen)] <- &model.Up{AID: aid, MID: mid}
			}

			cnt, err := s.arc.UpCount(context.Background(), mid)
			if err != nil {
				log.Error("newUp s.arc.UpCount(%d) error(%v)", mid, err)
				continue
			}
			if cnt == 1 { //新手投下自己的第一个稿件
				log.Info("databus newUp mid(%d) aid(%d) count(%d)", mid, aid, cnt)
				s.newUpQueue[s.shardingQueueIndex(strconv.FormatInt(mid, 10), s.databusQueueLen)] <- &model.Up{AID: aid, MID: mid}
			}
			if cnt >= 5 { //进阶任务视频投稿超过5个
				log.Info("databus oldUp mid(%d) aid(%d) count(%d)", mid, aid, cnt)
				s.oldUpQueue[s.shardingQueueIndex(strconv.FormatInt(mid, 10), s.databusQueueLen)] <- &model.Up{AID: aid, MID: mid}
			}
		}
	}

	for _, c := range s.newUpQueue {
		close(c)
	}
	for _, c := range s.oldUpQueue {
		close(c)
	}
	for _, c := range s.mobileUpQueue {
		close(c)
	}
	s.wg.Done()
}
