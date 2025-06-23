package awsRedis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	config "github.com/BeeTechHub/go-common/aws/config"
	"github.com/BeeTechHub/go-common/configs"
	"github.com/BeeTechHub/go-common/logger"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elasticache"
	"github.com/redis/go-redis/v9"
)

// "github.com/aws/aws-sdk-go/service/elasticache"

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

func initClientAws(cacheClusterName, cacheUrl string) (*RedisClientWrapper, error) {
	fmt.Println("Redis client start...")
	// Set up AWS session and Elasticache client
	sess := config.GetAWSSession()
	svc := elasticache.New(sess)

	clusterName := cacheClusterName
	// Get the Elasticache Redis cluster's endpoint address and port
	result, err := svc.DescribeCacheClusters(&elasticache.DescribeCacheClustersInput{
		CacheClusterId: aws.String(clusterName),
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
			fmt.Println("Get redis cluster error:%s", err.Error())
		}
	}

	// print the endpoint of the cluster
	if len(result.CacheClusters) <= 0 || result.CacheClusters[0] == nil {
		errMessage := fmt.Sprintf("Missing elasticache cluster with name: %s", clusterName)
		fmt.Println(errMessage)
		return nil, errors.New(errMessage)
	}

	// Create a Redis client and connect to the cluster
	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cacheUrl, 6379),
	})

	return &RedisClientWrapper{redisClient}, nil
}

func InitRedis(cacheClusterName, cacheUrl string) (*RedisClientWrapper, error) {
	if configs.GetCacheHost() == "local" {
		return initClientLocal()
	} else {
		return initClientAws(cacheClusterName, cacheUrl)
	}
}

func (redisClient *RedisClientWrapper) SetDataToCache(key string, value string, exprire time.Duration) error {
	if redisClient == nil {
		return errors.New("get value from redis error because redis client nil")
	}
	logger.Debugf("set data to redis with key:%s value:%s", key, value)

	result, err := redisClient.Client.Set(context.Background(), key, value, exprire).Result()
	if err != nil {
		logger.Warnf("Set Data To redis Cache error:%s", err.Error())
		return err
	}
	logger.Infof("SetDataToCache sucessfully with key:%s result:%s expireTime:%s", key, result, exprire)
	return nil
}

func (redisClient *RedisClientWrapper) GetValueFromKey(key string) (*string, error) {
	if redisClient == nil {
		return nil, errors.New("get value from redis error because redis client nil")
	}

	data, err := redisClient.Client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		logger.Warnf("Get Data From redis Cache error:%s", err.Error())
		return nil, err
	}

	logger.Debugf("GetValueFromKey with key:%s value:%s", key, data)
	return &data, nil
}

func (redisClient *RedisClientWrapper) DeleteDataFromKeys(keys []string) error {
	if redisClient == nil {
		return errors.New("get value from redis error because redis client nil")
	}
	pipe := redisClient.Client.Pipeline()
	for _, key := range keys {
		if len(key) <= 0 {
			continue
		}
		pipe.Del(context.Background(), key)
	}
	result, _err := pipe.Exec(context.Background())
	if _err != nil {
		return _err
	}
	logger.Debugf("DeleteDataFromKey with result:%s", result)
	return nil
}

func (redisClient *RedisClientWrapper) FlushAllAsync() error {
	if redisClient == nil {
		return errors.New("get value from redis error because redis client nil")
	}
	data := redisClient.Client.FlushAllAsync(context.Background())
	logger.Debugf("FlushAllAsync with result:%s", data)
	return nil
}

func (redisClient *RedisClientWrapper) SetNXDataToCache(key string, value string, exprire time.Duration) (bool, error) {
	if redisClient == nil {
		return false, errors.New("get value from redis error because redis client nil")
	}
	logger.Debugf("set data to redis with key:%s value:%s", key, value)

	result, err := redisClient.Client.SetNX(context.Background(), key, value, exprire).Result()
	if err != nil {
		logger.Warnf("Set Data To redis Cache error:%s", err.Error())
		return false, err
	}
	logger.Debugf("SetDataToCache sucessfully with key:%s result:%s expireTime:%s", key, result, exprire)
	return result, nil
}

