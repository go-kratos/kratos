package dao

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
)

// NickNameLog NickNameLog
func (d *Dao) NickNameLog(c context.Context, nickNameReq *model.NickNameReq) (res *model.NickNameLogRes, err error) {
	nowYear := time.Now().Year()
	index1 := "log_user_action_14_" + strconv.Itoa(nowYear)
	index2 := "log_user_action_14_" + strconv.Itoa(nowYear-1)

	r := d.es.NewRequest("log_user_action").Fields("str_0", "str_1").Index(index1, index2)
	r.Order("ctime", elastic.OrderDesc).Order("mid", elastic.OrderDesc).Pn(nickNameReq.Page).Ps(nickNameReq.Size)

	if nickNameReq.Mid != 0 {
		r.WhereEq("mid", nickNameReq.Mid)
	}

	if nickNameReq.From != 0 && nickNameReq.To != 0 {
		ftm := time.Unix(nickNameReq.From, 0)
		sf := ftm.Format("2006-01-02 15:04:05")

		ttm := time.Unix(nickNameReq.To, 0)
		tf := ttm.Format("2006-01-02 15:04:05")

		r.WhereRange("ctime", sf, tf, elastic.RangeScopeLoRo)
	}

	esres := new(model.NickESRes)
	if err = r.Scan(context.TODO(), &esres); err != nil {
		log.Error("nickNameLog search error(%v)", err)
		return
	}

	var nickNames = make([]*model.NickNameInfo, 0)
	for _, value := range esres.Result {
		ulr := model.NickNameInfo{OldName: value.OldName, NewName: value.NewName}
		nickNames = append(nickNames, &ulr)
	}
	res = &model.NickNameLogRes{Page: esres.Page, Result: nickNames}
	return
}

type userLogsExtra struct {
	EncryptTel   string `json:"tel"`
	EncryptEmail string `json:"email"`
}

// UserBindLog User bind log
func (d *Dao) UserBindLog(c context.Context, userActLogReq *model.UserBindLogReq) (res *model.UserBindLogRes, err error) {
	e := d.es
	nowYear := time.Now().Year()

	var count = 2 //默认查询两年
	//2016年就有了手机历史记录，此处需要循环建立索引 , 2018年才有邮箱这个功能
	if userActLogReq.Action == "telBindLog" {
		count = nowYear - 2015
	}
	if userActLogReq.Action == "emailBindLog" {
		count = nowYear - 2017
	}
	indexs := make([]string, count)
	for i := 0; i < count; i++ {
		indexs[i] = "log_user_action_54_" + strconv.Itoa(nowYear-i)
	}

	r := e.NewRequest("log_user_action").Fields("mid", "str_0", "extra_data", "ctime").Index(indexs...)
	r.Order("ctime", elastic.OrderDesc).Order("mid", elastic.OrderDesc).Pn(userActLogReq.Page).Ps(userActLogReq.Size)

	if userActLogReq.Mid != 0 {
		r.WhereEq("mid", userActLogReq.Mid)
	}

	if userActLogReq.Query != "" {
		hash := sha1.New()
		hash.Write([]byte(userActLogReq.Query))
		telHash := base64.StdEncoding.EncodeToString(hash.Sum(d.hashSalt))
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
		return
	}

	var userBindLogs = make([]*model.UserBindLog, 0)
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
			email, err = d.decrypt(userLogExtra.EncryptEmail)
			if err != nil {
				log.Error("EncryptEmail decode err(%v)", err)
				continue
			}
		}
		if userLogExtra.EncryptTel != "" {
			tel, err = d.decrypt(userLogExtra.EncryptTel)
			if err != nil {
				log.Error("EncryptTel decode err(%v)", err)
				continue
			}
		}
		ulr := model.UserBindLog{Mid: value.Mid, Email: email, Phone: tel, Time: value.CTime}
		userBindLogs = append(userBindLogs, &ulr)
	}
	res = &model.UserBindLogRes{Page: esres.Page, Result: userBindLogs}
	return
}
