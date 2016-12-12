// Package aredis is for initializing the redis requirements, the pool and connection
package aredis

import (
	"common/assert"
	"common/config"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/redis.v5"
)

var (
	// Client the actual pool to use with redis
	Client *redis.Client
	once   = &sync.Once{}
)

// Initialize try to create a redis pool
func Initialize() {
	once.Do(func() {
		Client = redis.NewClient(
			&redis.Options{
				Network:  config.Config.Redis.Network,
				Addr:     config.Config.Redis.Address,
				Password: config.Config.Redis.Password,
				PoolSize: config.Config.Redis.Size,
				DB:       config.Config.Redis.Database,
			},
		)
		// PING the server to make sure every thing is fine
		assert.Nil(Client.Ping().Err())
		logrus.Debug("redis is ready.")
	})
}

// StoreKey is a simple key value store with timeout
func StoreKey(key, data string, expire time.Duration) error {
	return Client.Set(key, data, expire).Err()
}

// StoreHashKey is a simple function to set hash key
func StoreHashKey(key, subkey, data string, expire time.Duration) error {
	err := Client.HSet(key, subkey, data).Err()
	if err == nil {
		err = Client.Expire(key, expire).Err()
	}

	return err
}

// GetKey Get a key from redis
func GetKey(key string, touch bool, expire time.Duration) (string, error) {
	cmd := Client.Get(key)
	if err := cmd.Err(); err != nil {
		return "", err
	}

	if touch {
		bCmd := Client.Expire(key, expire)
		if err := bCmd.Err(); err != nil {
			return "", err
		}
	}
	return cmd.Val(), nil
}

// GetHashKey return a key from a hash
func GetHashKey(key, subkey string, touch bool, expire time.Duration) (string, error) {
	cmd := Client.HGet(key, subkey)
	if err := cmd.Err(); err != nil {
		return "", err
	}
	if touch {
		bCmd := Client.Expire(key, expire)
		if err := bCmd.Err(); err != nil {
			return "", err
		}
	}
	return cmd.Val(), nil
}

// RemoveKey for removing a key in redis
func RemoveKey(key string) error {
	bCmd := Client.Del(key)
	return bCmd.Err()
}

// GetExpire return the expire of a key
func GetExpire(key string) (time.Duration, error) {
	eCmd := Client.TTL(key)
	return eCmd.Val(), eCmd.Err()
}
