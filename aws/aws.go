package aws

import (
	awsConfig "go-common/aws/config"
	awsRedis "go-common/aws/redis"
)

func InitAws() {
	awsConfig.InitAws()
}

func InitRedis(cacheClusterName, cacheUrl string) (*awsRedis.RedisClientWrapper, error) {
	return awsRedis.InitRedis(cacheClusterName, cacheUrl)
}
