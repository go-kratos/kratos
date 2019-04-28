package xregex

/*
golang version regex parser
refer to: https://github.com/aristotle9/as3cc/tree/master/java-template/src/org/lala/lex/utils/parser
*/

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	_initial   = "INITIAL"
	_deadState = 0xFFFFFFFF
	_maxValue  = 0x7fffffffffffffff
)

var (
	errEOF = errors.New("已经到达末尾")
)

// Lexer golang lexter
type lexer struct {
	transTable   []*stateTransItem
	finalTable   map[int64]int64
	initialTable map[string]int64
	inputTable   []*rangeItem
	start        int64
	oldStart     int64
	tokenName    string
	yyText       interface{}
	yy           interface{}
	ended        bool
	initialInput int64
	initialState string
	line         int64
	column       int64
	advanced     bool
	source       string
}

func newLexer() (lx *lexer) {
	lx = &lexer{}
	lx.transTable = []*stateTransItem{
		{false, []int64{0xFFFFFFFF, 0x3, 0x2, 0x1},
			[]*rangeItem{{0, 32, 0}, {33, 33, 1},
				{34, 34, 2}, {35, 35, 3}}},
		{false,
			[]int64{0xFFFFFFFF, 0xF, 0xE, 0xD, 0xC, 0xB, 0xA, 0x9, 0x8, 0x7, 0x6, 0x5,
				0x4},
			[]*rangeItem{{0, 0, 0}, {1, 1, 1},
				{2, 2, 2}, {3, 3, 3}, {4, 4, 4},
				{5, 5, 5}, {6, 6, 6}, {7, 7, 7},
				{8, 28, 8}, {29, 29, 9}, {30, 30, 10},
				{31, 31, 11}, {32, 32, 12},
				{33, 35, 0}}},
		{false, []int64{0xFFFFFFFF, 0xF, 0xE, 0xD, 0x8, 0x12, 0x11, 0x10},
			[]*rangeItem{{0, 0, 0}, {1, 1, 1},
				{2, 2, 2}, {3, 3, 3}, {4, 7, 4},
				{8, 8, 5}, {9, 9, 6}, {10, 27, 4},
				{28, 28, 7}, {29, 32, 4},
				{33, 35, 0}}},
		{false, []int64{0xFFFFFFFF, 0x16, 0x15, 0x14, 0x13},
			[]*rangeItem{{0, 21, 0}, {22, 24, 1},
				{25, 25, 2}, {26, 26, 3}, {27, 27, 4},
				{28, 35, 0}}},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil},
		{false,
			[]int64{0xFFFFFFFF, 0x1F, 0x17, 0xE, 0x1D, 0x1C, 0x1B, 0x1A, 0x19, 0x1E, 0x21,
				0x20, 0x18},
			[]*rangeItem{{0, 0, 0}, {1, 1, 1},
				{2, 9, 2}, {10, 11, 3}, {12, 12, 4},
				{13, 13, 5}, {14, 14, 6}, {15, 15, 7},
				{16, 16, 8}, {17, 18, 2}, {19, 19, 9},
				{20, 20, 10}, {21, 21, 11}, {22, 23, 2},
				{24, 24, 12}, {25, 32, 2},
				{33, 35, 0}}},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil}, {true, nil, nil},
		{false, []int64{0xFFFFFFFF, 0x14},
			[]*rangeItem{{0, 25, 0}, {26, 26, 1},
				{27, 35, 0}}},
		{true, nil, nil},
		{false, []int64{0xFFFFFFFF, 0x16},
			[]*rangeItem{{0, 21, 0}, {22, 24, 1},
				{25, 35, 0}}},
		{true, nil, nil},
		{false, []int64{0xFFFFFFFF, 0x22},
			[]*rangeItem{{0, 22, 0}, {23, 24, 1},
				{25, 35, 0}}},
		{false, []int64{0xFFFFFFFF, 0x23},
			[]*rangeItem{{0, 10, 0}, {11, 11, 1},
				{12, 12, 0}, {13, 14, 1}, {15, 17, 0},
				{18, 18, 1}, {19, 19, 0}, {20, 20, 1},
				{21, 21, 0}, {22, 24, 1},
				{25, 35, 0}}},
		{false, []int64{0xFFFFFFFF, 0x24},
			[]*rangeItem{{0, 10, 0}, {11, 11, 1},
				{12, 12, 0}, {13, 14, 1}, {15, 17, 0},
				{18, 18, 1}, {19, 19, 0}, {20, 20, 1},
				{21, 21, 0}, {22, 24, 1},
				{25, 35, 0}}},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil}, {true, nil, nil},
		{true, nil, nil},
		{false, []int64{0xFFFFFFFF, 0x25},
			[]*rangeItem{{0, 22, 0}, {23, 24, 1},
				{25, 35, 0}}},
		{false, []int64{0xFFFFFFFF, 0x26},
			[]*rangeItem{{0, 10, 0}, {11, 11, 1},
				{12, 12, 0}, {13, 14, 1}, {15, 17, 0},
				{18, 18, 1}, {19, 19, 0}, {20, 20, 1},
				{21, 21, 0}, {22, 24, 1},
				{25, 35, 0}}},
		{false, []int64{0xFFFFFFFF, 0x27},
			[]*rangeItem{{0, 10, 0}, {11, 11, 1},
				{12, 12, 0}, {13, 14, 1}, {15, 17, 0},
				{18, 18, 1}, {19, 19, 0}, {20, 20, 1},
				{21, 21, 0}, {22, 24, 1},
				{25, 35, 0}}},
		{true, nil, nil}, {true, nil, nil},
		{false, []int64{0xFFFFFFFF, 0x28},
			[]*rangeItem{{0, 10, 0}, {11, 11, 1},
				{12, 12, 0}, {13, 14, 1}, {15, 17, 0},
				{18, 18, 1}, {19, 19, 0}, {20, 20, 1},
				{21, 21, 0}, {22, 24, 1},
				{25, 35, 0}}},
		{false, []int64{0xFFFFFFFF, 0x29},
			[]*rangeItem{{0, 10, 0}, {11, 11, 1},
				{12, 12, 0}, {13, 14, 1}, {15, 17, 0},
				{18, 18, 1}, {19, 19, 0}, {20, 20, 1},
				{21, 21, 0}, {22, 24, 1},
				{25, 35, 0}}},
		{true, nil, nil}}
	lx.finalTable = make(map[int64]int64)
	lx.finalTable[0x4] = 0x0
	lx.finalTable[0x5] = 0x4
	lx.finalTable[0x6] = 0x1
	lx.finalTable[0x7] = 0x2
	lx.finalTable[0x8] = 0x1C
	lx.finalTable[0x9] = 0x3
	lx.finalTable[0xA] = 0x6
	lx.finalTable[0xB] = 0x5
	lx.finalTable[0xC] = 0xA
	lx.finalTable[0xD] = 0x1C
	lx.finalTable[0xE] = 0x12
	lx.finalTable[0xF] = 0x1B
	lx.finalTable[0x10] = 0x8
	lx.finalTable[0x11] = 0x7
	lx.finalTable[0x12] = 0x9
	lx.finalTable[0x13] = 0xE
	lx.finalTable[0x14] = 0xD
	lx.finalTable[0x15] = 0xB
	lx.finalTable[0x16] = 0xC
	lx.finalTable[0x17] = 0x1A
	lx.finalTable[0x18] = 0x1A
	lx.finalTable[0x19] = 0x1A
	lx.finalTable[0x1A] = 0x1A
	lx.finalTable[0x1B] = 0x16
	lx.finalTable[0x1C] = 0x17
	lx.finalTable[0x1D] = 0x13
	lx.finalTable[0x1E] = 0x15
	lx.finalTable[0x1F] = 0x18
	lx.finalTable[0x20] = 0x14
	lx.finalTable[0x21] = 0x19
	lx.finalTable[0x25] = 0xF
	lx.finalTable[0x26] = 0x10
	lx.finalTable[0x29] = 0x11
	lx.inputTable = []*rangeItem{{0, 8, 17}, {9, 9, 26},
		{10, 10, 0}, {11, 12, 17}, {13, 13, 0},
		{14, 31, 17}, {32, 32, 26}, {33, 39, 17},
		{40, 40, 31}, {41, 41, 5}, {42, 42, 32},
		{43, 43, 30}, {44, 44, 25}, {45, 45, 28},
		{46, 46, 2}, {47, 47, 1}, {48, 48, 24},
		{49, 55, 23}, {56, 57, 22}, {58, 62, 17},
		{63, 63, 29}, {64, 64, 17}, {65, 70, 18},
		{71, 90, 17}, {91, 91, 6}, {92, 92, 3},
		{93, 93, 8}, {94, 94, 9}, {95, 96, 17},
		{97, 97, 18}, {98, 98, 14}, {99, 99, 20},
		{100, 100, 11}, {101, 101, 18}, {102, 102, 13},
		{103, 109, 17}, {110, 110, 21}, {111, 113, 17},
		{114, 114, 12}, {115, 115, 10}, {116, 116, 19},
		{117, 117, 15}, {118, 118, 17}, {119, 119, 10},
		{120, 120, 16}, {121, 122, 17}, {123, 123, 4},
		{124, 124, 7}, {125, 125, 27}, {126, 65535, 17}}
	lx.initialTable = make(map[string]int64)
	lx.initialTable["REPEAT"] = 0x1
	lx.initialTable["BRACKET"] = 0x2
	lx.initialTable["INITIAL"] = 0x3
	return
}

