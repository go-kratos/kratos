package actriearea

import (
	"bytes"
	"container/list"
	"strings"
	"unicode"
)

type trieNode struct {
	// relation data
	isEnd bool               // 是否是结束节点
	index int                // 第几个节点
	fail  *trieNode          // 失败跳转节点
	child map[rune]*trieNode // 节点下面的子节点

	// extra data
	fid   int64   // 这个几点对应的filter id
	level int8    // 过滤等级
	tpIds []int64 // 这个节点的过滤有有哪些分区 type id
}

func newTrieNode() *trieNode {
	return &trieNode{
		fail:  nil,
		index: -1,
		child: make(map[rune]*trieNode),
	}
}

type Matcher struct {
	size        int
	root        *trieNode
	indexLenMap map[int]int
}

func NewMatcher() *Matcher {
	return &Matcher{
		size:        0,
		root:        newTrieNode(),
		indexLenMap: make(map[int]int),
	}
}

// Rule define filter rule
type Rule struct {
	Msg   string  `json:"msg"`
	Level int8    `json:"level"`
	Tpids []int64 `json:"tpids"`
	Fid   int64   `json:"fid"`
	Mode  int8    `json:"mode"`
}

type MatchHits struct {
	Fid       int64
	Area      string
	Mode      int8
	Level     int8
	Rule      string
	Word      string
	TypeIDs   []int64
	Positions [][]int
}

func (m *Matcher) Insert(s string, level int8, tpIds []int64, fid int64) {
	m.insert(s, level, tpIds, fid)
}

func (m *Matcher) Build() {
	m.build()
}

func (m *Matcher) Filter(in string, tpID int64, critical int8) (hits []*MatchHits) {
	var p *trieNode
	curNode := m.root
	re := []rune(in)
	byteMsg := []byte(in)
	buf := bytes.NewBuffer(byteMsg)
	bufLen := buf.Len()
	byteSize := 0
	for i := 0; i < bufLen; i++ {
		if v, vSize, err := buf.ReadRune(); err == nil {
			v = unicode.ToLower(v)
			// add out byte first
			byteSize = byteSize + vSize
			for curNode.child[v] == nil && curNode != m.root {
				curNode = curNode.fail
			}
			curNode = curNode.child[v]
			if curNode == nil {
				curNode = m.root
			}
			p = curNode
			if p != m.root && p.isEnd {
				for _, id := range p.tpIds {
					if id == 0 || id == tpID {
						start := i - m.indexLenMap[p.index] + 1
						hitWord := string(re[start : i+1])
						hit := &MatchHits{Level: p.level, Fid: p.fid, Mode: 1}
						hit.TypeIDs = append(hit.TypeIDs, []int64{0, tpID}...)
						hit.Rule = hitWord
						if p.level > critical {
							var positions = make([][]int, 1)
							positions[0] = []int{byteSize - len([]byte(hitWord)), byteSize}
							hit.Positions = positions
						}
						hits = append(hits, hit)
						p = p.fail
						break
					}
				}
			}
		}
	}
	return
}

// GetWhiteHits .
// tpID < 0 == all , tpID = 0 全站分区 , tpID > 0 子分区
func (m *Matcher) GetWhiteHits(in string, tpID int64) (hits []*MatchHits) {
	curNode := m.root
	var p *trieNode = nil
	re := []rune(in)
	vs := []rune(strings.ToLower(in))
	byteOut := []byte{}
	for i, v := range vs {
		byteOut = append(byteOut, []byte(string(re[i]))...)
		for curNode.child[v] == nil && curNode != m.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = m.root
		}
		p = curNode
		if p != m.root && p.isEnd {
			if tpID < 0 {
				start := i - m.indexLenMap[p.index] + 1
				hitWord := string(re[start : i+1])
				hitByte := []byte(hitWord)
				var positions = make([][]int, 1)
				positions[0] = []int{len(byteOut) - len(hitByte), len(byteOut)}
				hit := &MatchHits{
					Positions: positions,
					Level:     p.level,
					Rule:      hitWord,
					Fid:       p.fid,
					Mode:      1,
					TypeIDs:   p.tpIds,
				}
				hits = append(hits, hit)
				// p = p.fail
			} else {
				for _, id := range p.tpIds {
					if id == 0 || id == tpID {
						start := i - m.indexLenMap[p.index] + 1
						hitWord := string(re[start : i+1])
						hitByte := []byte(hitWord)
						var positions = make([][]int, 1)
						positions[0] = []int{len(byteOut) - len(hitByte), len(byteOut)}
						hit := &MatchHits{
							Positions: positions,
							Level:     p.level,
							Rule:      hitWord,
							Fid:       p.fid,
							Mode:      1,
							TypeIDs:   p.tpIds,
						}
						hits = append(hits, hit)
						// p = p.fail
						break
					}
				}
			}
		}
	}
	return
}

func (m *Matcher) ArtFilter(in string, tpid int64) (fmsg []string) {
	curNode := m.root
	var p *trieNode = nil
	vs := []rune(strings.ToLower(in))
	for i, v := range vs {
		for curNode.child[v] == nil && curNode != m.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = m.root
		}
		p = curNode
		if p != m.root && p.isEnd {
			for _, id := range p.tpIds {
				if id == 0 || id == tpid {
					start := i - m.indexLenMap[p.index] + 1
					end := i + 1
					fmsg = append(fmsg, string(vs[start:end]))
					// p = p.fail
					break
				}
			}
		}
	}
	return
}

func (m *Matcher) Test(in string) (level int8, rs []*Rule) {
	curNode := m.root
	var p *trieNode = nil
	vs := []rune(strings.ToLower(in))
	for i, v := range vs {
		for curNode.child[v] == nil && curNode != m.root {
			curNode = curNode.fail
		}
		curNode = curNode.child[v]
		if curNode == nil {
			curNode = m.root
		}
		p = curNode
		if p != m.root && p.isEnd {
			start := i - m.indexLenMap[p.index] + 1
			if p.level > level {
				level = p.level
			}
			r := &Rule{
				Msg:   string(vs[start : i+1]),
				Level: p.level,
				Tpids: p.tpIds,
				Fid:   p.fid,
				Mode:  1,
			}
			rs = append(rs, r)
			// p = p.fail
		}
	}
	return
}

func (m *Matcher) build() {
	ll := list.New()
	ll.PushBack(m.root)
	for ll.Len() > 0 {
		temp := ll.Remove(ll.Front()).(*trieNode)
		var p *trieNode = nil
		for i, v := range temp.child {
			if temp == m.root {
				v.fail = m.root
			} else {
				p = temp.fail
				for p != nil {
					if p.child[i] != nil {
						v.fail = p.child[i]
						break
					}
					p = p.fail
				}
				if p == nil {
					v.fail = m.root
				}
			}
			ll.PushBack(v)
		}
	}
}

func (m *Matcher) insert(s string, level int8, tpIDs []int64, fid int64) {
	curNode := m.root
	for _, v := range s {
		if curNode.child[v] == nil {
			curNode.child[v] = newTrieNode()
		}
		curNode = curNode.child[v]
	}
	curNode.isEnd = true

	// extra data
	m.indexLenMap[m.size] = len([]rune(s))
	curNode.level = level
	curNode.fid = fid
	curNode.tpIds = tpIDs
	curNode.index = m.size

	m.size++
}
