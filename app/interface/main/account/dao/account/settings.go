package account

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	xhttp "net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/interface/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	updateFaceURI = "/api/user/updateFace"
	identifyURI   = "/api/internal/identify/info"

	updateBirthURI = "/api/internal/member/updateBirthday"
	updateSexURI   = "/api/internal/member/updateSex"
	// updatePerson   = "/api/internal/member/updatePerson"
	// updateUnameURI = "/api/internal/member/updateUname"
)

// UpdateBirthday is
func (d *Dao) UpdateBirthday(c context.Context, mid int64, ip, birthday string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("birthday", birthday)
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.accCo+updateBirthURI, ip, params, &res); err != nil {
		log.Error("UpdateBirthday url(%s) error(%v)", updateBirthURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("UpdateBirthday url(%s) code(%v)", updateBirthURI+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// UpdateSign update sign.
// func (d *Dao) UpdateSign(c context.Context, mid int64, ip, sign string) (err error) {
// 	params := url.Values{}
// 	params.Set("mid", strconv.FormatInt(mid, 10))
// 	params.Set("user_sign", sign)
// 	var res struct {
// 		Code int `json:"code"`
// 	}
// 	if err = d.client.Post(c, d.accCo+updateSignURI, ip, params, &res); err != nil {
// 		log.Error("UpdateSign url(%s) error(%v)", updateBirthURI+"?"+params.Encode(), err)
// 		return
// 	}
// 	if res.Code != 0 {
// 		log.Error("UpdateSign url(%s) code(%v)", updateBirthURI+"?"+params.Encode(), res.Code)
// 		err = ParseJavaCode(res.Code)
// 	}
// 	return
// }

// UpdateSex update sex
func (d *Dao) UpdateSex(c context.Context, mid, sex int64, ip string) (err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("sex", strconv.FormatInt(sex, 10))
	var res struct {
		Code int `json:"code"`
	}
	if err = d.client.Post(c, d.accCo+updateSexURI, ip, params, &res); err != nil {
		log.Error("UpdateSex url(%s) error(%v)", updateSexURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("UpdateSex url(%s) error(%v)", updateSexURI+"?"+params.Encode(), res.Code)
		err = ecode.Int(res.Code)
	}
	return
}

// UpdatePerson update sex
// func (d *Dao) UpdatePerson(c context.Context, mid int64, birthday, ip string, datingType, marital, place int64) (err error) {
// 	params := url.Values{}
// 	params.Set("mid", strconv.FormatInt(mid, 10))
// 	params.Set("birthday", birthday)
// 	params.Set("datingtype", strconv.FormatInt(datingType, 10))
// 	params.Set("marital", strconv.FormatInt(marital, 10))
// 	params.Set("place", strconv.FormatInt(place, 10))
// 	var res struct {
// 		Code int `json:"code"`
// 	}
// 	if err = d.client.Post(c, d.accCo+updatePerson, ip, params, &res); err != nil {
// 		log.Error("UpdatePerson url(%s) error(%v)", updatePerson+"?"+params.Encode(), err)
// 		return
// 	}
// 	if res.Code != 0 {
// 		log.Error("UpdatePerson url(%s) error(%v)", updatePerson+"?"+params.Encode(), res.Code)
// 		err = ecode.Int(res.Code)
// 	}
// 	return
// }

// UpdateUname update uname.
// func (d *Dao) UpdateUname(c context.Context, mid int64, ip, uname string, isUpNickFree bool) (err error) {
// 	params := url.Values{}
// 	params.Set("mid", strconv.FormatInt(mid, 10))
// 	params.Set("uname", uname)
// 	params.Set("isupnick_free", strconv.FormatBool(isUpNickFree))
// 	var res struct {
// 		Code int `json:"code"`
// 	}
// 	if err = d.client.Post(c, d.accCo+updateUnameURI, ip, params, &res); err != nil {
// 		log.Error("UpdateUname url(%s) error(%v)", updateUnameURI+"?"+params.Encode(), err)
// 		return
// 	}
// 	if res.Code != 0 {
// 		log.Error("UpdateUname url(%s) code(%v)", updateUnameURI+"?"+params.Encode(), res.Code)
// 		err = ParseJavaCode(res.Code)
// 	}
// 	return
// }

// UpdateFace update face
func (d *Dao) UpdateFace(c context.Context, mid int64, face string, ip, cookie, accessKey string) (faceURL string, err error) {
	var (
		params = url.Values{}
		fu     = url.Values{}
		req    *xhttp.Request
	)
	log.Info("start UpdateFace %d,cookie: %v", mid, cookie)
	params.Set("mid", strconv.FormatInt(mid, 10))
	params.Set("ts", strconv.FormatInt(time.Now().Unix(), 10))
	params.Set("access_key", accessKey)
	var url = d.accCom + updateFaceURI + "?" + d.sign(&params)
	fu.Set("face", face)
	if req, err = xhttp.NewRequest("POST", url, bytes.NewBuffer([]byte(fu.Encode()))); err != nil {
		log.Error("http.NewRequest, err(%v)", err)
		return
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	var res struct {
		Code int `json:"code"`
		Data struct {
			Face string `json:"face"`
		} `json:"data"`
		Message string `json:"message"`
	}
	if err = d.client.Do(c, req, &res); err != nil {
		log.Error("UpdateFace url(%v), err(%v)", url, err)
		return
	}
	if res.Code != 0 {
		log.Error("UpdateFace failed, url(%v) res(%v)", url, res)
		err = ParseJavaCode(res.Code)
		return
	}
	faceURL = res.Data.Face
	log.Info("end success UpdateFace %d, %v", mid, faceURL)
	return
}

// ParseJavaCode is
func ParseJavaCode(code int) (err error) {
	switch code {
	case -617:
		err = ecode.UpdateUnameHadLocked
	case -618:
		err = ecode.UpdateUnameRepeated
	case -655:
		err = ecode.UpdateUnameMoneyIsNot
	case -610:
		err = ecode.MemberBlocked
	case -602:
		err = ecode.UpdateUnameTooShort
	case -603:
		err = ecode.UpdateUnameSensitive
	case -601:
		err = ecode.UpdateUnameTooLong
	case -605:
		err = ecode.UpdateUnameFormat
	case -1001:
		err = ecode.MemberSignHasEmoji
	case -40012:
		err = ecode.UpdateFaceFormat
	case -40013:
		err = ecode.UpdateFaceSize
	default:
		err = ecode.Int(code)
	}
	return
}

func (d *Dao) sign(params *url.Values) (query string) {
	params.Set("appkey", d.c.App.Key)
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	var b bytes.Buffer
	b.WriteString(tmp)
	b.WriteString(d.c.App.Secret)
	mh := md5.Sum(b.Bytes())
	// query
	var qb bytes.Buffer
	qb.WriteString(tmp)
	qb.WriteString("&sign=")
	qb.WriteString(hex.EncodeToString(mh[:]))
	query = qb.String()
	return
}

// IdentifyInfo get identify info by calling api.
func (d *Dao) IdentifyInfo(c context.Context, mid int64, ip string) (idt *model.IdentifyInfo, err error) {
	params := url.Values{}
	params.Set("mid", strconv.FormatInt(mid, 10))
	var res struct {
		Code int                 `json:"code"`
		Data *model.IdentifyInfo `json:"data"`
	}
	err = d.client.Get(c, d.accCo+identifyURI, ip, params, &res)
	if err != nil {
		log.Error("dao.client.Get(%s) error(%v)", d.accCo+identifyURI+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("identifyInfo url(%s) error(%v)", d.accCo+identifyURI+"?"+params.Encode(), err)
		err = ecode.Int(res.Code)
		return
	}
	idt = res.Data
	return
}
