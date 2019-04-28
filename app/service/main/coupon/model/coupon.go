package model

import (
	"fmt"
	"strings"

	"go-common/library/time"
)

// coupon_batch_info表 product_limit_renewal字段.
const (
	ProdLimRenewalAll int8 = iota
	ProdLimRenewalAuto
	ProdLimRenewalNotAuto
)

// coupon_batch_info表 product_limit_renewal字段.
const (
	None           int8 = 0
	ProdLimMonth1       = 1
	ProdLimMonth3       = 3
	ProdLimMonth12      = 12
)

const (
	// CardSalt .
	CardSalt = "7RbjA6mpSz9DYQ0n"
)

// CardType table:coupon_user_card field:card_type
const (
	CardType1 int8 = iota
	CardType3
	CardType12
)

// CardState table:coupon_user_card field:state
const (
	CardStateNotOpen int8 = iota
	CardStateOpened
	CardStateUsed
)

// product limit map .
var (
	ProdLimMonthMap   = map[int8]string{None: "", ProdLimMonth1: "月度", ProdLimMonth3: "季度", ProdLimMonth12: "年度"}
	ProdLimRenewalMap = map[int8]string{ProdLimRenewalAll: "", ProdLimRenewalAuto: "自动续期", ProdLimRenewalNotAuto: "非自动续期"}
)

// MapFullAmount .
var MapFullAmount = map[int8]float64{
	CardType1:  25,
	CardType3:  68,
	CardType12: 233,
}

