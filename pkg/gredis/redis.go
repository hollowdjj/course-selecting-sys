package gredis

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/hollowdjj/course-selecting-sys/pkg/logging"
	"github.com/hollowdjj/course-selecting-sys/pkg/setting"
	"time"
)

var Rdb *redis.Client

//连接redis
func SetUp() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:        setting.RedisSetting.Host,
		Password:    setting.RedisSetting.Password,
		DB:          0, // use default DB
		IdleTimeout: setting.RedisSetting.IdleTimeout,
	})
	if err := Rdb.Ping().Err(); err != nil {
		logging.Fatal("connect redis fail: %v", err)
	}
}

func Set(key string, data interface{}, expiredTime time.Duration) error {
	value, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return Rdb.Set(key, value, expiredTime).Err()
}

func Get(key string) ([]byte, error) {
	ret, err := Rdb.Get(key).Bytes()
	return ret, err
}

func Exist(key string) (bool, error) {
	//key存在返回1，不存在返回0
	v, err := Rdb.Exists(key).Result()
	if err != nil {
		return false, err
	}

	return v > 0, nil
}

func Delete(key string) (bool, error) {
	//返回被删除的key的数量
	v, err := Rdb.Del(key).Result()
	if err != nil {
		return false, err
	}

	return v > 0, nil
}

func Close() error {
	return Rdb.Close()
}
