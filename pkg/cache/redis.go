package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Stream-I-T-Consulting/stream-http-service-go/config"
	"github.com/Stream-I-T-Consulting/stream-http-service-go/pkg/utils/color"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

type Cache struct {
	redis *redis.Client
	// redisCluster *redis.ClusterClient
	tags []string

	prefix  string
	expired time.Duration
}

func Initialize(config *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Redis.RedisHost, config.Redis.RedisPort),
		Username: config.Redis.RedisUsername,
		Password: config.Redis.RedisPassword,
	})
}

func NewCacher(redisClient *redis.Client, opts ...Option) *Cache {
	if !fiber.IsChild() {
		log.Println("Redis connecting...")
	}

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	if !fiber.IsChild() {
		log.Println("Redis client connected", color.Format(color.GREEN, "successfully!"))
	}

	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}

	return &Cache{
		redis:   redisClient,
		prefix:  o.prefix,
		expired: o.expired,
	}
}

type Options struct {
	prefix  string
	expired time.Duration
}

type Option func(*Options)

func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.prefix = prefix
	}
}

func WithExpired(exp time.Duration) Option {
	return func(o *Options) {
		o.expired = exp
	}
}

func (c *Cache) Tag(tag ...string) *Cache {
	c.tags = tag
	return c
}

func (c *Cache) Get(ctx context.Context, key string, val interface{}) error {
	if len(c.prefix) > 0 {
		key = c.prefix + ":" + key
	}

	jsonStr, err := c.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		fmt.Println("c.redis.Get err:", err)
		return err
	}

	err = json.Unmarshal([]byte(jsonStr), &val)
	if err != nil {
		fmt.Println("json.Unmarshal err:", err)
		return err
	}

	return nil
}

func (c *Cache) Set(ctx context.Context, key string, val interface{}) error {
	if len(c.prefix) > 0 {
		key = c.prefix + ":" + key
	}

	_, err := c.redis.TxPipelined(ctx, func(p redis.Pipeliner) error {
		for _, v := range c.tags {
			err := p.SAdd(ctx, c.prefix+":"+v, key).Err()
			if err != nil {
				fmt.Println(fmt.Errorf("p.SAdd err %v", err))
				return err
			}
		}

		value, err := json.Marshal(val)
		if err != nil {
			fmt.Println("json.Marshal err:", err)
			return err
		}

		err = p.Set(ctx, key, string(value), c.expired).Err()
		if err != nil {
			fmt.Println("p.Set err:", err)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Cache) Flush(ctx context.Context) error {
	_, err := c.redis.TxPipelined(ctx, func(p redis.Pipeliner) error {
		for _, v := range c.tags {
			members, err := c.redis.SMembers(ctx, c.prefix+":"+v).Result()
			if err != nil {
				fmt.Println("c.redis.SMembers err:", err)
				return err
			}

			if len(members) > 0 {
				err = p.Del(ctx, members...).Err()
				if err != nil {
					fmt.Println("p.Del err:", err)
					return err
				}
			}

			err = p.Del(ctx, c.prefix+":"+v).Err()
			if err != nil {
				fmt.Println("p.Del err:", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Close() {
	c.redis.Close()
}
