package service

import (
	"context"
	"fmt"
	"net/url"

	"go-common/app/admin/ep/saga/conf"
	"go-common/app/admin/ep/saga/model"
	"go-common/app/admin/ep/saga/service/wechat"
	"go-common/library/log"

	"github.com/pkg/errors"
)

const qyWechatURL = "https://qyapi.weixin.qq.com"

// CollectWachatUsers send required wechat visible users stored in memcache by email
func (s *Service) CollectWachatUsers(c context.Context) (err error) {
	var (
		contactInfo *model.ContactInfo
		userMap     = make(map[string]model.RequireVisibleUser)
		user        string
	)

	if err = s.dao.RequireVisibleUsersRedis(c, &userMap); err != nil {
		log.Error("get require visible user error(%v)", err)
		return
	}

	for k, v := range userMap {
		if contactInfo, err = s.dao.QueryUserByID(k); err == nil {
			if contactInfo.VisibleSaga {
				continue
			}
		}
		user += v.UserName + " ,  " + v.NickName + "\n"
	}

	/*content := fmt.Sprintf("\n\n邮箱前缀     昵称\n\n%s", user)
	for _, addr := range conf.Conf.Property.ReportRequiredVisible.AlertAddrs {
		if err = mail.SendMail2(addr, "需添加的企业微信名单", content); err != nil {
			return
		}
	}*/
	if err = s.dao.DeleteRequireVisibleUsersRedis(c); err != nil {
		log.Error("Delete require visible user error(%v)", err)
		return
	}

	return
}

// SyncContacts sync the wechat contacts 更新企业微信列表（用户信息和saga信息）
func (s *Service) SyncContacts(c context.Context) (err error) {
	var (
		w = wechat.New(s.dao)
	)

	if err = w.SyncContacts(c); err != nil {
		return
	}
	return
}

// QueryContacts query machine logs.
func (s *Service) QueryContacts(c context.Context, queryRequest *model.Pagination) (p *model.PaginateContact, err error) {
	var (
		total    int64
		contacts []*model.ContactInfo
	)
	fmt.Print(queryRequest.PageNum)
	if total, contacts, err = s.dao.FindContacts(queryRequest.PageNum, queryRequest.PageSize); err != nil {
		return
	}
	fmt.Print(queryRequest.PageNum)
	p = &model.PaginateContact{
		PageNum:  queryRequest.PageNum,
		PageSize: queryRequest.PageSize,
		Total:    total,
		Contacts: contacts,
	}
	return
}

// QueryContactLogs query contact logs.
func (s *Service) QueryContactLogs(c context.Context, queryRequest *model.QueryContactLogRequest) (p *model.PaginateContactLog, err error) {
	var (
		total       int64
		machineLogs []*model.AboundContactLog
	)
	if total, machineLogs, err = s.dao.FindMachineLogs(queryRequest); err != nil {
		return
	}
	p = &model.PaginateContactLog{
		PageNum:     queryRequest.PageNum,
		PageSize:    queryRequest.PageSize,
		Total:       total,
		MachineLogs: machineLogs,
	}
	return
}

// Wechat ...
func (s *Service) Wechat() *wechat.Wechat {
	return wechat.New(s.dao)
}

// CreateWechat ...
func (s *Service) CreateWechat(c context.Context, req *model.CreateChatReq, username string) (resp *model.CreateChatResp, err error) {
	var (
		token     string
		userIDs   []string
		ownerInfo *model.ContactInfo
		w         = wechat.New(s.dao)
	)
	u := qyWechatURL + "/cgi-bin/appchat/create"
	params := url.Values{}
	wechatInfo := &model.WechatCreateLog{
		Name:   req.Name,
		Owner:  req.Owner,
		ChatID: req.ChatID,
		Cuser:  username,
		Status: 1,
	}

	//获取企业token
	if token, err = w.AccessToken(c, conf.Conf.Property.Wechat); err != nil {
		return
	}
	params.Set("access_token", token)

	//get owner and users id
	if ownerInfo, err = s.dao.QueryUserByUserName(req.Owner); err != nil {
		return
	}
	if userIDs, err = s.QueryUserIds(req.UserList); err != nil {
		return
	}

	req.Owner = ownerInfo.UserID
	req.UserList = userIDs
	if err = s.dao.PostJSON(c, u, "", params, &resp, req); err != nil {
		return
	}

	//add create wechat info to database
	if err = s.dao.AddWechatCreateLog(wechatInfo); err != nil {
		return
	}

	resp = &model.CreateChatResp{
		ChatID: wechatInfo.ChatID,
	}
	return
}

