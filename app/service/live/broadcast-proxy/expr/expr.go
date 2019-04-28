package expr

import "reflect"

// A Var identifies a variable, e.g., x.
type Var string

// A literal is a numeric constant, e.g., 3.141.
type literal struct {
	value interface{}
}

// An Expr is an arithmetic expression.
type Expr interface {
	// Eval returns the value of this Expr in the environment env.
	Eval(env Env) reflect.Value
	// Check reports errors in this Expr and adds its Vars to the set.
	Check(vars map[Var]interface{}) error
}

// A unary represents a unary operator expression, e.g., -x.
type unary struct {
	op string // one of '+', '-', '!', '~'
	x  Expr
}

// A binary represents a binary operator expression, e.g., x+y.
type binary struct {
	op   string
	x, y Expr
}

// A call represents a function call expression, e.g., sin(x).
type call struct {
	fn   string // one of "pow", "sin", "sqrt"
	args []Expr
}
