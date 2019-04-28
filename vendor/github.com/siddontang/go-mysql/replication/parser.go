package replication

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/juju/errors"
)

type BinlogParser struct {
	format *FormatDescriptionEvent

	tables map[uint64]*TableMapEvent

	// for rawMode, we only parse FormatDescriptionEvent and RotateEvent
	rawMode bool

	parseTime bool

	// used to start/stop processing
	stopProcessing uint32

	useDecimal bool
}

func NewBinlogParser() *BinlogParser {
	p := new(BinlogParser)

	p.tables = make(map[uint64]*TableMapEvent)
	// p.stop = make(uint32)

	return p
}

func (p *BinlogParser) Stop() {
	atomic.StoreUint32(&p.stopProcessing, 1)
}

func (p *BinlogParser) Resume() {
	atomic.StoreUint32(&p.stopProcessing, 0)
}

func (p *BinlogParser) Reset() {
	p.format = nil
}

type OnEventFunc func(*BinlogEvent) error

func (p *BinlogParser) ParseFile(name string, offset int64, onEvent OnEventFunc) error {
	f, err := os.Open(name)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	b := make([]byte, 4)
	if _, err = f.Read(b); err != nil {
		return errors.Trace(err)
	} else if !bytes.Equal(b, BinLogFileHeader) {
		return errors.Errorf("%s is not a valid binlog file, head 4 bytes must fe'bin' ", name)
	}

	if offset < 4 {
		offset = 4
	} else if offset > 4 {
		//  FORMAT_DESCRIPTION event should be read by default always (despite that fact passed offset may be higher than 4)
		if _, err = f.Seek(4, os.SEEK_SET); err != nil {
			return errors.Errorf("seek %s to %d error %v", name, offset, err)
		}

		p.getFormatDescriptionEvent(f, onEvent)
	}

	if _, err = f.Seek(offset, os.SEEK_SET); err != nil {
		return errors.Errorf("seek %s to %d error %v", name, offset, err)
	}

	return p.ParseReader(f, onEvent)
}

func (p *BinlogParser) getFormatDescriptionEvent(r io.Reader, onEvent OnEventFunc) error {
	_, err := p.parseSingleEvent(&r, onEvent)
	return err
}

func (p *BinlogParser) parseSingleEvent(r *io.Reader, onEvent OnEventFunc) (bool, error) {
	var err error
	var n int64

	headBuf := make([]byte, EventHeaderSize)

	if _, err = io.ReadFull(*r, headBuf); err == io.EOF {
		return true, nil
	} else if err != nil {
		return false, errors.Trace(err)
	}

	var h *EventHeader
	h, err = p.parseHeader(headBuf)
	if err != nil {
		return false, errors.Trace(err)
	}

	if h.EventSize <= uint32(EventHeaderSize) {
		return false, errors.Errorf("invalid event header, event size is %d, too small", h.EventSize)
	}

	var buf bytes.Buffer
	if n, err = io.CopyN(&buf, *r, int64(h.EventSize)-int64(EventHeaderSize)); err != nil {
		return false, errors.Errorf("get event body err %v, need %d - %d, but got %d", err, h.EventSize, EventHeaderSize, n)
	}

	data := buf.Bytes()
	rawData := data

	eventLen := int(h.EventSize) - EventHeaderSize

	if len(data) != eventLen {
		return false, errors.Errorf("invalid data size %d in event %s, less event length %d", len(data), h.EventType, eventLen)
	}

	var e Event
	e, err = p.parseEvent(h, data)
	if err != nil {
		if _, ok := err.(errMissingTableMapEvent); ok {
			return false, nil
		}
		return false, errors.Trace(err)
	}

	if err = onEvent(&BinlogEvent{rawData, h, e}); err != nil {
		return false, errors.Trace(err)
	}

	return false, nil
}

func (p *BinlogParser) ParseReader(r io.Reader, onEvent OnEventFunc) error {

	for {
		if atomic.LoadUint32(&p.stopProcessing) == 1 {
			break
		}

		done, err := p.parseSingleEvent(&r, onEvent)
		if err != nil {
			if _, ok := err.(errMissingTableMapEvent); ok {
				continue
			}
			return errors.Trace(err)
		}

		if done {
			break
		}
	}

	return nil
}

func (p *BinlogParser) SetRawMode(mode bool) {
	p.rawMode = mode
}

func (p *BinlogParser) SetParseTime(parseTime bool) {
	p.parseTime = parseTime
}

func (p *BinlogParser) SetUseDecimal(useDecimal bool) {
	p.useDecimal = useDecimal
}

