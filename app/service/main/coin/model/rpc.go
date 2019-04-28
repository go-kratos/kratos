package model

// ArchiveUserCoins resp user coins.
type ArchiveUserCoins struct {
	Multiply int64 `json:"multiply"`
}

// ArgCoinInfo arg coin info.
type ArgCoinInfo struct {
	Mid      int64
	Aid      int64
	AvType   int64
	Business string
	RealIP   string
}

// ArgAddCoin arg add coin.
type ArgAddCoin struct {
	Mid      int64
	UpMid    int64
	MaxCoin  int64
	Aid      int64
	AvType   int64
	Business string
	Multiply int64
	RealIP   string
	// archive only
	TypeID  int16
	PubTime int64
}

// ArgModifyCoin rpc arg ,modify user coins.
type ArgModifyCoin struct {
	Mid       int64   `json:"mid" form:"mid" validate:"required"`
	Count     float64 `json:"count" form:"count" validate:"required"`
	Reason    string  `json:"reason" form:"reason" validate:"required"`
	IP        string  `json:"ip"`
	Operator  string  `json:"operator" form:"operator"`
	CheckZero int8    `json:"check_zore" form:"check_zero"`
}

// ArgList rpc arg list.
type ArgList struct {
	Mid      int64
	TP       int64
	Business string
}

// ArgLog arg log
type ArgLog struct {
	Mid       int64
	Recent    bool
	Translate bool
}

// ArgAddUserCoinExp .
type ArgAddUserCoinExp struct {
	Mid          int64
	Business     int64
	BusinessName string
	Number       int64
	RealIP       string
}

// ArgMid .
type ArgMid struct {
	Mid    int64
	RealIP string
}
