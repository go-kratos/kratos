package service

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"testing"
)

var key = []byte("bili_account_enc")

func TestService_cbcEncrypt(t *testing.T) {
	tel := []byte("18612340123")
	cipherText, _ := cbcEncrypt(tel)
	fmt.Println(len(cipherText))
	fmt.Println(cipherText)
}

func TestService_cbcDecrypt(t *testing.T) {
	cipherText := []byte{115, 201, 179, 163, 254, 77, 59, 220, 62, 178, 19, 241, 165, 28, 249, 249}
	fmt.Println(string(cipherText))
	tel, _ := cbcDecrypt(cipherText)
	fmt.Println(string(tel))
}

func TestService_cfbEncrypt(t *testing.T) {
	tel := []byte("18612340123")
	cipherText, _ := cfbEncrypt(tel)
	fmt.Println(len(cipherText))
	fmt.Println(cipherText)
}

func TestService_cfbDecrypt(t *testing.T) {
	cipherText := []byte{48, 35, 4, 24, 80, 204, 171, 226, 219, 74, 16, 95, 138, 184, 249, 205}
	fmt.Println(string(cipherText))
	tel, _ := cfbDecrypt(cipherText)
	fmt.Println(string(tel))
}

func cbcEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func cbcDecrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func cfbEncrypt(origData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	cfb := cipher.NewCFBEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	cfb.XORKeyStream(crypted, origData)
	return crypted, nil
}

func cfbDecrypt(crypted []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	cfb := cipher.NewCFBDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	cfb.XORKeyStream(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}
