package service

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"

	pb "go-common/app/service/main/sms/api"
	"go-common/app/service/main/sms/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const mobilePattern = "^((13[0-9])|(14[1,4,5,6,7,8])|(15[^4])|(16[5-7])|(17[0-8])|(18[0-9])|(19[1,8,9]))\\d{8}$"

var mobileReg, _ = regexp.Compile(mobilePattern)

// Send send sms
func (s *Service) Send(ctx context.Context, req *pb.SendReq) (res *pb.SendReply, err error) {
	tpl := s.template[req.Tcode]
	if tpl == nil {
		err = ecode.SmsTemplateNotExist
		return
	}
	if req.Mid != 0 && req.Mobile != "" {
		err = ecode.SmsSendBothMidAndMobile
		return
	}
	if req.Mid == 0 && req.Mobile == "" {
		err = ecode.RequestErr
		return
	}
	if req.Mobile != "" {
		if req.Country == "" {
			req.Country = model.CountryChina
		}
		if req.Country == model.CountryChina {
			if match := mobileReg.MatchString(req.Mobile); !match {
				err = ecode.SmsMobilePatternErr
				return
			}
		}
	}
	send := &model.ModelSend{
		Mid:     strconv.FormatInt(req.Mid, 10),
		Mobile:  req.Mobile,
		Type:    tpl.Stype,
		Country: req.Country,
		Code:    tpl.Code,
		Content: tpl.Template,
	}
	m := make(map[string]interface{})
	if req.Tparam != "" {
		if err = json.Unmarshal([]byte(req.Tparam), &m); err != nil {
			log.Error("json.Unmarshal (%v) error(%v)", req.Tparam, err)
			return
		}
	}
	for _, k := range tpl.Param {
		if m[k] == "" || m[k] == nil {
			err = ecode.SmsTemplateParamNotEnough
			return
		}
		var v string
		switch m[k].(type) {
		case string:
			v = m[k].(string)
		case float64:
			v = strconv.FormatFloat(m[k].(float64), 'f', -1, 64)
		default:
			err = ecode.SmsTemplateParamIllegal
			return
		}
		send.Content = strings.Replace(send.Content, "#["+k+"]", v, -1)
	}
	err = s.dao.PubSingle(ctx, send)
	return
}

// SendBatch send sms batch
func (s *Service) SendBatch(ctx context.Context, req *pb.SendBatchReq) (res *pb.SendBatchReply, err error) {
	var (
		tpl = s.template[req.Tcode]
		mbs = make([]string, 0)
	)
	if tpl == nil {
		err = ecode.SmsTemplateNotExist
		return
	}
	if tpl.Stype != model.TypeActSms {
		err = ecode.SmsTemplateNotAct
		return
	}
	if len(req.Mids) == 0 && len(req.Mobiles) == 0 {
		err = ecode.RequestErr
		return
	}
	if len(req.Mids) > 0 && len(req.Mobiles) > 0 {
		err = ecode.SmsSendBothMidAndMobile
		return
	}
	if len(req.Mids)+len(req.Mobiles) > 100 {
		err = ecode.SmsSendBatchOverLimit
		return
	}
	for _, v := range req.Mobiles {
		if match := mobileReg.MatchString(v); match {
			mbs = append(mbs, v)
		}
	}
	if len(req.Mobiles) > 0 && len(mbs) == 0 {
		err = ecode.SmsMobilePatternErr
		return
	}
	m := make(map[string]interface{})
	if req.Tparam != "" {
		if err = json.Unmarshal([]byte(req.Tparam), &m); err != nil {
			log.Error("json.Unmarshal (%v) error(%v)", string(req.Tparam), err)
		}
	}
	send := &model.ModelSend{Type: model.TypeActBatch, Code: tpl.Code, Content: tpl.Template}
	for _, k := range tpl.Param {
		if m[k] == "" || m[k] == nil {
			err = ecode.SmsTemplateParamNotEnough
			return
		}
		var v string
		switch m[k].(type) {
		case string:
			v = m[k].(string)
		case float64:
			v = strconv.FormatFloat(m[k].(float64), 'f', -1, 64)
		default:
			err = ecode.SmsTemplateParamIllegal
			return
		}
		send.Content = strings.Replace(send.Content, "#["+k+"]", v, -1)
	}
	send.Mid = xstr.JoinInts(req.Mids)
	send.Mobile = strings.Join(mbs, ",")
	err = s.dao.PubBatch(ctx, send)
	return
}
