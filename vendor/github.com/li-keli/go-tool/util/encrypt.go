package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"fmt"
	"github.com/satori/go.uuid"
	"strings"
)

type AesEncrypt struct {
	StrKey string
}

func (this *AesEncrypt) getKey() []byte {
	keyLen := len(this.StrKey)
	if keyLen < 16 {
		panic("res key 长度不能小于16")
	}
	arrKey := []byte(this.StrKey)
	if keyLen >= 32 {
		//取前32个字节
		return arrKey[:32]
	}
	if keyLen >= 24 {
		//取前24个字节
		return arrKey[:24]
	}
	//取前16个字节
	return arrKey[:16]
}

//加密字符串
func (this *AesEncrypt) Encrypt(strBytes []byte) ([]byte, error) {
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	encrypted := make([]byte, len(string(strBytes)))
	aesBlockEncrypter, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesEncrypter := cipher.NewCFBEncrypter(aesBlockEncrypter, iv)
	aesEncrypter.XORKeyStream(encrypted, []byte(strBytes))
	return encrypted, nil
}

//解密字符串
func (this *AesEncrypt) Decrypt(src []byte) (strBytes []byte, err error) {
	defer func() {
		//错误处理
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	key := this.getKey()
	var iv = []byte(key)[:aes.BlockSize]
	decrypted := make([]byte, len(src))
	var aesBlockDecrypter cipher.Block
	aesBlockDecrypter, err = aes.NewCipher([]byte(key))
	if err != nil {
		return []byte(""), err
	}
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(decrypted, src)
	return decrypted, nil
}

// MD5加密
func ToMd5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}

// generate UUID
func GetNewUUID() string {
	uuids, _ := uuid.NewV4()
	return strings.Replace(uuids.String(), "-", "", -1)
}
