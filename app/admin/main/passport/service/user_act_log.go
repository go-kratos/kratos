package service

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/admin/main/passport/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

type userLogsExtra struct {
	EncryptTel   string `json:"tel"`
	EncryptEmail string `json:"email"`
}

// UserBindLog User bind log
func (s *Service) UserBindLog(c context.Context, userActLogReq *model.UserBindLogReq) (res *model.UserBindLogRes, err error) {
	e := s.dao.EsCli
	nowYear := time.Now().Year()
	index1 := "log_user_action_54_" + strconv.Itoa(nowYear)
	index2 := "log_user_action_54_" + strconv.Itoa(nowYear-1)

	r := e.NewRequest("log_user_action").Fields("mid", "str_0", "extra_data", "ctime").Index(index1, index2)
	r.Order("ctime", elastic.OrderDesc).Order("mid", elastic.OrderDesc).Pn(userActLogReq.Page).Ps(userActLogReq.Size)

	if userActLogReq.Mid != 0 {
		r.WhereEq("mid", userActLogReq.Mid)
	}

	if userActLogReq.Query != "" {
		hash := sha1.New()
		hash.Write([]byte(userActLogReq.Query))
		telHash := base64.StdEncoding.EncodeToString(hash.Sum(s.hashSalt))
		r.WhereEq("str_0", telHash)
	}

	if userActLogReq.Action != "" {
		r.WhereEq("action", userActLogReq.Action)
	}

	if userActLogReq.From != 0 && userActLogReq.To != 0 {
		ftm := time.Unix(userActLogReq.From, 0)
		sf := ftm.Format("2006-01-02 15:04:05")

		ttm := time.Unix(userActLogReq.To, 0)
		tf := ttm.Format("2006-01-02 15:04:05")

		r.WhereRange("ctime", sf, tf, elastic.RangeScopeLoRo)
	}

	esres := new(model.EsRes)
	if err = r.Scan(context.Background(), &esres); err != nil {
		log.Error("userActLogs search error(%v)", err)
	}

	var users = make([]*model.UserBindLog, 0)
	for _, value := range esres.Result {

		var email, tel string
		//var model.UserBindLog
		userLogExtra := userLogsExtra{}
		err = json.Unmarshal([]byte(value.ExtraData), &userLogExtra)
		if err != nil {
			log.Error("cannot convert json(%s) to struct,err(%+v) ", value.ExtraData, err)
			continue
		}
		if userLogExtra.EncryptEmail != "" {
			email, err = s.decrypt(userLogExtra.EncryptEmail)
			if err != nil {
				log.Error("EncryptEmail decode err(%v)", err)
				continue
			}
		}
		if userLogExtra.EncryptTel != "" {
			tel, err = s.decrypt(userLogExtra.EncryptTel)
			if err != nil {
				log.Error("EncryptTel decode err(%v)", err)
				continue
			}
		}
		ulr := model.UserBindLog{Mid: value.Mid, EMail: email, Phone: tel, Time: value.CTime}
		users = append(users, &ulr)
	}
	res = &model.UserBindLogRes{Page: esres.Page, Result: users}
	return
}

// DecryptBindLog decrypt bind log
func (s *Service) DecryptBindLog(c context.Context, reqParams *model.DecryptBindLogParam) (res map[string]string, err error) {
	if len(reqParams.EncryptText) == 0 {
		return make(map[string]string), nil
	}
	res = make(map[string]string, len(reqParams.EncryptText))
	for _, v := range reqParams.EncryptText {
		var tel string
		if tel, err = s.decrypt(v); err != nil {
			return
		}
		res[v] = tel
	}
	return
}
