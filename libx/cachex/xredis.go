package cachex

import (
	"encoding/json"
	"github.com/fushuilu/golibrary/lerror"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Cache interface {
	// Deprecated: 微信接口专用
	Get(key string) interface{} // wechat cache 中要求返回一个能转为 string 的对象
	// Deprecated: 微信接口专用
	Set(key string, val interface{}, timeout time.Duration) error
	IsExist(key string) bool
	Delete(key string) error
}

type XRedis struct {
	conn *redis.Pool
}

type RedisOpts struct {
	Host        string `yml:"host" json:"host"`
	Password    string `yml:"password" json:"password"`
	Database    int    `yml:"database" json:"database"`
	MaxIdle     int    `yml:"max_idle" json:"max_idle"`
	MaxActive   int    `yml:"max_active" json:"max_active"`
	IdleTimeout int32  `yml:"idle_timeout" json:"idle_timeout"` //second
}

func NewXRedis(opts *RedisOpts) *XRedis {
	pool := &redis.Pool{
		MaxActive:   opts.MaxActive,
		MaxIdle:     opts.MaxIdle,
		IdleTimeout: time.Second * time.Duration(opts.IdleTimeout),
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", opts.Host,
				redis.DialDatabase(opts.Database),
				redis.DialPassword(opts.Password),
			)
		},
		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := conn.Do("PING")
			return lerror.NewIf(err, "redis 连接失败")
		},
	}
	return &XRedis{conn: pool}
}

// https://www.redis.net.cn/order/3539.html
/// 以秒为单位，返回给定 key 的剩余生存时间(TTL, time to live)
func (x *XRedis) Ttl(key string) int64 {

	c := x.conn.Get()
	defer c.Close()

	num, _ := redis.Int64(c.Do("TTL", key))
	return num
}

// Redis Expire 命令用于设置 key 的过期时间。key 过期后将不再可用。
/// 设置成功返回 1 。 当 key 不存在或者不能为 key 设置过期时间时返回 0
func (x *XRedis) Expire(key string, seconds int64) bool {
	c := x.conn.Get()
	defer c.Close()

	i, _ := redis.Int64(c.Do("EXPIRE", key, seconds))
	return i == 1
}

//SetConn 设置conn
func (x *XRedis) SetConn(conn *redis.Pool) {
	x.conn = conn
}

// Deprecated: 微信接口专用
// Get 获取一个字符串类型值
func (x *XRedis) Get(key string) interface{} {
	//c := x.conn.Get()
	//defer c.Close()
	//
	//reply, err := c.Do("GET", key) // 默认 []uint8
	//rst, _ := redis.String(reply, err)
	//return rst

	conn := x.conn.Get()
	defer conn.Close()

	var data []byte
	var err error
	if data, err = redis.Bytes(conn.Do("GET", key)); err != nil {
		return nil
	}
	var reply interface{}
	if err = json.Unmarshal(data, &reply); err != nil {
		return nil
	}

	return reply
}

// Set 设置一个值
func (x *XRedis) Set(key string, val interface{}, timeout time.Duration) (err error) {
	//c := x.conn.Get()
	//defer c.Close()
	//_, err = c.Do("SETEX", key, int64(timeout/time.Second), data)
	//return
	conn := x.conn.Get()
	defer conn.Close()

	var data []byte
	if data, err = json.Marshal(val); err != nil {
		return
	}

	_, err = conn.Do("SETEX", key, int64(timeout/time.Second), data)

	return
}

func (x *XRedis) SetString(key, content string, seconds int64) error {
	c := x.conn.Get()
	defer c.Close()

	_, err := c.Do("SETEX", key, seconds, content)
	return err
}

func (x *XRedis) GetString(key string) (string, error) {
	c := x.conn.Get()
	defer c.Close()

	return redis.String(c.Do("GET", key))
}

//IsExist 判断key是否存在
func (x *XRedis) IsExist(key string) bool {
	c := x.conn.Get()
	defer c.Close()

	exists, _ := redis.Bool(c.Do("EXISTS", key))
	return exists
}

//Delete 删除
func (x *XRedis) Delete(key string) error {
	c := x.conn.Get()
	defer c.Close()

	if _, err := c.Do("DEL", key); err != nil {
		return err
	}
	return nil
}

func (x *XRedis) Int(key string) (int, error) {
	c := x.conn.Get()
	defer c.Close()

	return redis.Int(c.Do("GET", key))
}
func (x *XRedis) Conn() redis.Conn {
	return x.conn.Get()
}