func (lx *lexer) setSource(src string) {
	if src != "" {
		lx.source = src
	}
	lx.ended = false
	lx.start = 0
	lx.oldStart = 0
	lx.line = 1
	lx.column = 0
	lx.advanced = true
	lx.tokenName = ""
	lx.yy = nil
	lx.initialState = _initial
	lx.initialInput = lx.initialTable[lx.initialState]
}

func (lx *lexer) getToken() (string, error) {
	var err error
	if lx.advanced {
		lx.tokenName, err = lx.next()
		lx.advanced = false
	}
	return lx.tokenName, err
}

func (lx *lexer) getPositionInfo() string {
	return fmt.Sprintf("row(%d) column(%d)", lx.line, lx.column)
}

func (lx *lexer) next() (ret string, err error) {
	for {
		var (
			nextState         int64
			ch                int64
			och               = _maxValue
			next              = lx.start
			curState          = lx.transTable[0].toStates[lx.initialInput]
			lastFinalState    = int64(_deadState)
			lastFinalPosition = lx.start
		)
		for {
			if next < int64(len(lx.source)) {
				ch = int64(lx.source[next])
				// 计算行、列的位置
				if och != _maxValue {
					if ch == 0x0d { // \r符号
						lx.column = 0
						lx.line++
					} else if ch == 0x0a { // \n
						if och != 0x0d { // != \r
							lx.column = 0
							lx.line++
						}
					} else {
						lx.column++
					}
				}
				och = int(ch)
				if nextState, err = lx.trans(curState, ch); err != nil {
					return
				}
			} else {
				nextState = _deadState
			}
			//OK
			if nextState == _deadState {
				if lx.start == lastFinalPosition {
					if lx.start == int64(len(lx.source)) {
						if !lx.ended {
							lx.ended = true
							return "<$>", nil
						}
						return "", errEOF
					}
					return "", fmt.Errorf("意外的字符(line:%d,col:%d) of %s", lx.line, lx.column, lx.source)
				}
				lx.yyText = lx.source[lx.start:lastFinalPosition]
				lx.oldStart = lx.start
				lx.start = lastFinalPosition
				fIndex := lx.finalTable[lastFinalState]
				switch fIndex {
				case 0x0:
					return "*", nil
				case 0x1:
					return "+", nil
				case 0x2:
					return "?", nil
				case 0x3:
					return "|", nil
				case 0x4:
					return "(", nil
				case 0x5:
					return ")", nil
				case 0x6:
					if err = lx.begin("BRACKET"); err != nil {
						return
					}
					return "[", nil
				case 0x7:
					return "^", nil
				case 0x8:
					return "-", nil
				case 0x9:
					if err = lx.begin("INITIAL"); err != nil {
						return
					}
					return "]", nil
				case 0xA:
					if err = lx.begin("REPEAT"); err != nil {
						return
					}
					return "{", nil
				case 0xB:
					return ",", nil
				case 0xC:
					if lx.yyText, err = strconv.ParseInt(lx.yyText.(string), 10, 64); err != nil {
						return
					}
					return "d", nil
				case 0xE:
					if err = lx.begin("INITIAL"); err != nil {
						return
					}
					return "}", nil
				case 0xF:
					var tmp int64
					if tmp, err = strconv.ParseInt(lx.yyText.(string)[2:4], 8, 64); err != nil {
						return
					}
					lx.yyText = string(tmp)
					return "c", nil
				case 0x10:
					var tmp int64
					if tmp, err = strconv.ParseInt(lx.yyText.(string)[2:4], 16, 64); err != nil {
						return
					}
					lx.yyText = string(tmp)
					return "c", nil
				case 0x11:
					var tmp int64
					if tmp, err = strconv.ParseInt(lx.yyText.(string)[2:6], 16, 64); err != nil {
						return
					}
					lx.yyText = string(tmp)
					return "c", nil
				case 0x12:
					return "escc", nil
				case 0x13:
					lx.yyText = "\r"
					return "c", nil
				case 0x14:
					lx.yyText = "\n"
					return "c", nil
				case 0x15:
					lx.yyText = "\t"
					return "c", nil
				case 0x16:
					lx.yyText = "\b"
					return "c", nil
				case 0x17:
					lx.yyText = "\f"
					return "c", nil
				case 0x18:
					lx.yyText = "/"
					return "c", nil
				case 0x19:
					return "escc", nil
				case 0x1A:
					lx.yyText = lx.yyText.(string)[1:2]
					return "c", nil
				case 0x1B:
					return "/", nil
				case 0x1C:
					return "c", nil
				}
				break
			} else {
				next++
				if _, ok := lx.finalTable[nextState]; ok {
					lastFinalState = nextState
					lastFinalPosition = next
				}
				curState = nextState
			}
		}
	}
}

