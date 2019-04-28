package dao

import "fmt"

const (
	//logicNoOp  = 0
	logicOpAnd = 1
	logicOpOr  = 2
)

const (
	//AppID appid
	AppID = "main.archive.up-service"
	//UatToken uat token
	UatToken = "5f1660060bb011e8865c66d44b23cda7"
	//TreeID uat treeid
	TreeID = "15572"
)

// Condition 生成的condition为
// [before] [key] 	[operator] 	? 	[after]
// 如
// 	[(]		[id] 	[=] 		? 	[and]
// 			[ctime] 	[>] 	? 	[)]
// 			[order by] [time] 	nil []  //如果value是nil，不会设置placeholder ?
type Condition struct {
	logicOp  int
	Before   string
	Key      string
	Operator string
	Value    interface{}
	After    string
}

//ConcatCondition concat conditions
func ConcatCondition(conditions ...Condition) (conditionStr string, args []interface{}, hasOperator bool) {
	hasOperator = false
	for _, c := range conditions {
		var questionMark = "?"
		if c.Value == nil {
			questionMark = ""
		}
		if c.Operator != "" {
			hasOperator = true
		}
		var logicOp = ""
		switch c.logicOp {
		case logicOpAnd:
			logicOp = " and "
		case logicOpOr:
			logicOp = " or "
		}
		conditionStr += fmt.Sprintf(" %s %s %s %s %s %s", logicOp, c.Before, c.Key, c.Operator, questionMark, c.After)
		if c.Value != nil {
			args = append(args, c.Value)
		}
	}
	return
}

//AndCondition and condition
func AndCondition(conditions ...Condition) (result []Condition) {
	return addLogicOperator(logicOpAnd, conditions...)
}

//OrCondition or condition
func OrCondition(conditions ...Condition) (result []Condition) {
	return addLogicOperator(logicOpOr, conditions...)
}

func addLogicOperator(operator int, conditions ...Condition) (result []Condition) {
	var isFirst = true
	for _, v := range conditions {
		if isFirst {
			isFirst = false
		} else {
			v.logicOp = operator
		}
		result = append(result, v)
	}
	return
}

// Split split num by size; size should be positive
func Split(start int, end int, size int, f func(start int, end int)) {
	for s := start; s < end; s += size {
		e := s + size
		if e > end {
			e = end
		}
		f(s, e)
	}
}
