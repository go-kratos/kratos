package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net/url"
	"strings"
)

// Sign fn
func Sign(params url.Values) (query string, err error) {
	if len(params) == 0 {
		return
	}
	if params.Get("appkey") == "" {
		err = fmt.Errorf("utils http get must have parameter appkey")
		return
	}
	if params.Get("appsecret") == "" {
		err = fmt.Errorf("utils http get must have parameter appsecret")
		return
	}
	if params.Get("sign") != "" {
		err = fmt.Errorf("utils http get must have not parameter sign")
		return
	}
	// sign
	secret := params.Get("appsecret")
	params.Del("appsecret")
	tmp := params.Encode()
	if strings.IndexByte(tmp, '+') > -1 {
		tmp = strings.Replace(tmp, "+", "%20", -1)
	}
	mh := md5.Sum([]byte(tmp + secret))
	params.Set("sign", hex.EncodeToString(mh[:]))
	query = params.Encode()
	return
}

//DeDuplicationSlice for del repeat element
func DeDuplicationSlice(a []int64) (b []int64) {
	if len(a) == 0 {
		return
	}
	isHas := make(map[int64]bool)
	b = make([]int64, 0)
	for _, v := range a {
		if ok := isHas[v]; !ok {
			isHas[v] = true
			b = append(b, v)
		}
	}
	return
}

//ContainAll all element of a  contain in the b.
func ContainAll(a []int64, b []int64) bool {
	isHas := make(map[int64]bool)
	for _, k := range b {
		isHas[k] = true
	}
	for _, v := range a {
		if !isHas[v] {
			return false
		}
	}
	return true
}

//ContainAtLeastOne fn
func ContainAtLeastOne(a []int64, b []int64) bool {
	if len(a) == 0 {
		return true
	}
	isHas := make(map[int64]bool)
	for _, k := range b {
		isHas[k] = true
	}
	for _, v := range a {
		if isHas[v] {
			return true
		}
	}
	return false
}

//ElementInSlice fn
func ElementInSlice(a int64, b []int64) bool {
	if len(b) == 0 {
		return false
	}
	for _, v := range b {
		if a == v {
			return true
		}
	}
	return false
}

//RandomSliceKeys for get random keys from slice by rand.
func RandomSliceKeys(start int, end int, count int, seed int64) []int {
	if end < start || (end-start) < count {
		return nil
	}
	nums := make([]int, 0)
	r := rand.New(rand.NewSource(seed))
	for len(nums) < count {
		num := r.Intn((end - start)) + start
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}
