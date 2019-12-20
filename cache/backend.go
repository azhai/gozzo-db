package cache

import (
	"fmt"
	"strconv"

	"github.com/azhai/gozzo-db/schema"
	"github.com/azhai/gozzo-utils/common"
	"github.com/azhai/gozzo-utils/rdspool"
	"github.com/gomodule/redigo/redis"
)

// 缓存接口
type ICacheBackend interface {
	Connect(params schema.ConnParams) error
	Close() error
	ClearAll() error
	GetName() string
	GetTimeout() int
	Set(key string, value interface{}) (int, error)
	Get(key string) (interface{}, error)
	GetInt(key string) (int, error)
	GetString(key string) (string, error)
	GetAll() (interface{}, error)
	AddFlash(messages ...string) (int, error)
	GetFlashes(n int) ([]string, error)
}

// Redis哈希表缓存
type RedisBackend struct {
	Name    string
	Timeout int
	*rdspool.RedisHash
}

// 连接Redis
func NewRedisBackend(name string, timeout int) *RedisBackend {
	return &RedisBackend{Name: name, Timeout: timeout}
}

// 连接Redis
func ConnectRedisPool(params schema.ConnParams) *rdspool.RedisPool {
	addr := params.Concat(params.Host, params.StrPort())
	db, _ := strconv.Atoi(params.Database)
	return rdspool.NewRedisPool(addr, params.Password, db)
}

func (b *RedisBackend) SetRedisInst(pool rdspool.Redis) error {
	b.RedisHash = rdspool.NewRedisHash(pool, b.Name, b.Timeout)
	return nil
}

func (b *RedisBackend) Connect(params schema.ConnParams) error {
	return b.SetRedisInst(ConnectRedisPool(params))
}

func (b *RedisBackend) Close() error {
	return b.RedisHash.Inst.Close()
}

func (b *RedisBackend) ClearAll() error {
	_, err := b.RedisHash.DoWith("DEL")
	return err
}

func (b *RedisBackend) GetName() string {
	return b.Name
}

func (b *RedisBackend) AddFlash(messages ...string) (int, error) {
	key := fmt.Sprintf("flash:%s", b.Name)
	args := append([]interface{}{key}, common.StrToList(messages)...)
	return redis.Int(b.RedisHash.Inst.Do("RPUSH", args...))
}

// 数量n为最大取出多少条消息，-1表示所有消息
func (b *RedisBackend) GetFlashes(n int) ([]string, error) {
	key := fmt.Sprintf("flash:%s", b.Name)
	return redis.Strings(b.RedisHash.Inst.Do("LRANGE", key, 0, n))
}
