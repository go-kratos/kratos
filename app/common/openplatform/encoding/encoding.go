package encoding

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

//EncryptConfig 加密配置
type EncryptConfig struct {
	Key string
	IV  string
}

//JSON JSON格式化增加对[]byte的支持
func JSON(o interface{}) string {
	var v []byte
	switch o1 := o.(type) {
	case [][]byte:
		o2 := make([]interface{}, len(o1))
		for i := 0; i < len(o1); i++ {
			o2[i] = string(o1[i])
		}
		v, _ = json.Marshal(o2)
	case []interface{}:
		for i := 0; i < len(o1); i++ {
			if b, ok := o1[i].([]byte); ok {
				o1[i] = string(b)
			}
		}
		v, _ = json.Marshal(o1)
	default:
		v, _ = json.Marshal(o)
	}
	if len(v) > 0 {
		return string(v)
	}
	return ""
}

//填充到BlockSize整数倍长度，如果正好就是对的长度，再多填充一个BlockSize长度
func pad(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

//去除填充的字节
func unpad(src []byte) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return []byte{0}, nil
	}
	unpadding := int(src[length-1])
	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}
	return src[:(length - unpadding)], nil
}

//Encrypt 加密
func Encrypt(text string, c *EncryptConfig) (string, error) {
	if "" == text {
		return "", nil
	}
	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return "", err
	}
	msg := pad([]byte(text))
	ciphertext := make([]byte, len(msg))
	mode := cipher.NewCBCEncrypter(block, []byte(c.IV))
	mode.CryptBlocks(ciphertext, msg)
	finalMsg := base64.StdEncoding.EncodeToString(ciphertext)
	return finalMsg, nil
}

//Decrypt 解密
func Decrypt(text string, c *EncryptConfig) (string, error) {
	if "" == text {
		return "", nil
	}
	block, err := aes.NewCipher([]byte(c.Key))
	if err != nil {
		return "", err
	}
	decodedMsg, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}
	if (len(decodedMsg) % aes.BlockSize) != 0 {
		return "", errors.New("blocksize must be multipe of decoded message length")
	}
	msg := decodedMsg
	mode := cipher.NewCBCDecrypter(block, []byte(c.IV))
	mode.CryptBlocks(msg, msg)
	unpadMsg, err := unpad(msg)
	if err != nil {
		return "", err
	}
	return string(unpadMsg), nil
}
