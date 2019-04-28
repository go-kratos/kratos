package wechat

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/tool/saga/conf"
	"go-common/app/tool/saga/dao"
	"go-common/app/tool/saga/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// Wechat 企业微信应用
type Wechat struct {
	dao     *dao.Dao
	saga    *model.AppConfig
	contact *model.AppConfig
}

// New create an new wechat work
func New(d *dao.Dao) (w *Wechat) {
	w = &Wechat{
		dao:     d,
		saga:    conf.Conf.Property.Wechat,
		contact: conf.Conf.Property.Contact,
	}

	return w
}

// NewTxtNotify create wechat format text notification
func (w *Wechat) NewTxtNotify(content string) (txtMsg *model.TxtNotification) {
	return &model.TxtNotification{
		Notification: model.Notification{
			MsgType: "text",
			AgentID: w.saga.AppID,
		},
		Body: model.Text{
			Content: content,
		},
		Safe: 0,
	}
}

// AccessToken get access_token from cache first, if not found, get it via wechat api.
func (w *Wechat) AccessToken(c context.Context, app *model.AppConfig) (token string, err error) {
	var (
		key    string
		expire int32
	)

	key = fmt.Sprintf("appid_%d", app.AppID)
	if token, err = w.dao.AccessToken(c, key); err != nil {
		log.Warn("AccessToken: failed to get access_token from cache, appId (%d), error (%s)", app.AppID, err.Error())
		if token, expire, err = w.dao.WechatAccessToken(c, app.AppSecret); err != nil {
			err = errors.Wrapf(err, "AccessToken: both mc and api can't provide access_token, appId(%d)", app.AppID)
			return
		}

		// 通过API获取到了，缓存一波
		err = w.dao.SetAccessToken(c, key, token, expire)
		return
	}

	if token == "" {
		if token, expire, err = w.dao.WechatAccessToken(c, app.AppSecret); err != nil {
			return
		}

		// 通过API获取到了，缓存一波
		err = w.dao.SetAccessToken(c, key, token, expire)
	}

	return
}

// PushMsg push text message via wechat notification api with access_token.
func (w *Wechat) PushMsg(c context.Context, userNames []string, content string) (err error) {
	var (
		token       string
		userIds     string
		invalidUser string
		txtMsg      = w.NewTxtNotify(content)
	)

	if token, err = w.AccessToken(c, w.saga); err != nil {
		return
	}

	if token == "" {
		err = errors.Errorf("PushMsg: get access token failed, it's empty. appid (%d), secret (%s)", w.saga.AppID, w.saga.AppSecret)
		return
	}

	if userIds, err = w.UserIds(userNames); err != nil {
		return
	}
	txtMsg.ToUser = userIds

	if invalidUser, err = w.dao.WechatPushMsg(c, token, txtMsg); err != nil {
		if err = w.addRequireVisible(c, invalidUser); err != nil {
			log.Error("PushMsg add userID (%s) in cache, error(%s)", invalidUser, err.Error())
		}
		return
	}
	return
}

// UserIds query user ids for user name list
func (w *Wechat) UserIds(userNames []string) (ids string, err error) {
	ids, err = w.dao.UserIds(userNames)
	return
}

// addRequireVisible update wechat require visible users in memcache
func (w *Wechat) addRequireVisible(c context.Context, userIDs string) (err error) {
	var (
		contactInfo *model.ContactInfo
		userID      string
		alreadyIn   bool
	)

	users := strings.Split(userIDs, "|")
	for _, userID = range users {

		if alreadyIn, err = w.alreadyInCache(c, userID); err != nil || alreadyIn {
			continue
		}

		if contactInfo, err = w.dao.QueryUserByID(userID); err != nil {
			log.Error("no such userID (%s) in db, error(%s)", userID, err.Error())
			return
		}

		if err = w.dao.SetRequireVisibleUsers(c, contactInfo); err != nil {
			log.Error("failed set to cache userID (%s) username (%s), err (%s)", userID, contactInfo.UserName, err.Error())
			return
		}
	}
	return
}

// alreadyInCache check user is or not in the memcache
func (w *Wechat) alreadyInCache(c context.Context, userID string) (alreadyIn bool, err error) {
	var (
		userMap = make(map[string]model.RequireVisibleUser)
	)

	if err = w.dao.RequireVisibleUsers(c, &userMap); err != nil {
		log.Error("get userID (%s) from cache error(%s)", userID, err.Error())
		return
	}

	for k, v := range userMap {
		if userID == k {
			log.Info("(%s) is already exist in cache, value(%v)", k, v)
			alreadyIn = true
			return
		}
	}
	return
}
