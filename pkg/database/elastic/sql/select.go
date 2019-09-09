package sql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xwb1989/sqlparser"
)

// Condition for sql where.
type Condition struct {
	Field string
	Expr  string
	Value string
}

const (
	_defaultQuerySize = "10"
)

func handleSelect(sel *sqlparser.Select) (dsl string, esType string, where []*Condition, err error) {

	// Handle where
	// top level node pass in an empty interface
	// to tell the children this is root
	// is there any better way?
	var rootParent sqlparser.Expr
	var defaultQueryMapStr = `{"bool" : {"must": [{"match_all" : {}}]}}`
	var queryMapStr string

	// use may not pass where clauses
	if sel.Where != nil {
		queryMapStr, where, err = handleSelectWhere(&sel.Where.Expr, true, &rootParent)
		if err != nil {
			return "", "", nil, err
		}
	}
	if queryMapStr == "" {
		queryMapStr = defaultQueryMapStr
	}

	//TODO support multiple tables
	//for i, fromExpr := range sel.From {
	//	fmt.Printf("the %d of from is %#v\n", i, sqlparser.String(fromExpr))
	//}

	//Handle from
	// if len(sel.From) != 1 {
	// 	return "", "", errors.New("elasticsql: multiple from currently not supported")
	// }
	esType = sqlparser.String(sel.From)
	esType = strings.Replace(esType, "`", "", -1)

	queryFrom, querySize := "0", _defaultQuerySize

	aggFlag := false
	// if the request is to aggregation
	// then set aggFlag to true, and querySize to 0
	// to not return any query result

	var aggStr string
	if len(sel.GroupBy) > 0 || checkNeedAgg(sel.SelectExprs) {
		aggFlag = true
		querySize = "0"
		aggStr, err = buildAggs(sel)
		if err != nil {
			//aggStr = ""
			return "", "", where, err
		}
	}

	// Handle limit
	if sel.Limit != nil {
		if sel.Limit.Offset != nil {
			queryFrom = sqlparser.String(sel.Limit.Offset)
		}
		querySize = sqlparser.String(sel.Limit.Rowcount)
		if querySize == "0" && !aggFlag {
			querySize = _defaultQuerySize
		}
	}

	// Handle order by
	// when executating aggregations, order by is useless
	var orderByArr []string
	if aggFlag == false {
		for _, orderByExpr := range sel.OrderBy {
			orderField := strings.Replace(sqlparser.String(orderByExpr.Expr), "`", "", -1)
			nested := strings.Split(orderField, ".")
			var orderByStr string
			if len(nested) > 1 {
				orderByStr = fmt.Sprintf(`{"%v":{"mode":"max","order":"%v","nested_path":"%s"}}`, orderField, orderByExpr.Direction, nested[0])
			} else {
				orderByStr = fmt.Sprintf(`{"%v": "%v"}`, orderField, orderByExpr.Direction)
			}
			orderByArr = append(orderByArr, orderByStr)
		}
	}

	resultMap := make(map[string]interface{})
	resultMap["query"] = queryMapStr
	resultMap["from"] = queryFrom
	resultMap["size"] = querySize
	if len(aggStr) > 0 {
		resultMap["aggregations"] = aggStr
	}

	if len(orderByArr) > 0 {
		resultMap["sort"] = fmt.Sprintf("[%v]", strings.Join(orderByArr, ","))
	}

	// keep the travesal in order, avoid unpredicted json
	var keySlice = []string{"query", "from", "size", "sort", "aggregations"}
	var resultArr []string
	for _, mapKey := range keySlice {
		if val, ok := resultMap[mapKey]; ok {
			resultArr = append(resultArr, fmt.Sprintf(`"%v" : %v`, mapKey, val))
		}
	}

	dsl = "{" + strings.Join(resultArr, ",") + "}"
	return dsl, esType, where, nil
}

