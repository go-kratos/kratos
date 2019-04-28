package model

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/smartystreets/goconvey/convey"
	"go-common/app/admin/main/aegis/model/common"
	"go-common/app/admin/main/aegis/model/resource"
	xtime "go-common/library/time"
)

func Test_SubmitOption(t *testing.T) {
	sopt := &SubmitOptions{
		EngineOption: EngineOption{
			BaseOptions: common.BaseOptions{
				BusinessID: 1,
				FlowID:     1,
				OID:        "oid-123456",
				RID:        1,
				Uname:      "cxf",
			},
			Forbid: &Forbid{
				Comment: "dqqq",
				Reason:  "我的",
			},
			TaskID: 1,
			Result: &resource.Result{
				Attribute:    3,
				Note:         "哎",
				RejectReason: "",
				ReasonID:     3,
				AttributeList: map[string]int8{
					"no_comment": 1,
					"no_forward": 0,
				},
			},
			ExtraData: map[string]interface{}{
				"mailto": "chenxuefeng",
				"mail":   1,
				"mid":    10086,
				"notify": 1,
			},
		},
		Binds: []int64{2, 3},
	}

	mapformat := map[string]*ActionParam{
		"auditor": {
			Value: "Uname",
		},

		"oid": {
			Value: "OID",
		},

		"uid": {
			Value: "ExtraData.mid",
		},
		"forbid_params": {
			Value: "Forbid",
		},

		"no_comment": {
			Value: "Result.AttributeList.no_comment",
		},
		"no_forward": {
			Value: "Result.AttributeList.no_forward",
		},

		"notify": {
			Value: "ExtraData.notify",
		},
		"reason": {
			Value:   "Result.RejectReason",
			Default: " ",
		},
	}

	ot := reflect.TypeOf(*sopt)
	ov := reflect.ValueOf(*sopt)

	params := make(map[string]interface{})
	for k, v := range mapformat {
		SubReflect(ot, ov, k, strings.Split(v.Value, "."), v.Default, params)
	}
	fmt.Println("params:", params)

	values := url.Values{}

	for k, v := range params {
		values.Set(k, fmt.Sprint(v))
	}
	fmt.Println("values:", values.Encode())

	t.Fail()
}

func Test_reflectMap(t *testing.T) {
	mv := map[string]int8{"1": 1, "2": 2}
	v := reflect.ValueOf(mv)

	fmt.Println("a:", v.MapIndex(reflect.ValueOf("1")))
	fmt.Println("b:", v.MapIndex(reflect.ValueOf("b")))

	type A struct {
		Name string
		Age  int
	}
	type B struct {
		Info  *A
		Extra string
	}

	b := &B{
		Info: &A{
			Name: "name",
			Age:  1,
		},
		Extra: "extra",
	}

	vbp := reflect.ValueOf(b)
	vb := vbp.Elem()

	fmt.Println("nb:", vb.FieldByName("Info").Elem().FieldByName("Name"))

	bs, err := json.Marshal(vbp)
	fmt.Println("bs1:", string(bs))
	fmt.Println("err1:", err)

	bs, err = json.Marshal(vbp.Interface())
	fmt.Println("bs2:", string(bs))
	fmt.Println("er2:", err)

	t.Fail()
}

