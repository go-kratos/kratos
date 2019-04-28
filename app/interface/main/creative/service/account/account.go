package account

import (
	"context"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"time"

	accmdl "go-common/app/interface/main/creative/model/account"
	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	xtime "go-common/library/time"
	"sync"
)

const creatorMID = 37090048

// MyInfo get use level.
func (s *Service) MyInfo(c context.Context, mid int64, ip string, now time.Time) (m *accmdl.MyInfo, err error) {
	var (
		idCkCode = 0
		res      *account.Profile
	)
	// init myinfo
	m = &accmdl.MyInfo{
		Mid:        mid,
		Level:      1,
		Activated:  true,
		Deftime:    xtime.Time(now.Add(time.Duration(14400) * time.Second).Unix()),
		DeftimeEnd: xtime.Time(now.Add(time.Duration(15*24) * time.Hour).Unix()),
		DeftimeMsg: ecode.String(ecode.VideoupDelayTimeErr.Error()).Message(),
		IdentifyInfo: &accmdl.IdentifyInfo{
			Code: idCkCode,
			Msg:  accmdl.IdentifyEnum[idCkCode],
		},
		UploadSize: map[string]bool{
			"4-8":  false,
			"8-16": false,
		},
		DmSubtitle: true,
	}
	if res, err = s.acc.Profile(c, mid, ip); err != nil {
		log.Error("s.acc.Profile (%d) error(%v)", mid, err)
		return
	}
	if res == nil {
		err = ecode.NothingFound
		return
	}
	if res.Silence == 1 {
		m.Banned = true
	}
	if res.EmailStatus == 1 || res.TelStatus == 1 {
		m.Activated = true
	}
	m.Name = res.Name
	m.Face = res.Face
	m.Level = int(res.Level)
	if _, ok := s.exemptIDCheckUps[mid]; ok {
		log.Info("s.exemptIDCheckUps (%d) info(%v)", mid, s.exemptIDCheckUps)
		idCkCode = 0
	} else {
		idCkCode, _ = s.acc.IdentifyInfo(c, mid, 0, ip)
	}
	m.IdentifyInfo.Code = idCkCode
	m.IdentifyInfo.Msg = accmdl.IdentifyEnum[idCkCode] // NOTE: exist check
	m = s.ignoreLevelUpAndActivated(m, mid, s.exemptZeroLevelAndAnswerUps)
	if _, ok := s.uploadTopSizeUps[mid]; ok {
		m.UploadSize["4-8"] = true
		m.UploadSize["8-16"] = true
	}
	return
}

// type为12的白名单人员，免账号激活，免账号升级到1级
func (s *Service) ignoreLevelUpAndActivated(m *accmdl.MyInfo, mid int64, exemptZeroLevelAndAnswerUps map[int64]int64) *accmdl.MyInfo {
	if len(exemptZeroLevelAndAnswerUps) > 0 {
		_, ok := exemptZeroLevelAndAnswerUps[mid]
		if ok && (m.Level < 1 || !m.Activated) {
			m.Activated = true
			m.Level = 1
		}
	}
	return m
}

// UpInfo get video article pic blink up infos.
func (s *Service) UpInfo(c context.Context, mid int64, ip string) (v *accmdl.UpInfo, err error) {
	cache := true
	if v, err = s.acc.UpInfoCache(c, mid); err != nil {
		s.pCacheMiss.Incr("upinfo_cache")
		err = nil
		cache = false
	} else if v != nil {
		s.pCacheHit.Incr("upinfo_cache")
		return
	}
	s.pCacheMiss.Incr("upinfo_cache")
	var (
		g   = &errgroup.Group{}
		ctx = context.TODO()
	)
	v = &accmdl.UpInfo{
		Archive: accmdl.NotUp,
		Article: accmdl.NotUp,
		Pic:     accmdl.NotUp,
		Blink:   accmdl.NotUp,
	}
	g.Go(func() error {
		blikCnt, _ := s.acc.Blink(ctx, mid, ip)
		if blikCnt > 0 {
			v.Blink = accmdl.IsUp
		}
		return nil
	})
	g.Go(func() error {
		picCnt, _ := s.acc.Pic(c, mid, ip)
		if picCnt > 0 {
			v.Pic = accmdl.IsUp
		}
		return nil
	})
	g.Go(func() error {
		arcCnt, _ := s.archive.UpCount(c, mid)
		if arcCnt > 0 {
			v.Archive = accmdl.IsUp
		}
		return nil
	})
	g.Go(func() error {
		isAuthor, _ := s.article.IsAuthor(c, mid, ip)
		if isAuthor {
			v.Article = accmdl.IsUp
		}
		return nil
	})
	g.Wait()
	if cache {
		s.addCache(func() {
			s.acc.AddUpInfoCache(context.Background(), mid, v)
		})
	}
	return
}

// RecFollows 推荐关注
func (s *Service) RecFollows(c context.Context, mid int64) (fs []*accmdl.Friend, err error) {
	var (
		fsMids           = []int64{creatorMID}
		shouldFollowMids = []int64{}
		infos            = make(map[int64]*account.Info)
	)
	ip := metadata.String(c, metadata.RemoteIP)
	if shouldFollowMids, err = s.acc.ShouldFollow(c, mid, fsMids, ip); err != nil {
		log.Error("s.acc.ShouldFollow mid(%d)|ip(%s)|error(%v)", mid, ip, err)
		return
	}
	if len(shouldFollowMids) > 0 {
		if infos, err = s.acc.Infos(c, shouldFollowMids, ip); err != nil {
			log.Error("s.acc.Infos mid(%d)|ip(%s)|shouldFollowMids(%+v)|error(%v)", mid, ip, shouldFollowMids, err)
			return
		}
		for mid, info := range infos {
			f := &accmdl.Friend{
				Mid:  mid,
				Sign: info.Sign,
				Face: info.Face,
				Name: info.Name,
			}
			if mid == creatorMID {
				f.Comment = "关注创作中心，获取更多创作情报"
				f.ShouldFollow = 0
			}
			fs = append(fs, f)
		}
	}
	return
}

// Infos 获取多个UP主的信息
func (s *Service) Infos(c context.Context, mids []int64, ip string) (res map[int64]*account.Info, err error) {
	return s.acc.Infos(c, mids, ip)
}

// Relations 批量获取mid与其它用户的关系
func (s *Service) Relations(c context.Context, mid int64, fids []int64, ip string) (res map[int64]int, err error) {
	return s.acc.Relations(c, mid, fids, ip)
}

// FRelations 获取用户与mid的关系（Relations的反向）
func (s *Service) FRelations(c context.Context, mid int64, fids []int64, ip string) (res map[int64]int, err error) {
	var (
		g  errgroup.Group
		sm sync.RWMutex
	)
	res = make(map[int64]int)
	for _, v := range fids {
		g.Go(func() error {
			var r map[int64]int
			if r, err = s.acc.Relations(c, v, []int64{mid}, ip); err != nil {
				return err
			}
			sm.Lock()
			res[v] = r[mid]
			sm.Unlock()
			return nil
		})
	}
	if err = g.Wait(); err != nil {
		log.Error("s.FRelations(%d,%v) error(%v)", mid, fids, err)
	}
	return
}

// Cards 批量获取用户的Card
func (s *Service) Cards(c context.Context, mids []int64, ip string) (cards map[int64]*account.Card, err error) {
	return s.acc.Cards(c, mids, ip)
}
