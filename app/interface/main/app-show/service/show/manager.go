package show

import (
	"context"

	"go-common/app/interface/main/app-card/model/card/operate"
)

func (s *Service) cardSetChange(c context.Context, ids ...int64) (cardm map[int64]*operate.Card, aids []int64, upid int64) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if cs, ok := s.cardSetCache[id]; ok {
			card := &operate.Card{}
			card.FromCardSet(cs)
			cardm[id] = card
			upid = card.ID
			for _, item := range card.Items {
				switch cs.Type {
				case "up_rcmd_new":
					aids = append(aids, item.ID)
				}
			}
		}
	}
	return
}

func (s *Service) eventTopicChange(c context.Context, ids ...int64) (cardm map[int64]*operate.Card) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if st, ok := s.eventTopicCache[id]; ok {
			card := &operate.Card{}
			card.FromEventTopic(st)
			cardm[id] = card
		}
	}
	return
}