func (redisClient *RedisClientWrapper) DeleteDataFromKey(key string) error {
	if redisClient == nil {
		return errors.New("get value from redis error because redis client nil")
	}
	//logger.Debugf("set data to redis with key:%s value:%s", key, value)

	_, err := redisClient.Client.Del(context.Background(), key).Result()
	if err != nil {
		logger.Warnf("Delete Data From redis Cache error:%s", err.Error())
		return err
	}
	//logger.Infof("SetDataToCache sucessfully with key:%s result:%s expireTime:%s", key, result, exprire)
	return nil
}

func (redisClient *RedisClientWrapper) ZScore(key string, member string) (*float64, error) {
	if redisClient == nil {
		return nil, errors.New("get value from redis error because redis client nil")
	}
	//logger.Debugf("set data to redis with key:%s value:%s", key, value)

	data, err := redisClient.Client.ZScore(context.Background(), key, member).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		logger.Warnf("ZScore Data From redis Cache error:%s", err.Error())
		return nil, err
	}

	//logger.Infof("SetDataToCache sucessfully with key:%s result:%s expireTime:%s", key, result, exprire)
	return &data, nil
}

func (redisClient *RedisClientWrapper) ZAdd(key string, member string, score float64) error {
	if redisClient == nil {
		return errors.New("get value from redis error because redis client nil")
	}
	//logger.Debugf("set data to redis with key:%s value:%s", key, value)

	_, err := redisClient.Client.ZAdd(context.Background(), key, redis.Z{Score: score, Member: member}).Result()
	if err != nil {
		logger.Warnf("ZAdd Data To redis Cache error:%s", err.Error())
		return err
	}
	//logger.Infof("SetDataToCache sucessfully with key:%s result:%s expireTime:%s", key, result, exprire)
	return nil
}

func (redisClient *RedisClientWrapper) ZRem(key string, members ...interface{}) error {
	if redisClient == nil {
		return errors.New("get value from redis error because redis client nil")
	}
	//logger.Debugf("set data to redis with key:%s value:%s", key, value)

	_, err := redisClient.Client.ZRem(context.Background(), key, members).Result()
	if err != nil {
		logger.Warnf("ZRem Data To redis Cache error:%s", err.Error())
		return err
	}
	//logger.Infof("SetDataToCache sucessfully with key:%s result:%s expireTime:%s", key, result, exprire)
	return nil
}

func (redisClient *RedisClientWrapper) ZRangeByScore(key string, min float64, max float64) ([]string, error) {
	if redisClient == nil {
		return nil, errors.New("get value from redis error because redis client nil")
	}
	//logger.Debugf("set data to redis with key:%s value:%s", key, value)

	components, err := redisClient.Client.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{Min: strconv.FormatFloat(min, 'f', 0, 64), Max: strconv.FormatFloat(max, 'f', 0, 64)}).Result()
	if err != nil {
		logger.Warnf("ZRangeByScore Data From redis Cache error:%s", err.Error())
		return nil, err
	}
	//logger.Infof("SetDataToCache sucessfully with key:%s result:%s expireTime:%s", key, result, exprire)
	return components, nil
}

func (redisClient *RedisClientWrapper) SubscribeChannel(channelName string) (<-chan *redis.Message, error) {
	if redisClient == nil {
		return nil, errors.New("get value from redis error because redis client nil")
	}

	pubsub := redisClient.Client.Subscribe(context.Background(), channelName)
	//defer pubsub.Close()

	ch := pubsub.Channel()
	fmt.Println("Subscribed to %s", channelName)

	return ch, nil
}

func (redisClient *RedisClientWrapper) Publish(channelName string, message interface{}) error {
	if redisClient == nil {
		return errors.New("get value from redis error because redis client nil")
	}
	// Publish a message
	err := redisClient.Client.Publish(context.Background(), channelName, message).Err()
	if err != nil {
		logger.Infof("Publish Data To redis Cache error:%s", err.Error())
		return err
	}

	return nil
}

func (redisClient *RedisClientWrapper) SubscribeChannel2(channelName string) *redis.PubSub {
	if redisClient == nil {
		return nil
	}

	return redisClient.Client.Subscribe(context.Background(), channelName)
}
