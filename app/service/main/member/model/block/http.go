package block

// ParamValidator .
type ParamValidator interface {
	Validate() bool
}

// ParamInfo .
type ParamInfo struct {
	MID int64 `form:"mid"`
}

// Validate .
func (p *ParamInfo) Validate() bool {
	return p.MID > 0
}

// ParamBatchInfo .
type ParamBatchInfo struct {
	MIDs []int64 `form:"mids,split"`
}

// Validate .
func (p *ParamBatchInfo) Validate() bool {
	if len(p.MIDs) == 0 || len(p.MIDs) > 20 {
		return false
	}
	return true
}

// ParamBatchDetail .
type ParamBatchDetail struct {
	MIDs []int64 `form:"mids,split"`
}

// Validate .
func (p *ParamBatchDetail) Validate() bool {
	if len(p.MIDs) == 0 || len(p.MIDs) > 20 {
		return false
	}
	return true
}

// ParamBlock .
type ParamBlock struct {
	MID        int64       `form:"mid"`
	Source     BlockSource `form:"source"`
	Area       BlockArea   `form:"area"`
	Action     BlockAction `form:"action"`
	Duration   int64       `form:"duration"` // unix time
	StartTime  int64       `form:"start_time"`
	OperatorID int         `form:"op_id"`
	Operator   string      `form:"operator"`
	Reason     string      `form:"reason"`
	Comment    string      `form:"comment"`
	Notify     bool        `form:"notify"`
}

// Validate .
func (p *ParamBlock) Validate() bool {
	if p.MID <= 0 {
		return false
	}
	if !p.Source.Contain() {
		return false
	}
	if p.Action != BlockActionLimit && p.Action != BlockActionForever {
		return false
	}
	if p.StartTime <= 0 {
		return false
	}
	if p.Action == BlockActionLimit {
		if p.Duration <= 0 {
			return false
		}
	}
	return true
}

// ParamBatchBlock .
type ParamBatchBlock struct {
	MIDs       []int64     `form:"mids,split"`
	Source     BlockSource `form:"source"`
	Area       BlockArea   `form:"area"`
	Action     BlockAction `form:"action"`
	Duration   int64       `form:"duration"` // unix time
	StartTime  int64       `form:"start_time"`
	OperatorID int         `form:"op_id"`
	Operator   string      `form:"operator"`
	Reason     string      `form:"reason"`
	Comment    string      `form:"comment"`
	Notify     bool        `form:"notify"`
}

// Validate .
func (p *ParamBatchBlock) Validate() bool {
	if len(p.MIDs) == 0 || len(p.MIDs) > 20 {
		return false
	}
	if !p.Source.Contain() {
		return false
	}
	if p.Action != BlockActionLimit && p.Action != BlockActionForever {
		return false
	}
	if p.StartTime <= 0 {
		return false
	}
	if p.Action == BlockActionLimit {
		if p.Duration <= 0 {
			return false
		}
	}
	return true
}

// ParamRemove .
type ParamRemove struct {
	MID        int64       `form:"mid"`
	Source     BlockSource `form:"source"`
	OperatorID int         `form:"op_id"`
	Operator   string      `form:"operator"`
	Reason     string      `form:"reason"`
	Comment    string      `form:"comment"`
	Notify     bool        `form:"notify"`
}

// Validate .
func (p *ParamRemove) Validate() bool {
	if p.MID <= 0 {
		return false
	}
	if !p.Source.Contain() {
		return false
	}
	return true
}

// ParamBatchRemove .
type ParamBatchRemove struct {
	MIDs       []int64     `form:"mids,split"`
	Source     BlockSource `form:"source"`
	OperatorID int         `form:"op_id"`
	Operator   string      `form:"operator"`
	Reason     string      `form:"reason"`
	Comment    string      `form:"comment"`
	Notify     bool        `form:"notify"`
}

// Validate .
func (p *ParamBatchRemove) Validate() bool {
	if len(p.MIDs) == 0 || len(p.MIDs) > 20 {
		return false
	}
	if !p.Source.Contain() {
		return false
	}
	return true
}
