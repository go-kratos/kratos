package bvc

import (
	"context"
	"fmt"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"
	"math"
	"net/url"
	"strconv"
)

const livePlayURLHost = "live-playurl.bilibili.co"
const livePlayURLAddr = "http://live-playurl.bilibili.co/room/v1/Room/playUrl"
const livePlayURLMultiAddr = "http://%s/stream/v1/multiplayurl"

type PlayUrlItem struct {
	RoomId         int64             `json:"room_id"`
	CurrentQuality int64             `json:"current_quality"`
	AcceptQuality  []int64           `json:"accept_quality"`
	Url            map[string]string `json:"url"`
}
type respPlayUrl struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data []*PlayUrlItem `json:"data"`
}

// 历史原因，单个的接口和视频云@yuanfeng确认数据结构不统一，后面统一，暂时先不用
// @see http://live-hub.bilibili.co/doc/playurl/
func (c *Client) GetPlayUrl(ctx *bm.Context, cid, quality, ptype, needHttps, unicomFree int64, platform string) (resp *respPlayUrl) {
	cli := bm.NewClient(c.getConf())

	param := url.Values{}

	param.Set("cid", strconv.FormatInt(cid, 10))
	param.Set("quality", strconv.FormatInt(quality, 10))
	param.Set("platform", platform)
	param.Set("ptype", strconv.FormatInt(ptype, 10))
	param.Set("https_url_req", strconv.FormatInt(needHttps, 10))
	param.Set("unicom_free", strconv.FormatInt(unicomFree, 10))

	err := cli.Get(ctx, livePlayURLAddr, "", param, &resp)
	if err != nil {
		log.Error("call %s error(%v)", livePlayURLAddr, err)
		return
	}
	if resp.Code != ecode.OK.Code() {
		log.Error("call %s response error, code(%d) message(%s)", livePlayURLAddr, resp.Code, resp.Msg)
		return
	}
	return
}

// 批量接口，一次最好20个
func (c *Client) GetPlayUrlMulti(ctx context.Context, roomIds []int64,
	needHttps, quality, build int64, platform string) (result map[int64]*PlayUrlItem) {
	cli := bm.NewClient(c.getConf())
	result = make(map[int64]*PlayUrlItem)

	userIp := metadata.String(ctx, metadata.RemoteIP)
	if userIp == "127.0.0.1" || c.mock == "1" {
		userIp = "222.73.196.18"
	}
	bvcHost := c.getBvcHost(livePlayURLHost)

	if quality == 0 {
		quality = 4
	}

	chunkSize := 20
	lens := len(roomIds)
	if lens <= 0 {
		return
	}
	// 批次
	chunkNum := int(math.Ceil(float64(lens) / float64(chunkSize)))
	wg, _ := errgroup.WithContext(ctx)
	// 空间换时间
	chunkResult := make([][]*PlayUrlItem, chunkNum)
	for i := 1; i <= chunkNum; i++ {
		x := i
		wg.Go(func() error {
			param := url.Values{}
			chunkRoomIds := make([]int64, 20)
			if x == chunkNum {
				chunkRoomIds = roomIds[(x-1)*chunkSize:]
			} else {
				chunkRoomIds = roomIds[(x-1)*chunkSize : x*chunkSize]
			}
			param.Set("room_ids", xstr.JoinInts(chunkRoomIds))
			param.Set("platform", platform)
			param.Set("https_url_req", strconv.FormatInt(needHttps, 10))
			param.Set("quality", strconv.FormatInt(quality, 10))
			param.Set("build", strconv.FormatInt(build, 10))
			param.Set("req_biz", "live-app-interface")

			var resp respPlayUrl
			bvcUrl := fmt.Sprintf(livePlayURLMultiAddr, bvcHost)
			req, err := cli.NewRequest("GET", bvcUrl, userIp, param)
			if err != nil {
				log.Error("GetPlayUrlMulti client.NewRequest: get error(%v)", err)
				return nil
			}
			//req.Header.Add("x-backend-bili-real-ip", userIp)
			if err := cli.Do(ctx, req, &resp); err != nil {
				log.Error("GetPlayUrlMulti client.Do: get error(%v)", err)
				return nil
			}
			if resp.Code != 0 {
				log.Error("GetPlayUrlMulti call %s response error, code(%d) message(%s) param(%+v) room_ids(%+v)", bvcUrl, resp.Code, resp.Msg, param, chunkRoomIds)
				return nil
			}
			if resp.Data == nil || len(resp.Data) <= 0 {
				log.Warn("GetPlayUrlMulti call %s response empty, code(%d) message(%s) param(%+v) room_ids(%+v)", bvcUrl, resp.Code, resp.Msg, param, chunkRoomIds)
				return nil
			}

			//chunkResult = append(chunkResult, resp.Data)
			chunkResult[x-1] = resp.Data

			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		log.Error("GetPlayUrlMulti get playurl wait error(%+v)", err)
		return nil
	}
	// 整理数据
	for _, chunkItemList := range chunkResult {
		for _, item := range chunkItemList {
			if item != nil {
				result[item.RoomId] = item
			}
		}
	}

	return
}
