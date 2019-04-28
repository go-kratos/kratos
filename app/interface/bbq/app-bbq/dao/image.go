package dao

import (
	"context"
	image "go-common/app/service/bbq/video-image/api/grpc/v1"
	"go-common/library/log"
)

//Upload .
func (d *Dao) Upload(c context.Context, fileName string, filePath string, file []byte) (location string, err error) {
	imageReq := &image.ImgUploadRequest{
		Filename: fileName,
		Dir:      filePath,
		File:     file,
	}
	imageRes, err := d.imageClient.ImgUpload(c, imageReq)
	if err != nil {
		log.Errorv(c, log.KV("event", "grpc/imageupload"), log.KV("err", err))
		return
	}
	if imageRes != nil {
		location = imageRes.Location
	}
	return
}
