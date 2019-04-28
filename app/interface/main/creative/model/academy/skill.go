package academy

import (
	"reflect"

	"go-common/app/service/main/archive/api"
	"go-common/library/time"
)

//Occupation for occupation.
type Occupation struct {
	ID           int64  `json:"id"`
	Rank         int64  `json:"rank"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	MainStep     string `json:"main_step"`
	MainSoftWare string `json:"main_software"`
	Logo         string `json:"logo"`
}

//Skill for Skill.
type Skill struct {
	ID   int64  `json:"id"`
	OID  int64  `json:"oid"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

//SkillArc for Skill.
type SkillArc struct {
	ID   int64 `json:"id"`
	AID  int64 `json:"aid"`
	Type int   `json:"type"`
	PID  int64 `json:"pid"`
	SkID int64 `json:"skid"`
	SID  int64 `json:"sid"`
}

//ArcMeta for skill arc meta.
type ArcMeta struct {
	AID      int64             `json:"aid"`
	MID      int64             `json:"mid"`
	Cover    string            `json:"cover"`
	Title    string            `json:"title"`
	Type     string            `json:"type"`
	Duration int64             `json:"duration,omitempty"`
	PlayTime time.Time         `json:"play_time"` //历史课程上次学习时间
	Watch    int8              `json:"watch"`     //标记是否观看过
	ArcStat  *api.Stat         `json:"arc_stat,omitempty"`
	Skill    *SkillArc         `json:"-"`
	Business int8              `json:"business"`
	Tags     map[string][]*Tag `json:"tags,omitempty"`
}

//ArcList for archive list.
type ArcList struct {
	Items []*ArcMeta   `json:"items"`
	Page  *ArchivePage `json:"page"`
}

//NewbCourseList for NewbCourse list.
type NewbCourseList struct {
	Items []*ArcMeta `json:"items"`
	Title string     `json:"title"`
	TID   int64      `json:"tid"`
}

//Play for academy play list.
type Play struct {
	MID      int64     `json:"mid"`
	AID      int64     `json:"aid"`
	Business int8      `json:"business"`
	Watch    int8      `json:"watch"`
	CTime    time.Time `json:"ctime"`
	MTime    time.Time `json:"mtime"`
}

//SearchKeywords for academy h5 search keywords.
type SearchKeywords struct {
	ID       int64             `json:"id"`
	Rank     int64             `json:"rank"`
	ParentID int64             `json:"parent_id"`
	State    int8              `json:"state"`
	Name     string            `json:"name"`
	Comment  string            `json:"comment"`
	CTime    string            `json:"-"`
	MTime    string            `json:"-"`
	Count    int               `json:"count,omitempty"`
	Children []*SearchKeywords `json:"children,omitempty"`
}

//Trees for generate tree data set
// data - db result set
// idFieldStr - primary key in table map to struct
// pidFieldStr - top parent id in table map to struct
// chFieldStr - struct child nodes
func Trees(data interface{}, idFieldStr, pidFieldStr, chFieldStr string) (res []interface{}) {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return
	}

	sli := reflect.ValueOf(data)
	top := make(map[int64]interface{})
	res = make([]interface{}, 0, sli.Len())
	for i := 0; i < sli.Len(); i++ {
		v := sli.Index(i).Interface()
		if reflect.TypeOf(v).Kind() != reflect.Ptr {
			continue
		}

		if reflect.ValueOf(v).IsNil() {
			continue
		}

		getValue := reflect.ValueOf(v).Elem()
		getType := reflect.TypeOf(v).Elem()
		pid := getValue.FieldByName(pidFieldStr).Interface().(int64)
		if _, ok := getType.FieldByName(pidFieldStr); ok && pid == 0 {
			id := getValue.FieldByName(idFieldStr).Interface().(int64)
			top[id] = v
			res = append(res, v)
		}
	}

	for i := 0; i < sli.Len(); i++ {
		v := sli.Index(i).Interface()
		if reflect.TypeOf(v).Kind() != reflect.Ptr {
			continue
		}

		if reflect.ValueOf(v).IsNil() {
			continue
		}

		pid := reflect.ValueOf(v).Elem().FieldByName(pidFieldStr).Interface().(int64)
		if pid == 0 {
			continue
		}

		if p, ok := top[pid]; ok {
			ch := reflect.ValueOf(p).Elem().FieldByName(chFieldStr)
			ch.Set(reflect.Append(ch, reflect.ValueOf(v)))
		}
	}
	return
}
