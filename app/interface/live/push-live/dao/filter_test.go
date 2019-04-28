package dao

import (
	"context"
	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/live/push-live/model"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"math/rand"
	"testing"
)

func setRedis(conn redis.Conn, key string, value interface{}, ttl int32) error {
	_, err := conn.Do("SET", key, value, "EX", ttl)
	return err
}

// delRedis
func delRedis(conn redis.Conn, key string) error {
	_, err := conn.Do("DEL", key)
	return err
}

func filterClean(f *Filter, mids []int64) {
	for _, mid := range mids {
		keys := []string{
			GetDailyLimitKey(mid),
			GetIntervalKey(mid),
		}
		for _, key := range keys {
			delRedis(f.conn, key)
		}
	}
	f.Done()
	f = nil
}

func TestDao_needSmooth(t *testing.T) {
	initd()
	Convey("test business need smooth", t, func() {
		var (
			business int
			b        bool
		)
		Convey("test need smooth", func() {
			business = rand.Intn(100)
			b = d.needSmooth(business)
			So(b, ShouldEqual, true)

			business = 111
			b = d.needSmooth(business)
			So(b, ShouldEqual, true)
		})
		Convey("test no need smooth", func() {
			business = 101
			b = d.needSmooth(business)
			So(b, ShouldEqual, false)
		})
	})
}

func TestDao_needLimit(t *testing.T) {
	initd()
	Convey("test business need limit", t, func() {
		var (
			business int
			b        bool
		)
		Convey("test need limit", func() {
			business = rand.Intn(110)
			b = d.needSmooth(business)
			So(b, ShouldEqual, true)
		})
		Convey("test no need smooth", func() {
			business = 111
			b = d.needSmooth(business)
			So(b, ShouldEqual, true)
		})
	})
}

func TestDao_NewFilterChain(t *testing.T) {
	initd()
	Convey("test new filter chain", t, func() {
		var (
			f    *Filter
			conf *FilterConfig
			fc   FilterChain
			err  error
		)

		Convey("test business no need to filter", func() {
			conf = &FilterConfig{
				Business: 111,
			}
			f, err = d.NewFilter(conf)
			So(err, ShouldBeNil)
			fc = d.NewFilterChain(f)
			So(len(fc), ShouldEqual, 0)
		})

		Convey("test business with filter", func() {
			conf = &FilterConfig{
				Business: 101,
			}
			f, err = d.NewFilter(conf)
			So(err, ShouldBeNil)
			fc = d.NewFilterChain(f)
			So(len(fc), ShouldEqual, 1)

			// both business and IntervalExpired is necessary
			conf = &FilterConfig{
				Business:        rand.Intn(100),
				IntervalExpired: rand.Int31(),
			}
			f, err = d.NewFilter(conf)
			So(err, ShouldBeNil)
			fc = d.NewFilterChain(f)
			So(len(fc), ShouldEqual, 2)
		})
	})
}

func TestDao_dailyLimit(t *testing.T) {
	var (
		ctx = context.Background()
		mid int64
		key string
		f   *Filter
		err error
		b   bool
	)
	Convey("test daily limit filter", t, func() {
		mid = rand.Int63n(999999)
		key = GetDailyLimitKey(mid)
		filterConf := &FilterConfig{}
		f, err = d.NewFilter(filterConf)
		So(err, ShouldBeNil)
		log.Info("TestDao_dailyLimit mid(%d), key(%s), filter(%v)", mid, key, f)
		// del key first
		err = delRedis(f.conn, key)
		So(err, ShouldBeNil)
		// try get with nil return
		b, err = f.dailyLimitFilter(ctx, mid)
		So(b, ShouldEqual, false)
		So(err, ShouldBeNil)
		// then set a valid value
		setRedis(f.conn, key, rand.Intn(4)+1, 30)
		b, err = f.dailyLimitFilter(ctx, mid)
		So(b, ShouldEqual, false)
		So(err, ShouldBeNil)
		// then set value should be filtered
		setRedis(f.conn, key, -1, 30)
		b, err = f.dailyLimitFilter(ctx, mid)
		So(b, ShouldEqual, true)
		So(err, ShouldBeNil)
		// then test with conn error
		delRedis(f.conn, key)
		f.Done()
		b, err = f.dailyLimitFilter(ctx, mid)
		So(err, ShouldNotBeNil)
	})
}

