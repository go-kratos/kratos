package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/identify/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

func (s *Service) queryCookieDeleted() {
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Duration(s.c.CheckConf.Ticker))
	for {
		now := time.Now()
		for i := 0; i < 2; i++ {
			wg.Add(1)
			go func(i int, now time.Time) {
				defer wg.Done()
				s.readCookieData(dateFormat(addMonth(now, -i)))
			}(i, now)
		}
		wg.Wait()
		<-ticker.C
	}
}

func (s *Service) readCookieData(moth string) {
	var start int64
	for {
		cookies, err := s.d.CookieDeleted(context.Background(), start, s.c.CheckConf.Count, moth)
		if err != nil {
			log.Error("fail to get CookieDeleted error(%+v)", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if len(cookies) == 0 {
			log.Info("check cookie(%s) finished!", moth)
			return
		}
		maxID := int64(0)
		for _, a := range cookies {
			s.cookieCh[a.ID%int64(s.c.CheckConf.ChanNum)] <- a
			if maxID < a.ID {
				maxID = a.ID
			}
		}
		start = maxID
	}
}

func (s *Service) queryTokenDeleted() {
	var wg sync.WaitGroup
	ticker := time.NewTicker(time.Duration(s.c.CheckConf.Ticker))
	for {
		now := time.Now()
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func(i int, now time.Time) {
				defer wg.Done()
				s.readTokenData(dateFormat(addMonth(now, -i)))
			}(i, now)
		}
		wg.Wait()
		<-ticker.C
	}
}

func (s *Service) readTokenData(moth string) {
	var start int64
	for {
		tokens, err := s.d.TokenDeleted(context.Background(), start, s.c.CheckConf.Count, moth)
		if err != nil {
			log.Error("fail to get TokenDeleted error(%+v)", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if len(tokens) == 0 {
			log.Info("check token(%s) finished!", moth)
			return
		}
		maxID := int64(0)
		for _, a := range tokens {
			s.tokenCh[a.ID%int64(s.c.CheckConf.ChanNum)] <- a
			if maxID < a.ID {
				maxID = a.ID
			}
		}
		start = maxID
	}
}

func (s *Service) checkCookie(c chan *model.AuthCookie) {
	count := 0
	for {
		cookie, ok := <-c
		if !ok {
			log.Error("cookieChan closed")
			return
		}
		count = count + 1
		if count%10000 == 0 {
			count = 0
			time.Sleep(100 * time.Millisecond)
		}
		// auth
		for {
			res, err := s.d.CookieCache(context.Background(), cookie.Session)
			if err != nil {
				log.Error("fail to get cookie(%+v) cache from auth , error(%+v)", cookie, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if res != nil {
				prom.BusinessErrCount.Incr("auth:cacheNotDeleted")
				log.Error("auth cache not deleted, session(%s) mid(%d)", cookie.Session, cookie.Mid)
				for {
					err = s.d.DelCookieCache(context.Background(), cookie.Session)
					if err == nil {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
			break
		}
		// identify
		s.checkIdentifyCache(cookie.Session, cookie.Mid)
	}
}

func (s *Service) checkToken(c chan *model.AuthToken) {
	count := 0
	for {
		token, ok := <-c
		if !ok {
			log.Error("tokenChan closed")
			return
		}
		count = count + 1
		if count%10000 == 0 {
			count = 0
			time.Sleep(100 * time.Millisecond)
		}
		// auth
		for {
			res, err := s.d.TokenCache(context.Background(), token.Token)
			if err != nil {
				log.Error("fail to get token(%+v) cache from auth , error(%+v)", token, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if res != nil {
				prom.BusinessErrCount.Incr("auth:cacheNotDeleted")
				log.Error("auth cache not deleted, token(%s) mid(%d)", token.Token, token.Mid)
				for {
					err = s.d.DelTokenCache(context.Background(), token.Token)
					if err == nil {
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
			}
			break
		}
		// identify
		s.checkIdentifyCache(token.Token, token.Mid)
	}
}

func (s *Service) checkIdentifyCache(k string, mid int64) {
	for name, p := range s.poolm {
		mcc, ok := s.c.Memcaches[name]
		if !ok || mcc == nil {
			return
		}
		key := mcc.Prefix + k
		for {
			conn := p.Get(context.Background())
			res, err := conn.Get(key)
			conn.Close()
			if err != nil {
				if err == memcache.ErrNotFound {
					break
				}
				log.Error("fail to get cache(%s) from identify , error(%+v)", key, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if res != nil {
				prom.BusinessErrCount.Incr("identify:cacheNotDeleted")
				log.Error("identify cache not deleted, key(%s) mid(%d) cache(%s)", key, mid, name)
				for {
					conn := p.Get(context.Background())
					err := conn.Delete(key)
					conn.Close()
					if err == nil || err == memcache.ErrNotFound {
						break
					}
					log.Error("dao.DelCache(%s) error(%+v)", key, err)
					time.Sleep(100 * time.Millisecond)
				}
			}
			break
		}
	}
}

func dateFormat(t time.Time) string {
	return t.Format("200601")
}

func addMonth(t time.Time, delta int) time.Time {
	if delta == 0 {
		return t
	}
	year, month, _ := t.Date()
	thisMonthFirstDay := time.Date(year, month, 1, 1, 1, 1, 1, t.Location())
	return thisMonthFirstDay.AddDate(0, delta, 0)
}
