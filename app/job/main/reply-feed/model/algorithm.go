package model

import (
	"math"
	"time"
)

// Algorithm Algorithm
type Algorithm interface {
	Score(stat *ReplyStat) *ReplyScore
	Slots() []int
	Name() string
}

/*
Origin Algorithm.
*/

// Origin origin algorithm
type Origin struct {
	name  string
	slots []int
}

// NewOrigin NewLikeDesc Aogorithm
func NewOrigin(name string, slots []int) *Origin {
	return &Origin{
		name:  name,
		slots: slots,
	}
}

// Name get name
func (o *Origin) Name() string {
	return o.name
}

// Slots get slots
func (o *Origin) Slots() []int {
	return o.slots
}

// Score calc score
func (o *Origin) Score(stat *ReplyStat) (rs *ReplyScore) {
	rs = new(ReplyScore)
	rs.RpID = stat.RpID
	if stat.Like < 0 || stat.Reply < 0 || stat.Report < 0 || stat.Hate < 0 {
		return
	}
	score := int64(100 * ((stat.Like + 2) / (stat.Hate + 4 + stat.Report)))
	score = score<<32 | (int64(stat.Reply) & 0xFFFFFFFF)
	rs.Score = float64(score)
	return
}

/*
LikeDesc Algorithm order by like desc.
*/

// LikeDesc like desc
type LikeDesc struct {
	name  string
	slots []int
}

// NewLikeDesc NewLikeDesc Aogorithm
func NewLikeDesc(name string, slots []int) *LikeDesc {
	return &LikeDesc{
		name:  name,
		slots: slots,
	}
}

// Name get name
func (l *LikeDesc) Name() string {
	return l.name
}

// Slots get slots
func (l *LikeDesc) Slots() []int {
	return l.slots
}

// Score calc score
func (l *LikeDesc) Score(stat *ReplyStat) (rs *ReplyScore) {
	rs = new(ReplyScore)
	rs.RpID = stat.RpID
	rs.Score = float64(stat.Like)
	return
}

/*
WilsonLHRR wilson algorithm
like,reply
hate,report
*/

// WilsonLHRR WilsonLHRR
type WilsonLHRR struct {
	name   string
	slots  []int
	weight *WilsonLHRRWeight
}

// NewWilsonLHRR NewWilsonLHRR
func NewWilsonLHRR(name string, slots []int, weight *WilsonLHRRWeight) *WilsonLHRR {
	return &WilsonLHRR{
		name:   name,
		slots:  slots,
		weight: weight,
	}
}

// Name get name
func (w *WilsonLHRR) Name() string {
	return w.name
}

// Slots get slots
func (w *WilsonLHRR) Slots() []int {
	return w.slots
}

// Score calc score
func (w *WilsonLHRR) Score(stat *ReplyStat) (rs *ReplyScore) {
	rs = new(ReplyScore)
	rs.RpID = stat.RpID
	if stat.Like < 0 || stat.Reply < 0 || stat.Report < 0 || stat.Hate < 0 {
		return
	}
	ups := float64(stat.Like)*w.weight.Like + float64(stat.Reply)*w.weight.Reply
	downs := float64(stat.Hate)*w.weight.Hate + float64(stat.Report)*w.weight.Report
	n := ups + downs
	if n == 0 {
		return
	}
	z := float64(2)
	p := ups / n
	rs.Score = (p + math.Pow(z, 2)/(2*n) - (z/(2*n))*math.Sqrt(4*n*(1-p)*p+math.Pow(z, 2))) / (1 + math.Pow(z, 2)/n)
	return
}

/*
WilsonLHRRFluid wilson algorightm dynamic score by time
like,reply
hate,report
*/

// WilsonLHRRFluid WilsonLHRRFluid
type WilsonLHRRFluid struct {
	name   string
	slots  []int
	weight *WilsonLHRRFluidWeight
}

// NewWilsonLHRRFluid NewWilsonLHRRFluid
func NewWilsonLHRRFluid(name string, slots []int, weight *WilsonLHRRFluidWeight) *WilsonLHRRFluid {
	return &WilsonLHRRFluid{
		name:   name,
		slots:  slots,
		weight: weight,
	}
}

// Name get name
func (w *WilsonLHRRFluid) Name() string {
	return w.name
}

// Slots get slots
func (w *WilsonLHRRFluid) Slots() []int {
	return w.slots
}

func coolDownFunc(weight *WilsonLHRRFluidWeight, duration float64) (coefficient float64) {
	return 1.5 - (0.5 / (1 + math.Exp(-weight.Slope*(duration))))
}

// Score calc score
func (w *WilsonLHRRFluid) Score(stat *ReplyStat) (rs *ReplyScore) {
	rs = new(ReplyScore)
	rs.RpID = stat.RpID
	if stat.Like < 0 || stat.Reply < 0 || stat.Report < 0 || stat.Hate < 0 {
		return
	}
	ups := float64(stat.Like)*w.weight.Like + float64(stat.Reply)*w.weight.Reply
	downs := float64(stat.Hate)*w.weight.Hate + float64(stat.Report)*w.weight.Report
	n := ups + downs
	if n == 0 {
		return
	}
	z := float64(2)
	p := ups / n
	coefficient := coolDownFunc(w.weight, float64(time.Now().Unix()-int64(stat.ReplyTime))/86400)
	rs.Score = ((p + math.Pow(z, 2)/(2*n) - (z/(2*n))*math.Sqrt(4*n*(1-p)*p+math.Pow(z, 2))) / (1 + math.Pow(z, 2)/n)) * coefficient
	return
}
