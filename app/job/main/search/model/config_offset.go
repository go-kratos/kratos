package model

import (
	"time"
)

// LoopOffset single table offset
type LoopOffset struct {
	IsLoop          bool
	OffsetID        int64
	OffsetTime      string
	TempOffsetID    int64
	TempOffsetTime  string
	RecoverID       int64
	RecoverTime     string
	TempRecoverID   int64
	TempRecoverTime string
	ReviewID        int64
	ReviewTime      int64
}

// SetLoop .
func (lo *LoopOffset) SetLoop(isLoop bool) {
	lo.IsLoop = isLoop
}

// SetReview .
func (lo *LoopOffset) SetReview(rid int64, rtime int64) {
	lo.ReviewID = rid
	lo.ReviewTime = rtime
}

// SetOffset .
func (lo *LoopOffset) SetOffset(id int64, t string) {
	if id != 0 {
		lo.OffsetID = id
	}
	if t != "" {
		lo.OffsetTime = t
		if !lo.IsLoop {
			if local, err := time.LoadLocation("Local"); err == nil {
				if t2, e := time.ParseInLocation("2006-01-02 15:04:05", t, local); e == nil && t2.Unix()-lo.ReviewTime > 0 {
					lo.OffsetTime = time.Unix(t2.Unix()-lo.ReviewTime, 0).Format("2006-01-02 15:04:05") //往前推ReviewTime
				}
			}
		}
	}
}

// SetTempOffset .
func (lo *LoopOffset) SetTempOffset(id int64, time string) {
	if id != 0 {
		lo.TempOffsetID = id
	}
	if time != "" {
		lo.TempOffsetTime = time
	}
}

// SetRecoverOffset .
func (lo *LoopOffset) SetRecoverOffset(recoverID int64, recoverTime string) {
	if recoverID >= 0 {
		lo.RecoverID = recoverID
	}
	if recoverTime != "" {
		lo.RecoverTime = recoverTime
	}
}

// SetRecoverTempOffset .
func (lo *LoopOffset) SetRecoverTempOffset(recoverID int64, recoverTime string) {
	if recoverID >= 0 {
		lo.TempRecoverID = recoverID
	}
	if recoverTime != "" {
		lo.TempRecoverTime = recoverTime
	}
}

// LoopOffsets more tables offset
type LoopOffsets map[int]*LoopOffset

// SetLoops .
func (los LoopOffsets) SetLoops(i int, isLoop bool) {
	if _, ok := los[i]; ok {
		los[i].IsLoop = isLoop
	}
}

// SetOffsets .
func (los LoopOffsets) SetOffsets(i int, id int64, time string) {
	if id != 0 {
		los[i].OffsetID = id
	}
	if time != "" {
		los[i].OffsetTime = time
	}
}

// SetTempOffsets .
func (los LoopOffsets) SetTempOffsets(i int, id int64, time string) {
	if id != 0 {
		los[i].TempOffsetID = id
	}
	if time != "" {
		los[i].TempOffsetTime = time
	}
}

// SetRecoverOffsets .
func (los LoopOffsets) SetRecoverOffsets(i int, recoverID int64, recoverTime string) {
	if recoverID >= 0 {
		los[i].RecoverID = recoverID
	}
	if recoverTime != "" {
		los[i].RecoverTime = recoverTime
	}
}

// SetRecoverTempOffsets .
func (los LoopOffsets) SetRecoverTempOffsets(i int, recoverID int64, recoverTime string) {
	if recoverID >= 0 {
		los[i].TempRecoverID = recoverID
	}
	if recoverTime != "" {
		los[i].TempRecoverTime = recoverTime
	}
}
