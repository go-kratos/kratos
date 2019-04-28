package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"go-common/app/job/main/passport-sns/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

const (
	_asoAccountSnsTable         = "aso_account_sns"
	_insertAction               = "insert"
	_updateAction               = "update"
	_deleteAction               = "delete"
	_mySQLErrCodeDuplicateEntry = 1062
)

type asoAccountSnsBMsg struct {
	Action    string
	Table     string
	New       *model.AsoAccountSns
	Old       *model.AsoAccountSns
	Timestamp int64
}

func (s *Service) asoBinLogConsume() {
	s.group.New = func(msg *databus.Message) (res interface{}, err error) {
		bmsg := new(model.BMsg)
		if err = json.Unmarshal(msg.Value, bmsg); err != nil {
			log.Error("json.Unmarshal(%s) error(%+v)", string(msg.Value), err)
			return
		}
		log.Info("receive msg action(%s) table(%s) key(%s) partition(%d) offset(%d) timestamp(%d) New(%s) Old(%s)",
			bmsg.Action, bmsg.Table, msg.Key, msg.Partition, msg.Offset, msg.Timestamp, string(bmsg.New), string(bmsg.Old))
		if bmsg.Table == _asoAccountSnsTable {
			asoAccountSnsBMsg := &asoAccountSnsBMsg{
				Action:    bmsg.Action,
				Table:     bmsg.Table,
				Timestamp: msg.Timestamp,
			}
			newAccountSns := new(model.AsoAccountSns)
			if err = json.Unmarshal(bmsg.New, newAccountSns); err != nil {
				log.Error("json.Unmarshal(%s) error(%+v)", string(bmsg.New), err)
				return
			}
			asoAccountSnsBMsg.New = newAccountSns
			if bmsg.Action == _updateAction {
				oldAccountSns := new(model.AsoAccountSns)
				if err = json.Unmarshal(bmsg.Old, oldAccountSns); err != nil {
					log.Error("json.Unmarshal(%s) error(%+v)", string(bmsg.Old), err)
					return
				}
				asoAccountSnsBMsg.Old = oldAccountSns
			}
			return asoAccountSnsBMsg, nil
		}
		return
	}
	s.group.Split = func(msg *databus.Message, data interface{}) int {
		if t, ok := data.(*asoAccountSnsBMsg); ok {
			return int(t.New.Mid)
		}
		return 0
	}
	s.group.Do = func(msgs []interface{}) {
		for _, m := range msgs {
			if msg, ok := m.(*asoAccountSnsBMsg); ok {
				for {
					if err := s.handleAsoAccountSns(msg); err != nil {
						log.Error("fail to handleAsoAccountSns msg(%+v) new(%+v) error(%+v)", msg, msg.New, err)
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}
			}
		}
	}
	// start the group
	s.group.Start()
	log.Info("s.group.Start()")
}

func (s *Service) handleAsoAccountSns(msg *asoAccountSnsBMsg) (err error) {
	switch msg.Action {
	case _insertAction:
		if msg.New.QQOpenid != "" {
			if err = s.addSnsQQ(msg.New); err != nil {
				return
			}
		}
		if msg.New.SinaUID != 0 {
			if err = s.addSnsWeibo(msg.New); err != nil {
				return
			}
		}
	case _updateAction:
		if msg.New.QQOpenid != msg.Old.QQOpenid {
			if msg.New.QQOpenid == "" {
				if _, err = s.d.DelSnsUser(context.Background(), msg.New.Mid, model.PlatformQQ); err != nil {
					return
				}
				s.cache.Do(context.Background(), func(c context.Context) {
					s.d.DelSnsCache(c, msg.New.Mid, model.PlatformQQStr)
				})
			}
			if msg.New.QQOpenid != "" {
				var user *model.SnsUser
				if user, err = s.d.SnsUserByMid(context.Background(), msg.New.Mid, model.PlatformQQ); err != nil {
					return
				}
				if user == nil {
					if err = s.addSnsQQ(msg.New); err != nil {
						return
					}
				} else {
					if err = s.updateSnsQQ(msg.New); err != nil {
						return
					}
				}
			}
			return
		}
		if msg.New.SinaUID != msg.Old.SinaUID {
			if msg.New.SinaUID == 0 {
				if _, err = s.d.DelSnsUser(context.Background(), msg.New.Mid, model.PlatformWEIBO); err != nil {
					return
				}
				s.cache.Do(context.Background(), func(c context.Context) {
					s.d.DelSnsCache(c, msg.New.Mid, model.PlatformWEIBOStr)
				})
			}
			if msg.New.SinaUID != 0 {
				var user *model.SnsUser
				if user, err = s.d.SnsUserByMid(context.Background(), msg.New.Mid, model.PlatformWEIBO); err != nil {
					return
				}
				if user == nil {
					if err = s.addSnsWeibo(msg.New); err != nil {
						return
					}
				} else {
					if err = s.updateSnsWeibo(msg.New); err != nil {
						return
					}
				}
			}
			return
		}
		if msg.New.SinaAccessExpires != msg.Old.SinaAccessExpires {
			return s.updateSns(msg.New.Mid, msg.New.SinaAccessExpires, strconv.FormatInt(msg.New.SinaUID, 10), msg.New.SinaAccessToken, model.PlatformWEIBO)
		}
		if msg.New.QQAccessExpires != msg.Old.QQAccessExpires {
			var qqUnionID string
			if qqUnionID, err = s.d.QQUnionID(context.Background(), msg.New.QQOpenid); err != nil {
				return
			}
			if qqUnionID == "" {
				log.Error("update qq expires, qqUnionID is null, oldMsg(%+v) newMsg(%+v)", msg.Old, msg.New)
				return
			}
			return s.updateSns(msg.New.Mid, msg.New.QQAccessExpires, qqUnionID, msg.New.QQAccessToken, model.PlatformQQ)
		}
	case _deleteAction:
		if _, err = s.d.DelSnsUser(context.Background(), msg.New.Mid, model.PlatformQQ); err != nil {
			return
		}
		if _, err = s.d.DelSnsUser(context.Background(), msg.New.Mid, model.PlatformWEIBO); err != nil {
			return
		}
		s.cache.Do(context.Background(), func(c context.Context) {
			s.d.DelSnsCache(c, msg.New.Mid, model.PlatformQQStr)
		})
		s.cache.Do(context.Background(), func(c context.Context) {
			s.d.DelSnsCache(c, msg.New.Mid, model.PlatformWEIBOStr)
		})
	}
	return
}

func (s *Service) updateSns(mid, expires int64, unionID, token string, platform int) (err error) {
	tx, err := s.d.BeginSnsTran(context.Background())
	if err != nil {
		log.Error("s.d.BeginTran error(%+v)", err)
		return
	}

	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()

	snsUser := &model.SnsUser{
		Mid:      mid,
		UnionID:  unionID,
		Platform: platform,
		Expires:  expires,
	}
	if _, err = s.d.TxUpdateSnsUserExpires(tx, snsUser); err != nil {
		return
	}

	snsToken := &model.SnsToken{
		Mid:      mid,
		Platform: platform,
		Token:    token,
		Expires:  expires,
	}
	if _, err = s.d.TxUpdateSnsToken(tx, snsToken); err != nil {
		return
	}
	s.cache.Do(context.Background(), func(c context.Context) {
		proto := &model.SnsProto{
			Mid:      snsUser.Mid,
			Platform: int32(platform),
			UnionID:  snsUser.UnionID,
			Expires:  snsUser.Expires,
		}
		s.d.SetSnsCache(c, mid, parsePlatformStr(platform), proto)
	})
	return
}

func (s *Service) addSnsQQ(sns *model.AsoAccountSns) (err error) {
	unionID, err := s.d.GetUnionIDCache(context.Background(), sns.QQOpenid)
	if err != nil || unionID == "" {
		unionID, err = s.d.QQUnionID(context.Background(), sns.QQOpenid)
		if err != nil || unionID == "" {
			return
		}
	}

	tx, err := s.d.BeginSnsTran(context.Background())
	if err != nil {
		log.Error("s.d.BeginTran error(%+v)", err)
		return
	}

	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()

	snsUser := &model.SnsUser{
		Mid:      sns.Mid,
		UnionID:  unionID,
		Platform: model.PlatformQQ,
		Expires:  sns.QQAccessExpires,
	}
	if _, err = s.d.TxAddSnsUser(tx, snsUser); err != nil {
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				u, e := s.d.SnsUserByMid(context.Background(), snsUser.Mid, snsUser.Platform)
				if e != nil {
					return e
				}
				err = nil
				if u != nil && u.Mid == snsUser.Mid {
					return
				}
				log.Error("add sns qq duplicate (%+v)", snsUser)
				return
			}
		}
		return
	}

	snsOpenID := &model.SnsOpenID{
		Mid:      sns.Mid,
		OpenID:   sns.QQOpenid,
		UnionID:  unionID,
		AppID:    model.OldAppID,
		Platform: model.PlatformQQ,
	}
	if _, err = s.d.TxAddSnsOpenID(tx, snsOpenID); err != nil {
		return
	}

	snsToken := &model.SnsToken{
		Mid:      sns.Mid,
		OpenID:   sns.QQOpenid,
		UnionID:  unionID,
		Platform: model.PlatformQQ,
		Token:    sns.QQAccessToken,
		Expires:  sns.QQAccessExpires,
		AppID:    model.OldAppID,
	}
	if _, err = s.d.TxAddSnsToken(tx, snsToken); err != nil {
		return
	}

	proto := &model.SnsProto{
		Mid:      snsUser.Mid,
		Platform: model.PlatformQQ,
		UnionID:  snsUser.UnionID,
		Expires:  snsUser.Expires,
	}
	s.cache.Do(context.Background(), func(c context.Context) {
		s.d.SetSnsCache(c, sns.Mid, model.PlatformQQStr, proto)
	})
	return
}

