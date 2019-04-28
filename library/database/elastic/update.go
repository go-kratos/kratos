package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/library/ecode"
)

// Update elastic upsert
type Update struct {
	*Elastic
	business string
	data     map[string][]interface{}
	insert   bool
}

// NewUpdate new a request every update
func (e *Elastic) NewUpdate(business string) *Update {
	return &Update{
		Elastic:  e,
		business: business,
		data:     make(map[string][]interface{}),
	}
}

// IndexByMod index by mod
func (us *Update) IndexByMod(prefix string, val, mod int64) string {
	tmp := mod - 1
	var digit int
	for tmp > 0 {
		tmp /= 10
		digit++
	}
	format := fmt.Sprintf("%s_%%0%dd", prefix, digit)
	return fmt.Sprintf(format, val%mod)
}

// IndexByTime index by time
func (us *Update) IndexByTime(prefix string, typ indexType, t time.Time) (index string) {
	year := t.Format("2006")
	month := t.Format("01")
	switch typ {
	case IndexTypeYear:
		index = strings.Join([]string{prefix, year}, "_")
	case IndexTypeMonth:
		index = strings.Join([]string{prefix, year, month}, "_")
	case IndexTypeDay:
		day := t.Format("02")
		index = strings.Join([]string{prefix, year, month, day}, "_")
	case IndexTypeWeek:
		index = strings.Join([]string{prefix, year, month, weeks[t.Day()/8]}, "_")
	}
	return
}

// AddData add data items to request 'data' param
func (us *Update) AddData(index string, data interface{}) *Update {
	if data == nil {
		return us
	}
	us.data[index] = append(us.data[index], data)
	return us
}

// HasData weather data is empty or not
func (us *Update) HasData() bool {
	if us.data == nil {
		return false
	}
	return len(us.data) > 0
}

// Insert set insert flag, it means 'replace'
func (us *Update) Insert() *Update {
	us.insert = true
	return us
}

// Do post a request
func (us *Update) Do(ctx context.Context) (err error) {
	data, err := json.Marshal(us.data)
	if err != nil {
		return
	}
	params := url.Values{}
	params.Add("business", us.business)
	params.Add("data", string(data))
	params.Add("insert", fmt.Sprintf("%t", us.insert))
	response := new(response)
	if err = us.client.Post(ctx, us.c.Host+_pathUpsert, "", params, &response); err != nil {
		return
	}
	if !ecode.Int(response.Code).Equal(ecode.OK) {
		err = ecode.Int(response.Code)
	}
	return
}

// Params get query parameters
func (us *Update) Params() string {
	data, _ := json.Marshal(us.data)
	return fmt.Sprintf("business=%s&insert=%t&data=%s", us.business, us.insert, data)
}
