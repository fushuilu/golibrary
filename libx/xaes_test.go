package cmn

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const aesKey = "1234567890123456" // 16位

func TestXAes(t *testing.T) {

	xAes, err := NewXAes(aesKey, aesKey)
	assert.Nil(t, err)

	plantText := "hello world"

	encrypt, err := xAes.Encrypt(plantText)
	assert.Nil(t, err)
	fmt.Println("encryptText:", encrypt)

	decode, err := xAes.Decrypt(encrypt)
	assert.Nil(t, err)
	assert.Equal(t, plantText, decode)
}

func TestAesEcb(t *testing.T) {

	key := []byte(aesKey)
	data := []byte("hello")

	encrypt := AesEcbEncrypt(data, key)
	fmt.Println("AesEcbEncrypt:", string(encrypt))

	decrypt := AesEcbDecrypt(encrypt, key)
	assert.Equal(t, decrypt, data)
}