func (s *Service) addSnsWeibo(sns *model.AsoAccountSns) (err error) {
	tx, err := s.d.BeginSnsTran(context.Background())
	if err != nil {
		log.Error("s.d.BeginTran error(%+v)", err)
		return
	}

	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()

	snsUser := &model.SnsUser{
		Mid:      sns.Mid,
		UnionID:  strconv.FormatInt(sns.SinaUID, 10),
		Platform: model.PlatformWEIBO,
		Expires:  sns.SinaAccessExpires,
	}
	if _, err = s.d.TxAddSnsUser(tx, snsUser); err != nil {
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				u, e := s.d.SnsUserByMid(context.Background(), snsUser.Mid, snsUser.Platform)
				if e != nil {
					return e
				}
				err = nil
				if u != nil && u.Mid == snsUser.Mid {
					return
				}
				log.Error("add sns weibo duplicate (%+v)", snsUser)
				return
			}
		}
		return
	}

	snsOpenID := &model.SnsOpenID{
		Mid:      sns.Mid,
		OpenID:   strconv.FormatInt(sns.SinaUID, 10),
		UnionID:  strconv.FormatInt(sns.SinaUID, 10),
		AppID:    model.OldAppID,
		Platform: model.PlatformWEIBO,
	}
	if _, err = s.d.TxAddSnsOpenID(tx, snsOpenID); err != nil {
		return
	}

	snsToken := &model.SnsToken{
		Mid:      sns.Mid,
		OpenID:   strconv.FormatInt(sns.SinaUID, 10),
		UnionID:  strconv.FormatInt(sns.SinaUID, 10),
		Platform: model.PlatformWEIBO,
		Token:    sns.SinaAccessToken,
		Expires:  sns.SinaAccessExpires,
		AppID:    model.OldAppID,
	}
	if _, err = s.d.TxAddSnsToken(tx, snsToken); err != nil {
		return
	}

	proto := &model.SnsProto{
		Mid:      snsUser.Mid,
		Platform: model.PlatformWEIBO,
		UnionID:  snsUser.UnionID,
		Expires:  snsUser.Expires,
	}
	s.cache.Do(context.Background(), func(c context.Context) {
		s.d.SetSnsCache(c, sns.Mid, model.PlatformWEIBOStr, proto)
	})
	return
}