// if the where is empty, need to check whether to agg or not
func checkNeedAgg(sqlSelect sqlparser.SelectExprs) bool {
	for _, v := range sqlSelect {
		expr, ok := v.(*sqlparser.AliasedExpr)
		if !ok {
			// no need to handle, star expression * just skip is ok
			continue
		}

		//TODO more precise
		if _, ok := expr.Expr.(*sqlparser.FuncExpr); ok {
			return true
		}
	}
	return false
}

func buildNestedFuncStrValue(nestedFunc *sqlparser.FuncExpr) (string, error) {
	return "", errors.New("elasticsql: unsupported function" + nestedFunc.Name.String())
}

func handleSelectWhereAndExpr(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, []*Condition, error) {
	andExpr := (*expr).(*sqlparser.AndExpr)
	leftExpr := andExpr.Left
	rightExpr := andExpr.Right
	var cds []*Condition
	leftStr, cds, err := handleSelectWhere(&leftExpr, false, expr)
	if err != nil {
		return "", cds, err
	}
	rightStr, cdss, err := handleSelectWhere(&rightExpr, false, expr)
	if err != nil {
		return "", cdss, err
	}
	cds = append(cds, cdss...)
	// not toplevel
	// if the parent node is also and, then the result can be merged

	var resultStr string
	if leftStr == "" || rightStr == "" {
		resultStr = leftStr + rightStr
	} else {
		resultStr = leftStr + `,` + rightStr
	}

	if _, ok := (*parent).(*sqlparser.AndExpr); ok {
		return resultStr, cds, nil
	}
	return fmt.Sprintf(`{"bool" : {"must" : [%v]}}`, resultStr), cds, nil
}

func handleSelectWhereOrExpr(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, []*Condition, error) {
	orExpr := (*expr).(*sqlparser.OrExpr)
	leftExpr := orExpr.Left
	rightExpr := orExpr.Right
	var cds []*Condition
	leftStr, cds, err := handleSelectWhere(&leftExpr, false, expr)
	if err != nil {
		return "", cds, err
	}

	rightStr, cdss, err := handleSelectWhere(&rightExpr, false, expr)
	if err != nil {
		return "", cdss, err
	}
	cds = append(cds, cdss...)

	var resultStr string
	if leftStr == "" || rightStr == "" {
		resultStr = leftStr + rightStr
	} else {
		resultStr = leftStr + `,` + rightStr
	}

	// not toplevel
	// if the parent node is also or node, then merge the query param
	if _, ok := (*parent).(*sqlparser.OrExpr); ok {
		return resultStr, cds, nil
	}

	return fmt.Sprintf(`{"bool" : {"should" : [%v]}}`, resultStr), cds, nil
}

func buildComparisonExprRightStr(expr sqlparser.Expr) (string, bool, error) {
	var rightStr string
	var err error
	var missingCheck = false
	switch expr.(type) {
	case *sqlparser.SQLVal:
		rightStr = sqlparser.String(expr)
		rightStr = strings.Trim(rightStr, `'`)
	case *sqlparser.GroupConcatExpr:
		return "", missingCheck, errors.New("elasticsql: group_concat not supported")
	case *sqlparser.FuncExpr:
		// parse nested
		funcExpr := expr.(*sqlparser.FuncExpr)
		rightStr, err = buildNestedFuncStrValue(funcExpr)
		if err != nil {
			return "", missingCheck, err
		}
	case *sqlparser.ColName:
		if sqlparser.String(expr) == "missing" {
			missingCheck = true
			return "", missingCheck, nil
		}

		return "", missingCheck, errors.New("elasticsql: column name on the right side of compare operator is not supported")
	case sqlparser.ValTuple:
		rightStr = sqlparser.String(expr)
	default:
		// cannot reach here
	}
	return rightStr, missingCheck, err
}

