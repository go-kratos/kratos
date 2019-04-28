package xregex

/*
golang version regex parser
refer to: https://github.com/aristotle9/as3cc/tree/master/java-template/src/org/lala/lex/utils/parser
*/

import (
	"container/list"
	"fmt"
	"strings"
)

// Parser golang regex parser
type Parser struct {
	actionTable []map[int64]int64
	gotoTable   map[int64]map[int64]int64
	prodList    []*productionItem
	inputTable  map[string]int64
	codes       []interface{}
	duplCount   int64
	num         int64
	flags       string
	flagsC      string
}

// Info print parser info.
func (p *Parser) Info() (s string) {
	s += fmt.Sprintf("actionTable len:%d,", len(p.actionTable))
	s += fmt.Sprintf("prodList len:%d,", len(p.prodList))
	s += fmt.Sprintf("inputTable len:%d,", len(p.inputTable))
	s += fmt.Sprintf("codes len:%d\n", len(p.codes))
	return
}

// New return regex parser instance.
func New() (p *Parser) {
	p = &Parser{}
	var tmp map[int64]int64
	p.actionTable = make([]map[int64]int64, 0)
	p.actionTable = append(p.actionTable, nil)
	tmp = make(map[int64]int64)
	tmp[0x20] = 0x4
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x6] = 0x1
	tmp[0x7] = 0x2
	tmp[0x8] = 0x1C
	tmp[0x2A] = 0x24
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x27] = 0x6
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x16] = 0x8
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x1D] = 0xA
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x0] = 0x5
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x2C] = 0x22
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x2F
	tmp[0x2A] = 0x24
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x1D] = 0xA
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0x21
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x21
	tmp[0xA] = 0x21
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x1D] = 0xA
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x0] = 0x7
	tmp[0x1] = 0x7
	tmp[0x2] = 0x7
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0x7
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x7
	tmp[0xA] = 0x7
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0xF] = 0x7
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x21] = 0x18
	tmp[0x1D] = 0x7
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x2C] = 0x22
	tmp[0x8] = 0x1C
	tmp[0xA] = 0x31
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x1D] = 0xA
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0x23
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x23
	tmp[0xA] = 0x23
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x23
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x1D] = 0x23
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0x25
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x25
	tmp[0xA] = 0x25
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x25
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x1D] = 0x25
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0x27
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x27
	tmp[0xA] = 0x27
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x27
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x1D] = 0x27
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x0] = 0x9
	tmp[0x1] = 0x9
	tmp[0x2] = 0x9
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0x9
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x9
	tmp[0xA] = 0x9
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0xF] = 0x9
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x21] = 0x18
	tmp[0x1D] = 0x9
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x28] = 0x26
	tmp[0x22] = 0x2C
	tmp[0x1C] = 0x28
	tmp[0x5] = 0x30
	tmp[0x18] = 0x45
	tmp[0x1B] = 0x2E
	tmp[0xC] = 0x2E
	tmp[0xD] = 0x2A
	tmp[0xE] = 0x37
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x3] = 0x19
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x21] = 0x18
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0x29
	tmp[0x8] = 0x1C
	tmp[0x9] = 0x29
	tmp[0xA] = 0x29
	tmp[0x2B] = 0x20
	tmp[0x24] = 0x1E
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x29
	tmp[0x15] = 0xE
	tmp[0x17] = 0xC
	tmp[0x1A] = 0x16
	tmp[0x2A] = 0x24
	tmp[0x1D] = 0x29
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x1F] = 0x4F
	tmp[0x25] = 0x55
	tmp[0x13] = 0x3F
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x29] = 0x5B
	tmp[0x25] = 0x57
	tmp[0x1E] = 0x4B
	tmp[0x26] = 0x59
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x1E] = 0x4D
	tmp[0x13] = 0x41
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x0] = 0xB
	tmp[0x1] = 0xB
	tmp[0x2] = 0xB
	tmp[0x3] = 0x1B
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0xB
	tmp[0x8] = 0x1C
	tmp[0x9] = 0xB
	tmp[0xA] = 0xB
	tmp[0xB] = 0x1B
	tmp[0xC] = 0x2E
	tmp[0xD] = 0x2A
	tmp[0xE] = 0x39
	tmp[0xF] = 0xB
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x16] = 0x8
	tmp[0x17] = 0xC
	tmp[0x18] = 0x39
	tmp[0x19] = 0x47
	tmp[0x1A] = 0x16
	tmp[0x1B] = 0x2E
	tmp[0x1C] = 0x28
	tmp[0x1D] = 0xB
	tmp[0x20] = 0xB
	tmp[0x21] = 0x18
	tmp[0x22] = 0x2C
	tmp[0x23] = 0x53
	tmp[0x24] = 0x1E
	tmp[0x27] = 0x6
	tmp[0x28] = 0x26
	tmp[0x2A] = 0x24
	tmp[0x2B] = 0x20
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0xC] = 0x35
	tmp[0x1B] = 0x49
	p.actionTable = append(p.actionTable, tmp)
	tmp = make(map[int64]int64)
	tmp[0x0] = 0xD
	tmp[0x1] = 0xD
	tmp[0x2] = 0xD
	tmp[0x3] = 0xD
	tmp[0x4] = 0x2E
	tmp[0x5] = 0x30
	tmp[0x7] = 0xD
	tmp[0x8] = 0x1C
	tmp[0x9] = 0xD
	tmp[0xA] = 0xD
	tmp[0xB] = 0xD
	tmp[0xC] = 0x2E
	tmp[0xD] = 0x2A
	tmp[0xE] = 0xD
	tmp[0xF] = 0xD
	tmp[0x10] = 0x10
	tmp[0x11] = 0x12
	tmp[0x12] = 0x14
	tmp[0x14] = 0x1A
	tmp[0x15] = 0xE
	tmp[0x16] = 0x8
	tmp[0x17] = 0xC
	tmp[0x18] = 0xD
	tmp[0x1A] = 0x16
	tmp[0x1B] = 0x2E
	tmp[0x1C] = 0x28
	tmp[0x1D] = 0xD
	tmp[0x20] = 0xD
	tmp[0x21] = 0x18
	tmp[0x22] = 0x2C
	tmp[0x24] = 0x1E
	tmp[0x27] = 0x6
	tmp[0x28] = 0x26
	tmp[0x2A] = 0x24
	tmp[0x2B] = 0x20
	tmp[0x2C] = 0x22
	p.actionTable = append(p.actionTable, tmp)
	p.gotoTable = make(map[int64]map[int64]int64)
	tmp = make(map[int64]int64)
	tmp[0xB] = 0x33
	tmp[0x3] = 0x1F
	p.gotoTable[0x18] = tmp
	tmp = make(map[int64]int64)
	tmp[0x0] = 0xF
	p.gotoTable[0x13] = tmp
	tmp = make(map[int64]int64)
	tmp[0x0] = 0x11
	tmp[0x1] = 0x15
	tmp[0x2] = 0x17
	tmp[0x14] = 0x2B
	tmp[0x7] = 0x2B
	tmp[0x9] = 0x2B
	tmp[0xA] = 0x2B
	tmp[0x1D] = 0x2B
	tmp[0xF] = 0x3D
	p.gotoTable[0x14] = tmp
	tmp = make(map[int64]int64)
	tmp[0x16] = 0x43
	p.gotoTable[0x15] = tmp
	tmp = make(map[int64]int64)
	tmp[0x0] = 0x13
	tmp[0x1] = 0x13
	tmp[0x2] = 0x13
	tmp[0x3] = 0x1D
	tmp[0x7] = 0x13
	tmp[0x20] = 0x51
	tmp[0x9] = 0x13
	tmp[0xA] = 0x13
	tmp[0xB] = 0x1D
	tmp[0xE] = 0x3B
	tmp[0xF] = 0x13
	tmp[0x14] = 0x13
	tmp[0x18] = 0x3B
	tmp[0x1D] = 0x13
	p.gotoTable[0x16] = tmp
	tmp = make(map[int64]int64)
	tmp[0x9] = 0x2D
	tmp[0xA] = 0x2D
	tmp[0x14] = 0x2D
	tmp[0x1D] = 0x2D
	tmp[0x7] = 0x2D
	p.gotoTable[0x17] = tmp
	p.prodList = []*productionItem{{0x19, 0x2}, {0x13, 0x1},
		{0x13, 0x4}, {0x15, 0x2}, {0x15, 0x0},
		{0x14, 0x3}, {0x14, 0x3}, {0x14, 0x2},
		{0x14, 0x2}, {0x14, 0x2}, {0x14, 0x2},
		{0x14, 0x3}, {0x14, 0x4}, {0x14, 0x2},
		{0x14, 0x1}, {0x17, 0x3}, {0x17, 0x4},
		{0x17, 0x5}, {0x17, 0x4}, {0x18, 0x4},
		{0x18, 0x2}, {0x18, 0x1}, {0x18, 0x3},
		{0x16, 0x1}, {0x16, 0x1}}
	p.inputTable = make(map[string]int64)
	p.inputTable["/"] = 0x2
	p.inputTable["["] = 0x9
	p.inputTable["escc"] = 0x12
	p.inputTable["]"] = 0xA
	p.inputTable["("] = 0x4
	p.inputTable["{"] = 0xC
	p.inputTable["^"] = 0xB
	p.inputTable[")"] = 0x5
	p.inputTable["|"] = 0x3
	p.inputTable["?"] = 0x8
	p.inputTable["}"] = 0xE
	p.inputTable["+"] = 0x7
	p.inputTable["*"] = 0x6
	p.inputTable[","] = 0xF
	p.inputTable["<$>"] = 0x1
	p.inputTable["-"] = 0x11
	p.inputTable["c"] = 0x10
	p.inputTable["d"] = 0xD
	return
}

