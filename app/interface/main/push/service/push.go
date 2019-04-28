package service

import (
	"context"
	"net/url"
	"strconv"

	pushmdl "go-common/app/service/main/push/model"
	"go-common/library/log"
)

const (
	_testTokenURL = "http://api.bilibili.co/x/internal/push-service/test/token"
)

// TestToken for test via push token.
func (s *Service) TestToken(ctx context.Context, info *pushmdl.PushInfo, token string) (err error) {
	params := url.Values{}
	params.Add("app_id", strconv.FormatInt(info.APPID, 10))
	params.Add("alert_title", info.Title)
	params.Add("alert_body", info.Summary)
	params.Add("token", token)
	params.Add("link_type", strconv.FormatInt(int64(info.LinkType), 10))
	params.Add("link_value", info.LinkValue)
	params.Add("sound", strconv.Itoa(info.Sound))
	params.Add("vibration", strconv.Itoa(info.Vibration))
	params.Add("expire_time", strconv.FormatInt(int64(info.ExpireTime), 10))
	params.Add("image_url", info.ImageURL)
	if err = s.httpClient.Post(ctx, _testTokenURL, "", params, nil); err != nil {
		log.Error("s.TestToken(%+v) error(%v)", info, err)
	}
	return
}
