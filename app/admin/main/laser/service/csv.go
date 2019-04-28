package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"go-common/app/admin/main/laser/model"
	"strconv"
	"time"
)

var (
	csvMetaNodes = []model.CsvMetaNode{
		{Index: 0, Name: "日期", DataCode: 25},
		{Index: 1, Name: "操作人", DataCode: 26},
		{Index: 2, Name: "总操视频量", DataCode: model.TotalVideo},
		{Index: 3, Name: "总操作次数", DataCode: model.TotalVideoOper},
		{Index: 4, Name: "开放浏视频量", DataCode: model.OpenVideo},
		{Index: 5, Name: "开放浏览操作次数", DataCode: model.OpenVideoOper},
		{Index: 6, Name: "会员可视频量", DataCode: model.VipAccessVideo},
		{Index: 7, Name: "会员可见操作次数", DataCode: model.VipAccessVideoOper},
		{Index: 8, Name: "打视频量", DataCode: model.RejectVideo},
		{Index: 9, Name: "打回操作次数", DataCode: model.RejectVideoOper},
		{Index: 10, Name: "锁视频量", DataCode: model.LockVideo},
		{Index: 11, Name: "锁定操作次数", DataCode: model.LockVideoOper},
		{Index: 12, Name: "通过视频总时长", DataCode: model.PassVideoTotalDuration},
		{Index: 13, Name: "未通过视频总时长", DataCode: model.FailVideoTotalDuration},
		{Index: 14, Name: "视频提交到进入待审平均响应时间", DataCode: model.WaitAuditAvgTime},
		{Index: 15, Name: "视频提交到进入待审时间", DataCode: model.WaitAuditDuration},
		{Index: 16, Name: "视频提交到进入待审次数", DataCode: model.WaitAuditOper},
	}
)

// FormatCSV format to csv data
func FormatCSV(records [][]string) (data []byte, err error) {
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

func formatAuditCargo(wrappers []*model.CargoViewWrapper, lineWidth int) (data [][]string) {
	size := len(wrappers)
	if size <= 0 {
		return
	}

	data = make([][]string, lineWidth+1)
	index := 0
	data[index] = []string{"username", "审核时间段", "接收量", "完成量"}
	for _, v1 := range wrappers {
		for k2, v2 := range v1.Data {
			data[index+1] = []string{
				v1.Username,
				fmt.Sprintf("%s %d:00:00", v1.Date, k2),
				strconv.FormatInt(v2.ReceiveValue, 10),
				strconv.FormatInt(v2.AuditValue, 10),
			}
			index = index + 1
		}
	}
	return
}

func formatVideoAuditStat(statViewExts []*model.StatViewExt, lineWidth int) (data [][]string) {
	if lineWidth <= 0 {
		return
	}
	data = make([][]string, lineWidth+1)
	index := 0
	rowHeight := len(csvMetaNodes)
	titles := make([]string, rowHeight)
	cursorMap := make(map[int]int)
	for _, v := range csvMetaNodes {
		titles[v.Index] = v.Name
		cursorMap[v.DataCode] = v.Index
	}
	data[index] = titles
	for _, v1 := range statViewExts {
		date := time.Unix(v1.Date, 0).Format("2006-01-02")
		for _, v2 := range v1.Wraps {
			name := v2.Uname
			tempRows := make([]string, rowHeight)
			tempRows = append([]string{date, name}, tempRows[0:]...)
			for _, v3 := range v2.Stats {
				if cursor, ok := cursorMap[v3.DataCode]; ok {
					if v3.DataCode == model.WaitAuditAvgTime || v3.DataCode == model.WaitAuditDuration || v3.DataCode == model.PassVideoTotalDuration || v3.DataCode == model.FailVideoTotalDuration {
						tempRows[cursor] = fmt.Sprintf("%d:%d:%d", v3.Value/3600, v3.Value%3600/60, v3.Value%3600%60/1)
					} else {
						tempRows[cursor] = strconv.FormatInt(v3.Value, 10)
					}
				}
			}
			data[index+1] = tempRows
			index = index + 1
		}
	}
	return
}
