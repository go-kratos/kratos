package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// const .
const (
	BlockSizeMIN = 16 // AES-128
	BlockSizeMID = 24 // AES-192
	BlockSizeMAX = 32 // AES-256
)

// Main mainsite cryptor
type Main struct {
	publicKey, privateKey []byte
}

// NewMain is.
func NewMain(pub, priv string) (e *Main) {
	return &Main{
		publicKey:  []byte(pub),
		privateKey: []byte(priv),
	}
}

// IMGEncrypt rsa + AES-128
func (e *Main) IMGEncrypt(raw []byte) (data []byte, err error) {
	if len(raw) == 0 {
		return
	}
	var (
		hash             = md5.New()
		randToken        []byte
		rsaData, aesData []byte
		aesBase64        []byte
		rsaBase64Str     string
		buf              bytes.Buffer
	)
	if _, err = hash.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10))); err != nil {
		err = errors.WithStack(err)
		return
	}
	randToken = []byte(hex.EncodeToString(hash.Sum(nil)))
	if rsaData, err = e.RsaEncrypt(randToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	if aesData, err = e.AESEncrypt(randToken, raw, BlockSizeMIN); err != nil {
		err = errors.WithStack(err)
		return
	}
	rsaBase64Str = base64.StdEncoding.EncodeToString(rsaData)
	aesBase64 = make([]byte, base64.StdEncoding.EncodedLen(len(aesData)))
	base64.StdEncoding.Encode(aesBase64, aesData)

	fmt.Fprintf(&buf, "%04x", len(rsaBase64Str))
	buf.Write([]byte(rsaBase64Str))
	buf.Write(aesBase64)
	data = buf.Bytes()
	return
}

// IMGDecrypt rsa + AES-128
func (e *Main) IMGDecrypt(raw []byte) (data []byte, err error) {
	if len(raw) == 0 {
		return
	}
	var (
		rsaLen            int64
		rsaRandToken      []byte
		randToken         []byte
		aesBase64         []byte
		aesData           []byte
		base64DecodedSize int
	)
	if rsaLen, err = strconv.ParseInt(string(raw[:4]), 16, 64); err != nil {
		err = errors.WithStack(err)
		return
	}
	if rsaRandToken, err = base64.StdEncoding.DecodeString(string(raw[4 : 4+rsaLen])); err != nil {
		err = errors.WithStack(err)
		return
	}
	if randToken, err = e.RsaDecrypt(rsaRandToken); err != nil {
		err = errors.WithStack(err)
		return
	}
	aesBase64 = raw[4+rsaLen:]
	aesData = make([]byte, base64.StdEncoding.DecodedLen(len(aesBase64)))
	if base64DecodedSize, err = base64.StdEncoding.Decode(aesData, aesBase64); err != nil {
		err = errors.WithStack(err)
		return
	}
	if data, err = e.AESDecrypt(randToken, aesData[:base64DecodedSize], BlockSizeMIN); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// RsaEncrypt rsa encrypt.
func (e *Main) RsaEncrypt(text []byte) (data []byte, err error) {
	var (
		block *pem.Block
	)
	block, _ = pem.Decode(e.publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	var (
		pubInterface interface{}
		pub          *rsa.PublicKey
	)
	if pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	pub = pubInterface.(*rsa.PublicKey)
	if data, err = rsa.EncryptPKCS1v15(rand.Reader, pub, text); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// RsaDecrypt rsa decrypt.
func (e *Main) RsaDecrypt(text []byte) (data []byte, err error) {
	var (
		block *pem.Block
	)
	block, _ = pem.Decode(e.privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	var (
		privateKey *rsa.PrivateKey
	)
	if privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		err = errors.WithStack(err)
		return
	}
	if data, err = rsa.DecryptPKCS1v15(rand.Reader, privateKey, text); err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// AESEncrypt AES-128, AES-192, or AES-256 encrypt.
// mod 16, 24, or 32 bytes
func (e *Main) AESEncrypt(key, text []byte, mod int) (data []byte, err error) {
	var (
		block cipher.Block
	)
	if block, err = aes.NewCipher(key[:mod]); err != nil {
		err = errors.WithStack(err)
		return
	}
	msg := pkPadding(text, block.BlockSize())
	ciphertext := make([]byte, len(msg))
	cbc := cipher.NewCBCEncrypter(block, key[mod:])
	cbc.CryptBlocks(ciphertext, []byte(msg))
	return ciphertext, nil
}

// AESDecrypt AES-128, AES-192, or AES-256 decrypt.
// mod 16, 24, or 32 bytes
func (e *Main) AESDecrypt(key, text []byte, mod int) (data []byte, err error) {
	var (
		block cipher.Block
	)
	if block, err = aes.NewCipher(key[:mod]); err != nil {
		err = errors.WithStack(err)
		return
	}
	blockModel := cipher.NewCBCDecrypter(block, key[mod:])
	ciphertext := make([]byte, len(text))
	blockModel.CryptBlocks(ciphertext, text)
	ciphertext = pkUnPadding(ciphertext)
	return ciphertext, nil
}

func pkPadding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(0)}, padding)
	return append(src, padtext...)
}

func pkUnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