// QueryUserIds ...
func (s *Service) QueryUserIds(userNames []string) (userIds []string, err error) {
	var (
		userName    string
		contactInfo *model.ContactInfo
	)
	if len(userNames) == 0 {
		err = errors.Errorf("UserIds: userNames is empty!")
		return
	}

	for _, userName = range userNames {
		if contactInfo, err = s.dao.QueryUserByUserName(userName); err != nil {
			err = errors.Wrapf(err, "UserIds: no such user (%s) in db, err (%s)", userName, err.Error())
			return
		}

		log.Info("UserIds: username (%s), userid (%s)", userName, contactInfo.UserID)
		if contactInfo.UserID != "" {
			userIds = append(userIds, contactInfo.UserID)
		}
	}
	return
}

// QueryWechatCreateLog ...
func (s *Service) QueryWechatCreateLog(c context.Context, req *model.Pagination, username string) (resp *model.CreateChatLogResp, err error) {
	var (
		logs             []*model.WechatCreateLog
		logsResp         []*model.CreateChatLog
		total            int
		wechatCreateInfo *model.WechatCreateLog
	)

	if logs, total, err = s.dao.QueryWechatCreateLog(true, req, wechatCreateInfo); err != nil {
		return
	}
	for _, log := range logs {
		createChatlog := &model.CreateChatLog{}
		if log.Cuser == username {
			createChatlog.Buttons = append(createChatlog.Buttons, "WECHAT_TEST")
		}
		createChatlog.WechatCreateLog = log
		logsResp = append(logsResp, createChatlog)
	}
	resp = &model.CreateChatLogResp{
		Total:      total,
		Pagination: req,
		Logs:       logsResp,
	}
	return
}

// WechatParams ...
func (s *Service) WechatParams(c context.Context, chatid string) (resp *model.GetChatResp, err error) {
	var (
		w     = wechat.New(s.dao)
		token string
	)

	if token, err = w.AccessToken(c, conf.Conf.Property.Wechat); err != nil {
		return
	}

	u := qyWechatURL + "/cgi-bin/appchat/get"
	params := url.Values{}
	params.Set("access_token", token)
	params.Set("chatid", chatid)
	err = s.dao.WechatParams(c, u, params, &resp)
	return
}

// SendGroupWechat ...
func (s *Service) SendGroupWechat(c context.Context, req *model.SendChatReq) (resp *model.ChatResp, err error) {
	var (
		token       string
		w           = wechat.New(s.dao)
		total       int
		getChatResp *model.GetChatResp
		owner       string
		contentDB   = req.Text.Content
	)
	u := qyWechatURL + "/cgi-bin/appchat/send"
	params := url.Values{}

	if token, err = w.AccessToken(c, conf.Conf.Property.Wechat); err != nil {
		return
	}
	params.Set("access_token", token)

	if err = s.dao.PostJSON(c, u, "", params, &resp, req); err != nil {
		return
	}
	if len(contentDB) > model.MaxWechatLen {
		contentDB = contentDB[:model.MaxWechatLen]
	}
	chatLog := &model.WechatChatLog{
		ChatID:  req.ChatID,
		MsgType: req.MsgType,
		Content: contentDB,
		Safe:    req.Safe,
		Status:  1,
	}
	if err = s.dao.CreateChatLog(chatLog); err != nil {
		return
	}

	info := &model.WechatCreateLog{
		ChatID: req.ChatID,
	}
	if _, total, err = s.dao.QueryWechatCreateLog(false, nil, info); err != nil {
		return
	}

	if total == 0 {
		getChatResp, _ = s.WechatParams(c, req.ChatID)
		owner = getChatResp.ChatInfo.Owner
		contactInfo, _ := s.dao.QueryUserByID(owner)
		wechatInfo := &model.WechatCreateLog{
			Name:   getChatResp.ChatInfo.Name,
			Owner:  contactInfo.UserName,
			ChatID: req.ChatID,
			Status: 2,
		}
		if err = s.dao.AddWechatCreateLog(wechatInfo); err != nil {
			return
		}
	}
	return
}

// SendWechat ...
func (s *Service) SendWechat(c context.Context, req *model.SendMessageReq) (resp *model.ChatResp, err error) {
	var (
		w = wechat.New(s.dao)
	)
	err = w.PushMsg(c, req.Touser, req.Content)
	return
}

// UpdateWechat ...
func (s *Service) UpdateWechat(c context.Context, req *model.UpdateChatReq) (resp *model.ChatResp, err error) {
	var (
		token string
		w     = wechat.New(s.dao)
	)
	u := qyWechatURL + "/cgi-bin/appchat/update"
	params := url.Values{}

	if token, err = w.AccessToken(c, conf.Conf.Property.Wechat); err != nil {
		return
	}
	params.Set("access_token", token)
	if err = s.dao.PostJSON(c, u, "", params, &resp, req); err != nil {
		return
	}
	return
}

// SyncWechatContacts ...
func (s *Service) SyncWechatContacts(c context.Context) (message string, err error) {
	var (
		w = wechat.New(s.dao)
	)
	if err = w.AnalysisContacts(c); err != nil {
		return
	}
	message = "同步完成"
	return
}
