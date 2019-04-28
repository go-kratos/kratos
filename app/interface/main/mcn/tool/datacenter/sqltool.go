package datacenter

import (
	"encoding/json"
	"fmt"
	"go-common/library/log"
	"strings"
	"text/scanner"
)

// operator
const (
	opIn   = "in"
	opNin  = "nin"
	opLike = "like"
	opLte  = "lte" // <=
	opLt   = "lt"  // <
	opGte  = "gte" // >=
	opGt   = "gt"  // >
	opNull = "null"
)

// value for Null operator
const (
	IsNull    = 1
	IsNotNull = -1
)

// value for sort
const (
	Desc = -1
	Asc  = 1
)

//ConditionMapType condition map
type ConditionMapType map[string]map[string]interface{}

//ConditionType condition's in map
type ConditionType map[string]interface{}

//ConditionIn in
func ConditionIn(v ...interface{}) ConditionType {
	return ConditionType{
		opIn: v,
	}
}

func conditionHelper(k string, v interface{}) ConditionType {
	return ConditionType{
		k: v,
	}
}

//ConditionLte <=
func ConditionLte(v interface{}) ConditionType {
	return conditionHelper(opLte, v)
}

//ConditionLt <
func ConditionLt(v interface{}) ConditionType {
	return conditionHelper(opLt, v)
}

//ConditionGte >=
func ConditionGte(v interface{}) ConditionType {
	return conditionHelper(opGte, v)
}

//ConditionGt >
func ConditionGt(v interface{}) ConditionType {
	return conditionHelper(opGt, v)
}

//SortType sort
type SortType map[string]int

//Query query
type Query struct {
	selection []map[string]string
	// <field, <operator, value> >
	where map[string]map[string]interface{}
	sort  map[string]int
	limit map[string]int
	err   error
}

const (
	keyField = "name"
	keyAs    = "as"
)

func makeField(field string) map[string]string {
	return map[string]string{keyField: field}
}

func makeFieldAs(field, as string) map[string]string {
	return map[string]string{keyField: field, keyAs: as}
}

//Select select fields, use similar as sql
func (q *Query) Select(fields string) *Query {
	if q.err != nil {
		return q
	}
	var fieldsAll = strings.Split(fields, ",")
	for _, v := range fieldsAll {
		var s scanner.Scanner
		s.Init(strings.NewReader(v))
		var tokens []string
		for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
			txt := s.TokenText()
			tokens = append(tokens, txt)
		}
		switch len(tokens) {
		case 1:
			if tokens[0] == "*" {
				q.selection = []map[string]string{}
				return q
			}
			q.selection = append(q.selection, makeField(tokens[0]))
		case 2:
			q.selection = append(q.selection, makeFieldAs(tokens[0], tokens[1]))
		case 3:
			q.selection = append(q.selection, makeFieldAs(tokens[0], tokens[2]))
		}
	}
	return q
}

//Where where condition, see test for examples
func (q *Query) Where(conditions ...ConditionMapType) *Query {
	if q.err != nil {
		return q
	}
	if q.where == nil {
		q.where = make(ConditionMapType, len(conditions))
	}
	for _, mapData := range conditions {
		for k1, v1 := range mapData {
			if q.where[k1] == nil {
				q.where[k1] = make(map[string]interface{})
			}
			// combine all pair of map[string]interface{}(v1) into q.where[k1]
			for k2, v2 := range v1 {
				q.where[k1][k2] = v2
			}
		}
	}
	return q
}

//Order order field, use similar as sql
func (q *Query) Order(sort string) *Query {
	if q.err != nil {
		return q
	}
	var fields = strings.Split(sort, ",")
	if q.sort == nil {
		q.sort = make(map[string]int, len(fields))
	}
	for _, v := range fields {
		var s scanner.Scanner
		s.Init(strings.NewReader(v))
		var tokens []string
		for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
			txt := s.TokenText()
			tokens = append(tokens, txt)
		}
		switch len(tokens) {
		case 1:
			q.sort[tokens[0]] = Asc
		case 2:
			var order = Asc
			switch strings.ToLower(tokens[1]) {
			case "asc":
				order = Asc
			case "desc":
				order = Desc
			}
			q.sort[tokens[0]] = order
		default:
			q.err = fmt.Errorf("parse order fail, [%s]", sort)
			log.Error("%s", q.err)
			return q
		}
	}
	return q
}

//Limit limit, same as sql
func (q *Query) Limit(limit, offset int) *Query {
	if q.err != nil {
		return q
	}
	if q.limit == nil {
		q.limit = make(map[string]int, 2)
	}

	q.limit["limit"] = limit
	q.limit["skip"] = offset
	return q
}

//String to string
func (q *Query) String() (res string) {
	if q.err != nil {
		return q.err.Error()
	}
	var resultMap = map[string]interface{}{}

	if q.selection != nil {
		resultMap["select"] = q.selection
	}
	if q.where != nil {
		resultMap["where"] = q.where
	}

	if q.sort != nil {
		resultMap["sort"] = q.sort
	}

	if q.limit != nil {
		resultMap["page"] = q.limit
	}
	resBytes, _ := json.Marshal(resultMap)
	res = string(resBytes)
	return
}

//Error return error if get error
func (q *Query) Error() error {
	return q.err
}
