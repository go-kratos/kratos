package timemachine

import (
	"time"

	"go-common/app/interface/main/activity/conf"
	"go-common/library/cache/memcache"
	"go-common/library/sync/pipeline/fanout"

	"go-common/library/database/hbase.v2"
)

// Dao .
type Dao struct {
	c     *conf.Config
	hbase *hbase.Client
	mc    *memcache.Pool
	cache *fanout.Fanout
	//limiter     *rate.Limiter
	mcTmExpire int32
	//tmProcStart int64
	//tmProcStop  int64
}

// New .
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c:     c,
		hbase: hbase.NewClient(c.Hbase),
		mc:    memcache.NewPool(c.TimeMc.Timemachine),
		cache: fanout.New("timemachine", fanout.Worker(8), fanout.Buffer(10240)),
	}
	d.mcTmExpire = int32(time.Duration(c.TimeMc.TmExpire) / time.Second)
	//d.limiter = rate.NewLimiter(1000, 100)
	//go d.startTmproc(context.Background())
	return d
}

// StartTmproc start time machine proc
//func (d *Dao) startTmproc(c context.Context) {
//	if env.DeployEnv != env.DeployEnvPre {
//		return
//	}
//	for {
//		time.Sleep(time.Second)
//		if d.tmProcStart != 0 {
//			go func() {
//				// scan key
//				max := 10000000000
//				step := max / 10000
//				prefix := step - 1
//				for i := 0; i < max; i += step {
//					time.Sleep(10 * time.Millisecond)
//					startRow := fmt.Sprintf("%0*d", 10, i)
//					endRow := fmt.Sprintf("%0*d", 10, i+prefix)
//					if err := d.timemachineScan(c, startRow, endRow); err != nil {
//						log.Error("startTmproc timemachineScan startRow(%s) endRow(%s) error(%v)", startRow, endRow, err)
//						continue
//					}
//					log.Info("startTmproc finish startRow(%s) endRow(%s)", startRow, endRow)
//				}
//			}()
//			break
//		}
//	}
//}

// StartTmProc start time machine proc.
//func (d *Dao) StartTmProc(c context.Context) {
//	atomic.StoreInt64(&d.tmProcStart, 1)
//}

// StopTmproc stop time machine proc.
//func (d *Dao) StopTmproc(c context.Context) {
//	atomic.StoreInt64(&d.tmProcStop, 1)
//}
