package model

import (
	"fmt"
)

var (
	// segmentLength 分段长度，根据视频时长做分段，单位：毫秒
	segmentLength = int64(6 * 60 * 1000)

	_defaultSeg = &Segment{Start: 0, End: DefaultVideoEnd, Cnt: 1, Num: 1, Duration: 0}
	// <d p="弹幕ID,弹幕属性,播放时间,弹幕模式,字体大小,颜色,发送时间,弹幕池,用户hash id">弹幕内容</d>
	_xmlSegFmt = `<d p="%d,%d,%d,%d,%d,%d,%d,%d,%s">%s</d>`
	// DefaultPage default page info
	DefaultPage = &Page{Num: 1, Size: DefaultVideoEnd, Total: 1}

	_xmlSegHeader = `<?xml version="1.0" encoding="UTF-8"?><i><oid>%d</oid><ps>%d</ps><pe>%d</pe><pc>%d</pc><pn>%d</pn><state>%d</state><real_name>%d</real_name>`
)

// const variable
const (
	// DefaultVideoEnd 当视频时长不存在或者为0时的默认视频结尾时间点
	DefaultVideoEnd = 10 * 60 * 60 * 1000
	// DefaultPageSize 默认分段长度
	DefaultPageSize = 60 * 6 * 1000
)

// Page dm page info
type Page struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
}

// Segment dm segment struct
type Segment struct {
	Start    int64 `json:"ps"`       // 分段起始时间
	End      int64 `json:"pe"`       // 分段结束时间
	Cnt      int64 `json:"cnt"`      // 总分段数
	Num      int64 `json:"num"`      // 当前第几段
	Duration int64 `json:"duration"` // 视频总时长
}

// ToXMLHeader convert segment to xml header format.
func (s *Segment) ToXMLHeader(oid int64, state, realname int32) string {
	return fmt.Sprintf(_xmlSegHeader, oid, s.Start, s.End, s.Cnt, s.Num, state, realname)
}

// SegmentInfo get segment info by start time and video duration.
func SegmentInfo(ps, duration int64) (s *Segment) {
	var cnt, num, pe int64
	if duration == 0 {
		s = _defaultSeg
		return
	}
	cnt = duration / DefaultPageSize
	if duration%DefaultPageSize > 0 {
		cnt++
	}
	for i := int64(0); i < cnt; i++ {
		if ps >= i*DefaultPageSize && ps < (i+1)*DefaultPageSize {
			ps = i * DefaultPageSize
			pe = (i + 1) * DefaultPageSize
			num = i + 1
		}
	}
	if pe > duration {
		pe = duration
	}
	s = &Segment{
		Start:    ps,
		End:      pe,
		Cnt:      cnt,
		Num:      num,
		Duration: duration,
	}
	return
}

// ToXMLSeg convert dm struct to xml.
func (d *DM) ToXMLSeg() (s string) {
	if d.Content == nil {
		return
	}
	msg := d.Content.Msg
	if d.ContentSpe != nil {
		msg = d.ContentSpe.Msg
	}
	if len(msg) == 0 {
		return
	}
	if d.Pool == PoolSpecial {
		msg = ""
	}
	s = fmt.Sprintf(_xmlSegFmt, d.ID, d.Attr, d.Progress, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, hash(d.Mid, uint32(d.Content.IP)), xmlReplace([]byte(msg)))
	return
}

// ToElem convert dm struct to element.
func (d *DM) ToElem() (e *Elem) {
	if d.Content == nil {
		return
	}
	msg := d.Content.Msg
	if d.ContentSpe != nil {
		msg = d.ContentSpe.Msg
	}
	if len(msg) == 0 {
		return
	}
	if d.Pool == PoolSpecial {
		msg = ""
	}
	// "弹幕ID,弹幕属性,播放时间,弹幕模式,字体大小,颜色,发送时间,弹幕池,用户hash id
	e = &Elem{
		Attribute: fmt.Sprintf(`%d,%d,%d,%d,%d,%d,%d,%d,%s`, d.ID, d.Attr, d.Progress, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, hash(d.Mid, uint32(d.Content.IP))),
		Content:   msg,
	}
	return
}

// SegmentPoint 根据当前段数和视频总时长计算分段的起始时间点
func SegmentPoint(num, duration int64) (ps, pe int64) {
	if duration == 0 {
		ps = 0
		pe = DefaultVideoEnd
		return
	}
	pe = num * segmentLength
	ps = pe - segmentLength
	if pe > duration {
		pe = duration
	}
	if ps < 0 {
		ps = 0
	}
	return
}
