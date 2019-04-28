package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"go-common/app/job/main/videoup/model/archive"
	"go-common/app/job/main/videoup/model/message"
	"go-common/app/job/main/videoup/model/redis"
	"go-common/library/log"
)

func (s *Service) syncRetry(c context.Context, aid, mid int64, action, route, content string) (err error) {
	retry := &redis.RetryJSON{}
	retry.Action = action
	retry.Data.Aid = aid
	retry.Data.Route = route
	retry.Data.Mid = mid
	retry.Data.Content = content
	if action == redis.ActionForVideocovers && (content == "" || strings.Contains(content, "/bfs/archive")) {
		return
	}
	s.redis.PushFail(c, retry)
	return
}

func (s *Service) retryproc() {
	defer s.wg.Done()
	for {
		if s.closed {
			return
		}
		var (
			c     = context.TODO()
			bs    []byte
			err   error
			retry = &redis.RetryJSON{}
		)
		bs, err = s.redis.PopFail(c)
		if err != nil || bs == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		msg := &redis.Retry{}
		if err = json.Unmarshal(bs, msg); err != nil {
			log.Error("json.Unretry syncmarshal(%s) error(%v)", bs, err)
			continue
		}
		log.Info("retry %s %s", msg.Action, bs)
		if err = json.Unmarshal(bs, retry); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg, err)
			continue
		}
		s.promRetry.Incr(msg.Action)
		switch msg.Action {
		case redis.ActionForBvcCapable:
			var a *archive.Archive
			if a, err = s.arc.Archive(c, retry.Data.Aid); err != nil {
				log.Error("retry bvcCapable archive(%d) error(%v)", retry.Data.Aid, err)
				continue
			}
			log.Info("retry aid(%d) syncBVC  bvcCapable", a.Aid)
			s.syncBVC(c, a)
		case redis.ActionForSendOpenMsg:
			s.sendAuditMsg(c, retry.Data.Route, retry.Data.Aid)
		case redis.ActionForSendBblog:
			var a *archive.Archive
			dynamic := ""
			if a, err = s.arc.Archive(c, retry.Data.Aid); err != nil {
				log.Error("retry sendBblog archive(%d) error(%v)", retry.Data.Aid, err)
				s.syncRetry(c, retry.Data.Aid, retry.Data.Mid, redis.ActionForSendBblog, "", "")
				continue
			}
			if add, _ := s.arc.Addit(c, a.Aid); add != nil {
				dynamic = add.Dynamic
			}
			s.sendBblog(&archive.Result{Aid: a.Aid, Mid: a.Mid, Dynamic: dynamic})
		case redis.ActionForVideoshot:
			imgs := strings.Split(retry.Data.Content, ",")
			s.videoshotAdd(retry.Data.Aid, retry.Data.Route, imgs)
		case redis.ActionForVideocovers:
			if retry.Data.Mid > 20 {
				continue
			}
			var a *archive.Archive
			if a, err = s.arc.Archive(c, retry.Data.Aid); err != nil {
				s.syncRetry(c, retry.Data.Aid, retry.Data.Mid+1, redis.ActionForVideocovers, retry.Data.Content, retry.Data.Content)
				continue
			}
			log.Info("retryproc videocoverCopy aid(%d) old cover is(%s) retry count is(%d)", retry.Data.Aid, a.Cover, retry.Data.Mid)
			if strings.Index(a.Cover, "//") == 0 {
				a.Cover = "http:" + a.Cover
			}
			if strings.Index(a.Cover, "/bfs") == 0 && !strings.Contains(a.Cover, "bfs/archive") {
				a.Cover = "http://i0.hdslb.com" + a.Cover
			}
			if a == nil || a.Cover == "" || strings.Contains(a.Cover, "bfs/archive") || (!strings.HasPrefix(a.Cover, "http://") && !strings.HasPrefix(a.Cover, "https://")) {
				log.Error("retryproc videocoverCopy aid(%d) cover is(%s)", retry.Data.Aid, a.Cover)
				continue
			}
			s.videocoverCopy(retry.Data.Aid, retry.Data.Mid, a)
		case redis.ActionForPostFirstRound:
			pmsg := &message.Videoup{}
			if err = json.Unmarshal([]byte(retry.Data.Content), pmsg); err != nil {
				log.Error("retryproc postFirstRound json.Unmarshal(%s) error(%v)", retry.Data.Content, err)
				continue
			}
			s.sendPostFirstRound(c, retry.Data.Route, pmsg.Aid, pmsg.Filename, pmsg.AdminChange)
		default:
			log.Warn("retryproc unknown action(%s) message(%+v)", msg.Action, retry)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

//QueueProc .
func (s *Service) QueueProc() {
	defer s.wg.Done()
	for {
		if s.closed {
			return
		}

		var (
			c   = context.TODO()
			bs  []byte
			err error
		)
		bs, err = s.redis.PopQueue(c, message.RouteVideoshotpv)
		if err != nil || bs == nil {
			time.Sleep(5 * time.Second)
			continue
		}
		m := &message.BvcVideo{}
		if err = json.Unmarshal(bs, m); err != nil {
			log.Error("QueueProc json.Unmarshal(%v) error(%v)", string(bs), err)
			continue
		}
		log.Info("queue proc pop %+v", m)
		switch m.Route {
		case message.RouteVideoshotpv:
			err = s.videoshotPv(c, m)
		default:
			log.Warn("QueueProc unknown route(%s) message(%s)", m.Route, m.Route)
		}
		if err != nil {
			log.Error("QueueProc  error(%+v)", err)
		}
		time.Sleep(10 * time.Millisecond)

	}
}
