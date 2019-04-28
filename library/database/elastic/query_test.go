package elastic

import (
	"context"
	"testing"
	"time"

	"go-common/library/queue/databus/report"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_UserActionLog(t *testing.T) {
	Convey("elastic user action log", t, func() {
		type page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Total int `json:"total"`
		}
		var res struct {
			Page   *page                   `json:"page"`
			Result []*report.UserActionLog `json:"result"`
		}
		r := NewElastic(nil).NewRequest("log_user_action").Index("log_user_action_21_2018_06_1623")
		err := r.Scan(context.TODO(), &res)
		So(err, ShouldBeNil)
		t.Logf("query page(%+v)", res.Page)
		t.Logf("query result(%v)", len(res.Result))
	})
}

func Test_Scan(t *testing.T) {
	Convey("scan", t, func() {
		type page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Total int `json:"total"`
		}
		type ticketData struct {
			ID       int    `json:"id"`
			Name     string `json:"name"`
			City     int    `json:"city"`
			Province int    `json:"province"`
		}
		var res struct {
			Page   *page         `json:"page"`
			Result []*ticketData `json:"result"`
		}
		r := NewElastic(nil).NewRequest("ticket_venue").Fields("id", "city").Index("ticket_venue").WhereEq("city", 310100).Order("ctime", OrderDesc).Order("id", OrderAsc).
			WhereOr("province", 310000).WhereRange("id", 1, 2000000, RangeScopeLcRc).WhereLike([]string{"name"}, []string{"梅赛德斯奔驰文化中心"}, true, LikeLevelHigh)
		err := r.Scan(context.TODO(), &res)
		So(err, ShouldBeNil)
		t.Logf("query page(%+v)", res.Page)
		t.Logf("query result(%v)", len(res.Result))
	})
}

func Test_Index(t *testing.T) {
	Convey("example index by mod", t, func() {
		oid := int64(8888888)
		r := NewElastic(nil).NewRequest("LOG_USER_COIN")
		r.IndexByMod("log_user_action", oid, 100)
		q, _ := r.q.string()
		So(q, ShouldContainSubstring, `"from":"log_user_action_88"`)
	})
	Convey("example index by time", t, func() {
		Convey("indexByTime by week", func() {
			// 按周划分索引的方式
			r := NewElastic(nil).NewRequest("LOG_USER_COIN")
			r.IndexByTime("log_user_action", IndexTypeWeek, time.Date(2018, 8, 25, 2, 0, 0, 0, time.Local), time.Date(2018, 9, 24, 2, 0, 0, 0, time.Local))
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_08_2431")
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_09_0107")
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_09_0815")
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_09_1623")
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_09_2431")
		})
		Convey("indexByTime by year", func() {
			// 按年划分索引的方式
			r := NewElastic(nil).NewRequest("LOG_USER_COIN")
			r.IndexByTime("log_user_action", IndexTypeYear, time.Date(2017, 8, 25, 2, 0, 0, 0, time.Local), time.Date(2018, 9, 24, 2, 0, 0, 0, time.Local))
			So(r.q.From, ShouldContainSubstring, "log_user_action_2017")
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018")
		})
		Convey("indexByTime by month", func() {
			// 按月划分索引的方式
			r := NewElastic(nil).NewRequest("LOG_USER_COIN")
			r.IndexByTime("log_user_action", IndexTypeMonth, time.Date(2018, 8, 25, 2, 0, 0, 0, time.Local), time.Date(2018, 9, 24, 2, 0, 0, 0, time.Local))
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_08")
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_09")
		})
		Convey("indexByTime by day", func() {
			// 按天划分索引的方式
			r := NewElastic(nil).NewRequest("LOG_USER_COIN")
			r.IndexByTime("log_user_action", IndexTypeDay, time.Date(2018, 9, 23, 4, 0, 0, 0, time.Local), time.Date(2018, 9, 24, 2, 0, 0, 0, time.Local))
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_09_23")
			So(r.q.From, ShouldContainSubstring, "log_user_action_2018_09_24")
		})
	})
}

func Test_Combo_Simple(t *testing.T) {
	Convey("combo", t, func() {
		e := NewElastic(nil)
		// 实现: (tid in (2,3,4,5,9,20,21,22)) && (tid in (35,38)) && (tid in (15,17,18))
		cmbA := &Combo{}
		cmbA.ComboIn([]map[string][]interface{}{
			{"tid": {2, 3, 4, 5, 9, 20, 21, 22}},
		}).MinIn(1).MinAll(1)

		cmbB := &Combo{}
		cmbB.ComboIn([]map[string][]interface{}{
			{"tid": {35, 38}},
		}).MinIn(1).MinAll(1)

		cmbC := &Combo{}
		cmbC.ComboIn([]map[string][]interface{}{
			{"tid": {15, 17, 18}},
		}).MinIn(1).MinAll(1)

		r := e.NewRequest("").WhereCombo(cmbA, cmbB, cmbC)
		q, _ := r.q.string()
		So(q, ShouldContainSubstring, `"where":{"combo":[{"eq":null,"in":[{"tid":[2,3,4,5,9,20,21,22]}],"range":null,"min":{"eq":0,"in":1,"range":0,"min":1}},{"eq":null,"in":[{"tid":[35,38]}],"range":null,"min":{"eq":0,"in":1,"range":0,"min":1}},{"eq":null,"in":[{"tid":[15,17,18]}],"range":null,"min":{"eq":0,"in":1,"range":0,"min":1}}]}`)
	})
}

