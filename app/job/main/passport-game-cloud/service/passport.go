package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/passport-game-cloud/model"
	"go-common/library/log"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	_userTableUpdateDuration    = time.Second
	_mySQLErrCodeDuplicateEntry = 1062

	_tokenTableUpdateRetryCount = 3
	_tokenTableUpdateDuration   = time.Second
	_tokenTablePrefix           = "aso_app_perm"

	_tokenCacheRetryCount    = 3
	_tokenCacheRetryDuration = time.Second

	_notifyGameRetryCount    = 3
	_notifyGameRetryDuration = time.Second
)

func (s *Service) processUserInfo(bmsg *model.BMsg) {
	n := new(model.Info)
	if err := json.Unmarshal(bmsg.New, n); err != nil {
		log.Error("failed to parse binlog new, json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
		return
	}
	s.asoAccountInterval.Prom(context.TODO(), bmsg.MTS)
	switch bmsg.Action {
	case "insert":
		s.addMemberInit(context.TODO(), n)
		s.delInfoCache(context.TODO(), n.Mid)
	case "update":
		old := new(model.Info)
		if err := json.Unmarshal(bmsg.Old, old); err != nil {
			log.Error("failed to parse binlog old, json.Unmarshal(%s) error(%v)", string(bmsg.Old), err)
			return
		}
		if n.Equals(old) {
			return
		}
		s.delInfoCache(context.TODO(), n.Mid)
	case "delete":
		s.delInfoCache(context.TODO(), n.Mid)
	}
}

// sub log encryption UPDATE INSERT
func (s *Service) processAsoAccSub(msg *model.PMsg) {
	n := msg.Data
	flag := msg.Flag
	s.transInterval.Prom(context.TODO(), msg.MTS)
	switch msg.Action {
	case "insert":
		s.addAsoAccount(context.TODO(), n)
	case "update":
		s.updateAsoAccount(context.TODO(), n, flag)
	case "delete":
		s.delAsoAccount(context.TODO(), n.Mid)
	}
	s.delInfoCache(context.TODO(), n.Mid)
	s.notifyGame(context.TODO(), n.Mid, "", _updateUserInfo)
}

func (s *Service) processToken(bmsg *model.BMsg) {
	n := new(model.Perm)
	if err := json.Unmarshal(bmsg.New, n); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(bmsg.New), err)
		return
	}
	isGame := false
	for _, id := range s.gameAppIDs {
		if n.AppID == id {
			isGame = true
			break
		}
	}
	if !isGame {
		return
	}
	s.tokenInterval.Prom(context.TODO(), bmsg.MTS)
	switch bmsg.Action {
	case "insert":
		s.addToken(context.TODO(), n)
		s.setTokenCache(context.TODO(), n)
	case "update":
		old := new(model.Perm)
		if err := json.Unmarshal(bmsg.Old, old); err != nil {
			log.Error("failed to parse binlog old, json.Unmarshal(%s) error(%v)", string(bmsg.Old), err)
			return
		}
		if n.Equals(old) {
			return
		}
		s.updateToken(context.TODO(), n)
		s.setTokenCache(context.TODO(), n)
	case "delete":
		s.delToken(context.TODO(), n.AccessToken)
		s.delTokenCache(context.TODO(), n.AccessToken)
	}
}

func (s *Service) addAsoAccount(c context.Context, a *model.AsoAccount) (err error) {
	for {
		_, err = s.d.AddAsoAccount(c, a)
		if err == nil {
			break
		}
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				log.Error("failed to add aso because of duplicate entry, value is (%v), error(%v)", a, err)
				return
			}
		}
		time.Sleep(_userTableUpdateDuration)
	}
	return
}

func (s *Service) updateAsoAccount(c context.Context, a *model.AsoAccount, flag int) (err error) {
	for {
		_, err = s.UpdateAsoAccount(c, a, flag)
		if err == nil {
			break
		}
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				log.Error("failed to update aso because of duplicate entry, value is (%v), error(%v)", a, err)
				return
			}
		}
		time.Sleep(_userTableUpdateDuration)
	}
	return
}

// UpdateAsoAccount update aso account.
func (s *Service) UpdateAsoAccount(c context.Context, t *model.AsoAccount, flag int) (affected int64, err error) {
	if affected, err = s.d.UpdateAsoAccount(c, t); err != nil {
		return
	}
	if flag != 1 {
		return
	}
	mid := t.Mid
	var tokens []string
	if tokens, err = s.d.Tokens(c, mid); err != nil {
		return
	}
	for _, token := range tokens {
		s.delToken(c, token)
		s.delTokenCache(c, token)
		s.notifyGame(c, mid, token, _changePwd)
	}
	return
}

func (s *Service) delAsoAccount(c context.Context, mid int64) (err error) {
	for {
		if _, err = s.d.DelAsoAccount(c, mid); err == nil {
			break
		}
		time.Sleep(_userTableUpdateDuration)
	}
	return
}

func (s *Service) addMemberInit(c context.Context, a *model.Info) (err error) {
	for {
		memberInfo := &model.Info{
			Mid:   a.Mid,
			Uname: a.Uname,
			Face:  "",
		}
		err = s.addMemberInfo(c, memberInfo)
		if err == nil {
			break
		}
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				log.Error("failed to add member because of duplicate entry, error(%v)", err)
				return
			}
		}
		time.Sleep(_userTableUpdateDuration)
	}
	return
}

func (s *Service) setTokenCache(c context.Context, t *model.Perm) (err error) {
	for i := 0; i < _tokenCacheRetryCount; i++ {
		if err = s.d.SetTokenCache(c, t); err == nil {
			break
		}
		time.Sleep(_tokenCacheRetryDuration)
	}
	return
}

func (s *Service) delTokenCache(c context.Context, accessToken string) (err error) {
	for i := 0; i < _tokenCacheRetryCount; i++ {
		if err = s.d.DelTokenCache(c, accessToken); err == nil {
			break
		}
		time.Sleep(_tokenCacheRetryDuration)
	}
	return
}

func (s *Service) addToken(c context.Context, t *model.Perm) (err error) {
	for i := 0; i < _tokenTableUpdateRetryCount; i++ {
		_, err = s.d.AddToken(c, t)
		if err == nil {
			return
		}
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				log.Error("failed to add token because of duplicate entry, error(%v)", err)
				return
			}
		}
		time.Sleep(_tokenTableUpdateDuration)
	}
	return
}

func (s *Service) updateToken(c context.Context, t *model.Perm) (err error) {
	for i := 0; i < _tokenTableUpdateRetryCount; i++ {
		if _, err = s.d.UpdateToken(c, t); err == nil {
			return
		}
		time.Sleep(_tokenTableUpdateDuration)
	}
	return
}

func (s *Service) delToken(c context.Context, accessToken string) (err error) {
	for i := 0; i < _tokenTableUpdateRetryCount; i++ {
		if _, err = s.d.DelToken(c, accessToken); err == nil {
			return
		}
		time.Sleep(_tokenTableUpdateDuration)
	}
	return
}

func (s *Service) notifyGame(c context.Context, mid int64, token, action string) (err error) {
	for i := 0; i < _notifyGameRetryCount; i++ {
		if err = s.d.NotifyGame(c, mid, token, action); err == nil {
			return
		}
		time.Sleep(_notifyGameRetryDuration)
	}
	return
}
