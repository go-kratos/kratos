package reply

import (
	"context"
	"strconv"

	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/queue/databus"
)

// DatabusDao DatabusDao
type DatabusDao struct {
	topic   string
	databus *databus.Databus
}

type kafkadata struct {
	Op      string `json:"op,omitempty"`
	Mid     int64  `json:"mid,omitempty"`
	Adid    int64  `json:"adid,omitempty"`
	Oid     int64  `json:"oid,omitempty"`
	Rpid    int64  `json:"rpid,omitempty"`
	Root    int64  `json:"root,omitempty"`
	Dialog  int64  `json:"dialog,omitempty"`
	Remark  string `json:"remark,omitempty"`
	Adname  string `json:"adname,omitempty"`
	Mtime   int64  `json:"mtime,omitempty"`
	Action  int8   `json:"action,omitempty"`
	Sort    int8   `json:"sort,omitempty"`
	Tp      int8   `json:"tp,omitempty"`
	Moral   int    `json:"moral,omitempty"`
	Notify  bool   `json:"notify,omitempty"`
	Top     uint32 `json:"top,omitempty"`
	Ftime   int64  `json:"ftime,omitempty"`
	State   int8   `json:"state,omitempty"`
	Audit   int8   `json:"audit,omitempty"`
	Reason  int8   `json:"reason,omitempty"`
	Content string `json:"content,omitempty"`
	FReason int8   `json:"freason,omitempty"`
	Assist  bool   `json:"assist,omitempty"`
	Count   int    `json:"count,omitempty"`
	Floor   int    `json:"floor,omitempty"`
	IsUp    bool   `json:"is_up,omitempty"`
}

// NewDatabusDao new ReplyKafkaDao and return.
func NewDatabusDao(c *databus.Config) (dao *DatabusDao) {
	dao = &DatabusDao{
		topic:   c.Topic,
		databus: databus.New(c),
	}
	return
}

// PubEvent pub reply event.
func (dao *DatabusDao) push(c context.Context, key string, value interface{}) error {
	return dao.databus.Send(c, key, value)
}

