package expr

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"
)

const (
	TokenEOF = -(iota + 1)
	TokenIdent
	TokenInt
	TokenFloat
	TokenOperator
)

type lexer struct {
	scan  scanner.Scanner
	token rune
	text  string
}

func (lex *lexer) getToken() rune {
	return lex.token
}

func (lex *lexer) getText() string {
	return lex.text
}

func (lex *lexer) next() {
	token := lex.scan.Scan()
	text := lex.scan.TokenText()
	switch token {
	case scanner.EOF:
		lex.token = TokenEOF
		lex.text = text
	case scanner.Ident:
		lex.token = TokenIdent
		lex.text = text
	case scanner.Int:
		lex.token = TokenInt
		lex.text = text
	case scanner.Float:
		lex.token = TokenFloat
		lex.text = text
	case '+', '-', '*', '/', '%', '~':
		lex.token = TokenOperator
		lex.text = text
	case '&', '|', '=':
		var buffer bytes.Buffer
		lex.token = TokenOperator
		buffer.WriteRune(token)
		next := lex.scan.Peek()
		if next == token {
			buffer.WriteRune(next)
			lex.scan.Scan()
		}
		lex.text = buffer.String()
	case '>', '<', '!':
		var buffer bytes.Buffer
		lex.token = TokenOperator
		buffer.WriteRune(token)
		next := lex.scan.Peek()
		if next == '=' {
			buffer.WriteRune(next)
			lex.scan.Scan()
		}
		lex.text = buffer.String()
	default:
		if token >= 0 {
			lex.token = token
			lex.text = text
		} else {
			msg := fmt.Sprintf("got unknown token:%q, text:%s", lex.token, lex.text)
			panic(lexPanic(msg))
		}
	}
	//fmt.Printf("token:%d, text:%s\n", lex.token, lex.text)
}

type lexPanic string

// describe returns a string describing the current token, for use in errors.
func (lex *lexer) describe() string {
	switch lex.token {
	case TokenEOF:
		return "end of file"
	case TokenIdent:
		return fmt.Sprintf("identifier %s", lex.getText())
	case TokenInt, TokenFloat:
		return fmt.Sprintf("number %s", lex.getText())
	}
	return fmt.Sprintf("%q", rune(lex.getToken())) // any other rune
}

func precedence(token rune, text string) int {
	if token == TokenOperator {
		switch text {
		case "~", "!":
			return 9
		case "*", "/", "%":
			return 8
		case "+", "-":
			return 7
		case ">", ">=", "<", "<=":
			return 6
		case "!=", "==", "=":
			return 5
		case "&":
			return 4
		case "|":
			return 3
		case "&&":
			return 2
		case "||":
			return 1
		default:
			msg := fmt.Sprintf("unknown operator:%s", text)
			panic(lexPanic(msg))
		}
	}
	return 0
}

// ---- parser ----
type ExpressionParser struct {
	expression Expr
	variable   map[string]struct{}
}

func NewExpressionParser() *ExpressionParser {
	return &ExpressionParser{
		expression: nil,
		variable:   make(map[string]struct{}),
	}
}

// Parse parses the input string as an arithmetic expression.
//
//   expr = num                         a literal number, e.g., 3.14159
//        | id                          a variable name, e.g., x
//        | id '(' expr ',' ... ')'     a function call
//        | '-' expr                    a unary operator ( + - ! )
//        | expr '+' expr               a binary operator ( + - * / && & || | == )
//
func (parser *ExpressionParser) Parse(input string) (err error) {
	defer func() {
		switch x := recover().(type) {
		case nil:
			// no panic
		case lexPanic:
			err = fmt.Errorf("%s", x)
		default:
			// unexpected panic: resume state of panic.
			panic(x)
		}
	}()
	lex := new(lexer)
	lex.scan.Init(strings.NewReader(input))
	lex.scan.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats
	lex.scan.IsIdentRune = parser.isIdentRune
	lex.next() // initial lookahead
	parser.expression = nil
	parser.variable = make(map[string]struct{})
	e := parser.parseExpr(lex)
	if lex.token != scanner.EOF {
		return fmt.Errorf("unexpected %s", lex.describe())
	}
	parser.expression = e
	return nil
}

func (parser *ExpressionParser) GetExpr() Expr {
	return parser.expression
}

func (parser *ExpressionParser) GetVariable() []string {
	variable := make([]string, 0, len(parser.variable))
	for v := range parser.variable {
		if v != "true" && v != "false" {
			variable = append(variable, v)
		}
	}
	return variable
}

func (parser *ExpressionParser) isIdentRune(ch rune, i int) bool {
	return ch == '$' || ch == '_' || unicode.IsLetter(ch) || unicode.IsDigit(ch) && i > 0
}

func (parser *ExpressionParser) parseExpr(lex *lexer) Expr {
	return parser.parseBinary(lex, 1)
}

// binary = unary ('+' binary)*
// parseBinary stops when it encounters an
// operator of lower precedence than prec1.
func (parser *ExpressionParser) parseBinary(lex *lexer, prec1 int) Expr {
	lhs := parser.parseUnary(lex)
	for prec := precedence(lex.getToken(), lex.getText()); prec >= prec1; prec-- {
		for precedence(lex.getToken(), lex.getText()) == prec {
			op := lex.getText()
			lex.next() // consume operator
			rhs := parser.parseBinary(lex, prec+1)
			lhs = binary{op, lhs, rhs}
		}
	}
	return lhs
}

// unary = '+' expr | primary
func (parser *ExpressionParser) parseUnary(lex *lexer) Expr {
	if lex.getToken() == TokenOperator {
		op := lex.getText()
		if op == "+" || op == "-" || op == "~" || op == "!" {
			lex.next()
			return unary{op, parser.parseUnary(lex)}
		} else {
			msg := fmt.Sprintf("unary got unknown operator:%s", lex.getText())
			panic(lexPanic(msg))
		}
	}
	return parser.parsePrimary(lex)
}

// primary = id
//         | id '(' expr ',' ... ',' expr ')'
//         | num
//         | '(' expr ')'
func (parser *ExpressionParser) parsePrimary(lex *lexer) Expr {
	switch lex.token {
	case TokenIdent:
		id := lex.getText()
		lex.next()
		if lex.token != '(' {
			parser.variable[id] = struct{}{}
			return Var(id)
		}
		lex.next() // consume '('
		var args []Expr
		if lex.token != ')' {
			for {
				args = append(args, parser.parseExpr(lex))
				if lex.token != ',' {
					break
				}
				lex.next() // consume ','
			}
			if lex.token != ')' {
				msg := fmt.Sprintf("got %q, want ')'", lex.token)
				panic(lexPanic(msg))
			}
		}
		lex.next() // consume ')'
		return call{id, args}

	case TokenFloat:
		f, err := strconv.ParseFloat(lex.getText(), 64)
		if err != nil {
			panic(lexPanic(err.Error()))
		}
		lex.next() // consume number
		return literal{value: f}

	case TokenInt:
		i, err := strconv.ParseInt(lex.getText(), 10, 64)
		if err != nil {
			panic(lexPanic(err.Error()))
		}
		lex.next() // consume number
		return literal{value: i}

	case '(':
		lex.next() // consume '('
		e := parser.parseExpr(lex)
		if lex.token != ')' {
			msg := fmt.Sprintf("got %s, want ')'", lex.describe())
			panic(lexPanic(msg))
		}
		lex.next() // consume ')'
		return e
	}
	msg := fmt.Sprintf("unexpected %s", lex.describe())
	panic(lexPanic(msg))
}
