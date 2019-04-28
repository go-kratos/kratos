package v1

import (
	"context"
	v1pb "go-common/app/interface/live/app-ucenter/api/http/v1"
	"go-common/app/interface/live/app-ucenter/conf"
	"go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// RoomAdminService struct
type RoomAdminService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	// dao *dao.Dao
	conn v1.RoomAdminClient
}

//NewRoomAdminService init
func NewRoomAdminService(c *conf.Config) (s *RoomAdminService) {
	s = &RoomAdminService{
		conf: c,
	}
	conn, err := v1.NewXuserRoomAdminClient(conf.Conf.Warden)
	if err != nil {
		panic(err)
	}
	s.conn = conn
	return s
}

// History 相关服务

// ShowEntry implementation
// 根据登录态获取功能入口是否显示, 需要登录态
// `method:"GET" midware:"auth"`
func (s *RoomAdminService) ShowEntry(ctx context.Context, req *v1pb.ShowEntryReq) (resp *v1pb.ShowEntryResp, err error) {
	resp = &v1pb.ShowEntryResp{}
	mid, _ := metadata.Value(ctx, "mid").(int64)

	if mid <= 0 {
		err = ecode.NoLogin
		return
	}

	if err != nil {
		return
	}

	ret, err := s.conn.IsAny(ctx, &v1.RoomAdminShowEntryReq{
		Uid: mid,
	})
	log.Info("call IsAny mid(%v) ret(%v)", mid, ret)

	resp.HasAdmin = ret.HasAdmin

	return
}

// SearchForAdmin implementation
// 查询需要添加的房管
// `method:"POST" midware:"auth"`
func (s *RoomAdminService) SearchForAdmin(ctx context.Context, req *v1pb.RoomAdminSearchForAdminReq) (resp *v1pb.RoomAdminSearchForAdminResp, err error) {
	resp = &v1pb.RoomAdminSearchForAdminResp{}

	mid, _ := metadata.Value(ctx, "mid").(int64)

	keyWord := req.GetKeyWord()

	if keyWord == "" {
		err = ecode.ParamInvalid
		return
	}
	ret, err := s.conn.SearchForAdmin(ctx, &v1.RoomAdminSearchForAdminReq{
		Uid:     mid,
		KeyWord: keyWord,
	})
	log.Info("call SearchForAdmin mid(%v) keyword (%v) ret(%v)", mid, keyWord, ret)

	if err != nil {
		return
	}

	if ret == nil {
		log.Info("call SearchForAdmin nil mid(%v) keyword (%v) err (%v)", mid, keyWord, err)
		return
	}

	if ret.Data == nil {
		log.Info("SearchForAdmin(%v) return nil (%v)", keyWord, ret)
		return
	}

	for _, v := range ret.Data {
		resp.Data = append(resp.Data, &v1pb.RoomAdminSearchForAdminResp_Data{
			Uid:       v.Uid,
			IsAdmin:   v.IsAdmin,
			Uname:     v.Uname,
			Face:      v.Face,
			MedalName: v.MedalName,
			Level:     v.Level,
		})
	}
	return
}

// IsAny implementation
// 根据登录态获取功能入口是否显示, 需要登录态
// `method:"GET" midware:"auth"`
func (s *RoomAdminService) IsAny(ctx context.Context, req *v1pb.ShowEntryReq) (resp *v1pb.ShowEntryResp, err error) {
	resp = &v1pb.ShowEntryResp{}
	mid, _ := metadata.Value(ctx, "mid").(int64)

	if mid <= 0 {
		err = ecode.NoLogin
		return
	}

	ret, err := s.conn.IsAny(ctx, &v1.RoomAdminShowEntryReq{
		Uid: mid,
	})
	log.Info("call IsAny mid(%v) ret(%v)", mid, ret)

	if err != nil {
		return
	}

	if ret == nil {
		log.Info("call IsAny nil mid(%v) err (%v)", mid, err)
		return
	}

	resp.HasAdmin = ret.HasAdmin
	return
}

// GetByUid implementation
// 获取用户拥有的的所有房管身份
// `method:"GET" midware:"auth"`
func (s *RoomAdminService) GetByUid(ctx context.Context, req *v1pb.RoomAdminGetByUidReq) (resp *v1pb.RoomAdminGetByUidResp, err error) {
	resp = &v1pb.RoomAdminGetByUidResp{}

	mid, _ := metadata.Value(ctx, "mid").(int64)

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	ret, err := s.conn.GetByUid(ctx, &v1.RoomAdminGetByUidReq{
		Uid:  mid,
		Page: page,
	})
	log.Info("call GetByUid mid(%v) page (%v) ret(%v)", mid, page, ret)

	if err != nil {
		return
	}

	if ret == nil {
		log.Info("call GetByUid nil mid(%v) err (%v)", mid, err)
		return
	}

	if nil != ret.Page {
		resp.Page = &v1pb.RoomAdminGetByUidResp_Page{
			Page:       ret.GetPage().GetPage(),
			PageSize:   ret.GetPage().GetPageSize(),
			TotalPage:  ret.GetPage().GetTotalPage(),
			TotalCount: ret.GetPage().GetTotalCount(),
		}
	}

	if nil != ret.Data {
		for _, v := range ret.Data {
			resp.Data = append(resp.Data, &v1pb.RoomAdminGetByUidResp_Data{
				Uid:         v.Uid,
				Roomid:      v.Roomid,
				AnchorId:    v.AnchorId,
				Uname:       v.Uname,
				AnchorCover: v.AnchorCover,
				Ctime:       v.Ctime,
			})
		}
	}

	return
}

