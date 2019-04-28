package http

import (
	"bytes"
	"encoding/csv"
	"strconv"

	"go-common/app/admin/main/creative/model/whitelist"
	"go-common/library/log"
)

// FormatCSV  format csv data.
func FormatCSV(records [][]string) (res []byte) {
	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Error("error writing record to csv:", err)
			return
		}
	}
	w.Flush()
	res = buf.Bytes()
	return
}

func formatWhilteList(wl []*whitelist.Whitelist) (data [][]string, err error) {
	if len(wl) < 0 {
		return
	}
	data = append(data, []string{"MID", "昵称", "AdminMID", "备注", "粉丝数", "等级", "创建时间"})
	for _, v := range wl {
		var fields []string
		fields = append(fields, strconv.FormatInt(int64(v.MID), 10))
		fields = append(fields, v.Name)
		fields = append(fields, strconv.FormatInt(int64(v.AdminMID), 10))
		fields = append(fields, v.Comment)
		fields = append(fields, strconv.FormatInt(int64(v.Fans), 10))
		fields = append(fields, strconv.FormatInt(int64(v.CurrentLevel), 10))
		fields = append(fields, v.Ctime)
		data = append(data, fields)
	}
	return
}
