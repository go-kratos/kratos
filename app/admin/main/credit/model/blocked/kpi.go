package blocked

import (
	"encoding/json"
	"strconv"

	"go-common/app/admin/main/credit/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// const kpi
const (
	RateS = int8(1)
	RateA = int8(2)
	RateB = int8(3)
	RateC = int8(4)
	RateD = int8(5)
)

// KPI is blocked_kpi model.
type KPI struct {
	ID        int64      `gorm:"column:id" json:"id"`
	UID       int64      `gorm:"column:mid" json:"uid"`
	Day       xtime.Time `gorm:"column:day" json:"day"`
	Rate      int8       `gorm:"column:rate" json:"rate"`
	Rank      int64      `gorm:"column:rank" json:"rank"`
	RankPer   int16      `gorm:"column:rank_per" json:"rank_per"`
	RankTotal int64      `gorm:"column:rank_total" json:"rank_total"`
	HStatus   int8       `gorm:"column:handler_status" json:"handler_status"`
	PStatus   int8       `gorm:"column:pendent_status" json:"pendent_status"`
	CTime     xtime.Time `gorm:"column:ctime" json:"-"`
	MTime     xtime.Time `gorm:"column:mtime" json:"-"`
}

// KPIDesc is blocked_kpi_desc model.
type KPIDesc struct {
	ID        string `json:"id"`
	UID       string `json:"uid"`
	Day       string `json:"day"`
	Rate      string `json:"rate"`
	RankPer   string `json:"rank_per"`
	RankTotal string `json:"rank_total"`
	PStatus   string `json:"pendent_status"`
}

// KPIList is kpi list.
type KPIList struct {
	Count int64  `json:"count"`
	Pn    int    `json:"pn"`
	Ps    int    `json:"ps"`
	List  []*KPI `json:"list"`
}

// TableName kpi tablename
func (*KPI) TableName() string {
	return "blocked_kpi"
}

// RateDesc is rate desc.
func RateDesc(rate int8) string {
	switch rate {
	case RateS:
		return "S"
	case RateA:
		return "A"
	case RateB:
		return "B"
	case RateC:
		return "C"
	case RateD:
		return "D"
	}
	return ""
}

// DealKPI deal with kpi data.
func DealKPI(kpis []*KPI) (data [][]string, err error) {
	var kpiDescs []*KPIDesc
	for _, v := range kpis {
		kpiDesc := &KPIDesc{
			ID:        strconv.FormatInt(v.ID, 10),
			UID:       strconv.FormatInt(v.UID, 10),
			Day:       v.Day.Time().Format(model.TimeFormatDay),
			Rate:      RateDesc(v.Rate),
			RankPer:   strconv.FormatInt(int64(v.RankPer), 10),
			RankTotal: strconv.FormatInt(int64(v.RankTotal), 10),
			PStatus:   strconv.FormatInt(int64(v.PStatus), 10),
		}
		kpiDescs = append(kpiDescs, kpiDesc)
	}
	kpiMap, _ := json.Marshal(kpiDescs)
	var objmap []map[string]string
	if err = json.Unmarshal(kpiMap, &objmap); err != nil {
		log.Error("Unmarshal(%s) error(%v)", string(kpiMap), err)
		return
	}
	data = append(data, []string{"ID", "MID", "KPI日期", "评级", "排名百分比", "参与总数", "挂件发放状态"})
	for _, v := range objmap {
		var fields []string
		fields = append(fields, v["id"])
		fields = append(fields, v["uid"])
		fields = append(fields, v["day"])
		fields = append(fields, v["rate"])
		fields = append(fields, v["rank_per"])
		fields = append(fields, v["rank_total"])
		fields = append(fields, v["pendent_status"])
		data = append(data, fields)
	}
	return
}

// KPIPoint is blocked_kpi_point model.
type KPIPoint struct {
	ID           int64      `gorm:"column:id" json:"id"`
	UID          int64      `gorm:"column:mid" json:"uid"`
	Day          xtime.Time `gorm:"column:day" json:"day"`
	Point        int16      `gorm:"column:point" json:"point"`
	ActiveDays   int16      `gorm:"column:active_days" json:"active_days"`
	VoteTotal    int        `gorm:"column:vote_total" json:"vote_total"`
	VoteRadio    int16      `gorm:"column:vote_radio" json:"vote_radio"`
	BlockedTotal int        `gorm:"column:blocked_total" json:"blocked_total"`
	CTime        xtime.Time `gorm:"column:ctime" json:"-"`
	MTime        xtime.Time `gorm:"column:mtime" json:"-"`
}

// TableName kpi_point tablename
func (*KPIPoint) TableName() string {
	return "blocked_kpi_point"
}

// KPIPointList is kpi_point list.
type KPIPointList struct {
	Count int         `json:"total_count"`
	Order string      `json:"order"`
	Sort  string      `json:"sort"`
	PN    int         `json:"pn"`
	PS    int         `json:"ps"`
	IDs   []int64     `json:"-"`
	List  []*KPIPoint `json:"list"`
}
