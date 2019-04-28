package v1

import (
	"context"
	"encoding/json"
	"go-common/library/ecode"
	"regexp"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/status"

	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
	rspb "go-common/app/service/live/xlottery/api/grpc/v1"
	"go-common/library/log"
)

// CapsuleService struct
type CapsuleService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	capsuleCli *rspb.Client
}

//NewCapsuleService init
func NewCapsuleService(c *conf.Config) (s *CapsuleService) {
	s = &CapsuleService{
		conf: c,
	}
	var svc *rspb.Client
	var err error
	log.Info("CapsuleService Init: %+v", s.conf.CapsuleClient)
	if svc, err = rspb.NewClient(s.conf.CapsuleClient); err != nil {
		panic(err)
	}
	s.capsuleCli = svc
	return s
}

// GetCoinList implementation
func (s *CapsuleService) GetCoinList(ctx context.Context, req *v1pb.GetCoinListReqAdmin) (resp *v1pb.GetCoinListRespAdmin, err error) {
	var RPCResponse *rspb.GetCoinListResp
	RPCReq := &rspb.GetCoinListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if RPCReq.Page < 1 {
		RPCReq.Page = 1
	}
	if RPCReq.PageSize < 1 {
		RPCReq.PageSize = 16
	}

	RPCResponse, err = s.capsuleCli.GetCoinList(ctx, RPCReq)
	if err != nil {
		return
	}
	coinListResp := make([]*v1pb.GetCoinListRespAdmin_List, 0)
	for _, singleCoinList := range RPCResponse.List {
		coinList := &v1pb.GetCoinListRespAdmin_List{}
		coinList.Id = singleCoinList.Id
		coinList.Status = singleCoinList.Status
		coinList.Title = singleCoinList.Title
		coinList.GiftConfig = singleCoinList.GiftConfig
		coinList.GiftType = singleCoinList.GiftType
		coinList.ChangeNum = singleCoinList.ChangeNum
		startTimeUnix := time.Unix(singleCoinList.StartTime, 0)
		coinList.StartTime = startTimeUnix.Format("2006-01-02 15:04:05")
		endTimeUnix := time.Unix(singleCoinList.EndTime, 0)
		coinList.EndTime = endTimeUnix.Format("2006-01-02 15:04:05")
		areaIds := make([]*v1pb.GetCoinListRespAdmin_List_AreaIds, 0)
		for _, areaId := range singleCoinList.AreaIds {
			singleAreaId := &v1pb.GetCoinListRespAdmin_List_AreaIds{}
			singleAreaId.IsAll = areaId.IsAll
			singleAreaId.List = areaId.List
			singleAreaId.ParentId = areaId.ParentId
			areaIds = append(areaIds, singleAreaId)
		}
		coinList.AreaIds = areaIds
		coinListResp = append(coinListResp, coinList)
	}
	resp = &v1pb.GetCoinListRespAdmin{
		Total:     RPCResponse.Total,
		TotalPage: RPCResponse.TotalPage,
		List:      coinListResp,
	}
	return
}

