package service

import (
	"context"
	"time"

	"go-common/app/job/main/member/model/queue"
	"go-common/library/log"
)

// Item is
type Item struct {
	Mid    int64
	Time   time.Time
	Action string
}

// Compare is
func (i *Item) Compare(other queue.Item) int {
	o := asItem(other)
	if o == nil {
		return -1
	}
	if i.Time.Equal(o.Time) {
		return 0
	}
	if i.Time.After(o.Time) {
		return 1
	}
	return -1
}

// HashCode is
func (i *Item) HashCode() int64 {
	return i.Mid
}

func asItem(in queue.Item) *Item {
	o, ok := in.(*Item)
	if !ok {
		return nil
	}
	return o
}

func asItems(in []queue.Item) []*Item {
	out := make([]*Item, 0, len(in))
	for _, i := range in {
		item := asItem(i)
		if item == nil {
			continue
		}
		out = append(out, item)
	}
	return out
}

func (s *Service) cachedelayproc(ctx context.Context) {
	fiveSeconds := time.Second * 5
	t := time.NewTicker(fiveSeconds)

	delayed := func(t time.Time) bool {
		top := asItem(s.cachepq.Peek())
		if top == nil {
			log.Info("Empty cache queue top at: %v", t)
			return false
		}
		if t.Sub(top.Time) < fiveSeconds {
			log.Info("Top item is in five seconds, skip and waiting for next tick")
			return false
		}
		return true
	}

	for ti := range t.C {
		if !delayed(ti) {
			continue
		}

		for {
			qitems, err := s.cachepq.Get(1)
			if err != nil {
				log.Error("Failed to get queue items from cache queue: %+v", err)
				return
			}
			items := asItems(qitems)
			for _, it := range items {
				log.Info("Notify purge cache in delay queue with mid: %d", it.Mid)
				s.dao.NotifyPurgeCache(ctx, it.Mid, it.Action)
			}
			if s.cachepq.Empty() || !delayed(time.Now()) {
				break
			}
		}
	}
}
