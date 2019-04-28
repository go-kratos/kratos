package web

import (
	"context"

	"go-common/library/log"
)

const _ugcIncre = "web_goblin"

// UgcSearch ugc insert .
func (d *Dao) UgcSearch(ctx context.Context, data map[string]interface{}) (err error) {
	insert := d.ela.NewUpdate(_ugcIncre).Insert()
	insert.AddData(_ugcIncre, data)
	if err = insert.Do(ctx); err != nil {
		log.Error("insert.Do  error(%v)", err)
	}
	return
}