// UpdateCoinConfig implementation
func (s *CapsuleService) UpdateCoinConfig(ctx context.Context, req *v1pb.UpdateCoinConfigReqAdmin) (resp *v1pb.UpdateCoinConfigRespAdmin, err error) {
	var RPCResponse *rspb.UpdateCoinConfigResp

	loc, _ := time.LoadLocation("Local")
	var endTime, startTime time.Time
	startTime, err = time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, loc)
	if err != nil {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return nil, err
	}
	endTime, err = time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, loc)
	if err != nil {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return nil, err
	}
	var areaIds []*rspb.UpdateCoinConfigReq_AreaIds
	err = json.Unmarshal([]byte(req.AreaIds), &areaIds)
	if err != nil {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return nil, err
	}
	if len(areaIds) == 0 {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return nil, err
	}
	var isHasArea bool
	for _, areaInfo := range areaIds {
		if areaInfo.IsAll == 1 {
			isHasArea = true
			break
		}
		if len(areaInfo.List) > 0 {
			isHasArea = true
			break
		}
	}
	if !isHasArea {
		err = ecode.Error(ecode.XLotteryCapsuleAreaParamErr, "必须设置生效分区")
		return nil, err
	}
	var giftIds = make([]int64, 0)
	if req.GiftType == 3 {
		if req.GiftConfig == "" {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return nil, err
		}
		giftConfig := strings.Split(req.GiftConfig, ",")
		for _, giftStr := range giftConfig {
			var giftId int64
			giftId, err = strconv.ParseInt(giftStr, 10, 64)
			if err != nil {
				err = ecode.Error(ecode.InvalidParam, "参数错误")
				return nil, err
			}
			giftIds = append(giftIds, giftId)
		}
		if len(giftIds) == 0 {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return nil, err
		}
		match, _ := regexp.Compile("^([0-9]+,)*[0-9]+$")
		if !match.MatchString(req.GiftConfig) {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return nil, err
		}
	}
	RPCReq := &rspb.UpdateCoinConfigReq{
		Id:        req.Id,
		Title:     req.Title,
		ChangeNum: req.ChangeNum,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
		Status:    req.Status,
		GiftType:  req.GiftType,
		GiftIds:   giftIds,
		AreaIds:   areaIds,
	}
	RPCResponse, err = s.capsuleCli.UpdateCoinConfig(ctx, RPCReq)
	if err != nil {
		return
	}
	resp = &v1pb.UpdateCoinConfigRespAdmin{
		Status: RPCResponse.Status,
	}
	return
}

// UpdateCoinStatus implementation
func (s *CapsuleService) UpdateCoinStatus(ctx context.Context, req *v1pb.UpdateCoinStatusReqAdmin) (resp *v1pb.UpdateCoinStatusRespAdmin, err error) {
	var RPCResponse *rspb.UpdateCoinStatusResp
	RPCReq := &rspb.UpdateCoinStatusReq{
		Id:     req.Id,
		Status: req.Status,
	}
	if RPCReq.Status != 1 {
		RPCReq.Status = 0
	}
	RPCResponse, err = s.capsuleCli.UpdateCoinStatus(ctx, RPCReq)
	if err != nil {
		return
	}
	resp = &v1pb.UpdateCoinStatusRespAdmin{
		Status: RPCResponse.Status,
	}
	return
}

// DeleteCoin implementation
func (s *CapsuleService) DeleteCoin(ctx context.Context, req *v1pb.DeleteCoinReqAdmin) (resp *v1pb.DeleteCoinRespAdmin, err error) {
	var RPCResponse *rspb.DeleteCoinResp
	RPCReq := &rspb.DeleteCoinReq{
		Id: req.Id,
	}
	RPCResponse, err = s.capsuleCli.DeleteCoin(ctx, RPCReq)
	if err != nil {
		return
	}
	resp = &v1pb.DeleteCoinRespAdmin{
		Status: RPCResponse.Status,
	}
	return
}

// GetPoolList implementation
// `internal:"true"`
func (s *CapsuleService) GetPoolList(ctx context.Context, req *v1pb.GetPoolListReqAdmin) (resp *v1pb.GetPoolListRespAdmin, err error) {
	var RPCResponse *rspb.GetPoolListResp
	RPCReq := &rspb.GetPoolListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if RPCReq.Page < 1 {
		RPCReq.Page = 1
	}
	if RPCReq.PageSize < 1 {
		RPCReq.PageSize = 16
	}
	RPCResponse, err = s.capsuleCli.GetPoolList(ctx, RPCReq)
	if err != nil {
		return
	}
	poolList := make([]*v1pb.GetPoolListRespAdmin_List, 0)
	for _, singleList := range RPCResponse.List {
		poolInfo := &v1pb.GetPoolListRespAdmin_List{}
		poolInfo.Id = singleList.Id
		startTimeUnix := time.Unix(singleList.StartTime, 0)
		poolInfo.StartTime = startTimeUnix.Format("2006-01-02 15:04:05")
		endTimeUnix := time.Unix(singleList.EndTime, 0)
		poolInfo.EndTime = endTimeUnix.Format("2006-01-02 15:04:05")
		poolInfo.Title = singleList.Title
		poolInfo.CoinTitle = singleList.CoinTitle
		poolInfo.Status = singleList.Status
		poolInfo.Rule = singleList.Rule
		poolInfo.CoinId = singleList.CoinId
		poolList = append(poolList, poolInfo)
	}
	resp = &v1pb.GetPoolListRespAdmin{
		Total:     RPCResponse.Total,
		TotalPage: RPCResponse.TotalPage,
		List:      poolList,
	}
	return
}

