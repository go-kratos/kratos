package service

import (
	"context"
	"github.com/json-iterator/go"
	"go-common/app/common/openplatform/random"
	grpcv1 "go-common/app/service/openplatform/anti-fraud/api/grpc/v1"
	"go-common/app/service/openplatform/anti-fraud/api/http/v1"
	"go-common/app/service/openplatform/anti-fraud/conf"
	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
	"time"
)

//RiskCheck 风险检测
func (s *Service) RiskCheck(c context.Context, req *grpcv1.RiskCheckRequest) (res *grpcv1.RiskCheckResponse, err error) {
	res = new(grpcv1.RiskCheckResponse)
	res.Success = model.ServerNormal

	//售卖时间限制
	if err = s.d.CheckSalesTime(c, req.MID, req.ItemID, req.OpenTime, conf.Conf.Limit.SaleTimeOut); err != nil {
		res.Rank = model.RankAbnormal
		res.Method = model.MethodBan
		res.Desc = model.CheckSaleErr
		return res, nil
	}

	if req.Voucher != "" {
		if err = s.d.CheckVoucher(c, req.MID, req.Voucher, model.VoucherTypeCheck); err != nil {
			return
		}
		res.Desc = model.CheckPass
		return res, nil
	}

	if err = s.d.CheckGeetestCount(c); err != nil {
		log.Info("极验调用达到上限，放行")
		res.Success = model.ServerOutage
		return res, nil
	}
	//ip变更限制
	if err := s.d.CheckIPChange(c, req.MID, req.IP, conf.Conf.Limit.IPChangeInterval); err != nil {
		res.Rank = model.RankDoubt
		res.Method = model.MethodGeetest
		res.Desc = model.CheckIPChange
		return res, nil
	}

	//mid限制
	midKey := model.GetCreateMIDLimitKey(req.MID, req.ItemID)
	if err := s.d.CheckLimitNum(c, midKey, conf.Conf.Limit.MIDCreateMax, conf.Conf.Limit.MIDCreateTimeOut); err != nil {
		res.Rank = model.RankDoubt
		res.Method = model.MethodGeetest
		res.Desc = model.CheckMidEnough
		return res, nil
	}

	//ip限制,绕过白名单
	for _, i := range conf.Conf.Limit.IPWhiteList {
		if i == req.IP {
			res.Desc = model.CheckPass
			return
		}
	}
	ipKey := model.GetCreateIPLimitKey(req.IP, req.ItemID)
	if err := s.d.CheckLimitNum(c, ipKey, conf.Conf.Limit.IPCreateMax, conf.Conf.Limit.IPCreateTimeOut); err != nil {
		res.Rank = model.RankDoubt
		res.Method = model.MethodGeetest
		res.Desc = model.CheckIPEnough
		return res, nil
	}
	res.Desc = model.CheckPass
	return
}

//RiskCheckV2 风险检测
func (s *Service) RiskCheckV2(c context.Context, req *v1.RiskCheckV2Request) (res *v1.RiskCheckV2Response, err error) {
	rule := conf.Conf.Rule[strconv.FormatInt(req.CustomerID, 10)]
	if rule == nil {
		err = ecode.AntiCustomerErr
		return
	}

	res = new(v1.RiskCheckV2Response)
	res.RiskID = random.Uniqid(19)
	res.RiskLevel = model.RiskLevelNormal
	if s.d.CheckBlack(c, req.CustomerID, req.UID, req.UserClientIP) != nil {
		res.RiskLevel = model.RiskLevelSuperHigh
		res.Desc = "黑名单用户"
	}
	return
}

// PayShield 同步支付风控
func (s *Service) PayShield(c context.Context, params *v1.RiskCheckV2Request, res *v1.RiskCheckV2Response) {
	reqData := new(model.ReqData)
	if params.ItemInfo != "" {
		var itemInfo []*model.ItemInfo
		jsoniter.Unmarshal([]byte(params.ItemInfo), &itemInfo)
		for _, value := range itemInfo {
			reqData.ItemID = append(reqData.ItemID, value.ID)
		}
	}
	if params.AddrInfo != "" {
		addrInfo := new(model.AddrInfo)
		jsoniter.Unmarshal([]byte(params.AddrInfo), addrInfo)
		reqData.AddrID = addrInfo.ID
	}
	if params.BuyerInfo != "" {
		buyerInfo := new(model.BuyerInfo)
		jsoniter.Unmarshal([]byte(params.BuyerInfo), buyerInfo)
		reqData.BuyerID = buyerInfo.ID
	}

	extShield := new(model.ExtShield)
	extShield.OrderID = res.RiskID
	extShield.RiskLevel = res.RiskLevel
	extShield.ShieldResult = 1
	extShield.ShieldMsg = res.Desc
	extShield.Source = "ticket_shield"

	data := new(model.ShieldData)
	data.CustomerID = params.CustomerID
	data.UID = strconv.Itoa(int(params.UID))
	data.TraceID = strconv.Itoa(int(res.RiskID))
	data.Timestamp = time.Now().Unix()
	data.UserClientIp = params.UserClientIP
	data.DeviceID = ""
	data.SourceIP = ""
	data.InterfaceName = params.InterfaceName
	data.PayChannel = ""
	data.ReqData = reqData
	data.ExtShield = extShield

	s.d.AddPayData(data)
}

