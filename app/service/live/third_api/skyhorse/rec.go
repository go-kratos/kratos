package skyhorse

import (
	"context"
	"net/url"
	"strconv"

	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	_skyHorseURLAddr = "http://data.bilibili.co/recommand"
	_skyHorseCnt     = 6
	_skyHorseFrom    = "29004"
)

type SkyHorseItem struct {
	Tid        int    `json:"tid"`
	Id         int    `json:"id"`
	GotoType   string `json:"goto_type"`
	Source     string `json:"source"`
	TrackId    string `json:"track_id"`
	AvFeature  string `json:"av_feature"`
	RcmdReason string `json:"rcmd_reason"`
}

type skyHorseRecResp struct {
	Code  int             `json:"code"`
	Error string          `json:"error"`
	Data  []*SkyHorseItem `json:"data"`
}

func (c *Client) GetSkyHorseRec(ctx context.Context, mid int64, buvid string, build int64, plat string,
	duplicateItem []int64, strongLen int, timeout int) (skyHorseRec *skyHorseRecResp, err error) {
	cli := bm.NewClient(c.getConf())

	param := url.Values{}
	requestCnt := _skyHorseCnt - strongLen

	if requestCnt <= 0 {
		requestCnt = 6
	}
	param.Set("cmd", "live")
	param.Set("from", _skyHorseFrom)
	param.Set("request_cnt", strconv.Itoa(requestCnt))
	param.Set("mid", strconv.FormatInt(mid, 10))
	param.Set("buvid", buvid)
	param.Set("build", strconv.FormatInt(build, 10))
	param.Set("plat", plat)
	param.Set("timeout", strconv.Itoa(timeout))
	param.Set("duplicates", xstr.JoinInts(duplicateItem))

	skyHorseRec = &skyHorseRecResp{}
	err = cli.Get(ctx, _skyHorseURLAddr, "", param, skyHorseRec)

	if err != nil {
		log.Error("[GetSkyHorseRec]error:%+v=", err)
		return
	}

	if skyHorseRec.Code != ecode.OK.Code() {
		err = ecode.Int(skyHorseRec.Code)
		log.Error("[getSkyHorseRoomList]getSkyHorseList error:%+v,code:%d,msg:%s", err, skyHorseRec.Code, skyHorseRec.Error)
	}

	return
}
