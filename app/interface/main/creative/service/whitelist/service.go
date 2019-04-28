package whitelist

import (
	"context"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/dao/whitelist"
	accmdl "go-common/app/interface/main/creative/model/account"
	"go-common/app/interface/main/creative/service"
	"go-common/library/log"
	"time"
)

const (
	_creator         = 0
	_playerAttention = 2
	_taskMiddleUp    = 3
	_taskSmallUp     = 4
	_taskNotUp       = 5
)

//Service struct
type Service struct {
	c               *conf.Config
	Creator         map[int64]int64
	PlayerAttention map[int64]int64
	list            map[int]map[int64]int64
	wl              *whitelist.Dao
}

//New get service
func New(c *conf.Config, rpcdaos *service.RPCDaos) *Service {
	s := &Service{
		c:  c,
		wl: whitelist.New(c),
	}
	s.loadAll()
	s.loadCreator()
	s.loadSwitch()
	go s.loadproc()
	return s
}

func (s *Service) loadCreator() {
	s.Creator = s.list[_creator]
}

func (s *Service) loadSwitch() {
	s.PlayerAttention = s.list[_playerAttention]
}

func (s *Service) loadAll() {
	wlList, err := s.wl.List(context.TODO())
	if err != nil {
		return
	}
	temp := make(map[int]map[int64]int64)
	for _, v := range wlList {
		if _, ok := temp[v.Tp]; !ok {
			temp[v.Tp] = make(map[int64]int64)
		}
		temp[v.Tp][v.Mid] = v.Mid
	}
	s.list = temp
}

// loadproc
func (s *Service) loadproc() {
	for {
		time.Sleep(5 * time.Minute)
		s.loadAll()
		s.loadCreator()
		s.loadSwitch()
	}
}

// UploadInfoForCreator fn, 判断创作姬的能否进入app投稿的权限
func (s *Service) UploadInfoForCreator(mf *accmdl.MyInfo, mid int64) (uploadinfo map[string]interface{}) {
	uploadinfo = make(map[string]interface{})
	uploadinfo["info"] = 1
	uploadinfo["reason"] = "账号已经过校验，可以投稿。"
	if mf == nil {
		log.Error("accmdl.MyInfo is nil mid(%d)", mid)
		return
	}
	if mf.Banned {
		uploadinfo["info"] = 0
		uploadinfo["reason"] = "您的账号已被禁用，无法投稿。"
	}
	if !mf.Activated {
		uploadinfo["info"] = 0
		uploadinfo["reason"] = "您的账号尚未激活，无法投稿。"
	}
	if mf.IdentifyInfo.Code != 0 {
		uploadinfo["info"] = 0
		uploadinfo["reason"] = "投稿需要进行实名制登记，请先在PC上进行登记。"
	}
	var ok bool
	if _, ok = s.Creator[mf.Mid]; !ok && (mf.Level < 4) {
		uploadinfo["info"] = 0
		uploadinfo["reason"] = "亲爱的用户，应用内测期间仅开放给部分用户。敬请谅解。"
	}
	return
}

// Viewinfo fn, 判断创作姬的能否进入app查看的权限
func (s *Service) Viewinfo(mf *accmdl.MyInfo) (uploadinfo map[string]interface{}) {
	uploadinfo = make(map[string]interface{})
	var ok bool
	uploadinfo["info"] = 0
	uploadinfo["reason"] = "亲爱的用户，应用内测期间仅开放给部分用户。敬请谅解。"
	if _, ok = s.Creator[mf.Mid]; ok {
		uploadinfo["info"] = 1
		uploadinfo["reason"] = ""
	}
	return
}

// UploadInfoForMainApp fn, 判断主APP的能否进入app投稿的权限
func (s *Service) UploadInfoForMainApp(mf *accmdl.MyInfo, plat string, mid int64) (uploadinfo map[string]interface{}, white int) {
	uploadinfo = make(map[string]interface{})
	uploadinfo["info"] = 1
	uploadinfo["reason"] = "账号已经过校验，可以投稿。"
	uploadinfo["url"] = ""
	if mf == nil {
		log.Error("accmdl.MyInfo is nil mid(%d)", mid)
		return
	}
	if mf.Banned {
		uploadinfo["info"] = 0
		uploadinfo["reason"] = "您的账号已被禁用，无法投稿。"
	}
	if !mf.Activated {
		uploadinfo["info"] = 0
		uploadinfo["reason"] = "您的账号尚未激活，无法投稿。"
	}
	if mf.IdentifyInfo.Code != 0 {
		uploadinfo["info"] = 0
		if plat == "android" {
			uploadinfo["reason"] = "投稿需要进行实名认证，请前往“设置”-“账号与隐私”中进行认证"
		} else if plat == "ios" {
			uploadinfo["reason"] = "投稿需要进行实名认证，请前往“设置”-“安全隐私”中进行认证"
		} else {
			uploadinfo["reason"] = "投稿需要进行实名制登记，请先在PC上进行登记。"
		}
		uploadinfo["url"] = s.c.H5Page.Passport
	}
	if mf.Level < 1 {
		uploadinfo["info"] = 0
		uploadinfo["reason"] = "LV1以上用户才能投稿哦，请在头像下方进入答题升级吧～"
	}
	white, _ = uploadinfo["info"].(int)
	return
}

//TaskWhiteList 任务系统白名单 0-关闭 1-开启
func (s *Service) TaskWhiteList(mid int64) (res int8) {
	blackList := []int64{208259, 2, 9099524}
	for _, v := range blackList {
		if v == mid {
			return 0
		}
	}
	// internal whiteList
	if s.c.Whitelist.DataMids != nil {
		for _, m := range s.c.Whitelist.DataMids {
			if m == mid {
				return 1
			}
		}
	}
	if s.c.TaskCondition.WhiteSwitch {
		if s.list != nil {
			if _, ok := s.list[_taskNotUp][mid]; ok {
				return 1
			}
			if _, ok := s.list[_taskSmallUp][mid]; ok {
				return 1
			}
			if _, ok := s.list[_taskMiddleUp][mid]; ok {
				return 1
			}
		}
	} else {
		return 0
	}
	return 0
}
