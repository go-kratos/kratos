package service

import (
	"bytes"
	"encoding/csv"
	"strconv"
	"time"

	"go-common/app/admin/main/up-rating/model"
)

var (
	_avCategory = map[int]string{
		0:   "默认",
		1:   "动画",
		3:   "音乐",
		129: "舞蹈",
		4:   "游戏",
		36:  "科技",
		160: "生活",
		119: "鬼畜",
		155: "时尚",
		23:  "电影",
		11:  "电视剧",
		13:  "番剧",
		167: "国创",
		165: "广告",
		5:   "娱乐",
		177: "纪录片",
		181: "影视",
	}

	_scoreField = map[model.ScoreType]string{
		model.Magnetic:   "magnetic_score",
		model.Creativity: "creativity_score",
		model.Influence:  "influence_score",
		model.Credit:     "credit_score",
	}
)

// FormatCSV format to csv data
func formatCSV(records [][]string) (data []byte, err error) {
	buf := new(bytes.Buffer)

	// add utf bom
	if len(records) > 0 {
		buf.WriteString("\xEF\xBB\xBF")
	}

	w := csv.NewWriter(buf)
	err = w.WriteAll(records)
	if err != nil {
		return
	}

	data = buf.Bytes()
	return
}

func tagDesc(tagID int) string {
	if v, ok := _avCategory[tagID]; ok {
		return v
	}
	return _avCategory[0]
}

func formatScores(ratings []*model.RatingInfo) (data [][]string) {
	if len(ratings) <= 0 {
		return
	}
	data = make([][]string, len(ratings)+1)
	data[0] = []string{"月份", "UID", "昵称", "分区", "总分", "创作力", "影响力", "信用分", "投稿量", "粉丝量"}
	for i, v := range ratings {
		data[i+1] = []string{
			v.Date,
			strconv.FormatInt(v.Mid, 10),
			v.NickName,
			tagDesc(v.TagID),
			strconv.FormatInt(v.MagneticScore, 10),
			strconv.FormatInt(v.CreativityScore, 10),
			strconv.FormatInt(v.InfluenceScore, 10),
			strconv.FormatInt(v.CreditScore, 10),
			strconv.FormatInt(v.TotalAvs, 10),
			strconv.FormatInt(v.TotalFans, 10),
		}
	}

	return
}

func scoreField(st model.ScoreType) string {
	if v, ok := _scoreField[st]; ok {
		return v
	}
	return _scoreField[model.Magnetic]
}

func cDateStr(cdate time.Time) string {
	return cdate.Format("2006-01-02")
}

func prevComputation(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month()-1, 1, 0, 0, 0, 0, time.Local)
}

func formatStatis(list []*model.RatingStatis, ctype int64) (data [][]string) {
	if len(list) <= 0 {
		return
	}
	data = make([][]string, len(list)+1)
	data[0] = []string{"分数段", "本月", "占比", "对比", "占比", "平均分"}
	switch ctype {
	case 0:
		data[0] = append(data[0], []string{"创造力", "影响力", "信用分"}...)
	case 1:
		data[0] = append(data[0], []string{"平均稿件数", "平均播放量", "平均互动量"}...)
	case 2:
		data[0] = append(data[0], []string{"平均粉丝量"}...)
	}
	for i, v := range list {
		data[i+1] = []string{
			v.Tips,
			strconv.FormatInt(v.Ups, 10),
			v.Proportion,
			strconv.FormatInt(v.Compare, 10),
			v.ComparePropor,
			strconv.FormatInt(v.Score, 10),
		}
		switch ctype {
		case 0:
			data[i+1] = append(data[i+1], []string{strconv.FormatInt(v.CreativityScore, 10), strconv.FormatInt(v.InfluenceScore, 10), strconv.FormatInt(v.CreditScore, 10)}...)
		case 1:
			data[i+1] = append(data[i+1], []string{strconv.FormatInt(v.Avs, 10), strconv.FormatInt(v.Play, 10), strconv.FormatInt(v.Coin, 10)}...)
		case 2:
			data[i+1] = append(data[i+1], []string{strconv.FormatInt(v.Fans, 10)}...)
		}
	}
	return
}
