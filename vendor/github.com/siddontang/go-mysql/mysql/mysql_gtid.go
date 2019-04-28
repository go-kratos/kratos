package mysql

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/satori/go.uuid"
	"github.com/siddontang/go/hack"
)

// Like MySQL GTID Interval struct, [start, stop), left closed and right open
// See MySQL rpl_gtid.h
type Interval struct {
	// The first GID of this interval.
	Start int64
	// The first GID after this interval.
	Stop int64
}

// Interval is [start, stop), but the GTID string's format is [n] or [n1-n2], closed interval
func parseInterval(str string) (i Interval, err error) {
	p := strings.Split(str, "-")
	switch len(p) {
	case 1:
		i.Start, err = strconv.ParseInt(p[0], 10, 64)
		i.Stop = i.Start + 1
	case 2:
		i.Start, err = strconv.ParseInt(p[0], 10, 64)
		i.Stop, err = strconv.ParseInt(p[1], 10, 64)
		i.Stop = i.Stop + 1
	default:
		err = errors.Errorf("invalid interval format, must n[-n]")
	}

	if err != nil {
		return
	}

	if i.Stop <= i.Start {
		err = errors.Errorf("invalid interval format, must n[-n] and the end must >= start")
	}

	return
}

func (i Interval) String() string {
	if i.Stop == i.Start+1 {
		return fmt.Sprintf("%d", i.Start)
	} else {
		return fmt.Sprintf("%d-%d", i.Start, i.Stop-1)
	}
}

type IntervalSlice []Interval

func (s IntervalSlice) Len() int {
	return len(s)
}

func (s IntervalSlice) Less(i, j int) bool {
	if s[i].Start < s[j].Start {
		return true
	} else if s[i].Start > s[j].Start {
		return false
	} else {
		return s[i].Stop < s[j].Stop
	}
}

func (s IntervalSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s IntervalSlice) Sort() {
	sort.Sort(s)
}

func (s IntervalSlice) Normalize() IntervalSlice {
	var n IntervalSlice
	if len(s) == 0 {
		return n
	}

	s.Sort()

	n = append(n, s[0])

	for i := 1; i < len(s); i++ {
		last := n[len(n)-1]
		if s[i].Start > last.Stop {
			n = append(n, s[i])
			continue
		} else {
			stop := s[i].Stop
			if last.Stop > stop {
				stop = last.Stop
			}
			n[len(n)-1] = Interval{last.Start, stop}
		}
	}

	return n
}

// Return true if sub in s
func (s IntervalSlice) Contain(sub IntervalSlice) bool {
	j := 0
	for i := 0; i < len(sub); i++ {
		for ; j < len(s); j++ {
			if sub[i].Start > s[j].Stop {
				continue
			} else {
				break
			}
		}
		if j == len(s) {
			return false
		}

		if sub[i].Start < s[j].Start || sub[i].Stop > s[j].Stop {
			return false
		}
	}

	return true
}

func (s IntervalSlice) Equal(o IntervalSlice) bool {
	if len(s) != len(o) {
		return false
	}

	for i := 0; i < len(s); i++ {
		if s[i].Start != o[i].Start || s[i].Stop != o[i].Stop {
			return false
		}
	}

	return true
}

func (s IntervalSlice) Compare(o IntervalSlice) int {
	if s.Equal(o) {
		return 0
	} else if s.Contain(o) {
		return 1
	} else {
		return -1
	}
}

// Refer http://dev.mysql.com/doc/refman/5.6/en/replication-gtids-concepts.html
type UUIDSet struct {
	SID uuid.UUID

	Intervals IntervalSlice
}

func ParseUUIDSet(str string) (*UUIDSet, error) {
	str = strings.TrimSpace(str)
	sep := strings.Split(str, ":")
	if len(sep) < 2 {
		return nil, errors.Errorf("invalid GTID format, must UUID:interval[:interval]")
	}

	var err error
	s := new(UUIDSet)
	if s.SID, err = uuid.FromString(sep[0]); err != nil {
		return nil, errors.Trace(err)
	}

	// Handle interval
	for i := 1; i < len(sep); i++ {
		if in, err := parseInterval(sep[i]); err != nil {
			return nil, errors.Trace(err)
		} else {
			s.Intervals = append(s.Intervals, in)
		}
	}

	s.Intervals = s.Intervals.Normalize()

	return s, nil
}

func NewUUIDSet(sid uuid.UUID, in ...Interval) *UUIDSet {
	s := new(UUIDSet)
	s.SID = sid

	s.Intervals = in
	s.Intervals = s.Intervals.Normalize()

	return s
}

func (s *UUIDSet) Contain(sub *UUIDSet) bool {
	if !bytes.Equal(s.SID.Bytes(), sub.SID.Bytes()) {
		return false
	}

	return s.Intervals.Contain(sub.Intervals)
}

func (s *UUIDSet) Bytes() []byte {
	var buf bytes.Buffer

	buf.WriteString(s.SID.String())

	for _, i := range s.Intervals {
		buf.WriteString(":")
		buf.WriteString(i.String())
	}

	return buf.Bytes()
}