func (p *Parser) put(data interface{}) error {
	if fmt.Sprint(data) == "dupl" {
		p.duplCount++
		if p.duplCount > 100 {
			return fmt.Errorf("dupl commands over 100")
		}
	}
	return nil
}

func (p *Parser) parse(lx *lexer) (ret interface{}, err error) {
	var (
		act         int64
		token       string
		tokenID     int64
		state       int64
		actObj      int64
		stateStack  = list.New()
		outputStack []interface{}
	)
	stateStack.PushBack(uint(0))
	for {
		if token, err = lx.getToken(); err != nil {
			return
		}
		tokenID = p.inputTable[token]
		state = int64(stateStack.Front().Value.(uint))
		actObj = p.actionTable[tokenID][state]
		if actObj == 0 {
			err = fmt.Errorf("Parse Error: %s", lx.getPositionInfo())
			return
		}
		act = actObj
		if act == 1 {
			ret = outputStack[len(outputStack)-1]
			return
		} else if (act & 1) == 1 {
			outputStack = append(outputStack, lx.yyText)
			stateStack.PushFront(uint(act)>>1 - 1)
			lx.advanced = true
		} else if (act & 1) == 0 {
			pi := uint(act) >> 1
			length := p.prodList[pi].bodyLength
			var result interface{}
			if length > 0 {
				result = outputStack[int64(len(outputStack))-length]
			}
			switch pi {
			case 0x1:
				result = p.codes
			case 0x2:
				result = p.codes
			case 0x3:
				p.flagsC = outputStack[len(outputStack)-1].(string)
				if !strings.Contains(p.flags, p.flagsC) {
					err = fmt.Errorf("Flag Repeated:%s", lx.getPositionInfo())
					return
				}
				if !strings.Contains("igm", p.flagsC) {
					err = fmt.Errorf("Unknow Flag:%s", lx.getPositionInfo())
					return
				}
				p.flags += p.flagsC
			case 0x4:
			case 0x5:
				if err = p.put("or"); err != nil {
					return
				}
			case 0x6:
			case 0x7:
			case 0x8:
				if err = p.put("star"); err != nil {
					return
				}
			case 0x9:
				if err = p.put("more"); err != nil {
					return
				}
			case 0xA:
				if err = p.put("ask"); err != nil {
					return
				}
			case 0xB:
				if err = p.put([]interface{}{"include", outputStack[len(outputStack)-2]}); err != nil {
					return
				}
			case 0xC:
				if err = p.put([]interface{}{"exclude", outputStack[len(outputStack)-2]}); err != nil {
					return
				}
			case 0xD:
				if err = p.put("cat"); err != nil {
					return
				}
			case 0xE:
				if err = p.put("single"); err != nil {
					return
				}
			case 0xF:
				p.num = (outputStack[len(outputStack)-2]).(int64) - 1
				for p.num > 0 {
					if err = p.put("dupl"); err != nil {
						return
					}
					p.num--
				}
				p.num = (outputStack[len(outputStack)-2]).(int64) - 1
				for p.num > 0 {
					if err = p.put("cat"); err != nil {
						return
					}
					p.num--
				}
			case 0x10:
				if err = p.put("ask"); err != nil {
					return
				}
				p.num = (outputStack[len(outputStack)-2]).(int64) - 1
				for p.num > 0 {
					if err = p.put("dupl"); err != nil {
						return
					}
					p.num--
				}
				p.num = (outputStack[len(outputStack)-2]).(int64) - 1
				for p.num > 0 {
					if err = p.put("cat"); err != nil {
						return
					}
					p.num--
				}
			case 0x11:
				p.num = (outputStack[len(outputStack)-4]).(int64) - 1
				for p.num > 0 {
					if err = p.put("dupl"); err != nil {
						return
					}
					p.num--
				}
				p.num = (outputStack[len(outputStack)-2]).(int64) - (outputStack[len(outputStack)-4]).(int64)
				if p.num > 0 {
					if err = p.put("dupl"); err != nil {
						return
					}
					if err = p.put("ask"); err != nil {
						return
					}
					for p.num > 1 {
						if err = p.put("dupl"); err != nil {
							return
						}
						p.num--
					}
				}
				p.num = (outputStack[len(outputStack)-2]).(int64) - 1
				for p.num > 0 {
					if err = p.put("cat"); err != nil {
						return
					}
					p.num--
				}
			case 0x12:
				p.num = (outputStack[len(outputStack)-3]).(int64)
				for p.num > 0 {
					if err = p.put("dupl"); err != nil {
						return
					}
					p.num--
				}
				if err = p.put("star"); err != nil {
					return
				}
				p.num = (outputStack[len(outputStack)-3]).(int64)
				for p.num > 0 {
					if err = p.put("cat"); err != nil {
						return
					}
					p.num--
				}
			case 0x13:
				p.num = (outputStack[len(outputStack)-4]).(int64)
				result = p.num + 1
				if err = p.put([]interface{}{"c", (outputStack[len(outputStack)-3])}); err != nil {
					return
				}
				if err = p.put([]interface{}{"c", (outputStack[len(outputStack)-1])}); err != nil {
					return
				}
				if err = p.put("range"); err != nil {
					return
				}
			case 0x14:
				p.num = (outputStack[len(outputStack)-2]).(int64)
				result = p.num + 1
				if err = p.put("single"); err != nil {
					return
				}
			case 0x15:
				if err = p.put("single"); err != nil {
					return
				}
				result = int64(1)
			case 0x16:
				if err = p.put([]interface{}{"c", (outputStack[len(outputStack)-3])}); err != nil {
					return
				}
				if err = p.put([]interface{}{"c", (outputStack[len(outputStack)-1])}); err != nil {
					return
				}
				if err = p.put("range"); err != nil {
					return
				}
				result = int64(1)
			case 0x17:
				if err = p.put([]interface{}{"c", (outputStack[len(outputStack)-1])}); err != nil {
					return
				}
			case 0x18:
				if "\\c" == outputStack[len(outputStack)-1].(string) {
					err = fmt.Errorf("Control Character:%s", lx.getPositionInfo())
					return
				}
				if err = p.put([]interface{}{"escc", (outputStack[len(outputStack)-1])}); err != nil {
					return
				}
			}
			/** actions applying end **/
			var j int64
			var e *list.Element
			for j < length {
				e = stateStack.Front()
				stateStack.Remove(e)
				outputStack = outputStack[:len(outputStack)-1]
				j++
			}
			state = int64(stateStack.Front().Value.(uint))
			actObj = p.gotoTable[p.prodList[pi].headerID][state]
			if actObj == 0 {
				err = fmt.Errorf("goto error! %s", lx.getPositionInfo())
				return
			}
			act = actObj
			stateStack.PushFront(uint(act)>>1 - 1)
			outputStack = append(outputStack, result)
		}
	}
}

// Parse parse input string.
func (p *Parser) Parse(source string) (result interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("panic:%v", e)
			return
		}
	}()
	if source == "" {
		return
	}
	lx := newLexer()
	lx.setSource(source)
	return p.parse(lx)
}
