package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var (
	// ErrPublicKeyError public key error
	ErrPublicKeyError = errors.New("public key error")
	// ErrPrivateKeyError private key error
	ErrPrivateKeyError = errors.New("private key error")
	// ErrUnsupportedPubKey Unsupported key error
	ErrUnsupportedPubKey = errors.New("unsupported public key")
	// ErrUnsupportedPrivKey Private key error
	ErrUnsupportedPrivKey = errors.New("unsupported private key")
)

func pubKey(publicKey []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, ErrPublicKeyError
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, ErrUnsupportedPubKey
	}
	return pub, nil
}

func privKey(privateKey []byte) (*rsa.PrivateKey, error) {
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
		return nil, ErrUnsupportedPrivKey
	}
	return priv, nil
}

func rsaEncryptPKCS8(pub *rsa.PublicKey, origData []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

func rsaDecryptPKCS8(priv *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
