package music

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/music"
	"go-common/library/ecode"

	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_audioListURI = "/x/internal/v1/audio/songs/batch"
)

// Audio fn
func (d *Dao) Audio(c context.Context, ids []int64, level int, ip string) (au map[int64]*music.Audio, err error) {
	params := url.Values{}
	params.Set("ids", xstr.JoinInts(ids))
	params.Set("level", strconv.Itoa(level)) //0、只返回基本信息 1、会增加up主名称、播放数、评论数
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		res struct {
			Code int                    `json:"code"`
			Msg  string                 `json:"msg"`
			Data map[int64]*music.Audio `json:"data"`
		}
		query, _ = tool.Sign(params)
		url      string
	)
	url = d.audioListURL + "?" + query
	log.Info("music query url(%s)", url)
	// new requests
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error("http.NewRequest(%s) error(%v)|ip(%s)", url, err, ip)
		err = ecode.CreativeElecErr
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s)|error(%v)|ip(%s)", url, err, ip)
		err = ecode.CreativeMusicErr
		return
	}
	if res.Code != 0 {
		log.Error("music url(%s)|res(%v)|ip(%s)|code(%d)", url, res, ip, res.Code)
		err = ecode.CreativeMusicErr
		return
	}
	au = res.Data
	return
}
