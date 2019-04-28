package dao

import (
	"context"

	"go-common/library/database/bfs"
	"go-common/library/log"
)

// Upload .
func (d *Dao) Upload(c context.Context, bucket, fileName, contentType string, file []byte) (err error) {
	if _, err = d.bfsCli.Upload(c, &bfs.Request{
		Bucket:      bucket,
		ContentType: contentType,
		Filename:    fileName,
		File:        file,
	}); err != nil {
		log.Error("Upload(err:%v)", err)
	}
	return
}