func (s *UUIDSet) AddInterval(in IntervalSlice) {
	s.Intervals = append(s.Intervals, in...)
	s.Intervals = s.Intervals.Normalize()
}

func (s *UUIDSet) String() string {
	return hack.String(s.Bytes())
}

func (s *UUIDSet) encode(w io.Writer) {
	w.Write(s.SID.Bytes())
	n := int64(len(s.Intervals))

	binary.Write(w, binary.LittleEndian, n)

	for _, i := range s.Intervals {
		binary.Write(w, binary.LittleEndian, i.Start)
		binary.Write(w, binary.LittleEndian, i.Stop)
	}
}

func (s *UUIDSet) Encode() []byte {
	var buf bytes.Buffer

	s.encode(&buf)

	return buf.Bytes()
}

func (s *UUIDSet) decode(data []byte) (int, error) {
	if len(data) < 24 {
		return 0, errors.Errorf("invalid uuid set buffer, less 24")
	}

	pos := 0
	var err error
	if s.SID, err = uuid.FromBytes(data[0:16]); err != nil {
		return 0, err
	}
	pos += 16

	n := int64(binary.LittleEndian.Uint64(data[pos : pos+8]))
	pos += 8
	if len(data) < int(16*n)+pos {
		return 0, errors.Errorf("invalid uuid set buffer, must %d, but %d", pos+int(16*n), len(data))
	}

	s.Intervals = make([]Interval, 0, n)

	var in Interval
	for i := int64(0); i < n; i++ {
		in.Start = int64(binary.LittleEndian.Uint64(data[pos : pos+8]))
		pos += 8
		in.Stop = int64(binary.LittleEndian.Uint64(data[pos : pos+8]))
		pos += 8
		s.Intervals = append(s.Intervals, in)
	}

	return pos, nil
}

func (s *UUIDSet) Decode(data []byte) error {
	n, err := s.decode(data)
	if n != len(data) {
		return errors.Errorf("invalid uuid set buffer, must %d, but %d", n, len(data))
	}
	return err
}

type MysqlGTIDSet struct {
	Sets map[string]*UUIDSet
}

func ParseMysqlGTIDSet(str string) (GTIDSet, error) {
	s := new(MysqlGTIDSet)
	s.Sets = make(map[string]*UUIDSet)
	if str == "" {
		return s, nil
	}

	sp := strings.Split(str, ",")

	//todo, handle redundant same uuid
	for i := 0; i < len(sp); i++ {
		if set, err := ParseUUIDSet(sp[i]); err != nil {
			return nil, errors.Trace(err)
		} else {
			s.AddSet(set)
		}

	}
	return s, nil
}

func DecodeMysqlGTIDSet(data []byte) (*MysqlGTIDSet, error) {
	s := new(MysqlGTIDSet)

	if len(data) < 8 {
		return nil, errors.Errorf("invalid gtid set buffer, less 4")
	}

	n := int(binary.LittleEndian.Uint64(data))
	s.Sets = make(map[string]*UUIDSet, n)

	pos := 8

	for i := 0; i < n; i++ {
		set := new(UUIDSet)
		if n, err := set.decode(data[pos:]); err != nil {
			return nil, errors.Trace(err)
		} else {
			pos += n

			s.AddSet(set)
		}
	}
	return s, nil
}

func (s *MysqlGTIDSet) AddSet(set *UUIDSet) {
	if set == nil {
		return
	}
	sid := set.SID.String()
	o, ok := s.Sets[sid]
	if ok {
		o.AddInterval(set.Intervals)
	} else {
		s.Sets[sid] = set
	}
}

func (s *MysqlGTIDSet) Update(GTIDStr string) error {
	uuidSet, err := ParseUUIDSet(GTIDStr)
	if err != nil {
		return err
	}

	s.AddSet(uuidSet)

	return nil
}

func (s *MysqlGTIDSet) Contain(o GTIDSet) bool {
	sub, ok := o.(*MysqlGTIDSet)
	if !ok {
		return false
	}

	for key, set := range sub.Sets {
		o, ok := s.Sets[key]
		if !ok {
			return false
		}

		if !o.Contain(set) {
			return false
		}
	}

	return true
}

func (s *MysqlGTIDSet) Equal(o GTIDSet) bool {
	sub, ok := o.(*MysqlGTIDSet)
	if !ok {
		return false
	}

	for key, set := range sub.Sets {
		o, ok := s.Sets[key]
		if !ok {
			return false
		}

		if !o.Intervals.Equal(set.Intervals) {
			return false
		}
	}

	return true

}

func (s *MysqlGTIDSet) String() string {
	var buf bytes.Buffer
	sep := ""
	for _, set := range s.Sets {
		buf.WriteString(sep)
		buf.WriteString(set.String())
		sep = ","
	}

	return hack.String(buf.Bytes())
}

func (s *MysqlGTIDSet) Encode() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.LittleEndian, uint64(len(s.Sets)))

	for i, _ := range s.Sets {
		s.Sets[i].encode(&buf)
	}

	return buf.Bytes()
}
