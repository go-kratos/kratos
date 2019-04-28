package service

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/upload/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/queue/databus/report"
)

const (
	_genImageContentType = "image/png"
)

// GenImageUpload generate watermark image by text and upload it.
func (s *Service) GenImageUpload(ctx context.Context, uploadKey string, wmKey, wmText string, distance int, vertical bool) (res *model.ResultWm, err error) {
	var image []byte
	res = new(model.ResultWm)
	key, secret, bucket := s.authorizeInfo(uploadKey)
	if image, res.Height, res.Width, res.Md5, err = s.bfs.GenImage(ctx, wmKey, wmText, distance, vertical); err != nil {
		return
	}
	res.Location, _, err = s.bfs.Upload(ctx, key, secret, _genImageContentType, bucket, "", "", image)
	return
}

// Upload upload by key and secret.
func (s *Service) Upload(ctx context.Context, uploadKey, uploadToken, contentType string, data []byte) (result *model.Result, err error) {
	if !s.verify(uploadKey, uploadToken) {
		err = ecode.AccessDenied
		return
	}
	key, secret, bucket := s.authorizeInfo(uploadKey)
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}
	location, etag, err := s.bfs.Upload(ctx, key, secret, contentType, bucket, "", "", data)
	if err != nil {
		return
	}
	result = &model.Result{
		Location: location,
		Etag:     etag,
	}
	return
}

// authorizeInfo get authorize info by upload key.
func (s *Service) authorizeInfo(uploadKey string) (key, secret, bucket string) {
	key = s.c.BfsBucket.Key
	secret = s.c.BfsBucket.Sercet
	bucket = s.c.BfsBucket.Bucket
	for _, a := range s.c.Auths {
		if a.AppKey == uploadKey && a.BfsBucket != nil {
			key = a.BfsBucket.Key
			secret = a.BfsBucket.Sercet
			bucket = a.BfsBucket.Bucket
			break
		}
	}
	return
}

func (s *Service) verify(key, token string) bool {
	var (
		expire, now, delta int64
		err                error
	)
	for _, auth := range s.c.Auths {
		if key == auth.AppKey {
			srcs := strings.Split(token, ":")
			if len(srcs) != 2 {
				log.Error("verify error: len(srcs) != 2")
				return false
			}
			if expire, err = strconv.ParseInt(srcs[1], 10, 64); err != nil {
				log.Error("verify error: expire not int (%v)", err)
				return false
			}
			if s.gen(auth.AppKey, auth.AppSercet, expire) != srcs[0] {
				log.Error("verify error: s.gen(%s,%s,%d) != srcs[0]:%s", auth.AppKey, auth.AppSercet, expire, srcs[0])
				return false
			}
			now = time.Now().Unix()
			// > ±15 min is forbidden
			if expire > now {
				delta = expire - now
			} else {
				delta = now - expire
			}
			if delta > 900 {
				log.Error("verify error: delta > 900")
				return false
			}
			return true
		}
	}
	return false
}

func (s *Service) gen(key, sercet string, now int64) string {
	sha1 := sha1.New()
	sha1.Write([]byte(fmt.Sprintf("i love bilibili %s:%d", sercet, now)))
	return hex.EncodeToString(sha1.Sum([]byte("")))
}

// UploadRecord .
func (s *Service) UploadRecord(ctx context.Context, action model.UploadActionType, mid int64, up *model.UploadParam, data []byte) (result *model.Result, err error) {
	var (
		location string
		etag     string
		b        *model.Bucket
		ok       bool
	)
	if b, ok = s.bucketCache[up.Bucket]; !ok {
		err = ecode.BfsUplaodBucketNotExist
		log.Error("read bucket items failed: (%s)", up.Bucket)
		return
	}
	// content-type
	if up.ContentType == "" {
		up.ContentType = http.DetectContentType(data)
	}
	// check limit if dir not null
	if up.Dir != "" {
		up.Dir = strings.Trim(up.Dir, "/")
		//todo: dir limit if conf exist
		if err = s.dirlimit(up.Bucket, up.Dir, up.ContentType, data); err != nil {
			return
		}
	}
	log.Info("upload record params:%+v", up)
	if up.WmKey != "" || up.WmText != "" {
		var image []byte
		if image, err = s.bfs.Watermark(ctx, data, up.ContentType, up.WmKey, up.WmText, up.WmPaddingX, up.WmPaddingY, up.WmScale); err != nil {
			log.Error("Upload.Watermark data length(%d) params(%+v) error(%v)", len(data), up, err)
		} else {
			data = image
		}
	}
	if location, etag, err = s.bfs.Upload(ctx, b.KeyID, b.KeySecret, up.ContentType, up.Bucket, up.Dir, up.FileName, data); err != nil {
		return
	}
	uri, err := url.Parse(location)
	if err != nil {
		return
	}
	//todo: add user report
	// http://info.bilibili.co/pages/viewpage.action?pageId=8474335
	uInfo := &report.UserInfo{
		Mid:      mid,
		Platform: "bfs-upload-interface",
		Build:    1,
		Business: 61, //bfs-upload
		Action:   action.String(),
		Ctime:    time.Now(),
		Index:    []interface{}{location, up.Bucket, up.Dir, uri.Path}, // path is filename
		Content: map[string]interface{}{
			"upload_param": up,
		},
	}
	report.User(uInfo)

	result = &model.Result{
		Location: location,
		Etag:     etag,
	}
	return
}

