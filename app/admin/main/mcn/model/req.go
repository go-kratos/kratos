package model

import (
	"time"

	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

// MCNSignEntryReq req .
type MCNSignEntryReq struct {
	MCNMID      int64         `json:"mcn_mid" validate:"min=1"`
	BeginDate   string        `json:"begin_date" validate:"required"` // 0000-00-00
	EndDate     string        `json:"end_date" validate:"required"`   // 0000-00-00
	SignPayInfo []*SignPayReq `json:"sign_pay_info"`
	Permits     *Permits      `json:"permits"`
	UserName    string
	UID         int64
	Permission  uint32
}

// AttrPermitSet set Permission.
func (req *MCNSignEntryReq) AttrPermitSet() {
	req.Permission = req.Permits.GetAttrPermitVal()
}

// MCNSignPermissionReq .
type MCNSignPermissionReq struct {
	SignID     int64    `json:"sign_id" validate:"required"`
	Permits    *Permits `json:"permits"`
	Permission uint32
	UserName   string
	UID        int64
}

// AttrPermitSet set Permission.
func (req *MCNSignPermissionReq) AttrPermitSet() {
	req.Permission = req.Permits.GetAttrPermitVal()
}

// MCNUPPermitStateReq .
type MCNUPPermitStateReq struct {
	State MCNUPPermissionState `form:"state" validate:"required"`
	PageArg
}

// MCNUPPermitOPReq .
type MCNUPPermitOPReq struct {
	ID           int64                 `json:"id" validate:"min=1"`
	Action       MCNUPPermissionAction `json:"action" validate:"min=1"`
	RejectReason string                `json:"reject_reason"`
	UserName     string
	UID          int64
}

// ParseTime .
func (req *MCNSignEntryReq) ParseTime() (stime, etime xtime.Time, err error) {
	var st, et time.Time
	if st, err = time.ParseInLocation(TimeFormatDay, req.BeginDate, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) error(%+v)", req.BeginDate, err)
		return
	}
	if et, err = time.ParseInLocation(TimeFormatDay, req.EndDate, time.Local); err != nil {
		err = errors.Errorf("time.ParseInLocation(%s) error(%+v)", req.EndDate, err)
		return
	}
	stime = xtime.Time(st.Unix())
	etime = xtime.Time(et.Unix())
	return
}

// SignPayReq  .
type SignPayReq struct {
	DueDate  string `json:"due_date" validate:"required"` // 0000-00-00
	PayValue int64  `json:"pay_value" validate:"min=1"`   // thousand bit
}

// MCNSignInfoReq req
type MCNSignInfoReq struct {
	SignID int64 `form:"sign_id" validate:"min=1"`
}

// MCNSignStateReq req .
type MCNSignStateReq struct {
	State MCNSignState `form:"state" validate:"min=0"`
	PageArg
}

// MCNSignStateOpReq .
type MCNSignStateOpReq struct {
	SignID       int64         `json:"sign_id" validate:"min=1"`
	Action       MCNSignAction `json:"action" validate:"min=0"`
	RejectReason string        `json:"reject_reason"`
	UserName     string
	UID          int64
}

// MCNUPStateReq req .
type MCNUPStateReq struct {
	State MCNUPState `form:"state"  validate:"min=0"`
	PageArg
}

// MCNUPStateOpReq req .
type MCNUPStateOpReq struct {
	SignUpID     int64       `json:"sign_up_id" validate:"min=1"`
	Action       MCNUPAction `json:"action" validate:"min=0"`
	RejectReason string      `json:"reject_reason"`
	UserName     string
	UID          int64
}

// MCNListReq req .
type MCNListReq struct {
	McnCommonReq
	Permits
	ExpireSign       bool         `form:"expire_sign"`
	ExpirePay        bool         `form:"expire_pay"`
	FansNumMin       int64        `form:"fans_num_min"`
	FansNumMax       int64        `form:"fans_num_max"`
	State            MCNSignState `form:"state" default:"-1"`
	SortUP           string       `form:"sort_up"`
	SortAllFans      string       `form:"sort_all_fans"`
	SortRiseFans     string       `form:"sort_rise_fans"`
	SortTrueRiseFans string       `form:"sort_true_rise_fans"`
	SortCheatFans    string       `form:"sort_cheat_fans"`
	Order            string       `form:"order" default:"s.mtime"`
	Sort             string       `form:"sort"  default:"DESC"`
	PageArg
	ExportArg
}

// MCNPayEditReq req .
type MCNPayEditReq struct {
	ID       int64  `json:"id" validate:"min=1"`
	MCNMID   int64  `json:"mcn_mid" validate:"min=1"`
	SignID   int64  `json:"sign_id" validate:"min=1"`
	DueDate  string `json:"due_date" validate:"required"`
	PayValue int64  `json:"pay_value" validate:"min=1"`
	UserName string
	UID      int64
}

// MCNPayStateEditReq req .
type MCNPayStateEditReq struct {
	ID       int64 `json:"id" validate:"min=1"`
	MCNMID   int64 `json:"mcn_mid" validate:"min=1"`
	SignID   int64 `json:"sign_id" validate:"min=1"`
	State    int8  `json:"state"`
	UserName string
	UID      int64
}

// MCNStateEditReq req .
type MCNStateEditReq struct {
	ID       int64         `json:"id" validate:"min=1"`
	MCNMID   int64         `json:"mcn_mid" validate:"min=1"`
	Action   MCNSignAction `json:"action"`
	State    MCNSignState
	UserName string
	UID      int64
}

