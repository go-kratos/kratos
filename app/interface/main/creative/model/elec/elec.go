package elec

import (
	"encoding/json"
	"time"

	xtime "go-common/library/time"
)

// UserState user elec state.
type UserState struct {
	ID     string `json:"-"`
	Mid    string `json:"mid"`
	State  string `json:"state"`
	Reason string `json:"reason"`
	Count  string `json:"-"`
	CTime  string `json:"-"`
	MTime  string `json:"-"`
}

// UserInfo user elec info.
type UserInfo struct {
	ID     int64      `json:"-"`
	Mid    int64      `json:"mid"`
	State  int16      `json:"state"`
	Reason string     `json:"reason"`
	Count  int16      `json:"-"`
	CTime  xtime.Time `json:"-"`
	MTime  xtime.Time `json:"-"`
}

// ArcState arc elec info.
type ArcState struct {
	Show   bool            `json:"show"`
	State  int16           `json:"state"`
	Total  int             `json:"total"`
	Count  int             `json:"count"`
	Reason string          `json:"reason"`
	List   json.RawMessage `json:"list,omitempty"`
	User   json.RawMessage `json:"user,omitempty"`
}

// Notify up-to-date info to user
type Notify struct {
	Content string `json:"content"`
}

// EleRelation get elec relation.
type EleRelation struct {
	RetList []struct {
		Mid    int64 `json:"mid"`
		IsElec bool  `json:"is_elec"`
	} `json:"ret_list"`
}

// Status elec setting.
type Status struct {
	Specialday int8 `json:"display_specialday"`
}

// Rank up rank.
type Rank struct {
	MID      int64  `json:"mid"`
	PayMID   int64  `json:"pay_mid"`
	Rank     int64  `json:"rank"`
	Uname    string `json:"uname"`
	Avatar   string `json:"avatar"`
	IsFriend bool   `json:"isfriend"`
	MTime    string `json:"mtime"`
}

// BillList daily bill list.
type BillList struct {
	List       []*Bill `json:"list"`
	TotalCount int     `json:"totalCount"`
	Pn         int     `json:"pn"`
	Ps         int     `json:"ps"`
}

// Bill bill detail.
type Bill struct {
	ID            int64      `json:"id"`
	MID           int64      `json:"mid"`
	ChannelType   int8       `json:"channelType"`
	ChannelTyName string     `json:"channelTypeName"`
	AddNum        float32    `json:"addNum"`
	ReduceNum     float32    `json:"reduceNum"`
	WalletBalance float32    `json:"walletBalance"`
	DateVersion   string     `json:"dateVersion"`
	Weekday       string     `json:"weekday"`
	Remark        string     `json:"remark"`
	MonthBillResp *MonthBill `json:"monthBillResp"`
}

// MonthBill month bill.
type MonthBill struct {
	LastMonthNum float32 `json:"last_month_num"`
	ServiceNum   float32 `json:"service_num"`
	BkNum        float32 `json:"bk_num"`
}

// Balance get battery balance.
type Balance struct {
	Ts             string       `json:"ts"`
	BrokerageAudit int8         `json:"brokerage_audit"`
	BpayAcc        *BpayAccount `json:"bpay_account"`
	Wallet         *Wall        `json:"wallet"`
}

// BpayAccount shell detail.
type BpayAccount struct {
	Brokerage float32 `json:"brokerage"`
	DefaultBp float32 `json:"default_bp"`
}

// Wall wallet detail.
type Wall struct {
	MID            int64   `json:"mid"`
	Balance        float32 `json:"balance"`
	SponsorBalance float32 `json:"sponsorBalance"`
	Ver            int32   `json:"-"`
}

// ChargeBill daily bill for app charge.
type ChargeBill struct {
	List  []*Bill `json:"list"`
	Pager struct {
		Current int `json:"current"`
		Size    int `json:"size"`
		Total   int `json:"total"`
	} `json:"pager"`
}

// RecentElec recent detail for app.
type RecentElec struct {
	AID     int64   `json:"aid"`
	MID     int64   `json:"mid"`
	ElecNum float32 `json:"elec_num"`
	Title   string  `json:"title"`
	Uname   string  `json:"uname"`
	Avatar  string  `json:"avatar"`
	OrderNO string  `json:"-"`
	CTime   string  `json:"ctime"`
}

// RecentElecList recent list for app.
type RecentElecList struct {
	List  []*RecentElec `json:"list"`
	Pager struct {
		Current int `json:"current"`
		Size    int `json:"size"`
		Total   int `json:"total"`
	} `json:"pager"`
}

// RemarkList remark list.
type RemarkList struct {
	List  []*Remark `json:"list"`
	Pager struct {
		Current int `json:"current"`
		Size    int `json:"size"`
		Total   int `json:"total"`
	} `json:"pager"`
}

// Remark remark detail.
type Remark struct {
	ID          int64      `json:"id"`
	AID         int64      `json:"aid"`
	MID         int64      `json:"mid"`
	ReplyMID    int64      `json:"reply_mid"`
	ElecNum     int64      `json:"elec_num"`
	State       int8       `json:"state"`
	Msg         string     `json:"msg"`
	ArcName     string     `json:"aname"`
	Uname       string     `json:"uname"`
	Avator      string     `json:"avator"`
	ReplyName   string     `json:"reply_name"`
	ReplyAvator string     `json:"reply_avator"`
	ReplyMsg    string     `json:"reply_msg"`
	CTime       xtime.Time `json:"ctime"`
	ReplyTime   xtime.Time `json:"reply_time"`
}

// Earnings for elec.
type Earnings struct {
	State     int8    `json:"state"`
	Balance   float32 `json:"balance"`
	Brokerage float32 `json:"brokerage"`
}

// Weekday get day.
func Weekday(t time.Time) (w string) {
	switch t.Weekday().String() {
	case "Monday":
		w = "周一"
	case "Tuesday":
		w = "周二"
	case "Wednesday":
		w = "周三"
	case "Thursday":
		w = "周四"
	case "Friday":
		w = "周五"
	case "Saturday":
		w = "周六"
	case "Sunday":
		w = "周日"
	}
	return
}
