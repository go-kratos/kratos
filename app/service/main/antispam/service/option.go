package service

import (
	"time"

	"go-common/app/service/main/antispam/conf"
)

// NewOption .
func NewOption(config *conf.Config) *Option {
	opt := &Option{
		MaxSenderNum:           config.ServiceOption.MaxSenderNum,
		MinKeywordLen:          config.ServiceOption.MinKeywordLen,
		MaxExportRows:          config.ServiceOption.MaxExportRows,
		MaxRegexpCountsPerArea: config.ServiceOption.MaxRegexpCountsPerArea,
		MaxSpawnGoroutines:     config.ServiceOption.MaxSpawnGoroutines,

		DefaultChanSize:   config.ServiceOption.DefaultChanSize,
		AsyncTaskChanSize: config.ServiceOption.AsyncTaskChanSize,

		DefaultExpireSec:       config.ServiceOption.DefaultExpireSec,
		RuleDefaultExpireSec:   config.ServiceOption.RuleDefaultExpireSec,
		RegexpDefaultExpireSec: config.ServiceOption.RegexpDefaultExpireSec,
	}
	opt.Scheduler = &SchedulerOption{
		GcInterval:        time.Duration(config.ServiceOption.GcOpt.IntervalSec),
		GcMaxRowsPerQuery: config.ServiceOption.GcOpt.MaxRowsPerQuery,

		RefreshTrieIntervalSec:    time.Duration(config.ServiceOption.RefreshTrieIntervalSec),
		RefreshRulesIntervalSec:   time.Duration(config.ServiceOption.RefreshRulesIntervalSec),
		RefreshRegexpsIntervalSec: time.Duration(config.ServiceOption.RefreshRegexpsIntervalSec),

		BuildTrieIntervalMinute:  time.Duration(config.ServiceOption.BuildTrieIntervalMinute),
		BuildTrieMaxRowsPerQuery: config.ServiceOption.BuildTrieMaxRowsPerQuery,
	}
	if opt.AsyncTaskChanSize == 0 {
		opt.AsyncTaskChanSize = 500
	}
	return opt
}

// Option .
type Option struct {

	// MinKeywordLen specify the minimum length
	// a keyword should satify
	MinKeywordLen int

	// MaxSenderNum limit the length of
	// keyword's sender list
	MaxSenderNum int64
	// MaxExportRows specify the max rows
	// when export keywords as excel
	MaxExportRows int64
	// MaxRegexpCounts specify the max counts
	// of regexps inside the extract pipeline
	MaxRegexpCountsPerArea int64

	DefaultExpireSec       int64
	RuleDefaultExpireSec   int64
	RegexpDefaultExpireSec int64

	DefaultChanSize    int64
	MaxSpawnGoroutines int64
	AsyncTaskChanSize  int64

	Scheduler *SchedulerOption
}

// SchedulerOption .
type SchedulerOption struct {
	BuildTrieIntervalMinute  time.Duration
	BuildTrieMaxRowsPerQuery int64

	RefreshTrieIntervalSec     time.Duration
	RefreshTrieMaxRowsPerQuery int64

	RefreshRulesIntervalSec   time.Duration
	RefreshRegexpsIntervalSec time.Duration

	// GcInterval specify how often to
	// expire the useless keywords
	GcInterval        time.Duration
	GcMaxRowsPerQuery int64
}
