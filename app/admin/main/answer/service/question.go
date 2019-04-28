package service

import (
	"context"
	"encoding/csv"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/answer/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// QuestionList .
func (s *Service) QuestionList(c context.Context, arg *model.ArgQue) (res *model.QuestionPage, err error) {
	res = &model.QuestionPage{}
	if res.Total, err = s.dao.QuestionCount(c, arg); err != nil {
		return
	}
	res.Items = []*model.QuestionDB{}
	if res.Total > 0 {
		if res.Items, err = s.dao.QuestionList(c, arg); err != nil {
			return
		}
	}
	return
}

// UpdateStatus update question state
func (s *Service) UpdateStatus(c context.Context, qid int64, state int8, operator string) (err error) {
	var (
		r int64
		q *model.Question
	)
	if q, err = s.dao.QueByID(c, qid); err != nil {
		log.Error("dao QueByID(%d) error(%v)", qid, err)
		return
	}
	if q == nil || q.State == state {
		return
	}
	if r, err = s.dao.UpdateStatus(c, state, qid, operator); err != nil || r != 1 {
		return
	}
	return
}

// BatchUpdateState bacth update question state.
func (s *Service) BatchUpdateState(c context.Context, qids []int64, state int8, operator string) (err error) {
	for _, id := range qids {
		s.UpdateStatus(c, id, state, operator)
	}
	return
}

// Types question type
func (s *Service) Types(c context.Context) (res []*model.TypeInfo, err error) {
	return s.dao.Types(c)
}

// ReadCsv read csv file
func (s *Service) ReadCsv(f multipart.File, h *multipart.FileHeader) (rs [][]string, err error) {
	r := csv.NewReader(f)
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("upload question ReadCsv error(%v)", err)
			break
		}
		if len(record) == model.ArgsCount {
			rs = append(rs, record)
		}
	}
	return
}

// UploadQsts upload questions
func (s *Service) UploadQsts(c context.Context, f multipart.File, h *multipart.FileHeader, operator string) (msg string, err error) {
	defer f.Close()
	if h != nil && !strings.HasSuffix(h.Filename, ".csv") {
		msg = "not csv file."
		return
	}
	sz, ok := f.(model.Sizer)
	if !ok {
		msg = "get file size faild."
		return
	}
	size := sz.Size()
	if size > model.FileMaxSize {
		msg = "file size more than 2M."
		return
	}
	rs, err := s.ReadCsv(f, h)
	log.Info("file %s, len(%d)", h.Filename, len(rs))
	if len(rs) == 0 || len(rs) > model.MaxCount {
		msg = "file size count is 0 or more than " + strconv.FormatInt(model.MaxCount, 10)
		return
	}
	for _, r := range rs {
		typeID, err := strconv.ParseInt(r[0], 10, 8)
		if err != nil {
			log.Error("strconv.ParseInt(%+v) err(%v)", r[0], err)
		}
		if err == nil {
			q := &model.QuestionDB{
				TypeID:   int8(typeID),
				Question: r[1],
				Ans1:     r[2],
				Ans2:     r[3],
				Ans3:     r[4],
				Ans4:     r[5],
				Operator: operator,
			}
			if err = s.QuestionAdd(c, q); err != nil {
				log.Error("s.QuestionAdd(%+v) error(%v)", q, err)
			}
		}
	}
	return
}

// QuestionAdd add register question
func (s *Service) QuestionAdd(c context.Context, q *model.QuestionDB) (err error) {
	if len(q.Question) < model.MinQuestion || len(q.Question) > model.MaxQuestion {
		err = ecode.QuestionStrNotAllow
		return
	}
	if len(q.Ans1) < model.MinAns || len(q.Ans1) > model.MaxAns ||
		len(q.Ans2) < model.MinAns || len(q.Ans2) > model.MaxAns ||
		len(q.Ans3) < model.MinAns || len(q.Ans3) > model.MaxAns ||
		len(q.Ans4) < model.MinAns || len(q.Ans4) > model.MaxAns {
		err = ecode.QuestionAnsNotAllow
		return
	}
	if q.Tips != "" && (len(q.Tips) < model.MinTips || len(q.Tips) > model.MaxTips) {
		err = ecode.QuestionTipsNotAllow
		return
	}
	if q.TypeID <= 0 {
		err = ecode.QuestionTypeNotAllow
		return
	}
	// only sourport text question
	q.MediaType = model.TextMediaType
	q.State = model.PassCheck
	q.Ctime = time.Now()
	if _, err = s.dao.QuestionAdd(c, q); err != nil {
		return
	}
	qid := q.ID
	s.eventChan.Save(func() {
		s.CreateBFSImg(context.Background(), []int64{qid})
	})
	return
}

func (s *Service) loadtypes(c context.Context) (t map[int64]*model.TypeInfo, err error) {
	var tys []*model.TypeInfo
	tys, err = s.dao.Types(c)
	if err != nil {
		log.Error("s.dao.Types error(%v)", err)
		return
	}
	t = make(map[int64]*model.TypeInfo)
	for _, v := range tys {
		if v.Parentid == 0 && t[v.ID] == nil {
			t[v.ID] = &model.TypeInfo{ID: v.ID, Name: v.Name, Subs: []*model.SubType{}}
		} else if t[v.Parentid] != nil {
			t[v.Parentid].Subs = append(t[v.Parentid].Subs, &model.SubType{ID: v.ID, Name: v.Name, LabelName: v.LabelName})
		}
	}
	return
}

