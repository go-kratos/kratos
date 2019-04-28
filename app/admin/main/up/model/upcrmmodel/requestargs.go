package upcrmmodel

import (
	"go-common/app/admin/main/up/model/datamodel"
	"go-common/app/admin/main/up/util"
	"go-common/library/time"
	xtime "time"
)

const (
	//CompareTypeNothing 0
	CompareTypeNothing = 0
	//CompareType7day 1
	CompareType7day = 1
	//CompareType30day 2
	CompareType30day = 2
	//CompareTypeMonthFirstDay 3
	CompareTypeMonthFirstDay = 3
)

const (

	//AttrBitVideo video
	// see http://info.bilibili.co/pages/viewpage.action?pageId=9830931
	AttrBitVideo = 0
	//AttrBitAudio audio
	AttrBitAudio = 1
	//AttrBitArticle article
	AttrBitArticle = 2
	//AttrBitPhoto photo
	AttrBitPhoto = 3
	//AttrBitSign sign
	AttrBitSign = 4
	//AttrBitGrowup growup
	AttrBitGrowup = 5
	//AttrBitVerify verify
	AttrBitVerify = 6
)

var (
	//AttrGroup1 筛选用第一组attr
	AttrGroup1 = map[int]int{AttrBitVideo: 0, AttrBitAudio: 0, AttrBitArticle: 0, AttrBitPhoto: 0}
	//AttrGroup2 筛选用第二组attr， 两组之间的关系是与
	AttrGroup2 = map[int]int{AttrBitSign: 0, AttrBitGrowup: 0, AttrBitVerify: 0}
)

// ScoreQueryArgs ------------------------- requests ------------------------
type ScoreQueryArgs struct {
	ScoreType   int    `form:"score_type"`
	CompareType int    `form:"compare_type"`
	Export      string `form:"export"`
}

//ScoreQueryUpArgs arg
type ScoreQueryUpArgs struct {
	Mid  int64  `form:"mid" validate:"required"`
	Date string `form:"date"`
}

//ScoreQueryUpHistoryArgs  arg
type ScoreQueryUpHistoryArgs struct {
	Mid       int64  `form:"mid" validate:"required"`
	ScoreType int    `form:"score_type"`
	Day       int    `form:"day" default:"7"`
	Date      string `form:"date"`
}

//PlayQueryArgs arg
type PlayQueryArgs struct {
	Mid          int64 `form:"mid" validate:"required"`
	BusinessType int   `form:"business_type"`
}

//InfoQueryArgs arg
type InfoQueryArgs struct {
	Mid int64 `form:"mid" validate:"required"`
}

//CreditLogQueryArgs  arg
type CreditLogQueryArgs struct {
	Mid   int64 `form:"mid" validate:"required"`
	Limit int   `form:"limit"`
}

//UpRankQueryArgs  arg
type UpRankQueryArgs struct {
	Type int `form:"type" validate:"required"`
	Page int `form:"page"` // (从1开始)
	Size int `form:"size"` // 1 ~ 50
}

//InfoAccountInfoArgs  arg
type InfoAccountInfoArgs struct {
	Mids string `form:"mids" validate:"required"`
}

//InfoSearchArgs  arg
type InfoSearchArgs struct {
	AccountState   int    `json:"account_state"`
	Activity       int    `json:"activity"`
	Attrs          UpAttr `json:"attrs"`
	FirstDateBegin string `json:"first_date_begin"`
	FirstDateEnd   string `json:"first_date_end"`
	Mid            int64  `json:"mid"`
	Order          struct {
		Field string `json:"field"`
		Order string `json:"order"`
	}
	Page int `json:"page"`
	Size int `json:"size"`
}

//TestGetViewBaseArgs test arg
type TestGetViewBaseArgs struct {
	Mid int64 `form:"mid" validate:"required"`
}

// ------------------------- results ------------------------

//ScoreSection struct
type ScoreSection struct {
	Section int `json:"-"`
	Value   int
	Percent int
}

//ScoreQueryResult result
type ScoreQueryResult struct {
	CompareAxis []ScoreSection `json:"compareAxis"`
	XAxis       []string       `json:"xAxis"`
	YAxis       []ScoreSection `json:"yAxis"`
}

