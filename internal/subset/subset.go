package subset

import (
	"errors"
	"hash/fnv"
	"sort"
	"strconv"
)

type member interface {
	String() string
}

var errEmptyCircle = errors.New("empty circle")

// Subset returns a deterministic subset for the given select key.
func Subset[M member](selectKey string, inss []M, num int) []M {
	if num <= 0 || len(inss) <= num {
		return inss
	}

	circle := newConsistent[M]()
	circle.set(inss)
	out, err := circle.getN(selectKey, num)
	if err != nil {
		return inss
	}
	return out
}

type consistent[M member] struct {
	circle       map[uint32]M
	members      map[string]bool
	sortedHashes []uint32
}

func newConsistent[M member]() *consistent[M] {
	return &consistent[M]{
		circle:  make(map[uint32]M),
		members: make(map[string]bool),
	}
}

func (c *consistent[M]) set(inss []M) {
	for _, ins := range inss {
		if c.members[ins.String()] {
			continue
		}
		c.add(ins)
	}
}

func (c *consistent[M]) add(ins M) {
	for i := 0; i < 160; i++ {
		c.circle[hashKey(strconv.Itoa(i)+ins.String())] = ins
	}
	c.members[ins.String()] = true
	c.updateSortedHashes()
}

func (c *consistent[M]) getN(key string, n int) ([]M, error) {
	if len(c.circle) == 0 {
		return nil, errEmptyCircle
	}
	if len(c.members) < n {
		n = len(c.members)
	}
	offset := c.search(hashKey(key))
	out := make([]M, 0, n)
	for i := offset; ; i++ {
		if i >= len(c.sortedHashes) {
			i = 0
		}
		ins := c.circle[c.sortedHashes[i]]
		if !contains(out, ins) {
			out = append(out, ins)
			if len(out) == n {
				return out, nil
			}
		}
		if i == offset-1 || (offset == 0 && i == len(c.sortedHashes)-1) {
			break
		}
	}
	return out, nil
}

func (c *consistent[M]) search(key uint32) int {
	i := sort.Search(len(c.sortedHashes), func(i int) bool {
		return c.sortedHashes[i] > key
	})
	if i >= len(c.sortedHashes) {
		i = 0
	}
	return i
}

func (c *consistent[M]) updateSortedHashes() {
	hashes := c.sortedHashes[:0]
	for key := range c.circle {
		hashes = append(hashes, key)
	}
	sort.Slice(hashes, func(i, j int) bool {
		return hashes[i] < hashes[j]
	})
	c.sortedHashes = hashes
}

func contains[M member](inss []M, ins M) bool {
	for _, item := range inss {
		if item.String() == ins.String() {
			return true
		}
	}
	return false
}

func hashKey(key string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(key))
	return h.Sum32()
}
