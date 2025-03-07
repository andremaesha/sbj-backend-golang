package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Database interface {
	Table(string) Table
	Client() Client
}

type Table interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string, time.Duration) error
	HashSet(context.Context, int, string, any) error
	Del(context.Context, string) (int64, error)
	Exists(context.Context, string) (int64, error)
	Expire(context.Context, string, time.Duration) (bool, error)
}

type Client interface {
	Database() Database
	Disconnect() error
	Ping(context.Context) error
}

type redisClient struct {
	client *redis.Client
}

type redisDatabase struct {
	client *redis.Client
}

type redisTable struct {
	client *redis.Client
	prefix string
}

func NewClient(addr, username, password string, db int) (Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Username: username,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &redisClient{client: client}, nil
}

func (rc *redisClient) Database() Database {
	return &redisDatabase{client: rc.client}
}

func (rc *redisClient) Disconnect() error {
	return rc.client.Close()
}

func (rc *redisClient) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}

func (rd *redisDatabase) Table(prefix string) Table {
	return &redisTable{client: rd.client, prefix: prefix}
}

func (rd *redisDatabase) Client() Client {
	return &redisClient{client: rd.client}
}

func (rt *redisTable) Get(ctx context.Context, key string) (string, error) {
	return rt.client.Get(ctx, rt.prefix+key).Result()
}

func (rt *redisTable) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	return rt.client.Set(ctx, rt.prefix+key, value, expiration).Err()
}

func (rt *redisTable) HashSet(ctx context.Context, expire int, sessionId string, data any) error {
	timeExpire := time.Minute * time.Duration(expire)

	err := rt.client.HSet(ctx, rt.prefix+sessionId, data).Err()
	if err != nil {
		return err
	}

	return rt.client.Expire(ctx, rt.prefix+sessionId, timeExpire).Err()
}

func (rt *redisTable) Del(ctx context.Context, key string) (int64, error) {
	return rt.client.Del(ctx, rt.prefix+key).Result()
}

func (rt *redisTable) Exists(ctx context.Context, key string) (int64, error) {
	return rt.client.Exists(ctx, rt.prefix+key).Result()
}

func (rt *redisTable) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return rt.client.Expire(ctx, rt.prefix+key, expiration).Result()
}
