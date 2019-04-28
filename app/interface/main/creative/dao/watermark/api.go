package watermark

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/tool"
	"go-common/app/interface/main/creative/model/watermark"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_genWm = "/x/internal/image/gen"
	// upload_key	true	string	上传的key(需要拿到分配到key,确定上传的bucket)
	// uploadKey = "creative"
	// wm_key	   true	string	水印key，新业务需要提前向bfs申请（未上线)
	wmKey = "creative"
	// wm_text	   true	string	水印文字，限制不超过20个字符
	// wm_vertical true	bool	水印图片文字，是否水平排列
	wmVertical = "true"
	// wm_scale    true float	水印文字宽度和水印图片宽度比值(0,1]，默认1
	wmScale = 1
	//水印图片到文字的距离. 默认是1
	distance = 1
)

// GenWm set wm by api.
func (d *Dao) GenWm(c context.Context, mid int64, wmText, ip string) (gm *watermark.GenWatermark, err error) {
	params := url.Values{}
	params.Set("wm_key", wmKey)
	params.Set("wm_text", wmText)
	params.Set("wm_vertical", wmVertical)
	params.Set("wm_scale", strconv.Itoa(wmScale))
	params.Set("distance", strconv.Itoa(distance))
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("appkey", conf.Conf.App.Key)
	params.Set("appsecret", conf.Conf.App.Secret)
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	var (
		req *http.Request
		res struct {
			Code int                    `json:"code"`
			Data watermark.GenWatermark `json:"data"`
		}
		query, _ = tool.Sign(params)
		genWmURL = d.genWm + "?" + query
	)
	log.Info("genWm url(%v)", genWmURL)
	if req, err = http.NewRequest("POST", genWmURL, nil); err != nil {
		log.Error("genWm url(%s) error(%v)", genWmURL, err)
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", ip)
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("genWm url(%s) response(%v) error(%v)", genWmURL, res, err)
		return
	}
	if res.Code != 0 {
		log.Error("genWm url(%s) res(%v)", genWmURL, res)
		err = ecode.Int(res.Code)
		return
	}
	gm = &watermark.GenWatermark{
		Location: res.Data.Location,
		MD5:      res.Data.MD5,
		Width:    res.Data.Width,
		Height:   res.Data.Height,
	}
	return
}
