package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"strings"
	"time"

	"go-common/app/infra/canal/conf"
	"go-common/app/infra/canal/infoc"
	"go-common/app/infra/canal/model"
	"go-common/library/log"
	"go-common/library/queue/databus"

	"github.com/pkg/errors"
	"github.com/siddontang/go-mysql/canal"
)

var (
	errInvalidAction = errors.New("invalid rows action")
	errInvalidUpdate = errors.New("invalid update rows event")
	errBinlogFormat  = errors.New("binlog format failed")
)

type producer interface {
	Rows(int64)
	Send(context.Context, string, interface{}) error
	Close()
	Name() string
}

type databusP struct {
	group, topic string
	*databus.Databus
}

func (d *databusP) Rows(b int64) {
	// ignore
}

func (d *databusP) Send(c context.Context, key string, data interface{}) error {
	return d.Databus.Send(c, key, data)
}

func (d *databusP) Name() string {
	return fmt.Sprintf("databus:group(%s)topic(%s)", d.group, d.topic)
}

func (d *databusP) Close() {
	d.Databus.Close()
}

// infocP infoc producer
type infocP struct {
	taskID string
	*infoc.Infoc
}

// Rows rows
func (i *infocP) Rows(b int64) {
	i.Infoc.Rows(b)
}

// Send send msg
func (i *infocP) Send(c context.Context, key string, data interface{}) error {
	return i.Infoc.Send(c, key, data)
}

// Name infoc name
func (i *infocP) Name() string {
	return fmt.Sprintf("infoc(%s)", i.taskID)
}

// Close close infoc
func (i *infocP) Close() {
	i.Infoc.Flush()
	i.Infoc.Close()
}

// Target databus target
type Target struct {
	producers []producer
	eventLen  uint32
	events    []chan *canal.RowsEvent
	db        *conf.Database

	closed bool
}

// NewTarget new databus target
func NewTarget(db *conf.Database) (t *Target) {
	t = &Target{
		db:       db,
		eventLen: uint32(len(db.CTables)),
	}
	t.events = make([]chan *canal.RowsEvent, t.eventLen)
	if db.Databus != nil {
		t.producers = append(t.producers, &databusP{group: db.Databus.Group, topic: db.Databus.Topic, Databus: databus.New(db.Databus)})
	}
	if db.Infoc != nil {
		t.producers = append(t.producers, &infocP{taskID: db.Infoc.TaskID, Infoc: infoc.New(db.Infoc)})
	}
	for i := 0; i < int(t.eventLen); i++ {
		ch := make(chan *canal.RowsEvent, 1024)
		t.events[i] = ch
		go t.proc(ch)
	}
	return
}

// compare check if the binlog event is needed
// check the table name and schame
func (t *Target) compare(schame, table, action string) bool {
	if t.db.Schema == schame {
		for _, ctb := range t.db.CTables {
			for _, tb := range ctb.Tables {
				if table == tb {
					for _, act := range ctb.OmitAction {
						if act == action { // NOTE: omit action
							return false
						}
					}
					return true
				}
			}
		}
	}
	return false
}

// send send rows event into event chans
// and hash by table%concurrency.
func (t *Target) send(ev *canal.RowsEvent) {
	yu := crc32.ChecksumIEEE([]byte(ev.Table.Name))
	t.events[yu%t.eventLen] <- ev
}

func (t *Target) close() {
	for _, p := range t.producers {
		p.Close()
	}
	t.closed = true
}

// proc aync method for transfer the binlog data
// when connection is bad, just refresh it with retry
func (t *Target) proc(ch chan *canal.RowsEvent) {
	type pData struct {
		datas    []*model.Data
		producer producer
	}
	var (
		err         error
		normalDatas []*pData
		errorDatas  []*pData
		ev          *canal.RowsEvent
	)
	for {
		if t.closed {
			return
		}
		if len(errorDatas) != 0 {
			normalDatas = errorDatas
			errorDatas = errorDatas[0:0]
			time.Sleep(time.Second)
		} else {
			ev = <-ch
			var datas []*model.Data
			if datas, err = makeDatas(ev, t.db.TableMap); err != nil {
				log.Error("makeData(%v) error(%v)", ev, err)
				continue
			}
			normalDatas = normalDatas[0:0]
			for _, p := range t.producers {
				p.Rows(int64(len(datas)))
				normalDatas = append(normalDatas, &pData{datas: datas, producer: p})
				if stats != nil {
					stats.Incr("send_counter", p.Name(), ev.Table.Schema, tblReplacer.ReplaceAllString(ev.Table.Name, ""), ev.Action)
				}
			}
		}
		for _, pd := range normalDatas {
			var eDatas []*model.Data
			for _, data := range pd.datas {
				if err = pd.producer.Send(context.TODO(), data.Key, data); err != nil {
					// retry pub error data
					eDatas = append(eDatas, data)
					continue
				}
				log.Info("%s pub(key:%s, value:%+v) succeed", pd.producer.Name(), data.Key, data)
			}
			if len(eDatas) > 0 {
				errorDatas = append(errorDatas, &pData{datas: eDatas, producer: pd.producer})
				if stats != nil && ev != nil {
					stats.Incr("retry_counter", pd.producer.Name(), ev.Table.Schema, tblReplacer.ReplaceAllString(ev.Table.Name, ""), ev.Action)
				}
				log.Error("%s scheme(%s) pub fail,add to retry", pd.producer.Name(), ev.Table.Schema)
			}
		}
	}
}

