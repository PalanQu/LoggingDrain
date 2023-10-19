package loggingdrain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisPersistence struct {
	addr       string
	password   string
	db         int
	serviceKey string
	rdb        RedisClient
}

// RedisClient for mock testing
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

var _ RedisClient = &redis.Client{}

func NewRedisPersistence(addr, password string, db int, serviceKey string) *RedisPersistence {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisPersistence{
		addr:       addr,
		password:   password,
		db:         db,
		serviceKey: serviceKey,
		rdb:        rdb,
	}
}

var _ persistenceHandler = &RedisPersistence{}

func (p *RedisPersistence) Save(ctx context.Context, template *TemplateMiner) error {
	b, err := json.Marshal(template)
	if err != nil {
		return errInternal(err)
	}
	if err := p.rdb.Set(ctx, p.serviceKey, string(b), 0).Err(); err != nil {
		return errInternal(err)
	}
	return nil
}

func (p *RedisPersistence) Load(ctx context.Context) (*TemplateMiner, error) {
	val, err := p.rdb.Get(ctx, p.serviceKey).Result()
	if err != nil {
		return nil, errInternal(err)
	}
	miner := TemplateMiner{}
	if err := json.Unmarshal([]byte(val), &miner); err != nil {
		return nil, errInternal(err)
	}
	return &miner, nil
}

func (p *RedisPersistence) Subscribe(ctx context.Context) *redis.PubSub {
	return p.rdb.Subscribe(ctx, p.serviceKey)
}
