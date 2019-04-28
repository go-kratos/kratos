package dao

import (
	"context"
	"net/http"

	"go-common/library/database/bfs"
	"go-common/library/log"
)

// BfsData get bfs data
func (d *Dao) BfsData(c context.Context, bfsURL string) (bs []byte, err error) {
	var (
		req *http.Request
	)
	if req, err = http.NewRequest(http.MethodGet, bfsURL, nil); err != nil {
		log.Error("NewRequest(bfsURL:%v),error(%v)", bfsURL, err)
		return
	}
	if bs, err = d.httpCli.Raw(c, req); err != nil {
		log.Error("Raw(bfsURL:%v),error(%v)", bfsURL, err)
		return
	}
	return
}

// BfsDmUpload .
func (d *Dao) BfsDmUpload(c context.Context, fileName string, bs []byte) (location string, err error) {
	if location, err = d.bfsCli.Upload(c, &bfs.Request{
		Bucket:      d.conf.Bfs.Dm,
		Filename:    fileName,
		ContentType: "application/json",
		File:        bs,
	}); err != nil {
		log.Error("bfs.BfsDmUpload error(%v)", err)
		return
	}
	return
}
