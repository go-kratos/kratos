package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/interface/main/push-archive/model"
	"go-common/library/conf/env"
	"go-common/library/log"
)

const (
	_relationMidTable     = "user_relation_mid_"
	_relationTagUserTable = "user_relation_tag_user_"

	_retry = 3

	_relationStatusDeleted = 1 // 取消关注

	_relationTagSpecial = int64(-10) // 特殊关注的tag
)

func (s *Service) consumeRelationproc() {
	defer s.wg.Done()
	var err error
	for {
		msg, ok := <-s.relationSub.Messages()
		if !ok {
			log.Warn("s.RelationSub has been closed.")
			return
		}
		msg.Commit()
		time.Sleep(time.Millisecond)
		s.relMo++
		log.Info("consume relation key(%s) offset(%d) message(%s)", msg.Key, msg.Offset, msg.Value)
		m := new(model.Message)
		if err = json.Unmarshal(msg.Value, &m); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", msg.Value, err)
			continue
		}

		switch {
		case strings.HasPrefix(m.Table, _relationMidTable):
			err = s.relationMid(m.Action, m.New, m.Old)
		case strings.HasPrefix(m.Table, _relationTagUserTable):
			err = s.relationTagUser(m.Action, m.New, m.Old)
		default:
			continue
		}
		if err != nil {
			log.Error("consumeRelationproc data(%s) error(%+v)", msg.Value, err)
			if env.DeployEnv == env.DeployEnvProd {
				s.dao.WechatMessage(fmt.Sprintf("push-archive sync relation fail error(%v)", err))
			}
		}
	}
}

func (s *Service) addFans(upper, fans int64, tp int) (err error) {
	for i := 0; i < _retry; i++ {
		if err = s.dao.AddFans(context.TODO(), upper, fans, tp); err == nil {
			break
		}
		log.Info("retry s.dao.AddFans(%d,%d,%d)", upper, fans, tp)
	}
	if err != nil {
		log.Error("s.dao.AddFans(%d,%d,%d) error(%v)", upper, fans, tp, err)
	}
	return
}

func (s *Service) delFans(upper, fans int64) (err error) {
	for i := 0; i < _retry; i++ {
		if err = s.dao.DelFans(context.TODO(), upper, fans); err == nil {
			break
		}
		log.Info("retry s.dao.DelFans(%d,%d)", upper, fans)
	}
	if err != nil {
		log.Error("s.dao.DelFans(%d,%d) error(%v)", upper, fans, err)
	}
	return
}

func (s *Service) delSpecialAttention(upper, fans int64) (err error) {
	for i := 0; i < _retry; i++ {
		if err = s.dao.DelSpecialAttention(context.TODO(), upper, fans); err == nil {
			break
		}
		log.Info("retry s.dao.DelSpecialAttention(%d,%d)", upper, fans)
	}
	if err != nil {
		log.Error("s.dao.DelSpecialAttention(%d,%d) error(%v)", upper, fans, err)
	}
	return
}

// relationMid .
func (s *Service) relationMid(action string, nwMsg, oldMsg []byte) (err error) {
	n := &model.Relation{}
	if err = json.Unmarshal(nwMsg, n); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nwMsg, err)
		return
	}
	switch action {
	case _insertAct:
		if !n.Following() {
			return
		}
		err = s.addFans(n.Fid, n.Mid, model.RelationAttention)
	case _updateAct:
		o := &model.Relation{}
		if err = json.Unmarshal(oldMsg, o); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", oldMsg, err)
			return
		}
		if n.Status == o.Status && n.Attribute == o.Attribute {
			return
		}
		if n.Status == _relationStatusDeleted || !n.Following() {
			err = s.delFans(n.Fid, n.Mid) // 删除粉丝关系
		} else {
			err = s.addFans(n.Fid, n.Mid, model.RelationAttention) // 增加、更新关注数据
		}
	}
	if err != nil {
		log.Error("s.relationMid(%s,%s) error(%v)", nwMsg, oldMsg, err)
	}
	return
}

// relationTagUser .
func (s *Service) relationTagUser(action string, nwMsg, oldMsg []byte) (err error) {
	n := &model.RelationTagUser{}
	if err = json.Unmarshal(nwMsg, n); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", nwMsg, err)
		return
	}
	tagB, _ := base64.StdEncoding.DecodeString(n.Tag)
	n.Tag = string(tagB)
	switch action {
	case _insertAct:
		if !n.HasTag(_relationTagSpecial) {
			return
		}
		err = s.addFans(n.Fid, n.Mid, model.RelationSpecial)
	case _updateAct:
		o := &model.RelationTagUser{}
		if err = json.Unmarshal(oldMsg, o); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", oldMsg, err)
			return
		}
		tagB, _ = base64.StdEncoding.DecodeString(o.Tag)
		o.Tag = string(tagB)
		nt := n.HasTag(_relationTagSpecial)
		ot := o.HasTag(_relationTagSpecial)
		if nt && !ot {
			err = s.addFans(n.Fid, n.Mid, model.RelationSpecial)
		} else if !nt && ot {
			err = s.delSpecialAttention(n.Fid, n.Mid)
		}
	case _deleteAct:
		err = s.delSpecialAttention(n.Fid, n.Mid)
	}
	if err != nil {
		log.Error("s.relationTagUser(%s,%s) error(%v)", action, nwMsg, err)
	}
	return
}
