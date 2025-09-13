package redisclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	rdc "github.com/go-redis/cache/v8"
	rd "github.com/go-redis/redis/v8"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"go.uber.org/zap"
)

type redisStruct struct {
	Url    string
	Port   int
	Client *rd.Client
	Cache  *rdc.Cache
}

func newRedisClient(url string, port int) RedisClientItf {

	redisClient := rd.NewClient(&rd.Options{
		Addr:     fmt.Sprintf("%s:%d", url, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	//check connection
	pong, err := redisClient.Ping(context.TODO()).Result()
	if err != nil {
		logger.Log().Error("redis connection error: ", zap.Error(err))
		return nil
	}
	logger.Log().Info("redis connected: ", zap.String("result", pong))

	redisCache := rdc.New(&rdc.Options{
		Redis: redisClient,
	})

	redisObj := &redisStruct{
		Url:    url,
		Port:   port,
		Client: redisClient,
		Cache:  redisCache,
		//Context: ctx,
	}

	return redisObj
}

// func CloseRedis(ctx context.Context) error {
// 	belogger.Log(ctx).Info("Closing redis connection START")
// 	defer belogger.Log(ctx).Info("Closing redis connection END")
// 	if redisObj != nil {
// 		return redisObj.Client.Close()
// 	}
// 	return nil
// }

func (obj *redisStruct) SetValue(ctx context.Context, key string, value interface{}, timeout int, writeIfNotSet bool) error {
	cacheItem := &rdc.Item{
		Ctx:   ctx,
		Key:   key,
		Value: value,
		TTL:   time.Duration(timeout) * time.Millisecond,
		SetNX: writeIfNotSet,
	}
	err := obj.Cache.Set(cacheItem)
	return err
}

func (obj *redisStruct) GetValue(ctx context.Context, key string, value interface{}) error {
	return obj.Cache.Get(ctx, key, &value)
}

func (obj *redisStruct) DeleteKey(ctx context.Context, key string) error {
	return obj.Cache.Delete(ctx, key)
}

func (obj *redisStruct) GetTTL(ctx context.Context, key string) int {
	dur := obj.Client.TTL(ctx, key)
	if dur.Err() != nil {
		return 0
	}
	return int(dur.Val().Milliseconds())
}

func (obj *redisStruct) KeyExists(ctx context.Context, key string) bool {
	return obj.Cache.Exists(ctx, key)
}

func (obj *redisStruct) SetRedisHash(ctx context.Context, key string, kvpairs map[string]string) error {
	if kvpairs == nil {
		return fmt.Errorf("nothing to hash")
	}
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("key cannot be blank")
	}
	for k, v := range kvpairs {
		val := obj.Client.HSet(ctx, key, k, v)
		if val.Err() != nil {
			return val.Err()
		}
	}
	return nil
}

func (obj *redisStruct) GetRedisHashValue(ctx context.Context, key string, args ...string) (map[string]string, error) {
	if strings.TrimSpace(key) == "" {
		return nil, fmt.Errorf("key cannot be blank")
	}
	field := ""
	if len(args) > 0 {
		field = args[0]
	}
	if strings.TrimSpace(field) == "" {
		kmap := obj.Client.HGetAll(ctx, key)
		rmap, err := kmap.Result()
		if err != nil {
			return nil, err
		}
		return rmap, nil
	} else {
		kval := obj.Client.HGet(ctx, key, field)
		rval, err := kval.Result()
		if err != nil {
			return nil, err
		}
		m := make(map[string]string)
		m[field] = rval
		return m, nil
	}
}

func (obj *redisStruct) DeleteRedisHash(ctx context.Context, key string, fields ...string) error {
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("key cannot be blank")
	}

	kval := obj.Client.HDel(ctx, key, fields...)
	_, err := kval.Result()
	if err != nil {
		return err
	}
	return nil
}
