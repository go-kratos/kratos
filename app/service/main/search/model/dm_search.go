package model

import (
	"fmt"
	"strconv"
)

// DmSearchParams .
type DmSearchParams struct {
	Bsp        *BasicSearchParams
	Oid        int64  `form:"oid" params:"oid" default:"-1"`
	Mid        int64  `form:"mid" params:"mid" default:"-1"`
	Mode       int    `form:"mode" params:"mode" default:"-1"`
	Pool       int    `form:"pool" params:"pool" default:"-1"`
	Progress   int    `form:"progress" params:"progress" default:"-1"`
	States     []int  `form:"states,split" params:"states"`
	Type       int    `form:"type" params:"type" default:"-1"`
	AttrFormat []int  `form:"attr_format,split" params:"attr_format"`
	CtimeFrom  string `form:"ctime_from" params:"ctime_from"`
	CtimeTo    string `form:"ctime_to" params:"ctime_to"`
}

// DmUptParams .
type DmUptParams struct {
	ID    int64 `json:"id"`
	Oid   int64 `json:"oid"`
	Field map[string]interface{}
}

// IndexName .
func (m *DmUptParams) IndexName() string {
	return "dm_search_" + strconv.FormatInt(m.Oid%1000, 10)
}

// IndexType .
func (m *DmUptParams) IndexType() string {
	return "base"
}

// IndexID .
func (m *DmUptParams) IndexID() string {
	return fmt.Sprintf("%d", m.ID)
}

// PField .
func (m *DmUptParams) PField() map[string]interface{} {
	return m.Field
}
