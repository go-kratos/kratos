package model

import (
	"strconv"
	"time"
)

/* 奖项信息 */

// SaveAwardArg for saving an award
type SaveAwardArg struct {
	// identity
	AwardID int64 `json:"award_id"`
	// properties
	AddAwardArg
}

// AddAwardArg for adding an award
type AddAwardArg struct {
	AwardName     string           `json:"award_name"`     // 专项奖名称
	DisplayStatus int              `json:"display_status"` // 信息是否完整 1=不完整，不前台展示 2=完整，前台展示
	CycleStart    int64            `json:"cycle_start" `   // 评选开始
	CycleEnd      int64            `json:"cycle_end"`      // 评选结束
	AnnounceDate  int64            `json:"announce_date"`  // 公示日期
	Divisions     []*AwardDivision `json:"divisions"`      // 赛区
	Prizes        []*AwardPrize    `json:"prizes"`         // 奖项
	Resources     *AwardResource   `json:"resources"`      // 物料
}

// AwardResource model
type AwardResource struct {
	Rule   string     `json:"rule"`
	Detail string     `json:"detail"`
	QA     []*AwardQA `json:"qa"`
}

// AwardQA model
type AwardQA struct {
	Index int    `json:"-"`
	Q     string `json:"q"`
	A     string `json:"a"`
}

// AwardDivision model
type AwardDivision struct {
	// identity
	AwardID    int64 `json:"-"`
	DivisionID int64 `json:"-"`
	// properties
	DivisionName string `json:"division_name"`
	TagID        int64  `json:"tag_id"`
	Tag          string `json:"tag"`
}

// AwardPrize model
type AwardPrize struct {
	// identity
	AwardID int64 `json:"-"`
	PrizeID int64 `json:"-"`
	// properties
	Bonus int `json:"bonus"`
	Quota int `json:"quota"`
}

// AwardListModel .
type AwardListModel struct {
	ID        int64  `json:"id"`
	AwardID   int64  `json:"award_id"`   // 专项奖ID
	AwardName string `json:"award_name"` // 专项奖名称
	//DisplayStatus   int      `json:"display_status"`   // 是否前台展示(1不展示,2展示)
	CycleStart      int64    `json:"cycle_start"`      // 评选开始
	CycleEnd        int64    `json:"cycle_end"`        // 评选结束
	TotalQuota      int      `json:"total_quota"`      // 总中奖名额
	TotalBonus      int      `json:"total_bonus"`      // 奖金总金额
	AnnounceDate    int64    `json:"announce_date"`    // 公示时间
	OpenStatus      int      `json:"deliver_status"`   // 发奖情况(1未发奖,2已发奖)
	OpenTime        int64    `json:"deliver_time"`     // 发奖时间
	CTime           int64    `json:"created_at"`       // 创建时间
	CreatedBy       string   `json:"created_by"`       // 创建人
	SelectionStatus int      `json:"selection_status"` // 评奖状态(1未评奖,2已评奖)
	DivisionNames   []string `json:"division_names"`   // 分赛区名称列表
	Tags            []string `json:"tags"`             // 分区名称列表
}

// Award model
type Award struct {
	ID              int64     `json:"id"`
	AwardID         int64     `json:"award_id"`         // 专项奖ID
	AwardName       string    `json:"award_name"`       // 专项奖名称
	DisplayStatus   int       `json:"display_status"`   // 是否前台展示(1不展示,2展示)
	CycleStart      time.Time `json:"-"`                // 评选周期开始
	CycleEnd        time.Time `json:"-"`                // 评选周期结束
	CycleStartTS    int64     `json:"cycle_start"`      // 评选开始
	CycleEndTS      int64     `json:"cycle_end"`        // 评选结束
	TotalQuota      int       `json:"total_quota"`      // 总中奖名额
	TotalBonus      int       `json:"total_bonus"`      // 奖金总金额
	AnnounceDate    time.Time `json:"-"`                // 公示日期
	AnnounceDateTS  int64     `json:"announce_date"`    // 公示时间
	OpenStatus      int       `json:"deliver_status"`   // 发奖情况(1未发奖,2已发奖)
	OpenTime        time.Time `json:"-"`                // 开奖时间
	OpenTimeTS      int64     `json:"deliver_time"`     // 发奖时间
	CTime           time.Time `json:"-"`                // 创建时间
	CTimeTS         int64     `json:"created_at"`       // 创建时间
	CreatedBy       string    `json:"created_by"`       // 创建人
	SelectionStatus int       `json:"selection_status"` // 评奖状态(1未评奖,2已评奖)
	IncentiveStart  int64     `json:"incentive_start"`
	IncentiveEnd    int64     `json:"incentive_end"`
}

