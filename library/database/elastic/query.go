package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go-common/library/ecode"
	httpx "go-common/library/net/http/blademaster"
	timex "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	// OrderAsc order ascend
	OrderAsc = "asc"
	// OrderDesc order descend
	OrderDesc = "desc"

	// RangeScopeLoRo left open & right open
	RangeScopeLoRo rangeScope = "( )"
	// RangeScopeLoRc left open & right close
	RangeScopeLoRc rangeScope = "( ]"
	// RangeScopeLcRo left close & right open
	RangeScopeLcRo rangeScope = "[ )"
	// RangeScopeLcRc lect close & right close
	RangeScopeLcRc rangeScope = "[ ]"

	// NotTypeEq not type eq
	NotTypeEq notType = "eq"
	// NotTypeIn not type in
	NotTypeIn notType = "in"
	// NotTypeRange not type range
	NotTypeRange notType = "range"

	// LikeLevelHigh wildcard keyword
	LikeLevelHigh likeLevel = "high"
	// LikeLevelMiddle ngram(1,2)
	LikeLevelMiddle likeLevel = "middle"
	// LikeLevelLow match split word
	LikeLevelLow likeLevel = "low"

	// IndexTypeYear index by year
	IndexTypeYear indexType = "year"
	// IndexTypeMonth index by month
	IndexTypeMonth indexType = "month"
	// IndexTypeWeek index by week
	IndexTypeWeek indexType = "week"
	// IndexTypeDay index by day
	IndexTypeDay indexType = "day"

	// EnhancedModeGroupBy group by mode
	EnhancedModeGroupBy enhancedType = "group_by"
	// EnhancedModeDistinct distinct mode
	EnhancedModeDistinct enhancedType = "distinct"
	// EnhancedModeSum sum mode
	EnhancedModeSum enhancedType = "sum"
)

type (
	notType      string
	rangeScope   string
	likeLevel    string
	indexType    string
	enhancedType string
)

var (
	_defaultHost = "http://manager.bilibili.co"
	_pathQuery   = "/x/admin/search/query"
	_pathUpsert  = "/x/admin/search/upsert"

	_defaultHTTPClient = &httpx.ClientConfig{
		App: &httpx.App{
			Key:    "3c4e41f926e51656",
			Secret: "26a2095b60c24154521d24ae62b885bb",
		},
		Dial:    timex.Duration(time.Second),
		Timeout: timex.Duration(time.Second),
	}

	// for index by week
	weeks = map[int]string{0: "0107", 1: "0815", 2: "1623", 3: "2431"}
)

// Config Elastic config
type Config struct {
	Host       string
	HTTPClient *httpx.ClientConfig
}

// response query elastic response
type response struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

type query struct {
	Fields          []string            `json:"fields"`
	From            string              `json:"from"`
	OrderScoreFirst bool                `json:"order_score_first"`
	OrderRandomSeed string              `json:"order_random_seed"`
	Highlight       bool                `json:"highlight"`
	Pn              int                 `json:"pn"`
	Ps              int                 `json:"ps"`
	Order           []map[string]string `json:"order,omitempty"`
	Where           where               `json:"where,omitempty"`
}

