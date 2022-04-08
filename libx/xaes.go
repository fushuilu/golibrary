package cmn

import (
	"bytes"
	"crypto/aes"
	"encoding/hex"
	"errors"
	"github.com/wumansgy/goEncrypt"
)

// key: 16, 24, or 32 bytes
func AesEcbEncrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	data = PKCS7Padding(data, block.BlockSize())
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Encrypt(decrypted[bs:be], data[bs:be])
	}

	return decrypted
}

func AesEcbDecrypt(data, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	decrypted := make([]byte, len(data))
	size := block.BlockSize()

	for bs, be := 0, size; bs < len(data); bs, be = bs+size, be+size {
		block.Decrypt(decrypted[bs:be], data[bs:be])
	}

	return PKCS7UnPadding(decrypted)
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

type XAes struct {
	Ok         bool   // 是否已经设置 iv, key
	encryptKey []byte // 16字节
	encryptIv  []byte
}

func NewXAes(key, iv string) (aes XAes, err error) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		err = errors.New("key len not 16/24/32")
		return
	}
	ivLen := len(iv)
	if ivLen != 16 && ivLen != 24 && ivLen != 32 {
		err = errors.New("iv len not 16/24/32")
		return
	}

	return XAes{encryptKey: []byte(key), encryptIv: []byte(iv), Ok: true}, nil
}

func (x *XAes) Encrypt(plainText string) (string, error) {
	if plainText == "" {
		return "", nil
	}
	encryptKey, err := goEncrypt.AesCbcEncrypt([]byte(plainText), x.encryptKey, x.encryptIv)
	if err != nil {
		return "", err
	}
	v := hex.EncodeToString(encryptKey)
	return v, nil
}

func (x *XAes) Decrypt(cryptText string) (string, error) {
	if cryptText == "" {
		return "", nil
	}
	bytesK, err := hex.DecodeString(cryptText)
	if err != nil {
		return "", err
	}
	decryptK, err := goEncrypt.AesCbcDecrypt(bytesK, x.encryptKey, x.encryptIv)
	if err != nil {
		return "", err
	}
	return string(decryptK), nil
}
