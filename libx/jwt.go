package cmn

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type Jwt struct {
	secret string
}

func NewJwt(secret string) Jwt {
	return Jwt{secret: secret}
}
// Authorization: Bearer <token>
func (j *Jwt) Encode(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// https://godoc.org/github.com/dgrijalva/jwt-go#example-Parse--Hmac
// https://self-issued.info/docs/draft-ietf-oauth-json-web-token.html
func (j *Jwt) Decode(content string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(content, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return jwt.MapClaims{}, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return jwt.MapClaims{}, errors.New("登录凭证验证失败")
}