func (q *query) string() (string, error) {
	var (
		sli []string
		m   = make(map[string]bool)
	)
	for _, i := range strings.Split(q.From, ",") {
		if m[i] {
			continue
		}
		m[i] = true
		sli = append(sli, i)
	}
	q.From = strings.Join(sli, ",")
	bs, err := json.Marshal(q)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

type where struct {
	GroupBy  string                      `json:"group_by,omitempty"`
	Like     []whereLike                 `json:"like,omitempty"`
	Eq       map[string]interface{}      `json:"eq,omitempty"`
	Or       map[string]interface{}      `json:"or,omitempty"`
	In       map[string][]interface{}    `json:"in,omitempty"`
	Range    map[string]string           `json:"range,omitempty"`
	Combo    []*Combo                    `json:"combo,omitempty"`
	Not      map[notType]map[string]bool `json:"not,omitempty"`
	Enhanced []interface{}               `json:"enhanced,omitempty"`
}

type whereLike struct {
	Fields []string  `json:"kw_fields"`
	Words  []string  `json:"kw"`
	Or     bool      `json:"or"`
	Level  likeLevel `json:"level"`
}

// Combo mix eq & in & range
type Combo struct {
	EQ       []map[string]interface{}   `json:"eq,omitempty"`
	In       []map[string][]interface{} `json:"in,omitempty"`
	Range    []map[string]string        `json:"range,omitempty"`
	NotEQ    []map[string]interface{}   `json:"not_eq,omitempty"`
	NotIn    []map[string][]interface{} `json:"not_in,omitempty"`
	NotRange []map[string]string        `json:"not_range,omitempty"`
	Min      struct {
		EQ       int `json:"eq,omitempty"`
		In       int `json:"in,omitempty"`
		Range    int `json:"range,omitempty"`
		NotEQ    int `json:"not_eq,omitempty"`
		NotIn    int `json:"not_in,omitempty"`
		NotRange int `json:"not_range,omitempty"`
		Min      int `json:"min"`
	} `json:"min"`
}

// ComboEQ .
func (cmb *Combo) ComboEQ(eq []map[string]interface{}) *Combo {
	cmb.EQ = append(cmb.EQ, eq...)
	return cmb
}

// ComboRange .
func (cmb *Combo) ComboRange(r []map[string]string) *Combo {
	cmb.Range = append(cmb.Range, r...)
	return cmb
}

// ComboIn .
func (cmb *Combo) ComboIn(in []map[string][]interface{}) *Combo {
	cmb.In = append(cmb.In, in...)
	return cmb
}

// ComboNotEQ .
func (cmb *Combo) ComboNotEQ(eq []map[string]interface{}) *Combo {
	cmb.NotEQ = append(cmb.NotEQ, eq...)
	return cmb
}

// ComboNotRange .
func (cmb *Combo) ComboNotRange(r []map[string]string) *Combo {
	cmb.NotRange = append(cmb.NotRange, r...)
	return cmb
}

// ComboNotIn .
func (cmb *Combo) ComboNotIn(in []map[string][]interface{}) *Combo {
	cmb.NotIn = append(cmb.NotIn, in...)
	return cmb
}

// MinEQ .
func (cmb *Combo) MinEQ(min int) *Combo {
	cmb.Min.EQ = min
	return cmb
}

// MinIn .
func (cmb *Combo) MinIn(min int) *Combo {
	cmb.Min.In = min
	return cmb
}

// MinRange .
func (cmb *Combo) MinRange(min int) *Combo {
	cmb.Min.Range = min
	return cmb
}

// MinNotEQ .
func (cmb *Combo) MinNotEQ(min int) *Combo {
	cmb.Min.NotEQ = min
	return cmb
}

// MinNotIn .
func (cmb *Combo) MinNotIn(min int) *Combo {
	cmb.Min.NotIn = min
	return cmb
}

// MinNotRange .
func (cmb *Combo) MinNotRange(min int) *Combo {
	cmb.Min.NotRange = min
	return cmb
}

// MinAll .
func (cmb *Combo) MinAll(min int) *Combo {
	cmb.Min.Min = min
	return cmb
}

type groupBy struct {
	Mode  enhancedType        `json:"mode"`
	Field string              `json:"field"`
	Order []map[string]string `json:"order"`
}

type enhance struct {
	Mode  enhancedType        `json:"mode"`
	Field string              `json:"field"`
	Order []map[string]string `json:"order,omitempty"`
	Size  int                 `json:"size,omitempty"`
}

// Elastic clastic instance
type Elastic struct {
	c      *Config
	client *httpx.Client
}

// NewElastic .
func NewElastic(c *Config) *Elastic {
	if c == nil {
		c = &Config{
			Host:       _defaultHost,
			HTTPClient: _defaultHTTPClient,
		}
	}
	return &Elastic{
		c:      c,
		client: httpx.NewClient(c.HTTPClient),
	}
}

// Request request to elastic
type Request struct {
	*Elastic
	q        query
	business string
}

// NewRequest new a request every search query
func (e *Elastic) NewRequest(business string) *Request {
	return &Request{
		Elastic:  e,
		business: business,
		q: query{
			Fields:          []string{},
			Highlight:       false,
			OrderScoreFirst: true,
			OrderRandomSeed: "",
			Pn:              1,
			Ps:              10,
		},
	}
}

// Fields add query fields
func (r *Request) Fields(fields ...string) *Request {
	r.q.Fields = append(r.q.Fields, fields...)
	return r
}

// Index add query index
func (r *Request) Index(indexes ...string) *Request {
	r.q.From = strings.Join(indexes, ",")
	return r
}

// IndexByMod index by mod
func (r *Request) IndexByMod(prefix string, val, mod int64) *Request {
	tmp := mod - 1
	var digit int
	for tmp > 0 {
		tmp /= 10
		digit++
	}
	format := fmt.Sprintf("%s_%%0%dd", prefix, digit)
	r.q.From = fmt.Sprintf(format, val%mod)
	return r
}

// IndexByTime index by time
func (r *Request) IndexByTime(prefix string, typ indexType, begin, end time.Time) *Request {
	var (
		buf     bytes.Buffer
		index   string
		indexes = make(map[string]struct{})
	)
	for {
		year := begin.Format("2006")
		month := begin.Format("01")
		switch typ {
		case IndexTypeYear:
			index = strings.Join([]string{prefix, year}, "_")
		case IndexTypeMonth:
			index = strings.Join([]string{prefix, year, month}, "_")
		case IndexTypeDay:
			day := begin.Format("02")
			index = strings.Join([]string{prefix, year, month, day}, "_")
		case IndexTypeWeek:
			index = strings.Join([]string{prefix, year, month, weeks[begin.Day()/8]}, "_")
		}
		if begin.After(end) && begin.Day() != end.Day() {
			break
		}
		indexes[index] = struct{}{}
		begin = begin.AddDate(0, 0, 1)
	}
	for i := range indexes {
		buf.WriteString(i)
		buf.WriteString(",")
	}
	r.q.From = strings.TrimSuffix(buf.String(), ",")
	return r
}

// OrderScoreFirst switch for order score first
func (r *Request) OrderScoreFirst(v bool) *Request {
	r.q.OrderScoreFirst = v
	return r
}

// OrderRandomSeed switch for order random
func (r *Request) OrderRandomSeed(v string) *Request {
	r.q.OrderRandomSeed = v
	return r
}

// Highlight switch from highlight
func (r *Request) Highlight(v bool) *Request {
	r.q.Highlight = v
	return r
}

// Pn page number
func (r *Request) Pn(v int) *Request {
	r.q.Pn = v
	return r
}

// Ps page size
func (r *Request) Ps(v int) *Request {
	r.q.Ps = v
	return r
}

// Order filed sort
func (r *Request) Order(field, sort string) *Request {
	if sort != OrderAsc {
		sort = OrderDesc
	}
	r.q.Order = append(r.q.Order, map[string]string{field: sort})
	return r
}

// WhereEq where qual
func (r *Request) WhereEq(field string, eq interface{}) *Request {
	if r.q.Where.Eq == nil {
		r.q.Where.Eq = make(map[string]interface{})
	}
	r.q.Where.Eq[field] = eq
	return r
}

// WhereOr where or
func (r *Request) WhereOr(field string, or interface{}) *Request {
	if r.q.Where.Or == nil {
		r.q.Where.Or = make(map[string]interface{})
	}
	r.q.Where.Or[field] = or
	return r
}

// WhereIn where in
func (r *Request) WhereIn(field string, in interface{}) *Request {
	if r.q.Where.In == nil {
		r.q.Where.In = make(map[string][]interface{})
	}
	switch v := in.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, string:
		r.q.Where.In[field] = append(r.q.Where.In[field], v)
	case []int:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int64:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []string:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int8:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int16:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []int32:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint8:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint16:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint32:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	case []uint64:
		for _, i := range v {
			r.q.Where.In[field] = append(r.q.Where.In[field], i)
		}
	}
	return r
}