// CouponChangeLog coupon change log.
type CouponChangeLog struct {
	ID          int64     `json:"-"`
	CouponToken string    `json:"coupon_token"`
	Mid         int64     `json:"mid"`
	State       int8      `json:"state"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

// CouponPageResp coupon page.
type CouponPageResp struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Time  int64  `json:"time"`
	RefID int64  `json:"ref_id"`
	Tips  string `json:"tips"`
	Count int64  `json:"count"`
}

// CouponOrder coupon order info.
type CouponOrder struct {
	ID           int64     `json:"id"`
	OrderNo      string    `json:"order_no"`
	Mid          int64     `json:"mid"`
	Count        int64     `json:"count"`
	State        int8      `json:"state"`
	CouponType   int8      `json:"coupon_type"`
	ThirdTradeNo string    `json:"third_trade_no"`
	Remark       string    `json:"remark"`
	Tips         string    `json:"tips"`
	UseVer       int64     `json:"use_ver"`
	Ver          int64     `json:"ver"`
	Ctime        time.Time `json:"ctime"`
	Mtime        time.Time `json:"mtime"`
}

// CouponOrderLog coupon order log.
type CouponOrderLog struct {
	ID      int64     `json:"id"`
	OrderNo string    `json:"order_no"`
	Mid     int64     `json:"mid"`
	State   int8      `json:"state"`
	Ctime   time.Time `json:"ctime"`
	Mtime   time.Time `json:"mtime"`
}

// CouponBalanceChangeLog coupon balance change log.
type CouponBalanceChangeLog struct {
	ID            int64     `json:"id"`
	OrderNo       string    `json:"order_no"`
	Mid           int64     `json:"mid"`
	BatchToken    string    `json:"batch_token"`
	Balance       int64     `json:"balance"`
	ChangeBalance int64     `json:"change_balance"`
	ChangeType    int8      `json:"change_type"`
	Ctime         time.Time `json:"ctime"`
	Mtime         time.Time `json:"mtime"`
}

// CouponCartoonPageResp coupon cartoon page.
type CouponCartoonPageResp struct {
	Count       int64             `json:"count"`
	CouponCount int64             `json:"coupon_count"`
	List        []*CouponPageResp `json:"list"`
}

// CouponBatchInfo coupon batch info.
type CouponBatchInfo struct {
	ID             int64     `json:"id"`
	AppID          int64     `json:"app_id"`
	Name           string    `json:"name"`
	BatchToken     string    `json:"batch_token"`
	MaxCount       int64     `json:"max_count"`
	CurrentCount   int64     `json:"current_count"`
	LimitCount     int64     `json:"limit_count"`
	StartTime      int64     `json:"start_time"`
	ExpireTime     int64     `json:"expire_time"`
	ExpireDay      int64     `json:"expire_day"`
	Ver            int64     `json:"ver"`
	Ctime          time.Time `json:"ctime"`
	Mtime          time.Time `json:"mtime"`
	FullAmount     float64   `json:"full_amount"`
	Amount         float64   `json:"amount"`
	State          int8      `json:"state"`
	CouponType     int8      `json:"coupon_type"`
	PlatformLimit  string    `json:"platform_limit"`
	ProdLimMonth   int8      `json:"product_limit_month"`
	ProdLimRenewal int8      `json:"product_limit_Renewal"`
}

// CouponAllowancePanelInfo allowance coupon panel info.
type CouponAllowancePanelInfo struct {
	CouponToken         string  `json:"coupon_token"`
	Amount              float64 `json:"coupon_amount"`
	State               int32   `json:"state"`
	FullLimitExplain    string  `json:"full_limit_explain"`
	ScopeExplain        string  `json:"scope_explain"`
	FullAmount          float64 `json:"full_amount"`
	CouponDiscountPrice float64 `json:"coupon_discount_price"`
	StartTime           int64   `json:"start_time"`
	ExpireTime          int64   `json:"expire_time"`
	Selected            int8    `json:"selected"`
	DisablesExplains    string  `json:"disables_explains"`
	OrderNO             string  `json:"order_no"`
	Name                string  `json:"name"`
	Usable              int8    `json:"usable"`
}

// CouponTipInfo coupon tip info.
type CouponTipInfo struct {
	CouponTip  string                    `json:"coupon_tip"`
	CouponInfo *CouponAllowancePanelInfo `json:"coupon_info"`
}

// CouponAllowanceChangeLog coupon allowance change log.
type CouponAllowanceChangeLog struct {
	ID          int64     `json:"-"`
	CouponToken string    `json:"coupon_token"`
	OrderNO     string    `json:"order_no"`
	Mid         int64     `json:"mid"`
	State       int8      `json:"state"`
	ChangeType  int8      `json:"change_type"`
	Ctime       time.Time `json:"ctime"`
	Mtime       time.Time `json:"mtime"`
}

//CouponReceiveLog receive log.
type CouponReceiveLog struct {
	ID          int64  `json:"id"`
	Appkey      string `json:"appkey"`
	OrderNo     string `json:"order_no"`
	Mid         int64  `json:"mid"`
	CouponToken string `json:"coupon_token"`
	CouponType  int8   `json:"coupon_type"`
}

//CouponAllowancePanelResp def.
type CouponAllowancePanelResp struct {
	Usables  []*CouponAllowancePanelInfo `json:"usables"`
	Disables []*CouponAllowancePanelInfo `json:"disables"`
	Using    []*CouponAllowancePanelInfo `json:"using"`
}

// SalaryCouponForThirdResp resp.
type SalaryCouponForThirdResp struct {
	Amount      float64 `json:"amount"`
	FullAmount  float64 `json:"full_amount"`
	Description string  `json:"description"`
}

// ScopeExplainFmt get scope explain fmt.
func (c *CouponAllowancePanelInfo) ScopeExplainFmt(pstr string, prodLimMonth, prodLimRenewal int8, platMap map[string]string) {
	var (
		ps                                  []string
		plats, scope, scopePlat, limr, limm string
	)
	if len(pstr) == 0 && prodLimMonth == 0 && prodLimRenewal == 0 {
		c.ScopeExplain = ScopeNoLimit
		return
	}
	if len(pstr) > 0 {
		ps = strings.Split(pstr, ",")
		for _, v := range ps {
			plats += platMap[v] + ","
		}
	}
	if len(plats) > 0 {
		plats = plats[:len(plats)-1]
		scopePlat = fmt.Sprintf(ScopePlatFmt, plats)
	}
	limr = ProdLimRenewalMap[prodLimRenewal]
	limm = ProdLimMonthMap[prodLimMonth]
	scope = scopePlat + fmt.Sprintf(ScopeProductFmt, limr, limm)
	c.ScopeExplain = scope
}

// PlatfromLimitExplain platform limit explain.
func PlatfromLimitExplain(pstr string, platMap map[string]string) string {
	var (
		ps    []string
		plats string
	)
	if len(pstr) == 0 {
		return ""
	}
	if len(pstr) > 0 {
		ps = strings.Split(pstr, ",")
		for _, v := range ps {
			plats += platMap[v] + ","
		}
	}
	if len(plats) > 0 {
		plats = plats[:len(plats)-1]
	}
	return plats
}

// PrizeCards struct .
type PrizeCards struct {
	List []*PrizeCardRep `json:"list"`
}

// PrizeCardRep struct .
type PrizeCardRep struct {
	CardType      int8   `json:"card_type"`
	State         int8   `json:"state"`
	OriginalPrice int64  `json:"original_price,omitempty"`
	CouponAmount  int64  `json:"coupon_amount,omitempty"`
	DiscountRate  string `json:"discount_rate,omitempty"`
}

// CouponUserCard struct .
type CouponUserCard struct {
	MID         int64  `json:"mid"`
	CardType    int8   `json:"card_type"`
	State       int8   `json:"state"`
	BatchToken  string `json:"batch_token"`
	CouponToken string `json:"coupon_token"`
	ActID       int64  `json:"act_id"`
}