// makeDatas parse the binlog event and return the model.Data struct
// a little bit cautious about the binlog type
// if the type is update:
//     the old value and new value will alternate appearing in the event.Rows
func makeDatas(e *canal.RowsEvent, tbMap map[string]*conf.Addition) (datas []*model.Data, err error) {
	var (
		rowsLen     = len(e.Rows)
		firstRowLen = len(e.Rows[0])
		lenCol      = len(e.Table.Columns)
	)
	if rowsLen == 0 || firstRowLen == 0 || firstRowLen != lenCol {
		log.Error("rows length(%d) first row length(%d) columns length(%d)", rowsLen, firstRowLen, lenCol)
		err = errBinlogFormat
		return
	}
	datas = make([]*model.Data, 0, rowsLen)
	switch e.Action {
	case canal.InsertAction, canal.DeleteAction:
		for _, values := range e.Rows {
			var keys []string
			data := &model.Data{
				Action: e.Action,
				Table:  e.Table.Name,
				// the first primary key as the kafka key
				Key: fmt.Sprint(values[0]),
				New: make(map[string]interface{}, lenCol),
			}
			for i, c := range e.Table.Columns {
				if c.IsUnsigned {
					values[i] = unsignIntCase(values[i])
				}
				if strings.Contains(c.RawType, "binary") {
					if bs, ok := values[i].(string); ok {
						values[i] = base64.StdEncoding.EncodeToString([]byte(bs))
					}
				}
				data.New[c.Name] = values[i]
			}
			// set kafka key and remove omit columns data
			addition, ok := tbMap[e.Table.Name]
			if ok {
				for _, omit := range addition.OmitField {
					delete(data.New, omit)
				}
				for _, primary := range addition.PrimaryKey {
					if _, ok := data.New[primary]; ok {
						keys = append(keys, fmt.Sprint(data.New[primary]))
					}
				}
			}
			if len(keys) != 0 {
				data.Key = strings.Join(keys, ",")
			}
			datas = append(datas, data)
		}
	case canal.UpdateAction:
		if rowsLen%2 != 0 {
			err = errInvalidUpdate
			return
		}
		for i := 0; i < rowsLen; i += 2 {
			var keys []string
			data := &model.Data{
				Action: e.Action,
				Table:  e.Table.Name,
				// the first primary key as the kafka key
				Key: fmt.Sprint(e.Rows[i][0]),
				Old: make(map[string]interface{}, lenCol),
				New: make(map[string]interface{}, lenCol),
			}
			for j, c := range e.Table.Columns {
				if c.IsUnsigned {
					e.Rows[i][j] = unsignIntCase(e.Rows[i][j])
					e.Rows[i+1][j] = unsignIntCase(e.Rows[i+1][j])
				}
				if strings.Contains(c.RawType, "binary") {
					if bs, ok := e.Rows[i][j].(string); ok {
						e.Rows[i][j] = base64.StdEncoding.EncodeToString([]byte(bs))
					}
					if bs, ok := e.Rows[i+1][j].(string); ok {
						e.Rows[i+1][j] = base64.StdEncoding.EncodeToString([]byte(bs))
					}
				}
				data.Old[c.Name] = e.Rows[i][j]
				data.New[c.Name] = e.Rows[i+1][j]
			}
			// set kafka key and remove omit columns data
			addition, ok := tbMap[e.Table.Name]
			if ok {
				for _, omit := range addition.OmitField {
					delete(data.New, omit)
					delete(data.Old, omit)
				}
				for _, primary := range addition.PrimaryKey {
					if _, ok := data.New[primary]; ok {
						keys = append(keys, fmt.Sprint(data.New[primary]))
					}
				}
			}
			if len(keys) != 0 {
				data.Key = strings.Join(keys, ",")
			}
			datas = append(datas, data)
		}
	default:
		err = errInvalidAction
	}
	return
}

func unsignIntCase(i interface{}) (v interface{}) {
	switch si := i.(type) {
	case int8:
		v = uint8(si)
	case int16:
		v = uint16(si)
	case int32:
		v = uint32(si)
	case int64:
		v = uint64(si)
	default:
		v = i
	}
	return
}
