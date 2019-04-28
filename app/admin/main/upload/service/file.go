package service

import (
	"context"
	"net/http"
	"time"

	"go-common/app/admin/main/upload/model"

	"github.com/pkg/errors"
)

// UploadAdminRecord upload file to bfs
func (s *Service) UploadAdminRecord(ctx context.Context, action string, up *model.UploadParam, data []byte) (result *model.UploadResult, err error) {
	var (
		location, etag string
		b              *model.Bucket
		ok             bool
	)
	if b, ok = s.bucketCache[up.Bucket]; !ok {
		err = errors.Errorf("read bucket items failed: (%s)", up.Bucket)
		return
	}
	// auth calc.
	up.Auth = s.dao.Bfs.Authorize(b.KeyID, b.KeySecret, http.MethodPut, up.Bucket, up.FileName, time.Now().Unix())
	if up.ContentType == "" {
		up.ContentType = http.DetectContentType(data)
	}
	if location, etag, err = s.dao.Bfs.Upload(ctx, up, data); err != nil {
		return
	}
	result = &model.UploadResult{
		Location: location,
		Etag:     etag,
	}
	return
}