// MCNRenewalReq req .
type MCNRenewalReq struct {
	ID           int64         `json:"id" validate:"min=1"`
	MCNMID       int64         `json:"mcn_mid" validate:"min=1"`
	BeginDate    string        `json:"begin_date" validate:"required"` // 0000-00-00
	EndDate      string        `json:"end_date" validate:"required"`   // 0000-00-00
	ContractLink string        `json:"contract_link" validate:"required"`
	SignPayInfo  []*SignPayReq `json:"sign_pay_info"`
	Permits      Permits       `json:"permits"`
	Permission   uint32
	UserName     string
	UID          int64
}

// AttrPermitSet set Permission.
func (req *MCNRenewalReq) AttrPermitSet() {
	req.Permission = req.Permits.GetAttrPermitVal()
}

// MCNInfoReq req .
type MCNInfoReq struct {
	McnCommonReq
	ID int64 `form:"id"`
}

// MCNUPListReq req .
type MCNUPListReq struct {
	SignID                     int64      `form:"sign_id" validate:"required"`
	DataType                   int8       `form:"data_type" validate:"min=1"`
	State                      MCNUPState `form:"state" default:"-1"`
	ActiveTID                  int64      `form:"active_tid"`
	FansNumMin                 int64      `form:"fans_num_min"`
	FansNumMax                 int64      `form:"fans_num_max"`
	UPMID                      int64      `form:"up_mid"`
	SortFansCount              string     `form:"sort_fans_count"`
	SortFansCountActive        string     `form:"sort_fans_count_active"`
	SortFansIncreaseAccumulate string     `form:"sort_fans_increase_accumulate"`
	SortArchiveCount           string     `form:"sort_archive_count"`
	SortPlayCount              string     `form:"sort_play_count"`
	SortPubPrice               string     `form:"sort_pub_price"`
	UpType                     int8       `form:"up_type" default:"-1"`
	Order                      string     `form:"order" default:"u.mtime"`
	Sort                       string     `form:"sort"  default:"DESC"`
	Permits
	PageArg
	ExportArg
}

// MCNUPStateEditReq req .
type MCNUPStateEditReq struct {
	ID       int64       `json:"id" validate:"required"`
	SignID   int64       `json:"sign_id" validate:"required"`
	MCNMID   int64       `json:"mcn_mid" validate:"required"`
	UPMID    int64       `json:"up_mid" validate:"required"`
	Action   MCNUPAction `json:"action"`
	State    MCNUPState
	UserName string
	UID      int64
}

// MCNUPRecommendReq req .
type MCNUPRecommendReq struct {
	TID            int64                `form:"tid"`
	UpMid          int64                `form:"up_mid"`
	FansMin        int64                `form:"fans_min"`
	FansMax        int64                `form:"fans_max"`
	PlayMin        int64                `form:"play_min"`
	PlayMax        int64                `form:"play_max"`
	PlayAverageMin int64                `form:"play_average_min"`
	PlayAverageMax int64                `form:"play_average_max"`
	State          MCNUPRecommendState  `form:"state"`
	Source         MCNUPRecommendSource `form:"source"`
	Order          string               `form:"order" default:"mtime"`
	Sort           string               `form:"sort"  default:"DESC"`
	PageArg
	ExportArg
}

// MCNCheatListReq req .
type MCNCheatListReq struct {
	McnCommonReq
	UPMID int64 `form:"up_mid"`
	PageArg
}

// MCNCheatUPListReq struct .
type MCNCheatUPListReq struct {
	UPMID int64 `form:"up_mid" validate:"required"`
	PageArg
}

// MCNImportUPInfoReq struct .
type MCNImportUPInfoReq struct {
	McnCommonReq
	UPMID int64 `form:"up_mid" validate:"required"`
}

// MCNImportUPRewardSignReq struct .
type MCNImportUPRewardSignReq struct {
	SignID   int64 `json:"sign_id" validate:"required"`
	UPMID    int64 `json:"up_mid" validate:"required"`
	UserName string
	UID      int64
}

// RecommendUpReq req .
type RecommendUpReq struct {
	UpMid    int64 `json:"up_mid" validate:"min=1"`
	UserName string
	UID      int64
}

// MCNIncreaseListReq struct .
type MCNIncreaseListReq struct {
	McnCommonReq
	DataType  int8  `form:"data_type"`
	ActiveTID int64 `form:"active_tid" default:"65535"`
	PageArg
}

// RecommendStateOpReq .
type RecommendStateOpReq struct {
	UpMids   []int64              `json:"up_mids"`
	Action   MCNUPRecommendAction `json:"action" validate:"min=1"`
	UserName string
	UID      int64
}

// McnGetRankReq req to 获取排行
type McnGetRankReq struct {
	McnCommonReq
	Tid      int16    `form:"tid"` // 分区 1累计，2昨日，3上周，4上月 0全部
	DataType DataType `form:"data_type"`
	PageArg
	ExportArg
}

// McnCommonReq common mcn
type McnCommonReq struct {
	SignID int64 `form:"sign_id"`
	MCNMID int64 `form:"mcn_mid"`
}

// TotalMcnDataReq .
type TotalMcnDataReq struct {
	Date xtime.Time `form:"date" validate:"required"`
}
