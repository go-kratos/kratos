package model

import (
	pb "go-common/app/service/main/tv/api"
)

type CreateQrReq struct {
	Mid        int64  `form:"mid" validate:"required"`
	Pid        int32  `form:"pid" validate:"required"`
	BuyNum     int32  `form:"buy_num" validate:"required"`
	AppChannel string `form:"app_channel"`
	Guid       string
}

func (in *CreateQrReq) CopyIntoPbCreateOrReq(out *pb.CreateQrReq) {
	out.Mid = in.Mid
	out.Pid = in.Pid
	out.BuyNum = in.BuyNum
	out.Guid = in.Guid
	out.AppChannel = in.AppChannel
}

type CreateGuestQrReq struct {
	Pid        int32  `form:"pid" validate:"required"`
	BuyNum     int32  `form:"buy_num" validate:"required"`
	AppChannel string `form:"app_channel"`
	Guid       string
}

func (in *CreateGuestQrReq) CopyIntoPbCreateGuestQrReq(out *pb.CreateGuestQrReq) {
	out.Pid = in.Pid
	out.BuyNum = in.BuyNum
	out.Guid = in.Guid
	out.AppChannel = in.AppChannel
}

type CreateOrderReq struct {
	Token       string `form:"token" validate:"required"`
	Platform    int8
	PaymentType string
}

func (in *CreateOrderReq) CopyIntoPbCreateOrderReq(out *pb.CreateOrderReq) {
	out.Token = in.Token
	out.Platform = in.Platform
	out.PaymentType = in.PaymentType
}

type CreateGuestOrderReq struct {
	Token       string `form:"token" validate:"required"`
	Platform    int8
	PaymentType string
}

func (in *CreateGuestOrderReq) CopyIntoPbCreateGuestOrderReq(out *pb.CreateGuestOrderReq) {
	out.Token = in.Token
	out.Platform = in.Platform
	out.PaymentType = in.PaymentType
}

type WxContractCallbackReq struct {
	ContractId              string `json:"contract_id" url:"contract_id"  validate:"required"`
	ContractCode            string `json:"contract_code" url:"contract_code"  validate:"required"`
	ChangeType              string `json:"change_type" url:"contract_id"  validate:"required"`
	ContractTerminationMode string `json:"contract_termination_mode" url:"contract_termination_mode,omitempty"`
	Sign                    string `json:"sign" url:"sign"  validate:"required"`
}

func (in *WxContractCallbackReq) CopyIntoPbWxContractCallbackReq(out *pb.WxContractCallbackReq) {
	out.ContractId = in.ContractId
	out.ContractCode = in.ContractCode
	out.ChangeType = in.ChangeType
	out.ContractTerminationMode = in.ContractTerminationMode
	out.Sign = in.Sign
}

type YstPayCallbackReq struct {
	SeqNo      string `json:"seqno" validate:"required"`
	TraceNo    string `json:"traceno" validate:"required"`
	TradeState string `json:"trade_state" validate:"required"`
	ContractId string `json:"contract_id"`
	Sign       string `json:"sign" validate:"required"`
}

func (in *YstPayCallbackReq) CopyIntoPbPayCallbackReq(out *pb.PayCallbackReq) {
	out.SeqNo = in.SeqNo
	out.TraceNo = in.TraceNo
	out.TradeState = in.TradeState
	out.ContractId = in.ContractId
	out.Sign = in.Sign
}

type YstUserInfoReq struct {
	Mid  int64  `form:"mid" validate:"required"`
	Sign string `form:"sign" validate:"required"`
}
