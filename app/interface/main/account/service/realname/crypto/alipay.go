package crypto

import (
	stdcrypto "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"net/url"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// Alipay alipay cryptor
type Alipay struct {
	aliPub   []byte
	biliPriv []byte
}

// NewAlipay is.
func NewAlipay(aliPub, biliPriv string) (a *Alipay) {
	return &Alipay{
		aliPub:   ParsePublicKey(aliPub),
		biliPriv: ParsePrivateKey(biliPriv),
	}
}

func (e *Alipay) splitData(originalData []byte, packageSize int) (r [][]byte) {
	var src = make([]byte, len(originalData))
	copy(src, originalData)

	r = make([][]byte, 0)
	if len(src) <= packageSize {
		return append(r, src)
	}
	for len(src) > 0 {
		var p = src[:packageSize]
		r = append(r, p)
		src = src[packageSize:]
		if len(src) <= packageSize {
			r = append(r, src)
			break
		}
	}
	return r
}

// EncryptParam rsa encrypt.
func (e *Alipay) EncryptParam(p url.Values) (ep string, err error) {
	var (
		pubInterface interface{}
		pub          *rsa.PublicKey
		data         []byte
		block        *pem.Block
	)
	block, _ = pem.Decode(e.aliPub)
	if block == nil {
		err = errors.New("private key error")
		return
	}
	if pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		err = errors.WithStack(err)
		return
	}
	pub = pubInterface.(*rsa.PublicKey)
	var sd = e.splitData([]byte(p.Encode()), pub.N.BitLen()/8-11)
	for _, d := range sd {
		var pd []byte
		if pd, err = rsa.EncryptPKCS1v15(rand.Reader, pub, d); err != nil {
			err = errors.WithStack(err)
			return
		}
		data = append(data, pd...)
	}
	ep = base64.StdEncoding.EncodeToString(data)
	return
}

// SignParam sign alipay param
func (e *Alipay) SignParam(p url.Values) (sign string, err error) {
	if p == nil {
		p = make(url.Values)
	}
	var pList = make([]string, 0)
	for key := range p {
		var value = strings.TrimSpace(p.Get(key))
		if len(value) > 0 {
			pList = append(pList, key+"="+value)
		}
	}
	sort.Strings(pList)
	var src = strings.Join(pList, "&")

	var h = stdcrypto.SHA256.New()
	if _, err = h.Write([]byte(src)); err != nil {
		err = errors.WithStack(err)
		return
	}
	var (
		hashed = h.Sum(nil)
		pri    *rsa.PrivateKey
		data   []byte
		block  *pem.Block
	)
	block, _ = pem.Decode(e.biliPriv)
	if block == nil {
		err = errors.New("private key error")
		return
	}
	if pri, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		err = errors.WithStack(err)
		return
	}
	if data, err = rsa.SignPKCS1v15(rand.Reader, pri, stdcrypto.SHA256, hashed); err != nil {
		err = errors.WithStack(err)
		return
	}
	sign = base64.StdEncoding.EncodeToString(data)
	return
}
