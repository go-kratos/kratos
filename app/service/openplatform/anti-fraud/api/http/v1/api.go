package v1

// RiskCheckV2Request .
type RiskCheckV2Request struct {
	UID           int64  `json:"uid" form:"uid" validate:"required"`
	UserClientIP  string `json:"userClientIp" form:"userClientIp" validate:"required"`
	InterfaceName string `json:"interfaceName" form:"interfaceName" validate:"required"`
	InterfaceDesc string `json:"interfaceDesc" form:"interfaceDesc" validate:"required"`
	CustomerID    int64  `json:"customerID" form:"customerId" validate:"required"`
	DeviceInfo    string `json:"deviceInfo" form:"deviceInfo" validate:"required"`
	ItemInfo      string `json:"itemInfo" form:"itemInfo" validate:"required"`
	BuyerInfo     string `json:"buyerInfo" form:"buyerInfo"`
	AddrInfo      string `json:"addrInfo" form:"addrInfo"`
	Voucher       string `json:"voucher" form:"voucher"`
	ReqData       string `json:"reqData" form:"reqData" validate:"required"`
	ExtraData     string `json:"extraData" form:"extraData"`
}

// RiskCheckV2Response .
type RiskCheckV2Response struct {
	RiskID    int64  `json:"riskId"`
	RiskLevel int64  `json:"riskLevel"`
	Method    string `json:"method"`
	Desc      string `json:"desc"`
}

// IPListRequest .
type IPListRequest struct {
}

// IPListResponse .
type IPListResponse struct {
	List []*IPListDetail `json:"list"`
}

// IPListDetail .
type IPListDetail struct {
	IP        string `json:"ip"`
	Num       int64  `json:"num"`
	Timestamp int64  `json:"timestamp"`
}

// UIDListRequest .
type UIDListRequest struct {
}

// UIDListResponse .
type UIDListResponse struct {
	List []*UIDListDetail `json:"list"`
}

// UIDListDetail .
type UIDListDetail struct {
	UID       string `json:"uid"`
	Num       int64  `json:"num"`
	Timestamp int64  `json:"timestamp"`
}

// IPDetailRequest .
type IPDetailRequest struct {
	IP        string `json:"ip" form:"ip" validate:"required"`
	Timestamp int64  `json:"timestamp" form:"timestamp" validate:"required"`
}

// IPDetailResponse .
type IPDetailResponse struct {
	List []*ListDetail `json:"list"`
}

// UIDDetailRequest .
type UIDDetailRequest struct {
	UID       string `json:"uid" form:"uid" validate:"required"`
	Timestamp int64  `json:"timestamp" form:"timestamp" validate:"required"`
}

// UIDDetailResponse .
type UIDDetailResponse struct {
	List []*ListDetail `json:"list"`
}

// ListDetail .
type ListDetail struct {
	UID string `json:"uid"`
	IP  string `json:"ip"`
}

// IPBlackRequest .
type IPBlackRequest struct {
	IP         string `json:"ip" form:"ip" validate:"required"`
	CustomerID int64  `json:"customer_id" form:"customer_id" validate:"required"`
	Minute     int64  `json:"minute" form:"minute" validate:"required"`
}

// IPBlackResponse .
type IPBlackResponse struct {
}

// UIDBlackRequest .
type UIDBlackRequest struct {
	UID        string `json:"uid" form:"uid" validate:"required"`
	CustomerID int64  `json:"customer_id" form:"customer_id" validate:"required"`
	Minute     int64  `json:"minute" form:"minute" validate:"required"`
}

// UIDBlackResponse .
type UIDBlackResponse struct {
}
