package datamodel

type commonArg struct {
	Mid int64 `form:"mid"`
}

//GetFansSummaryArg arg to get fans
type GetFansSummaryArg struct {
	Mid int64 `form:"mid"`
}

//FansSummaryResult result for fans result
type FansSummaryResult struct {
	FanSummary FanSummaryData `json:"fan_summary"`
}

const (
	//DataType30Day 30 day
	DataType30Day = 1
	//DataTypeMonth by month
	DataTypeMonth = 2
)

//GetRelationFansHistoryArg arg
type GetRelationFansHistoryArg struct {
	Mid      int64 `form:"mid"`
	DataType int   `form:"data_type"`
}

//GetRelationFansHistoryResult relation fan history
type GetRelationFansHistoryResult struct {
	RelationFanHistoryData
}

// GetRelationFansMonthArg arg
type GetRelationFansMonthArg = GetFansSummaryArg

// GetRelationFansMonthResult relation fan history
type GetRelationFansMonthResult struct {
	RelationFanHistoryData
}

//GetUpArchiveInfoArg arg
type GetUpArchiveInfoArg struct {
	Mids     string `form:"mids" validate:"required"`
	DataType int    `form:"data_type"`
}

//GetUpArchiveInfoResult result, key = mid, value = data
type GetUpArchiveInfoResult = map[int64]*UpArchiveData

//GetUpArchiveTagInfoArg tag info
type GetUpArchiveTagInfoArg = commonArg

//GetUpArchiveTagInfoResult resutl
type GetUpArchiveTagInfoResult = []*ViewerTagData

//GetUpArchiveTypeInfoArg arg to get type
type GetUpArchiveTypeInfoArg = commonArg

//GetUpViewInfoArg get up view info data
type GetUpViewInfoArg = commonArg
