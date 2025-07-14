package aws

import (
	awsConfig "github.com/BeeTechHub/go-common/aws/config"
	awsRedis "github.com/BeeTechHub/go-common/aws/redis"
)

func InitAws() {
	awsConfig.InitAws()
}

func InitRedis(cacheClusterName string) (*awsRedis.RedisClientWrapper, error) {
	return awsRedis.InitRedis(cacheClusterName)
}