// RecoverFixDialogIdx ...
func (dao *DatabusDao) RecoverFixDialogIdx(c context.Context, oid int64, tp int8, root int64) {
	var message = map[string]interface{}{}
	message["action"] = "fix_dialog"
	message["data"] = kafkadata{
		Oid:  oid,
		Tp:   tp,
		Root: root,
	}
	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// RecoverFolderIdx ...
func (dao *DatabusDao) RecoverFolderIdx(c context.Context, oid int64, tp int8, root int64) {
	var message = map[string]interface{}{}
	message["action"] = "folder"
	message["data"] = kafkadata{
		Op:   "re_idx",
		Oid:  oid,
		Tp:   tp,
		Root: root,
	}
	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// RecoverDialogIdx Recover dialog index
func (dao *DatabusDao) RecoverDialogIdx(c context.Context, oid int64, tp int8, root, dialog int64) {
	var message = map[string]interface{}{}
	message["action"] = "idx_dialog"
	message["data"] = kafkadata{
		Oid:    oid,
		Tp:     tp,
		Root:   root,
		Dialog: dialog,
	}
	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// RecoverFloorIdx RecoverFloorIdx
func (dao *DatabusDao) RecoverFloorIdx(c context.Context, oid int64, tp int8, num int, isFloor bool) {
	var (
		message = map[string]interface{}{}
	)
	message["action"] = "idx_floor"
	if isFloor {
		message["data"] = kafkadata{
			Oid:   oid,
			Tp:    tp,
			Floor: num,
		}
	} else {
		message["data"] = kafkadata{
			Oid:   oid,
			Tp:    tp,
			Count: num,
		}
	}
	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AddTop AddTop
func (dao *DatabusDao) AddTop(c context.Context, oid int64, tp int8, top uint32) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "add_top"
	message["data"] = kafkadata{
		Oid: oid,
		Tp:  tp,
		Top: top,
	}
	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AddReply push event message into kafka.
func (dao *DatabusDao) AddReply(c context.Context, oid int64, rp *reply.Reply) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "add"
	message["data"] = rp
	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AddSpam push event message into kafka.
func (dao *DatabusDao) AddSpam(c context.Context, oid, mid int64, isUp bool, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "spam"
	data := kafkadata{
		Mid:  mid,
		IsUp: isUp,
		Tp:   tp,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AddReport push event message into kafka.
func (dao *DatabusDao) AddReport(c context.Context, oid, rpID int64, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "rpt"
	data := kafkadata{
		Oid:  oid,
		Rpid: rpID,
		Tp:   tp,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// Like push event message into kafka.
func (dao *DatabusDao) Like(c context.Context, oid, rpID, mid int64, action int8, ts int64) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "act"
	data := kafkadata{
		Oid:    oid,
		Mid:    mid,
		Rpid:   rpID,
		Action: action,
		Op:     "like",
		Mtime:  ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// Hate push event message into kafka.
func (dao *DatabusDao) Hate(c context.Context, oid, rpID, mid int64, action int8, ts int64) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "act"
	data := kafkadata{
		Oid:    oid,
		Mid:    mid,
		Rpid:   rpID,
		Action: action,
		Op:     "hate",
		Mtime:  ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// RecoverIndex push event message into kafka.
func (dao *DatabusDao) RecoverIndex(c context.Context, oid int64, tp, sort int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "re_idx"
	data := kafkadata{
		Oid:  oid,
		Tp:   tp,
		Sort: sort,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// RecoverIndexByRoot push event message into kafka.
func (dao *DatabusDao) RecoverIndexByRoot(c context.Context, oid, root int64, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "re_rt_idx"
	data := kafkadata{
		Oid:  oid,
		Tp:   tp,
		Root: root,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// Hide push event message into kafka.
func (dao *DatabusDao) Hide(c context.Context, oid, rpID int64, tp int8, ts int64) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "up"
	data := kafkadata{
		Oid:   oid,
		Rpid:  rpID,
		Mtime: ts,
		Tp:    tp,
		Op:    "hide",
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// Show push event message into kafka.
func (dao *DatabusDao) Show(c context.Context, oid, rpID int64, tp int8, ts int64) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "up"
	data := kafkadata{
		Op:    "show",
		Oid:   oid,
		Tp:    tp,
		Rpid:  rpID,
		Mtime: ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// Delete push event message into kafka.
func (dao *DatabusDao) Delete(c context.Context, mid, oid, rpID int64, ts int64, tp int8, assist bool) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:     "del_up",
		Oid:    oid,
		Mid:    mid,
		Rpid:   rpID,
		Tp:     tp,
		Mtime:  ts,
		Assist: assist,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminEdit push event message into kafka.
func (dao *DatabusDao) AdminEdit(c context.Context, oid, rpID int64, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:   "edit",
		Oid:  oid,
		Tp:   tp,
		Rpid: rpID,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminAddTop push event message into kafka.
func (dao *DatabusDao) AdminAddTop(c context.Context, adid, oid, rpID, ts int64, act, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:     "top_add",
		Oid:    oid,
		Adid:   adid,
		Rpid:   rpID,
		Tp:     tp,
		Action: act,
		Mtime:  ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// UpperAddTop push event message into kafka.
func (dao *DatabusDao) UpperAddTop(c context.Context, mid, oid, rpID, ts int64, act, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "up"
	data := kafkadata{
		Op:     "top_add",
		Oid:    oid,
		Tp:     tp,
		Mid:    mid,
		Rpid:   rpID,
		Action: act,
		Mtime:  ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminDelete push event message into kafka.
func (dao *DatabusDao) AdminDelete(c context.Context, adid, oid, rpID, ftime int64, moral int, notify bool, adname, remark string, ts int64, tp, reason, freason int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:      "del",
		Adid:    adid,
		Rpid:    rpID,
		Oid:     oid,
		Moral:   moral,
		Notify:  notify,
		Tp:      tp,
		Adname:  adname,
		Remark:  remark,
		Mtime:   ts,
		Ftime:   ftime,
		Reason:  reason,
		FReason: freason,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminDeleteByReport push event message into kafka.
func (dao *DatabusDao) AdminDeleteByReport(c context.Context, adid, oid, rpID, mid, ftime int64, moral int, notify bool, adname, remark string, ts int64, tp, audit, reason int8, content string, freason int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:      "del_rpt",
		Adid:    adid,
		Oid:     oid,
		Rpid:    rpID,
		Mid:     mid,
		Moral:   moral,
		Tp:      tp,
		Notify:  notify,
		Adname:  adname,
		Remark:  remark,
		Mtime:   ts,
		Ftime:   ftime,
		Audit:   audit,
		Reason:  reason,
		Content: content,
		FReason: freason,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminRecover push event message into kafka.
func (dao *DatabusDao) AdminRecover(c context.Context, adid, oid, rpID int64, remark string, ts int64, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:     "re",
		Adid:   adid,
		Oid:    oid,
		Rpid:   rpID,
		Tp:     tp,
		Remark: remark,
		Mtime:  ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminPass push pass event message into kafka.
func (dao *DatabusDao) AdminPass(c context.Context, adid, oid, rpID int64, remark string, ts int64, tp int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:     "pass",
		Adid:   adid,
		Oid:    oid,
		Rpid:   rpID,
		Tp:     tp,
		Remark: remark,
		Mtime:  ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminStateSet push event message into kafka.
func (dao *DatabusDao) AdminStateSet(c context.Context, adid, oid, rpID, ts int64, tp, state int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:    "stateset",
		Adid:  adid,
		Oid:   oid,
		Rpid:  rpID,
		Tp:    tp,
		State: state,
		Mtime: ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminTransfer push event message into kafka.
func (dao *DatabusDao) AdminTransfer(c context.Context, adid, oid, rpID, ts int64, tp, audit int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	message["data"] = kafkadata{
		Op:    "transfer",
		Adid:  adid,
		Oid:   oid,
		Rpid:  rpID,
		Tp:    tp,
		Audit: audit,
		Mtime: ts,
	}

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminIgnore push event message into kafka.
func (dao *DatabusDao) AdminIgnore(c context.Context, adid, oid, rpID, ts int64, tp, audit int8) {

	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:    "ignore",
		Adid:  adid,
		Oid:   oid,
		Rpid:  rpID,
		Tp:    tp,
		Audit: audit,
		Mtime: ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}

// AdminReportRecover push event message into kafka.
func (dao *DatabusDao) AdminReportRecover(c context.Context, adid, oid, rpID int64, remark string, ts int64, tp, audit int8) {
	var (
		message = make(map[string]interface{})
	)
	message["action"] = "admin"
	data := kafkadata{
		Op:     "rpt_re",
		Adid:   adid,
		Oid:    oid,
		Rpid:   rpID,
		Tp:     tp,
		Audit:  audit,
		Remark: remark,
		Mtime:  ts,
	}
	message["data"] = data

	key := strconv.FormatInt(oid, 10)
	dao.push(c, key, message)
}
