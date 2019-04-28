package dao

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func (d *Dao) berserkerSign(ak, sk, dt, ver string) string {
	str := fmt.Sprintf("%sappKey%stimestamp%sversion%s%s", sk, ak, dt, ver, sk)
	b := md5.Sum([]byte(str))
	sign := hex.EncodeToString(b[:])
	return sign
}