// UpdatePool implementation
// `method:"POST" internal:"true"`
func (s *CapsuleService) UpdatePool(ctx context.Context, req *v1pb.UpdatePoolReqAdmin) (resp *v1pb.UpdatePoolRespAdmin, err error) {
	var RPCResponse *rspb.UpdatePoolResp
	loc, _ := time.LoadLocation("Local")
	startTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.StartTime, loc)
	if err != nil {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return nil, err
	}
	endTime, err := time.ParseInLocation("2006-01-02 15:04:05", req.EndTime, loc)
	if err != nil {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return nil, err
	}
	RPCReq := &rspb.UpdatePoolReq{
		Id:        req.Id,
		CoinId:    req.CoinId,
		Title:     req.Title,
		StartTime: startTime.Unix(),
		EndTime:   endTime.Unix(),
		Rule:      req.Rule,
	}
	RPCResponse, err = s.capsuleCli.UpdatePool(ctx, RPCReq)
	if err != nil {
		return
	}
	resp = &v1pb.UpdatePoolRespAdmin{
		Status: RPCResponse.Status,
	}
	return
}

// DeletePool implementation
func (s *CapsuleService) DeletePool(ctx context.Context, req *v1pb.DeletePoolReqAdmin) (resp *v1pb.DeletePoolRespAdmin, err error) {
	var RPCResponse *rspb.DeletePoolResp
	RPCReq := &rspb.DeletePoolReq{
		Id: req.Id,
	}
	RPCResponse, err = s.capsuleCli.DeletePool(ctx, RPCReq)
	if err != nil {
		return
	}
	resp = &v1pb.DeletePoolRespAdmin{
		Status: RPCResponse.Status,
	}

	return
}

// UpdatePoolStatus implementation
func (s *CapsuleService) UpdatePoolStatus(ctx context.Context, req *v1pb.UpdatePoolStatusReqAdmin) (resp *v1pb.UpdatePoolStatusRespAdmin, err error) {
	resp = &v1pb.UpdatePoolStatusRespAdmin{}
	data, err := s.capsuleCli.UpdatePoolStatus(ctx, &rspb.UpdatePoolStatusReq{Id: req.Id, Status: req.Status})
	if err != nil {
		return
	}
	resp.Status = data.Status
	return
}

