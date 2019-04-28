package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

//var key = []byte("bili_account_enc")

// CBCEncrypt aes cbc encrypt
func (s *Service) CBCEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.aesKey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, s.aesKey[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

// CBCDecrypt aes cbc decrypt
func (s *Service) CBCDecrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.aesKey)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, s.aesKey[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

// PKCS5Padding padding
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS5UnPadding unpadding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
