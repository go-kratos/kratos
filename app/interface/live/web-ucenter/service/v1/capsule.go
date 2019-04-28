package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc/status"

	"go-common/library/log"

	"github.com/pkg/errors"

	v1pb "go-common/app/interface/live/web-ucenter/api/http/v1"
	"go-common/app/interface/live/web-ucenter/conf"
	capsuledao "go-common/app/interface/live/web-ucenter/dao/capsule"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// CapsuleService struct
type CapsuleService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao   *capsuledao.Dao
	infoc *infoc.Infoc
}

//NewCapsuleService init
func NewCapsuleService(c *conf.Config) (s *CapsuleService) {
	s = &CapsuleService{
		conf: c,
		dao:  capsuledao.New(c),
	}
	if c.Infoc != nil && c.Infoc.CapsuleInfoc != nil {
		s.infoc = infoc.New(c.Infoc.CapsuleInfoc)
	}
	return s
}

// GetDetail implementation
func (s *CapsuleService) GetDetail(ctx context.Context, req *v1pb.CapsuleGetDetailReq) (resp *v1pb.CapsuleGetDetailResp, err error) {
	resp = &v1pb.CapsuleGetDetailResp{}
	uid, ok := metadata.Value(ctx, metadata.Mid).(int64)
	if !ok {
		err = errors.Wrap(err, "未取到uid")
		return
	}
	data, err := s.dao.GetDetail(ctx, uid, req.From)
	if err != nil {
		return
	}
	if data.Normal != nil {
		normal := &v1pb.CapsuleGetDetailResp_CapsuleInfo{}
		normal.Status = data.Normal.Status
		normal.Coin = data.Normal.Coin
		normal.Change = data.Normal.Change
		normal.Rule = data.Normal.Rule
		if data.Normal.Progress != nil {
			normal.Progress = &v1pb.Progress{}
			normal.Progress.Now = data.Normal.Progress.Now
			normal.Progress.Max = data.Normal.Progress.Max
		}
		glen := len(data.Normal.Gift)
		normal.Gift = make([]*v1pb.CapsuleGetDetailResp_Gift, glen)
		for i := 0; i < glen; i++ {
			gift := &v1pb.CapsuleGetDetailResp_Gift{}
			gift.Name = data.Normal.Gift[i].Name
			gift.WebImage = data.Normal.Gift[i].WebImage
			gift.MobileImage = data.Normal.Gift[i].MobileImage
			gift.Image = data.Normal.Gift[i].Image
			if data.Normal.Gift[i].Usage != nil {
				gift.Usage = &v1pb.Usage{}
				gift.Usage.Text = data.Normal.Gift[i].Usage.Text
				gift.Usage.Url = data.Normal.Gift[i].Usage.Url
			}
			normal.Gift[i] = gift
		}
		glen = len(data.Normal.List)
		normal.List = make([]*v1pb.CapsuleGetDetailResp_List, glen)
		for i := 0; i < glen; i++ {
			info := &v1pb.CapsuleGetDetailResp_List{}
			info.Name = data.Normal.List[i].Name
			info.Num = data.Normal.List[i].Num
			info.Date = data.Normal.List[i].Date
			info.Gift = data.Normal.List[i].Gift
			normal.List[i] = info
		}
		resp.Normal = normal
	}
	if data.Colorful != nil {
		colorful := &v1pb.CapsuleGetDetailResp_CapsuleInfo{}
		colorful.Status = data.Colorful.Status
		colorful.Coin = data.Colorful.Coin
		colorful.Change = data.Colorful.Change
		colorful.Rule = data.Colorful.Rule
		if data.Colorful.Progress != nil {
			colorful.Progress = &v1pb.Progress{}
			colorful.Progress.Now = data.Colorful.Progress.Now
			colorful.Progress.Max = data.Colorful.Progress.Max
		}
		glen := len(data.Colorful.Gift)
		colorful.Gift = make([]*v1pb.CapsuleGetDetailResp_Gift, glen)
		for i := 0; i < glen; i++ {
			gift := &v1pb.CapsuleGetDetailResp_Gift{}
			gift.Name = data.Colorful.Gift[i].Name
			gift.WebImage = data.Colorful.Gift[i].WebImage
			gift.MobileImage = data.Colorful.Gift[i].MobileImage
			gift.Image = data.Colorful.Gift[i].Image
			if data.Colorful.Gift[i].Usage != nil {
				gift.Usage = &v1pb.Usage{}
				gift.Usage.Text = data.Colorful.Gift[i].Usage.Text
				gift.Usage.Url = data.Colorful.Gift[i].Usage.Url
			}
			colorful.Gift[i] = gift
		}
		glen = len(data.Colorful.List)
		colorful.List = make([]*v1pb.CapsuleGetDetailResp_List, glen)
		for i := 0; i < glen; i++ {
			info := &v1pb.CapsuleGetDetailResp_List{}
			info.Name = data.Colorful.List[i].Name
			info.Num = data.Colorful.List[i].Num
			info.Date = data.Colorful.List[i].Date
			info.Gift = data.Colorful.List[i].Gift
			colorful.List[i] = info
		}
		resp.Colorful = colorful
	}
	return
}

