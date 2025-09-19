package awsRedis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	config "github.com/BeeTechHub/go-common/aws/config"
	"github.com/BeeTechHub/go-common/configs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/redis/go-redis/v9"
)

// "github.com/aws/aws-sdk-go/service/elasticache"

var nilClientError = errors.New("Access redis failed because redis client nil")

type RedisClientWrapper struct {
	Client *redis.Client
}

func initClientLocal() (*RedisClientWrapper, error) {
	fmt.Println("Initializing Redis client for local connection...")

	// Connect to the local Redis server
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Default local Redis address
		Password: "",               // No password for local Redis by default
		DB:       0,                // Default DB
	})

	// Test the connection
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		fmt.Println("Failed to connect to local Redis: %v\n", err)
		return nil, err
	} else {
		fmt.Println("Connected to local Redis successfully!")
		return &RedisClientWrapper{redisClient}, nil
	}
}

func initClientAws(cacheClusterName string) (*RedisClientWrapper, error) {
	fmt.Println("Redis client start...")
	// Set up AWS session and Elasticache client
	sess := config.GetAWSSession()
	svc := elasticache.New(sess)

	clusterName := cacheClusterName
	// Get the Elasticache Redis cluster's endpoint address and port
	result, err := svc.DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{
		CacheClusterId:    aws.String(clusterName),
		ShowCacheNodeInfo: aws.Bool(true),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case elasticache.ErrCodeCacheClusterNotFoundFault:
				fmt.Println(elasticache.ErrCodeCacheClusterNotFoundFault, aerr.Error())
			case elasticache.ErrCodeInvalidParameterValueException:
				fmt.Println(elasticache.ErrCodeInvalidParameterValueException, aerr.Error())
			case elasticache.ErrCodeInvalidParameterCombinationException:
				fmt.Println(elasticache.ErrCodeInvalidParameterCombinationException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and Message from an error.
			fmt.Printf("Get redis cluster error:%s", err.Error())
		}

		return nil, err
	}

	// print the endpoint of the cluster
	if len(result.CacheClusters) <= 0 || result.CacheClusters[0] == nil ||
		len(result.CacheClusters[0].CacheNodes) <= 0 || result.CacheClusters[0].CacheNodes[0] == nil {
		errMessage := fmt.Sprintf("Missing elasticache cluster with name: %s", clusterName)
		fmt.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	endpoint := result.CacheClusters[0].CacheNodes[0].Endpoint.Address
	if endpoint == nil {
		errMessage := fmt.Sprintf("Missing elasticache cluster address with name: %s", clusterName)
		fmt.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	port := result.CacheClusters[0].CacheNodes[0].Endpoint.Port
	if port == nil {
		errMessage := fmt.Sprintf("Missing elasticache cluster port with name: %s", clusterName)
		fmt.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	// Create a Redis client and connect to the cluster
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", *endpoint, *port),
	})

	_, err = redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("Failed to ping Redis: %w", err)
	}

	return &RedisClientWrapper{redisClient}, nil
}

func InitRedis(cacheClusterName string) (*RedisClientWrapper, error) {
	if configs.GetCacheHost() == "local" {
		return initClientLocal()
	} else {
		return initClientAws(cacheClusterName)
	}
}

func (redisClient RedisClientWrapper) SetDataToCache(key string, value string, exprire time.Duration) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.Set(context.Background(), key, value, exprire).Result()
	if err != nil {
		return err
	}

	return nil
}

func (redisClient RedisClientWrapper) GetValueFromKey(key string) (*string, error) {
	if redisClient.Client == nil {
		return nil, nilClientError
	}

	data, err := redisClient.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, nil
}

func (redisClient RedisClientWrapper) DeleteDataFromKeys(keys []string) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	pipe := redisClient.Client.Pipeline()
	for _, key := range keys {
		if len(key) <= 0 {
			continue
		}
		pipe.Del(context.Background(), key)
	}

	_, _err := pipe.Exec(context.Background())
	if _err != nil {
		return _err
	}

	return nil
}