func Test_Combo_Complicated(t *testing.T) {
	Convey("combo not", t, func() {
		e := NewElastic(nil)
		cmb := &Combo{}
		cmb.ComboIn([]map[string][]interface{}{
			{"tid": {1, 2}},
			{"tid_type": {2, 3}},
		}).ComboRange([]map[string]string{
			{"id": "(10,20)"},
		}).ComboNotEQ([]map[string]interface{}{
			{"aid": 122},
			{"id": 677},
		}).MinIn(1).MinRange(1).MinNotEQ(1).MinAll(1)
		r := e.NewRequest("").WhereCombo(cmb)
		q, _ := r.q.string()
		So(q, ShouldContainSubstring, `"where":{"combo":[{"in":[{"tid":[1,2]},{"tid_type":[2,3]}],"range":[{"id":"(10,20)"}],"not_eq":[{"aid":122},{"id":677}],"min":{"in":1,"range":1,"not_eq":1,"min":1}}]}}`)
	})
	Convey("combo", t, func() {
		e := NewElastic(nil)
		// 实现:
		// (aid=122 or id=677) && (tid in (1,2,3,21) or tid_type in (1,2,3)) && (id > 10) &&
		// (aid=88 or fid=99) && (mid in (11,33) or id in (22,33)) && (2 < cid <= 10  || sid > 10)
		cmbA := &Combo{}
		cmbA.ComboEQ([]map[string]interface{}{
			{"aid": 122},
			{"id": 677},
		}).ComboIn([]map[string][]interface{}{
			{"tid": {1, 2, 3, 21}},
			{"tid_type": {1, 2, 3}},
		}).ComboRange([]map[string]string{
			{"id": "(10,)"},
		}).MinIn(1).MinEQ(1).MinRange(1).MinAll(1)

		cmbB := &Combo{}
		cmbB.ComboEQ([]map[string]interface{}{
			{"aid": 88},
			{"fid": 99},
		}).ComboIn([]map[string][]interface{}{
			{"mid": {11, 33}},
			{"id": {22, 33}},
		}).ComboRange([]map[string]string{
			{"cid": "(2,4]"},
			{"sid": "(10,)"},
		}).MinEQ(1).MinIn(2).MinRange(1).MinAll(1)

		r := e.NewRequest("").WhereCombo(cmbA, cmbB)
		q, _ := r.q.string()
		So(q, ShouldContainSubstring, `"where":{"combo":[{"eq":[{"aid":122},{"id":677}],"in":[{"tid":[1,2,3,21]},{"tid_type":[1,2,3]}],"range":[{"id":"(10,)"}],"min":{"eq":1,"in":1,"range":1,"min":1}},{"eq":[{"aid":88},{"fid":99}],"in":[{"mid":[11,33]},{"id":[22,33]}],"range":[{"cid":"(2,4]"},{"sid":"(10,)"}],"min":{"eq":1,"in":2,"range":1,"min":1}}]}`)
	})
}

func Test_Query(t *testing.T) {
	Convey("query conditions", t, func() {
		e := NewElastic(nil)

		r := e.NewRequest("").WhereEq("city", 310100)
		q, _ := r.q.string()
		So(q, ShouldContainSubstring, `"where":{"eq":{"city":310100}}`)

		r = e.NewRequest("").WhereLike([]string{"a", "b"}, []string{"c", "d"}, true, LikeLevelHigh).WhereLike([]string{"e", "f"}, []string{"g", "h"}, false, LikeLevelMiddle)
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"like":[{"kw_fields":["a","b"],"kw":["c","d"],"or":true,"level":"high"},{"kw_fields":["e","f"],"kw":["g","h"],"or":false,"level":"middle"}]}`)

		r = e.NewRequest("").WhereNot(NotTypeEq, "province")
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"not":{"eq":{"province":true}}}`)

		r = e.NewRequest("").WhereRange("id", 100, 500, RangeScopeLcRo).WhereRange("date", "2018-08-08 08:08:08", "2019-09-09 09:09:09", RangeScopeLoRo)
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"range":{"date":"(2018-08-08 08:08:08,2019-09-09 09:09:09)","id":"[100,500)"}}`)

		r = e.NewRequest("").WhereOr("city", 100000)
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"or":{"city":100000}}`)

		ids := []int64{100, 200, 300}
		r = e.NewRequest("").WhereIn("city", ids)
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"in":{"city":[100,200,300]}}`)

		strs := []string{"a"}
		r = e.NewRequest("").WhereIn("name", strs)
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"in":{"name":["a"]}}`)

		field := "a"
		order := []map[string]string{{"a": "asc"}}
		r = e.NewRequest("").GroupBy(EnhancedModeGroupBy, field, order)
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"enhanced":[{"mode":"group_by","field":"a","order":[{"a":"asc"}]}]}`)

		field = "a"
		r = e.NewRequest("").Sum(field)
		q, _ = r.q.string()
		So(q, ShouldContainSubstring, `"where":{"enhanced":[{"mode":"sum","field":"a"}]}`)
	})
}
