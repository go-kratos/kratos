package model

//风险常量
const (
	//ServerOutage 服务不可用
	ServerOutage = 0
	//ServerNormal 服务正常
	ServerNormal = 1

	//RankNormal 正常
	RankNormal = 0
	//RankAbnormal 不正常
	RankAbnormal = 1
	//RankDoubt 可疑
	RankDoubt = 2

	//MethodPass 通过
	MethodPass = 0
	//MethodBan 禁止
	MethodBan = 1
	//MethodGeetest 极验
	MethodGeetest = 2
	//MethodQuestion 答题
	MethodQuestion = 3

	//VoucherTypePull 凭证拉起
	VoucherTypePull = 1
	//VoucherTypeCheck 凭证验证
	VoucherTypeCheck = 2

	CheckPass      = "验证通过"
	CheckSaleErr   = "未到售卖时间"
	CheckMidEnough = "mid下单次数达到上限"
	CheckIPEnough  = "IP下单次数达到上限"
	CheckIPChange  = "用户网络环境变更"

	RiskLevelSuperHigh = 1
	RiskLevelHigh      = 2
	RiskLevelMiddle    = 3
	RiskLevelLight     = 4
	RiskLevelNormal    = 5
)

// DeviceInfo 设备信息
type DeviceInfo struct {
	UA       string `json:"ua"`
	Info     string `json:"info"`
	Type     string `json:"type"`
	Platform string `json:"platform"`
	Build    string `json:"build"`
}

// ItemInfo 商品
type ItemInfo struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	SaleTime int64  `json:"saleTime"`
	Count    int64  `json:"count"`
	Money    int64  `json:"money"`
}

// BuyerInfo 购买人
type BuyerInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PersonalID  string `json:"personalId"`
	IDCardFront string `json:"idCardFront"`
	IDCardBack  string `json:"idCardBack"`
}

// AddrInfo 收货地址
type AddrInfo struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	ProvID int64  `json:"provId"`
	Prov   string `json:"prov"`
	CityID int64  `json:"cityId"`
	City   string `json:"city"`
	AreaID int64  `json:"areaId"`
	Area   string `json:"area"`
	Addr   string `json:"addr"`
}

// ShieldData .
type ShieldData struct {
	CustomerID    int64      `json:"customerId"`
	UID           string     `json:"uid"`
	TraceID       string     `json:"traceId"`
	Timestamp     int64      `json:"timestamp"`
	UserClientIp  string     `json:"userClientIp"`
	DeviceID      string     `json:"deviceId"`
	SourceIP      string     `json:"sourceIp"`
	InterfaceName string     `json:"interfaceName"`
	PayChannel    string     `json:"payChannel"`
	ReqData       *ReqData   `json:"reqData"`
	ExtShield     *ExtShield `json:"extShield"`
}

// ReqData 业务方信息
type ReqData struct {
	ItemID  []int64 `json:"itemId"`
	AddrID  int64   `json:"addrId"`
	BuyerID int64   `json:"buyerId"`
}

// ExtShield .
type ExtShield struct {
	OrderID      int64  `json:"orderId"`
	RiskLevel    int64  `json:"riskLevel"`
	ShieldResult int64  `json:"shieldResult"`
	ShieldMsg    string `json:"shieldMsg"`
	Source       string `json:"source"`
}

// ShieldIPList .
type ShieldIPList struct {
	IP  string `json:"ip"`
	Num string `json:"num"`
}
