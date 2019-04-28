package http

import (
	"strconv"
	"strings"

	v1 "go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/app/interface/bbq/app-bbq/model"
	video "go-common/app/service/bbq/video/api/grpc/v1"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// @params SvListReq
// @router get /bbq/app-bbq/sv/list/
// @response VideoResponse
func svList(c *bm.Context) {
	b, _ := c.Get("BBQBase")
	mid, _ := c.Get("mid")
	arg := &v1.SvListReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	arg.Base = b.(*v1.Base)
	if mid != nil {
		arg.MID = mid.(int64)
	} else {
		arg.MID = 0
	}
	//获取deviceID
	deviceID := c.Request.Form.Get("device_id")
	log.Info("sv list Context [%+v]", c.Request.Header)
	log.Info("sv list Base [%+v]", arg.Base)
	arg.RemoteIP = metadata.String(c, metadata.RemoteIP)
	resp, err := srv.SvList(c, arg.PageSize, arg.MID, arg.Base, deviceID)
	c.JSON(resp, err)

	// 埋点
	if err != nil {
		return
	}
	var svidList []int64
	for _, v := range resp {
		svidList = append(svidList, v.SVID)
	}
	ext := struct {
		Svids []int64 `json:"svid_list"`
	}{
		Svids: svidList,
	}
	uiLog(c, model.ActionRecommend, ext)
}

func svDetail(c *bm.Context) {
	mid := int64(0)
	if res, ok := c.Get("mid"); ok {
		mid = res.(int64)
	}
	arg := &v1.SvDetailReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	c.JSON(srv.SvDetail(c, arg.SVID, mid))
}

//svStatistics 视频互动数据
func svStatistics(c *bm.Context) {
	arg := new(model.ParamStatistic)
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	ids := strings.Split(arg.SVIDs, ",")
	if len(ids) == 0 {
		err := ecode.RequestErr
		errors.Wrap(err, "svid解析为空")
		return
	}
	var svids []int64
	for _, v := range ids {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			errors.Wrap(err, "svid解析错误")
			return
		}
		svids = append(svids, id)
	}
	var mid int64
	if res, ok := c.Get("mid"); ok {
		mid = res.(int64)
	}
	c.JSON(srv.SvStatistics(c, mid, svids))
}

func svPlayList(c *bm.Context) {
	dev, _ := c.Get("device")
	mid, _ := c.Get("mid")
	arg := &v1.PlayListReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	arg.Device = dev.(*bm.Device)
	if mid != nil {
		arg.MID = mid.(int64)
	} else {
		arg.MID = 0
	}
	arg.RemoteIP = metadata.String(c, metadata.RemoteIP)

	cidStr := strings.Split(arg.CIDs, ",")
	var cids []int64
	for _, v := range cidStr {
		cid, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			errors.Wrap(err, "cid解析错误")
			return
		}
		cids = append(cids, cid)
	}
	c.JSON(srv.SvCPlays(c, cids, arg.MID))
}

func svRelList(c *bm.Context) {
	b, _ := c.Get("BBQBase")
	mid, _ := c.Get("mid")
	arg := &v1.SvRelReq{}
	if err := c.Bind(arg); err != nil {
		errors.Wrap(err, "参数验证失败")
		return
	}
	base := b.(*v1.Base)
	if mid != nil {
		arg.MID = mid.(int64)
	} else {
		arg.MID = 0
	}
	arg.BUVID = base.BUVID
	arg.APP = base.App
	arg.APPVersion = base.Version
	arg.QueryID = base.QueryID
	arg.Limit = 15
	arg.Offset = 0
	c.JSON(srv.SvRelRec(c, arg))
}

func svDel(c *bm.Context) {
	arg := new(video.VideoDeleteRequest)
	if err := c.Bind(arg); err != nil {
		return
	}

	if mid, _ := c.Get("mid"); mid != nil {
		arg.UpMid = mid.(int64)
	} else {
		arg.UpMid = 0
	}

	c.JSON(srv.SvDel(c, arg))
}
