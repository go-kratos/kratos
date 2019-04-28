package show

import (
	"context"
	"hash/crc32"
	"strconv"
	"time"

	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/card"
	"go-common/app/interface/main/app-show/model/feed"
	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
)

var (
	_emptyList = []*feed.Item{}
)

// FeedIndex feed index
func (s *Service) FeedIndex(c context.Context, mid, idx int64, plat int8, build, loginEvent int, lastParam, mobiApp, device, buvid string, now time.Time) (res []*feed.Item) {
	var (
		ps               = 10
		isIpad           = plat == model.PlatIPad
		cards, cardCache []*card.PopularCard
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
		res = _emptyList
		return
	}
	res = s.dealItem(c, plat, build, ps, cards, idx, lastParam, now)
	if len(res) == 0 {
		res = _emptyList
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
		feed:       res,
	}
	s.infocfeed(infoc)
	return
}

// dealItem feed item
func (s *Service) dealItem(c context.Context, plat int8, build, ps int, cards []*card.PopularCard, idx int64, lastParam string, now time.Time) (is []*feed.Item) {
	const _rankCount = 3
	var (
		uri map[int64]string
		// key
		max             = int64(100)
		_fTypeOperation = "operation"
		aids            []int64
		am              map[int64]*api.Arc
		feedcards       []*card.PopularCard
		err             error
	)
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
		tmp.Idx = cardIdx
		feedcards = append(feedcards, tmp)
		switch ca.Type {
		case model.GotoAv:
			aids = append(aids, ca.Value)
		}
		if len(feedcards) == ps {
			break
		}
	}
	if len(aids) != 0 {
		if am, err = s.arc.ArchivesPB(c, aids); err != nil {
			s.pMiss.Incr("popularcard_Archives")
			err = nil
		} else {
			s.pHit.Incr("popularcard_Archives")
		}
	}
	for _, ca := range feedcards {
		i := &feed.Item{}
		i.FromType = ca.FromType
		i.Idx = ca.Idx
		i.Pos = ca.Pos
		switch ca.Type {
		case model.GotoAv:
			a := am[ca.Value]
			isOsea := model.IsOverseas(plat)
			if a != nil && a.IsNormal() && (!isOsea || (isOsea && a.AttrVal(archive.AttrBitOverseaLock) == 0)) {
				i.FromPlayerAv(a, uri[a.Aid])
				i.FromRcmdReason(ca)
				// if tag, ok := s.hotArcTag[a.Aid]; ok {
				// 	i.Tag = &feed.Tag{TagID: tag.ID, TagName: tag.Name}
				// }
				i.Goto = ca.Type
				is = append(is, i)
			}
		case model.GotoRank:
			if rankAids := s.rankAidsCache; len(rankAids) >= _rankCount {
				i.FromRank(rankAids, s.rankScoreCache, s.rankArchivesCache)
				// i.Param = strconv.FormatInt(ca.ID, 10)
				if i.Goto != "" {
					is = append(is, i)
				}
			}
		case model.GotoHotTopic:
			if hotTopics := s.hottopicsCache; len(hotTopics) > 0 {
				i.FromHotTopic(hotTopics)
				is = append(is, i)
			}
		}
	}
	if rl := len(is); rl == 0 {
		is = _emptyList
		return
	}
	return
}