func handleSelectWhereComparisonExpr(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, []*Condition, error) {
	comparisonExpr := (*expr).(*sqlparser.ComparisonExpr)
	colName, ok := comparisonExpr.Left.(*sqlparser.ColName)
	if !ok {
		return "", nil, nil // continue
		// return "", nil, errors.New("elasticsql: invalid comparison expression, the left must be a column name")
	}
	var cds []*Condition
	colNameStr := sqlparser.String(colName)
	colNameStr = strings.Replace(colNameStr, "`", "", -1)
	rightStr, missingCheck, err := buildComparisonExprRightStr(comparisonExpr.Right)
	if err != nil {
		return "", nil, err
	}
	condition := &Condition{
		Field: colNameStr,
		Expr:  comparisonExpr.Operator,
		Value: rightStr,
	}
	cds = append(cds, condition)
	resultStr := ""
	nested := strings.Split(colNameStr, ".")
	var (
		nestedFlag bool
		nestedPath string
	)
	if len(nested) == 2 {
		nestedFlag = true
		nestedPath = nested[0]
	}
	switch comparisonExpr.Operator {
	case ">=":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"from" : "%v"}}}`, colNameStr, rightStr)
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"nested":{"path":"%s","query":{"range":{"%v":{"from":"%v"}}}}}`, nestedPath, colNameStr, rightStr)
		}
	case "<=":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"to" : "%v"}}}`, colNameStr, rightStr)
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"nested":{"path":"%s","query":{"range":{"%v":{"to":"%v"}}}}}`, nestedPath, colNameStr, rightStr)
		}
	case "=":
		// field is missing
		if missingCheck {
			resultStr = fmt.Sprintf(`{"missing":{"field":"%v"}}`, colNameStr)
		} else {
			resultStr = fmt.Sprintf(`{"term" : {"%v" : "%v"}}`, colNameStr, rightStr)
			if nestedFlag {
				resultStr = fmt.Sprintf(`{"nested":{"path":"%s","query":{"term":{"%v":"%v"}}}}`, nestedPath, colNameStr, rightStr)
			}
		}
	case ">":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"gt" : "%v"}}}`, colNameStr, rightStr)
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"nested":{"path":"%s","query":{"range":{"%v":{"gt":"%v"}}}}}`, nestedPath, colNameStr, rightStr)
		}
	case "<":
		resultStr = fmt.Sprintf(`{"range" : {"%v" : {"lt" : "%v"}}}`, colNameStr, rightStr)
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"nested":{"path":"%s","query":{"range":{"%v":{"lt":"%v"}}}}}`, nestedPath, colNameStr, rightStr)
		}
	case "!=":
		if missingCheck {
			resultStr = fmt.Sprintf(`{"bool" : {"must_not" : [{"missing":{"field":"%v"}}]}}`, colNameStr)
		} else {
			resultStr = fmt.Sprintf(`{"bool" : {"must_not" : {"term" : {"%v" : "%v"}}}}`, colNameStr, rightStr)
			if nestedFlag {
				resultStr = fmt.Sprintf(`{"bool":{"must_not":{"nested":{"path":"%s","query":{"term":{"%v":"%v"}}}}}}`, nestedPath, colNameStr, rightStr)
			}
		}
	case "in":
		// the default valTuple is ('1', '2', '3') like
		// so need to drop the () and replace ' to "
		rightStr = strings.Replace(rightStr, `'`, `"`, -1)
		rightStr = strings.Trim(rightStr, "(")
		rightStr = strings.Trim(rightStr, ")")
		resultStr = fmt.Sprintf(`{"terms" : {"%v" : [%v]}}`, colNameStr, rightStr)
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"nested":{"path":"%s","query":{"terms":{"%v":[%v]}}}}`, nestedPath, colNameStr, rightStr)
		}
	case "like":
		rightStr = strings.Replace(rightStr, `%`, ``, -1)
		resultStr = fmt.Sprintf(`{"common" : {"%v" : {"query" : "%v", "minimum_should_match" : "100%%"}}}`, colNameStr, string([]rune(rightStr)))
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"nested":{"path":"%s","query":{"common":{"%v":{"query":"%v","minimum_should_match":"100%%"}}}}}`, nestedPath, colNameStr, rightStr)
		}
	case "not like":
		rightStr = strings.Replace(rightStr, `%`, ``, -1)
		resultStr = fmt.Sprintf(`{"bool" : {"must_not" : {"common" : {"%v" : {"query" : "%v",  "minimum_should_match" : "100%%"}}}}}`, colNameStr, string([]rune(rightStr)))
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"bool":{"must_not":{"nested":{"path":"%s","query":{"common":{"%v":{"query":"%v","minimum_should_match":"100%%"}}}}}}}`, nestedPath, colNameStr, rightStr)
		}
	case "not in":
		// the default valTuple is ('1', '2', '3') like
		// so need to drop the () and replace ' to "
		rightStr = strings.Replace(rightStr, `'`, `"`, -1)
		rightStr = strings.Trim(rightStr, "(")
		rightStr = strings.Trim(rightStr, ")")
		resultStr = fmt.Sprintf(`{"bool" : {"must_not" : {"terms" : {"%v" : [%v]}}}}`, colNameStr, rightStr)
		if nestedFlag {
			resultStr = fmt.Sprintf(`{"bool":{"must_not":{"nested":{"path":"%s","query":{"terms":{"%v":[%v]}}}}}}`, nestedPath, colNameStr, rightStr)
		}
	}

	// the root node need to have bool and must
	if topLevel {
		resultStr = fmt.Sprintf(`{"bool" : {"must" : [%v]}}`, resultStr)
	}

	return resultStr, cds, nil
}

