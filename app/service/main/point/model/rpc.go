package model

// ArgRPCMid def.
type ArgRPCMid struct {
	Mid int64
}

// PointHistoryResp point history resp.
type PointHistoryResp struct {
	Phs   []*OldPointHistory
	Total int
}

//ArgRPCPointHistory def .
type ArgRPCPointHistory struct {
	Mid int64
	PS  int `form:"ps"`
	PN  int `form:"pn"`
}