// OpenCapsule implementation
func (s *CapsuleService) OpenCapsule(ctx context.Context, req *v1pb.CapsuleOpenCapsuleReq) (resp *v1pb.CapsuleOpenCapsuleResp, err error) {
	resp = &v1pb.CapsuleOpenCapsuleResp{}
	uid, ok := metadata.Value(ctx, metadata.Mid).(int64)
	if !ok {
		err = errors.Wrap(err, "未取到uid")
		return
	}
	data, err := s.dao.OpenCapsule(ctx, uid, req.Type, req.Count, req.Platform)
	if err != nil {
		return
	}
	if data.Info != nil {
		resp.Info = &v1pb.CapsuleOpenCapsuleResp_Info{}
		if data.Info.Colorful != nil {
			resp.Info.Colorful = &v1pb.CapsuleOpenCapsuleResp_CapsuleInfo{}
			resp.Info.Colorful.Change = data.Info.Colorful.Change
			resp.Info.Colorful.Coin = data.Info.Colorful.Coin
			if data.Info.Colorful.Progress != nil {
				resp.Info.Colorful.Progress = &v1pb.Progress{}
				resp.Info.Colorful.Progress.Max = data.Info.Colorful.Progress.Max
				resp.Info.Colorful.Progress.Now = data.Info.Colorful.Progress.Now
			}
		}
		if data.Info.Normal != nil {
			resp.Info.Normal = &v1pb.CapsuleOpenCapsuleResp_CapsuleInfo{}
			resp.Info.Normal.Change = data.Info.Normal.Change
			resp.Info.Normal.Coin = data.Info.Normal.Coin
			if data.Info.Normal.Progress != nil {
				resp.Info.Normal.Progress = &v1pb.Progress{}
				resp.Info.Normal.Progress.Max = data.Info.Normal.Progress.Max
				resp.Info.Normal.Progress.Now = data.Info.Normal.Progress.Now
			}
		}
	}
	resp.Status = data.Status
	resp.IsEntity = data.IsEntity
	resp.ShowTitle = data.ShowTitle
	resp.Text = data.Text
	l := len(data.Awards)
	resp.Awards = make([]*v1pb.CapsuleOpenCapsuleResp_Award, l)
	for i := 0; i < l; i++ {
		resp.Awards[i] = &v1pb.CapsuleOpenCapsuleResp_Award{}
		resp.Awards[i].Num = data.Awards[i].Num
		resp.Awards[i].Name = data.Awards[i].Name
		resp.Awards[i].Text = data.Awards[i].Text
		resp.Awards[i].MobileImage = data.Awards[i].MobileImage
		resp.Awards[i].WebImage = data.Awards[i].WebImage
		resp.Awards[i].Img = data.Awards[i].Img
		if data.Awards[i].Usage != nil {
			resp.Awards[i].Usage = &v1pb.Usage{}
			resp.Awards[i].Usage.Text = data.Awards[i].Usage.Text
			resp.Awards[i].Usage.Url = data.Awards[i].Usage.Url
		}
	}
	if s.infoc != nil {
		awards, err1 := json.Marshal(resp.Awards)
		if err1 != nil {
			log.Error("OpenCapsule OpenCapsule err")
			return
		}
		bmc, ok := ctx.(bm.Context)
		if !ok {
			return
		}
		header := bmc.Request.Header
		userip := header.Get("x-backend-bili-real-ip")
		if userip == "" {
			userip = bmc.Request.RemoteAddr
		}
		s.infoc.Infov(context.Background(), uid, userip, strconv.FormatInt(time.Now().Unix(), 10), awards, header.Get("platform"), header.Get("version"), header.Get("buvid"), bmc.Request.UserAgent(), bmc.Request.Referer())
	}
	return
}

// FormatErr format error msg
func (s *CapsuleService) FormatErr(statusCode *status.Status) (code int32, msg string) {
	gCode := statusCode.Code()
	fmt.Printf("FormatErr %d %s", gCode, statusCode.Message())
	code = 1
	if gCode == 2 {
		code, _ := strconv.Atoi(statusCode.Message())

		switch code {
		case -400:
			msg = "参数错误"
		case -401:
			msg = "扭蛋币不足"
		case -500:
			msg = "系统繁忙，请稍后再试"
		case -501:
			msg = "系统繁忙，请稍后再试"
		default:
			msg = "系统繁忙，请稍后再试"
		}
	} else {
		msg = "系统繁忙，请稍后再试"
	}
	return
}