func TestDao_intervalSmooth(t *testing.T) {
	var (
		ctx = context.Background()
		mid int64
		key string
		f   *Filter
		err error
		b   bool
	)
	Convey("test interval smooth filter", t, func() {
		mid = rand.Int63n(999999)
		key = GetIntervalKey(mid)
		// new filter
		task := &model.ApPushTask{
			LinkValue: "test",
		}
		fc := &FilterConfig{
			IntervalExpired: 300,
			Task:            task,
		}
		f, err = d.NewFilter(fc)
		So(err, ShouldBeNil)
		log.Info("TestDao_intervalSmooth mid(%d), key(%s), filter(%v)", mid, key, f)
		// del key first
		err = delRedis(f.conn, key)
		So(err, ShouldBeNil)
		// first setnx should success
		b, err = f.intervalSmoothFilter(ctx, mid)
		So(b, ShouldEqual, false)
		So(err, ShouldBeNil)
		// second setnx should fail
		b, err = f.intervalSmoothFilter(ctx, mid)
		So(b, ShouldEqual, true)
		So(err, ShouldBeNil)
		// test error
		delRedis(f.conn, key)
		f.Done()
		b, err = f.intervalSmoothFilter(ctx, mid)
		So(err, ShouldNotBeNil)
	})
}

func TestDao_BatchFilter(t *testing.T) {
	initd()
	Convey("test mids filter by different business", t, func() {
		var (
			ctx           = context.Background()
			business      int
			mids, resMids []int64
			conf          *FilterConfig
			fc            FilterChain
			f             *Filter
			err           error
		)

		Convey("test empty mids or filter chain", func() {
			// empty mids
			business = rand.Int()
			conf = &FilterConfig{Business: business}
			f, err = d.NewFilter(conf)
			So(err, ShouldBeNil)
			fc = d.NewFilterChain(f)
			resMids = f.BatchFilter(ctx, fc, mids)
			So(len(resMids), ShouldEqual, 0)

			// empty fc with business 111
			business = 111
			conf = &FilterConfig{Business: business}
			f, err = d.NewFilter(conf)
			So(err, ShouldBeNil)
			fc = d.NewFilterChain(f)
			So(len(fc), ShouldEqual, 0)
			resMids = f.BatchFilter(ctx, fc, mids)
			So(len(resMids), ShouldEqual, 0)

			// test business 101 limit filter case
			business = 101
			total := 10
			for i := 0; i < total; i++ {
				mids = append(mids, rand.Int63())
			}
			conf = &FilterConfig{Business: business}
			f, err = d.NewFilter(conf)
			So(err, ShouldBeNil)
			fc = d.NewFilterChain(f)
			So(len(fc), ShouldEqual, 1)
			resMids = f.BatchFilter(ctx, fc, mids)
			So(len(resMids), ShouldEqual, total)

			// clean test mids
			filterClean(f, mids)
		})

		Convey("test filter", func() {
			var b bool
			business = rand.Intn(99) + 1 // should through all filters
			total := 10
			shouldFilterMids := make([]int64, 0, total)
			for i := 0; i < total; i++ {
				mid := rand.Int63()
				if i%3 == 0 {
					shouldFilterMids = append(shouldFilterMids, mid)
				}
				mids = append(mids, mid)
			}
			task := &model.ApPushTask{
				LinkValue: "test",
			}
			conf = &FilterConfig{
				Business:        business,
				IntervalExpired: 300,
				DailyExpired:    300,
				Task:            task}
			f, err = d.NewFilter(conf)
			So(err, ShouldBeNil)
			fc = d.NewFilterChain(f)
			So(len(fc), ShouldEqual, 2)

			// init should filtered mids, half interval smooth and another daily limit
			for i, mid := range shouldFilterMids {
				if i%2 == 0 {
					b, err = f.intervalSmoothFilter(ctx, mid)
					So(b, ShouldEqual, false)
				} else {
					key := GetDailyLimitKey(mid)
					err = setRedis(f.conn, key, 0, int32(f.conf.DailyExpired))
				}
				So(err, ShouldBeNil)
			}

			// do filter
			resMids = f.BatchFilter(ctx, fc, mids)
			So(len(resMids), ShouldEqual, len(mids)-len(shouldFilterMids))

			// clean
			filterClean(f, mids)
		})
	})
}

func TestDao_BatchDecreaseLimit(t *testing.T) {
	initd()
	Convey("test batch decrease daily limit", t, func() {
		var (
			ctx               = context.Background()
			mids              []int64
			total, limitTotal int
			conf              *FilterConfig
			f                 *Filter
			err               error
		)
		total = 10
		for i := 0; i < total; i++ {
			mids = append(mids, rand.Int63())
		}
		log.Info("TestDao_BatchDecreaseLimit mids(%v)", mids)
		conf = &FilterConfig{
			DailyExpired: 300,
		}
		f, err = d.NewFilter(conf)
		So(err, ShouldBeNil)

		// do limit decrease
		limitTotal, err = f.BatchDecreaseLimit(ctx, mids)
		So(err, ShouldBeNil)
		So(limitTotal, ShouldEqual, total)

		// clean
		filterClean(f, mids)
	})
}