// Ticket2PayShield .
func (s *Service) Ticket2PayShield(c context.Context, req *grpcv1.RiskCheckRequest, res *grpcv1.RiskCheckResponse) {
	reqData := new(model.ReqData)
	reqData.ItemID = append(reqData.ItemID, req.ItemID)

	riskID := random.Uniqid(19)
	riskLevel := model.RiskLevelNormal

	switch res.Rank {
	case model.RankAbnormal:
		riskLevel = model.RiskLevelSuperHigh
	case model.RankDoubt:
		riskLevel = model.RiskLevelMiddle
	}

	extShield := new(model.ExtShield)
	extShield.OrderID = riskID
	extShield.RiskLevel = int64(riskLevel)
	extShield.ShieldResult = 1
	extShield.ShieldMsg = res.Desc
	extShield.Source = "ticket_shield"

	data := new(model.ShieldData)
	data.CustomerID = 10001
	data.UID = strconv.Itoa(int(req.MID))
	data.TraceID = strconv.Itoa(int(riskID))
	data.Timestamp = time.Now().Unix()
	data.UserClientIp = req.IP
	data.DeviceID = ""
	data.SourceIP = ""
	data.InterfaceName = "TICKET_ORDER"
	data.PayChannel = ""
	data.ReqData = reqData
	data.ExtShield = extShield

	s.d.AddPayData(data)
}

// IPList ip列表
func (s *Service) IPList(c context.Context, req *v1.IPListRequest) (res *v1.IPListResponse, err error) {
	res = new(v1.IPListResponse)

	currentTime := time.Now().Unix()
	mtime := time.Unix(currentTime-conf.Conf.Base.ShieldListTime, 0).Format("2006-01-02 15:04:05")
	list, err := s.d.ShieldIPList(c, mtime)

	if err != nil {
		return res, nil
	}

	for _, v := range list {
		ll := new(v1.IPListDetail)
		ll.IP = v.IP
		ll.Num = v.Num
		ll.Timestamp = currentTime
		res.List = append(res.List, ll)
	}

	return
}

// IPDetail ip详情列表
func (s *Service) IPDetail(c context.Context, req *v1.IPDetailRequest) (res *v1.IPDetailResponse, err error) {
	res = new(v1.IPDetailResponse)

	stime := time.Unix(req.Timestamp-conf.Conf.Base.ShieldListTime, 0).Format("2006-01-02 15:04:05")
	etime := time.Unix(req.Timestamp, 0).Format("2006-01-02 15:04:05")
	list, err := s.d.ShieldIPDetail(c, req.IP, stime, etime)

	if err != nil {
		return res, nil
	}

	for _, v := range list {
		ll := new(v1.ListDetail)
		ll.IP = v.IP
		ll.UID = v.UID
		res.List = append(res.List, ll)
	}

	return
}

// UIDList uid列表
func (s *Service) UIDList(c context.Context, req *v1.UIDListRequest) (res *v1.UIDListResponse, err error) {
	res = new(v1.UIDListResponse)

	currentTime := time.Now().Unix()
	mtime := time.Unix(currentTime-conf.Conf.Base.ShieldListTime, 0).Format("2006-01-02 15:04:05")
	list, err := s.d.ShieldUIDList(c, mtime)

	if err != nil {
		return res, nil
	}

	for _, v := range list {
		ll := new(v1.UIDListDetail)
		ll.UID = v.UID
		ll.Num = v.Num
		ll.Timestamp = currentTime
		res.List = append(res.List, ll)
	}

	return
}

// UIDDetail uid详情列表
func (s *Service) UIDDetail(c context.Context, req *v1.UIDDetailRequest) (res *v1.UIDDetailResponse, err error) {
	res = new(v1.UIDDetailResponse)

	stime := time.Unix(req.Timestamp-conf.Conf.Base.ShieldListTime, 0).Format("2006-01-02 15:04:05")
	etime := time.Unix(req.Timestamp, 0).Format("2006-01-02 15:04:05")
	list, err := s.d.ShieldUIDDetail(c, req.UID, stime, etime)

	if err != nil {
		return res, nil
	}

	for _, v := range list {
		ll := new(v1.ListDetail)
		ll.IP = v.IP
		ll.UID = v.UID
		res.List = append(res.List, ll)
	}

	return
}

// IPBlack 设置ip黑名单
func (s *Service) IPBlack(c context.Context, req *v1.IPBlackRequest) (res *v1.IPBlackResponse, err error) {
	res = new(v1.IPBlackResponse)

	key := model.GetIPBlackKey(req.CustomerID, req.IP)

	s.d.SetexRedisKey(c, key, req.Minute*60)
	return
}

// UIDBlack 设置uid黑名单
func (s *Service) UIDBlack(c context.Context, req *v1.UIDBlackRequest) (res *v1.UIDBlackResponse, err error) {
	res = new(v1.UIDBlackResponse)

	uid, err := strconv.ParseInt(req.UID, 10, 64)
	if err != nil {
		log.Warn("black uid err %s", req.UID)
		return res, nil
	}

	key := model.GetMIDBlackKey(req.CustomerID, uid)

	s.d.SetexRedisKey(c, key, req.Minute*60)
	return
}
