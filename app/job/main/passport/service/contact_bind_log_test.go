package service

import (
	"context"
	"testing"
	"time"

	"encoding/json"
	"go-common/app/job/main/passport/model"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"
)

var (
	pCfg = &databus.Config{
		Key:          "dbe67e6a4c36f877",
		Secret:       "8c775ea242caa367ba5c876c04576571",
		Group:        "Test1-MainCommonArch-P",
		Topic:        "test1",
		Action:       "pub",
		Name:         "databus",
		Proto:        "tcp",
		Addr:         "172.18.33.50:6205",
		Active:       10,
		Idle:         5,
		DialTimeout:  xtime.Duration(time.Second),
		WriteTimeout: xtime.Duration(time.Second),
		ReadTimeout:  xtime.Duration(time.Second),
		IdleTimeout:  xtime.Duration(time.Minute),
	}
)

func testPub(t *testing.T, d *databus.Databus) {
	tel := model.TelBindLog{ID: 2, Mid: 88883, Tel: "18817352650", Timestamp: 1500022511}
	da, _ := json.Marshal(&tel)
	c := &model.BMsg{Action: "insert", Table: "aso_telephone_bind_log", New: da}

	if err := d.Send(context.Background(), "test", c); err != nil {
		t.Errorf("d.Send(test) error(%v)", err)
	}
}

func TestDatabus(t *testing.T) {
	d := databus.New(pCfg)
	testPub(t, d)
	testPub(t, d)
	testPub(t, d)
	d.Close()
}

//var aesBlock, _ = aes.NewCipher([]byte("1234567890abcdef"))

//func TestEncode(t *testing.T) {
//	for a := 0; a < 1000; a++ {
//		go enconde()
//	}
//	time.Sleep(10000 * time.Second)
//
//}
//
//func enconde() {
//	for i := 0; i < 100; i++ {
//		key := []byte("1234567890abcdef")
//		origData := []byte(strconv.Itoa(rand.Intn(100)))
//		blockSize := aesBlock.BlockSize()
//		origData = PKCS7Padding(origData, blockSize)
//		blockMode := cipher.NewCBCEncrypter(aesBlock, key[:blockSize])
//		crypted := make([]byte, len(origData))
//		blockMode.CryptBlocks(crypted, origData)
//		fmt.Println(base64.StdEncoding.EncodeToString(crypted))
//	}
//}

//func TestDeode(t *testing.T) {
//	key := []byte("1234567890abcdef")
//	b,_:=base64.StdEncoding.DecodeString("29YQhqBb/J2XiBAj6bP3Zg==");
//	s,_:=AesDecrypt(b,key)
//	fmt.Print(string(s))
//}
//
//
//func AesDecrypt(crypted, key []byte) ([]byte, error) {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return nil, err
//	}
//	blockSize := block.BlockSize()
//	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
//	origData := make([]byte, len(crypted))
//	blockMode.CryptBlocks(origData, crypted)
//	origData = PKCS7UnPadding(origData)
//	return origData, nil
//}
//
//
//func PKCS7UnPadding(origData []byte) []byte {
//	length := len(origData)
//	unpadding := int(origData[length-1])
//	return origData[:(length - unpadding)]
//}
