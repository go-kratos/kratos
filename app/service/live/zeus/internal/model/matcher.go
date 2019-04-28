package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-common/app/service/live/zeus/expr"
	"go-common/library/log"
)

type Matcher struct {
	Group             map[string][]*MatcherBucket `json:"group"`
	VariableWhitelist []string                    `json:"variable_whitelist"`
}

type MatcherBucket struct {
	RuleExpr string      `json:"expr"`
	Extend   interface{} `json:"extend"`
	expr     expr.Expr
	variable []string
	config   string
}

func NewMatcher(config string) (*Matcher, error) {
	matcher := &Matcher{}
	if err := json.Unmarshal([]byte(config), matcher); err != nil {
		return nil, err
	}
	whitelist := make(map[string]struct{}, len(matcher.VariableWhitelist))
	for _, v := range matcher.VariableWhitelist {
		whitelist[v] = struct{}{}
	}
	parser := expr.NewExpressionParser()
	for _, bucket := range matcher.Group {
		for _, b := range bucket {
			if e := parser.Parse(b.RuleExpr); e != nil {
				msg := fmt.Sprintf("zeus compile error, rule:%s error:%s", b.RuleExpr, e.Error())
				return nil, errors.New(msg)
			}
			b.expr = parser.GetExpr()
			b.variable = parser.GetVariable()
			for _, v := range b.variable {
				if _, ok := whitelist[v]; !ok {
					msg := fmt.Sprintf("zeus check error, rule:%s error: variable %s not in whitelist", b.RuleExpr, v)
					return nil, errors.New(msg)
				}
			}
			if config, e := json.Marshal(b.Extend); e != nil {
				msg := fmt.Sprintf("zeus parse config error, rule:%s error:%s", b.RuleExpr, e.Error())
				return nil, errors.New(msg)
			} else {
				b.config = string(config)
			}
		}
	}
	return matcher, nil
}

func (m *Matcher) Match(group string, env expr.Env) (bool, string, error) {
	var bucket []*MatcherBucket
	var ok bool
	if bucket, ok = m.Group[group]; !ok {
		return false, "", errors.New("group not found")
	}
	isMatch := false
	extend := ""
	for _, b := range bucket {
		result, err := expr.EvalBool(b.expr, env)
		if err != nil {
			log.Error("zeus eval error, group:%s expr:%s error:%s", group, b.RuleExpr, err.Error())
			continue
		}
		if result {
			isMatch = true
			extend = b.config
			break
		}
	}
	return isMatch, extend, nil
}
