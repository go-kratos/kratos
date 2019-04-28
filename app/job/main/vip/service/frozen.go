package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/vip/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
)

func (s *Service) accloginproc() {
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("Runtime error caught: %+v", r)
			go s.accloginproc()
		}
	}()
	var (
		err     error
		msgChan = s.accLogin.Messages()
		msg     *databus.Message
		ok      bool
	)
	for {
		msg, ok = <-msgChan
		log.Info("login ip msg %+v", string(msg.Value))
		if !ok {
			log.Info("accLogin msgChan closed")
		}
		if err = msg.Commit(); err != nil {
			log.Error("msg.Commit err(%+v)", err)
		}
		m := &model.LoginLog{}
		if err = json.Unmarshal([]byte(msg.Value), m); err != nil {
			log.Error("json.Unmarshal(%v) err(%+v)", m, err)
			continue
		}
		s.Frozen(context.TODO(), m)
	}
}

// Frozen handle vip frozen logic.
func (s *Service) Frozen(c context.Context, ll *model.LoginLog) (err error) {
	var (
		lc  int64
		uvs *model.VipUserInfo
		ctx = context.TODO()
	)
	// 判定用户是否为vip
	if uvs, err = s.dao.VipStatus(ctx, ll.Mid); err != nil {
		log.Error("s.dao.VipStatus（%d）err(%+v)", ll.Mid, err)
		return
	}
	if uvs == nil || uvs.Status == model.VipStatusOverTime {
		log.Warn("user(%d) not vip.(%+v)", ll.Mid, uvs)
		return
	}
	// 判定是否为15分钟4次以上不同ip登录
	if err = s.dao.AddLogginIP(ctx, ll.Mid, ll.IP); err != nil {
		log.Error("s.dao.AddLogginIP(%d）err(%+v)", ll.Mid, err)
		return
	}
	if lc, err = s.dao.LoginCount(ctx, ll.Mid); err != nil {
		log.Error("s.dao.LoginCount(%d）err(%+v)", ll.Mid, err)
		return
	}
	if lc >= s.c.Property.FrozenLimit {
		if err = s.dao.Enqueue(ctx, ll.Mid, time.Now().Add(s.frozenDate).Unix()); err != nil {
			log.Error("enqueue error(%+v)", err)
		}
		if err = s.dao.SetVipFrozen(ctx, ll.Mid); err != nil {
			log.Error("set vip frozen err(%+v)", err)
		}
		//通知业务方清理缓存
		s.cleanCache(ll.Mid)
		if err = s.notifyOldVip(ll.Mid, -1); err != nil {
			log.Error("del vip java frozen err(%+v)", err)
		}
		log.Info("mid(%+v) frozen success", ll.Mid)
	}
	return
}

// unFrozenJob timing to unFrozen vip user
func (s *Service) unFrozenJob() {
	log.Info("unfrozen job start........................................")
	defer func() {
		if r := recover(); r != nil {
			r = errors.WithStack(r.(error))
			log.Error("recover panic  error(%+v)", r)
		}
		log.Info("unfrozen job end.............................")
	}()

	var (
		err  error
		mids []int64
		ctx  = context.TODO()
	)
	if mids, err = s.dao.Dequeue(ctx); err != nil {
		log.Error("s.dao.Dequeue err(%+v)", err)
		return
	}
	for _, mid := range mids {
		if err = s.dao.RemQueue(ctx, mid); err != nil {
			log.Error("s.dao.RemQueue(%d）err(%+v)", mid, err)
			continue
		}
		if err = s.dao.DelCache(ctx, mid); err != nil {
			log.Error("del cache mid(%+v) err(%+v)", mid, err)
		}
		if err = s.dao.DelVipFrozen(ctx, mid); err != nil {
			log.Error("del vip frozen err(%+v)", err)
		}
		if err = s.notifyOldVip(mid, 0); err != nil {
			log.Error("del vip java frozen err(%+v)", err)
		}
		s.cleanCache(mid)
		log.Info("mid(%+v) unfrozen success", mid)
	}
}

// FIXME AFTER REMOVE JAVA.
func (s *Service) notifyOldVip(mid, status int64) error {
	return s.dao.OldFrozenChange(mid, status)
}
