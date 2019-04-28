package service

import (
	"context"
	"strings"
	"sync/atomic"
	"time"

	"go-common/app/job/main/click/model"
	"go-common/library/log"
)

func (s *Service) isAllow(ctx context.Context, c *model.ClickMsg) (rtype int8, err error) {
	var (
		f        *model.Forbid
		ok       bool
		duration int64
	)
	rtype = model.LogTypeForNotUse
	// 自动播放的恶心逻辑 开始
	if c.Plat == model.PlatForAutoPlayAndroid || c.Plat == model.PlatForAutoPlayInlineAndroid || c.Plat == model.PlatForAutoPlayIOS || c.Plat == model.PlafForAutoPlayInlineIOS ||
		strings.Contains(c.UserAgent, "(inline_play_begin)") { // plat的逻辑更换为UA中添加(inline_play_begin)
		log.Warn("no count! hit autoplay plat(%d) aid(%d)", c.Plat, c.AID)
		rtype = model.LogTypeForInlineBegin // 2
		if c.Buvid != "" {
			if err = s.setRealDid(ctx, c.Buvid, c.AID, c.Did); err != nil {
				log.Error("s.setRealDid(%s, %s) error(%v)", c.Buvid, c.Did, err)
				return
			}
		}
		return
	}
	// 自动播放的恶心逻辑 结束
	if c.MID > 0 {
		if _, ok := s.forbidMids[c.MID]; ok {
			log.Warn("mid(%d) forbidden", c.MID)
			return
		}
	}
	if f, ok = s.forbids[c.AID][c.Plat]; ok {
		if f.Lv == -2 && (strings.HasSuffix(c.UserAgent, "(no_accesskey)") || c.MID == 0) {
			log.Warn("no count! hit no_accesskey! agent(%s) aid(%d) mid(%d) plat(%d)", c.UserAgent, c.AID, c.MID, c.Plat)
			return
		} else if f.Lv == -1 && c.MID == 0 {
			// 游客不计算点击数
			log.Warn("no count! hit forbid_lv(%d) mid(%d)", f.Lv, c.MID)
			return
		} else if f.Lv >= 0 && (c.Lv <= f.Lv || c.MID == 0) {
			// 游客和低于锁定等级的不计算点击数
			return
		}
	}
	if c.EpID > 0 {
		if err = s.checkEpAvRelation(ctx, c.AID, c.EpID, c.SeasonType); err != nil {
			log.Error("s.getOid(%d, %d, %d) error(%v)", c.AID, c.EpID, c.SeasonType)
			return
		}
	}
	if !s.canCount(ctx, c.AID, c.EpID, c.IP, c.STime, c.Did) {
		log.Warn("same ip(%s) and av(%d) replay", c.IP, c.AID)
		return
	}
	duration = s.ArcDuration(ctx, c.AID)
	if s.isReplay(ctx, c.MID, c.AID, c.Did, duration) {
		log.Warn("mid(%d) bvid(%s) aid(%d) epid(%d) isPGC(%v) gapTime(%d) is replay", c.MID, c.Did, c.AID, c.EpID, c.EpID > 0, duration)
		return
	}
	rtype = model.LogTypeForTurly // 1
	return
}

func (s *Service) checkEpAvRelation(ctx context.Context, aid, epid int64, seasonType int) (err error) {
	s.etamMutex.RLock()
	_, ok := s.eTam[epid]
	s.etamMutex.RUnlock()
	if ok {
		return
	}
	var isLegal bool
	if isLegal, err = s.db.IsLegal(ctx, aid, epid, seasonType); err != nil {
		// 接口请求失败时按ugc维度走
		err = nil
	} else if !isLegal {
		// epid not exist
		log.Error("aid(%d) epid(%d) type(%d) not exist", aid, epid, seasonType)
		return
	}
	s.etamMutex.Lock()
	s.eTam[epid] = aid
	s.etamMutex.Unlock()
	return
}