// WhereRange where range
func (r *Request) WhereRange(field string, start, end interface{}, scope rangeScope) *Request {
	if r.q.Where.Range == nil {
		r.q.Where.Range = make(map[string]string)
	}
	if start == nil {
		start = ""
	}
	if end == nil {
		end = ""
	}
	switch scope {
	case RangeScopeLoRo:
		r.q.Where.Range[field] = fmt.Sprintf("(%v,%v)", start, end)
	case RangeScopeLoRc:
		r.q.Where.Range[field] = fmt.Sprintf("(%v,%v]", start, end)
	case RangeScopeLcRo:
		r.q.Where.Range[field] = fmt.Sprintf("[%v,%v)", start, end)
	case RangeScopeLcRc:
		r.q.Where.Range[field] = fmt.Sprintf("[%v,%v]", start, end)
	}
	return r
}

// WhereNot where not
func (r *Request) WhereNot(typ notType, fields ...string) *Request {
	if r.q.Where.Not == nil {
		r.q.Where.Not = make(map[notType]map[string]bool)
	}
	if r.q.Where.Not[typ] == nil {
		r.q.Where.Not[typ] = make(map[string]bool)
	}
	for _, v := range fields {
		r.q.Where.Not[typ][v] = true
	}
	return r
}

// WhereLike where like
func (r *Request) WhereLike(fields, words []string, or bool, level likeLevel) *Request {
	if len(fields) == 0 || len(words) == 0 {
		return r
	}
	l := whereLike{Fields: fields, Words: words, Or: or, Level: level}
	r.q.Where.Like = append(r.q.Where.Like, l)
	return r
}