// Resign implementation
// 辞职房管
// `method:"POST" midware:"auth"`
func (s *RoomAdminService) Resign(ctx context.Context, req *v1pb.RoomAdminResignRoomAdminReq) (resp *v1pb.RoomAdminResignRoomAdminResp, err error) {
	resp = &v1pb.RoomAdminResignRoomAdminResp{}
	mid, _ := metadata.Value(ctx, "mid").(int64)
	roomId := req.GetRoomid()

	if roomId <= 0 {
		err = ecode.ParamInvalid
		return
	}

	ret, err := s.conn.Resign(ctx, &v1.RoomAdminResignRoomAdminReq{
		Roomid: roomId,
		Uid:    mid,
	})
	log.Info("call Resign mid(%v) room (%v) ret(%v)", mid, roomId, ret)

	if err != nil {
		return
	}

	return
}

// GetByAnchor implementation
// 获取主播拥有的的所有房管身份
// `method:"GET" midware:"auth"`
func (s *RoomAdminService) GetByAnchor(ctx context.Context, req *v1pb.RoomAdminGetByAnchorReq) (resp *v1pb.RoomAdminGetByAnchorResp, err error) {
	resp = &v1pb.RoomAdminGetByAnchorResp{}
	mid, _ := metadata.Value(ctx, "mid").(int64)

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}

	ret, err := s.conn.GetByAnchor(ctx, &v1.RoomAdminGetByAnchorReq{
		Page: page,
		Uid:  mid,
	})
	log.Info("call GetByAnchor mid(%v) page (%v) ret(%v)", mid, page, ret)

	if ret == nil {
		log.Info("call GetByAnchor nil mid(%v) err (%v)", mid, err)
		return
	}

	if err != nil {
		return
	}

	if nil != ret.GetPage() {
		resp.Page = &v1pb.RoomAdminGetByAnchorResp_Page{
			Page:       ret.GetPage().GetPage(),
			PageSize:   ret.GetPage().GetPageSize(),
			TotalPage:  ret.GetPage().GetTotalPage(),
			TotalCount: ret.GetPage().GetTotalCount(),
		}
	}

	if nil != ret.Data {
		for _, v := range ret.Data {
			resp.Data = append(resp.Data, &v1pb.RoomAdminGetByAnchorResp_Data{
				Uid:       v.GetUid(),
				Uname:     v.GetUname(),
				Face:      v.GetFace(),
				Ctime:     v.GetCtime(),
				MedalName: v.GetMedalName(),
				Level:     v.GetLevel(),
			})
		}
	}

	return
}

// Dismiss implementation
// 撤销房管
// `method:"POST" midware:"auth"`
func (s *RoomAdminService) Dismiss(ctx context.Context, req *v1pb.RoomAdminDismissAdminReq) (resp *v1pb.RoomAdminDismissAdminResp, err error) {
	resp = &v1pb.RoomAdminDismissAdminResp{}

	mid, _ := metadata.Value(ctx, "mid").(int64)
	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.ParamInvalid
		return
	}

	ret, err := s.conn.Dismiss(ctx, &v1.RoomAdminDismissAdminReq{
		Uid:      uid,
		AnchorId: mid,
	})

	log.Info("call Dismiss mid(%v) user (%v) ret(%v)", mid, uid, ret)

	if err != nil {
		return
	}

	return
}

// Appoint implementation
// 添加房管
// `method:"POST" midware:"auth"`
func (s *RoomAdminService) Appoint(ctx context.Context, req *v1pb.RoomAdminAddReq) (resp *v1pb.RoomAdminAddResp, err error) {
	resp = &v1pb.RoomAdminAddResp{}
	mid, _ := metadata.Value(ctx, "mid").(int64)

	uid := req.GetUid()
	if uid <= 0 {
		err = ecode.ParamInvalid
		return
	}

	ret, err := s.conn.Appoint(ctx, &v1.RoomAdminAddReq{
		Uid:      uid,
		AnchorId: mid,
	})
	log.Info("call Appoint mid(%v) uid (%v) ret(%v)", mid, uid, ret)

	if err != nil {
		log.Info("Appoint error statusCode(%v) ret(%v), err(%+v)", ret, err)
	}

	if nil != ret {
		resp = &v1pb.RoomAdminAddResp{
			Uid:    ret.GetUid(),
			Roomid: ret.GetRoomid(),
		}
		if nil != ret.Userinfo {
			userInfo := &v1pb.RoomAdminAddResp_UI{
				Uid:   ret.GetUserinfo().GetUid(),
				Uname: ret.GetUserinfo().GetUname(),
			}
			resp.Userinfo = userInfo
		}
	} else {
		log.Info("call appoint return nil uid (%v) anchorid (%v)", uid, mid)
	}

	return
}
