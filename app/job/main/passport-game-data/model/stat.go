package model

import (
	"fmt"
	"time"
)

type ProcStat struct {
	Cloud2Local *CompareProcStat `json:"cloud_2_local"`
	Local2Cloud *CompareProcStat `json:"local_2_cloud"`
}

// CompareProcStat status of compare proc.
type CompareProcStat struct {
	StartTime     string       `json:"start_time"`
	EndTime       string       `json:"end_time"`
	StepDuration  JsonDuration `json:"step_duration"`
	LoopDuration  JsonDuration `json:"loop_duration"`
	DelayDuration JsonDuration `json:"delay_duration"`

	BatchSize           int `json:"batch_size"`
	BatchMissRetryCount int `json:"batch_miss_retry_count"`

	Debug bool `json:"debug"`
	Fix   bool `json:"fix"`

	CurrentRangeStart        JSONTime `json:"current_range_start"`
	CurrentRangeEnd          JSONTime `json:"current_range_end"`
	CurrentRangeRecordsCount int      `json:"current_range_records_count"`
	TotalRangeRecordsCount   int      `json:"total_range_records_count"`
	DiffCount                int      `json:"diff_count"`

	Sleeping           bool   `json:"sleeping"`
	SleepFrom          string `json:"sleep_from,omitempty"`
	SleepSeconds       int64  `json:"sleep_seconds,omitempty"`
	SleepRemainSeconds int64  `json:"sleep_remain_seconds,omitempty"`
}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%s"`, time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(s), nil
}

type JsonDuration time.Duration

func (t JsonDuration) MarshalJSON() ([]byte, error) {
	s := fmt.Sprintf(`"%v"`, time.Duration(t))
	return []byte(s), nil
}
