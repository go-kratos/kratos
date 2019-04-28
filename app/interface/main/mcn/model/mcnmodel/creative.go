package mcnmodel

// CreativeCommonReq common creative
type CreativeCommonReq struct {
	UpMid  int64 `form:"mid" validate:"required"`
	McnMid int64
}

// AidCommonReq common aid
type AidCommonReq struct {
	CreativeCommonReq

	Aid int64 `form:"aid" validate:"required"`
}

// DataCommonReq request params
type DataCommonReq struct {
	CreativeCommonReq

	Type int8 `form:"type" validate:"required"`
}

// ArchivesReq request params
type ArchivesReq struct {
	CreativeCommonReq

	Pn      int64  `form:"pn"`
	Ps      int64  `form:"ps"`
	Order   string `form:"order"`
	TID     int64  `form:"tid"`
	KeyWord string `form:"keyword"`
	Status  string `form:"status" validate:"required"`
	Coop    int16  `form:"coop"`
}

// ArchiveHistoryListReq request params
type ArchiveHistoryListReq = AidCommonReq

// ArchiveVideosReq request params
type ArchiveVideosReq = AidCommonReq

// DataArchiveReq request params
type DataArchiveReq = AidCommonReq

// DataVideoQuitReq request params
type DataVideoQuitReq struct {
	CreativeCommonReq

	Cid int64 `form:"cid" validate:"required"`
}

// DanmuDistriReq .
type DanmuDistriReq struct {
	CreativeCommonReq

	Aid int64 `form:"aid" validate:"required"`
	Cid int64 `form:"cid" validate:"required"`
}

// DataBaseReq request params
type DataBaseReq = CreativeCommonReq

// DataTrendReq request params
type DataTrendReq = CreativeCommonReq

// DataActionReq request params
type DataActionReq = CreativeCommonReq

// DataFanReq request params
type DataFanReq = CreativeCommonReq

// DataPandectReq request params
type DataPandectReq = DataCommonReq

// DataSurveyReq request params
type DataSurveyReq = DataCommonReq

// DataPlaySourceReq request params
type DataPlaySourceReq = CreativeCommonReq

// DataPlayAnalysisReq request params
type DataPlayAnalysisReq struct {
	CreativeCommonReq

	Copyright int8 `form:"copyright"`
}

// DataArticleRankReq request params
type DataArticleRankReq = DataCommonReq