func (lx *lexer) begin(state string) error {
	return lx.setInitialState(state)
}

func (lx *lexer) setInitialState(state string) (err error) {
	if _, ok := lx.initialTable[state]; !ok {
		err = fmt.Errorf("未定义的初始状态:%s", state)
		return
	}
	lx.initialState = state
	lx.initialInput = lx.initialTable[state]
	return
}

func (lx *lexer) trans(curState, ch int64) (int64, error) {
	if ch < lx.inputTable[0].from || ch > lx.inputTable[len(lx.inputTable)-1].to {
		return 0, fmt.Errorf("line:%d,column:%d 输入字符超出范围", lx.line, lx.column)
	}
	if lx.transTable[curState].isDead {
		return _deadState, nil
	}
	pubInput := find(ch, lx.inputTable)
	innerInput := find(pubInput, lx.transTable[curState].transEdge)
	return lx.transTable[curState].toStates[innerInput], nil
}

func find(code int64, table []*rangeItem) int64 {
	var (
		max = len(table) - 1
		min int
		mid uint64
	)
	for {
		mid = uint64(min+max) >> 1
		if table[mid].from <= code {
			if table[mid].to >= code {
				return table[mid].value
			}
			min = int(mid) + 1
		} else {
			max = int(mid) - 1
		}
	}
}