func (s *Service) updateSnsQQ(sns *model.AsoAccountSns) (err error) {
	unionID, err := s.d.GetUnionIDCache(context.Background(), sns.QQOpenid)
	if err != nil || unionID == "" {
		unionID, err = s.d.QQUnionID(context.Background(), sns.QQOpenid)
		if err != nil || unionID == "" {
			return
		}
	}

	tx, err := s.d.BeginSnsTran(context.Background())
	if err != nil {
		log.Error("s.d.BeginTran error(%+v)", err)
		return
	}

	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()

	if _, err = s.d.TxUpdateSnsUser(tx, sns.Mid, sns.QQAccessExpires, unionID, model.PlatformQQ); err != nil {
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				u, e := s.d.SnsUserByMid(context.Background(), sns.Mid, model.PlatformQQ)
				if e != nil {
					return e
				}
				err = nil
				if u != nil && u.Mid == sns.Mid {
					return
				}
				log.Error("update sns qq duplicate, mid(%d) unionID(%s)", sns.Mid, unionID)
				return
			}
		}
		return
	}

	snsToken := &model.SnsToken{
		Mid:      sns.Mid,
		OpenID:   sns.QQOpenid,
		UnionID:  unionID,
		Platform: model.PlatformQQ,
		Token:    sns.QQAccessToken,
		Expires:  sns.QQAccessExpires,
		AppID:    model.OldAppID,
	}
	if _, err = s.d.TxAddSnsToken(tx, snsToken); err != nil {
		return
	}

	proto := &model.SnsProto{
		Mid:      sns.Mid,
		Platform: model.PlatformQQ,
		UnionID:  unionID,
		Expires:  sns.QQAccessExpires,
	}
	s.cache.Do(context.Background(), func(c context.Context) {
		s.d.SetSnsCache(c, sns.Mid, model.PlatformQQStr, proto)
	})
	return
}

