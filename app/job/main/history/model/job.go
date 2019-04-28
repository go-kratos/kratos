package model

// ArgMid arg.
type ArgMid struct {
	Mid int64
	Ts  int64
}

// Merge merge.
type Merge struct {
	Mid int64 `json:"mid"`
	Now int64 `json:"now"`
}
