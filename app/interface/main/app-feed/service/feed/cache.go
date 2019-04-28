package feed

import (
	"context"
	"hash/crc32"
	"math/rand"
	"time"

	"go-common/app/interface/main/app-card/model/card/ai"
	"go-common/app/interface/main/app-card/model/card/live"
	"go-common/app/interface/main/app-feed/model"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

func (s *Service) indexCache(c context.Context, mid int64, count int) (rs []*ai.Item, err error) {
	var (
		pos, nextPos int
	)
	cache := s.rcmdCache
	if len(cache) < count {
		return
	}
	if pos, err = s.rcmd.PositionCache(c, mid); err != nil {
		return
	}
	rs = make([]*ai.Item, 0, count)
	if pos < len(cache)-count-1 {
		nextPos = pos + count
		rs = append(rs, cache[pos:nextPos]...)
	} else if pos < len(cache)-1 {
		nextPos = count - (len(cache) - pos)
		rs = append(rs, cache[pos:]...)
		rs = append(rs, cache[:nextPos]...)
	} else {
		nextPos = count - 1
		rs = append(rs, cache[:nextPos]...)
	}
	s.addCache(func() {
		s.rcmd.AddPositionCache(context.Background(), mid, nextPos)
	})
	return
}

func (s *Service) recommendCache(count int) (rs []*ai.Item) {
	cache := s.rcmdCache
	index := len(cache)
	if count > 0 && count < index {
		index = count
	}
	rs = make([]*ai.Item, 0, index)
	for _, idx := range rand.Perm(len(cache))[:index] {
		rs = append(rs, cache[idx])
	}
	return
}

func (s *Service) group(mid int64, buvid string) (group int) {
	if mid == 0 && buvid == "" {
		group = -1
		return
	}
	if mid != 0 {
		if v, ok := s.groupCache[mid]; ok {
			group = v
			return
		}
		group = int(mid % 20)
		return
	}
	group = int(crc32.ChecksumIEEE([]byte(buvid)) % 20)
	return
}

func (s *Service) loadRcmdCache() {
	is, err := s.rcmd.RcmdCache(context.Background())
	if err != nil {
		log.Error("%+v", err)
	}
	if len(is) >= 50 {
		for _, i := range is {
			i.Goto = model.GotoAv
		}
		s.rcmdCache = is
		return
	}
	aids, err := s.rcmd.Hots(context.Background())
	if err != nil {
		log.Error("%+v", err)
	}
	if len(aids) == 0 {
		if aids, err = s.rcmd.RcmdAidsCache(context.Background()); err != nil {
			log.Error("%+v", err)
			return
		}
	}
	if len(aids) < 50 && len(s.rcmdCache) != 0 {
		return
	}
	s.addCache(func() {
		s.rcmd.AddRcmdAidsCache(context.Background(), aids)
	})
	if is, err = s.fromArchvies(aids); err != nil {
		log.Error("%+v", err)
		return
	}
	s.rcmdCache = is
}

func (s *Service) UpRcmdCache(c context.Context, is []*ai.Item) (err error) {
	if err = s.rcmd.AddRcmdCache(c, is); err != nil {
		log.Error("%+v", err)
	}
	return
}

func (s *Service) fromArchvies(aids []int64) (is []*ai.Item, err error) {
	var as map[int64]*archive.ArchiveWithPlayer
	if as, err = s.arc.ArchivesWithPlayer(context.Background(), aids, 0, "", 0, 0, 0, 0); err != nil {
		return
	}
	is = make([]*ai.Item, 0, len(aids))
	for _, aid := range aids {
		a, ok := as[aid]
		if !ok || a.Archive3 == nil || !a.IsNormal() {
			continue
		}
		is = append(is, &ai.Item{ID: aid, Goto: model.GotoAv, Archive: a.Archive3})
	}
	return
}

func (s *Service) rcmdproc() {
	for {
		time.Sleep(s.tick)
		s.loadRcmdCache()
	}
}

func (s *Service) loadRankCache() {
	rank, err := s.rank.AllRank(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.rankCache = rank
}

func (s *Service) rankproc() {
	for {
		time.Sleep(s.tick)
		s.loadRankCache()
	}
}

func (s *Service) loadConvergeCache() {
	converge, err := s.cvg.Cards(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.convergeCache = converge
}

func (s *Service) convergeproc() {
	for {
		time.Sleep(s.tick)
		s.loadConvergeCache()
	}
}

func (s *Service) loadDownloadCache() {
	download, err := s.gm.DownLoad(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.downloadCache = download
}

func (s *Service) downloadproc() {
	for {
		time.Sleep(s.tick)
		s.loadDownloadCache()
	}
}

func (s *Service) loadSpecialCache() {
	special, err := s.sp.Card(context.Background(), time.Now())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	var roomIDs []int64
	idm := map[int64]int64{}
	for _, sp := range special {
		if sp.Goto == model.GotoLive && sp.Pid != 0 {
			roomIDs = append(roomIDs, sp.Pid)
			idm[sp.Pid] = sp.ID
		}
	}
	room, err := s.lv.Rooms(context.Background(), roomIDs, "")
	if err != nil {
		log.Error("%+v", err)
	}
	if len(room) != 0 {
		for rid, id := range idm {
			if r, ok := room[rid]; !ok || r.LiveStatus != 1 {
				delete(special, id)
			}
		}
	}
	s.specialCache = special
}

func (s *Service) specialproc() {
	for {
		time.Sleep(s.tick)
		s.loadSpecialCache()
	}
}

func (s *Service) loadGroupCache() {
	group, err := s.rcmd.Group(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.groupCache = group
}

func (s *Service) groupproc() {
	for {
		time.Sleep(s.tick)
		s.loadGroupCache()
	}
}

func (s *Service) loadFollowModeList() {
	list, err := s.rcmd.FollowModeList(context.Background())
	if err != nil {
		log.Error("%+v", err)
		if list, err = s.rcmd.FollowModeListCache(context.Background()); err != nil {
			log.Error("%+v", err)
			return
		}
	} else {
		s.addCache(func() {
			s.rcmd.AddFollowModeListCache(context.Background(), list)
		})
	}
	log.Warn("loadFollowModeList list len(%d)", len(list))
	s.followModeList = list
}

func (s *Service) followModeListproc() {
	for {
		time.Sleep(s.tick)
		s.loadFollowModeList()
	}
}

func (s *Service) loadUpCardCache() {
	follow, err := s.card.Follow(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.followCache = follow
}

func (s *Service) upCardproc() {
	for {
		time.Sleep(s.tick)
		s.loadUpCardCache()
	}
}

func (s *Service) loadLiveCardCache() {
	liveCard, err := s.lv.Card(context.Background())
	if err != nil {
		log.Error("%+v", err)
		return
	}
	s.liveCardCache = liveCard
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

func (s *Service) liveCardproc() {
	for {
		time.Sleep(1 * time.Second)
		s.loadLiveCardCache()
	}
}

func (s *Service) loadABTestCache() {
	res, err := s.rsc.AbTest(context.Background(), _feedgroups)
	if err != nil {
		log.Error("resource s.rsc.AbTest error(%v)", err)
		return
	}
	s.abtestCache = res
	log.Info("loadAbTestCache cache success")
}

func (s *Service) loadABTestCacheProc() {
	for {
		time.Sleep(s.tick)
		s.loadABTestCache()
	}
}

func (s *Service) loadAutoPlayMid() {
	tmp := map[int64]struct{}{}
	for _, mid := range s.c.AutoPlayMids {
		tmp[mid] = struct{}{}
	}
	s.autoplayMidsCache = tmp
}
