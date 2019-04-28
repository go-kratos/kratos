package model

import (
	"strconv"

	"go-common/app/service/main/account/api"
)

// YstCreateOrderReq represents request of yst order creation.
type YstCreateOrderReq struct {
	SeqNo     string `json:"seqno" url:"seqno"`
	Source    string `json:"source" url:"source"`
	ProductId string `json:"product_id" url:"product_id"`
	Price     int32  `json:"price" url:"price"`
	Total     int32  `json:"total" url:"total"`
	BuyNum    int32  `json:"buy_num" url:"buy_num"`
	VideoType string `json:"video_type" url:"video_type"`
	PayType   string `json:"pay_type" url:"pay_type"`
	PayWay    string `json:"pay_way" url:"pay_way"`
	UserLogin string `json:"user_login" url:"user_login"`
	UserId    string `json:"user_id" url:"user_id"`
	GUID      string `json:"guid" url:"guid"`
	ClientIp  string `json:"client_ip" url:"client_ip"`
	Sign      string `json:"sign" url:"sign,omitempty"`
	LoginName string `json:"login_name" url:"login_name"`
}

// CopyFromPayOrder copies fields from pay order.
func (y *YstCreateOrderReq) CopyFromPayOrder(po *PayOrder) {
	y.SeqNo = po.OrderNo
	y.Source = "snm_bilibili"
	y.ProductId = po.ProductId
	y.Price = po.Money
	y.Total = po.PaymentMoney
	y.BuyNum = po.Quantity
	if po.OrderType == 0 {
		y.VideoType = "fvod"
	}
	if po.OrderType == 1 {
		y.VideoType = "svod"
	}
	if po.PaymentType == PaymentTypeAliPay {
		y.PayType = YstPayTypeAliPay
	}
	if po.PaymentType == PaymentTypeWechat {
		y.PayType = YstPayTypeWechat
	}
	y.PayWay = YstPayWayQr
}

// CopyFromAccount copies fields from account info.
func (y *YstCreateOrderReq) CopyFromAccount(account *api.Info) {
	y.LoginName = account.Name
	y.UserId = strconv.Itoa(int(account.Mid))
}

// YstCreateOrderReply represents response of yst order creation.
type YstCreateOrderReply struct {
	SeqNo        string `json:"seqno" url:"seqno"`
	TraceNo      string `json:"traceno" url:"traceno"`
	PayWary      string `json:"pay_wary" url:"pay_wary"`
	CodeUrl      string `json:"code_url" url:"code_url"`
	ContractCode string `json:"contract_code" url:"contract_code"`
	Price        int32  `json:"price" url:"price"`
	VideoType    string `json:"video_type" url:"video_type"`
	PayParam     string `json:"pay_param" url:"pay_param"`
	ResultCode   string `json:"result_code" url:"result_code"`
	ResultMsg    string `json:"result_msg" url:"result_msg"`
	Sign         string `json:"sign"`
}

// YstPayCallbackReq represents request of pay callback.
type YstPayCallbackReq struct {
	SeqNo      string `json:"seqno" url:"seqno"`
	TraceNo    string `json:"traceno" url:"traceno"`
	TradeState string `json:"trade_state" url:"trade_state"`
	ContractId string `json:"contract_id" url:"contract_id,omitempty"`
	Sign       string `json:"sign" url:"sign"`
}

// YstPayCallbackReply represents response of pay callback.
type YstPayCallbackReply struct {
	TraceNo string
	Result  string
	Msg     string
}

// PayInfo represents short pay details.
type PayInfo struct {
	OrderNo      string
	PaymentType  string
	CodeUrl      string
	PaymentMoney int32
}

// CopyFromPayOrder copies fields from pay order.
func (p *PayInfo) CopyFromPayOrder(po *PayOrder) {
	p.OrderNo = po.OrderNo
	p.PaymentType = po.PaymentType
	p.PaymentMoney = po.PaymentMoney
}

// YstRenewOrderReq.
type YstRenewOrderReq struct {
	SeqNo      string `json:"seqno" url:"seqno"`
	Source     string `json:"source" url:"source"`
	ProductId  string `json:"product_id" url:"product_id"`
	Price      int32  `json:"price" url:"price"`
	BuyNum     int32  `json:"buy_num" url:"buy_num"`
	Total      int32  `json:"total" url:"total"`
	VideoType  string `json:"video_type" url:"video_type"`
	PayType    string `json:"pay_type" url:"pay_type"`
	UserId     string `json:"user_id" url:"user_id"`
	ContractId string `json:"contract_id" url:"contract_id"`
	Sandbox    string `json:"sandbox" url:"sandbox"`
	ClientIp   string `json:"client_ip" url:"client_ip"`
	Sign       string `json:"sign"`
}

// YstRenewOrderReply.
type YstRenewOrderReply struct {
	SeqNo      string `json:"seqno" url:"seqno"`
	TraceNo    string `json:"traceno" url:"traceno"`
	Price      int32  `json:"price" url:"price"`
	VideoType  string `json:"video_type" url:"video_type"`
	ResultCode string `json:"result_code" url:"result_code"`
	ResultMsg  string `json:"result_msg" url:"result_msg"`
	Sign       string `json:"sign" `
}

// YstOrderState.
type YstOrderStateReq struct {
	SeqNo   string `json:"seqno" url:"seqno"`
	TraceNo string `json:"traceno" url:"traceno"`
	Sign    string `json:"sign" url:"sign"`
}

// YstOrderStateReply.
type YstOrderStateReply struct {
	SeqNo     string `json:"seqno" `
	TraceNo   string `json:"traceno" `
	PayStatus string `json:"pay_status"`
	Result    string `json:"result"`
	Msg       string `json:"msg"`
}

// YstUserInfoReq.
type YstUserInfoReq struct {
	Mid  int32  `json:"mid" url:"mid"`
	Sign string `json:"sign" url:"sign,omitempty"`
}

// WxContractCallbackReq.
type WxContractCallbackReq struct {
	ContractId              string `json:"contract_id" url:"contract_id"`
	ContractCode            string `json:"contract_code" url:"contract_code"`
	ChangeType              string `json:"change_type" url:"contract_id"`
	ContractTerminationMode string `json:"contract_termination_mode" url:"contract_termination_mode,omitempty"`
	Sign                    string `json:"sign" url:"sign"`
}

// WxContractCallbackReply.
type WxContractCallbackReply struct {
	ContractId string `json:"contract_id"`
	Result     string `json:"result"`
	Msg        string `json:"msg"`
}
