package v1

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	v1pb "go-common/app/admin/live/live-admin/api/http/v1"
	"go-common/app/admin/live/live-admin/conf"
	"go-common/app/admin/live/live-admin/dao"
	"go-common/app/admin/live/live-admin/model"
	"go-common/library/database/bfs"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
)

// UploadService struct
type UploadService struct {
	conf *conf.Config
	// optionally add other properties here, such as dao
	dao *dao.Dao
}

//NewUploadService init
func NewUploadService(c *conf.Config, d *dao.Dao) (s *UploadService) {
	s = &UploadService{
		conf: c,
		dao:  d,
	}
	return s
}

// File implementation
// `method:"POST" content-type:"multipart/form-data" midware:"guest"`
func (s *UploadService) File(ctx context.Context, req *v1pb.UploadFileReq) (resp *v1pb.UploadFileResp, err error) {
	_, ok := s.conf.Bucket[req.Bucket]
	if !ok {
		err = ecode.UploadBucketErr
		return
	}

	if !s.verify(ctx, req.Token) {
		err = ecode.Error(ecode.InvalidParam, "invalid upload token")
		return
	}

	bmc := ctx.(*blademaster.Context)

	file, _, err := bmc.Request.FormFile("file_up")
	if err != nil {
		log.Error("Parse file part failure: %v", err)
		err = ecode.Error(ecode.InvalidParam, "file not found")
		return
	}

	defer file.Close()

	fileData := new(bytes.Buffer)
	if _, err = io.Copy(fileData, file); err != nil {
		log.Error("Read file data failure: %v", err)
		err = ecode.UploadUploadErr
		return
	}

	resp = &v1pb.UploadFileResp{}

	bfsClient := bfs.New(nil)

	bfsReq := &bfs.Request{
		Bucket:      req.Bucket,
		Dir:         req.Dir,
		ContentType: req.ContentType,
		Filename:    req.Filename,
		File:        fileData.Bytes(),
		WMKey:       req.WmKey,
		WMText:      req.WmText,
		WMPaddingX:  req.WmPaddingX,
		WMPaddingY:  req.WmPaddingY,
		WMScale:     req.WmScale,
	}

	model.TweakWatermark(bfsReq)

	resp.Url, err = bfsClient.Upload(ctx, bfsReq)
	if err != nil {
		log.Error("Upload to bfs failure: %v", err)
		err = ecode.UploadUploadErr
		return
	}

	resp.Url = s.getPrivateURL(req.Bucket, resp.Url)

	return
}

func (s *UploadService) verify(ctx context.Context, token string) bool {
	return s.dao.VerifyUploadToken(ctx, token)
}

func (s *UploadService) getPrivateURL(bucket, fileURL string) string {
	bucketConf, ok := s.conf.Bucket[bucket]
	if !ok || !bucketConf.Private {
		return fileURL
	}
	bucketFindStr := fmt.Sprintf("/%s/", bucket)
	bucketIndex := strings.Index(fileURL, bucketFindStr)
	if bucketIndex < 0 {
		return fileURL
	}
	filename := fileURL[bucketIndex+len(bucketFindStr):]

	now := time.Now().Unix()
	mac := hmac.New(sha1.New, []byte(bucketConf.Secret))
	salt := fmt.Sprintf("GET\n%s\n%s\n%d\n", bucket, filename, now)
	mac.Write([]byte(salt))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	token := fmt.Sprintf("%s:%s:%d", bucketConf.Key, sign, now)
	v := url.Values{}
	v.Set("token", token)
	return fmt.Sprintf("%s?%s", fileURL, v.Encode())

}
