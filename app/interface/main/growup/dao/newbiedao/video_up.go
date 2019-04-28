package newbiedao

import (
	"context"
	"github.com/pkg/errors"
	"go-common/app/interface/main/growup/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"net/url"
	"strconv"
)

// GetVideoUp get video up
func (d *Dao) GetVideoUp(c context.Context, aid int64) (videoUpArchive *model.VideoUpArchive, err error) {
	uv := url.Values{}
	uv.Set("aid", strconv.FormatInt(aid, 10))
	videoUpRes := new(model.VideoUpRes)
	err = d.httpRead.Get(c, d.c.Host.VideoUpURI, metadata.String(c, metadata.RemoteIP), uv, videoUpRes)
	if err != nil {
		log.Error("s.dao.GetVideoUp error(%v)", err)
		return
	}
	if videoUpRes.Code != ecode.OK.Code() {
		err = errors.Wrap(ecode.Int(videoUpRes.Code), "get video up failed")
		log.Error("s.dao.GetVideoUp get video up failed, ecode: %d", videoUpRes.Code)
		return
	}
	if videoUpRes.Data == nil || videoUpRes.Data.Archive == nil {
		err = errors.Wrap(ecode.Int(videoUpRes.Code), "get video up nil")
		log.Error("s.dao.GetVideoUp get video up is empty")
		return
	}

	videoUpArchive = videoUpRes.Data.Archive
	return
}
