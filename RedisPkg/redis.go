package redisPkg
import (
	"context"
	"time"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)


func InitRedis() (*redis.Client,context.Context){

	rdb := redis.NewClient(&redis.Options{
		Addr:	os.Getenv("IP_REIS"),
		Password: os.Getenv("PASSWORD"),
		DB:	0,
	})
	ctx := context.Background()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v",err)
	}

	return rdb,ctx
}


// use val.getVal() to get vaule and val.Err() to get err
func	GetValueFromKey(rdb *redis.Client,ctx context.Context,key string) (*redis.StringCmd){

	val := rdb.Get(ctx,key)
	return val
}


func	SetNew(rdb *redis.Client,ctx context.Context,key string,value interface{}) error{
	// 1 day exp time
	exp := time.Duration(86400 * 30 * time.Second)
	return rdb.Set(ctx,key,value,exp).Err()
}

func	RemoveFromKey(rdb *redis.Client,ctx context.Context,key string) error{
		// rdb.Del(ctx, key).Result()
		return rdb.Del(ctx, key).Err()
}
