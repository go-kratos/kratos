package dao

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"go-common/library/database/bfs"
	"go-common/library/log"
)

const BUCKET = "live"
const FILE_PATH = "/data/www/cover/"

//上传至bfs相关接口

//ImgUpload  上传图片至bfs
func (d *Dao) ImgUpload(ctx context.Context, roomId int64, pic string, file []byte) (resp string, err error) {
	log.Info("ImgUpload_start")
	fileName := strconv.Itoa(int(roomId)) + ".jpg"
	reply, err := d.BfsClient.Upload(ctx, &bfs.Request{
		Bucket:      BUCKET,
		Dir:         "",
		ContentType: "",
		Filename:    fileName,
		File:        []byte(file),
	})
	if err != nil {
		log.Error("ImgUpload_bfs_Upload_failed,err:%v", err)
		return
	}
	resp = reply
	return
}

func (d *Dao) ImgDownload(ctx context.Context, pic string) (resp []byte, err error) {
	reply, err := http.Get(pic)
	if err != nil {
		log.Warn("ImgDownload_failed_err:%v", err)
		return
	}
	defer reply.Body.Close()
	if reply.StatusCode != 200 {
		err = errors.New("curl error http code not equal to 200")
		log.Warn("ImgDownload_failed_httpCode:%d", reply.StatusCode)
		return
	}
	resp, err = ioutil.ReadAll(reply.Body)
	if err != nil {
		log.Warn("ImgDownload_read_err:%v", err)
		return
	}
	return
}
