package appx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSeedSign(t *testing.T) {

	seed := "1-15"

	ss := seedSign{}
	sign := ss.Create(seed)

	err := ss.Invalid(seed, sign)
	assert.NoError(t, err)

	err = ss.MustNotExpired(sign, 15)
	assert.NoError(t, err)
}
