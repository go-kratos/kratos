package channel

import (
	"context"
	"go-common/app/interface/main/app-card/model"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-card/model/card/operate"
)

func (s *Service) convergeCard2(c context.Context, limit int, ids ...int64) (cardm map[int64]*operate.Card, aids, roomIDs, metaIDs []int64) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for i, id := range ids {
		if o, ok := s.convergeCardCache[id]; ok {
			card := &operate.Card{}
			card.FromConverge(o)
			cardm[id] = card
			for _, item := range card.Items {
				switch item.Goto {
				case model.GotoAv:
					if item.ID != 0 {
						aids = append(aids, item.ID)
					}
				case model.GotoLive:
					if item.ID != 0 {
						roomIDs = append(roomIDs, item.ID)
					}
				case model.GotoArticle:
					if item.ID != 0 {
						metaIDs = append(metaIDs, item.ID)
					}
				}
				if i == limit-1 {
					break
				}
			}
		}
	}
	return
}

func (s *Service) downloadCard(c context.Context, ids ...int64) (cardm map[int64]*operate.Card) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if o, ok := s.gameDownloadCache[id]; ok {
			card := &operate.Card{}
			card.FromDownload(o)
			cardm[id] = card
		}
	}
	return
}

func (s *Service) subscribeCard(c context.Context, ids ...int64) (cardm map[int64]*operate.Card, upIDs, tids []int64) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if o, ok := s.upCardCache[id]; ok {
			card := &operate.Card{}
			card.FromFollow(o)
			cardm[id] = card
			for _, item := range card.Items {
				switch item.Goto {
				case model.GotoMid:
					if item.ID != 0 {
						upIDs = append(upIDs, item.ID)
					}
				case model.GotoTag:
					if item.ID != 0 {
						tids = append(tids, item.ID)
					}
				}
			}
		}
	}
	return
}

func (s *Service) channelRcmdCard(c context.Context, ids ...int64) (cardm map[int64]*operate.Card, aids, tids []int64) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if o, ok := s.upCardCache[id]; ok {
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

func (s *Service) liveUpRcmdCard(c context.Context, ids ...int64) (cardm map[int64][]*live.Card, upIDs []int64) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64][]*live.Card, len(ids))
	for _, id := range ids {
		if card, ok := s.liveCardCache[id]; ok {
			cardm[id] = card
			for _, c := range card {
				if c.UID != 0 {
					upIDs = append(upIDs, c.UID)
				}
			}
		}
	}
	return
}

func (s *Service) specialCard(c context.Context, ids ...int64) (cardm map[int64]*operate.Card) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if o, ok := s.specialCardCache[id]; ok {
			card := &operate.Card{}
			card.FromSpecial(o)
			cardm[id] = card
		}
	}
	return
}

func (s *Service) topstickCard(c context.Context, ids ...int64) (cardm map[int64]*operate.Card) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if o, ok := s.specialCardCache[id]; ok {
			card := &operate.Card{}
			card.FromTopstick(o)
			cardm[id] = card
		}
	}
	return
}

func (s *Service) cardSetChange(c context.Context, ids ...int64) (cardm map[int64]*operate.Card, seasonIDs []int32) {
	if len(ids) == 0 {
		return
	}
	cardm = make(map[int64]*operate.Card, len(ids))
	for _, id := range ids {
		if cs, ok := s.cardSetCache[id]; ok {
			card := &operate.Card{}
			card.FromCardSet(cs)
			cardm[id] = card
			for _, item := range card.Items {
				switch cs.Type {
				case "pgcs_rcmd":
					seasonIDs = append(seasonIDs, int32(item.ID))
				}
			}
		}
	}
	return
}
