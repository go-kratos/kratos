package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/cache/redis"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_keyBankQuestions = "AntiFraud:BANK_%d:QUESTIONS"

	_keyQuestionFetchTime = "AntiFraud:USER_%s:ITEM_%d_%d_%s_PLATFORM:%d_Date"
	// 上次调起组件的id
	_keyComponentID = "AntiFraud:USER_%s:ITEM_%d_%d_%s_PLATFORM:%d_COMID"
	// 回答过的问题 id
	_keyAnsweredIds = "AntiFraud:USER_%s:ITEM_%d_%d_%s_PLATFORM:%d_ANSWER_IDS"
	// 组件内获取题目次数
	_keyBindBank = "AntiFraud:BIND_BANK_ITEM_%s"

	_keyComponentTimes = "AntiFraud:USER_%s:ITEM_%d_%d_%s_PLATFORM:%d_ANSWERTIMES"

	_keyBankID = "AntiFraud:BANKID_%d"
	//答案ids缓存
	_keyAnswerIds = "AntiFraud:AnswerIds_%d"

	// 图片id
	_keyPicID = "AntiFraud:AnswerAllPic"
	// 缓存所有图片id
	_keyPicIds = "AntiFraud:AnswerPic_%d"
	//默认100条数据
	_limit = 100
	// 缓存问题
	_keyQusInfo = "AntiFraud:QusInfo_%d"
	// 缓存图片
	_keyAnswerPicID = "AntiFraud:USER_%s:ITEM_%d_%d_%s_AnswerPic:%d_COMID"
	//答题日志list
	_keyAddLog = "AntiFraud:AddLog"
	//答案缓存
	_keyAnswer = "AntiFraud:Answer_%d"
	// 5分钟缓存
	_fiveMinuts = 5 * time.Minute
)

// PingRedis check redis connection
func (d *Dao) PingRedis(c context.Context) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()
	_, err = conn.Do("PING")
	return
}

// GetUserQuestionCache get
func (d *Dao) GetUserQuestionCache(c context.Context, args *model.ArgGetQuestion, bankID int64) {

	key := fmt.Sprintf("USER_%s_ITEM_%d_%d_%s_PLATFORM_%d_BANK_%d", args.UID, args.Source, args.TargetItemType, args.TargetItem, args.Platform, bankID)
	d.RedisDo(c, "GET", key)

}

// RedisDo redis cmd
func (d *Dao) RedisDo(c context.Context, cmd string, args ...interface{}) (reply interface{}, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	return conn.Do(cmd, args...)
}

// SetBankQuestionsCache 将题库下的全部问题存入缓存
func (d *Dao) SetBankQuestionsCache(c context.Context, bankID int64, ids []int64) (err error) {
	key := fmt.Sprintf(_keyBankQuestions, bankID)
	err = d.SetObj(c, key, ids, time.Hour)
	if err != nil {
		log.Error("d.SetBankQuestionsCache() error(%v)", err)
		return
	}

	return
}

// GetBankQuestionsCache 从缓存获取题库下的全部问题
func (d *Dao) GetBankQuestionsCache(c context.Context, bankID int64) (ids []int64) {
	key := fmt.Sprintf(_keyBankQuestions, bankID)
	reply, err := redis.Bytes(d.RedisDo(c, "GET", key))
	if err == redis.ErrNil {
		return
	}
	if err != nil {
		log.Error("查询 redis 出错 d.GetBankQuestionsCache(%d) error(%v)", bankID, err)
		return
	}
	err = json.Unmarshal(reply, &ids)
	if err != nil {
		return
	}

	return
}

// QusFetchTime 上次题目拉取时间
func (d *Dao) QusFetchTime(c context.Context, args *model.ArgGetQuestion) (ts int64) {
	ts, _ = redis.Int64(d.RedisDo(c, "GET", d.GetQusKey(_keyQuestionFetchTime, args)))
	return
}

// SetQusFetchTime 设置上次题目拉取时间
func (d *Dao) SetQusFetchTime(c context.Context, args *model.ArgGetQuestion, ts int64) (err error) {
	err = d.Setex(c, d.GetQusKey(_keyQuestionFetchTime, args), ts, _fiveMinuts)
	if err != nil {
		log.Error("d.SetQusFetchTime(%v, %d) error(%v)", args, ts, err)
	}
	return
}

// GetQusKey getkey
func (d *Dao) GetQusKey(format string, args *model.ArgGetQuestion) (s string) {
	s = fmt.Sprintf(format, args.UID, args.Source, args.TargetItemType, args.TargetItem, args.Platform)
	return
}

