package expr

import (
	"fmt"
	"testing"
)

func TestExpr(t *testing.T) {
	tests := []struct {
		expr string
		env  Env
		want string
	}{
		{"$1 > 80 && $2 <9", Env{"$1": 100, "$2": 2}, "true"},
		{"$1 % 10", Env{"$1": 104, "$2": 2}, "4"},
		{"pow(x, false) + pow(y, false)", Env{"x": 12, "y": 1}, "2"},
		{"pow(x, 3) + pow(y, 3)", Env{"x": 9, "y": 10}, "1729"},
		{"5.0 / 9 * (F - 32)", Env{"F": -40}, "-40"},
		{"5.0 / 9 * (F - 32)", Env{"F": 32}, "0"},
		{"5.0 / 9 * (F - 32)", Env{"F": 212}, "100"},
		////!-Eval
		//// additional tests that don't appear in the book
		{"-1 + -x", Env{"x": 1}, "-2"},
		{"-1 - x", Env{"x": 1}, "-2"},
		{"a >= 10", Env{"a": 15}, "true"},
		{"b >= sin(10) && a < 1", Env{"a": 9, "b": 10}, "false"},
		{"!!!true", Env{"a": 9, "b": 10}, "false"},
		//!+Eval
	}
	var prevExpr string
	parser := NewExpressionParser()
	for _, test := range tests {
		// Print expr only when it changes.
		if test.expr != prevExpr {
			t.Logf("\n%s\n", test.expr)
			prevExpr = test.expr
		}
		if err := parser.Parse(test.expr); err != nil {
			t.Error(err) // parse error
			continue
		}
		got := fmt.Sprintf("%v", parser.GetExpr().Eval(test.env))
		t.Logf("\t%v => %s\n", test.env, got)
		if got != test.want {
			t.Errorf("%s.Eval() in %v = %q, want %q\n",
				test.expr, test.env, got, test.want)
		}
	}
}
