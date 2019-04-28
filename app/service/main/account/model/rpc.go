package model

// ArgMid is.
type ArgMid struct {
	Mid int64
}

// ArgMids is.
type ArgMids struct {
	Mids []int64
}

// ArgNames is.
type ArgNames struct {
	Names []string
}

// ArgExp is.
type ArgExp struct {
	Mid      int64
	Exp      float64
	Operater string
	Operate  string
	Reason   string
	RealIP   string
}

// ArgMoral is.
type ArgMoral struct {
	Mid    int64
	Moral  float64
	Oper   string
	Reason string
	Remark string
	RealIP string
}

// ArgRelation is.
type ArgRelation struct {
	Mid, Owner int64
	RealIP     string
}

// ArgRelations is.
type ArgRelations struct {
	Mid    int64
	Owners []int64
	RealIP string
}

// ArgRichRelation is.
type ArgRichRelation struct {
	Owner  int64
	Mids   []int64
	RealIP string
}
