package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/pkg/errors"
)

// const .
const (
	BlockSizeMIN = 16 // AES-128
	BlockSizeMID = 24 // AES-192
	BlockSizeMAX = 32 // AES-256
)

// Realname realname cryptor
type Realname struct {
	publicKey, privateKey []byte
}

// NewRealname get a new instance of realname cryptor
func NewRealname(pub, priv string) (e *Realname) {
	return &Realname{
		publicKey:  []byte(pub),
		privateKey: []byte(priv),
	}
}

// RsaEncrypt rsa encrypt.
func (e *Realname) RsaEncrypt(text []byte) ([]byte, error) {
	block, _ := pem.Decode(e.publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, text)
}

// RsaDecrypt rsa decrypt.
func (e *Realname) RsaDecrypt(text []byte) ([]byte, error) {
	block, _ := pem.Decode(e.privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, text)
}

// CardEncrypt encrypt realname card code
func (e *Realname) CardEncrypt(text []byte) (data []byte, err error) {
	if len(text) == 0 {
		return []byte{}, nil
	}
	var (
		encryptedData []byte
	)
	if encryptedData, err = e.RsaEncrypt(text); err != nil {
		err = errors.Wrapf(err, "rsa encrypt(%s)", text)
		return
	}
	data = make([]byte, base64.StdEncoding.EncodedLen(len(encryptedData)))
	base64.StdEncoding.Encode(data, encryptedData)
	return
}

// CardDecrypt decrypt realname card code
func (e *Realname) CardDecrypt(data []byte) (text []byte, err error) {
	if len(data) == 0 {
		return []byte{}, nil
	}
	var (
		decryptedData []byte
		size          int
	)
	decryptedData = make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	if size, err = base64.StdEncoding.Decode(decryptedData, data); err != nil {
		err = errors.Wrapf(err, "base decode %s", data)
		return
	}
	if text, err = e.RsaDecrypt(decryptedData[:size]); err != nil {
		err = errors.Wrapf(err, "rsa decrypt %s", decryptedData)
		return
	}
	return
}
