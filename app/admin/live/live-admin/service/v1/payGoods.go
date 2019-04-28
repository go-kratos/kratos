package v1

import (
	"context"
	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
	"go-common/app/admin/live/live-admin/dao"
	v0av "go-common/app/service/live/av/api/liverpc/v0"
	"go-common/library/ecode"
	"go-common/library/log"
)

// PayGoodsService struct
type PayGoodsService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewPayGoodsService init
func NewPayGoodsService(c *conf.Config) (s *PayGoodsService) {
	s = &PayGoodsService{
		conf: c,
	}
	return s
}

// Add implementation
// * 生成一张付费直播票
func (s *PayGoodsService) Add(ctx context.Context, req *v1pb.PayGoodsAddReq) (resp *v1pb.PayGoodsAddResp, err error) {
	resp = &v1pb.PayGoodsAddResp{}
	log.Info("Add params:%v", req)
	r, err := dao.AvApi.V0PayGoods.Add(ctx, &v0av.PayGoodsAddReq{
		Platform:  req.Platform,
		Title:     req.Title,
		Type:      req.Type,
		Price:     req.Price,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IpLimit:   req.IpLimit,
	})
	if err != nil {
		log.Error("call av error,err:%v", err)
		return
	}
	if r.Code != 0 {
		log.Error("call av error,code:%v,msg:%v", r.Code, r.Msg)
		err = ecode.Error(ecode.Int(int(r.Code)), r.Msg)
		return
	}
	return
}

// Update implementation
// * 更新一张付费直播票
func (s *PayGoodsService) Update(ctx context.Context, req *v1pb.PayGoodsUpdateReq) (resp *v1pb.PayGoodsUpdateResp, err error) {
	resp = &v1pb.PayGoodsUpdateResp{}
	log.Info("Update params:%v", req)
	r, err := dao.AvApi.V0PayGoods.Update(ctx, &v0av.PayGoodsUpdateReq{
		Id:        req.Id,
		Platform:  req.Platform,
		Title:     req.Title,
		Type:      req.Type,
		Price:     req.Price,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IpLimit:   req.IpLimit,
	})
	if err != nil {
		log.Error("call av error,err:%v", err)
		return
	}
	if r.Code != 0 {
		log.Error("call av error,code:%v,msg:%v", r.Code, r.Msg)
		err = ecode.Error(ecode.Int(int(r.Code)), r.Msg)
		return
	}
	return
}

// GetList implementation
// * 获取付费直播票列表
func (s *PayGoodsService) GetList(ctx context.Context, req *v1pb.PayGoodsGetListReq) (resp *v1pb.PayGoodsGetListResp, err error) {
	resp = &v1pb.PayGoodsGetListResp{}
	r, err := dao.AvApi.V0PayGoods.GetList(ctx, &v0av.PayGoodsGetListReq{
		Id:       req.Id,
		Platform: req.Platform,
		Title:    req.Title,
		Type:     req.Type,
		IpLimit:  req.IpLimit,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		log.Error("call av error,err:%v", err)
		return
	}
	if r.Code != 0 {
		log.Error("call av error,code:%v,msg:%v", r.Code, r.Msg)
		err = ecode.Error(ecode.Int(int(r.Code)), r.Msg)
		return
	}
	data := r.Data
	resp.PageInfo = &v1pb.PayGoodsGetListResp_PageInfo{
		TotalCount: data.PageInfo.TotalCount,
		PageNum:    data.PageInfo.PageNum,
	}
	for _, v := range r.Data.GoodsInfo {
		tmp := &v1pb.PayGoodsGetListResp_GoodsInfo{
			Id:        v.Id,
			Title:     v.Title,
			Platform:  v.Platform,
			Type:      v.Type,
			Price:     v.Price,
			StartTime: v.StartTime,
			EndTime:   v.EndTime,
			IpLimit:   v.IpLimit,
			Status:    v.Status,
		}
		resp.GoodsInfo = append(resp.GoodsInfo, tmp)
	}
	return
}

// Close implementation
// * 关闭购票
func (s *PayGoodsService) Close(ctx context.Context, req *v1pb.PayGoodsCloseReq) (resp *v1pb.PayGoodsCloseResp, err error) {
	resp = &v1pb.PayGoodsCloseResp{}
	r, err := dao.AvApi.V0PayGoods.Close(ctx, &v0av.PayGoodsCloseReq{
		Id: req.Id,
	})
	if err != nil {
		log.Error("call av error,err:%v", err)
		return
	}
	if r.Code != 0 {
		log.Error("call av error,code:%v,msg:%v", r.Code, r.Msg)
		err = ecode.Error(ecode.Int(int(r.Code)), r.Msg)
		return
	}
	return
}

// Open implementation
// * 开启购票
func (s *PayGoodsService) Open(ctx context.Context, req *v1pb.PayGoodsOpenReq) (resp *v1pb.PayGoodsOpenResp, err error) {
	resp = &v1pb.PayGoodsOpenResp{}
	r, err := dao.AvApi.V0PayGoods.Open(ctx, &v0av.PayGoodsOpenReq{
		Id: req.Id,
	})
	if err != nil {
		log.Error("call av error,err:%v", err)
		return
	}
	if r.Code != 0 {
		log.Error("call av error,code:%v,msg:%v", r.Code, r.Msg)
		err = ecode.Error(ecode.Int(int(r.Code)), r.Msg)
		return
	}
	return
}
