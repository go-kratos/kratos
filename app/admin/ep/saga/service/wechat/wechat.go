package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/dao"
	"go-common/app/admin/ep/saga/model"
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

// NewTxtNotify create wechat format text notification 从配置初始化企业微信TxtNotification
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
	if token, err = w.dao.AccessTokenRedis(c, key); err != nil {
		log.Warn("AccessToken: failed to get access_token from cache, appId (%s), error (%s)", app.AppID, err.Error())
		//企业微信api获取公司token
		if token, expire, err = w.dao.WechatAccessToken(c, app.AppSecret); err != nil {
			err = errors.Wrapf(err, "AccessToken: both mc and api can't provide access_token, appId(%s)", app.AppID)
			return
		}

		// 通过API获取到了，缓存一波
		err = w.dao.SetAccessTokenRedis(c, key, token, expire)
		return
	}

	if token == "" {
		if token, expire, err = w.dao.WechatAccessToken(c, app.AppSecret); err != nil {
			return
		}

		// 通过API获取到了，缓存一波
		err = w.dao.SetAccessTokenRedis(c, key, token, expire)
	}

	return
}

// PushMsg push text message via wechat notification api with access_token.推送企业微信
func (w *Wechat) PushMsg(c context.Context, userNames []string, content string) (err error) {
	var (
		token         string
		userIds       string
		invalidUser   string
		userNamesByte []byte
		txtMsg        = w.NewTxtNotify(content)
		contentDB     = content
	)
	//获取企业token
	if token, err = w.AccessToken(c, w.saga); err != nil {
		return
	}

	if token == "" {
		err = errors.Errorf("PushMsg: get access token failed, it's empty. appid (%s), secret (%s)", w.saga.AppID, w.saga.AppSecret)
		return
	}
	//员工编号以竖线分隔
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
	if userNamesByte, err = json.Marshal(userNames); err != nil {
		return
	}
	if len(contentDB) > model.MaxWechatLen {
		contentDB = contentDB[:model.MaxWechatLen]
	}
	messageLog := &model.WechatMessageLog{
		Touser:  string(userNamesByte),
		Content: contentDB,
		Status:  1,
	}
	return w.dao.CreateMessageLog(messageLog)
}

// UserIds query user ids for user name list 查询员工编号
func (w *Wechat) UserIds(userNames []string) (ids string, err error) {
	if ids, err = w.dao.UserIds(userNames); err != nil {
		return
	}
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
		//查看是否缓存，缓存则继续
		if alreadyIn, err = w.alreadyInCache(c, userID); err != nil || alreadyIn {
			continue
		}
		//未缓存从数据库查询
		if contactInfo, err = w.dao.QueryUserByID(userID); err != nil {
			log.Error("no such userID (%s) in db, error(%s)", userID, err.Error())
			return
		}
		//数据库查询结果缓存
		if err = w.dao.SetRequireVisibleUsersRedis(c, contactInfo); err != nil {
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
	//查询所有的值
	if err = w.dao.RequireVisibleUsersRedis(c, &userMap); err != nil {
		log.Error("get userID (%s) from cache error(%s)", userID, err.Error())
		return
	}
	//匹配需要查询的用户id
	for k, v := range userMap {
		if userID == k {
			log.Info("(%s) is already exist in cache, value(%v)", k, v)
			alreadyIn = true
			return
		}
	}
	return
}
