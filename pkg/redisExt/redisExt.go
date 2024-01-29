package redisExt

import (
	"boilerplate-service/pkg/newRelicExt"
	"context"

	"fmt"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/redis/go-redis/v9"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"

	nrredis "github.com/newrelic/go-agent/v3/integrations/nrredis-v9"
)

type IRedisExt interface {
	Close() error
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Ping(ctx context.Context) *redis.StatusCmd

	// Redsync
	NewMutex(name string, options ...redsync.Option) *redsync.Mutex
}

type redisExt struct {
	client *redis.Client
	rs     *redsync.Redsync
}

func New(host, port, password string, db int) (IRedisExt, error) {
	opts := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	}
	client := redis.NewClient(opts)

	client.AddHook(nrredis.NewHook(opts))

	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}

	rs := redsync.New(goredis.NewPool(client))

	return &redisExt{client, rs}, nil
}

func (r *redisExt) Close() error {
	return r.client.Close()
}

func (r *redisExt) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	ctx = newrelic.NewContext(ctx, newRelicExt.GetTxnFromCtx(ctx))
	return r.client.Del(ctx, keys...)
}

func (r *redisExt) Get(ctx context.Context, key string) *redis.StringCmd {
	ctx = newrelic.NewContext(ctx, newRelicExt.GetTxnFromCtx(ctx))
	return r.client.Get(ctx, key)
}

func (r *redisExt) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	ctx = newrelic.NewContext(ctx, newRelicExt.GetTxnFromCtx(ctx))
	return r.client.Set(ctx, key, value, expiration)
}

func (r *redisExt) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd {
	ctx = newrelic.NewContext(ctx, newRelicExt.GetTxnFromCtx(ctx))
	return r.client.SetNX(ctx, key, value, expiration)
}

func (r *redisExt) NewMutex(name string, options ...redsync.Option) *redsync.Mutex {
	return r.rs.NewMutex(name, options...)
}

func (r *redisExt) Ping(ctx context.Context) *redis.StatusCmd {
	return r.client.Ping(ctx)
}
