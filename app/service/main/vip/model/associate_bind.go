package model

// BindInfo struct.
type BindInfo struct {
	Account *BindAccount `json:"account"`
	Outer   *BindOuter   `json:"outer"`
}

// BindAccount bind account.
type BindAccount struct {
	Mid  int64  `json:"mid"`
	Name string `json:"name"`
	Face string `json:"face"`
}

// BindOuter outer bind info.
type BindOuter struct {
	Tel       string `json:"tel"`
	BindState int32  `json:"bind_state"`
}

// ArgOpenBindByMid args.
type ArgOpenBindByMid struct {
	Mid       int64
	AppID     int64
	OutOpenID string
}
