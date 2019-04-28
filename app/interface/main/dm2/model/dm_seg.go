package model

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	_defaultSeg = &Segment{Start: 0, End: defaultVideoEnd, Cnt: 1, Num: 1, Duration: 0}
	// DefaultFlag default dm flag if bigdata downgrade.
	DefaultFlag = []byte(`{"rec_flag":2,"rec_text":"开启后，全站视频将按等级等优化弹幕","rec_switch":1,"dmflags":[]}`)
)

const (
	// segmentLength 分段长度,6分钟
	segmentLength = 60 * 6 * 1000
	// defaultVideoEnd 当视频时长不存在或者为0时的默认视频结尾时间点
	defaultVideoEnd = int64(10 * 60 * 60 * 1000)
	// DefaultVideoEnd 当视频时长不存在或者为0时的默认视频结尾时间点
	DefaultVideoEnd = int64(3 * 60 * 60 * 1000)
	// DefaultPageSize 默认分段长度
	DefaultPageSize = 60 * 6 * 1000
	// NotFound nothing found flag
	NotFound = int64(-1)
	// DefaultPage default page info

	// <d p="弹幕ID,弹幕属性,播放时间,弹幕模式,字体大小,颜色,发送时间,弹幕池,用户hash id">弹幕内容</d>
	_xmlSegFmt = `<d p="%d,%d,%d,%d,%d,%d,%d,%d,%s">%s</d>`
	// <d p="弹幕ID,弹幕属性,播放时间,弹幕模式,字体大小,颜色,发送时间,弹幕池,用户hash id,用户mid">弹幕内容</d>
	_xmlSegRealnameFmt = `<d p="%d,%d,%d,%d,%d,%d,%d,%d,%s,%d">%s</d>`
	_xmlSegHeader      = `<?xml version="1.0" encoding="UTF-8"?><i><oid>%d</oid><ps>%d</ps><pe>%d</pe><pc>%d</pc><pn>%d</pn><state>%d</state><real_name>%d</real_name>`
)

// JudgeSlice sort for dm judgement
type JudgeSlice []*DM

func (d JudgeSlice) Len() int           { return len(d) }
func (d JudgeSlice) Swap(i, j int)      { d[i], d[j] = d[j], d[i] }
func (d JudgeSlice) Less(i, j int) bool { return d[i].Progress < d[j].Progress }

// Segment dm segment struct
type Segment struct {
	Start    int64 `json:"ps"`       // 分段起始时间
	End      int64 `json:"pe"`       // 分段结束时间
	Cnt      int64 `json:"cnt"`      // 总分段数
	Num      int64 `json:"num"`      // 当前第几段
	Duration int64 `json:"duration"` // 视频总时长
}

// DMSegResp segment dm list response
type DMSegResp struct {
	Dms  []*Elem         `json:"dms"`
	Flag json.RawMessage `json:"flags,omitempty"`
}

// Page page info.
type Page struct {
	Num   int64 `json:"num"`
	Size  int64 `json:"size"`
	Total int64 `json:"total"`
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
	cnt = duration / segmentLength
	if duration%segmentLength > 0 {
		cnt++
	}
	for i := int64(0); i < cnt; i++ {
		if ps >= i*segmentLength && ps < (i+1)*segmentLength {
			ps = i * segmentLength
			pe = (i + 1) * segmentLength
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

// Encode dm ecode.
func Encode(flag, xml []byte) (res []byte) {
	var (
		fl = uint32(len(flag))
		xl = uint32(len(xml))
	)
	res = make([]byte, 4+fl+xl)
	binary.BigEndian.PutUint32(res[0:4], fl)
	copy(res[4:], flag)
	copy(res[4+fl:], xml)
	return
}

// Decode decode dm proto.
func Decode(buf []byte) (flag, xml []byte, err error) {
	var (
		zr *gzip.Reader
	)
	fl := binary.BigEndian.Uint32(buf[0:4])
	flag = buf[4 : 4+fl]
	if zr, err = gzip.NewReader(bytes.NewBuffer(buf[4+fl:])); err != nil {
		return
	}
	zr.Close()
	if xml, err = ioutil.ReadAll(zr); err != nil {
		return
	}
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
		Attribute: fmt.Sprintf(`"%d,%d,%d,%d,%d,%d,%d,%d,%s"`, d.ID, d.Attr, d.Progress, d.Content.Mode, d.Content.FontSize, d.Content.Color, d.Ctime, d.Pool, Hash(d.Mid, uint32(d.Content.IP))),
		Content:   msg,
	}
	return
}
