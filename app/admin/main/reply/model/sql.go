package model

import (
	"database/sql/driver"
	"encoding/binary"
)

// Int64Bytes implements the Scanner interface.
type Int64Bytes []int64

// Scan parse the data into int64 slice
func (is *Int64Bytes) Scan(src interface{}) (err error) {
	switch sc := src.(type) {
	case []byte:
		var res []int64
		for i := 0; i < len(sc) && i+8 <= len(sc); i += 8 {
			ui := binary.BigEndian.Uint64(sc[i : i+8])
			res = append(res, int64(ui))
		}
		*is = res
	}
	return
}

// Value marshal int64 slice to driver.Value,each int64 will occupy Fixed 8 bytes
func (is Int64Bytes) Value() (driver.Value, error) {
	return is.Bytes(), nil
}

// Bytes marshal int64 slice to bytes,each int64 will occupy Fixed 8 bytes
func (is Int64Bytes) Bytes() []byte {
	res := make([]byte, 0, 8*len(is))
	for _, i := range is {
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(i))
		res = append(res, bs...)
	}
	return res
}

// Evict get rid of the sepcified num  from the slice
func (is *Int64Bytes) Evict(e int64) (ok bool) {
	res := make([]int64, len(*is)-1)
	for _, v := range *is {
		if v != e {
			res = append(res, v)
		} else {
			ok = true
		}
	}
	*is = res
	return
}

// Exist judge the sepcified num is in the slice or not
func (is Int64Bytes) Exist(i int64) (e bool) {
	for _, v := range is {
		if v == i {
			e = true
			return
		}
	}
	return
}
