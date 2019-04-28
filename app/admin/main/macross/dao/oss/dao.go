package oss

import (
	"context"
	"io"
	"path"

	"go-common/app/admin/main/macross/conf"
	"go-common/library/log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// Dao macross dao
type Dao struct {
	// conf
	c      *conf.Config
	client *oss.Client
}

// New dao.
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		c: c,
	}
	client, err := oss.New(d.c.Oss.Endpoint, d.c.Oss.AccessKeyID, d.c.Oss.AccessKeySecret)
	if err != nil {
		panic(err)
	}
	d.client = client
	return
}

// Put put object into oss.
func (d *Dao) Put(c context.Context, rd io.Reader, apkName string) (uri string, err error) {
	bucket, _ := d.client.Bucket(d.c.Oss.Bucket)
	uri = path.Join(d.c.Oss.OriginDir, apkName)
	// put
	if err = bucket.PutObject(uri, rd); err != nil {
		log.Error("bucket.PutObject(%s) error(%v)", uri, err)
		return
	}
	uri = path.Join(d.c.Oss.Bucket, uri) // NOTE: begin with '/' when upload success
	return
}

// Publish publish object into oss.
func (d *Dao) Publish(objectKey string, destKey string) (uri string, err error) {
	bucket, _ := d.client.Bucket(d.c.Oss.Bucket)
	_, err = bucket.CopyObject(d.c.Oss.PublishDir+objectKey, d.c.Oss.PublishDir+destKey)
	if err != nil {
		log.Error("bucket.CopyObject(%s, %s) error(%v)", d.c.Oss.PublishDir+objectKey, d.c.Oss.PublishDir+destKey, err)
		return
	}
	uri = path.Join("/", d.c.Oss.PublishDir+destKey)
	return
}

// Close close kafka connection.
func (d *Dao) Close() {
}
