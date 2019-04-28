package dao

const (
	_setManual  = "REPLACE INTO ugc_archive (aid, manual, deleted) VALUES (?,?,?)"
	_needImport = 1
	_notDeleted = 0
)

// NeedImport sets the archive as manual, if it's deleted, we recover it
func (d *Dao) NeedImport(aid int64) (err error) {
	return d.DB.Exec(_setManual, aid, _needImport, _notDeleted).Error
}