//NewEmptyScoreQueryResult make new result
func NewEmptyScoreQueryResult() ScoreQueryResult {
	return ScoreQueryResult{
		CompareAxis: []ScoreSection{},
		XAxis:       []string{},
		YAxis:       []ScoreSection{},
	}
}

//ScoreInfo struct
type ScoreInfo struct {
	Current     int `json:"current"`
	DiffLastDay int `json:"diff_last_day"`
}

//ScoreQueryUpResult  result
type ScoreQueryUpResult struct {
	PrScore      ScoreInfo
	QualityScore ScoreInfo
	CreditScore  ScoreInfo
	Date         time.Time
}

//ScoreHistoryInfo struct
type ScoreHistoryInfo struct {
	Type  int         `json:"type"`
	Score []int       `json:"score"`
	Date  []time.Time `json:"date"`
}

//ScoreQueryUpHistoryResult  result
type ScoreQueryUpHistoryResult struct {
	ScoreData []ScoreHistoryInfo `json:"score_data"`
}

//PlayInfo struct
type PlayInfo struct {
	Type                int   `json:"type"`
	PlayCountAccumulate int64 `json:"play_count_accumulate"`
	PlayCountAvg        int64 `json:"play_count_avg"`
	PlayCountAvg90Day   int64 `json:"play_count_avg_90day"`
}

//PlayQueryResult result
type PlayQueryResult struct {
	ArticleCount30Day      int        `json:"article_count_30day"`
	ArticleCountAccumulate int        `json:"article_count_accumulate"`
	BusinessData           []PlayInfo `json:"business_data"`
}

//CastUpPlayInfoToPlayInfo cast
func CastUpPlayInfoToPlayInfo(info UpPlayInfo) (r PlayInfo) {
	r.Type = int(info.BusinessType)
	r.PlayCountAccumulate = info.PlayCountAccumulate
	r.PlayCountAvg = info.PlayCountAccumulate / info.ArticleCount
	r.PlayCountAvg90Day = info.PlayCount90Day / info.ArticleCount
	return
}

//UpAttr struct
type UpAttr struct {
	AttrVerify  int `json:"attr_verify"`
	AttrVideo   int `json:"attr_video"`
	AttrAudio   int `json:"attr_audio"`
	AttrArticle int `json:"attr_article"`
	AttrPhoto   int `json:"attr_photo"`
	AttrSign    int `json:"attr_sign"`
	AttrGrowup  int `json:"attr_growup"`
}

//InfoQueryResult result
type InfoQueryResult struct {
	ID                     uint32     `json:"-"`
	Mid                    int64      `json:"mid"`
	Name                   string     `json:"name"`
	Sex                    int8       `json:"sex"`
	JoinTime               time.Time  `json:"join_time"`
	FirstUpTime            time.Time  `json:"first_up_time"`
	Level                  int16      `json:"level"`
	FansCount              int        `json:"fans_count"`
	AccountState           int8       `json:"account_state"`
	Activity               int        `json:"activity"`
	ArticleCount30day      int        `json:"article_count_30day"`
	ArticleCountAccumulate int        `json:"article_count_accumulate"`
	VerifyType             int8       `json:"verify_type"`
	BusinessType           int8       `json:"business_type"`
	CreditScore            int        `json:"credit_score"`
	PrScore                int        `json:"pr_score"`
	QualityScore           int        `json:"quality_score"`
	ActiveTid              int64      `json:"active_tid"`
	ActiveSubtid           int64      `json:"active_subtid"`
	Region                 string     `json:"region"`
	Province               string     `json:"province"`
	Age                    int        `json:"age"`
	Attr                   int        `json:"-"`
	Attrs                  UpAttr     `json:"attrs"`
	Birthday               xtime.Time `json:"-"`
}

//CopyFromBaseInfo copy
func (i *InfoQueryResult) CopyFromBaseInfo(info UpBaseInfo) {
	i.ID = info.ID
	i.Mid = info.Mid
	i.Name = info.Name
	i.Sex = info.Sex
	i.JoinTime = info.JoinTime
	i.FirstUpTime = info.FirstUpTime
	i.Level = info.Level
	i.FansCount = info.FansCount
	i.AccountState = info.AccountState
	i.ArticleCount30day = info.ArticleCount30day
	i.ArticleCountAccumulate = info.ArticleCountAccumulate
	i.VerifyType = info.VerifyType
	i.BusinessType = info.BusinessType
	i.CreditScore = info.CreditScore
	i.ActiveTid = info.ActiveTid
	i.Birthday = info.Birthday
	i.Region = info.ActiveCity
	i.Province = info.ActiveProvince
	i.Attr = info.Attr
	i.PrScore = info.PrScore
	i.QualityScore = info.QualityScore
	i.Activity = info.Activity
}

