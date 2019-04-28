package model

// IndexFilter for dm_index_filter
type IndexFilter struct {
	ID       int64  `json:"id"`
	CID      int64  `json:"cid"`
	MID      int64  `json:"mid"`
	Filter   string `json:"filter"`
	Activate int8   `json:"activate"`
	Regex    int8   `json:"type"`
	Ctime    int64  `json:"ctime"`
}

// UserFilterList for member filter list
// to show global filter or video filter
type UserFilterList struct {
	Top     int8           `json:"accept_top"`
	Bottom  int8           `json:"accept_bottom"`
	Reverse int8           `json:"accept_reverse"`
	Filter  []*IndexFilter `json:"filter_list"`
}