// GetCapsuleInfo implementation
// `midware:"guest"`
func (s *CapsuleService) GetCapsuleInfo(ctx context.Context, req *v1pb.CapsuleGetCapsuleInfoReq) (resp *v1pb.CapsuleGetCapsuleInfoResp, err error) {
	resp = &v1pb.CapsuleGetCapsuleInfoResp{}
	uid, _ := metadata.Value(ctx, metadata.Mid).(int64)
	data, err := s.dao.GetCapsuleInfo(ctx, uid, req.Type, req.From)
	if err != nil || data == nil {
		return
	}
	resp.Coin = data.Coin
	resp.Rule = data.Rule
	if len(data.GiftFilter) > 0 {
		resp.GiftFilter = make([]*v1pb.CapsuleGetCapsuleInfoResp_GiftFilter, len(data.GiftFilter))
		for ix, award := range data.GiftFilter {
			resp.GiftFilter[ix] = &v1pb.CapsuleGetCapsuleInfoResp_GiftFilter{}
			resp.GiftFilter[ix].Id = award.Id
			resp.GiftFilter[ix].Name = award.Name
			resp.GiftFilter[ix].WebUrl = award.WebUrl
			resp.GiftFilter[ix].MobileUrl = award.MobileUrl
			if award.Usage != nil {
				resp.GiftFilter[ix].Usage = &v1pb.Usage{}
				resp.GiftFilter[ix].Usage.Text = award.Usage.Text
				resp.GiftFilter[ix].Usage.Url = award.Usage.Url
			}
		}
	}
	if len(data.GiftList) > 0 {
		resp.GiftList = make([]*v1pb.CapsuleGetCapsuleInfoResp_GiftList, len(data.GiftList))
		for ix, award := range data.GiftList {
			resp.GiftList[ix] = &v1pb.CapsuleGetCapsuleInfoResp_GiftList{}
			resp.GiftList[ix].Id = award.Id
			resp.GiftList[ix].Name = award.Name
			resp.GiftList[ix].Num = award.Num
			resp.GiftList[ix].WebUrl = award.WebUrl
			resp.GiftList[ix].MobileUrl = award.MobileUrl
			resp.GiftList[ix].Type = award.Type
			resp.GiftList[ix].Expire = award.Expire
			if award.Usage != nil {
				resp.GiftList[ix].Usage = &v1pb.Usage{}
				resp.GiftList[ix].Usage.Text = award.Usage.Text
				resp.GiftList[ix].Usage.Url = award.Usage.Url
			}
		}
	}
	return
}

// OpenCapsuleByType implementation
// `method:"POST" midware:"auth"`
func (s *CapsuleService) OpenCapsuleByType(ctx context.Context, req *v1pb.CapsuleOpenCapsuleByTypeReq) (resp *v1pb.CapsuleOpenCapsuleByTypeResp, err error) {
	resp = &v1pb.CapsuleOpenCapsuleByTypeResp{}
	uid, ok := metadata.Value(ctx, metadata.Mid).(int64)
	if !ok {
		err = errors.Wrap(err, "未取到uid")
		return
	}
	data, err := s.dao.OpenCapsuleByType(ctx, uid, req.Type, req.Count, req.Platform)
	if err != nil || data == nil {
		return
	}
	resp.IsEntity = data.IsEntity
	resp.Status = data.Status
	resp.Text = data.Text
	resp.Info = &v1pb.CapsuleOpenCapsuleByTypeResp_CapsuleInfo{Coin: 0}
	if data.Info != nil {
		resp.Info.Coin = data.Info.Coin
	}
	resp.Awards = make([]*v1pb.CapsuleOpenCapsuleByTypeResp_Award, len(data.Awards))
	for ix, award := range data.Awards {
		resp.Awards[ix] = &v1pb.CapsuleOpenCapsuleByTypeResp_Award{}
		resp.Awards[ix].Id = award.Id
		resp.Awards[ix].Name = award.Name
		resp.Awards[ix].Num = award.Num
		resp.Awards[ix].Text = award.Text
		resp.Awards[ix].WebUrl = award.WebUrl
		resp.Awards[ix].MobileUrl = award.MobileUrl
		resp.Awards[ix].Type = award.Type
		resp.Awards[ix].Expire = award.Expire
		if award.Usage != nil {
			resp.Awards[ix].Usage = &v1pb.Usage{}
			resp.Awards[ix].Usage.Text = award.Usage.Text
			resp.Awards[ix].Usage.Url = award.Usage.Url
		}
	}
	return
}
