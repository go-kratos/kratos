package service

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	model "go-common/app/job/main/reply/model/reply"
	"go-common/library/log"
)

const (
	_replySliceNum = 20000
)

func dialogMapByRoot(rootID int64, rps []*model.RpItem, oid int64, tp int8) (dialogMap map[int64][]*model.RpItem) {
	length := len(rps)
	dialogMap = make(map[int64][]*model.RpItem)
	// 根评论下没有评论
	if length == 0 {
		return
	}

	// 这里由于种种原因， 可能子评论的ID比父评论大，故按Floor排序
	// 按Floor严格排序，保证父评论也几乎是升序排列，提高后续的命中率
	sort.Slice(rps, func(i, j int) bool {
		return rps[i].Floor < rps[j].Floor
	})

	for i := 0; i < length; i++ {
		if rps[i] == nil {
			return
		}
		if rps[i].Parent == rootID {
			// 对话根评论
			continue
		}
		// i-1 must >= 0, 因为第一个元素必然是对话根评论which parrent=root, 出现这种情况说明数据是脏的
		if i-1 < 0 {
			log.Error("invalid data, rootID(%d), rpID(%d), parrent(%d) oid(%d) type(%d)", rootID, rps[i].ID, rps[i].Parent, oid, tp)
			return
		}

		if rps[i].Parent == rps[i-1].Parent {
			// 和上一条评论回复同一条评论的情况，按ID排序，所以这种情况的概率会很高
			rps[i].Next = rps[i-1].Next
		} else {
			var j int
			if sort.IsSorted(model.RpItems(rps)) {
				// 和上一条评论回复不同评论的情况, 其父评论一定在它之前
				j = sort.Search(i, func(n int) bool {
					return rps[n].ID >= rps[i].Parent
				})
			} else {
				for index := range rps[:i] {
					if rps[index].ID == rps[i].Parent {
						j = index
						break
					}
				}
			}
			// search 如果返回j==i说明没搜索,或者遍历到了最后一个, 这种情况说明数据是脏的
			if j == i {
				log.Error("invalid data, rootID(%d), rpID(%d), parrent(%d) oid(%d) type(%d)", rootID, rps[i].ID, rps[i].Parent, oid, tp)
				return nil
			}
			rps[i].Next = rps[j]
		}
	}
	tmp := new(struct {
		ID        int64
		DiaglogID int64
	})
	for i := 0; i < length; i++ {
		if rps[i] == nil {
			return
		}

		next := rps[i].Next
		if next == nil {
			// 如果是对话根评论
			dialogMap[rps[i].ID] = append(dialogMap[rps[i].ID], rps[i])
		} else if next.ID == tmp.ID {
			// 这里tmp缓存了上一个评论的父评论, 减少查找的次数
			// 如果跟上一条评论评论的是同一条评论,则可以直接加进上一个dialog
			dialogMap[tmp.DiaglogID] = append(dialogMap[tmp.DiaglogID], rps[i])
		} else {
			depth := 0
			for next.Next != nil {
				next = next.Next
				depth++
				if depth > 10000 {
					for i := range rps {
						log.Error("rp: %v", rps[i])
					}
					log.Error("recursive reach max depth")
					return nil
				}
			}
		}
	}
	return
}

func (s *Service) setDialogByRoot(c context.Context, oid int64, tp int8, rootID int64) (err error) {
	// 循环获取某个根评论下的所有子评论
	rps, err := s.dao.Reply.FixDialogGetRepliesByRoot(c, oid, tp, rootID)
	if err != nil {
		log.Error("fix dialog error (%v)", err)
		return
	}
	//根据所有子评论构造 key为二级父评论(即对话根评论), value为二级父评论下的所有子评论的map
	dialogMap := dialogMapByRoot(rootID, rps, oid, tp)
	for k, v := range dialogMap {
		ids := make([]int64, len(v))
		for i := range v {
			ids[i] = v[i].ID
		}
		s.dao.Reply.FixDialogSetDialogBatch(c, oid, k, ids)
	}
	return
}

// actionRecoverFixDialog fix dialog
func (s *Service) actionRecoverFixDialog(c context.Context, msg *consumerMsg) {
	var (
		err error
	)
	var d struct {
		Oid  int64 `json:"oid"`
		Tp   int8  `json:"tp"`
		Root int64 `json:"root"`
	}
	if err = json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	s.setDialogByRoot(c, d.Oid, d.Tp, d.Root)
}