// GetPoolPrize implementation
// `internal:"true"`
func (s *CapsuleService) GetPoolPrize(ctx context.Context, req *v1pb.GetPoolPrizeReqAdmin) (resp *v1pb.GetPoolPrizeRespAdmin, err error) {
	var RPCResponse *rspb.GetPoolPrizeResp
	RPCReq := &rspb.GetPoolPrizeReq{
		PoolId: req.PoolId,
	}
	RPCResponse, err = s.capsuleCli.GetPoolPrize(ctx, RPCReq)
	if err != nil {
		return
	}
	list := make([]*v1pb.GetPoolPrizeRespAdmin_List, 0)
	for _, poolDetail := range RPCResponse.List {
		singleDetail := &v1pb.GetPoolPrizeRespAdmin_List{}
		singleDetail.Id = poolDetail.Id
		singleDetail.Limit = poolDetail.Limit
		singleDetail.Loop = poolDetail.Loop
		singleDetail.Chance = poolDetail.Chance
		singleDetail.ProType = poolDetail.ProType
		singleDetail.JumpUrl = poolDetail.JumpUrl
		singleDetail.Description = poolDetail.Description
		singleDetail.MobileUrl = poolDetail.MobileUrl
		singleDetail.WebUrl = poolDetail.WebUrl
		singleDetail.Num = poolDetail.Num
		singleDetail.Type = poolDetail.Type
		singleDetail.ObjectId = poolDetail.ObjectId
		singleDetail.Expire = poolDetail.Expire
		singleDetail.PoolId = poolDetail.PoolId
		singleDetail.Name = poolDetail.Name
		singleDetail.Weight = poolDetail.Weight
		singleDetail.WhiteUids = ""
		wlen := len(poolDetail.WhiteUids)
		if wlen > 0 {
			wUids := make([]string, wlen)
			for ix, uid := range poolDetail.WhiteUids {
				wUids[ix] = strconv.FormatInt(uid, 10)
			}
			singleDetail.WhiteUids = strings.Join(wUids, ",")
		}
		list = append(list, singleDetail)
	}
	resp = &v1pb.GetPoolPrizeRespAdmin{
		List: list,
	}
	return
}

// UpdatePoolPrize implementation
// `method:"POST" internal:"true"`
func (s *CapsuleService) UpdatePoolPrize(ctx context.Context, req *v1pb.UpdatePoolPrizeReqAdmin) (resp *v1pb.UpdatePoolPrizeRespAdmin, err error) {
	chance, err := strconv.ParseFloat(req.Chance, 64)
	if err != nil || chance > 1 {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return
	}
	chance = chance*10000 + 0.5
	whiteUids := make([]int64, 0)
	if req.ProType == 1 {
		if chance < 1 {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return
		}
		req.Loop = 0
		req.Limit = 0
	} else if req.ProType == 2 {
		chance = 0
		if req.Loop < 1 || req.Limit < 1 {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return
		}
	} else if req.ProType == 3 {
		chance = 0
		if req.Loop < 1 || req.Limit < 1 {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return
		}
	} else if req.ProType == 4 {
		chance = 0
		req.Loop = 0
		req.Limit = 0
		if req.WhiteUids == "" {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return
		}
		wUids := strings.Split(req.WhiteUids, ",")
		for _, uidStr := range wUids {
			var uid int64
			uid, err = strconv.ParseInt(uidStr, 10, 64)
			if err != nil {
				err = ecode.Error(ecode.InvalidParam, "参数错误")
				return nil, err
			}
			whiteUids = append(whiteUids, uid)
		}
		wlen := len(whiteUids)
		if wlen == 0 || wlen > 1000 {
			err = ecode.Error(ecode.InvalidParam, "参数错误")
			return nil, err
		}
	} else {
		err = ecode.Error(ecode.InvalidParam, "参数错误")
		return
	}
	var RPCResponse *rspb.UpdatePoolPrizeResp
	RPCReq := &rspb.UpdatePoolPrizeReq{
		Id:          req.Id,
		PoolId:      req.PoolId,
		Type:        req.Type,
		Num:         req.Num,
		ObjectId:    req.ObjectId,
		Expire:      req.Expire,
		WebUrl:      req.WebUrl,
		MobileUrl:   req.MobileUrl,
		Description: req.Description,
		JumpUrl:     req.JumpUrl,
		ProType:     req.ProType,
		Chance:      int64(chance),
		Loop:        req.Loop,
		Limit:       req.Limit,
		Weight:      req.Weight,
		WhiteUids:   whiteUids,
	}
	RPCResponse, err = s.capsuleCli.UpdatePoolPrize(ctx, RPCReq)
	if err != nil {
		return
	}
	resp = &v1pb.UpdatePoolPrizeRespAdmin{
		Status:  RPCResponse.Status,
		PrizeId: RPCResponse.PrizeId,
	}
	return
}

