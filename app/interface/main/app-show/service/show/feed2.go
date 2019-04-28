package show

import (
	"context"
	"hash/crc32"
	"strconv"
	"time"

	cdm "go-common/app/interface/main/app-card/model"
	cardm "go-common/app/interface/main/app-card/model/card"
	operate "go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-card/model/card/rank"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/card"
	"go-common/app/interface/main/app-show/model/feed"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	relation "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

var (
	_emptyList2 = []cardm.Handler{}
)

// FeedIndex feed index
func (s *Service) FeedIndex2(c context.Context, mid, idx int64, plat int8, build, loginEvent int, lastParam, mobiApp, device, buvid string, now time.Time) (res []cardm.Handler, ver string, err error) {
	var (
		ps               = 10
		isIpad           = plat == model.PlatIPad
		cards, cardCache []*card.PopularCard
		infocs           []*feed.Item
		style            int8
	)
	if isIpad {
		ps = 20
	}
	var key int
	if mid > 0 {
		key = int((mid / 1000) % 10)
	} else {
		key = int((crc32.ChecksumIEEE([]byte(buvid)) / 1000) % 10)
	}
	cardCache = s.PopularCardTenList(c, key)
	if len(cardCache) > int(idx) {
		cards = cardCache[idx:]
	} else {
		err = ecode.AppNotData
		res = _emptyList2
		return
	}
	// HotDynamic====================
	// cards = append(cards[:0], append([]*card.PopularCard{&card.PopularCard{Type: model.GotoHotDynamic, ReasonType: 0, FromType: "recommend"}}, cards[0:]...)...)
	// HotDynamic====================
	//build
	if plat == model.PlatIPhone && build > 8230 || plat == model.PlatAndroid && build > 5345000 {
		// switch key {
		// // case 0, 3:
		// // 	style = cdm.HotCardStyleShowUp
		// case 2, 5:
		// 	style = cdm.HotCardStyleHideUp
		// default:
		// 	style = cdm.HotCardStyleOld
		// }
		style = cdm.HotCardStyleHideUp
	} else {
		style = cdm.HotCardStyleOld
	}
	//build
	res, infocs = s.dealItem2(c, plat, build, ps, cards, mid, idx, lastParam, now, style)
	ver = strconv.FormatInt(now.Unix(), 10)
	if len(res) == 0 {
		err = ecode.AppNotData
		res = _emptyList2
		return
	}
	//infoc
	infoc := &feedInfoc{
		mobiApp:    mobiApp,
		device:     device,
		build:      strconv.Itoa(build),
		now:        now.Format("2006-01-02 15:04:05"),
		loginEvent: strconv.Itoa(loginEvent),
		mid:        strconv.FormatInt(mid, 10),
		buvid:      buvid,
		page:       strconv.Itoa((int(idx) / ps) + 1),
		feed:       infocs,
	}
	s.infocfeed(infoc)
	return
}

// dealItem feed item
func (s *Service) dealItem2(c context.Context, plat int8, build, ps int, cards []*card.PopularCard, mid, idx int64, lastParam string, now time.Time, style int8) (is []cardm.Handler, infocs []*feed.Item) {
	var (
		max                           = int64(100)
		_fTypeOperation               = "operation"
		aids, avUpIDs, upIDs, rnUpIDs []int64
		am                            map[int64]*archive.ArchiveWithPlayer
		feedcards                     []*card.PopularCard
		err                           error
		rank                          *operate.Card
		accountm                      map[int64]*account.Card
		isAtten                       map[int64]int8
		statm                         map[int64]*relation.Stat
	)
	cardSet := map[int64]*operate.Card{}
	eventTopic := map[int64]*operate.Card{}
LOOP:
	for pos, ca := range cards {
		var cardIdx = idx + int64(pos+1)
		if cardIdx > max && ca.FromType != _fTypeOperation {
			continue
		}
		if config, ok := ca.PopularCardPlat[plat]; ok {
			for _, l := range config {
				if model.InvalidBuild(build, l.Build, l.Condition) {
					continue LOOP
				}
			}
		} else if ca.FromType == _fTypeOperation {
			continue LOOP
		}
		tmp := &card.PopularCard{}
		*tmp = *ca
		if plat == model.PlatIPhone && build > 8290 || plat == model.PlatAndroid && build > 5365000 {
			switch ca.Type {
			case model.GotoUpRcmdNew:
				tmp.Type = model.GotoUpRcmdNewV2
			}
		}
		tmp.Idx = cardIdx
		feedcards = append(feedcards, tmp)
		switch ca.Type {
		case model.GotoAv:
			aids = append(aids, ca.Value)
		case model.GotoRank:
			rank = &operate.Card{}
			rank.FromRank(s.rankCache2)
		case model.GotoUpRcmdNew, model.GotoUpRcmdNewV2:
			cardm, as, upid := s.cardSetChange(c, ca.Value)
			aids = append(aids, as...)
			for id, card := range cardm {
				cardSet[id] = card
			}
			rnUpIDs = append(rnUpIDs, upid)
		case model.GotoEventTopic:
			eventTopic = s.eventTopicChange(c, ca.Value)
		}
		if len(feedcards) == ps {
			break
		}
	}
	if len(aids) != 0 {
		var as map[int64]*api.Arc
		if as, err = s.arc.ArchivesPB(c, aids); err != nil {
			log.Error("%+v", err)
			err = nil
		} else {
			am = map[int64]*archive.ArchiveWithPlayer{}
			for _, a := range as {
				avUpIDs = append(avUpIDs, a.Author.Mid)
				am[a.Aid] = &archive.ArchiveWithPlayer{Archive3: archive.BuildArchive3(a)}
			}
		}
	}
	switch style {
	case cdm.HotCardStyleShowUp, cdm.HotCardStyleHideUp:
		upIDs = append(upIDs, avUpIDs...)
	}
	upIDs = append(upIDs, rnUpIDs...)
	avUpIDs = append(avUpIDs, rnUpIDs...)
	g, ctx := errgroup.WithContext(c)
	if len(avUpIDs) > 0 {
		g.Go(func() (err error) {
			if accountm, err = s.acc.Cards3(ctx, avUpIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return nil
		})
	}
	if len(upIDs) > 0 {
		g.Go(func() (err error) {
			if statm, err = s.reldao.Stats(ctx, upIDs); err != nil {
				log.Error("%+v", err)
				err = nil
			}
			return nil
		})
		if mid != 0 {
			g.Go(func() error {
				isAtten = s.acc.IsAttention(ctx, upIDs, mid)
				return nil
			})
		}
	}
	g.Wait()
	for _, ca := range feedcards {
		var (
			r        = ca.PopularCardToAiChange()
			main     interface{}
			cardType cdm.CardType
		)
		r.Style = style
		op := &operate.Card{}
		op.From(cdm.CardGt(r.Goto), r.ID, 0, plat, build)
		switch r.Style {
		case cdm.HotCardStyleShowUp, cdm.HotCardStyleHideUp:
			switch r.Goto {
			case model.GotoAv:
				cardType = cdm.SmallCoverV5
			}
		}
		switch r.Goto {
		case model.GotoAv:
			if a, ok := am[r.ID]; ok && (a.AttrVal(archive.AttrBitOverseaLock) == 0 || !model.IsOverseas(plat)) {
				main = map[int64]*archive.ArchiveWithPlayer{a.Aid: a}
				// op.Tid = a.Aid
				r.HideButton = true
				if (plat == model.PlatIPhone && build > 8290 || plat == model.PlatAndroid && build > 5365000) && cardType == cdm.SmallCoverV5 {
					op.Switch = cdm.SwitchCooperationShow
				} else {
					op.Switch = cdm.SwitchCooperationHide
				}
			}
		case model.GotoRank:
			ams := map[int64]*archive.ArchiveWithPlayer{}
			for aid, a := range s.rankArchivesCache {
				ams[aid] = &archive.ArchiveWithPlayer{Archive3: archive.BuildArchive3(a)}
			}
			main = map[cdm.Gt]interface{}{cdm.GotoAv: ams}
			op = rank
		case model.GotoHotTopic:
			main = s.hottopicsCache
		case model.GotoUpRcmdNew, model.GotoUpRcmdNewV2:
			main = am
			op = cardSet[r.ID]
		case model.GotoHotDynamic:
			main = s.dynamicHotCache
		case model.GotoEventTopic:
			op = eventTopic[r.ID]
		}
		h := cardm.Handle(plat, cdm.CardGt(r.Goto), cardType, cdm.ColumnSvrSingle, r, nil, isAtten, statm, accountm)
		if h == nil {
			continue
		}
		h.From(main, op)
		h.Get().FromType = ca.FromType
		h.Get().Idx = ca.Idx
		if h.Get().Right {
			h.Get().ThreePointWatchLater()
			is = append(is, h)
		}
		// infoc
		tinfo := &feed.Item{
			Goto:       ca.Type,
			Param:      strconv.FormatInt(ca.Value, 10),
			URI:        h.Get().URI,
			FromType:   ca.FromType,
			Idx:        h.Get().Idx,
			CornerMark: ca.CornerMark,
			CardStyle:  r.Style,
		}
		if r.RcmdReason != nil {
			tinfo.RcmdContent = r.RcmdReason.Content
		}
		if op != nil {
			switch r.Goto {
			case model.GotoEventTopic:
				tinfo.Item = append(tinfo.Item, &feed.Item{Param: op.URI, Goto: string(op.Goto)})
			default:
				for _, tmp := range op.Items {
					tinfo.Item = append(tinfo.Item, &feed.Item{Param: strconv.FormatInt(tmp.ID, 10), Goto: string(tmp.Goto)})
				}
			}
		}
		infocs = append(infocs, tinfo)
		if len(is) == ps {
			break
		}
	}
	rl := len(is)
	if rl == 0 {
		is = _emptyList2
		return
	}
	return
}

func (s *Service) RankCard() (ranks []*rank.Rank, aids []int64) {
	const _limit = 3
	ranks = make([]*rank.Rank, 0, _limit)
	aids = make([]int64, 0, _limit)
	for _, rank := range s.rankCache2 {
		ranks = append(ranks, rank)
		aids = append(aids, rank.Aid)
		if len(ranks) == _limit {
			break
		}
	}
	return
}
