package http

import (
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
)

func saveFiles(c *bm.Context) {
	c.JSON(nil, srv.SaveFiles(c))
}

func downloadStoryFile(c *bm.Context) {
	var (
		err  error
		data []byte
		code int
	)
	if data, err = srv.DownloadStoryFile(c); err != nil {
		log.Error("Download story file failed, error:%v", err)
		code = -1
	}
	contentType := " text/plain;charset:utf-8;"
	c.Writer.Header().Set("content-disposition", `attachment; filename=story.txt`)
	c.Bytes(code, contentType, data)
}

func downloadChangeFile(c *bm.Context) {
	var (
		err  error
		data []byte
		code int
	)
	if data, err = srv.DownloadChangeFile(c); err != nil {
		log.Error("Download change file failed, error:%v", err)
		code = -1
	}
	contentType := " text/plain;charset:utf-8;"
	c.Writer.Header().Set("content-disposition", `attachment; filename=change.txt`)
	c.Bytes(code, contentType, data)
}

func downloadIterationFile(c *bm.Context) {
	var (
		err  error
		data []byte
		code int
	)
	if data, err = srv.DownloadIterationFile(c); err != nil {
		log.Error("Download iteration file failed, error:%v", err)
		code = -1
	}
	contentType := " text/plain;charset:utf-8;"
	c.Writer.Header().Set("content-disposition", `attachment; filename=iteration.txt`)
	c.Bytes(code, contentType, data)
}

func downloadBugFile(c *bm.Context) {
	var (
		err  error
		data []byte
		code int
	)
	if data, err = srv.DownBugFile(c); err != nil {
		log.Error("Download bug file failed, error:%v", err)
		code = -1
	}
	contentType := " text/plain;charset:utf-8;"
	c.Writer.Header().Set("content-disposition", `attachment; filename=bug.txt`)
	c.Bytes(code, contentType, data)
}
