package http

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/aegis/model"
	"go-common/library/database/elastic"
	"go-common/library/log"
	"go-common/library/xstr"
)

//const
const (
	LogMinYear = 2018
	LogMinTime = "2018-11-01 10:00:00"
)

func setESParams(r *elastic.Request, args interface{}, cm model.EsCommon) {
	av := reflect.ValueOf(args)
	if av.Kind() == reflect.Ptr {
		av = av.Elem()
	}
	if av.Kind() != reflect.Struct {
		return
	}

	atp := av.Type()
	ranges := map[string]map[string]interface{}{}
	for i := atp.NumField() - 1; i >= 0; i-- {
		fdt := atp.Field(i)
		tag := fdt.Tag.Get("reflect")
		if tag == "ignore" || tag == "" {
			continue
		}

		fdv := av.Field(i)
		fdk := fdt.Type.Kind()
		if (fdk == reflect.Slice || fdk == reflect.String) && fdv.Len() == 0 {
			continue
		}

		//default处理
		omitdefault := strings.Index(tag, ",omitdefault")
		tag = strings.Replace(tag, ",omitdefault", "", -1)
		fdvv := fdv.Interface()
		fdvslice := false
		if omitdefault > -1 && fmt.Sprintf("%v", fdvv) == fdt.Tag.Get("default") {
			continue
		}

		//字段值处理，parse额外处理
		switch fdk {
		case reflect.Int64, reflect.Int32, reflect.Int, reflect.Int16, reflect.Int8:
			fdvv = fdv.Int()
		case reflect.String:
			if fdv.Len() == 0 {
				continue
			}

			v := fdv.String()
			parse := fdt.Tag.Get("parse")
			if parse == "int" {
				vi, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					log.Error("setESParams strconv.ParseInt(%s) error(%v)", v, err)
					continue
				}
				fdvv = vi
			} else if parse == "[]int" {
				vi, err := xstr.SplitInts(v)
				if err != nil {
					log.Error("setESParams xstr.SplitInts(%s) error(%v)", v, err)
					continue
				}
				fdvv = vi
				fdvslice = true
			} else {
				fdvv = v
			}
		case reflect.Slice:
			if fdv.Len() == 0 {
				continue
			}
			fdvslice = true
		default:
			log.Warn("setESParams not support kind(%s) for tag(%s)", fdk.String(), tag)
			continue
		}

		//范围处理
		from := strings.Index(tag, ",from")
		to := strings.Index(tag, ",to")
		if from > -1 {
			if _, exist := ranges[tag[:from]]; !exist {
				ranges[tag[:from]] = map[string]interface{}{}
			}
			ranges[tag[:from]]["from"] = fdvv
			continue
		}
		if to > -1 {
			if _, exist := ranges[tag[:to]]; !exist {
				ranges[tag[:to]] = map[string]interface{}{}
			}
			ranges[tag[:to]]["to"] = fdvv
			continue
		}

		if fdvslice {
			r.WhereIn(tag, fdvv)
		} else {
			r.WhereEq(tag, fdvv)
		}
	}

	for field, items := range ranges {
		r.WhereRange(field, items["from"], items["to"], elastic.RangeScopeLcRc)
	}

	r.Ps(cm.Ps).Pn(cm.Pn)
	order := []map[string]string{}
	if cm.Order != "" || cm.Sort != "" {
		r.Order(cm.Order, cm.Sort)
		order = append(order, map[string]string{cm.Order: cm.Sort})
	}
	if cm.Group != "" {
		r.GroupBy(elastic.EnhancedModeDistinct, cm.Group, order)
	}
}

//QueryLogSearch .
func (d *Dao) QueryLogSearch(c context.Context, args *model.ParamsQueryLog, cm model.EsCommon) (resp *model.SearchLogResult, err error) {
	var (
		min                int = LogMinYear
		max                int
		ctimefrom, ctimeto time.Time
	)
	//默认获取所有行为日志，确定了时间范围的，只查询该段范围内的日志
	if args.CtimeFrom != "" {
		ctimefrom, _ = time.ParseInLocation("2006-01-02 15:04:05", args.CtimeFrom, time.Local)
		if ctimefrom.Year() > min {
			min = ctimefrom.Year()
		}
	}
	if args.CtimeTo != "" {
		ctimeto, _ = time.ParseInLocation("2006-01-02 15:04:05", args.CtimeTo, time.Local)
		if ctimeto.Year() >= min {
			max = ctimeto.Year()
		}
	} else {
		max = time.Now().Year()
	}

	tmpl := ",log_audit_%d_%d"
	index := ""
	for i := min; i <= max; i++ {
		index += fmt.Sprintf(tmpl, args.Business, i)
	}
	index = strings.TrimLeft(index, ",")

	r := d.es.NewRequest("log_audit").Index(index).Fields(
		"uid",
		"uname",
		"oid",
		"type",
		"action",
		"str_0",
		"str_1",
		"str_2",
		"int_0",
		"int_1",
		"int_2",
		"ctime",
		"extra_data")
	setESParams(r, args, cm)
	err = r.Scan(c, &resp)
	return
}
