package dao

import (
	"context"

	"go-common/library/log"

	"github.com/tsuna/gohbase/hrpc"
)

var (
	prefix = "bucket_"
)

// CreateTable .
// TODO check namespace
func (d *Dao) CreateTable(c context.Context, table string) error {
	families := make(map[string]map[string]string)
	families["bfsfile"] = map[string]string{
		"BLOOMFILTER":         "ROW",
		"VERSIONS":            "1",
		"IN_MEMORY":           "false",
		"KEEP_DELETED_CELLS":  "false",
		"DATA_BLOCK_ENCODING": "NONE",
		"TTL":               "2147483647", // NOTE: 2147483647 is forever
		"COMPRESSION":       "NONE",
		"MIN_VERSIONS":      "0",
		"BLOCKCACHE":        "true",
		"BLOCKSIZE":         "65536",
		"REPLICATION_SCOPE": "0",
	}
	b := [][]byte{[]byte("1"), []byte("2"), []byte("3"), []byte("4"), []byte("5"), []byte("6"), []byte("7"),
		[]byte("8"), []byte("9"), []byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f")}
	ct := hrpc.NewCreateTable(c, []byte(prefix+table), families, hrpc.SplitKeys(b))
	err := d.hbase.CreateTable(ct)
	if err != nil {
		log.Error("CreateTable(),err:%+v", err)
		return err
	}
	return nil
}