func (s *Service) updateSnsWeibo(sns *model.AsoAccountSns) (err error) {
	tx, err := s.d.BeginSnsTran(context.Background())
	if err != nil {
		log.Error("s.d.BeginTran error(%+v)", err)
		return
	}

	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("tx.Rollback() error(%v)", err1)
			}
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
		}
	}()

	if _, err = s.d.TxUpdateSnsUser(tx, sns.Mid, sns.SinaAccessExpires, strconv.FormatInt(sns.SinaUID, 10), model.PlatformWEIBO); err != nil {
		switch nErr := errors.Cause(err).(type) {
		case *mysql.MySQLError:
			if nErr.Number == _mySQLErrCodeDuplicateEntry {
				u, e := s.d.SnsUserByMid(context.Background(), sns.Mid, model.PlatformWEIBO)
				if e != nil {
					return e
				}
				err = nil
				if u != nil && u.Mid == sns.Mid {
					return
				}
				log.Error("update sns weibo duplicate, mid(%d) unionID(%s)", sns.Mid, strconv.FormatInt(sns.SinaUID, 10))
				return
			}
		}
		return
	}

	snsToken := &model.SnsToken{
		Mid:      sns.Mid,
		OpenID:   strconv.FormatInt(sns.SinaUID, 10),
		UnionID:  strconv.FormatInt(sns.SinaUID, 10),
		Platform: model.PlatformWEIBO,
		Token:    sns.SinaAccessToken,
		Expires:  sns.SinaAccessExpires,
		AppID:    model.OldAppID,
	}
	if _, err = s.d.TxAddSnsToken(tx, snsToken); err != nil {
		return
	}

	proto := &model.SnsProto{
		Mid:      sns.Mid,
		Platform: model.PlatformWEIBO,
		UnionID:  strconv.FormatInt(sns.SinaUID, 10),
		Expires:  sns.SinaAccessExpires,
	}
	s.cache.Do(context.Background(), func(c context.Context) {
		s.d.SetSnsCache(c, sns.Mid, model.PlatformWEIBOStr, proto)
	})
	return
}

