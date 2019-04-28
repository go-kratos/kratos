package archive

// Result is archive model.
type Result struct {
	Aid       int64  `json:"aid"`
	Mid       int64  `json:"mid"`
	TypeID    int16  `json:"tid"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Cover     string `json:"cover"`
	Tag       string `json:"tag"`
	Duration  int64  `json:"duration"`
	Copyright int8   `json:"copyright"`
	Desc      string `json:"desc"`
	Round     int8   `json:"round"`
	Forward   int64  `json:"forward"`
	Attribute int32  `json:"attribute"`
	HumanRank int    `json:"humanrank"`
	Access    int16  `json:"access"`
	State     int8   `json:"state"`
	Reason    string `json:"reject_reason"`
	PTime     string `json:"ptime"`
	CTime     string `json:"ctime"`
	MTime     string `json:"mtime"`
	Dynamic   string `json:"dynamic"`
}

// IsNormal check archive is open.
func (a *Result) IsNormal() bool {
	return a.State >= StateOpen || a.State == StateForbidFixed
}

// NotAllowUp check archive is or not allow update state.
func (a *Result) NotAllowUp() bool {
	return a.State == StateForbidUpDelete || a.State == StateForbidLater || a.State == StateForbidLock || a.State == StateForbidPolice
}

// IsForbid check archive state forbid by admin or delete.
func (a *Result) IsForbid() bool {
	return a.State == StateForbidUpDelete || a.State == StateForbidRecicle || a.State == StateForbidPolice || a.State == StateForbidLock || a.State == StateForbidLater || a.State == StateForbidXcodeFail
}

// AttrVal get attribute value.
func (a *Result) AttrVal(bit uint) int32 {
	return (a.Attribute >> bit) & int32(1)
}

// AttrSet set attribute value.
func (a *Result) AttrSet(v int32, bit uint) {
	a.Attribute = a.Attribute&(^(1 << bit)) | (v << bit)
}

// WithAttr set attribute value with a attr value.
func (a *Result) WithAttr(attr Attr) {
	a.Attribute = a.Attribute | int32(attr)
}