func (s *Service) actionRecoverDialog(c context.Context, msg *consumerMsg) {
	var (
		ok  bool
		err error
	)
	var d struct {
		Oid    int64 `json:"oid"`
		Tp     int8  `json:"tp"`
		Root   int64 `json:"root"`
		Dialog int64 `json:"dialog"`
	}
	if err = json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	if ok, err = s.dao.Redis.ExpireDialogIndex(c, d.Dialog); err == nil && !ok {
		rps, err := s.dao.Reply.GetByDialog(c, d.Oid, d.Tp, d.Root, d.Dialog)
		if err != nil {
			return
		}
		err = s.dao.Redis.AddDialogIndex(c, d.Dialog, rps)
		if err != nil {
			log.Error("s.dao.Redis.AddDialogIndex() error (%v)", err)
			return
		}
	}
}

func (s *Service) acionRecoverFloorIdx(c context.Context, msg *consumerMsg) {
	var (
		err    error
		rCount int
		limit  int
		ok     bool
		sub    *model.Subject
	)
	var d struct {
		Oid   int64 `json:"oid"`
		Tp    int8  `json:"tp"`
		Count int   `json:"count"`
		Floor int   `json:"floor"`
	}
	if err = json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	sub, err = s.getSubject(c, d.Oid, d.Tp)
	if err != nil || sub == nil {
		log.Error("s.getSubject(%d,%d) failed!err:=%v", d.Oid, d.Tp, err)
		return
	}
	if sub.RCount == 0 {
		return
	}
	startFloor := sub.Count + 1
	if ok, err = s.dao.Redis.ExpireIndex(c, d.Oid, d.Tp, model.SortByFloor); err == nil && ok {
		startFloor, err = s.dao.Redis.MinScore(c, d.Oid, d.Tp, model.SortByFloor)
		if err != nil {
			log.Error("s.dao.Redis.MinScore(%d,%d) failed!err:=%v", d.Oid, d.Tp, err)
			return
		}
		if startFloor <= 1 {
			if startFloor != -1 {
				if err = s.dao.Redis.AddFloorIndexEnd(c, d.Oid, d.Tp); err != nil {
					log.Error("s.dao.Redis.AddFloorIndexEnd(%d, %d) error(%v)", d.Oid, d.Tp, err)
				}
			}
			return
		}
		if d.Count > 0 {
			rCount, err = s.dao.Redis.CountReplies(c, d.Oid, d.Tp, model.SortByFloor)
			if err != nil {
				log.Error("s.dao.Redis.CountReplies(%d,%d) failed!err:=%v", d.Oid, d.Tp, err)
				return
			}
		}
	}
	if d.Count > 0 {
		limit = d.Count - rCount
	} else if d.Floor > 0 {
		limit = startFloor - d.Floor
	} else {
		log.Warn("RecoverFloorByCount(%d,%d) count(%d) or floor(%d) invalid!", d.Oid, d.Tp, d.Floor, d.Count)
		return
	}
	limit += s.batchNumber
	if limit < (s.batchNumber / 2) {
		return
	} else if limit < s.batchNumber {
		limit = s.batchNumber
	}
	rs, err := s.dao.Reply.GetByFloorLimit(c, d.Oid, d.Tp, startFloor, limit)
	if err != nil {
		log.Error("s.dao.Reply.GetByFloorLimit(%d,%d) failed!err:=%v", d.Oid, d.Tp, err)
		return
	}
	if err = s.dao.Redis.AddFloorIndex(c, d.Oid, d.Tp, rs...); err != nil {
		log.Error("s.dao.Redis.AddFloorIndex(%d, %d) error(%v)", d.Oid, d.Tp, err)
	}
	if len(rs) < limit {
		if err = s.dao.Redis.AddFloorIndexEnd(c, d.Oid, d.Tp); err != nil {
			log.Error("s.dao.Redis.AddFloorIndexEnd(%d, %d) error(%v)", d.Oid, d.Tp, err)
		}
		return
	}
}

