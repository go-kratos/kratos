package service

import (
	"context"
	"math"
	"math/rand"
	"time"

	"go-common/app/service/openplatform/anti-fraud/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) beginTran(c context.Context) (*xsql.Tx, error) {
	return s.d.BeginTran(c)
}

// GetQusInfo 获取题目信息
func (s *Service) GetQusInfo(c context.Context, qsID int64) (res *model.Question, err error) {
	if res, err = s.d.GetQusInfo(c, qsID); err != nil {
		return
	}
	return
}

// GetAnswerList 答案列表
func (s *Service) GetAnswerList(c context.Context, qsID int64) (res []*model.Answer, err error) {
	if res, err = s.d.GetAnswerList(c, qsID); err != nil {
		return
	}
	return
}

// AddQus 添加题目
func (s *Service) AddQus(c context.Context, in *model.AddQus, answers []model.Answer) (res model.AddReturn, err error) {
	tx, err := s.beginTran(c)
	if err != nil {
		log.Error("AddQusQusbeginTran error(%v)", err)
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

	qb := &model.Question{
		QsID:       time.Now().UnixNano() / 1e6,
		QsName:     in.Name,
		QsDif:      in.Dif,
		QsBId:      in.BId,
		QsType:     in.Type,
		AnswerType: in.AnType,
	}
	res.ID, err = s.d.InsertQus(c, qb)
	if err != nil {
		log.Error("s.dao.AddQusBank(%v) error(%v)", qb)
		return
	}

	list := make([]*model.AnswerAdd, 0)
	for _, item := range answers {
		answer := &model.AnswerAdd{}
		answer.IsCorrect = item.IsCorrect
		answer.AnswerContent = item.AnswerContent
		answer.QsID = qb.QsID
		answer.AnswerID = time.Now().UnixNano() / 1e6
		time.Sleep(time.Millisecond)
		list = append(list, answer)
	}

	err = s.d.MultiAddAnwser(c, list)
	if err != nil {
		log.Error("s.dao.AnswerAdd error(%v)", err)
		return
	}

	_, err = s.d.UpdateQsBankCnt(c, in.BId)
	if err != nil {
		log.Error("s.dao.UpdateQsBankCnt error(%v)", err)
		return
	}

	return
}

// DelQus 删除题目
func (s *Service) DelQus(c context.Context, qid int64) (res bool, err error) {
	res = false
	tx, err := s.beginTran(c)
	if err != nil {
		log.Error("DelQus beginTran error(%v)", err)
		return
	}
	var rows int64

	info, err := s.GetQusInfo(c, qid)
	if err != nil || info == nil {
		err = ecode.QusIDInvalid
		return
	}

	cnt, err := s.d.CountBindItem(c, info.QsBId)
	if cnt > 0 || err != nil {
		err = ecode.BankUsing
	}

	rows, err = s.d.DelQus(c, qid)
	if err != nil || rows == 0 {
		tx.Rollback()
		log.Error("dao.Subject.DelQus error(%v) or rows==0", err)
		err = ecode.ParamInvalid
		return
	}

	rows, err = s.d.DelAnwser(c, qid)
	if err != nil || rows == 0 {
		tx.Rollback()
		log.Error("dao.Subject.Del Answer error(%v) or rows==0", err)
		err = ecode.ParamInvalid
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}

	_, err = s.d.UpdateQsBankCnt(c, info.QsBId)
	if err != nil {
		log.Error("s.dao.UpdateQsBankCnt error(%v)", err)
		return
	}

	res = true
	return
}

// GetQuslist 题目列表
func (s *Service) GetQuslist(c context.Context, pageNo int, pageSize int, qBid int64) (res []*model.QuestionAll, err error) {

	offset, limit := s.PageInfo(pageNo, pageSize)
	list, err := s.d.GetQusList(c, offset, limit, qBid)
	if err != nil {
		return
	}
	for _, item := range list {
		alist, err := s.d.GetAnswerList(c, item.QsID)
		if err != nil {
			break
		}
		data := &model.QuestionAll{
			Question:    *item,
			AnswersList: alist,
		}
		res = append(res, data)
	}

	return
}

// GetQusTotal 题目列表
func (s *Service) GetQusTotal(c context.Context, bid int64) (res int64, err error) {
	if res, err = s.d.GetQusCount(c, bid); err != nil {
		return
	}
	return
}

// UpdateQus 更新题库
func (s *Service) UpdateQus(c context.Context, in *model.ArgUpdateQus, answers []model.Answer) (res bool, err error) {
	res = false
	info, _ := s.GetQusBankInfo(c, in.BId)
	if info == nil {
		err = ecode.QusbNotFound
		return
	}
	var row int64
	tx, err := s.beginTran(c)
	if err != nil {
		log.Error("Del Qus beginTran error(%v)", err)
		return
	}

	if row, err = s.d.UpdateQus(c, in, answers); err != nil {
		tx.Rollback()
		log.Error("s.dao.DelQusBank(%v) error(%v)", err)
		err = ecode.UpdateError
		return
	}
	if row > 0 {
		res = true
	}

	_, err = s.d.DelAnwser(c, in.QsID)
	if err != nil {
		tx.Rollback()
		log.Error("s.dao.DelAnwser error(%v)", err)
		err = ecode.UpdateError
		return
	}

	list := make([]*model.AnswerAdd, 0)
	for _, item := range answers {
		answer := model.AnswerAdd{}
		answer.Answer = item
		answer.QsID = in.QsID
		if item.AnswerID != 0 {
			_, err = s.d.UpdateAnwser(c, &answer)
			if err != nil {
				tx.Rollback()
				log.Error("s.dao.UpdateAnwser error(%v)", err)
				return
			}
		} else {
			list = append(list, &answer)
		}

	}
	//批量插入
	if len(list) > 0 {
		err = s.d.MultiAddAnwser(c, list)
		if err != nil {
			tx.Rollback()
			log.Error("s.dao.MultiAddAnwser error(%v)", err)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		return
	}
	res = true

	_, err = s.d.UpdateQsBankCnt(c, in.BId)
	if err != nil {
		log.Error("s.dao.UpdateQsBankCnt error(%v)", err)
		return
	}

	err = s.d.DelQusCache(c, in.QsID)
	if err != nil {
		log.Error("s.dao.DelQusCache error(%v)", err)
		return
	}

	err = s.d.DelAnswerCache(c, in.QsID)
	if err != nil {
		log.Error("s.dao.DelAnswerCache error(%v)", err)
		return
	}
	return
}

// CheckAnswer 检查答案
func (s *Service) CheckAnswer(c context.Context, qusID int64, qusType int8, anlist []model.Answer) (etag string, err error) {
	var correct int8
	for _, item := range anlist {
		if item.AnswerContent == "" {
			etag = "answer_content is empty"
			err = ecode.ParamInvalid
			return
		}
		if item.IsCorrect > 1 || item.IsCorrect < 0 {
			etag = "invalid is_correct"
			err = ecode.AnswerError
			return
		}
		correct = item.IsCorrect + correct
		if item.QsID > 0 {
			if qusID == item.QsID {
				if res, _ := s.d.GetQusInfo(c, item.QsID); res == nil {
					etag = "invalid QsID"
					err = ecode.QusIDInvalid
					return
				}
				continue
			} else {
				etag = "invalid QsID"
				err = ecode.QusIDInvalid
				return
			}

		}
	}

	if qusType == model.MULTIPLECHOICE {
		if correct < 2 {
			etag = "invalid is_correct"
			err = ecode.AnswerError
			return
		}
	} else {
		if correct != 1 {
			etag = "invalid is_correct"
			err = ecode.AnswerError
			return
		}
	}

	return
}

// GetQuestion 获取题目
func (s *Service) GetQuestion(c context.Context, args *model.ArgGetQuestion) (q *model.GetQuestionItem, err error) {
	// 获取已绑定的题库

	bind, err := s.d.GetBindBankInfo(c, args.Source, args.TargetItemType, args.TargetItem)
	if err != nil {
		return
	}

	// 判断是不是答题组件
	cid, err := s.d.GetComponentID(c, args)
	if err != nil {
		log.Info("s.GetQuestion(%v) 两次获取题目间隔小于 cd 时间", args)
		err = ecode.GetComponentIDErr
		return
	}

	// 组件第一题答题
	if cid == 0 || cid != args.ComponentID {

		// 检查上次时间与这次的差值如果小于 cd 时间
		//if time.Now().Unix()-s.d.QusFetchTime(c, args) < bind.QuestionBank.CdTime {
		//	log.Info("s.GetQuestion(%v) 两次获取题目间隔小于 cd 时间", args)
		//	err = ecode.AnswerIntervalInvalid
		//	return
		//}

		//缓存新的组件
		if err = s.d.SetComponentID(c, args); err != nil {
			err = ecode.SetComponentIDErr
			return
		}

		//	把组件答题次数清零
		err = s.d.SetComponentTimes(c, args)
		if err != nil {
			log.Info("s.GetQuestion(%v) 两次获取题目间隔小于 cd 时间", args)
			err = ecode.SetComponentTimesErr
			return
		}

	}

	BankInfo, err := s.d.GetQusBankInfoCache(c, bind.QsBId)
	if err != nil {
		err = ecode.GetQusBankInfoCache
		return
	}

	//	同一组件
	times, err := s.d.GetComponentTimes(c, args)
	if err != nil {
		err = ecode.GetComponentTimesErr
		return
	}

	if times == BankInfo.MaxRetryTime {
		err = ecode.SameCompentErr
		return
	}

	// 获取题库包含的问题
	var qusIds []int64

	qusIds, err = s.d.GetQusIds(c, BankInfo.QsBId)
	if err != nil {
		err = ecode.GetQusIDsErr
		return
	}

	q = new(model.GetQuestionItem)
	// 随机取一道
	q.Question, err = s.randQuestion(c, qusIds, args)
	if err != nil {
		err = ecode.GetQusIDsErr
		return
	}
	//记录答题次数

	q.AllCnt = BankInfo.MaxRetryTime
	q.AnTime = times + 1

	//产生答题背景图片
	if q.QsType == 1 {
		q.QuestBkPic, err = s.d.GetRandPic(c, args)
	}

	if err != nil {
		log.Error("s.GetRandPic(%v) error(%v)", args, err)
		q = nil
		return
	}

	q.Answers, err = s.d.GetAnswersByCache(c, q.QsID)
	if err != nil {
		log.Error("s.GetQuestion(%v) error(%v)", args, err)
		q = nil
		return
	}

	// 组件题目+1
	err = s.d.IncrComponentTimes(c, args)
	if err != nil {
		q = nil
		return
	}

	// 记录上次获取题目时间
	err = s.d.SetQusFetchTime(c, args, time.Now().Unix())
	if err != nil {
		q = nil
		return
	}

	return
}

// UserAnswer 回答
func (s *Service) UserAnswer(c context.Context, check *model.ArgCheckAnswer) (res bool, err error) {

	// 检查题目合法不合法
	args := check.ArgGetQuestion
	bind, err := s.d.GetBindBankInfo(c, args.Source, args.TargetItemType, args.TargetItem)

	if err != nil {
		return
	}

	//答题记录
	insertID, err := s.d.AddUserAnwser(c, check, 1)
	if err != nil || insertID == 0 {
		log.Error("s.dao.AddUser Answser error(%v)", err)
		return
	}

	//检查这是第几次答题
	times, err := s.d.GetComponentTimes(c, &check.ArgGetQuestion)
	if err != nil {
		err = ecode.GetComponentTimesErr
		return
	}

	if times > bind.QuestionBank.MaxRetryTime {
		err = ecode.SameCompentErr
		return
	}

	qusInfo, err := s.d.GetCacheQus(c, check.QsID)

	if err != nil {
		err = ecode.AnswerError
		return
	}

	//检查xy坐标
	if qusInfo.QsType == 1 {
		answerPic, err1 := s.d.GetCacheAnswerPic(c, &check.ArgGetQuestion)
		if err1 != nil {
			err = ecode.AnswerError
			return
		}
		m := check.X - answerPic.X
		n := check.X - answerPic.X
		if int(math.Abs(float64(m))) > 75 || int(math.Abs(float64(n))) > 75 {
			err = ecode.AnswerPoiError
			return
		}
	}

	//检查答案
	ids, err := s.d.CorrectAnswerIds(c, check.QsID)
	if len(ids) < 1 || len(ids) != len(check.Answers) || err != nil {
		log.Info("s.CheckAnswer(%v) 答案错误", err)
		err = ecode.AnswerError
		return
	}
	num := 0
	for _, id := range ids {
		for _, answerid := range check.Answers {
			if answerid == id {
				num++
				continue
			}
		}
	}

	if num != len(ids) {
		log.Info("s.CheckAnswer(%v) 答案错误", err)
		err = ecode.AnswerError
		return
	}

	res = true
	return
}

// randQuestion 随机获取题目
func (s *Service) randQuestion(c context.Context, qs []int64, args *model.ArgGetQuestion) (q *model.Question, err error) {
	// 获取已经答过的题目 id
	usedIds := s.d.GetAnsweredID(c, args)
	var notUseQ []int64
	for _, id := range qs {
		used := false
		for _, i := range usedIds {
			if id == i {
				used = true

				break
			}
		}

		if !used {
			notUseQ = append(notUseQ, id)
		}
	}
	// 全部使用过了就重新开始随机
	if len(notUseQ) == 0 {
		err = s.d.RmAnsweredID(c, args)
		if err != nil {
			log.Error("d.randQuestion(%v, %v) error(%v)", qs, args, err)
			return
		}
		notUseQ = qs
	}

	qid := notUseQ[rand.Intn(len(notUseQ))]
	q, err = s.d.GetCacheQus(c, qid)
	if err != nil {
		log.Error("d.GetQusInfo(%v) error(%v)", qid, err)
		return
	}

	// 设置已经答过的题目 id
	err = s.d.SetAnsweredID(c, args, qid)
	if err != nil {
		log.Error("d.randQuestion(%v, %v) s.d.SetAnsweredID(%d) error(%v)", qs, args, qid, err)
		return
	}

	return
}
