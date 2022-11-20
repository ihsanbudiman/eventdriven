package redis_conn

import "github.com/go-redis/redis"

type RedisConn struct {
	rdb *redis.Client
}

func NewRedisConn() *RedisConn {

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisConn{
		rdb: rdb,
	}

}

func (r *RedisConn) Get(key string) (string, error) {
	return r.rdb.Get(key).Result()
}

func (r *RedisConn) Set(key string, value interface{}, expiration int) error {
	return r.rdb.Set(key, value, 0).Err()
}
