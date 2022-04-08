package cachex

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gomodule/redigo/redis"
)

type ShopCart struct {
	name   string
	userId int64
	redis  *XRedis
}

func NewShopCart(name string, userId int64, redis *XRedis) ShopCart {
	return ShopCart{
		name:   name,
		userId: userId,
		redis:  redis,
	}
}

func (sc *ShopCart) userKey() string {
	return fmt.Sprintf("card:%s:%d", sc.name, sc.userId)
}

func (sc *ShopCart) goodsKey(goodsId int64) string {
	return fmt.Sprintf("card:%s:%d:%d", sc.name, sc.userId, goodsId)
}

// 将商品添加到购物车
func (sc *ShopCart) Insert(goodsId int64, quantity int) error {
	c := sc.redis.Conn()
	defer c.Close()

	_, err := c.Do("HMSET", sc.userKey(), sc.goodsKey(goodsId), quantity)
	return err
}

// 修改购物车中商品的数量
func (sc *ShopCart) ChangeWith(goodsId int64, inc int) error {
	c := sc.redis.Conn()
	defer c.Close()

	_, err := c.Do("HINCRBY", sc.userKey(), sc.goodsKey(goodsId), inc)
	return err
}

// 清空购物车
func (sc *ShopCart) Delete() error {
	c := sc.redis.Conn()
	defer c.Close()

	_, err := c.Do("DEL", sc.userKey())
	return err
}

// 删除一个商品
func (sc *ShopCart) DeleteWith(goodsId int64) error {
	c := sc.redis.Conn()
	defer c.Close()

	_, err := c.Do("HDEL", sc.userKey(), sc.goodsKey(goodsId))
	return err
}

// 获取商品数量
func (sc *ShopCart) GetWith(goodsId int64) (int, error) {
	c := sc.redis.Conn()
	defer c.Close()

	return sc.int(c.Do("HMGET", sc.userKey(), sc.goodsKey(goodsId)))
}

func (sc *ShopCart) int(doRst interface{}, err error) (int, error) {

	switch doRst.(type) {
	case []interface{}:
		if len(doRst.([]interface{})) > 0 {
			return redis.Int(doRst.([]interface{})[0], err)
		}
	case interface{}:
		return redis.Int(doRst, err)
	default:
		fmt.Printf("rst:%+v; t:%s\n", doRst, reflect.TypeOf(doRst))
	}
	return 0, nil
}

type CartItem struct {
	GoodsId  int64
	Quantity int
}

var numReg = regexp.MustCompile("(\\d+)$")

func (sc *ShopCart) GetAll() ([]CartItem, error) {
	c := sc.redis.Conn()
	defer c.Close()

	if rst, err := c.Do("HGETALL", sc.userKey()); err != nil {
		return nil, err
	} else {
		slices, _ := redis.IntMap(rst, nil)
		items := make([]CartItem, len(slices))
		index := 0
		for key, q := range slices {
			goodsId, err := strconv.ParseInt(numReg.FindString(key), 10, 64)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("无法从 key 中提取出商品 id:%s", key))
			}
			items[index] = CartItem{Quantity: q, GoodsId: goodsId}
			index += 1
		}
		return items, nil
	}
}

// 商品数量
func (sc *ShopCart) Length() (int, error) {
	c := sc.redis.Conn()
	defer c.Close()

	return redis.Int(c.Do("HLEN", sc.userKey()))
}
