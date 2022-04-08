package cmn

import (
	"github.com/dgrijalva/jwt-go"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestJwt(t *testing.T) {
	jj := Jwt{secret: "9988"}

	claims := jwt.MapClaims{
		"hello": "go",
	}

	encode, err := jj.Encode(claims)
	assert.Nil(t, err)
	encode = "jwt" + encode

	decode, err := jj.Decode(encode[3:])
	assert.Nil(t, err)
	err = decode.Valid()
	assert.Nil(t, err)
	assert.Equal(t, "go", decode["hello"])
}
