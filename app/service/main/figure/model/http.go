package model

type ParamBatchInfo struct {
	MIDs []int64 `form:"mids,split"`
}