// Setex set
func (d *Dao) Setex(c context.Context, key string, data interface{}, exp time.Duration) error {

	log.Info(" d.setex(%s, %v, %d)", key, data, exp)
	_, err := d.RedisDo(c, "SETEX", key, int(exp/1e9), data)
	if err != nil {
		log.Error("d.Setex(%s, %v) error(%v)", key, data, err)
	}
	return err
}

// SetObj set
func (d *Dao) SetObj(c context.Context, key string, obj interface{}, exp time.Duration) error {
	log.Info(" d.setex(%s, %v, %d)", key, obj, exp)

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	_, err = d.RedisDo(c, "SETEX", key, int(exp/1e9), data)
	if err != nil {
		log.Error("d.Setex(%s, %v) error(%v)", key, string(data), err)
	}

	return err
}

//GetObj get
func (d *Dao) GetObj(c context.Context, key string, obj interface{}) (err error) {
	reply, err := redis.Bytes(d.RedisDo(c, "GET", key))

	if err != nil {
		return
	}

	err = json.Unmarshal(reply, obj)
	if err != nil {
		return
	}
	return
}

// GetAnsweredID 获取已回答问题
func (d *Dao) GetAnsweredID(c context.Context, args *model.ArgGetQuestion) (ids []int64) {
	key := d.GetQusKey(_keyAnsweredIds, args)
	ids, err := redis.Int64s(d.RedisDo(c, "SMEMBERS", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("d.GetAnsweredID(%v) error(%v)", args, err)
		}
	}

	return
}

// SetAnsweredID 设置已回答问题
func (d *Dao) SetAnsweredID(c context.Context, args *model.ArgGetQuestion, questionID int64) (err error) {
	key := d.GetQusKey(_keyAnsweredIds, args)
	_, err = d.RedisDo(c, "SADD", key, questionID)
	if err != nil {
		log.Error("d.SetAnsweredID(%v, %d) error(%v)", args, questionID, err)
	}
	_, err = d.RedisDo(c, "EXPIRE", key, int(_fiveMinuts/1e9))

	if err != nil {
		log.Error("d.SetAnsweredID expire (%v, %d) error(%v)", args, questionID, err)
	}

	return
}

// RmAnsweredID 删除已回答问题
func (d *Dao) RmAnsweredID(c context.Context, args *model.ArgGetQuestion) (err error) {
	key := d.GetQusKey(_keyAnsweredIds, args)
	_, err = d.RedisDo(c, "DEL", key)
	if err != nil {
		log.Error("d.SetAnsweredID(%v, %d) error(%v)", args, err)
	}
	return
}

// GetComponentID 获取组件id
func (d *Dao) GetComponentID(c context.Context, args *model.ArgGetQuestion) (cID int, err error) {
	key := d.GetQusKey(_keyComponentID, args)
	if cID, err = redis.Int(d.RedisDo(c, "GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("d.GetComponentID(%v) error(%v)", args, err)
		}
	}
	return
}

// SetComponentID 设置组件id
func (d *Dao) SetComponentID(c context.Context, args *model.ArgGetQuestion) (err error) {
	key := d.GetQusKey(_keyComponentID, args)
	err = d.Setex(c, key, args.ComponentID, _fiveMinuts)
	if err != nil {
		log.Error("d.SetComponentID(%v) error(%v)", args, err)
	}
	return
}

// GetComponentTimes 获取组件答题次数
func (d *Dao) GetComponentTimes(c context.Context, args *model.ArgGetQuestion) (cID int64, err error) {
	key := d.GetQusKey(_keyComponentTimes, args)
	if cID, err = redis.Int64(d.RedisDo(c, "GET", key)); err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("d.GetComponentTimes(%v) error(%v)", args, err)
		}
	}
	return
}

// SetComponentTimes 设置组件答题次数
func (d *Dao) SetComponentTimes(c context.Context, args *model.ArgGetQuestion) (err error) {
	key := d.GetQusKey(_keyComponentTimes, args)
	err = d.Setex(c, key, 0, _fiveMinuts)
	if err != nil {
		log.Error("d.SetComponentID(%v) error(%v)", args, err)
	}
	return
}

// IncrComponentTimes 组件计数
func (d *Dao) IncrComponentTimes(c context.Context, args *model.ArgGetQuestion) (err error) {
	key := d.GetQusKey(_keyComponentTimes, args)
	_, err = d.RedisDo(c, "INCR", key)
	if err != nil {
		log.Error("d.GetComponentTimes(%v) error(%v)", args, err)
	}
	return
}

// GetQusBankInfoCache get
func (d *Dao) GetQusBankInfoCache(c context.Context, qbid int64) (oi *model.QuestionBank, err error) {
	oi = &model.QuestionBank{}
	key := fmt.Sprintf(_keyBankID, qbid)
	err = d.GetObj(c, key, oi)

	if err == redis.ErrNil {
		err = nil
		oi, err = d.GetQusBankInfo(c, qbid)
		if err != nil {
			log.Error("d.GetQusBankInfoCache error(%v)", err)
			return
		}
		err = d.SetObj(c, key, oi, _fiveMinuts)
	}

	return
}

