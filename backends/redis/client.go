package redis

import (
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/moensch/confmgr/backends"
	"github.com/moensch/confmgr/vars"
	"log"
	"time"
)

type ConfigBackendRedisFactory struct {
	Pool *redis.Pool
}

func NewFactory() backend.ConfigBackendFactory {
	factory := &ConfigBackendRedisFactory{
		Pool: newRedisPool("tcp", ":6379"),
	}

	return factory
}

func (f *ConfigBackendRedisFactory) NewBackend() backend.ConfigBackend {
	// TODO: Make config passing work
	backend := &ConfigBackendRedis{}

	backend.Conn = f.Pool.Get()
	err := backend.Conn.Err()
	if err != nil {
		log.Printf("redis error: %s", err)
	}

	return backend
}

func newRedisPool(proto string, address string) *redis.Pool {
	log.Printf("Setting up redis pool for: %s:%s", proto, address)
	return &redis.Pool{
		MaxIdle:     5,
		MaxActive:   20,
		Wait:        true,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(proto, address)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

type ConfigBackendRedis struct {
	Conn redis.Conn
}

func (b ConfigBackendRedis) Check() error {
	return b.Conn.Err()
}

func (b ConfigBackendRedis) Close() {
	b.Conn.Close()
}

func (b ConfigBackendRedis) GetType(key string) (int, error) {
	var keytype int
	var err error

	redistype, err := redis.String(b.Conn.Do("TYPE", key))

	if err != nil {
		return keytype, err
	}

	switch redistype {
	case "none":
		return vars.TYPE_NOT_FOUND, err
	case "list":
		return vars.TYPE_LIST, err
	case "hash":
		return vars.TYPE_HASH, err
	case "string":
		return vars.TYPE_STRING, err
	default:
		return keytype, errors.New("Invalid redis key type")
	}
}

func (b ConfigBackendRedis) GetString(key string) (string, error) {
	value, err := redis.String(b.Conn.Do("GET", key))

	return value, err
}

func (b ConfigBackendRedis) SetString(key string, value string) error {
	_, err := b.Conn.Do("SET", key, value)
	if err != nil {
		return err
	}
	return err
}

func (b ConfigBackendRedis) DeleteKey(key string) error {
	_, err := b.Conn.Do("DEL", key)
	return err
}

func (b ConfigBackendRedis) GetHashField(key string, field string) (string, error) {
	value, err := redis.String(b.Conn.Do("HGET", key, field))

	return value, err
}

func (b ConfigBackendRedis) GetListIndex(key string, index int64) (string, error) {
	value, err := redis.String(b.Conn.Do("LINDEX", key, index))

	return value, err
}

func (b ConfigBackendRedis) GetHash(key string) (map[string]string, error) {
	var value map[string]string
	var err error

	value, err = redis.StringMap(b.Conn.Do("HGETALL", key))

	return value, err
}

func (b ConfigBackendRedis) SetHash(key string, value map[string]string) error {
	_, err := b.Conn.Do("DEL", key)
	if err != nil {
		return err
	}

	// Cannot easily do hmset
	for field, v := range value {
		_, err = b.Conn.Do("HSET", key, field, v)
		if err != nil {
			return err
		}
	}
	return err
}

func (b ConfigBackendRedis) ListKeys(filter string) ([]string, error) {
	if filter == "" {
		filter = "*"
	}

	value, err := redis.Strings(b.Conn.Do("KEYS", filter))

	return value, err
}

func (b ConfigBackendRedis) SetHashField(key string, field string, value string) error {
	keytype, err := b.GetType(key)
	if err != nil {
		return err
	}
	switch keytype {
	case vars.TYPE_NOT_FOUND:
		fallthrough
	case vars.TYPE_HASH:
		_, err = b.Conn.Do("HSET", key, field, value)
	default:
		return errors.New(fmt.Sprintf("Unsupported key type: %d", keytype))
	}

	return err
}

func (b ConfigBackendRedis) GetList(key string) ([]string, error) {
	var value []string
	var err error

	value, err = redis.Strings(b.Conn.Do("LRANGE", key, 0, -1))

	return value, err
}

func (b ConfigBackendRedis) SetList(key string, value []string) error {
	_, err := b.Conn.Do("DEL", key)
	if err != nil {
		return err
	}
	for _, entry := range value {
		_, err = b.Conn.Do("RPUSH", key, entry)
		if err != nil {
			return err
		}
	}
	return err
}

func (b ConfigBackendRedis) ListAppend(key string, value string) error {
	_, err := b.Conn.Do("RPUSH", key, value)
	return err
}

func (b ConfigBackendRedis) Exists(key string) (bool, error) {
	var exists bool
	var err error

	exists, err = redis.Bool(b.Conn.Do("EXISTS", key))

	return exists, err
}

func (b ConfigBackendRedis) HashFieldExists(key string, field string) (bool, error) {
	var exists bool
	var err error

	// First, check the actual key
	exists, err = b.Exists(key)
	if err != nil {
		return exists, err
	}

	// Now, check hash field
	exists, err = redis.Bool(b.Conn.Do("HEXISTS", key, field))

	return exists, err
}

func (b ConfigBackendRedis) ListIndexExists(key string, index int64) (bool, error) {
	var exists bool
	var err error

	if index < 0 {
		return false, err
	}

	// First, check the actual key
	exists, err = b.Exists(key)
	if err != nil {
		return exists, err
	}

	// Get list length
	length, err := redis.Int64(b.Conn.Do("LLEN", key))

	if err != nil {
		return exists, err
	}

	// index is zero based
	if length <= index {
		return false, err
	}
	return true, err
}
