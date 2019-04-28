package service

import (
	"context"
	"go-common/app/interface/live/push-live/dao"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/log"
	"strings"
	"sync"
	"time"
)

// MidFilter 收敛所有mid过滤逻辑入口
func (s *Service) midFilter(ml map[int64]bool, business int, task *model.ApPushTask) (midMap []int64) {
	var (
		mutex        sync.Mutex
		i            int
		midsList     [][]int64
		wg           sync.WaitGroup
		needDecrease = needDecrease(business)
		filterConf   = &dao.FilterConfig{
			Business:        business,
			IntervalExpired: s.safeGetExpired(),
			IntervalValue:   intervalValueByLinkValue(task.LinkValue),
			DailyExpired:    dailyExpired(time.Now()),
			Task:            task}
	)
	midMap = make([]int64, 0, len(ml))

	// split mids by limit
	mids := make([]int64, 0, s.c.Push.IntervalLimit)
	for mid := range ml {
		mids = append(mids, mid)
		i++
		if i == s.c.Push.IntervalLimit {
			i = 0
			midsList = append(midsList, mids)
			mids = make([]int64, 0, s.c.Push.IntervalLimit)
		}
	}
	if len(mids) > 0 {
		midsList = append(midsList, mids)
	}

	// filter goroutines
	for i := 0; i < len(midsList); i++ {
		wg.Add(1)
		go func(index int, mids []int64) {
			var (
				filteredMids []int64
				f            *dao.Filter
				err          error
				ctx          = context.TODO()
			)
			defer func() {
				log.Info("[service.mids|midFilter] BatchFilter before(%d), after(%d), task(%v), business(%d), err(%v)",
					len(mids), len(filteredMids), task, business, err)
				wg.Done()
			}()

			// new filter
			f, err = s.dao.NewFilter(filterConf)
			if err != nil {
				return
			}
			filteredMids = f.BatchFilter(ctx, s.dao.NewFilterChain(f), mids)
			if len(filteredMids) == 0 {
				f.Done()
				return
			}

			// after filter, do something
			if needDecrease {
				go f.BatchDecreaseLimit(ctx, filteredMids)
			}
			mutex.Lock()
			midMap = append(midMap, filteredMids...)
			mutex.Unlock()
		}(i, midsList[i])
	}
	wg.Wait()
	log.Info("[service.mids|midFilter] filtered task(%v), before(%d), after(%d), type(%d)",
		task, len(ml), len(midMap), business)
	return
}

// intervalValueByLinkValue get roomid by link value
func intervalValueByLinkValue(linkValue string) string {
	s := strings.Split(linkValue, ",")
	return s[0]
}

// needDecrease
func needDecrease(business int) bool {
	return business != model.ActivityBusiness
}
