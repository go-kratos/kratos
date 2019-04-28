package model

import (
	"encoding/json"
	// "fmt"
	// "time"

	"go-common/library/log"

	"github.com/pkg/errors"
)

// Trend enum
const (
	TrendHold = 0
	TrendUp   = 1
	TrendDown = 2
)

// Charge 充电，amount: 总充电数
func (e *RankElecPrepUPProto) Charge(payMID int64, amount int64, newElecer bool) {
	// 检查是否在榜，如果在榜直接更新
	ele := e.Find(payMID)
	if ele != nil {
		ele.Amount = amount
		e.update(ele)
		return
	}
	if newElecer {
		e.Count++
		e.CountUPTotalElec++
	}
	// 排行榜列表为最大长度 且 榜单末位的充电数 >= 新的充电数，返回
	if len(e.List) >= e.Size_ && e.List[len(e.List)-1].Amount >= amount {
		return
	}
	// 没有在榜且可以上榜则插入原榜
	newEle := &RankElecPrepElementProto{
		Rank:      -1,
		MID:       payMID,
		TrendType: TrendHold,
		Amount:    amount,
	}
	e.insert(newEle)
}

// UpdateMessage 更新留言
func (e *RankElecPrepUPProto) UpdateMessage(payMID int64, message string, hidden bool) {
	ele := e.Find(payMID)
	if ele != nil {
		ele.Message = &ElecMessageProto{
			Message: message,
			Hidden:  hidden,
		}
	}
}

// Find 获得榜单中payMID的排名信息，如果不存在则返回nil
func (e *RankElecPrepUPProto) Find(payMID int64) (ele *RankElecPrepElementProto) {
	for _, r := range e.List {
		if r.MID == payMID {
			ele = r
			return
		}
	}
	return nil
}

// 更新排名
func (e *RankElecPrepUPProto) update(ele *RankElecPrepElementProto) {
	for i := range e.List {
		if e.List[i] == nil {
			log.Error("ElecPrepUPRank: %s, index: %d, ele: %+v", e, i, ele)
		}
		if e.List[i].MID != ele.MID && e.List[i].Amount >= ele.Amount {
			continue
		}
		newRank := e.List[i].Rank
		if err := e.shift(i, ele.Rank-1); err != nil {
			log.Error("ElecPrepUPRank.update err: %+v, ele: %+v in rank: %s ", err, ele, e)
			return
		}
		ele.Rank = newRank
		e.List[i] = ele
		break
	}
}

// 插入新排名到榜单
func (e *RankElecPrepUPProto) insert(ele *RankElecPrepElementProto) {
	for i := range e.List {
		if e.List[i].Amount >= ele.Amount {
			continue
		}
		// 找到新排名的位置，并插入
		ele.Rank = e.List[i].Rank
		if len(e.List) >= e.Size_ { // 榜单已满
			if err := e.shift(i, len(e.List)-1); err != nil {
				log.Error("ElecPrepUPRank.insert err: %+v, ele: %+v in rank: %s ", err, ele, e)
				return
			}
		} else { // 榜单未满
			tailEle := e.List[len(e.List)-1]
			if err := e.shift(i, len(e.List)-1); err != nil {
				log.Error("ElecPrepUPRank.insert err: %+v, ele: %+v in rank: %s ", err, ele, e)
				return
			}
			tailEle.Rank++
			e.List = append(e.List, tailEle)
		}
		e.List[i] = ele
		break
	}
	// 插入到末位
	if ele.Rank < 0 {
		e.List = append(e.List, ele)
		ele.Rank = len(e.List)
	}
}

// shift 将排名从 fromRank 起后移一位到 toRank，原 toRank 会被丢弃，原 fromRank 将会空出
func (e *RankElecPrepUPProto) shift(fromIndex, toIndex int) (err error) {
	if fromIndex > toIndex {
		err = errors.Errorf("shift from(%d) > to(%d)", fromIndex, toIndex)
		return
	}
	if len(e.List)-1 < toIndex {
		err = errors.Errorf("shift out of range [%d,%d], just have: %d", fromIndex, toIndex, len(e.List))
		return
	}
	lastEle := e.List[fromIndex]
	e.List[fromIndex] = nil

	for i := fromIndex + 1; i <= toIndex; i++ {
		lastEle.Rank++
		e.List[i], lastEle = lastEle, e.List[i]
	}
	return
}

// Message binlog databus msg.
type Message struct {
	Action string          `json:"action"`
	Table  string          `json:"table"`
	New    json.RawMessage `json:"new"`
	Old    json.RawMessage `json:"old"`
}

// ElecUserSetting .
type ElecUserSetting int32

// ShowMessage 充电榜单是否显示留言
func (e ElecUserSetting) ShowMessage() bool {
	return (e & 0x1) > 0
}
