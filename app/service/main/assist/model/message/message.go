package message

// Relation str map to databus message
type Relation struct {
	Action string `json:"action"`
	Table  string `json:"table"`
	New    struct {
		Attr int   `json:"attribute"`
		FID  int64 `json:"fid"`
		MID  int64 `json:"mid"`
	} `json:"new"`
	Old struct {
		Attr int   `json:"attribute"`
		FID  int64 `json:"fid"`
		MID  int64 `json:"mid"`
	}
}
