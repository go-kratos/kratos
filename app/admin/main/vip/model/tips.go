package model

import "go-common/library/time"

// Tips def.
type Tips struct {
	ID        int64     `json:"id" form:"id"`
	Platform  int8      `json:"platform" form:"platform" validate:"required,min=1,gte=1"`
	Version   int64     `json:"version" form:"version"`
	Tip       string    `json:"tip" form:"tip" validate:"required"`
	Link      string    `json:"link" form:"link"`
	StartTime int64     `json:"start_time" form:"start_time" validate:"required,min=1,gte=1"`
	EndTime   int64     `json:"end_time" form:"end_time" validate:"required,min=1,gte=1"`
	Level     int8      `json:"level" form:"level" validate:"required,min=1,gte=1"`
	JudgeType int8      `json:"judge_type" form:"judge_type"`
	Operator  string    `json:"operator"`
	Deleted   int8      `json:"deleted"`
	Position  int8      `json:"position" form:"position" validate:"required,min=1,gte=1"`
	State     int8      `json:"state"`
	Ctime     time.Time `json:"ctime"`
	Mtime     time.Time `json:"mtime"`
}

// TipResp def.
type TipResp struct {
	ID          int64  `json:"id"`
	PlatformStr string `json:"platform_str"`
	JudgeBuild  int8   `json:"judge_build_str"`
	StateStr    string `json:"state_str"`
	State       int8   `json:"state"`
	Version     int64  `json:"version"`
	Tip         string `json:"tip"`
	Link        string `json:"link"`
	Operator    string `json:"operator"`
	Position    int8   `json:"position"`
	Ctime       int64  `json:"ctime"`
	Mtime       int64  `json:"mtime"`
}

// TipState tip state
func (t *Tips) TipState(stime, etime, now int64) {
	if stime > now {
		t.State = WaitShowTips
	} else if etime < now {
		t.State = ExpireTips
	} else {
		t.State = EffectiveTips
	}
}
