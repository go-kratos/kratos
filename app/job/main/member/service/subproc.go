package service

import (
	"context"
	"encoding/json"
	"runtime/debug"
	"strings"
	"time"

	"go-common/app/job/main/member/model"
	share "go-common/app/service/main/share/model"
	"go-common/library/log"
	"go-common/library/net/ip"

	"github.com/pkg/errors"
)

const (
	_DedeMember          = "dede_member"
	_AsoAccount          = "aso_account"
	_DedeMemberPerson    = "dede_member_person"
	_DedeMemberSpace     = "dede_member_space"
	_MemberBaseInfo      = "user_base_"
	_MemberMoral         = "user_moral"
	_memberExp           = "user_exp_"
	_DedeMemberTags      = "dede_member_tags"
	_dedeMemberMoral     = "dede_member_moral"
	_memberRealnameApply = "realname_apply"
	_memberRealnameInfo  = "realname_info"
	_retry               = 3
)

var (
	_shareVideoType = map[int]struct{}{
		1: {},
		3: {},
	}
)

// subproc databus sub
func (s *Service) subproc() {
	var err error
	var c = context.TODO()
	for res := range s.ds.Messages() {
		mu := &model.Message{}
		if err = json.Unmarshal(res.Value, mu); err != nil {
			log.Error("member-job,json.Unmarshal (%v) error(%v)", string(res.Value), err)
			continue
		}
		for i := 0; i < _retry; i++ {
			if strings.HasPrefix(mu.Table, _MemberBaseInfo) {
				var (
					rank = &struct {
						Mid  int64 `json:"mid"`
						Rank int64 `json:"rank"`
					}{}
				)
				oldRank := int64(5000)
				if mu.Old != nil && len(mu.Old) != 0 {
					if err = json.Unmarshal(mu.Old, rank); err != nil {
						log.Error("json.Unmarsha(%s) error(%v)", string(mu.Old), err)
						break
					}
					oldRank = rank.Rank
				}
				if err = json.Unmarshal(mu.New, rank); err != nil {
					log.Error("json.Unmarsha(%s) error(%v)", string(mu.New), err)
					break
				}
				newRank := rank.Rank
				if oldRank <= 5000 && newRank >= 10000 {
					if err = s.initExp(c, rank.Mid); err != nil {
						log.Error("s.initExp(%d) error(%v)", rank.Mid, err)
					}
				}
				if err = s.dao.DelBaseInfoCache(c, rank.Mid); err != nil {
					continue
				}
				// sync face to old account .
				s.updateAccFace(c, rank.Mid) //todo delete
				// TODO: with update face or name to purge cache at the same time
				if err = s.dao.NotifyPurgeCache(c, rank.Mid, model.ActUpdateFace); err != nil {
					log.Error("s.dao.NotifyPurgeCache(%d, %s) error(%v)", rank.Mid, model.ActUpdateFace, err)
					break
				}
				item := &Item{
					Mid:    rank.Mid,
					Time:   time.Now(),
					Action: model.ActUpdateUname,
				}
				if err = s.cachepq.Put(item); err != nil {
					log.Error("Failed to put into cachepq with item: %+v: %+v", item, err)
					err = nil
				}
				log.Info("Notify to purge cache with mid(%d) action(%s) message(old: %s, new: %s)", rank.Mid, model.ActUpdateFace, string(mu.Old), string(mu.New))
			} else if mu.Table == _MemberMoral {
				var p = &model.MemberMid{}
				if err = json.Unmarshal(mu.New, p); err != nil {
					log.Error("member-job,json.Unmarshal (%v) error(%v)", string(mu.New), err)
					break
				}
				if err = s.dao.DelMoralCache(c, p.Mid); err != nil {
					continue
				}
				if err = s.dao.NotifyPurgeCache(c, p.Mid, model.ActUpdateMoral); err != nil {
					log.Error("s.dao.NotifyPurgeCache(%d, %s) error(%v)", p.Mid, model.ActUpdateMoral, err)
					break
				}
			} else if mu.Table == _memberRealnameInfo || mu.Table == _memberRealnameApply {
				var p = &struct {
					Mid    int64                     `json:"mid"`
					Status model.RealnameApplyStatus `json:"status"`
				}{}
				if err = json.Unmarshal(mu.New, p); err != nil {
					log.Error("member-job,json.Unmarshal (%v) error(%+v)", string(mu.New), err)
					break
				}
				if err = s.dao.DeleteRealnameCache(c, p.Mid); err != nil {
					log.Error("Delete RealnameCache cache err : %+v", err)
					continue
				}
				if err = s.dao.NotifyPurgeCache(c, p.Mid, "updateRealname"); err != nil {
					log.Error("s.dao.NotifyPurgeCache(%d, %s) error(%+v)", p.Mid, "updateRealname", err)
					break
				}
				log.Info("Notify to purge realname cache with mid(%d) action(%s) message(old: %s, new: %s)", p.Mid, "updateRealname", string(mu.Old), string(mu.New))
				if p.Status.IsPass() {
					// 尝试补发一次经验
					s.addExp(context.TODO(), &model.AddExp{
						Mid:   p.Mid,
						IP:    ip.InternalIP(),
						Ts:    time.Now().Unix(),
						Event: "identify",
					})
				}
				if mu.Table == _memberRealnameInfo {
					s.syncParsedRealnameInfo(c, p.Mid)
				}
			} else if strings.HasPrefix(mu.Table, _memberExp) {
				var (
					p   = &model.MemberMid{}
					exp *model.NewExp
				)
				if err = json.Unmarshal(mu.New, p); err != nil {
					log.Error("s.subproc() table(%s) json.Unmarshal() error(%v)", mu.Table, err)
					break
				}
				if exp, err = s.dao.SelExp(c, p.Mid); err != nil {
					log.Error("s.dao.SelNewExp(%d) error(%v)", p.Mid, err)
					break
				}
				if err = s.dao.SetExpCache(c, exp.Mid, exp.Exp); err != nil {
					log.Error("Failed to set exp cache: %+v: %+v", exp, err)
					break
				}
				log.Info("s.dao.SetExpCache(%d) set exp cache complete. exp(%d)", exp.Mid, exp.Exp)

				expChange, levelChange := isExpAndLevelChange(mu)
				if expChange {
					log.Info("Notify to purge cache with mid(%d) action(%s) message(old: %s, new: %s)", p.Mid, model.ActUpdateExp, string(mu.Old), string(mu.New))
					if err = s.dao.NotifyPurgeCache(c, p.Mid, model.ActUpdateExp); err != nil {
						log.Error("s.dao.NotifyPurgeCache(%d, %s) error(%v)", p.Mid, model.ActUpdateExp, err)
						break
					}
				}
				if levelChange {
					log.Info("Notify to purge cache with mid(%d) action(%s) message(old: %s, new: %s)", p.Mid, model.ActUpdateLevel, string(mu.Old), string(mu.New))
					if err = s.dao.NotifyPurgeCache(c, p.Mid, model.ActUpdateLevel); err != nil {
						log.Error("s.dao.NotifyPurgeCache(%d, %s) error(%v)", p.Mid, model.ActUpdateLevel, err)
						break
					}
				}
			}
			if err == nil {
				break
			}
		}
		if err = res.Commit(); err != nil {
			log.Error("databus.Commit err(%v)", err)
		}
		log.Info("subproc key:%v,topic: %v, part:%v offset:%v,message %s,", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
	}
}

// subproc databus sub
func (s *Service) accSubproc() {
	var (
		err error
	)
	for res := range s.accDs.Messages() {
		mu := &model.Message{}
		ms := &struct {
			Mid int64 `json:"mid"`
		}{}
		if err = json.Unmarshal(res.Value, mu); err != nil {
			log.Error("member-job,json.Unmarshal (%v) error(%v)", string(res.Value), err)
			continue
		}
		if err = json.Unmarshal(mu.New, ms); err != nil {
			log.Error("json.Unmarsha(%s) error(%v)", string(mu.New), ms)
			continue
		}
		for num := 0; num < 3; num++ {
			switch {
			//sex,face,rank
			case strings.HasPrefix(mu.Table, _DedeMember) && len(mu.Table) < 14:
				//err = s.setFace(c, ms.Mid)
				//birthday,dating,place,marital
			case mu.Table == _DedeMemberPerson:
				//sign
			case mu.Table == _DedeMemberSpace:
				//err = s.setSign(c, ms.Mid)
				//tag
			case strings.HasPrefix(mu.Table, _DedeMemberTags):
				//moral
			case mu.Table == _dedeMemberMoral:
			}
			if err != nil {
				log.Error("accSubproc err(%v)", err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
		if err = res.Commit(); err != nil {
			log.Error("databus.Commit err(%v)", err)
		}
		log.Info("subproc key:%v,topic: %v, part:%v offset:%v,message %s,", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
	}
}

// subproc databus sub.
func (s *Service) passportSubproc() {
	var err error
	for res := range s.passortDs.Messages() {
		mu := &model.Message{}
		ms := &struct {
			Mid int64 `json:"mid"`
		}{}
		if err = json.Unmarshal(res.Value, mu); err != nil {
			log.Error("member-job,json.Unmarshal (%v) error(%v)", string(res.Value), err)
			continue
		}
		if err = json.Unmarshal(mu.New, ms); err != nil {
			log.Error("json.Unmarsha(%s) error(%v)", string(mu.New), ms)
			continue
		}

		for num := 0; num < 3; num++ {
			//name
			if mu.Table == _AsoAccount {
				if mu.Action != "delete" {
					err = s.setName(ms.Mid)
				}
			}
			if err != nil {
				log.Error("passportSubproc err(%v)", err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
		if err = res.Commit(); err != nil {
			log.Error("databus.Commit err(%v)", err)
		}
		log.Info("subproc key:%v,topic: %v, part:%v offset:%v,message %s,", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
	}
}

func (s *Service) logproc() {
	for {
		mu := <-s.logDatabus.Messages()
		l := &model.UserLog{}
		err := json.Unmarshal(mu.Value, l)
		if err != nil {
			log.Error("Failed to parse log databus message value: value(%s): err: %+v", string(mu.Value), err)
			continue
		}
		// send log to report
		s.dao.AddExpLog(context.TODO(), l)

		// send log to origin hbase
		// content := make(map[string][]byte, len(l.Content))
		// for k, v := range l.Content {
		// 	content[k] = []byte(v)
		// }
		// content["ip"] = []byte(l.IP)
		// for i := 0; i < 3; i++ {
		// 	err := s.dao.AddLog(context.TODO(), l.Mid, l.TS, content, model.TableExpLog)
		// 	if err == nil {
		// 		break
		// 	}
		// 	log.Error("addlog  mid %d err %v", l.Mid, err)
		// 	time.Sleep(time.Millisecond * 500)
		// }
		log.Info("consumer key:%s,message:%s", mu.Key, mu.Value)
		mu.Commit()
	}
}

func (s *Service) expproc() {
	for {
		mu := <-s.expDatabus.Messages()
		ex := new(model.AddExp)
		err := json.Unmarshal(mu.Value, ex)
		if err != nil {
			log.Error("s.expproc() json.Unmarshal error(%v)", err)
			mu.Commit()
			continue
		}
		try := 0
		success := false
		for {
			if err = s.addExp(context.TODO(), ex); err == nil {
				success = true
				break
			}
			try++
			if try > 3 {
				log.Error("Failed to add exp, try 3 times mid: %d error: %+v", ex.Mid, err)
				mu.Commit()
				break
			}
			time.Sleep(time.Millisecond * 500)
		}
		if !success {
			continue
		}

		// 如果是一个观看视频的消息就尝试补发一下登录奖励
		if ex.Event == "view" {
			s.addExp(context.TODO(), &model.AddExp{
				Mid:   ex.Mid,
				IP:    ex.IP,
				Ts:    ex.Ts,
				Event: "login",
			})
			s.recoverMoral(context.TODO(), ex.Mid)
		}

		log.Info("expproc consumer key:%s,value: %s", mu.Key, mu.Value)
		mu.Commit()
	}
}

func isVideoShare(shareType int) bool {
	_, ok := _shareVideoType[shareType]
	return ok
}

func (s *Service) shareMidproc() {
	for {
		mu := <-s.shareMidDatabus.Messages()
		sh := new(share.MIDShare)
		err := json.Unmarshal(mu.Value, sh)
		if err != nil {
			log.Error("s.shareMidproc() json.Unmarshal error(%v)", err)
			mu.Commit()
			continue
		}

		if !isVideoShare(sh.TP) {
			log.Warn("Not a video share, skip to add exp: %+v", sh)
			mu.Commit()
			continue
		}

		try := 0
		success := false
		ex := &model.AddExp{
			Event: "share",
			Mid:   sh.MID,
			IP:    ip.InternalIP(),
			Ts:    sh.Time,
		}
		for {
			if err = s.addExp(context.TODO(), ex); err == nil {
				success = true
				break
			}
			try++
			if try > 3 {
				log.Error("Failed to add share exp, try 3 times mid: %d error: %+v", ex.Mid, err)
				mu.Commit()
				break
			}
			time.Sleep(time.Millisecond * 500)
		}
		if !success {
			continue
		}
		log.Info("shareMidproc consumer key:%s,value: %s", mu.Key, mu.Value)
		mu.Commit()
	}
}

func (s *Service) realnameSubproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("realnameSubproc panic(%+v) :\n %s", x, debug.Stack())
			go s.realnameSubproc()
		}
	}()
	log.Info("realnameSubproc run")
	var (
		c   = context.TODO()
		err error
	)
	for res := range s.realnameDatabus.Messages() {
		msg := &model.Message{}
		if err = json.Unmarshal(res.Value, msg); err != nil {
			log.Error("member-job,json.Unmarshal (%v) error(%v)", string(res.Value), errors.WithStack(err))
			continue
		}
		switch msg.Table {
		case "dede_identification_card_apply":
			if msg.Action == "delete" {
				log.Error("dede_identification_card_apply got delete msg (%s)", msg.New)
				continue
			}
			ms := &model.RealnameApplyMessage{}
			if err = json.Unmarshal(msg.New, ms); err != nil {
				err = errors.Wrapf(err, "dede_identification_card_apply , %s", msg.New)
				log.Error("%+v", err)
				continue
			}
			log.Info("upsert realname apply : (%+v)", ms)
			if err = s.dao.UpdateRealnameFromMSG(c, ms); err != nil {
				log.Error("%+v", err)
				continue
			}
			if err = s.dao.DeleteRealnameCache(c, ms.MID); err != nil {
				log.Error("Delete RealnameApplyStatus cache err : %+v", err)
				continue
			}
			if err = s.dao.NotifyPurgeCache(c, ms.MID, "updateRealname"); err != nil {
				log.Error("s.dao.NotifyPurgeCache(%d, %s) error(%+v)", ms.MID, "updateRealname", err)
				continue
			}
			log.Info("Notify to purge realname cache with mid(%d) action(%s) message(old: %s, new: %s)", ms.MID, "updateRealname", string(msg.Old), string(msg.New))
		case "dede_identification_card_apply_img":
			if msg.Action == "delete" {
				log.Error("dede_identification_card_apply_img got delete msg (%s)", msg.New)
				continue
			}
			ms := &model.RealnameApplyImgMessage{}
			if err = json.Unmarshal(msg.New, ms); err != nil {
				err = errors.Wrapf(err, "dede_identification_card_apply_img , %s", msg.New)
				log.Error("%+v", err)
				continue
			}
			log.Info("upsert realname apply img : (%+v)", ms)
			if err = s.dao.UpsertRealnameApplyImg(c, ms); err != nil {
				log.Error("%+v", err)
				continue
			}
		}
		if err = res.Commit(); err != nil {
			err = errors.Wrapf(err, "realnameSubproc commit")
			log.Error("%+v", err)
		}
		log.Info("Realname subproc key:%v,topic: %v, part:%v offset:%v,message %s,", res.Key, res.Topic, res.Partition, res.Offset, res.Value)
	}
}

func (s *Service) syncParsedRealnameInfo(ctx context.Context, mid int64) {
	info, err := s.dao.RealnameInfo(ctx, mid)
	if err != nil {
		log.Error("Failed to fetch realname info with mid: %d: %+v", mid, err)
		return
	}
	if info.Country != model.RealnameCountryChina ||
		info.CardType != model.RealnameCardTypeIdentity {
		log.Info("Skip to sync parsed realname info with mid: %d", mid)
		return
	}

	card, err := info.DecryptedCard()
	if err != nil {
		log.Error("Failed to decrypt realname card with mid: %d: %+v", mid, err)
		return
	}
	birth, gender, err := ParseIdentity(card)
	if err != nil {
		log.Error("Failed to parse idenitfy with mid: %d: %+v", mid, err)
		return
	}
	// sync to hive
	log.Infov(ctx,
		log.KV("action", "Syning realname parsed info"),
		log.KV("mid", info.MID),
		log.KV("birthday", birth.Format("2006-01-02")),
		log.KV("status", info.Status),
		log.KV("gender", gender),
		log.KV("mtime", info.MTime.Format("2006-01-02 15:04:05")),
	)
	s.ParsedRealnameInfoc.Infov(ctx,
		birth.Format("2006-01-02"),
		info.MTime.Format("2006-01-02 15:04:05"),
		int64(info.Status),
		int64(info.MID),
		gender,
	)
}
