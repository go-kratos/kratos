package dao

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"go-common/app/service/main/account-recovery/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// GetMidInfo get mid info by more condition
func (d *Dao) GetMidInfo(c context.Context, qType string, qKey string) (v *model.MIDInfo, err error) {
	params := url.Values{}
	params.Set("q_type", qType)
	params.Set("q_key", qKey)
	res := new(struct {
		Code int           `json:"code"`
		Data model.MIDInfo `json:"data"`
	})
	if err = d.httpClient.Get(c, d.c.AccRecover.MidInfoURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("GetMidInfo HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("GetMidInfo server err_code %d,params: qType=%s,qKey=%s", res.Code, qType, qKey)
		err = ecode.ServerErr
		return
	}
	log.Info("GetMidInfo url=%v, params: qType=%s,qKey=%s, res: %+v", d.c.AccRecover.MidInfoURL, qType, qKey, res)
	return &res.Data, nil
}

// GetUserInfo get user info by mid
func (d *Dao) GetUserInfo(c context.Context, mid int64) (v *model.UserInfo, err error) {
	params := url.Values{}
	params.Add("mid", strconv.Itoa(int(mid)))
	res := new(struct {
		Code int            `json:"code"`
		Data model.UserInfo `json:"data"`
	})
	if err = d.httpClient.Get(c, d.c.AccRecover.GetUserInfoURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("GetUserInfo HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("GetUserInfo server err_code %d,params: mid=%d", res.Code, mid)
		err = ecode.ServerErr
		return
	}
	log.Info("GetUserInfo url=%v, params: mid=%d, res: %+v", d.c.AccRecover.GetUserInfoURL, mid, res)
	return &res.Data, nil
}

// UpdatePwd update password
func (d *Dao) UpdatePwd(c context.Context, mid int64, operator string) (user *model.User, err error) {
	params := url.Values{}
	params.Set("mid", strconv.Itoa(int(mid)))
	params.Set("operator", operator)
	res := new(struct {
		Code int        `json:"code"`
		Data model.User `json:"data"`
	})
	if err = d.httpClient.Post(c, d.c.AccRecover.UpPwdURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("UpdatePwd HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("UpdatePwd server err_code %d,params: mid=%d", res.Code, mid)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("UpdatePwd url=%v, params: mid=%d,operator=%s, res: %+v", d.c.AccRecover.UpPwdURL, mid, operator, res)
	return &res.Data, nil
}

// CheckSafe safe info
func (d *Dao) CheckSafe(c context.Context, mid int64, question int8, answer string) (check *model.Check, err error) {
	params := url.Values{}
	params.Add("mid", strconv.Itoa(int(mid)))
	params.Add("question", strconv.Itoa(int(question)))
	params.Add("answer", answer)
	res := new(struct {
		Code int         `json:"code"`
		Data model.Check `json:"data"`
	})
	if err = d.httpClient.Post(c, d.c.AccRecover.CheckSafeURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("CheckSafe HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("CheckSafe server err_code %d,params: mid=%d,question=%d,answer=%s", res.Code, mid, question, answer)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("CheckSafe url=%v, params: mid=%d,question=%d,answer=%s, res: %+v", d.c.AccRecover.CheckSafeURL, mid, question, answer, res)
	return &res.Data, nil
}

// GetUserType get user_type
func (d *Dao) GetUserType(c context.Context, mid int64) (gams []*model.Game, err error) {
	params := url.Values{}
	params.Add("mid", strconv.Itoa(int(mid)))
	res := new(struct {
		Code int           `json:"code"`
		Data []*model.Game `json:"items"`
	})
	if err = d.httpClient.Get(c, d.c.AccRecover.GameURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("GetUserType HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("GetUserType server err_code %d,params: mid=%d", res.Code, mid)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("GetUserType url=%v, params: mid=%d, res: %+v", d.c.AccRecover.GameURL, mid, res)
	return res.Data, nil
}

// CheckReg check reg info
func (d *Dao) CheckReg(c context.Context, mid int64, regTime int64, regType int8, regAddr string) (v *model.Check, err error) {
	params := url.Values{}
	params.Add("mid", strconv.Itoa(int(mid)))
	params.Add("reg_time", strconv.FormatInt(regTime, 10))
	params.Add("reg_type", strconv.Itoa(int(regType)))
	params.Add("reg_addr", regAddr)
	res := new(struct {
		Code int         `json:"code"`
		Data model.Check `json:"data"`
	})
	if err = d.httpClient.Post(c, d.c.AccRecover.CheckRegURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("CheckReg HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("CheckReg server err_code %d,params: mid=%d,regTime=%d,regType=%d,regAddr=%s", res.Code, mid, regTime, regType, regAddr)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("CheckReg url=%v, params: mid=%d,regTime=%d,regType=%d,regAddr=%s, res: %+v", d.c.AccRecover.CheckRegURL, mid, regTime, regType, regAddr, res)
	return &res.Data, nil
}

// UpdateBatchPwd batch update password
func (d *Dao) UpdateBatchPwd(c context.Context, mids string, operator string) (userMap map[string]*model.User, err error) {
	params := url.Values{}
	params.Set("mids", mids)
	params.Set("operator", operator)
	res := new(struct {
		Code int                    `json:"code"`
		Data map[string]*model.User `json:"data"`
	})
	if err = d.httpClient.Post(c, d.c.AccRecover.UpBatchPwdURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("UpdateBatchPwd HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("UpdateBatchPwd server err_code %d,params: mids=%s", res.Code, mids)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("UpdateBatchPwd url=%v, params: mid=%s,operator=%s, res: %+v", d.c.AccRecover.UpBatchPwdURL, mids, operator, res)
	return res.Data, nil
}

// CheckCard check card
func (d *Dao) CheckCard(c context.Context, mid int64, cardType int8, cardCode string) (ok bool, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("card_type", strconv.Itoa(int(cardType)))
	params.Set("card_code", cardCode)
	res := new(struct {
		Code int  `json:"code"`
		Data bool `json:"data"`
	})
	if err = d.httpClient.Get(c, d.c.AccRecover.CheckCardURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("CheckCard HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("CheckCard server err_code %d,params: mid=%d,cardType=%d,cardCode=%s", res.Code, mid, cardType, cardCode)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("CheckCard url=%v, params: mid=%d,cardType=%d,cardCode=%s, res: %+v", d.c.AccRecover.CheckCardURL, mid, cardType, cardCode, res)
	return res.Data, nil
}

// CheckPwds check pwd
func (d *Dao) CheckPwds(c context.Context, mid int64, pwds string) (v string, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("pwd", pwds)
	res := new(struct {
		Code int    `json:"code"`
		Data string `json:"data"`
	})
	if err = d.httpClient.Post(c, d.c.AccRecover.CheckPwdURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("CheckPwds HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("CheckPwds server err_code %d,params: mid=%d,pwds=%s", res.Code, mid, pwds)
		err = ecode.Int(res.Code)
		return
	}
	log.Info("CheckPwds url=%v, params: mid=%d,pwds=%s, res: %+v", d.c.AccRecover.CheckPwdURL, mid, pwds, res)
	return res.Data, nil
}

// GetLoginIPs get login ip
func (d *Dao) GetLoginIPs(c context.Context, mid int64, limit int64) (ipInfo []*model.LoginIPInfo, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("limit", strconv.FormatInt(limit, 10))
	res := new(struct {
		Code int                  `json:"code"`
		Data []*model.LoginIPInfo `json:"data"`
	})
	if err = d.httpClient.Get(c, d.c.AccRecover.GetLoginIPURL, metadata.String(c, metadata.RemoteIP), params, &res); err != nil {
		log.Error("GetLoginIPs HTTP request err %+v", err)
		return
	}
	if res.Code != 0 {
		log.Error("GetLoginIPs server err_code %d,params: mid=%d,limit=%d", res.Code, mid, limit)
		err = ecode.ServerErr
		return
	}
	log.Info("GetLoginIPs url=%v, params: mid=%d,limit=%d, res: %+v", d.c.AccRecover.GetLoginIPURL, mid, limit, res)
	return res.Data, nil
}

// GetAddrByIP get addr by ip
func (d *Dao) GetAddrByIP(c context.Context, mid int64, limit int64) (addrs string, err error) {
	ipInfo, err := d.GetLoginIPs(c, mid, limit)
	if err != nil || len(ipInfo) == 0 {
		return
	}
	var ipLen = len(ipInfo)
	ips := make([]string, 0, ipLen)
	//ip去重复和空串
	for i := 0; i < ipLen; i++ {
		if (i > 0 && ipInfo[i-1].LoginIP == ipInfo[i].LoginIP) || len(ipInfo[i].LoginIP) == 0 {
			continue
		}
		ips = append(ips, ipInfo[i].LoginIP)
	}
	ipMap, err := d.Infos(c, ips)
	i := 0
	for _, loc := range ipMap {
		if loc.Country != "" {
			addrs += loc.Country + "-"
		}
		if loc.Province != "" {
			addrs += loc.Province + "-"
		}
		if loc.City != "" {
			addrs += loc.City + "-"
		}
		addrs = strings.TrimRight(addrs, "-") + ","
		i++
		if i >= 3 {
			break
		}
	}
	addrs = strings.TrimRight(addrs, ",")
	return
}