func (redisClient RedisClientWrapper) FlushAllAsync() error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.FlushAllAsync(context.Background()).Result()
	return err
}

func (redisClient RedisClientWrapper) SetNXDataToCache(key string, value string, exprire time.Duration) (bool, error) {
	if redisClient.Client == nil {
		return false, nilClientError
	}

	result, err := redisClient.Client.SetNX(context.Background(), key, value, exprire).Result()
	if err != nil {
		return false, err
	}

	return result, nil
}

func (redisClient RedisClientWrapper) DeleteDataFromKey(key string) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.Del(context.Background(), key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (redisClient RedisClientWrapper) ZScore(key string, member string) (*float64, error) {
	if redisClient.Client == nil {
		return nil, nilClientError
	}

	data, err := redisClient.Client.ZScore(context.Background(), key, member).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &data, nil
}

func (redisClient RedisClientWrapper) ZAdd(key string, member string, score float64) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.ZAdd(context.Background(), key, redis.Z{Score: score, Member: member}).Result()
	if err != nil {
		return err
	}

	return nil
}

func (redisClient RedisClientWrapper) ZRem(key string, members ...interface{}) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.ZRem(context.Background(), key, members).Result()
	if err != nil {
		return err
	}

	return nil
}

func (redisClient RedisClientWrapper) ZRangeByScore(key string, min float64, max float64) ([]string, error) {
	if redisClient.Client == nil {
		return nil, nilClientError
	}

	components, err := redisClient.Client.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{Min: strconv.FormatFloat(min, 'f', 0, 64), Max: strconv.FormatFloat(max, 'f', 0, 64)}).Result()
	if err != nil {
		return nil, err
	}

	return components, nil
}

func (redisClient RedisClientWrapper) PublishToChannel(channelName string, message interface{}) error {
	if redisClient.Client == nil {
		return nilClientError
	}
	// Publish a message
	err := redisClient.Client.Publish(context.Background(), channelName, message).Err()
	if err != nil {
		return err
	}

	return nil
}

func (redisClient RedisClientWrapper) SubscribeChannel(channelName string) (*redis.PubSub, error) {
	if redisClient.Client == nil {
		return nil, nilClientError
	}

	return redisClient.Client.Subscribe(context.Background(), channelName), nil
}

func (redisClient RedisClientWrapper) Unlink(keys ...string) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.Unlink(context.Background(), keys...).Result()
	return err
}

func (redisClient RedisClientWrapper) AddToSet(folder string, keys ...any) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.SAdd(context.Background(), folder, keys...).Result()
	return err
}

func (redisClient RedisClientWrapper) RemoveFromSet(folder string, keys ...any) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.SRem(context.Background(), folder, keys...).Result()
	return err
}

func (redisClient RedisClientWrapper) RemoveASet(folder string) error {
	if redisClient.Client == nil {
		return nilClientError
	}

	keys, err := redisClient.Client.SMembers(context.Background(), folder).Result()
	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	_keys := make([]interface{}, len(keys))
	for i, key := range keys {
		_keys[i] = key
	}

	_, err = redisClient.Client.Del(context.Background(), folder).Result()
	if err != nil {
		return err
	}

	_, err = redisClient.Client.Unlink(context.Background(), keys...).Result()
	if err != nil {
		return err
	}

	/*_, err = redisClient.Client.SRem(context.Background(), folder, _keys...).Result()
	if err != nil {
		return err
	}*/

	return nil
}

func (redisClient RedisClientWrapper) FlushDBAsync() error {
	if redisClient.Client == nil {
		return nilClientError
	}

	_, err := redisClient.Client.FlushDBAsync(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

func (redisClient RedisClientWrapper) GetSetKeys(folder string) ([]string, error) {
	if redisClient.Client == nil {
		return nil, nilClientError
	}

	return redisClient.Client.SMembers(context.Background(), folder).Result()
}
