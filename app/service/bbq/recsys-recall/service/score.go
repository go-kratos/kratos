package service

import (
	"go-common/app/service/bbq/recsys-recall/dao"
	"math/rand"
)

// Scorer interface
type Scorer interface {
	doScore(uint64, ...interface{}) float32
}

// ScorerManager .
type ScorerManager struct {
	scorers map[string]Scorer
}

// NewScorerManager .
func NewScorerManager(d *dao.Dao) *ScorerManager {
	s := make(map[string]Scorer)

	s["default"] = &DefaultScorer{
		d: d,
	}

	return &ScorerManager{
		scorers: s,
	}
}

// DoScore .
func (sm *ScorerManager) DoScore(svid uint64, name string, params ...interface{}) float32 {
	if s, ok := sm.scorers[name]; ok && s != nil {
		return s.doScore(svid, params...)
	}
	return float32(0)
}

// DefaultScorer .
type DefaultScorer struct {
	d *dao.Dao
}

func (ds *DefaultScorer) doScore(svid uint64, params ...interface{}) float32 {
	return rand.Float32()
}