// GetBindBankInfo get
func (d *Dao) GetBindBankInfo(c context.Context, source, targetItemType int8, targetItem string) (bind *model.QuestionBankBind, err error) {

	bind = &model.QuestionBankBind{}
	key := fmt.Sprintf(_keyBindBank, targetItem)
	err = d.GetObj(c, key, bind)
	if err == redis.ErrNil {
		//err = nil
		binds, err1 := d.GetBindBank(c, source, targetItemType, []string{targetItem})
		if err1 != nil {
			log.Error("s.GetQuestion(%v) error(%v)", targetItem, err)
			err = err1
			return
		}
		if len(binds) < 1 {
			log.Warn("s.GetQuestion(%v) 未找到题库绑定关系", targetItem)
			err = ecode.BindBankNotFound
			return
		}

		bind = binds[0]
		if bind.QuestionBank == nil {
			log.Error("s.GetQuestion(%v) 未找到已绑定的题库", targetItem)
			err = ecode.QusbNotFound
			return
		}

		err = d.SetObj(c, key, bind, time.Hour)
		//return
	}

	return
}

// CorrectAnswerIds id
func (d *Dao) CorrectAnswerIds(c context.Context, qid int64) (ids []int64, err error) {
	key := fmt.Sprintf(_keyAnswerIds, qid)
	err = d.GetObj(c, key, &ids)
	if err == redis.ErrNil {
		err = nil
		answers, err1 := d.GetAnswerList(c, qid)
		for _, answer := range answers {
			if answer.IsCorrect == 1 {
				ids = append(ids, answer.AnswerID)
			}
		}
		if err1 != nil {
			err = err1
			log.Error("d.GetQusBankInfoCache error(%v)", err)
			return
		}
		err1 = d.SetObj(c, key, ids, time.Hour)
		err = err1
		return
	}

	return
}

// GetRandPic 背景图
func (d *Dao) GetRandPic(c context.Context, args *model.ArgGetQuestion) (oi *model.QuestBkPic, err error) {
	key := _keyPicID
	cnt, err := redis.Int(d.RedisDo(c, "LLen", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("GetStatisticsCache do(RPOP, %s) error(%v)", key, err)
		}
		return
	}

	if cnt == 0 {
		err1 := d.PushAllPic(c)
		if err1 != nil {
			if err1 == redis.ErrNil {
				err = nil
			} else {
				log.Error("PushAllPic do(RPOP, %s) error(%v)", key, err)
			}
			return
		}

	}

	id, err := redis.Int(d.RedisDo(c, "RPOP", key))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
		} else {
			log.Error("GetStatisticsCache do(RPOP, %s) error(%v)", key, err)
		}
		return
	}

	oi, err = d.GetPic(c, id)
	if err != nil {
		return
	}

	//缓存坐标
	picKey := d.GetQusKey(_keyAnswerPicID, args)
	err = d.SetObj(c, picKey, oi, time.Hour)
	if err != nil {
		return
	}
	return
}

// PushAllPic 事先放入redis list中
func (d *Dao) PushAllPic(c context.Context) (err error) {
	cnt, err := d.GetPicCount(c)
	limit := _limit
	page := (int)(cnt / limit)
	for i := 0; i <= page; i++ {
		key := _keyPicID
		ids, _ := d.GetAllPicIds(c, i*limit, limit)
		a := make([]interface{}, 0)
		a = append(a, key)
		for _, id := range ids {
			a = append(a, id)

		}
		if len(a) > 1 {
			_, err := d.RedisDo(c, "LPUSH", a...)
			if err != nil {
				log.Error("[PushAllPic]conn.Do(lpush, %s) error(%v)", key, err)
			}
		}

	}
	return
}

// GetPic 获取背景图
func (d *Dao) GetPic(c context.Context, id int) (oi *model.QuestBkPic, err error) {

	key := fmt.Sprintf(_keyPicIds, id)
	oi = &model.QuestBkPic{}
	err = d.GetObj(c, key, oi)
	if err == redis.ErrNil {
		err = nil
		picInfo, err1 := d.GetRandomPic(c, id)

		if err1 != nil {
			err = err1
			log.Error("d.GetPic error(%v)", err1)
			return
		}
		err1 = d.SetObj(c, key, picInfo, time.Hour)
		oi = picInfo
		err = err1
		return
	}

	return
}

