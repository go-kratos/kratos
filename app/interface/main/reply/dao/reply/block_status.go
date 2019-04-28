package reply

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/interface/main/reply/conf"
	"go-common/library/log"
	httpx "go-common/library/net/http/blademaster"
)

// BlockStatusDao BlockStatusDao
type BlockStatusDao struct {
	requestURL string
	httpClient *httpx.Client
}

// NewBlockStatusDao NewBlockStatusDao
func NewBlockStatusDao(c *conf.Config) *BlockStatusDao {
	d := &BlockStatusDao{
		httpClient: httpx.NewClient(c.HTTPClient),
		requestURL: c.Reply.BlockStatusURL,
	}
	return d
}

// BlockInfo BlockInfo
type BlockInfo struct {
	// block_status 变成了 PassTest,原因不可描述
	PassTest     int   `json:"block_status"` // 1: fail 0: succ
	ForeverBlock bool  `json:"blocked_forever"`
	BlockUntil   int64 `json:"blocked_end"`
	Moral        int   `json:"moral"`
}

// BlockStatusResp BlockStatusResp
type BlockStatusResp struct {
	Code int       `json:"code"`
	Data BlockInfo `json:"data"`
	Msg  string    `json:"msg"`
}

// BlockInfo BlockInfo
func (dao *BlockStatusDao) BlockInfo(c context.Context, mid int64) (*BlockInfo, error) {
	var res BlockStatusResp
	params := url.Values{}
	params.Set("mid", fmt.Sprintf("%d", mid))
	err := dao.httpClient.Get(c, dao.requestURL, "", params, &res)
	// WARNING: must put urlStr after httpClient.Get()
	// since params will be modified
	urlStr := dao.requestURL + "?" + params.Encode()
	if err != nil {
		log.Error("call 账号系统小黑屋(%s) error(%v)", urlStr, err)
		return nil, err
	}
	if res.Code != 0 {
		err = fmt.Errorf("call 账号系统小黑屋(%s) error, return code(%v)", urlStr, res.Code)
		log.Error("%v", err)
		return nil, err
	}
	log.Info("call 账号系统小黑屋(%s) successful. resp(%+v)", urlStr, res.Data)
	resData := res.Data
	if resData.PassTest != 0 && resData.PassTest != 1 {
		err = fmt.Errorf("call 账号系统小黑屋(%s) successful. but got error resp(blocked_status=%d), want 0 or 1", urlStr, resData.PassTest)
		log.Error("%v", err)
		return nil, err
	}
	return &res.Data, nil
}
