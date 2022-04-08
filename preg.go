package golibrary

import (
	"errors"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// 用户名 6~30 位
var usernameReg = regexp.MustCompile(`^[a-z][a-zA-Z0-9]{5,29}$`)

// 密码 6~30 位
var passwordReg = regexp.MustCompile(`^[a-zA-Z0-9.,]{6,100}$`)

// 是否为用户名
func IsUserName(username string) bool {
	return usernameReg.MatchString(username)
}

func IsPassword(password string) bool {
	return passwordReg.MatchString(password)
}

// 带区号的号码

func IsPhone(phone string) bool {
	return IsMobilePhone(phone) || IsZonePhone(phone)
}

var mobilePhoneReg = regexp.MustCompile(`^[1-9]\d{7,10}$`)

// 本地号码，不带区号
func IsMobilePhone(phone string) bool {
	if phone == "" {
		return false
	}
	return mobilePhoneReg.MatchString(phone)
}

var zonePhoneReg = regexp.MustCompile(`^[1-9]{2,5}-[0-9]\d{7,11}$`)
var zonePhoneReg2 = regexp.MustCompile(`^\([1-9]{2,5}\)[0-9]\d{7,11}$`) //(852)00000000

// 固话
var zoneTelReg = regexp.MustCompile(`^\d{4,5}-\d{7,8}$`)

// 带区号的
func IsZonePhone(phone string) bool {
	return zonePhoneReg.MatchString(phone) || zonePhoneReg2.MatchString(phone)
}

func IsZoneTel(tel string) bool {
	return zoneTelReg.MatchString(tel)
}
func IsZonePhoneOrTel(phone string) bool {
	return IsPhone(phone) || IsZoneTel(phone)
}

func PhoneExplode(zonePhone string) (zone, phone string, err error) {
	if IsMobilePhone(zonePhone) {
		return "", zonePhone, nil
	} else if IsZonePhone(zonePhone) {
		rst := strings.Split(zonePhone, "-")
		return rst[0], rst[1], nil
	}
	return "", "", errors.New("未知的手机号")
}

var emailReg = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmail(email string) bool {
	if email == "" {
		return false
	}
	return emailReg.MatchString(email)
}

func IsAccount(account string) bool {
	return IsEmail(account) || IsPhone(account)
}

func MustAccount(account string) (err error) {
	if IsEmail(account) || IsPhone(account) {
		return nil
	}
	return errors.New("账号格式错误:只支持邮箱或手机号")
}

func IsUrl(text string) bool {
	_, err := url.ParseRequestURI(text)
	if err != nil {
		return false
	}
	u, err := url.Parse(text)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}

var idReg = regexp.MustCompile(`^\d+$`)

func IsNumberString(text string) bool {
	return idReg.MatchString(text)
}

var letterReg = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

func IsLetter(text string) bool {
	return letterReg.MatchString(text)
}

var preNumberReg = regexp.MustCompile("^(\\d+)")

func GetPrefixNumber(text string) int {
	if text == "" {
		return 0
	}
	date, _ := strconv.Atoi(preNumberReg.FindString(text))
	return date
}