// PushAnswer push
func (d *Dao) PushAnswer(c context.Context, answer *model.ArgCheckAnswer, isCorrect int8) (affect int64, err error) {
	ids := xstr.JoinInts(answer.Answers)
	obj := model.AddLog{
		UID:       answer.UID,
		QsID:      answer.QsID,
		Platform:  answer.Platform,
		Source:    answer.Source,
		Ids:       ids,
		IsCorrect: isCorrect,
	}
	data, err := json.Marshal(obj)
	if err != nil {
		return
	}

	_, err = d.RedisDo(c, "LPUSH", _keyAddLog, data)
	if err != nil {
		log.Error("[PushAllPic]conn.Do(RPOP,%s) error(%v)", "_keyAddLog", err)

	}
	affect = obj.QsID
	return
}

// PopAnswer pop
func (d *Dao) PopAnswer(c context.Context) {

	for {
		reply, err := redis.Bytes(d.RedisDo(c, "RPOP", _keyAddLog))
		if len(reply) > 0 && err == nil {
			answer := &model.AddLog{}
			err = json.Unmarshal(reply, answer)
			if err == nil {
				_, err = d.db.Exec(c, _addUserAnswerSQL, answer.UID, answer.QsID, answer.Platform, answer.Source, answer.Ids, answer.IsCorrect)
				if err != nil {
					log.Error("d.PopAnswer error(%v)", err)
					time.Sleep(time.Second * 1)
					continue
				}
			}

		}
		time.Sleep(time.Second * 5)
	}

}

// GetCacheQus get
func (d *Dao) GetCacheQus(c context.Context, id int64) (oi *model.Question, err error) {
	key := fmt.Sprintf(_keyQusInfo, id)
	oi = &model.Question{}
	err = d.GetObj(c, key, oi)
	if err == redis.ErrNil {
		err = nil
		info, err1 := d.GetQusInfo(c, id)

		if err1 != nil {
			err = err1
			log.Error("d.GetPic error(%v)", err1)
			return
		}
		err1 = d.SetObj(c, key, info, time.Hour)
		oi = info
		err = err1
		return
	}
	return
}

// GetCacheAnswerPic get
func (d *Dao) GetCacheAnswerPic(c context.Context, args *model.ArgGetQuestion) (oi *model.QuestBkPic, err error) {
	picKey := d.GetQusKey(_keyAnswerPicID, args)
	oi = &model.QuestBkPic{}
	err = d.GetObj(c, picKey, oi)
	if err != nil {
		return
	}
	return
}

// DelTargetItemBindCache del cache
func (d *Dao) DelTargetItemBindCache(c context.Context, skuID string) (err error) {
	key := fmt.Sprintf(_keyBindBank, skuID)
	_, err = d.RedisDo(c, "DEL", key)
	if err != nil {
		if err != redis.ErrNil {
			log.Error("d.DelTargetItemBind(%v, %d) error(%v)", skuID, err)
			return
		}
		err = nil
	}
	return
}

// DelQusCache del cache
func (d *Dao) DelQusCache(c context.Context, id int64) (err error) {
	key := fmt.Sprintf(_keyQusInfo, id)
	_, err = d.RedisDo(c, "DEL", key)
	if err != nil {
		if err != redis.ErrNil {
			log.Error("d.DelQusCache(%v, %d) error(%v)", id, err)
			return
		}
		err = nil
	}
	return
}

// GetAnswersByCache cache
func (d *Dao) GetAnswersByCache(c context.Context, id int64) (oi []*model.Answer, err error) {
	key := fmt.Sprintf(_keyAnswer, id)
	err = d.GetObj(c, key, &oi)
	if err == redis.ErrNil {
		err = nil
		info, err1 := d.GetAnswerList(c, id)

		if err1 != nil {
			err = err1
			log.Error("d.GetPic error(%v)", err1)
			return
		}
		err1 = d.SetObj(c, key, info, time.Hour)
		oi = info
		err = err1
		return
	}
	return
}

// DelAnswerCache del cache
func (d *Dao) DelAnswerCache(c context.Context, id int64) (err error) {
	key := fmt.Sprintf(_keyAnswer, id)
	_, err = d.RedisDo(c, "DEL", key)
	if err != nil {
		if err != redis.ErrNil {
			log.Error("d.DelAnswerCache(%v, %d) error(%v)", id, err)
			return
		}
		err = nil
	}
	return
}

// DelQusBankCache del cache
func (d *Dao) DelQusBankCache(c context.Context, id int64) (err error) {
	key := fmt.Sprintf(_keyBankID, id)
	_, err = d.RedisDo(c, "DEL", key)
	if err != nil {
		if err != redis.ErrNil {
			log.Error("d.DelQusBankCache(%v, %d) error(%v)", id, err)
			return
		}
		err = nil
	}
	return
}
