package util

import "sort"
import (
	recsys "go-common/app/service/bbq/recsys/api/grpc/v1"
	"math"
)

//Pair A data structure to hold a key/value pair.
type Pair struct {
	Key   int64
	Value int64
}

//PairList A slice of Pairs that implements sort.Interface to sort by Value in descending order.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

//SortMapByValue A function to turn a map into a PairList, then sort and return it.
func SortMapByValue(m map[int64]int64) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

//PairStr A data structure to hold a key/value pair.
type PairStr struct {
	Key   string
	Value string
}

//PairStrList A slice of Pairs that implements sort.Interface to sort by Value in descending order.
type PairStrList []PairStr

func (p PairStrList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairStrList) Len() int           { return len(p) }
func (p PairStrList) Less(i, j int) bool { return p[i].Value > p[j].Value }

//SortStrMapByValue A function to turn a map into a PairList, then sort and return it.
func SortStrMapByValue(m map[string]string) PairStrList {
	p := make(PairStrList, len(m))
	i := 0
	for k, v := range m {
		p[i] = PairStr{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

//PairStrInt A data structure to hold a key/value pair.
type PairStrInt struct {
	Key   string
	Value int
}

//PairStrIntList A slice of Pairs that implements sort.Interface to sort by Value in descending order.
type PairStrIntList []PairStrInt

func (p PairStrIntList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairStrIntList) Len() int           { return len(p) }
func (p PairStrIntList) Less(i, j int) bool { return p[i].Value > p[j].Value }

//SortStrIntMapByValue A function to turn a map into a PairList, then sort and return it.
func SortStrIntMapByValue(m map[string]int) PairStrIntList {
	p := make(PairStrIntList, len(m))
	i := 0
	for k, v := range m {
		p[i] = PairStrInt{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

//Records sort
type Records []*recsys.RecsysRecord

func (a Records) Len() int           { return len(a) }
func (a Records) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Records) Less(i, j int) bool { return a[i].Score < a[j].Score }

//ScoreCount ...
func ScoreCount(count float64) (score float64) {
	//score = math.Min(float64(count), 3.0)
	maxCount := 10.0
	count = math.Min(count, maxCount)
	score = (1 + 0.1*count) / (1 + 0.1*maxCount)
	return
}

//ScoreTimeDiff ...
func ScoreTimeDiff(timeDiff float64) (score float64) {
	score = 1 - 0.56*math.Pow(timeDiff/3600, 0.06)
	return
}
