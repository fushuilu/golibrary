package verify

import (
	"github.com/fushuilu/golibrary"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGlobalDomainInvalid(t *testing.T) {

	Domains = []string{
		"a.com", "b.com",
	}

	err := URLsInvalid("http://a.com/abc", "http://b.com/abc")
	assert.Nil(t, err)

	err = URLsInvalid("http://a.com/abc", "http://c.com/abc")
	assert.Error(t, err)

	err = URLsInvalid("http://d.com/abc")
	assert.Error(t, err)

	data := []struct {
		Data string
		Rst  string
	}{
		{Data: "http://localhost:4201#/auth/token", Rst: "localhost:4201"},
		{Data: "https://demo.fushuilu.com/a/b/c.php", Rst: "demo.fushuilu.com"},
		{Data: "https://0768hx.com#/oauth/token", Rst: "0768hx.com"},
	}

	for _, v := range data {
		domain, err := golibrary.HttpURLHostPort(v.Data)
		assert.Nil(t, err, v.Data)
		assert.Equal(t, v.Rst, domain, v.Data)
	}
}
