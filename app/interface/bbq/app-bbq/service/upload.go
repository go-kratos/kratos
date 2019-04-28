package service

import (
	interface_video_v1 "go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/conf"
	"go-common/app/interface/bbq/app-bbq/model"
	service_video_v1 "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"strconv"
	"time"
)

const (
	//FILESIZE .
	FILESIZE = 2 * 1024 * 1024
)

//Upload ...
func (s *Service) Upload(c *bm.Context, mid int64, t int) (location string, err error) {
	var (
		fileName string
		filePath string
		file     []byte
	)
	c.Request.ParseMultipartForm(32 << 20)
	f, h, err := c.Request.FormFile("file")

	//参数判断
	if err != nil {
		//log.Errorv(c, log.KV("event", "service/upload/http/FormFile"), log.KV("err", err))
		err = ecode.FileNotExists
		return
	}
	defer f.Close()
	//文件大小
	if h.Size > FILESIZE {
		err = ecode.FileTooLarge
		return
	}

	tmp := make([]byte, h.Size)
	if _, err = f.Read(tmp); err != nil {
		log.Errorv(c, log.KV("event", "service/upload/ioreader/read"), log.KV("err", err))
		err = ecode.ServerErr
	}

	file = tmp
	fileName = strconv.FormatInt(mid, 10) + "_" + strconv.FormatInt(time.Now().Unix(), 10)
	filePath = "userface"
	switch t {
	case 1:
		filePath = "home_img"
	}

	if location, err = s.dao.Upload(c, fileName, filePath, file); err != nil {
		log.Errorv(c, log.KV("event", "service/upload/dao/upload"), log.KV("err", err))
		err = ecode.ServerErr
	}
	return
}

//PreUpload ..
func (s *Service) PreUpload(c *bm.Context, req *interface_video_v1.PreUploadRequest, mid int64) (rep *service_video_v1.PreUploadResponse, err error) {
	p := &service_video_v1.PreUploadRequest{
		Title:     req.Title,
		Mid:       mid,
		From:      model.FromBBQ,
		FileExt:   req.FileExt,
		Entension: req.Extension,
	}
	if rep, err = s.videoClient.PreUpload(c, p); err != nil {
		log.Errorw(c, "event", "s.videoClient.PreUpload err", "err", err)
		return
	}
	return
}

//CallBack ..
func (s *Service) CallBack(c *bm.Context, req *interface_video_v1.CallBackRequest, mid int64) (rep struct{}, err error) {
	p := &service_video_v1.CallBackRequest{
		Svid: req.Svid,
		Mid:  mid,
	}
	if _, err = s.videoClient.CallBack(c, p); err != nil {
		log.Errorw(c, "event", "s.videoClient.CallBack", "err", err)
		return
	}
	s.dao.MergeUploadReq(c, req.URL, req.UploadID, req.Profile, req.Svid, req.Auth)
	return
}

// VideoUploadCheck 创作中心白名单（估计用不到）
func (s *Service) VideoUploadCheck(c *bm.Context, mid int64) (response interface{}, err error) {
	isAllow := true
	msg := ""
	for _, v := range conf.Filter.MidFilter.White {
		if v == mid {
			isAllow = true
		}
	}

	for _, v := range conf.Filter.MidFilter.Black {
		if v == mid {
			isAllow = false
			msg = "您不在邀请名单之内，请期待我们的邀请"
		}
	}

	return &interface_video_v1.UploadCheckResponse{
		Msg:     msg,
		IsAllow: isAllow,
	}, nil
}

//HomeImg ..
func (s *Service) HomeImg(c *bm.Context, req *interface_video_v1.HomeImgRequest, mid int64) (rep struct{}, err error) {
	p := &service_video_v1.HomeImgCreateRequest{
		Svid:   req.SVID,
		Url:    req.URL,
		Width:  req.Width,
		Mid:    mid,
		Height: req.Height,
	}
	if _, err = s.videoClient.HomeImgCreate(c, p); err != nil {
		return
	}
	return
}