// WhereCombo where combo
func (r *Request) WhereCombo(cmb ...*Combo) *Request {
	r.q.Where.Combo = append(r.q.Where.Combo, cmb...)
	return r
}

// GroupBy where group by
func (r *Request) GroupBy(mode enhancedType, field string, order []map[string]string) *Request {
	for _, i := range order {
		for k, v := range i {
			if v != OrderAsc {
				i[k] = OrderDesc
			}
		}
	}
	r.q.Where.Enhanced = append(r.q.Where.Enhanced, groupBy{Mode: mode, Field: field, Order: order})
	return r
}

// Sum where enhance sum
func (r *Request) Sum(field string) *Request {
	r.q.Where.Enhanced = append(r.q.Where.Enhanced, enhance{Mode: EnhancedModeSum, Field: field})
	return r
}

// Scan parse the query response data
func (r *Request) Scan(ctx context.Context, result interface{}) (err error) {
	q, err := r.q.string()
	if err != nil {
		return
	}
	params := url.Values{}
	params.Add("business", r.business)
	params.Add("query", q)
	response := new(response)
	if err = r.client.Get(ctx, r.c.Host+_pathQuery, "", params, &response); err != nil {
		return
	}
	if !ecode.Int(response.Code).Equal(ecode.OK) {
		return ecode.Int(response.Code)
	}
	err = errors.Wrapf(json.Unmarshal(response.Data, &result), "scan(%s)", response.Data)
	return
}

// Params get query parameters
func (r *Request) Params() string {
	q, _ := r.q.string()
	return fmt.Sprintf("business=%s&query=%s", r.business, q)
}