// DeletePoolPrize implementation
// `method:"POST" internal:"true" `
func (s *CapsuleService) DeletePoolPrize(ctx context.Context, req *v1pb.DeletePoolPrizeReqAdmin) (resp *v1pb.DeletePoolPrizeRespAdmin, err error) {
	var RPCResponse *rspb.DeletePoolPrizeResp
	RPCReq := &rspb.DeletePoolPrizeReq{
		Id: req.Id,
	}
	RPCResponse, err = s.capsuleCli.DeletePoolPrize(ctx, RPCReq)
	if err != nil {
		return
	}
	resp = &v1pb.DeletePoolPrizeRespAdmin{
		Status: RPCResponse.Status,
	}
	return
}

// GetPrizeType implementation
// `internal:"true"`
func (s *CapsuleService) GetPrizeType(ctx context.Context, req *v1pb.GetPrizeTypeReq) (resp *v1pb.GetPrizeTypeResp, err error) {
	resp = &v1pb.GetPrizeTypeResp{}
	prizeType, err := s.capsuleCli.GetPrizeType(ctx, &rspb.GetPrizeTypeReq{})
	if err != nil || prizeType == nil {
		return
	}
	resp.List = make([]*v1pb.GetPrizeTypeResp_List, len(prizeType.List))
	for ix, data := range prizeType.List {
		resp.List[ix] = &v1pb.GetPrizeTypeResp_List{Type: data.Type, Name: data.Name}
	}
	return
}

// GetPrizeExpire implementation
// `internal:"true"`
func (s *CapsuleService) GetPrizeExpire(ctx context.Context, req *v1pb.GetPrizeExpireReq) (resp *v1pb.GetPrizeExpireResp, err error) {
	resp = &v1pb.GetPrizeExpireResp{}
	prizeExpire, err := s.capsuleCli.GetPrizeExpire(ctx, &rspb.GetPrizeExpireReq{})
	if err != nil || prizeExpire == nil {
		return
	}
	resp.List = make([]*v1pb.GetPrizeExpireResp_List, len(prizeExpire.List))
	for ix, data := range prizeExpire.List {
		resp.List[ix] = &v1pb.GetPrizeExpireResp_List{Expire: data.Expire, Name: data.Name}
	}
	return
}

// FormatErr format error msg
func (s *CapsuleService) FormatErr(statusCode *status.Status) (code int32, msg string) {
	gCode := statusCode.Code()
	code = 1
	if gCode == 2 {
		code, _ := strconv.Atoi(statusCode.Message())

		switch code {
		case -400:
			msg = "参数错误"
		case -401:
			msg = "上线失败，概率设置不为1"
		case -402:
			msg = "保底奖池的上下线时间不可修改"
		case -403:
			msg = "保底奖池不可下线"
		case -404:
			msg = "普通扭蛋币不可更改道具类型"
		case -500:
			msg = "内部错误"
		default:
			msg = "内部错误"
		}
	} else {
		msg = "内部错误"
	}
	return
}

// GetCouponList implementation
// `internal:"true"`
func (s *CapsuleService) GetCouponList(ctx context.Context, req *v1pb.CapsuleGetCouponListReq) (resp *v1pb.CapsuleGetCouponListResp, err error) {
	resp = &v1pb.CapsuleGetCouponListResp{}
	couponList, err := s.capsuleCli.GetCouponList(ctx, &rspb.CapsuleGetCouponListReq{Uid: req.Uid})
	if err != nil || couponList == nil {
		return
	}
	resp.List = make([]*v1pb.CapsuleGetCouponListResp_List, len(couponList.List))
	for ix, conpon := range couponList.List {
		resp.List[ix] = &v1pb.CapsuleGetCouponListResp_List{}
		resp.List[ix].Uid = req.Uid
		resp.List[ix].Status = conpon.Status
		resp.List[ix].AwardTime = conpon.AwardTime
		resp.List[ix].RetryTime = conpon.RetryTime
		resp.List[ix].AwardName = conpon.AwardName
		resp.List[ix].AwardCode = conpon.AwardCode
	}
	return
}
