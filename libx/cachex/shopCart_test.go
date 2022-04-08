package cachex

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShopCart(t *testing.T) {
	xredis := NewXRedis(&RedisOpts{
		Host:      ":6379",
		MaxIdle:   80,
		MaxActive: 1200,
	})

	cart := NewShopCart("app", 5, xredis)
	err := cart.Insert(15, 5)
	assert.Nil(t, err)

	// 查询不存在的商品
	quantity, err := cart.GetWith(1)
	assert.Error(t, err)
	assert.Equal(t, 0, quantity)

	quantity, err = cart.GetWith(15)
	assert.Nil(t, err)
	assert.Equal(t, 5, quantity)

	err = cart.ChangeWith(15, 1)
	assert.Nil(t, err)

	quantity, err = cart.GetWith(15)
	assert.Nil(t, err)
	assert.Equal(t, 6, quantity)

	err = cart.ChangeWith(15, -3)
	assert.Nil(t, err)

	quantity, err = cart.GetWith(15)
	assert.Nil(t, err)
	assert.Equal(t, 3, quantity)

	err = cart.Insert(2, 15)
	assert.Nil(t, err)

	err = cart.Insert(3, 99)
	assert.Nil(t, err)

	length, err := cart.Length()
	assert.Nil(t, err)
	assert.Equal(t, 3, length)

	err = cart.DeleteWith(3)
	assert.Nil(t, err)

	length, err = cart.Length()
	assert.Nil(t, err)
	assert.Equal(t, 2, length)

	all, err := cart.GetAll()
	assert.Nil(t, err)
	fmt.Println("rst", all)
}
