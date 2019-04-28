// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

package expr

import (
	"fmt"
)

//!+Check

func (v Var) Check(vars map[Var]interface{}) error {
	vars[v] = true
	return nil
}

func (literal) Check(vars map[Var]interface{}) error {
	return nil
}

func (u unary) Check(vars map[Var]interface{}) error {
	return u.x.Check(vars)
}

func (b binary) Check(vars map[Var]interface{}) error {
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

func (c call) Check(vars map[Var]interface{}) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("unknown function %q", c.fn)
	}
	if len(c.args) != arity {
		return fmt.Errorf("call to %s has %d args, want %d",
			c.fn, len(c.args), arity)
	}
	for _, arg := range c.args {
		if err := arg.Check(vars); err != nil {
			return err
		}
	}
	return nil
}

var numParams = map[string]int{"pow": 2, "sin": 1, "sqrt": 1}

//!-Check
