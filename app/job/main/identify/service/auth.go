package service

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"

	"go-common/app/job/main/identify/model"
	mdl "go-common/app/service/main/identify/model"
	"go-common/library/log"
	"go-common/library/queue/databus"
)

func (s *Service) new(msg *databus.Message) (interface{}, error) {
	bmsg := new(model.BMsg)
	if err := json.Unmarshal(msg.Value, bmsg); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(msg.Value), err)
		return nil, err
	}
	if bmsg.Action != "delete" {
		return bmsg, nil
	}
	if strings.HasPrefix(bmsg.Table, "user_token_") {
		t := new(model.AuthToken)
		if err := json.Unmarshal(bmsg.New, t); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
			return nil, err
		}
		log.Info("identifyconsumeproc table:%s key:%s partition:%d offset:%d", bmsg.Table, msg.Key, msg.Partition, msg.Offset)
		return t, nil
	} else if strings.HasPrefix(bmsg.Table, "user_cookie_") {
		t := new(model.AuthCookie)
		if err := json.Unmarshal(bmsg.New, t); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
			return nil, err
		}
		log.Info("identifyconsumeproc table:%s key:%s partition:%d offset:%d", bmsg.Table, msg.Key, msg.Partition, msg.Offset)
		return t, nil
	}
	return bmsg, nil
}

func (s *Service) spilt(msg *databus.Message, data interface{}) int {
	switch t := data.(type) {
	case *model.AuthToken:
		return int(t.Mid)
	case *model.AuthCookie:
		return int(t.Mid)
	default:
		return 0
	}
}

func (s *Service) processAuthBinlog2(bmsgs []interface{}) {
	for _, msg := range bmsgs {
		switch t := msg.(type) {
		case *model.AuthToken:
			var (
				bytes []byte
				err   error
			)
			if bytes, err = base64.StdEncoding.DecodeString(t.Token); err != nil {
				log.Error("cleanCookieCache base64 decode err %v", err)
				err = nil
				return
			}
			info := &mdl.IdentifyInfo{
				Mid:     t.Mid,
				Expires: int32(t.Expires),
			}

			if ok := isGameAppID(t.AppID); ok {
				continue
			}
			for {
				log.Info("auth service process token databus, key(%s)", hex.EncodeToString(bytes))
				if err := s.processIdentify("delete", hex.EncodeToString(bytes), info); err != nil {
					continue
				}
				break
			}
		case *model.AuthCookie:
			var (
				bytes     []byte
				bytesCSRF []byte
				err       error
			)
			if bytes, err = base64.StdEncoding.DecodeString(t.Session); err != nil {
				log.Error("cleanCookieCache base64 decode err %v", err)
				err = nil
				return
			}
			if bytesCSRF, err = base64.StdEncoding.DecodeString(t.CSRF); err != nil {
				log.Error("cleanCookieCache base64 decode err %v", err)
				err = nil
				return
			}
			info := &mdl.IdentifyInfo{
				Mid:     t.Mid,
				Csrf:    hex.EncodeToString(bytesCSRF),
				Expires: int32(t.Expires),
			}

			for {
				log.Info("auth service process cookie databus, key(%s)", string(bytes))
				if err := s.processIdentify("delete", string(bytes), info); err != nil {
					continue
				}
				break
			}
		default:
			return
		}
	}
}
