package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"go-common/app/job/main/growup/model"
)

// WriteCSV write data to csv
func WriteCSV(records [][]string, filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.WriteAll(records)
	w.Flush()
	return
}

// FormatCSV format to csv data
func FormatCSV(records [][]string) (data []byte, err error) {
	buf := new(bytes.Buffer)

	// add utf bom
	buf.WriteString("\xEF\xBB\xBF")

	w := csv.NewWriter(buf)
	err = w.WriteAll(records)
	if err != nil {
		return
	}

	data = buf.Bytes()
	return
}

func formatAvIncome(list []*model.AvIncome) (data [][]string) {
	if len(list) <= 0 {
		return
	}
	data = make([][]string, len(list)+1)
	data[0] = []string{"稿件id", "UP主UID", "稿件月收入", "累计收入", "分区id", "最后收入时间"}
	for i := 0; i < len(list); i++ {
		l := list[i]
		data[i+1] = []string{
			strconv.FormatInt(l.AvID, 10),
			strconv.FormatInt(l.MID, 10),
			fmt.Sprintf("%.2f", float64(l.Income)/float64(100)),
			fmt.Sprintf("%.2f", float64(l.TotalIncome)/float64(100)),
			strconv.FormatInt(l.TagID, 10),
			l.Date.Time().Format(_layout),
		}
	}
	return
}

func formatUpAccount(list []*model.UpAccount, month int) (data [][]string) {
	if len(list) <= 0 {
		return
	}
	data = make([][]string, len(list)+1)
	data[0] = []string{"UP主UID", "昵称", fmt.Sprintf("%d月有收入稿件数", month), fmt.Sprintf("%d月收入", month), "累计收入", "待结算收入", fmt.Sprintf("%d月收入-待结算", month)}
	for i := 0; i < len(list); i++ {
		l := list[i]
		data[i+1] = []string{
			strconv.FormatInt(l.MID, 10),
			l.Nickname,
			strconv.FormatInt(l.AvCount, 10),
			fmt.Sprintf("%.2f", float64(l.MonthIncome)/float64(100)),
			fmt.Sprintf("%.2f", float64(l.TotalIncome)/float64(100)),
			fmt.Sprintf("%.2f", float64(l.TotalUnwithdrawIncome)/float64(100)),
			fmt.Sprintf("%.2f", float64(l.MonthIncome-l.TotalUnwithdrawIncome)/float64(100)),
		}
	}
	return
}
