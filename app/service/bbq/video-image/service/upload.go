package service

import (
	"context"
	"go-common/app/service/bbq/video-image/api/grpc/v1"
	"go-common/library/database/bfs"
	"go-common/library/log"
)

const (
	//BUCKET ...
	BUCKET = "bbq"
)

//ImgUpload ...
func (s *Service) ImgUpload(ctx context.Context, req *v1.ImgUploadRequest) (rep *v1.ImgUploadResponse, err error) {
	log.Info("begin imgupload")
	rep = &v1.ImgUploadResponse{}
	conf := &bfs.Config{
		Host:       s.c.URL["api"],
		HTTPClient: s.c.BM.Client,
	}
	bfsClient := bfs.New(conf)

	log.Info("bfs begin......")
	bfsreq := &bfs.Request{
		Bucket:   BUCKET,
		Dir:      "/video-image/" + req.Dir,
		Filename: req.Filename,
		File:     req.File,
	}
	rep.Location, err = bfsClient.Upload(context.Background(), bfsreq)
	if err != nil {
		log.Error("bfs failed,err:%v", err)
	}
	return
}
