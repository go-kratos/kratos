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

// PayLiveService struct
type PayLiveService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
}

//NewPayLiveService init
func NewPayLiveService(c *conf.Config) (s *PayLiveService) {
	s = &PayLiveService{
		conf: c,
	}
	return s
}

// Add implementation
// `method:"POST" internal:"true"` 生成付费直播信息
func (s *PayLiveService) Add(ctx context.Context, req *v1pb.PayLiveAddReq) (resp *v1pb.PayLiveAddResp, err error) {
	resp = &v1pb.PayLiveAddResp{}
	log.Info("Add params:%v", req)
	r, err := dao.AvApi.V0PayLive.Add(ctx, &v0av.PayLiveAddReq{
		Platform:    req.Platform,
		RoomId:      req.RoomId,
		Title:       req.Title,
		Status:      req.Status,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		LiveEndTime: req.LiveEndTime,
		LivePic:     req.LivePic,
		AdPic:       req.AdPic,
		GoodsLink:   req.GoodsLink,
		GoodsId:     req.GoodsId,
		IpLimit:     req.IpLimit,
		BuyGoodsId:  req.BuyGoodsId,
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
// `method:"POST" internal:"true"` 更新付费直播信息
func (s *PayLiveService) Update(ctx context.Context, req *v1pb.PayLiveUpdateReq) (resp *v1pb.PayLiveUpdateResp, err error) {
	resp = &v1pb.PayLiveUpdateResp{}
	log.Info("Update params:%v", req)
	r, err := dao.AvApi.V0PayLive.Update(ctx, &v0av.PayLiveUpdateReq{
		LiveId:      req.LiveId,
		Platform:    req.Platform,
		RoomId:      req.RoomId,
		Title:       req.Title,
		Status:      req.Status,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		LiveEndTime: req.LiveEndTime,
		LivePic:     req.LivePic,
		AdPic:       req.AdPic,
		GoodsLink:   req.GoodsLink,
		GoodsId:     req.GoodsId,
		IpLimit:     req.IpLimit,
		BuyGoodsId:  req.BuyGoodsId,
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
// `method:"POST" internal:"true"` 获取付费直播列表
func (s *PayLiveService) GetList(ctx context.Context, req *v1pb.PayLiveGetListReq) (resp *v1pb.PayLiveGetListResp, err error) {
	resp = &v1pb.PayLiveGetListResp{}
	r, err := dao.AvApi.V0PayLive.GetList(ctx, &v0av.PayLiveGetListReq{
		RoomId:   req.RoomId,
		Title:    req.Title,
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
	resp.PageInfo = &v1pb.PayLiveGetListResp_PageInfo{
		TotalCount: data.PageInfo.TotalCount,
		PageNum:    data.PageInfo.PageNum,
	}
	for _, v := range r.Data.GoodsInfo {
		tmp := &v1pb.PayLiveGetListResp_GoodsInfo{
			LiveId:        v.LiveId,
			Platform:      v.Platform,
			RoomId:        v.RoomId,
			Title:         v.Title,
			Status:        v.Status,
			PayLiveStatus: v.PayLiveStatus,
			StartTime:     v.StartTime,
			EndTime:       v.EndTime,
			LiveEndTime:   v.LiveEndTime,
			LivePic:       v.LivePic,
			AdPic:         v.AdPic,
			GoodsLink:     v.GoodsLink,
			GoodsId:       v.GoodsId,
			IpLimit:       v.IpLimit,
			BuyGoodsId:    v.BuyGoodsId,
		}
		resp.GoodsInfo = append(resp.GoodsInfo, tmp)
	}
	return
}

// Close implementation
// `method:"POST" internal:"true"` 关闭鉴权
func (s *PayLiveService) Close(ctx context.Context, req *v1pb.PayLiveCloseReq) (resp *v1pb.PayLiveCloseResp, err error) {
	resp = &v1pb.PayLiveCloseResp{}
	r, err := dao.AvApi.V0PayLive.Close(ctx, &v0av.PayLiveCloseReq{
		LiveId: req.LiveId,
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
// `method:"POST" internal:"true"` 开启鉴权
func (s *PayLiveService) Open(ctx context.Context, req *v1pb.PayLiveOpenReq) (resp *v1pb.PayLiveOpenResp, err error) {
	resp = &v1pb.PayLiveOpenResp{}
	r, err := dao.AvApi.V0PayLive.Open(ctx, &v0av.PayLiveOpenReq{
		LiveId: req.LiveId,
	})
	if err != nil {
		log.Error("call av error,err:%v", err)
		return
	}
	if err != nil || r.Code != 0 {
		log.Error("call av error,code:%v,msg:%v", r.Code, r.Msg)
		err = ecode.Error(ecode.Int(int(r.Code)), r.Msg)
		return
	}
	return
}
