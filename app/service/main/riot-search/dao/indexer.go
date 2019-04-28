package dao

import (
	"github.com/go-ego/riot/types"
)

// Insert doc into index...
func (d *Dao) Insert(id uint64, content string, forceUpdate bool) {
	d.searcher.Index(id, types.DocData{Content: content}, forceUpdate)
}

// Flush force update data from cache to index
func (d *Dao) Flush() {
	d.searcher.Flush()
}

// Remove remove a doc from index
func (d *Dao) Remove(id uint64, forceUpdate bool) {
	d.searcher.RemoveDoc(id, forceUpdate)
}
