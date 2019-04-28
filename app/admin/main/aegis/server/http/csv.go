package http

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"

	"go-common/app/admin/main/aegis/model"
	"go-common/library/log"

	"github.com/pkg/errors"
)

// CSVContentType
var (
	CSVContentType        = []string{"application/csv"}
	_              Render = CSV{}
)

// Render http reponse render.
type Render interface {
	Render(http.ResponseWriter) error
	WriteContentType(w http.ResponseWriter)
}

// CSV str.
type CSV struct {
	Content []byte
	Title   string
}

// WriteContentType fn
func (j CSV) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, CSVContentType, j.Title, "csv")
}

func writeContentType(w http.ResponseWriter, value []string, title, ext string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
	header["Content-Disposition"] = append(header["Content-Disposition"], fmt.Sprintf("attachment; filename=%s.%s", title, ext))
}

// Render (JSON) writes data with json ContentType.
func (j CSV) Render(w http.ResponseWriter) (err error) {
	if _, err = w.Write(j.Content); err != nil {
		err = errors.WithStack(err)
	}
	return
}

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

func formatReport(wl map[string][]*model.ReportFlowItem) (data [][]string, err error) {
	if len(wl) < 0 {
		return
	}
	tab := []string{"审核员"}

	for k, v := range wl {
		var fields []string
		fields = append(fields, k)
		for _, item := range v {
			fields = append(fields, fmt.Sprintf("%d/%d", item.Out, item.In))
			if len(tab) != len(v)+1 {
				tab = append(tab, item.Hour)
			}
		}
		data = append(data, fields)
	}

	data = append([][]string{tab}, data...)
	return
}

func formatReportTaskSubmit(res *model.ReportSubmitRes) (data [][]string, err error) {
	if res == nil || len(res.Order) <= 0 || len(res.Header) <= 0 || len(res.Rows) <= 0 {
		return
	}

	llen := len(res.Order)
	header := make([]string, llen)
	for i, k := range res.Order {
		header[i] = res.Header[k]
	}
	data = [][]string{header}

	for _, list := range res.Rows {
		for _, one := range list {
			fields := make([]string, llen)
			for i, k := range res.Order {
				fields[i] = one[k]
			}
			data = append(data, fields)
		}
	}
	return
}
