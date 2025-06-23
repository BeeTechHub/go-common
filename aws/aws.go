package aws

import (
	awsConfig "github.com/BeeTechHub/go-common/aws/config"
	awsRedis "github.com/BeeTechHub/go-common/aws/redis"
)

func InitAws() {
	awsConfig.InitAws()
}

func InitRedis(cacheClusterName, cacheUrl string) (*awsRedis.RedisClientWrapper, error) {
	return awsRedis.InitRedis(cacheClusterName, cacheUrl)
}