// QuestionEdit .
func (s *Service) QuestionEdit(c context.Context, arg *model.QuestionDB) (aff int64, err error) {
	if aff, err = s.dao.QuestionEdit(c, arg); err != nil {
		return
	}
	s.eventChan.Save(func() {
		s.CreateBFSImg(context.Background(), []int64{arg.ID})
	})
	return
}

// LoadTypes .
func (s *Service) LoadTypes(c context.Context) (err error) {
	var allType = []*model.TypeInfo{
		{ID: 1, Parentid: 0, Name: "游戏"},
		{ID: 2, Parentid: 0, Name: "影视"},
		{ID: 3, Parentid: 0, Name: "科技"},
		{ID: 4, Parentid: 0, Name: "动画"},
		{ID: 5, Parentid: 0, Name: "艺术"},
		{ID: 6, Parentid: 0, Name: "流行前线"},
		{ID: 7, Parentid: 0, Name: "鬼畜"},
		{ID: 8, Parentid: 1, Name: "动作射击", LabelName: "游戏"},
		{ID: 9, Parentid: 1, Name: "冒险格斗", LabelName: "游戏"},
		{ID: 12, Parentid: 1, Name: "策略模拟 ", LabelName: "游戏"},
		{ID: 13, Parentid: 1, Name: "角色扮演 ", LabelName: "游戏"},
		{ID: 14, Parentid: 1, Name: "音乐体育 ", LabelName: "游戏"},
		{ID: 15, Parentid: 2, Name: "纪录片 ", LabelName: "影视"},
		{ID: 16, Parentid: 2, Name: "电影 ", LabelName: "影视"},
		{ID: 17, Parentid: 2, Name: "电视剧 ", LabelName: "影视"},
		{ID: 18, Parentid: 3, Name: "军事 ", LabelName: "科技"},
		{ID: 19, Parentid: 3, Name: "地理 ", LabelName: "科技"},
		{ID: 20, Parentid: 3, Name: "历史 ", LabelName: "科技"},
		{ID: 21, Parentid: 3, Name: "文学 ", LabelName: "科技"},
		{ID: 22, Parentid: 3, Name: "数学 ", LabelName: "科技"},
		{ID: 23, Parentid: 3, Name: "物理 ", LabelName: "科技"},
		{ID: 24, Parentid: 3, Name: "化学 ", LabelName: "科技"},
		{ID: 25, Parentid: 3, Name: "生物 ", LabelName: "科技"},
		{ID: 26, Parentid: 3, Name: "数码科技 ", LabelName: "科技"},
		{ID: 27, Parentid: 4, Name: "动画声优 ", LabelName: "动画"},
		{ID: 28, Parentid: 4, Name: "动漫内容 ", LabelName: "动画"},
		{ID: 29, Parentid: 5, Name: "ACG音乐 ", LabelName: "艺术"},
		{ID: 30, Parentid: 5, Name: "三次元音乐 ", LabelName: "艺术"},
		{ID: 31, Parentid: 5, Name: "绘画 ", LabelName: "艺术"},
		{ID: 32, Parentid: 6, Name: "娱乐 ", LabelName: "流行前线"},
		{ID: 33, Parentid: 6, Name: "时尚 ", LabelName: "流行前线"},
		{ID: 34, Parentid: 6, Name: "运动 ", LabelName: "流行前线"},
		{ID: 35, Parentid: 7, Name: "鬼畜 ", LabelName: "鬼畜"},
		{ID: 36, Parentid: 0, Name: "基础题", LabelName: "基础题"},
	}
	for _, v := range allType {
		if _, err := s.dao.TypeSave(context.Background(), v); err != nil {
			log.Error("s.dao.TypeSave(%+v) err(%v)", v, err)
		}
	}
	return
}

// LoadImg .
func (s *Service) LoadImg(c context.Context) (err error) {
	qss, err := s.dao.AllQS(c)
	if err != nil {
		log.Error("s.dao.AllQS() err(%v)", err)
	}
	for _, qs := range qss {
		lastID := qs.ID
		if err = s.eventChan.Save(func() {
			s.CreateBFSImg(context.Background(), []int64{lastID})
		}); err != nil {
			log.Error("s.CreateBFSImg(%d) err(%v)", lastID, err)
		}
	}
	return
}

// QueHistory .
func (s *Service) QueHistory(c context.Context, arg *model.ArgHistory) (res *model.HistoryPage, err error) {
	res = &model.HistoryPage{}
	if res.Total, err = s.dao.HistoryCount(c, arg); err != nil {
		return
	}
	res.Items = []*model.AnswerHistoryDB{}
	if res.Total > 0 {
		if res.Items, err = s.dao.QueHistory(c, arg); err != nil {
			return
		}
	}
	return
}

// History .
func (s *Service) History(c context.Context, arg *model.ArgHistory) (res *model.HistoryPage, err error) {
	if arg.Pn <= 0 || arg.Ps <= 0 {
		arg.Pn, arg.Ps = 1, 1000
	}
	res = &model.HistoryPage{}
	if res.Items, err = s.dao.HistoryES(c, arg); err != nil {
		return
	}
	res.Total = int64(len(res.Items))
	return
}
