package model

// ArchiveParam .
type ArchiveParam struct {
	Aids []int64 `json:"aids" form:"aids,split" validate:"min=1,max=30,dive,gt=0"`
}
