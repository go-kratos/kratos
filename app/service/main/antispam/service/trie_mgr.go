package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/util/trie"

	"go-common/library/log"
)

var (
	// ErrTrieNotFound .
	ErrTrieNotFound = errors.New("so such key in trie tree")
)

// Put .
func (ctrie *ConcurrentTrie) Put(key string, val interface{}) {
	ctrie.Lock()
	defer ctrie.Unlock()
	ctrie.trier.Put(key, val)
}

// Delete .
func (ctrie *ConcurrentTrie) Delete(key string) {
	ctrie.RLock()
	if ctrie.trier.Get(key) == nil {
		ctrie.RUnlock()
		return
	}
	ctrie.RUnlock()

	ctrie.Lock()
	defer ctrie.Unlock()
	ctrie.trier.Put(key, nil)
}

// KeywordLimitInfo the data stored in trie tree
type KeywordLimitInfo struct {
	KeywordID int64
	LimitType string
}

func (ctrie *ConcurrentTrie) find(content string) (string, *KeywordLimitInfo) {
	ctrie.RLock()
	defer ctrie.RUnlock()

	key, val := ctrie.trier.Find(content, "")
	if val == nil {
		return "", nil
	}
	if limitInfo, ok := val.(*KeywordLimitInfo); ok {
		return key, limitInfo
	}
	return "", nil
}

// NewConcurrentTrie .
func NewConcurrentTrie() *ConcurrentTrie {
	return &ConcurrentTrie{}
}

// ConcurrentTrie .
type ConcurrentTrie struct {
	sync.RWMutex
	trier trie.Trier
}

// LastModifiedTime .
func (mgr *TrieMgr) LastModifiedTime() *time.Time {
	if t, ok := mgr.lastModifiedTime.Load().(time.Time); ok {
		return &t
	}
	return nil
}

// UpdateLastModifiedTime .
func (mgr *TrieMgr) UpdateLastModifiedTime() {
	mgr.lastModifiedTime.Store(time.Now())
}

// Refresh .
func (mgr *TrieMgr) Refresh() {
	ks, _, err := mgr.service.GetKeywordsByCond(context.TODO(),
		&Condition{LastModifiedTime: mgr.LastModifiedTime()})
	if err != nil {
		return
	}

	mgr.UpdateLastModifiedTime()
	for _, k := range ks {
		if k.State == model.StateDeleted ||
			k.Tag == model.KeywordTagWhite {
			for _, ctrie := range mgr.tries[k.Area] {
				ctrie.Delete(k.Content)
			}
			continue
		}
		for tag, ctrie := range mgr.tries[k.Area] {
			// NOTICE: must exec `put` before exec `delete`
			if tag == k.Tag {
				ctrie.Put(k.Content, &KeywordLimitInfo{
					KeywordID: k.ID,
					LimitType: tag,
				})
			} else {
				ctrie.Delete(k.Content)
			}
		}
	}
	log.Info("refresh finished.")
}

func (ctrie *ConcurrentTrie) build(area string, tag string, s Service, defaultDBLimit int64) {
	newTrier := trie.New()
	var lastID int64
	for {
		ks, err := s.GetKeywordsByOffsetLimit(context.TODO(),
			&Condition{
				Area:   area,
				Offset: fmt.Sprintf("%d", lastID),
				Limit:  fmt.Sprintf("%d", defaultDBLimit),
				Tags:   []string{tag}})

		if err != nil || len(ks) == 0 {
			break
		}
		for _, k := range ks {
			if k.ID > lastID {
				lastID = k.ID
			}
			newTrier.Put(k.Content, &KeywordLimitInfo{
				KeywordID: k.ID,
				LimitType: tag,
			})
		}
		if len(ks) < int(defaultDBLimit) {
			break
		}
	}

	ctrie.Lock()
	defer ctrie.Unlock()
	ctrie.trier = newTrier
}

// Build .
func (mgr *TrieMgr) Build(dBLimit int64) {
	mgr.UpdateLastModifiedTime()
	for area, subMap := range mgr.tries {
		for tag, ctrie := range subMap {
			ctrie.build(area, tag, mgr.service, dBLimit)
		}
	}
	log.Info("build finished.")
}

// Get .
func (mgr *TrieMgr) Get(area string, content string) (string, *KeywordLimitInfo, error) {
	if ttr, ok := mgr.tries[area]; ok {
		for _, severity := range []string{
			model.KeywordTagBlack,
			model.KeywordTagRestrictLimit,
			model.KeywordTagDefaultLimit,
		} {
			if tr, ok := ttr[severity]; ok {
				if k, v := tr.find(content); v != nil {
					return k, v, nil
				}
			}
		}
	}
	return "", nil, ErrTrieNotFound
}

// NewTagConcurrentTrieMap .
func NewTagConcurrentTrieMap() map[string]*ConcurrentTrie {
	return map[string]*ConcurrentTrie{
		model.KeywordTagDefaultLimit:  NewConcurrentTrie(),
		model.KeywordTagRestrictLimit: NewConcurrentTrie(),
		model.KeywordTagBlack:         NewConcurrentTrie(),
	}
}

// NewTrieMgr .
func NewTrieMgr(service Service) *TrieMgr {
	return &TrieMgr{
		service: service,
		tries: map[string]map[string]*ConcurrentTrie{
			model.AreaReply:      NewTagConcurrentTrieMap(),
			model.AreaIMessage:   NewTagConcurrentTrieMap(),
			model.AreaLiveDM:     NewTagConcurrentTrieMap(),
			model.AreaMainSiteDM: NewTagConcurrentTrieMap(),
		},
	}
}

// TrieMgr .
type TrieMgr struct {
	service          Service
	tries            map[string]map[string]*ConcurrentTrie // area:tag:trie, TODO: pick a better representation
	lastModifiedTime atomic.Value
}
