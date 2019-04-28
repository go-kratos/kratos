package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
)

const (
	_tagIncome = "%d年%d月%d日收入消息统计(运营标签)"
)

// SendTagIncomeByHTTP exec http
func (s *Service) SendTagIncomeByHTTP(c context.Context, year int, month int, day int) (err error) {
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	return s.sendTagIncome(c, t)
}

func (s *Service) sendTagIncome(c context.Context, date time.Time) (err error) {
	// ready
	err = GetTaskService().TaskReady(c, date.Format("2006-01-02"), TaskCreativeIncome)
	if err != nil {
		return
	}

	allAvInfos, err := s.getDayAllAvInfo(c, date)
	if err != nil {
		log.Error("s.sendTagIncome getDayAllAvInfo error(%v)", err)
		return
	}
	allUpInfos, err := s.getDayAllUpInfo(c, date)
	if err != nil {
		log.Error("s.sendTagIncome getDayAllUpInfo error(%v)", err)
		return
	}
	totalIncome, err := s.calTotalIncome(c, allUpInfos)
	if err != nil {
		log.Error("s.sendTagIncome s.totalInfo error(%v)", err)
		return
	}
	mi, err := s.topThirtyMID(c, allUpInfos)
	if err != nil {
		log.Error("s.sendTagIncome s.topThirtyMID error(%v)", err)
		return
	}
	ai, err := s.topThirtyAV(c, allAvInfos)
	if err != nil {
		log.Error("s.sendTagIncome s.topThirtyAV error(%v)", err)
		return
	}
	tag, err := s.tagInfo(c, date)
	if err != nil {
		log.Error("s.sendTagIncome s.tagInfo error(%v)", err)
		return
	}
	f, err := os.Create("income.csv")
	if err != nil {
		log.Error("s.sendTagIncome create income.csv error(%v)", err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	data := [][]string{
		{"昨日UP主人数", "昨日稿件总数", "昨日收入总金额(元)"},
		{strconv.Itoa(len(allUpInfos)), strconv.Itoa(len(allAvInfos)), strconv.FormatFloat(float64(totalIncome)/100, 'f', 2, 32)},
	}
	mt := []string{"MID", "昨日收入(元)", "累计收入(元)"}
	data = append(data, mt)
	for _, o := range mi {
		m := []string{strconv.FormatInt(o.MID, 10), strconv.FormatFloat(float64(o.Income)/100, 'f', 2, 32), strconv.FormatFloat(float64(o.TotalIncome)/100, 'f', 2, 32)}
		data = append(data, m)
	}
	at := []string{"AVID", "昨日收入(元)", "累计收入(元)"}
	data = append(data, at)
	for _, o := range ai {
		a := []string{strconv.FormatInt(o.AVID, 10), strconv.FormatFloat(float64(o.Income)/100, 'f', 2, 32), strconv.FormatFloat(float64(o.TotalIncome)/100, 'f', 2, 32)}
		data = append(data, a)
	}

	tt := []string{"标签", "稿件数", "昨日收入(元)", "累计收入(元)"}
	data = append(data, tt)
	for _, o := range tag {
		t := []string{o.Tag, strconv.Itoa(o.AVCount), strconv.FormatFloat(float64(o.Income)/100, 'f', 2, 32), strconv.FormatFloat(float64(o.TotalIncome)/100, 'f', 2, 32)}
		data = append(data, t)
	}
	w.WriteAll(data)
	w.Flush()
	rfile, err := os.Open("income.csv")
	if err != nil {
		log.Error("s.sendTagIncome open income.csv error(%v)", err)
		return
	}
	defer rfile.Close()
	var body string
	r := csv.NewReader(rfile)
	for {
		var strs []string
		strs, err = r.Read()
		if err == io.EOF {
			break
		}
		body += fmt.Sprintf("<tr><td>%s</td></tr>", strings.Join(strs, "</td><td>"))
	}
	var send []string
	for _, v := range s.conf.Mail.Send {
		if v.Type == 1 {
			send = v.Addr
		}
	}
	err = s.email.SendMail(date, fmt.Sprintf("<table border='1'>%s</table>", body), _tagIncome, send...)
	if err != nil {
		log.Error("s.sendTagIncome send email error(%v)", err)
		return
	}
	err = os.Remove("income.csv")
	if err != nil {
		log.Error("s.sendTagIncome remove income.csv error(%v)", err)
	}
	return
}

// get av_income one date all av.
func (s *Service) getDayAllAvInfo(c context.Context, date time.Time) (allInfo []*model.IncomeInfo, err error) {
	from, limit := int64(0), int64(2000)
	var partInfo []*model.IncomeInfo
	for {
		partInfo, err = s.dao.TotalIncome(c, date, from, limit)
		if err != nil {
			log.Error("s.totalInfo dao.TotalIncome error(%v)", err)
			return
		}
		for _, info := range partInfo {
			if info.IsDeleted == 0 {
				allInfo = append(allInfo, info)
			}
		}
		if int64(len(partInfo)) < limit {
			break
		}
		from = partInfo[len(partInfo)-1].ID
	}
	return
}

func (s *Service) getDayAllUpInfo(c context.Context, date time.Time) (allInfo []*model.MIDInfo, err error) {
	from, limit := int64(0), int64(2000)
	var partInfo []*model.MIDInfo
	for {
		partInfo, err = s.dao.GetUpIncome(c, date, from, limit)
		if err != nil {
			log.Error("s.getDayAllUpInfo dao.GetUpIncome error(%v)", err)
			return
		}
		for _, info := range partInfo {
			if info.IsDeleted == 0 {
				allInfo = append(allInfo, info)
			}
		}
		if int64(len(partInfo)) < limit {
			break
		}
		from = partInfo[len(partInfo)-1].ID
	}
	return
}

// calculate totalIncome.
func (s *Service) calTotalIncome(c context.Context, upInfos []*model.MIDInfo) (totalIncome int64, err error) {
	for _, v := range upInfos {
		totalIncome += v.Income
	}
	return
}

// get av income top 30.
func (s *Service) topThirtyAV(c context.Context, avInfos []*model.IncomeInfo) (vs []*model.AVIDInfo, err error) {
	vs = make([]*model.AVIDInfo, 0)
	/*	tv, err := s.dao.GetAV(c, date)
		if err != nil {
			log.Error("s.topThirtyAV dao.GetAV error(%v)", err)
			return
		}*/
	sort.SliceStable(avInfos, func(i, j int) bool {
		return avInfos[i].Income > avInfos[j].Income
	})
	cnt := 30
	if len(avInfos) < 30 {
		cnt = len(avInfos)
	}
	for i := 0; i < cnt; i++ {
		a := &model.AVIDInfo{AVID: avInfos[i].AVID, Income: avInfos[i].Income, TotalIncome: avInfos[i].TotalIncome}
		vs = append(vs, a)
	}
	return
}

// get top30 av
func (s *Service) topThirtyMID(c context.Context, upInfos []*model.MIDInfo) (ms []*model.MIDInfo, err error) {
	sort.SliceStable(upInfos, func(i, j int) bool {
		return upInfos[i].Income > upInfos[j].Income

	})
	cnt := 30
	if len(upInfos) < 30 {
		cnt = len(upInfos)
	}

	for i := 0; i < cnt; i++ {
		ms = append(ms, upInfos[i])
	}
	return
}

// tagInfo get tag info
func (s *Service) tagInfo(c context.Context, date time.Time) (infos []*model.TagInfo, err error) {
	tagAvs, err := s.getTagAvs(c, date)
	if err != nil {
		log.Error("s.tagInfo getTagAvs error(%v)", err)
		return
	}
	infos, err = s.handleTag(c, tagAvs, date)
	if err != nil {
		log.Error("s.tagInfo handleTag error(%v)", err)
		return
	}
	return
}

func (s *Service) getTagAvs(c context.Context, date time.Time) (allInfo []*model.AvIncome, err error) {
	var (
		from, limit int64
		info        []*model.AvIncome
	)
	from, limit = 0, 2000
	for {
		info, err = s.tag.GetTagAvInfo(c, date, from, limit)
		if err != nil {
			log.Error("s.getTagAvs tag.GetTagAvInfo error(%v)", err)
			return
		}
		for _, i := range info {
			if i.IsDeleted == 0 {
				allInfo = append(allInfo, i)
			}
		}
		if int64(len(info)) < limit {
			break
		}
		from = info[len(info)-1].ID
	}
	return
}

func (s *Service) handleTag(c context.Context, allInfo []*model.AvIncome, date time.Time) (results []*model.TagInfo, err error) {
	tagMap := make(map[int64]*model.TagInfo) // key-value: tagID-TagInfo
	var tagIDs []int64                       // tagID
	for _, info := range allInfo {
		if _, ok := tagMap[info.TagID]; !ok {
			a := &model.TagInfo{}
			tagMap[info.TagID] = a
			tagIDs = append(tagIDs, info.TagID)
		}
		tagMap[info.TagID].AVCount++
		tagMap[info.TagID].Income += int64(info.Income)
	}
	if len(tagIDs) <= 0 {
		return
	}

	tags, err := s.dao.GetTagTotalIncome(c, tagIDs)
	if err != nil {
		log.Error("s.handleTag dao.GetTagTotalIncome error(%v)", err)
		return
	}
	infoMap := make(map[int64]*model.TagInfo) // key-value: tagID-info
	for _, tag := range tags {
		if _, ok := infoMap[tag.ID]; !ok {
			a := &model.TagInfo{ID: tag.ID, Tag: tag.Tag, TotalIncome: tag.TotalIncome}
			infoMap[tag.ID] = a
		}
	}

	for tagID, tagInfo := range tagMap {
		_, ok := infoMap[tagID]
		if ok {
			tagInfo.TotalIncome = infoMap[tagID].TotalIncome
			tagInfo.Tag = infoMap[tagID].Tag
			results = append(results, tagInfo)
		} else {
			log.Error("s.handleTag tagID:%d not exist in tag_info", tagID)
		}
	}
	return
}
