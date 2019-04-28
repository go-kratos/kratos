package reply

import (
	"database/sql/driver"
	"encoding/binary"
)

//Ints be used to MySql\Protobuf varbinary converting.
type Ints []int64

// MarshalTo  marshal int64 slice to bytes,each int64 will occupy Fixed 8 bytes.
//if the argument data not supplied with the full size,it will return the actual written size
func (is Ints) MarshalTo(data []byte) (int, error) {
	for i, n := range is {
		start := i * 8
		end := (i + 1) * 8

		if len(data) < end {
			return start, nil
		}
		bs := data[start:end]
		binary.BigEndian.PutUint64(bs, uint64(n))
	}

	return 8 * len(is), nil
}

// Size return the total size it will occupy in bytes
func (is Ints) Size() int {
	return len(is) * 8
}

// Unmarshal parse the data into int64 slice
func (is *Ints) Unmarshal(data []byte) error {
	return is.Scan(data)
}

// Scan parse the data into int64 slice
func (is *Ints) Scan(src interface{}) (err error) {
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
func (is Ints) Value() (driver.Value, error) {
	return is.Bytes(), nil
}

// Bytes marshal int64 slice to bytes,each int64 will occupy Fixed 8 bytes
func (is Ints) Bytes() []byte {
	res := make([]byte, 0, 8*len(is))
	for _, i := range is {
		bs := make([]byte, 8)
		binary.BigEndian.PutUint64(bs, uint64(i))
		res = append(res, bs...)
	}
	return res
}

// Evict get rid of the sepcified num  from the slice
func (is *Ints) Evict(e int64) (ok bool) {
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
func (is Ints) Exist(i int64) (e bool) {
	for _, v := range is {
		if v == i {
			e = true
			return
		}
	}
	return
}
