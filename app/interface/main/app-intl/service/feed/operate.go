package feed

import (
	"context"
	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card/operate"
)

// channelRcmdCard is.
func (s *Service) channelRcmdCard(c context.Context, ids ...int64) (cardm map[int64]*operate.Card, aids, tids []int64) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if o, ok := s.followCache[id]; ok {
			card := &operate.Card{}
			card.FromFollow(o)
			cardm[id] = card
			switch card.Goto {
			case model.GotoAv:
				if card.ID != 0 {
					aids = append(aids, card.ID)
				}
				if card.Tid != 0 {
					tids = append(tids, card.Tid)
				}
			}
		}
	}
	return
}
