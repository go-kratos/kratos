package http

import (
	"io/ioutil"
	"net/http"

	grpc "go-common/app/service/bbq/video-image/api/grpc/v1"
	"go-common/app/service/bbq/video-image/api/http/v1"
	"go-common/app/service/bbq/video-image/conf"
	"go-common/app/service/bbq/video-image/service"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"

	"github.com/pkg/errors"
)

var (
	srv *service.Service
)

// Init init
func Init(c *conf.Config, s *service.Service) {
	srv = s
	engine := bm.DefaultServer(c.BM.Server)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("bm Start error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	e.Register(register)
	g := e.Group("/bbq/internal/image")
	{
		g.POST("/upload", imageUpload)
		g.POST("/video_cover/score", videoCoverScore)
		g.POST("/upload/v2", upload)
	}
}

func ping(c *bm.Context) {
	if err := srv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(map[string]interface{}{}, nil)
}

func videoCoverScore(c *bm.Context) {
	args := &v1.ScoreRequest{}
	if err := c.Bind(args); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.Request.ParseMultipartForm(1 << 22)
	f, h, err := c.Request.FormFile("file")
	//参数判断
	if err != nil {
		log.Errorv(c, log.KV("log", "get file error"), log.KV("err", err))
		err = ecode.FileNotExists
		return
	}
	defer f.Close()
	//文件大小
	if h.Size > (1 << 22) {
		log.Errorv(c, log.KV("log", "get file error"), log.KV("err", "file size too large"))
		err = ecode.FileTooLarge
		return
	}

	c.JSON(srv.VideoCoverScore(c, args.Name, f))
}

//internal upload
func upload(c *bm.Context) {
	req := &grpc.ImgUploadRequest{}
	req.Filename = c.Request.PostFormValue("filename")
	req.Dir = c.Request.PostFormValue("dir")
	req.File = []byte(c.Request.PostFormValue("file"))
	c.JSON(srv.ImgUpload(c, req))
}
func imageUpload(c *bm.Context) {
	req := &grpc.ImgUploadRequest{}
	req.Filename = c.Request.PostFormValue("filename")
	req.Dir = c.Request.PostFormValue("dir")

	c.Request.ParseMultipartForm(1 << 22)
	f, h, err := c.Request.FormFile("file")
	//参数判断
	if err != nil {
		log.Errorv(c, log.KV("log", "get file error"), log.KV("err", err))
		err = ecode.FileNotExists
		return
	}
	defer f.Close()
	//文件大小
	if h.Size > (1 << 22) {
		log.Errorv(c, log.KV("log", "get file error"), log.KV("err", "file size too large"))
		err = ecode.FileTooLarge
		return
	}
	req.File, err = ioutil.ReadAll(f)
	// req.File = []byte(c.Request.PostFormValue("file"))
	//参数判断
	if err != nil {
		log.Errorv(c, log.KV("log", "get file error"), log.KV("err", err))
		err = ecode.FileNotExists
		return
	}

	c.JSON(srv.ImgUpload(c, req))
}
