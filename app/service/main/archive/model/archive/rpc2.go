package archive

// ArgAid2 ArgAid2
type ArgAid2 struct {
	Aid    int64
	RealIP string
}

// ArgCid2 ArgCid2
type ArgCid2 struct {
	Aid    int64
	Cid    int64
	RealIP string
}

// ArgVideo2 ArgVideo2
type ArgVideo2 struct {
	Aid, Cid int64
	RealIP   string
}

// ArgAids2 ArgAids2
type ArgAids2 struct {
	Aids   []int64
	RealIP string
}

// ArgPage2 ArgPage2
type ArgPage2 struct {
	Aid       int64
	Mid       int64
	AccessKey string
	RealIP    string
}

// ArgVideoshot2 ArgVideoshot2
type ArgVideoshot2 struct {
	Cid    int64
	Count  int
	RealIP string
}

// ArgStat2 ArgStat2
type ArgStat2 struct {
	Aid    int64
	Field  int
	Value  int
	RealIP string
}

// ArgAidMid2 ArgAidMid2
type ArgAidMid2 struct {
	Aid    int64
	Mid    int64
	RealIP string
}

// ArgUpArcs2 ArgUpArcs2
type ArgUpArcs2 struct {
	Mid    int64
	Pn, Ps int
	RealIP string
}

// ArgUpCount2 ArgUpCount2
type ArgUpCount2 struct {
	Mid int64
}

// ArgUpsArcs2 ArgUpsArcs2
type ArgUpsArcs2 struct {
	Mids   []int64
	Pn, Ps int
	RealIP string
}

// ArgMovie2 ArgMovie2
type ArgMovie2 struct {
	MovieId int64
	RealIP  string
}

// ArgMid2 ArgMid2
type ArgMid2 struct {
	Mid    int64
	RealIP string
}

// ArgRank2 ArgRank2
type ArgRank2 struct {
	Rid    int16
	Type   int8
	Pn, Ps int
	RealIP string
}

// ArgRanks2 ArgRanks2
type ArgRanks2 struct {
	Rids   []int16
	Type   int8
	Pn, Ps int
	RealIP string
}

// ArgRankTop2 ArgRankTop2
type ArgRankTop2 struct {
	ReID   int16
	Pn, Ps int
}

// ArgRankAll2 ArgRankAll2
type ArgRankAll2 struct {
	Pn, Ps int
}

// ArgRankTopsCount2 ArgRankTopsCount2
type ArgRankTopsCount2 struct {
	ReIDs []int16
}

// ArgCIDs2 ArgCIDs2
type ArgCIDs2 struct {
	Cids []int64
}

// ArgEp2 ArgEp2
type ArgEp2 struct {
	EpIDs []int64
	Tp    int8
}

// const action type
const (
	CacheAdd    = "add"
	CacheUpdate = "update"
	CacheDelete = "delete"
)

// ArgCache2 ArgCache2
type ArgCache2 struct {
	Aid    int64
	Tp     string
	OldMid int64
}

// ArgFieldCache2 ArgFieldCache2
type ArgFieldCache2 struct {
	Aid       int64
	TypeID    int16
	OldTypeID int16
}
