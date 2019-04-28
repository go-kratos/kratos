package service

import (
	"context"
	"io/ioutil"
	"time"

	"go-common/app/job/main/passport-game-data/conf"
	"go-common/app/job/main/passport-game-data/dao"
	"go-common/library/log"
)

const (
	_defaultDelayDuration = time.Minute * 0
	_defaultStepDuration  = time.Minute * 15
	_defaultLoopDuration  = time.Second * 3

	_defaultBatchSize           = 1000
	_defaultBatchMissRetryCount = 3

	_timeFormat = "2006-01-02 15:04:05"
)

var (
	_loc = time.Now().Location()
)

// Service service.
type Service struct {
	c *conf.Config
	d *dao.Dao

	// init cloud
	ic *initCloudConfig

	// c2l
	c2lC *compareConfig

	// l2c
	l2cC *compareConfig
}

type compareConfig struct {
	On             bool
	OffsetFilePath string
	UseOldOffset   bool

	End bool

	StartTime     time.Time
	EndTime       time.Time
	DelayDuration time.Duration
	StepDuration  time.Duration
	LoopDuration  time.Duration

	BatchSize           int
	BatchMissRetryCount int

	Debug bool
	Fix   bool

	// runtime
	st, ed     time.Time
	rangeCount int
	totalCount int
	diffCount  int

	sleeping        bool
	sleepingSeconds int64
	sleepFromTs     int64
}

func newCompareConfigFrom(c *conf.CompareConfig) (cc *compareConfig) {
	st, err := time.ParseInLocation(_timeFormat, c.StartTime, _loc)
	if err != nil {
		log.Error("failed to parse end time, time.ParseInLocation(%s, %s, %v), error(%v)", _timeFormat, c.StartTime, _loc, err)
		return
	}
	ed, err := time.ParseInLocation(_timeFormat, c.EndTime, _loc)
	if err != nil {
		log.Error("failed to parse end time, time.ParseInLocation(%s, %s, %v), error(%v)", _timeFormat, c.EndTime, _loc, err)
		return
	}

	cc = &compareConfig{
		On:    c.On,
		Debug: c.Debug,

		OffsetFilePath: c.OffsetFilePath,
		UseOldOffset:   c.UseOldOffset,
		End:            c.End,

		StartTime: st,
		EndTime:   ed,

		DelayDuration: time.Duration(c.DelayDuration),
		StepDuration:  time.Duration(c.StepDuration),
		LoopDuration:  time.Duration(c.LoopDuration),

		BatchSize:           c.BatchSize,
		BatchMissRetryCount: c.BatchMissRetryCount,

		Fix: c.Fix,
	}

	if cc.UseOldOffset {
		if oldOffset, err := parseOldOffset(cc.OffsetFilePath); err == nil {
			cc.StartTime = oldOffset
		}
	}

	cc.fix()
	return
}

func parseOldOffset(path string) (oldOffset time.Time, err error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("failed to read old offset, ioutil.ReadFile(%s) error(%v) skip", path, err)
		return
	}

	if oldOffset, err = time.ParseInLocation(_timeFormat, string(data), _loc); err != nil {
		log.Error("failed to parse offset, time.ParseInLocation(%s, %s, %v) error(%v)", _timeFormat, string(data), _loc, err)
	}
	return
}

func (cc *compareConfig) fix() {
	if int64(cc.DelayDuration) < 0 {
		cc.DelayDuration = _defaultDelayDuration
	}
	if int64(cc.StepDuration) < 0 {
		cc.StepDuration = _defaultStepDuration
	}
	if int64(cc.LoopDuration) < 0 {
		cc.LoopDuration = _defaultLoopDuration
	}

	if cc.BatchSize <= 0 {
		cc.BatchSize = _defaultBatchSize
	}
	if cc.BatchMissRetryCount < 0 {
		cc.BatchMissRetryCount = _defaultBatchMissRetryCount
	}
}

// New new a service instance.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c: c,
		d: dao.New(c),
	}
	if c.Compare.Cloud2Local.On {
		s.c2lC = newCompareConfigFrom(c.Compare.Cloud2Local)
		go s.cloud2localcompareproc()
	}
	if c.Compare.Local2Cloud.On {
		s.l2cC = newCompareConfigFrom(c.Compare.Local2Cloud)
		go s.local2cloudcompareproc()
	}
	return
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.d.Ping(c)
	return
}

// Close close service.
func (s *Service) Close() (err error) {
	s.d.Close()
	return
}