func (p *BinlogParser) parseHeader(data []byte) (*EventHeader, error) {
	h := new(EventHeader)
	err := h.Decode(data)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (p *BinlogParser) parseEvent(h *EventHeader, data []byte) (Event, error) {
	var e Event

	if h.EventType == FORMAT_DESCRIPTION_EVENT {
		p.format = &FormatDescriptionEvent{}
		e = p.format
	} else {
		if p.format != nil && p.format.ChecksumAlgorithm == BINLOG_CHECKSUM_ALG_CRC32 {
			data = data[0 : len(data)-4]
		}

		if h.EventType == ROTATE_EVENT {
			e = &RotateEvent{}
		} else if !p.rawMode {
			switch h.EventType {
			case QUERY_EVENT:
				e = &QueryEvent{}
			case XID_EVENT:
				e = &XIDEvent{}
			case TABLE_MAP_EVENT:
				te := &TableMapEvent{}
				if p.format.EventTypeHeaderLengths[TABLE_MAP_EVENT-1] == 6 {
					te.tableIDSize = 4
				} else {
					te.tableIDSize = 6
				}
				e = te
			case WRITE_ROWS_EVENTv0,
				UPDATE_ROWS_EVENTv0,
				DELETE_ROWS_EVENTv0,
				WRITE_ROWS_EVENTv1,
				DELETE_ROWS_EVENTv1,
				UPDATE_ROWS_EVENTv1,
				WRITE_ROWS_EVENTv2,
				UPDATE_ROWS_EVENTv2,
				DELETE_ROWS_EVENTv2:
				e = p.newRowsEvent(h)
			case ROWS_QUERY_EVENT:
				e = &RowsQueryEvent{}
			case GTID_EVENT:
				e = &GTIDEvent{}
			case BEGIN_LOAD_QUERY_EVENT:
				e = &BeginLoadQueryEvent{}
			case EXECUTE_LOAD_QUERY_EVENT:
				e = &ExecuteLoadQueryEvent{}
			case MARIADB_ANNOTATE_ROWS_EVENT:
				e = &MariadbAnnotateRowsEvent{}
			case MARIADB_BINLOG_CHECKPOINT_EVENT:
				e = &MariadbBinlogCheckPointEvent{}
			case MARIADB_GTID_LIST_EVENT:
				e = &MariadbGTIDListEvent{}
			case MARIADB_GTID_EVENT:
				ee := &MariadbGTIDEvent{}
				ee.GTID.ServerID = h.ServerID
				e = ee
			default:
				e = &GenericEvent{}
			}
		} else {
			e = &GenericEvent{}
		}
	}

	if err := e.Decode(data); err != nil {
		return nil, &EventError{h, err.Error(), data}
	}

	if te, ok := e.(*TableMapEvent); ok {
		p.tables[te.TableID] = te
	}

	if re, ok := e.(*RowsEvent); ok {
		if (re.Flags & RowsEventStmtEndFlag) > 0 {
			// Refer https://github.com/alibaba/canal/blob/38cc81b7dab29b51371096fb6763ca3a8432ffee/dbsync/src/main/java/com/taobao/tddl/dbsync/binlog/event/RowsLogEvent.java#L176
			p.tables = make(map[uint64]*TableMapEvent)
		}
	}

	return e, nil
}

// Given the bytes for a a binary log event: return the decoded event.
// With the exception of the FORMAT_DESCRIPTION_EVENT event type
// there must have previously been passed a FORMAT_DESCRIPTION_EVENT
// into the parser for this to work properly on any given event.
// Passing a new FORMAT_DESCRIPTION_EVENT into the parser will replace
// an existing one.
func (p *BinlogParser) Parse(data []byte) (*BinlogEvent, error) {
	rawData := data

	h, err := p.parseHeader(data)

	if err != nil {
		return nil, err
	}

	data = data[EventHeaderSize:]
	eventLen := int(h.EventSize) - EventHeaderSize

	if len(data) != eventLen {
		return nil, fmt.Errorf("invalid data size %d in event %s, less event length %d", len(data), h.EventType, eventLen)
	}

	e, err := p.parseEvent(h, data)
	if err != nil {
		return nil, err
	}

	return &BinlogEvent{rawData, h, e}, nil
}

func (p *BinlogParser) newRowsEvent(h *EventHeader) *RowsEvent {
	e := &RowsEvent{}
	if p.format.EventTypeHeaderLengths[h.EventType-1] == 6 {
		e.tableIDSize = 4
	} else {
		e.tableIDSize = 6
	}

	e.needBitmap2 = false
	e.tables = p.tables
	e.parseTime = p.parseTime
	e.useDecimal = p.useDecimal

	switch h.EventType {
	case WRITE_ROWS_EVENTv0:
		e.Version = 0
	case UPDATE_ROWS_EVENTv0:
		e.Version = 0
	case DELETE_ROWS_EVENTv0:
		e.Version = 0
	case WRITE_ROWS_EVENTv1:
		e.Version = 1
	case DELETE_ROWS_EVENTv1:
		e.Version = 1
	case UPDATE_ROWS_EVENTv1:
		e.Version = 1
		e.needBitmap2 = true
	case WRITE_ROWS_EVENTv2:
		e.Version = 2
	case UPDATE_ROWS_EVENTv2:
		e.Version = 2
		e.needBitmap2 = true
	case DELETE_ROWS_EVENTv2:
		e.Version = 2
	}

	return e
}
