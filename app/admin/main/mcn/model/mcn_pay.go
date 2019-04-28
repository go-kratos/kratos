package model

import xtime "go-common/library/time"

// MCNSignPay State .
const (
	MCNStateNoPay = int8(1)
	MCNStatePayed = int8(2)
	MCNStateDeled = int8(100)
)

// MCNSign struct .
type MCNSign struct {
	ID                 int64        `json:"id"`
	MCNMID             int64        `json:"mcn_mid"`
	MCNName            string       `json:"mcn_name"`
	CompanyName        string       `json:"company_name"`
	CompanyLicenseID   string       `json:"company_license_id"`
	CompanyLicenseLink string       `json:"company_license_link"`
	ContractLink       string       `json:"contract_link"`
	ContactName        string       `json:"contact_name"`
	ContactTitle       string       `json:"contact_title"`
	ContactIdcard      string       `json:"contact_idcard"`
	ContactPhone       string       `json:"contact_phone"`
	BeginDate          xtime.Time   `json:"begin_date"`
	EndDate            xtime.Time   `json:"end_date"`
	State              MCNSignState `json:"state"`
	RejectTime         xtime.Time   `json:"reject_time"`
	RejectReason       string       `json:"reject_reason"`
	Ctime              xtime.Time   `json:"ctime"`
	Mtime              xtime.Time   `json:"mtime"`
	Permission         uint32       `json:"permission"`
	Permits            *Permits     `json:"permits"` // 权限集合
}

// AttrPermitVal get Permission all.
func (n *MCNSign) AttrPermitVal() {
	n.Permits = &Permits{}
	n.Permits.SetAttrPermitVal(n.Permission)
}

// MCNSignPay struct .
type MCNSignPay struct {
	ID       int64  `json:"id"`
	MID      int64  `json:"mid"`
	SignID   int64  `json:"sign_id"`
	DueDate  string `json:"due_date"`
	PayValue int64  `json:"pay_value"`
	State    int8   `json:"state"`
	Note     string `json:"note"`
	Ctime    string `json:"ctime"`
	Mtime    string `json:"mtime"`
}

// MCNUP struct .
type MCNUP struct {
	SignID           int64      `json:"sign_id"`
	MCNMID           int64      `json:"mcn_mid"`
	UPMID            int64      `json:"up_mid"`
	BeginDate        xtime.Time `json:"begin_date"`
	EndDate          xtime.Time `json:"end_date"`
	ContractLink     string     `json:"contract_link"`
	UPAuthLink       string     `json:"up_auth_link"`
	RejectReason     string     `json:"reject_reason"`
	RejectTime       xtime.Time `json:"reject_time"`
	State            MCNUPState `json:"state"`
	StateChangeTime  xtime.Time `json:"state_change_time"`
	Ctime            xtime.Time `json:"ctime"`
	Mtime            xtime.Time `json:"mtime"`
	UpType           int8       `json:"up_type"`
	SiteLink         string     `json:"site_link"`
	ConfirmTime      xtime.Time `json:"confirm_time"`
	Permission       uint32     `json:"permission"`
	PublicationPrice int64      `json:"publication_price"`
	Permits          *Permits   `json:"permits"` // 权限集合
}

// AttrPermitVal get Permission all.
func (n *MCNUP) AttrPermitVal() {
	n.Permits = &Permits{}
	n.Permits.SetAttrPermitVal(n.Permission)
}
