package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/live/push-live/conf"
	"go-common/app/interface/live/push-live/dao"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

var (
	_limitDecreaseUUIDKey = "ld:%s" // 接口请求防重复key

	errLimitRequestRepeat = errors.New("limit decrease request repeat")
	errConvertMidString   = errors.New("convert mid string error")
	errConvertBusiness    = errors.New("convert business error")
)

// Service struct
type Service struct {
	c               *conf.Config
	dao             *dao.Dao
	liveStartSub    *databus.Databus
	liveCommonSub   *databus.Databus
	wg              sync.WaitGroup
	closeCh         chan bool
	pushTypes       []string
	intervalExpired int32
	mutex           sync.RWMutex
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:             c,
		dao:           dao.New(c),
		liveStartSub:  databus.New(c.LiveRoomSub),
		liveCommonSub: databus.New(c.LiveCommonSub),
		closeCh:       make(chan bool),
		pushTypes:     make([]string, 0, 4),
		mutex:         sync.RWMutex{},
	}
	s.wg.Add(1)
	go s.loadPushConfig()

	for i := 0; i < c.Push.ConsumerProcNum; i++ {
		s.wg.Add(1)
		go s.liveMessageConsumeproc()
	}

	return s
}

// loadPushConfig Load push config
func (s *Service) loadPushConfig() {
	var ctx = context.TODO()
	defer s.wg.Done()
	for {
		select {
		case _, ok := <-s.closeCh:
			if !ok {
				log.Info("[service.push|loadPushConfig] s.loadPushConfig is closed by closeCh")
				return
			}
		default:
		}
		// get push delay time
		interval, err := s.dao.GetPushInterval(ctx)
		if err != nil || interval < 0 {
			time.Sleep(time.Duration(time.Minute))
			continue
		}
		s.mutex.Lock()
		s.intervalExpired = interval
		s.mutex.Unlock()

		// get push options
		types, err := s.dao.GetPushConfig(ctx)
		if err != nil || len(types) == 0 {
			time.Sleep(time.Duration(time.Minute))
			continue
		}
		s.mutex.Lock()
		s.pushTypes = types
		s.mutex.Unlock()

		time.Sleep(time.Duration(time.Minute))
	}
}

// safeGetExpired
func (s *Service) safeGetExpired() int32 {
	s.mutex.RLock()
	expired := s.intervalExpired
	s.mutex.RUnlock()
	return expired
}

// LimitDecrease do mid string limit decrease
func (s *Service) LimitDecrease(ctx context.Context, business, targetID, uuid, midStr string) (err error) {
	var (
		f    *dao.Filter
		mids []int64
		b    int
	)

	// 判断请求是否重复
	err = s.limitDecreaseUnique(getUniqueKey(business, targetID, uuid))
	if err != nil {
		log.Error("[service.service|LimitDecrease] limitDecreaseUnique error(%v), uuid(%s), business(%s), targetID(%s), mid(%s)",
			err, uuid, business, targetID, midStr)
		return
	}

	b, err = strconv.Atoi(business)
	if err != nil {
		log.Error("[service.service|LimitDecrease] strconv business params error(%v)", err)
		err = errConvertBusiness
		return
	}
	filterConf := &dao.FilterConfig{
		Business:     b,
		DailyExpired: dailyExpired(time.Now())}

	// convert mid string to []int64
	mids, err = s.convertStrToInt64(midStr)
	if err != nil {
		log.Error("[service.service|LimitDecrease] convertStrToInt64 error(%v), business(%s), uuid(%s), mids(%s)",
			err, business, uuid, midStr)
		err = errConvertMidString
		return
	}

	// aysnc decrease limit
	f, err = s.dao.NewFilter(filterConf)
	if err != nil {
		log.Error("[service.service|LimitDecrease] new filter error(%v), business(%s), uuid(%s), mids(%v)",
			err, business, uuid, mids)
		return
	}
	go f.BatchDecreaseLimit(ctx, mids)
	return
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return nil
}

// Close Service
func (s *Service) Close() {
	close(s.closeCh)
	s.subClose()
	s.wg.Wait()
	s.dao.Close()
}

// subClose Close all sub channels
func (s *Service) subClose() {
	s.liveCommonSub.Close()
	s.liveStartSub.Close()
}

// dailyExpired
func dailyExpired(from time.Time) float64 {
	tm1 := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	return math.Floor(tm2.Sub(from).Seconds())
}

// convertStrToInt64 convert mid string to []int64 slice
func (s *Service) convertStrToInt64(m string) (mInts []int64, err error) {
	var (
		mSplit   []string
		errCount int
	)
	if m == "" {
		return
	}
	mSplit = strings.Split(m, ",")
	for _, mStr := range mSplit {
		mInt, convErr := strconv.Atoi(mStr)
		if convErr != nil {
			log.Error("[service.push|formatMidstr] convert mid(%v), error(%v)", mStr, convErr)
			errCount++
			continue
		}
		mInts = append(mInts, int64(mInt))
	}
	if errCount == len(mSplit) {
		err = fmt.Errorf("[service.push|formatMidstr] convert all mid failed, midstr(%s)", m)
	}
	return
}

// limitDecreaseUnique
func (s *Service) limitDecreaseUnique(key string) (err error) {
	var (
		conn  redis.Conn
		reply interface{}
	)
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	conn, err = redis.Dial(s.c.Redis.PushInterval.Proto, s.c.Redis.PushInterval.Addr, s.dao.RedisOption()...)
	if err != nil {
		log.Error("[service.service|limitDecreaseUnique] redis.Dial error(%v)", err)
		return
	}

	// redis cache exists judgement
	reply, err = conn.Do("SET", key, time.Now(), "EX", dailyExpired(time.Now()), "NX")
	if err != nil {
		return
	}
	// key exists
	if reply == nil {
		err = errLimitRequestRepeat
		return
	}
	return
}

// getUniqueKey get request unique key
func getUniqueKey(a, b, c string) string {
	return fmt.Sprintf(_limitDecreaseUUIDKey, a+b+c)
}