// GenStr .
func (v *Award) GenStr() {
	v.CycleStartTS = v.CycleStart.Unix()
	v.CycleEndTS = v.CycleEnd.Unix()
	v.AnnounceDateTS = v.AnnounceDate.Unix()
	v.OpenTimeTS = v.OpenTime.Unix()
	v.CTimeTS = v.CTime.Unix()
	v.IncentiveStart = time.Date(v.CycleEnd.Year(), time.Month(v.CycleEnd.Month()+1), 15, 0, 0, 0, 0, time.Local).Unix()
	v.IncentiveEnd = time.Date(v.CycleEnd.Year(), time.Month(v.CycleEnd.Month()+1), 29, 0, 0, 0, 0, time.Local).Unix()
}

// AwardDetail wrapper
type AwardDetail struct {
	Award     *Award           `json:"award"`
	Divisions []*AwardDivision `json:"divisions"`
	Prizes    []*AwardPrize    `json:"prizes"`
	Resources *AwardResource   `json:"resources"`
}

// AwardResult .
type AwardResult struct {
	AwardID      int64                  `json:"award_id"`
	OpenTime     int64                  `json:"deliver_time"`
	AnnounceDate int64                  `json:"announce_date"`
	CycleEnd     int64                  `json:"cycle_end"`
	Divisions    []*AwardDivisionResult `json:"divisions"`
}

// AwardDivisionResult .
type AwardDivisionResult struct {
	DivisionID   int64               `json:"-"`
	DivisionName string              `json:"division_name"`
	Prizes       []*AwardPrizeResult `json:"prizes"`
}

// AwardPrizeResult .
type AwardPrizeResult struct {
	PrizeID int64   `json:"-"`
	MIDs    []int64 `json:"mids"`
}

// AwardRecord model
type AwardRecord struct {
	// identity
	AwardID int64
	MID     int64
	// properties
	TagID int64
}

// AwardWinner model
type AwardWinner struct {
	// identity
	AwardID int64 `json:"award_id"`
	MID     int64 `json:"mid"`
	// properties
	DivisionID int64 `json:"division_id"` //赛区ID
	PrizeID    int64 `json:"prize_id"`    //专项奖奖项ID
	TagID      int64 `json:"-"`           //分区ID
	// derived
	Tag          string `json:"tag"`      //分区
	Nickname     string `json:"nickname"` //昵称
	Bonus        int    `json:"bonus"`    //专项奖奖项奖金
	DivisionName string `json:"division"` //专项奖赛区
}

// QueryAwardWinnerArg .
type QueryAwardWinnerArg struct {
	AwardID  int64  `form:"award_id" validate:"required"`
	MID      int64  `form:"mid"`
	Nickname string `form:"nickname"`
	TagID    int64  `form:"tag_id"`
	From     int    `form:"from" validate:"min=0" default:"0"`
	Limit    int    `form:"limit" validate:"min=1" default:"20"`
}

// AwardWinnerExportFields .
func AwardWinnerExportFields() []string {
	return []string{"UID", "昵称", "分赛区", "奖项类型", "奖金", "分区"}
}

// ExportStrings of an AwardWinner
func (v *AwardWinner) ExportStrings() []string {
	return []string{
		strconv.FormatInt(v.MID, 10),
		v.Nickname,
		v.DivisionName,
		"奖项" + strconv.FormatInt(v.PrizeID, 10),
		strconv.FormatInt(int64(v.Bonus/100), 10),
		v.Tag,
	}
}