func handleSelectWhere(expr *sqlparser.Expr, topLevel bool, parent *sqlparser.Expr) (string, []*Condition, error) {
	if expr == nil {
		return "", nil, errors.New("elasticsql: error expression cannot be nil here")
	}
	var cds []*Condition
	switch (*expr).(type) {
	case *sqlparser.AndExpr:
		return handleSelectWhereAndExpr(expr, topLevel, parent)

	case *sqlparser.OrExpr:
		return handleSelectWhereOrExpr(expr, topLevel, parent)
	case *sqlparser.ComparisonExpr:
		return handleSelectWhereComparisonExpr(expr, topLevel, parent)

	case *sqlparser.IsExpr:
		return "", cds, errors.New("elasticsql: is expression currently not supported")
	case *sqlparser.RangeCond:
		// between a and b
		// the meaning is equal to range query
		rangeCond := (*expr).(*sqlparser.RangeCond)
		colName, ok := rangeCond.Left.(*sqlparser.ColName)

		if !ok {
			return "", cds, errors.New("elasticsql: range column name missing")
		}

		colNameStr := sqlparser.String(colName)
		fromStr := strings.Trim(sqlparser.String(rangeCond.From), `'`)
		toStr := strings.Trim(sqlparser.String(rangeCond.To), `'`)

		resultStr := fmt.Sprintf(`{"range" : {"%v" : {"from" : "%v", "to" : "%v"}}}`, colNameStr, fromStr, toStr)
		if topLevel {
			resultStr = fmt.Sprintf(`{"bool" : {"must" : [%v]}}`, resultStr)
		}

		return resultStr, cds, nil

	case *sqlparser.ParenExpr:
		parentBoolExpr := (*expr).(*sqlparser.ParenExpr)
		boolExpr := parentBoolExpr.Expr

		// if paren is the top level, bool must is needed
		var isThisTopLevel = false
		if topLevel {
			isThisTopLevel = true
		}
		return handleSelectWhere(&boolExpr, isThisTopLevel, parent)
	case *sqlparser.NotExpr:
		return "", cds, errors.New("elasticsql: not expression currently not supported")
	}

	return "", cds, errors.New("elaticsql: logically cannot reached here")
}
