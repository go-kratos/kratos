package dao

import (
	"context"
	"time"

	"go-common/app/interface/main/kvo/model"

	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_upDocument  = "INSERT INTO document(check_sum,doc,ctime,mtime) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE doc=?"
	_getDocument = "SELECT check_sum,doc FROM document WHERE check_sum=?"
)

// Document get docuemtn
func (d *Dao) Document(ctx context.Context, checkSum int64) (doc *model.Document, err error) {
	row := d.getDocument.QueryRow(ctx, checkSum)
	doc = &model.Document{}
	err = row.Scan(&doc.CheckSum, &doc.Doc)
	if err != nil {
		if err == sql.ErrNoRows {
			doc = nil
			err = nil
			return
		}
		log.Error("row.scan err:%v", err)
	}
	return
}

// TxUpDocuement add a document
func (d *Dao) TxUpDocuement(ctx context.Context, tx *sql.Tx, checkSum int64, data string, now time.Time) (err error) {
	_, err = tx.Exec(_upDocument, checkSum, data, now, now, data)
	if err != nil {
		log.Error("db.exec err:%v", err)
	}
	return
}
