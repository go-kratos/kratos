package dao

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"
)

const (
	_gitUsersAPI = "http://git.bilibili.co/api/v4/users"
)

// GitLabFace  return face of gitlab.
func (d *Dao) GitLabFace(c context.Context, username string) (avatarURL string, err error) {
	params := url.Values{}
	params.Set("username", username)
	params.Set("private_token", conf.Conf.Gitlab.Token)
	var req *http.Request
	if req, err = http.NewRequest(http.MethodGet, _gitUsersAPI, strings.NewReader(params.Encode())); err != nil {
		log.Error("http.NewRequest(%s) error(%v)", username, err)
		return
	}
	res := make([]*ut.Image, 0)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("d.client.Do(%s) error(%v)", username, err)
		return
	}
	for _, v := range res {
		avatarURL = v.AvatarURL
	}
	return
}