// actionRecoverIndex recover index of archive's reply
func (s *Service) actionRecoverIndex(c context.Context, msg *consumerMsg) {
	var (
		err error
		ok  bool
		sub *model.Subject
		rs  []*model.Reply
	)
	var d struct {
		Oid  int64 `json:"oid"`
		Tp   int8  `json:"tp"`
		Sort int8  `json:"sort"`
	}
	if err = json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	if d.Oid <= 0 || !model.CheckSort(d.Sort) {
		log.Error("The structure of doActionRecoverIndex msg.Data(%s) was wrong", msg.Data)
		return
	}
	if ok, err = s.dao.Redis.ExpireIndex(c, d.Oid, d.Tp, d.Sort); err == nil && !ok {
		sub, err = s.getSubject(c, d.Oid, d.Tp)
		if err != nil || sub == nil {
			log.Error("s.getSubject failed , oid(%d,%d) err(%v)", d.Oid, d.Tp, err)
			return
		}
		if d.Sort == model.SortByFloor {
			rs, err = s.dao.Reply.GetAllInSlice(c, d.Oid, d.Tp, sub.Count, _replySliceNum)
			if err != nil {
				log.Error("dao.Reply.GetAllInSlice(%d, %d) error(%v)", d.Oid, d.Tp, err)
				return
			}
			// floor index
			if ok, err = s.dao.Redis.ExpireIndex(c, d.Oid, d.Tp, model.SortByFloor); err == nil && !ok {
				if err = s.dao.Redis.AddFloorIndex(c, d.Oid, d.Tp, rs...); err != nil {
					log.Error("s.dao.Redis.AddFloorIndex(%d, %d) error(%v)", d.Oid, d.Tp, err)
				}
				if err = s.dao.Redis.AddFloorIndexEnd(c, d.Oid, d.Tp); err != nil {
					log.Error("s.dao.Redis.AddFloorIndexEnd(%d, %d) error(%v)", d.Oid, d.Tp, err)
				}
			}
		} else if d.Sort == model.SortByLike {
			rs, err = s.dao.Reply.GetByLikeLimit(c, d.Oid, d.Tp, 30000)
			if err != nil {
				log.Error("dao.Reply.GetAllInSlice(%d, %d) error(%v)", d.Oid, d.Tp, err)
				return
			}
			// like index
			if ok, err = s.dao.Redis.ExpireIndex(c, d.Oid, d.Tp, model.SortByLike); err == nil && !ok {
				rpts, _ := s.dao.Report.GetMapByOid(c, d.Oid, d.Tp)
				if err = s.dao.Redis.AddLikeIndexBatch(c, d.Oid, d.Tp, rpts, rs...); err != nil {
					log.Error("s.dao.Redis.AddLikeIndexBatch(%d, %d) error(%v)", d.Oid, d.Tp, err)
				}
			}
		} else if d.Sort == model.SortByCount {
			rs, err = s.dao.Reply.GetByCountLimit(c, d.Oid, d.Tp, 20000)
			if err != nil {
				log.Error("dao.Reply.GetAllInSlice(%d, %d) error(%v)", d.Oid, d.Tp, err)
				return
			}
			// count index
			if ok, err = s.dao.Redis.ExpireIndex(c, d.Oid, d.Tp, model.SortByCount); err == nil && !ok {
				if err = s.dao.Redis.AddCountIndexBatch(c, d.Oid, d.Tp, rs...); err != nil {
					log.Error("s.dao.Redis.AddCountIndex(%d, %d) error(%v)", d.Oid, d.Tp, err)
				}
			}
		}
		// 回源index时把top缓存初始化，减少attr扫表慢查询
		if n := sub.TopCount(); n > 0 {
			for _, r := range rs {
				if r.IsTop() {
					top := model.SubAttrAdminTop
					if r.IsUpTop() {
						top = model.SubAttrUpperTop
					}
					err = sub.TopSet(r.RpID, top, 1)
					if err == nil {
						_, err = s.dao.Subject.UpMeta(c, d.Oid, d.Tp, sub.Meta, time.Now())
						if err != nil {
							log.Error("s.dao.Subject.UpMeta(%d,%d,%d) failed!err:=%v ", r.RpID, r.Oid, d.Tp, err)
						}
						s.dao.Mc.AddSubject(c, sub)
					}
					// get reply with content
					var rp *model.Reply
					rp, err = s.getReply(c, d.Oid, r.RpID)
					if err == nil && rp != nil {
						s.dao.Mc.AddTop(c, rp)
					}
					n--
				}
				if n == 0 {
					break
				}
			}
		}
	}
}

// actionRecoverRootIndex recover index of root reply
func (s *Service) actionRecoverRootIndex(c context.Context, msg *consumerMsg) {
	var (
		err error
		ok  bool
		rs  []*model.Reply
	)
	var d struct {
		Oid  int64 `json:"oid"`
		Tp   int8  `json:"tp"`
		Root int64 `json:"root"`
	}
	if err = json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	if d.Oid <= 0 || d.Root <= 0 {
		log.Error("The structure of doActionRecoverRootIndex msg.Data(%s) was wrong", msg.Data)
		return
	}
	if ok, err = s.dao.Redis.ExpireNewChildIndex(c, d.Root); err == nil && !ok {
		if rs, err = s.dao.Reply.GetAllByRoot(c, d.Oid, d.Root, d.Tp); err != nil {
			log.Error("dao.Reply.GetAllReply(%d, %d) error(%v)", d.Oid, d.Tp, err)
			return
		}
		if err = s.dao.Redis.AddNewChildIndex(c, d.Root, rs...); err != nil {
			log.Error("s.dao.Redis.AddFloorIndexByRoot(%d, %d) error(%v)", d.Oid, d.Tp, err)
			return
		}
	}
}
