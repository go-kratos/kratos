package bnj

// ResetMsg .
type ResetMsg struct {
	Mid int64 `json:"mid"`
	Ts  int64 `json:"ts"`
}

// Push .
type Push struct {
	Second        int64  `json:"second"`
	Name          string `json:"name"`
	TimelinePic   string `json:"timeline_pic,omitempty"`
	H5TimelinePic string `json:"h5_timeline_pic,omitempty"`
}
