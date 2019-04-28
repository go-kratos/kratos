package service

import (
	"context"
	"crypto/md5"
	"fmt"
	topic_v1 "go-common/app/service/bbq/topic/api"
	video_v1 "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/app/service/bbq/video/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"io"
	"net/url"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
)

//PreUpload ...
func (s *Service) PreUpload(c context.Context, req *video_v1.PreUploadRequest) (rep *video_v1.PreUploadResponse, err error) {
	var (
		p *video_v1.CreateIDResponse
		m *model.VideoRepository
	)
	r := &video_v1.CreateIDRequest{
		Mid: req.Mid,
	}
	if p, err = s.CreateID(c, r); err != nil {
		log.Errorw(c, "event", "CreateID err:", "err", err)
		return
	}
	//new topic
	tp := &topic_v1.VideoExtension{
		Svid:      p.NewId,
		Extension: req.Entension,
	}
	if _, err = s.topicClient.Register(c, tp); err != nil {
		log.Errorw(c, "event", "s.topicClient.Register err", "err", err)
		return
	}
	m = &model.VideoRepository{
		SVID:       p.NewId,
		From:       req.From,
		Title:      req.Title,
		MID:        req.Mid,
		SyncStatus: model.SourceRequest,
	}
	if err = s.dao.InsertVR(c, m); err != nil {
		log.Warnw(c, "event", "s.dao.InsertVR err", "err", err)
		return
	}
	fn := s.getFileName(p.NewId)
	rep = &video_v1.PreUploadResponse{
		Svid:      p.NewId,
		UposUri:   "upos://bbq/" + fn + "." + req.FileExt,
		EndPoint:  s.c.Upload.Endpoint.Main,
		EndPoints: []string{s.c.Upload.Endpoint.Main, s.c.Upload.Endpoint.BackUp},
		Auth:      s.getAuth(fn, req.FileExt),
	}
	return
}

func (s *Service) getAuth(fileName string, fileExt string) (auth string) {

	var (
		url    = url.Values{}
		h      = md5.New()
		md5key string
		t      string
		sign   string
	)
	t = strconv.FormatInt(time.Now().Unix(), 10)
	md5key = s.c.Upload.Auth.AK + t + "/bbq/" + fileName + "." + fileExt + s.c.Upload.Auth.SK
	io.WriteString(h, md5key)
	sign = fmt.Sprintf("%x", h.Sum(nil))

	url.Add("ak", s.c.Upload.Auth.AK)
	url.Add("timestamp", t)
	url.Add("sign", sign)

	return url.Encode()
}

func (s *Service) getFileName(svid int64) (fn string) {
	fn = fmt.Sprintf("%s%6s%2s%23s",
		s.c.Upload.File.Prefix,
		time.Now().Format("060102"),
		s.c.Upload.File.Line,
		"svid"+strconv.FormatInt(svid, 10),
	)
	return
}

//CallBack upload call back function
func (s *Service) CallBack(c context.Context, req *video_v1.CallBackRequest) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	var (
		pvr  *model.VideoRepository
		pvup *model.VideoUploadProcess
		rvr  *model.VideoRepository
	)
	pvr = &model.VideoRepository{
		SVID:       req.Svid,
		SyncStatus: model.SourceRequest,
	}
	if rvr, err = s.dao.QueryVR(c, pvr); err != nil {
		log.Warnw(c, "event", "s.dao.QueryVR err", "err", err)
		return
	}
	//callback legality verification
	if req.Mid != rvr.MID {
		err = ecode.UploadFailed
		return
	}
	pvup = &model.VideoUploadProcess{
		SVID:          req.Svid,
		Title:         rvr.Title,
		Mid:           rvr.MID,
		UploadStatus:  model.UploadStatusWaiting,
		RetryTimes:    0,
		HomeImgURL:    rvr.HomeImgURL,
		HomeImgWidth:  rvr.HomeImgWidth,
		HomeImgHeight: rvr.HomeImgHeight,
	}
	if err = s.dao.InsertOrUpdateVUP(c, pvup); err != nil {
		log.Warnw(c, "event", "s.dao.InsertOrUpdateVUP", "err", err)
		return
	}
	//ignore update result
	s.dao.UpdateVR(c, pvr)
	return
}

// ListPrepareVideo 获取prepare视频
func (s *Service) ListPrepareVideo(c context.Context, req *video_v1.PrepareVideoRequest) (res *video_v1.PrepareVideoResponse, err error) {
	res = new(video_v1.PrepareVideoResponse)
	list, err := s.dao.GetPrepareVUP(c, req.Mid)
	if err != nil {
		log.Warnw(c, "log", "get prepare vup fail", "mid", req.Mid)
		return
	}
	res.List = list
	return
}

//HomeImgCreate ..
func (s *Service) HomeImgCreate(c context.Context, req *video_v1.HomeImgCreateRequest) (res *empty.Empty, err error) {
	res = new(empty.Empty)
	p := &model.VideoRepository{
		SVID:          req.Svid,
		MID:           req.Mid,
		HomeImgHeight: req.Height,
		HomeImgURL:    req.Url,
		HomeImgWidth:  req.Width,
	}
	if err = s.dao.HomeImgCreate(c, p); err != nil {
		return
	}
	return
}