// UploadAdminRecord no dir limit upload method.
func (s *Service) UploadAdminRecord(ctx context.Context, action model.UploadActionType, up *model.UploadParam, data []byte) (result *model.Result, err error) {
	var (
		location, etag string
		b              *model.Bucket
		ok             bool
	)
	if b, ok = s.bucketCache[up.Bucket]; !ok {
		err = ecode.BfsUplaodBucketNotExist
		log.Error("read bucket items failed: (%s)", up.Bucket)
		return
	}
	if up.ContentType == "" {
		up.ContentType = http.DetectContentType(data)
	}
	if location, etag, err = s.bfs.Upload(ctx, b.KeyID, b.KeySecret, up.ContentType, up.Bucket, up.Dir, up.FileName, data); err != nil {
		return
	}
	uri, err := url.Parse(location)
	if err != nil {
		return
	}
	//todo: add user report
	// http://info.bilibili.co/pages/viewpage.action?pageId=8474335
	uInfo := &report.UserInfo{
		Mid:      0,
		Platform: "bfs-upload-interface",
		Build:    1,
		Business: 61,              //bfs-upload
		Action:   action.String(), //操作类型
		Ctime:    time.Now(),
		Index:    []interface{}{location, up.Bucket, up.Dir, uri.Path}, // path is filename in bfs
		Content: map[string]interface{}{
			"upload_param": up,
		},
	}
	report.User(uInfo)

	result = &model.Result{
		Location: location,
		Etag:     etag,
	}
	return
}

func (s *Service) dirlimit(bucket, dir, contentType string, data []byte) (err error) {
	var (
		width    int
		height   int
		dirlimit *model.DirConfig
		ok       bool
	)
	dir = strings.Trim(dir, "/")
	if dirlimit, ok = s.bucketCache[bucket].DirLimit[dir]; !ok {
		return
	}
	if len(dirlimit.Pic.AllowTypeSlice) != 0 {
		var isAllow bool
		for _, ctype := range dirlimit.Pic.AllowTypeSlice {
			if ctype == contentType {
				isAllow = true
			}
		}
		if !isAllow {
			log.Error("image content type illegal,bucket(%s),dir(%s),content type(%s)", bucket, dir, contentType)
			err = ecode.BfsUploadFileContentTypeIllegal
			return
		}
	}
	if dirlimit.Pic.FileSize > 0 && len(data) > dirlimit.Pic.FileSize {
		log.Error("data is too large, bucket(%s), dir(%s)", bucket, dir)
		err = ecode.FileTooLarge
		return
	}
	dataReader := bytes.NewReader(data)
	if width, height, err = Pixel(dataReader); err != nil {
		log.Error("get pixal error(%v), content-type maybe not support, bucket(%s), dir(%s)", err, bucket, dir)
		err = nil
		return
	}
	if (dirlimit.Pic.MinPixelWidthSize != 0 && width < dirlimit.Pic.MinPixelWidthSize) || (dirlimit.Pic.MaxPixelWidthSize != 0 && width > dirlimit.Pic.MaxPixelWidthSize) {
		log.Error("image width illegal, bucket(%s), dir(%s)", bucket, dir)
		err = ecode.BfsUploadFilePixelWidthIllegal
		return
	}
	if (dirlimit.Pic.MinPixelWidthSize != 0 && height < dirlimit.Pic.MinPixelHeightSize) || (dirlimit.Pic.MaxPixelWidthSize != 0 && height > dirlimit.Pic.MaxPixelHeightSize) {
		log.Error("image height illegal, bucket(%s), dir(%s)", bucket, dir)
		err = ecode.BfsUploadFilePixelHeightIllegal
		return
	}
	if dirlimit.Pic.MinAspectRatio != 0 && float64(width/height) < dirlimit.Pic.MinAspectRatio {
		log.Error("image MinAspectRatio illegal, bucket(%s), dir(%s)", bucket, dir)
		err = ecode.BfsUploadFilePixelAspectRatioIllegal
		return
	}
	if dirlimit.Pic.MaxAspectRatio != 0 && float64(width/height) > dirlimit.Pic.MaxAspectRatio {
		log.Error("image MaxAspectRatio illegal, bucket(%s), dir(%s)", bucket, dir)
		err = ecode.BfsUploadFilePixelAspectRatioIllegal
		return
	}
	return
}
