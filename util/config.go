package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenSymmetricKey   string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (config Config, err error) {
	//配置文件路径
	viper.AddConfigPath(path)
	//配置文件名
	viper.SetConfigName("app")
	//配置文件类型
	viper.SetConfigType("env")

	//自动匹配配置项
	viper.AutomaticEnv()

	//读取匹配项
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	//反射到结构体
	err = viper.Unmarshal(&config)

	return
}
