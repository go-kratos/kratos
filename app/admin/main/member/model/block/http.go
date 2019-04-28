package block

// ParamValidator .
type ParamValidator interface {
	Validate() bool
}

// ParamSearch .
type ParamSearch struct {
	MIDs []int64 `form:"mids,split"`
}

// Validate .
func (p *ParamSearch) Validate() bool {
	p.MIDs = intsSet(p.MIDs)
	if len(p.MIDs) == 0 || len(p.MIDs) > 200 {
		return false
	}
	return true
}

// ParamHistory .
type ParamHistory struct {
	MID  int64 `form:"mid"`
	Desc bool  `form:"desc"`
	PS   int   `form:"ps"`
	PN   int   `form:"pn"`
}

// Validate .
func (p *ParamHistory) Validate() bool {
	if p.MID <= 0 {
		return false
	}
	if p.PS <= 0 || p.PS > 100 {
		return false
	}
	if p.PN <= 0 {
		return false
	}
	return true
}

// ParamBatchBlock .
type ParamBatchBlock struct {
	MIDs      []int64        `form:"mids,split"`
	AdminID   int64          `form:"admin_id"`
	AdminName string         `form:"admin_name"`
	Source    BlockMgrSource `form:"source"`
	Area      BlockArea      `form:"area"`
	Reason    string         `form:"reason"`
	Comment   string         `form:"comment"`
	Action    BlockAction    `form:"action"`
	Duration  int64          `form:"duration"` // 单位：天
	Notify    bool           `form:"notify"`
}

// Validate .
func (p *ParamBatchBlock) Validate() bool {
	p.MIDs = intsSet(p.MIDs)
	if len(p.MIDs) == 0 || len(p.MIDs) > 200 {
		return false
	}
	if p.AdminID <= 0 {
		return false
	}
	if p.AdminName == "" {
		return false
	}
	if p.Source != BlockMgrSourceSys && p.Source != BlockMgrSourceCredit {
		return false
	}
	if !p.Area.Contain() {
		return false
	}
	if p.Comment == "" {
		return false
	}
	if p.Action != BlockActionForever && p.Action != BlockActionLimit {
		return false
	}
	if p.Action == BlockActionLimit {
		if p.Duration <= 0 {
			return false
		}
	}
	return true
}

// ParamBatchRemove .
type ParamBatchRemove struct {
	MIDs      []int64 `form:"mids,split"`
	AdminID   int64   `form:"admin_id"`
	AdminName string  `form:"admin_name"`
	Comment   string  `form:"comment"`
	Notify    bool    `form:"notify"`
}

// Validate .
func (p *ParamBatchRemove) Validate() bool {
	p.MIDs = intsSet(p.MIDs)
	if len(p.MIDs) == 0 || len(p.MIDs) > 200 {
		return false
	}
	if p.AdminID <= 0 {
		return false
	}
	if p.AdminName == "" {
		return false
	}
	if p.Comment == "" {
		return false
	}
	return true
}

func intsSet(ints []int64) (intSet []int64) {
	if len(ints) == 0 {
		return
	}
OUTER:
	for i := range ints {
		for ni := range intSet {
			if ints[i] == intSet[ni] {
				continue OUTER
			}
		}
		intSet = append(intSet, ints[i])
	}
	return
}
