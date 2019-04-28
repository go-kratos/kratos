package conf

// Rules 稿件导入规则
type Rules struct {
	Archive   *ArchiveRule
	Up        *UPRule
	Dimension *PageRule
}

// ArchiveRule .
type ArchiveRule struct {
	Titles    []string
	Contents  []string
	TID       []int32
	SubTID    []int32
	State     int
	NotAccess int
}

// UPRule .
type UPRule struct {
	UName []string
	MID   []int64
}

// PageRule .
type PageRule struct {
	MinX        int64
	MinY        int64
	MaxX        int64
	MaxY        int64
	MinYX       float32
	MaxYX       float32
	MaxDuration int64
	MinDuration int64
}
