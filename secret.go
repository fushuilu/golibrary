package golibrary

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"os"
)

func HashFilepathMd5(filePath string) (string, error) {
	var returnMD5String string
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}
	return HashFileMd5(file)

}

func HashFileMd5(file *os.File) (string, error) {
	var returnMD5String string
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

func HashMultipartFile(file multipart.File) (string, error) {
	var returnMD5String string
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String = hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}

// https://medium.com/@jcox250/password-hash-salt-using-golang-b041dc94cb72
// 返回一个 60 位长度的密码
func EncryptPassword(plainPassword string) string {
	data, _ := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.MinCost)
	return string(data)
}
func ComparePassword(plainPassword, hashPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
	if err != nil {
		return errors.New("密码错误")
	}
	return nil
}

// EncryptBytes encrypts <data> using MD5 algorithms.
func EncryptBytes(data []byte) (encrypt string, err error) {
	h := md5.New()
	if _, err = h.Write([]byte(data)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}


// EncryptBytes encrypts string <data> using MD5 algorithms.
func MD5(data string) (encrypt string, err error) {
	return EncryptBytes([]byte(data))
}