func (s *Service) fullSyncSns() {
	var (
		start   int64
		chanNum = int64(s.c.SyncConf.ChanNum)
	)
	for {
		log.Info("fullSyncSns, start %d", start)
		res, err := s.d.AsoAccountSns(context.Background(), start)
		if err != nil {
			log.Error("fail to get AsoAccountSns error(%+v)", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		for _, a := range res {
			s.snsChan[a.Mid%chanNum] <- a
		}
		if len(res) == 0 {
			log.Info("fullSyncSns finished! endID(%d)", start)
			break
		}
		start = res[len(res)-1].Mid
	}
}

func (s *Service) fullSyncSnsConsume(c chan *model.AsoAccountSns) {
	for {
		a, ok := <-c
		if !ok {
			log.Error("snsChan closed")
			return
		}
		for {
			if a.SinaUID != 0 {
				if err := s.addSnsWeibo(a); err != nil {
					continue
				}
			}
			if a.QQOpenid != "" {
				if err := s.addSnsQQ(a); err != nil {
					continue
				}
			}
			break
		}
	}
}

func (s *Service) checkAll() {
	ticker := time.NewTicker(time.Duration(s.c.SyncConf.CheckTicker))
	for {
		s.checkSnsUser()
		<-ticker.C
	}
}

func (s *Service) checkConsume(c chan *model.AsoAccountSns) {
	for {
		a, ok := <-c
		if !ok {
			log.Error("snsChan closed")
			return
		}
		for {
			if err := s.checkWeibo(a); err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if err := s.checkQQ(a); err != nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}
			break
		}
	}
}

func (s *Service) checkSnsUser() {
	var (
		start   int64
		chanNum = int64(s.c.SyncConf.ChanNum)
	)
	for {
		log.Info("checkSnsUser, start %d", start)
		res, err := s.d.AsoAccountSnsAll(context.Background(), start)
		if err != nil {
			log.Error("fail to get AsoAccountSns error(%+v)", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		for _, a := range res {
			s.checkChan[a.Mid%chanNum] <- a
		}
		if len(res) == 0 {
			log.Info("checkSnsUser finished! endID(%d)", start)
			break
		}
		start = res[len(res)-1].Mid
	}
}

func (s *Service) checkQQ(a *model.AsoAccountSns) (err error) {
	sns, err := s.d.SnsUserByMid(context.Background(), a.Mid, model.PlatformQQ)
	if err != nil {
		return
	}
	if a.QQOpenid == "" {
		if sns != nil {
			log.Error("qq not match, old is nil,old(%+v) new(%+v)", a, sns)
			if _, err = s.d.DelSnsUser(context.Background(), a.Mid, model.PlatformQQ); err != nil {
				return
			}
			s.cache.Do(context.Background(), func(c context.Context) {
				s.d.DelSnsCache(c, a.Mid, model.PlatformQQStr)
			})
		}
		return
	}
	qqUnionID, err := s.d.QQUnionID(context.Background(), a.QQOpenid)
	if err != nil || qqUnionID == "" {
		return
	}
	if sns == nil {
		log.Error("qq not match, new not exists, old(%+v), ", a)
		if _, err = s.d.AddSnsUser(context.Background(), a.Mid, a.QQAccessExpires, qqUnionID, model.PlatformQQ); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == _mySQLErrCodeDuplicateEntry {
					log.Error("checkQQ, add qq duplicate, mid(%d) unionid(%s)", a.Mid, qqUnionID)
					err = nil
					//var (
					//	u   *model.SnsUser
					//	aso *model.AsoAccountSns
					//)
					//if u, err = s.d.SnsUserByUnionID(context.Background(), qqUnionID, model.PlatformQQ); err != nil {
					//	return
					//}
					//if _, err = s.d.DelSnsUser(context.Background(), u.Mid, model.PlatformQQ); err != nil {
					//	return
					//}
					//if _, err = s.d.AddSnsUser(context.Background(), a.Mid, a.QQAccessExpires, qqUnionID, model.PlatformQQ); err != nil {
					//	return
					//}
					//if aso, err = s.d.AsoAccountSnsByMid(context.Background(), u.Mid); err != nil {
					//	return
					//}
					//if qqUnionID, err = s.d.QQUnionID(context.Background(), a.QQOpenid); err != nil {
					//	return
					//}
					//if _, err = s.d.AddSnsUser(context.Background(), aso.Mid, aso.QQAccessExpires, qqUnionID, model.PlatformQQ); err != nil {
					//	return
					//}
				}
			}
			return
		}
	} else if sns.UnionID != qqUnionID {
		log.Error("qq not match, new(%s), mid(%d) oldOpenID(%s) oldUnionID(%s) ", sns.UnionID, a.Mid, a.QQOpenid, qqUnionID)
		if _, err = s.d.UpdateSnsUser(context.Background(), a.Mid, a.QQAccessExpires, qqUnionID, model.PlatformQQ); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == _mySQLErrCodeDuplicateEntry {
					log.Error("checkQQ, update qq duplicate, mid(%d) unionid(%s)", a.Mid, qqUnionID)
					err = nil
					//var (
					//	u   *model.SnsUser
					//	aso *model.AsoAccountSns
					//)
					//if u, err = s.d.SnsUserByUnionID(context.Background(), qqUnionID, model.PlatformQQ); err != nil {
					//	return
					//}
					//if _, err = s.d.DelSnsUser(context.Background(), u.Mid, model.PlatformQQ); err != nil {
					//	return
					//}
					//if _, err = s.d.UpdateSnsUser(context.Background(), a.Mid, a.QQAccessExpires, qqUnionID, model.PlatformQQ); err != nil {
					//	return
					//}
					//if aso, err = s.d.AsoAccountSnsByMid(context.Background(), u.Mid); err != nil {
					//	return
					//}
					//if qqUnionID, err = s.d.QQUnionID(context.Background(), a.QQOpenid); err != nil {
					//	return
					//}
					//if _, err = s.d.AddSnsUser(context.Background(), aso.Mid, aso.QQAccessExpires, qqUnionID, model.PlatformQQ); err != nil {
					//	return
					//}
				}
			}
			return
		}
	} else if sns.Expires != a.QQAccessExpires {
		log.Error("qq expires not match, new(%+v), old(%+v)", sns, a)
		if err = s.updateSns(a.Mid, a.QQAccessExpires, a.QQOpenid, a.QQAccessToken, model.PlatformQQ); err != nil {
			return
		}
		return
	}
	s.cache.Do(context.Background(), func(c context.Context) {
		proto := &model.SnsProto{
			Mid:      a.Mid,
			Platform: int32(model.PlatformQQ),
			UnionID:  qqUnionID,
			Expires:  a.QQAccessExpires,
		}
		s.d.SetSnsCache(c, a.Mid, model.PlatformQQStr, proto)
	})
	return
}