// 播放计数主方法
func (s *Service) countClick(ctx context.Context, msg *model.ClickMsg, i int64) (err error) {
	var (
		ci  *model.ClickInfo
		ok  bool
		now = time.Now().Unix()
	)
	if msg == nil {
		log.Info("svr close s.aidMap length is %d", len(s.aidMap[i]))
		for _, ci = range s.aidMap[i] {
			// 反正这些视频点击数很多,不差这些,照顾小透明
			if ci.GetSum() > 100 {
				continue
			}
			s.upClick(ctx, ci)
		}
		return
	}
	idx := msg.AID % s.c.ChanNum
	if atomic.LoadInt64(&s.lockedMap[idx]) == _locked {
		log.Info("locking aidMap[%d] current length(%d)", idx, len(s.aidMap[idx]))
		for _, ci := range s.aidMap[idx] {
			if ci.GetSum() > 100 {
				continue
			}
			s.upClick(ctx, ci)
			delete(s.aidMap[idx], ci.Aid)
			time.Sleep(1 * time.Millisecond)
		}
		atomic.StoreInt64(&s.lockedMap[idx], _unLock)
		log.Info("unlocked aidMap[%d] current length(%d)", idx, len(s.aidMap[idx]))
	}
	if ci, ok = s.aidMap[idx][msg.AID]; !ok {
		if ci, err = s.db.Click(ctx, msg.AID); err != nil {
			log.Error("s.db.Click(%d) error(%v)", msg.AID, err)
			return
		}
		if ci == nil {
			if _, err = s.db.AddClick(ctx, msg.AID, 0, 0, 0, 0, 0, 0); err != nil {
				log.Error("s.db.AddClick(%d) error(%v)", msg.AID, err)
				return
			}
			ci = &model.ClickInfo{Aid: msg.AID}
		}
		ci.Ready(now)
		s.aidMap[idx][msg.AID] = ci
	}
	switch msg.Plat {
	case 0:
		ci.Web++
	case 1:
		ci.H5++
	case 2:
		ci.Outer++
	case 3:
		ci.Ios++
	case 4:
		ci.Android++
	case 5:
		ci.AndroidTV++
	}
	if ci.Sum == 0 || now-ci.LastChangeTime > s.c.LastChangeTime {
		if err = s.upClick(ctx, ci); err != nil {
			log.Error("s.upClick(%v) error(%v)", ci, err)
			return
		}
		log.Info("truly add click message(%+v)", msg)
		if now-ci.LastChangeTime > s.c.ReleaseTime || ci.Aid == s.c.BnjMainAid || ci.NeedRelease() {
			delete(s.aidMap[idx], msg.AID)
			return
		}
		ci.Ready(now)
	}
	return
}

func (s *Service) upClick(c context.Context, ci *model.ClickInfo) (err error) {
	if _, err = s.db.UpClick(c, ci); err != nil {
		log.Error("s.db.UpClick(%+v) error(%v)", ci, err)
		return
	}
	s.busChan <- &model.StatMsg{AID: ci.Aid, Click: int(ci.Sum + ci.GetSum())}
	// 拜年祭需求 单品视频的播放数算进主视频
	if _, ok := s.bnjListAidMap[ci.Aid]; ok {
		bnj := &model.ClickInfo{
			Aid:       s.c.BnjMainAid,
			Web:       ci.Web,
			H5:        ci.H5,
			Outer:     ci.Outer,
			Ios:       ci.Ios,
			Android:   ci.Android,
			AndroidTV: ci.AndroidTV,
		}
		if _, err = s.db.UpClick(c, bnj); err != nil {
			log.Error("s.db.UpClick(%d) error(%v)", s.c.BnjMainAid, err)
			return
		}
		log.Info("bnjaid(%d) Forced to increase click(%v) by relateAid(%d)", s.c.BnjMainAid, bnj, ci.Aid)
	}
	// 拜年祭需求 单品视频的播放数算进主视频
	return
}

// SetSpecial http set special aid click
func (s *Service) SetSpecial(c context.Context, aid, num int64, tp string) (err error) {
	_, err = s.db.UpSpecial(c, aid, tp, num)
	return
}
