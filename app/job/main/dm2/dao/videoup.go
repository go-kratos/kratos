package dao

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/job/main/dm2/model"
	"go-common/library/log"
)

const (
	_archiveURL = "/videoup/view"
)

func (d *Dao) archiveURI() string {
	return d.conf.Host.Videoup + _archiveURL
}

// Videos 根据aid获取分批信息，包含未开放浏览分批.
func (d *Dao) Videos(c context.Context, aid int64) (videos []*model.Video, err error) {
	var (
		uri    = d.archiveURI()
		params = url.Values{}
		res    struct {
			Code    int64  `json:"code"`
			Message string `json:"message"`
			Data    *struct {
				Archive *model.Archive `json:"archive"`
				Videos  []*model.Video `json:"videos"`
			} `json:"data"`
		}
	)
	params.Set("aid", fmt.Sprint(aid))
	if err = d.httpCli.Get(c, uri, "", params, &res); err != nil {
		return
	}
	if res.Code != 0 {
		err = fmt.Errorf("aid:%d,res.Code:%d", aid, res.Code)
		log.Error("url(%s) aid(%d) res(%v)", uri+"?"+params.Encode(), aid, res)
		return
	}
	if res.Data == nil || res.Data.Archive == nil {
		log.Error("url(%s) aid(%d) res(%v)", uri+"?"+params.Encode(), aid, res)
		return
	}
	if len(res.Data.Videos) == 0 {
		log.Error("url(%s) aid(%d) videos is empty", uri+"?"+params.Encode(), aid)
		return
	}
	for _, v := range res.Data.Videos {
		v.Mid = res.Data.Archive.Mid
		videos = append(videos, v)
	}
	return
}