// CalculateAttr 根据attr来计算各个attr_xx的属性
func (i *InfoQueryResult) CalculateAttr() {
	// todo 计算attr属性
	if util.IsBitSet(i.Attr, AttrBitVideo) {
		i.Attrs.AttrVideo = 1
	}
	if util.IsBitSet(i.Attr, AttrBitAudio) {
		i.Attrs.AttrAudio = 1
	}
	if util.IsBitSet(i.Attr, AttrBitArticle) {
		i.Attrs.AttrArticle = 1
	}
	if util.IsBitSet(i.Attr, AttrBitPhoto) {
		i.Attrs.AttrPhoto = 1
	}
	if util.IsBitSet(i.Attr, AttrBitSign) {
		i.Attrs.AttrSign = 1
	}
	if util.IsBitSet(i.Attr, AttrBitGrowup) {
		i.Attrs.AttrGrowup = 1
	}
	if util.IsBitSet(i.Attr, AttrBitVerify) {
		i.Attrs.AttrVerify = 1
	}

	if !i.Birthday.IsZero() {
		i.Age = int(xtime.Since(i.Birthday).Hours() / float64(24*365))
		if i.Age < 0 {
			i.Age = 0
		}
	}
}

//CreditLogInfo struct
type CreditLogInfo struct {
	Time time.Time `json:"time"`
	Log  string    `json:"log"`
}

//CreditLogUpResult  result
type CreditLogUpResult struct {
	Logs []CreditLogInfo `json:"logs"`
}

//UpRankInfo struct
type UpRankInfo struct {
	InfoQueryResult
	Rank         int       `json:"rank"`
	Value        uint      `json:"value"`
	Value2       int       `json:"value_2"`
	CompleteTime time.Time `json:"complete_time"`
	RankType     int16     `json:"-"`
}

//CopyFromUpRank copy
func (u *UpRankInfo) CopyFromUpRank(upRank *UpRank) {
	u.Value = upRank.Value
	u.Value2 = upRank.Value2
	u.RankType = upRank.Type
}

//UpRankQueryResult  result
type UpRankQueryResult struct {
	Result []*UpRankInfo `json:"result"`
	Date   time.Time     `json:"date"`
	PageInfo
}

//PageInfo page info
type PageInfo struct {
	TotalCount int `json:"total_count"`
	Size       int `json:"size"`
	Page       int `json:"page"`
}

//InfoSearchResult  result
type InfoSearchResult struct {
	Result []*InfoQueryResult `json:"result"`
	PageInfo
}

//UpInfoWithViewerData up data with view data
type UpInfoWithViewerData struct {
	Mid         int64                      `json:"mid"`
	UpBaseInfo  *InfoQueryResult           `json:"up_base_info"`
	ViewerTrend *datamodel.ViewerTrendInfo `json:"viewer_trend"`
	ViewerArea  *datamodel.ViewerAreaInfo  `json:"viewer_area"`
	ViewerBase  *datamodel.ViewerBaseInfo  `json:"viewer_base"`
	UpPlayInfo  *UpPlayInfo                `json:"up_play_info"`
}

//UpInfoWithViewerResult info result
type UpInfoWithViewerResult struct {
	Result []*UpInfoWithViewerData `json:"result"`
	PageInfo
}

// -------------

const (
	// FlagUpBaseData up base info
	FlagUpBaseData = 1
	// FlagUpPlayData up play info
	FlagUpPlayData = 1 << 1
	// FlagViewData view base data flag
	FlagViewData = 1 << 2
)

//UpInfoWithViewerArg arg
type UpInfoWithViewerArg struct {
	Mids  string `form:"mids"`
	Sort  string `form:"sort" default:"fans_count"`
	Order string `form:"order" default:"desc"`
	Page  int    `form:"page" default:"1"`
	Size  int    `form:"size" default:"20"`
	// 需要的信息
	Flag int `form:"flag" default:"0"`
}