func Test_Report(t *testing.T) {
	time1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-11-29 11:11:28", time.Local)
	time2 := time1.Add(+10 * time.Minute)
	time3 := time1.Add(+2 * time.Hour)
	time4 := time1.Add(+4 * time.Hour)

	btime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-11-29 11:00:00", time.Local)
	etime, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-11-29 20:00:00", time.Local)
	/*
		1.同一小时内同用户
		2.同一小时内不同用户
	*/
	metas := []*ReportMeta{
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":4,"dy":1,"rl":2}`,
			UID:     1,
			Uname:   "1",
		},
		{Mtime: xtime.Time(time2.Unix()),
			Content: `{"ds":4,"rl":3}`,
			UID:     1,
			Uname:   "1",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":4}`,
			UID:     2,
			Uname:   "2",
		},
		{Mtime: xtime.Time(time2.Unix()),
			Content: `{"ds":0,"st_3_225":1,"ut":0}`,
			UID:     2,
			Uname:   "2",
		},
		{Mtime: xtime.Time(time3.Unix()),
			Content: `{"ds":13,"dy":5,"rl":6,"rs_-1":1,"st_3_736":1,"st_4_0":2,"st_4_736":1,"ut":0}`,
			UID:     1,
			Uname:   "1",
		},
		{Mtime: xtime.Time(time4.Unix()),
			Content: `{"ds":2,"rs_1":1,"st_3_2":1,"ut":9}`,
			UID:     2,
			Uname:   "2",
		},
	}
	opt := &OptReport{
		Btime: btime,
		Etime: etime,
	}

	mnames := map[int64]string{
		1: "1",
		2: "2",
	}
	tempres := Gentempres(opt, mnames, metas)
	bs1, _ := json.Marshal(tempres)
	fmt.Printf("1: %s\n", string(bs1))
	res := Genres(opt, tempres, mnames)
	bs2, _ := json.Marshal(res)
	fmt.Printf("2: %s\n", string(bs2))
	form := Genform(res)
	bs3, _ := json.Marshal(form)
	fmt.Printf("3: %s\n", string(bs3))
	t.Fail()
}

func Test_MemberStat(t *testing.T) {
	time1, _ := time.ParseInLocation("2006-01-02 15:04:05", "2018-11-30 11:11:28", time.Local)
	/*
		1.处理量
		2.处理率
		3.通过率
		4.平均耗时
	*/
	metas := []*ReportMeta{
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":10,"dy":5,"rl":5,"rs_-1":1,"st_3_1":2,"st_4_0":2,"st_4_736":1,"ut":10}`,
			UID:     1,
			Uname:   "1",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":5,"dy":5,"rl":5,"rs_-1":1,"rs_0":3,"st_3_1":2,"st_4_0":2,"st_4_736":1,"ut":10}`,
			UID:     1,
			Uname:   "1",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":5,"dy":0,"rl":0,"rs_-1":1,"rs_0":3,"st_3_1":1,"st_4_0":2,"st_4_736":1,"ut":10}`,
			UID:     1,
			Uname:   "1",
		},

		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":0,"st_3_0":1,"ut":0} `,
			UID:     2,
			Uname:   "2",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":0,"st_3_0":1,"ut":0}`,
			UID:     2,
			Uname:   "2",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":5,"rs_0":5,"st_3_2":5,"ut":21}`,
			UID:     2,
			Uname:   "2",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":2}`,
			UID:     3,
			Uname:   "3",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":0,"rl":2}`,
			UID:     3,
			Uname:   "3",
		},
		{Mtime: xtime.Time(time1.Unix()),
			Content: `{"ds":0,"st_4_4":4,"ut":0} `,
			UID:     4,
			Uname:   "4",
		},
	}

	res, _ := GenMemberStat(metas, 0)
	bs, _ := json.Marshal(res)
	fmt.Printf("%s\n", string(bs))
	t.Fail()
}

func Test_Column(t *testing.T) {
	column := `[
		{
			"name":"id",
			"chname":"呵呵"
		},{
			"name":"id2",
			"chname":"呵呵",
			"enum":{
				"1":"上",
				"2":"下"
			}
		}
	]`

	cs := []*Column{}
	e := json.Unmarshal([]byte(column), &cs)
	fmt.Println(e)
	fmt.Printf("cs:%+v\n", cs)
	t.Fail()
}

func TestLogFieldTemp(t *testing.T) {
	s := LogFieldTemp(LogFieldPID, 1, 0, true)
	convey.Convey("LogFieldTemp", t, func() {
		convey.So(s, convey.ShouldNotEqual, "")
	})
}

func TestGetEmptyInfo(t *testing.T) {
	s := GetEmptyInfo()
	convey.Convey("GetEmptyInfo", t, func() {
		convey.So(s, convey.ShouldNotBeNil)
	})
}

func TestEmptyListItem(t *testing.T) {
	s := EmptyListItem()
	convey.Convey("EmptyListItem", t, func() {
		convey.So(s, convey.ShouldNotBeNil)
	})
}
