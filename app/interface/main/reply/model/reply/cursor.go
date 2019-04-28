package reply

import (
	"errors"
	"fmt"
	"sort"

	"go-common/library/log"
)

// RootReplyListHeader RootReplyListHeader
type RootReplyListHeader struct {
	TopAdmin *Reply
	TopUpper *Reply
	Hots     []*Reply
}

// RootReplyList RootReplyList
type RootReplyList struct {
	Subject *Subject

	Header *RootReplyListHeader

	TopAdmin *Reply
	TopUpper *Reply
	Hots     []*Reply

	Roots                          []*Reply
	CursorRangeMax, CursorRangeMin int64
}

// CursorParams CursorParams
type CursorParams struct {
	Mid        int64
	Oid        int64
	OTyp       int8
	Sort       int8
	HTMLEscape bool
	IP         string
	RootID     int64
	Cursor     *Cursor
	HotSize    int
	ShowFolded bool
}

// Comp Comp
type Comp func(x, y int64) bool

// order
var (
	OrderDESC = func(x, y int64) bool { return x > y }
	OrderASC  = func(x, y int64) bool { return x < y }
)

// SortArr SortArr
func SortArr(arr []int64, cmp Comp) []int64 {
	if cmp == nil {
		log.Warn("unexpected: cmp is nil")
		cmp = OrderDESC
	}
	less := func(i, j int) bool { return cmp(arr[i], arr[j]) }
	if sort.SliceIsSorted(arr, less) {
		return arr
	}
	sort.Slice(arr, less)
	return arr
}

// Len Len
func (c *Cursor) Len() int {
	return c.length
}

// Current Current
func (c *Cursor) Current() int64 {
	if c.maxID == 0 {
		return c.minID
	}
	return c.maxID
}

// Latest Latest
func (c *Cursor) Latest() bool {
	return c.maxID == 0 && c.minID == 0
}

// Descrease Descrease
func (c *Cursor) Descrease() bool {
	return c.maxID > 0
}

// Increase Increase
func (c *Cursor) Increase() bool {
	return c.minID > 0
}

// Max return maxID
func (c *Cursor) Max() int64 {
	return c.maxID
}

// Min return minID
func (c *Cursor) Min() int64 {
	return c.minID
}

// Sort Sort
func (c *Cursor) Sort(arr []int64) []int64 {
	return SortArr(arr, c.comp)
}

// String String
func (c *Cursor) String() string {
	return fmt.Sprintf("current: %d, growIncr: %v, size: %d", c.Current(), c.Increase(), c.length)
}

// NewCursor NewCursor
func NewCursor(maxID int64, minID int64, size int, cmp Comp) (*Cursor, error) {
	if maxID < 0 || minID < 0 || cmp == nil {
		return nil, fmt.Errorf("either max_id(%d) or min_id(%d) < 0 or cmp = nil", maxID, minID)
	}
	if (minID * maxID) != 0 {
		return nil, errors.New("both max_id and max_id > 0")
	}
	return &Cursor{
		length: size,
		comp:   cmp,
		maxID:  maxID,
		minID:  minID,
	}, nil
}

// Cursor Cursor
type Cursor struct {
	maxID, minID int64
	length       int
	comp         Comp
}
