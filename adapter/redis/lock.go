package redis

import (
	"context"
	"sync"
	"time"
)

const times = 10

var mutex sync.Mutex

func Lock(key string) error {
	mutex.Lock()
	defer mutex.Unlock()
	err := GetDefaultRedisClient().SetNX(context.Background(), key, true, 10*time.Second).Err()
	if err != nil {
		for i := 0; i < times; i++ {
			time.Sleep(2 * time.Millisecond)
			err = GetDefaultRedisClient().SetNX(context.Background(), key, true, 10*time.Second).Err()
			if err == nil {
				return nil
			}
		}
	}
	return nil
}

func Unlock(key string) error {
	mutex.Lock()
	defer mutex.Unlock()
	GetDefaultRedisClient().Del(context.Background(), key)
	return nil
}
