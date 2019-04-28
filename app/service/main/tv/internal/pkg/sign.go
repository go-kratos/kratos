package pkg

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"go-common/library/log"

	"github.com/google/go-querystring/query"
)

const (
	_sign = "sign"
)

// Signer.
type Signer struct {
	Key string
}

func (s *Signer) struct2Values(in interface{}) (out url.Values, err error) {
	out, err = query.Values(in)
	if err != nil {
		return
	}
	out.Del(_sign)
	return
}

func (s *Signer) encodeUrlValues(values url.Values) (res string) {
	keys := make([]string, 0)
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	kvSli := make([]string, 0)
	for _, key := range keys {
		kvStr := fmt.Sprintf("%s=%s", key, values.Get(key))
		kvSli = append(kvSli, kvStr)
	}
	paramsStr := strings.Join(kvSli, "&")
	return strings.ToLower(paramsStr) + "&key=" + s.Key
}

func (s *Signer) escapeStr(str string) (res string) {
	return s.adapt(url.QueryEscape(str))
}

func (s *Signer) md5(str string) (res string) {
	var (
		buf bytes.Buffer
	)
	buf.WriteString(str)
	hexMd5 := md5.Sum(buf.Bytes())
	return strings.ToUpper(hex.EncodeToString(hexMd5[:]))
}

// Sign makes sign for yst request.
func (s *Signer) Sign(in interface{}) (sign string, err error) {
	var (
		values url.Values
	)
	if values, err = s.struct2Values(in); err != nil {
		log.Error("signer.struct2Values(%+v) err(%+v)", in, err)
		return
	}
	encodedVals := s.encodeUrlValues(values)
	escapedVals := s.escapeStr(encodedVals)
	sign = s.md5(escapedVals)
	log.Info("Signer.Sign(%+v) sign(%s)", in, sign)
	return
}

// NOTE: 经核对，云视听只有 `*` 不需要转义
func (s *Signer) adapt(str string) string {
	res := str
	res = strings.Replace(res, "%2A", "*", -1)
	return res
}
