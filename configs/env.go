package configs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BeeTechHub/go-common/logger"

	"github.com/joho/godotenv"
)

func GetEnv(envName string) string {
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env name:" + envName)
	}

	return os.Getenv(envName)
}

func GetEnvFromOS(envName string) string {
	result := os.Getenv(envName)
	fmt.Printf("GetEnvFromOS with name:%s result:%s", envName, result)
	if len(result) <= 0 {
		return GetEnv(envName)
	}
	return result
}

func GetInt64EnvFromOS(envName string) (int64, error) {
	val := GetEnvFromOS(envName)
	ret, err := strconv.ParseInt(val, 10, 64)
	return ret, err
}

func GetBoolEnvFromOS(envName string) (bool, error) {
	val := GetEnvFromOS(envName)
	ret, err := strconv.ParseBool(val)
	return ret, err
}

func GetCacheHost() string {
	return GetEnvFromOS("CACHE_HOST")
}