func (s *Service) checkWeibo(a *model.AsoAccountSns) (err error) {
	sns, err := s.d.SnsUserByMid(context.Background(), a.Mid, model.PlatformWEIBO)
	if err != nil {
		return
	}
	if a.SinaUID == 0 {
		if sns != nil {
			log.Error("weibo not match, old is nil,old(%+v) new(%+v)", a, sns)
			if _, err = s.d.DelSnsUser(context.Background(), a.Mid, model.PlatformWEIBO); err != nil {
				return
			}
			s.cache.Do(context.Background(), func(c context.Context) {
				s.d.DelSnsCache(c, a.Mid, model.PlatformWEIBOStr)
			})
		}
		return
	}
	if sns == nil {
		log.Error("weibo not match, new not exists, old(%+v)", a)
		if _, err = s.d.AddSnsUser(context.Background(), a.Mid, a.SinaAccessExpires, strconv.FormatInt(a.SinaUID, 10), model.PlatformWEIBO); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == _mySQLErrCodeDuplicateEntry {
					log.Error("checkWeibo, add weibo duplicate, mid(%d) unionid(%s)", a.Mid, strconv.FormatInt(a.SinaUID, 10))
					err = nil
					//var (
					//	u   *model.SnsUser
					//	aso *model.AsoAccountSns
					//)
					//if u, err = s.d.SnsUserByUnionID(context.Background(), strconv.FormatInt(a.SinaUID, 10), model.PlatformWEIBO); err != nil {
					//	return
					//}
					//if _, err = s.d.DelSnsUser(context.Background(), u.Mid, model.PlatformWEIBO); err != nil {
					//	return
					//}
					//if _, err = s.d.AddSnsUser(context.Background(), a.Mid, a.SinaAccessExpires, strconv.FormatInt(a.SinaUID, 10), model.PlatformWEIBO); err != nil {
					//	return
					//}
					//if aso, err = s.d.AsoAccountSnsByMid(context.Background(), u.Mid); err != nil {
					//	return
					//}
					//if _, err = s.d.AddSnsUser(context.Background(), aso.Mid, aso.SinaAccessExpires, strconv.FormatInt(aso.SinaUID, 10), model.PlatformWEIBO); err != nil {
					//	return
					//}
				}
			}
			return
		}
	} else if sns.UnionID != strconv.FormatInt(a.SinaUID, 10) {
		log.Error("weibo not match, new(%s), mid(%d) old(%d)", sns.UnionID, a.Mid, a.SinaUID)
		if _, err = s.d.UpdateSnsUser(context.Background(), a.Mid, a.SinaAccessExpires, strconv.FormatInt(a.SinaUID, 10), model.PlatformWEIBO); err != nil {
			switch nErr := errors.Cause(err).(type) {
			case *mysql.MySQLError:
				if nErr.Number == _mySQLErrCodeDuplicateEntry {
					log.Error("checkWeibo, update weibo duplicate, mid(%d) unionid(%s)", a.Mid, strconv.FormatInt(a.SinaUID, 10))
					err = nil
					//var (
					//	u   *model.SnsUser
					//	aso *model.AsoAccountSns
					//)
					//if u, err = s.d.SnsUserByUnionID(context.Background(), strconv.FormatInt(a.SinaUID, 10), model.PlatformWEIBO); err != nil {
					//	return
					//}
					//if _, err = s.d.DelSnsUser(context.Background(), u.Mid, model.PlatformWEIBO); err != nil {
					//	return
					//}
					//if _, err = s.d.UpdateSnsUser(context.Background(), a.Mid, a.SinaAccessExpires, strconv.FormatInt(a.SinaUID, 10), model.PlatformWEIBO); err != nil {
					//	return
					//}
					//if aso, err = s.d.AsoAccountSnsByMid(context.Background(), u.Mid); err != nil {
					//	return
					//}
					//if _, err = s.d.AddSnsUser(context.Background(), aso.Mid, aso.SinaAccessExpires, strconv.FormatInt(aso.SinaUID, 10), model.PlatformWEIBO); err != nil {
					//	return
					//}
				}
			}
			return
		}
	} else if sns.Expires != a.SinaAccessExpires {
		log.Error("weibo expires not match, new(%+v), old(%+v)", sns, a)
		if err = s.updateSns(a.Mid, a.SinaAccessExpires, strconv.FormatInt(a.SinaUID, 10), a.SinaAccessToken, model.PlatformWEIBO); err != nil {
			return
		}
		return
	}
	s.cache.Do(context.Background(), func(c context.Context) {
		proto := &model.SnsProto{
			Mid:      a.Mid,
			Platform: int32(model.PlatformWEIBO),
			UnionID:  strconv.FormatInt(a.SinaUID, 10),
			Expires:  a.SinaAccessExpires,
		}
		s.d.SetSnsCache(c, a.Mid, model.PlatformWEIBOStr, proto)
	})
	return
}
