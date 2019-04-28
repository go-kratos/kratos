package archive

// ArgAid ArgAid
type ArgAid struct {
	Aid    int64
	RealIP string
}

// ArgCid ArgCid
type ArgCid struct {
	Cid    int64
	RealIP string
}

// ArgAids ArgAids
type ArgAids struct {
	Aids   []int64
	RealIP string
}

// ArgPage ArgPage
type ArgPage struct {
	Aid       int64
	Mid       int64
	AccessKey string
	RealIP    string
}

// ArgVideoshot ArgVideoshot
type ArgVideoshot struct {
	Cid    int64
	Count  int
	RealIP string
}

// ArgStat ArgStat
type ArgStat struct {
	Aid    int64
	Field  int
	Value  int
	RealIP string
}

// ArgTag ArgTag
type ArgTag struct {
	Aid    int64
	Tag    string
	RealIP string
}

// ArgPlayer ArgPlayer
type ArgPlayer struct {
	Aids     []int64
	Qn       int
	Platform string
	RealIP   string
	Fnval    int
	Fnver    int
	Build    int
	// 非必传
	Session   string
	ForceHost int
}
