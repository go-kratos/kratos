package model

// ArgRegionTotal arg region total
type ArgRegionTotal struct {
	RealIP string
}

// ArgRegion3 arg region
type ArgRegion3 struct {
	RegionID int32
	Pn       int
	Ps       int
	RealIP   string
}

// ArgRegionTag3 arg region tag
type ArgRegionTag3 struct {
	TagID    int64
	RegionID int32
	Pn       int
	Ps       int
	RealIP   string
}

// ArgRegions3 arg regions
type ArgRegions3 struct {
	RegionIDs []int32
	Count     int
	RealIP    string
}
