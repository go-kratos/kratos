package income

// ArchiveChargeRatio av charge ratio
type ArchiveChargeRatio struct {
	ID         int64
	ArchiveID  int64
	Ratio      int64
	AdjustType int
	CType      int
}

// UpChargeRatio up charge ratio
type UpChargeRatio struct {
	ID         int64
	MID        int64
	Ratio      int64
	AdjustType int
	CType      int
}
