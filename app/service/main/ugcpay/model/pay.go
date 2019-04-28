package model

// PayCallbackMSG .
type PayCallbackMSG struct {
	CustomerID     int64  `json:"customerId"`     //业务id
	ServiceType    int64  `json:"serviceType"`    //业务方业务类型
	TXID           int64  `json:"txId"`           //支付平台支付id
	OrderID        string `json:"orderId"`        //业务方订单id
	DeviceType     int64  `json:"deviceType"`     //支付设备渠道类型，  1 pc 2 webapp 3 app 4jsapi 5 server 6小程序支付 7聚合二维码支付
	PayStatus      string `json:"payStatus"`      //支付状态，FINISHED（交易成功）|SUCCESS（成功）|REFUND(退款中)|PAYING（支付中）|CLOSED（关闭）|NOT_PAY（未支付）|FAIL(支付失败)|WITHDRAW(支付撤销) 支付回调暂时仅通知 SUCCESS, 其他状态不通知。
	PayChannelID   int64  `json:"payChannelId"`   //支付渠道id, 用户实际选择的支付实体渠道。(payChannel 代表笼统的微信、支付宝等第三方渠道， payChannelId 代表实际签约的实体渠道 id)
	PayChannel     string `json:"payChannel"`     //支付渠道，alipay(支付宝)、wechat(微信) ,paypal(paypal), iap(In App Purchase)、qpay(QQ支付)、huabei(花呗支付)、ali_bank（网银支付）、bocom（交行信用卡支付）、bp（B币支付）
	PayChannelName string `json:"payChannelName"` //支付渠道名称 如支付宝、微信、PayPal、IAP、QQ、花呗、网银支付、B币支付
	PayAccount     string `json:"payAccount"`     //支付渠道账号
	PayBank        string `json:"payBank"`        //支付银行
	FeeType        string `json:"feeType"`        //货币类型，默认人民币CNY
	PayAmount      int64  `json:"payAmount"`      //实际支付金额
	PayMsgContent  string `json:"payMsgContent"`  //支付返回的额外信息，json格式字符串，比如：payCounponAmount：使用B币券金额（单位 分），payBpAmount：B币金额（分）
	ExtData        string `json:"extData"`        //支付请求时的扩展json串
	ExpiredTime    int64  `json:"expiredTime"`    //IAP代扣过期时间，毫秒值,业务方需要判断expiredTime的值，因为重复通知返回的expiredTime是一样的
	OrderPayTime   string `json:"orderPayTime"`   //订单支付时间，格式：0000-00-00 00:00:00
	Timestamp      string `json:"timestamp"`      //请求时间戳，毫秒
	TraceID        string `json:"traceId"`        //追踪id
	SignType       string `json:"signType"`       //签名类型 ，默认MD5
	Sign           string `json:"sign"`           //签名（应当支持支付平台这边新增返回字段）
}

// IsSuccess 是否为成功支付内容
func (p *PayCallbackMSG) IsSuccess() bool {
	return p.PayStatus == PayStateFinished || p.PayStatus == PayStateSuccess
}

// PayRefundCallbackMSG .
type PayRefundCallbackMSG struct {
	CustomerID int64                  `json:"customerId"`      //业务id
	OrderID    string                 `json:"orderId"`         //业务方订单id
	TXID       int64                  `json:"txId"`            //支付平台支付id
	List       []PayRefundCallbackEle `json:"batchRefundList"` //
}

// PayRefundCallbackEle .
type PayRefundCallbackEle struct {
	RefundStatus     string `json:"refundStatus"`
	RefundStatusDesc string `json:"refundStatusDesc"`
	RefundEndTime    string `json:"refundEndTime"`
}

// IsSuccess 是否为成功退款内容
func (p *PayRefundCallbackEle) IsSuccess() bool {
	return p.RefundStatus == PayStateRefund
}

// RechargeShellCallbackMSG .
type RechargeShellCallbackMSG struct {
	CustomerID   int64  `json:"customerId"`
	Status       string `json:"status"`
	ThirdOrderNo string `json:"thirdOrderNo"`
	MID          int64  `json:"mid"`
}

// PayQuery .
type PayQuery struct {
	Orders []*PayOrder `json:"orders"`
}

// PayOrder .
type PayOrder struct {
	TXID          int64  `json:"txId"`
	OrderID       string `json:"orderId"`
	PayStatus     string `json:"payStatus"`
	PayStatusDesc string `json:"payStatusDesc"`
	FailReason    string `json:"failReason"`
}
