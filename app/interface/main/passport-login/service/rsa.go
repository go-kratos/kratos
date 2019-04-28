package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"go-common/app/interface/main/passport-login/model"
)

var (
	// ErrPrivateKeyError private key error
	ErrPrivateKeyError = errors.New("private key error")
	// ErrUnsupportedPriKey Private key error
	ErrUnsupportedPriKey = errors.New("unsupported private key")
)

func priKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, ErrPrivateKeyError
	}
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	priv, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrUnsupportedPriKey
	}
	return priv, nil
}

func rsaDecryptPKCS8(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// RSAKey get rsa pub key and seconds ts hash.
func (s *Service) RSAKey(c context.Context) *model.RSAKey {
	return &model.RSAKey{
		Hash: TsSeconds2Hash(time.Now().Unix()),
		Key:  publicKey,
	}
}
