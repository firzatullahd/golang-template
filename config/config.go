package config

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	AppName       string `mapstructure:"APP_NAME"`
	AppEnv        string `mapstructure:"APP_ENV"`
	AppPort       string `mapstructure:"APP_PORT"`
	JWTSecretKey  string `mapstructure:"JWT_SECRET_KEY"`
	DBName        string `mapstructure:"DB_NAME"`
	DBHost        string `mapstructure:"DB_HOST"`
	DBPort        string `mapstructure:"DB_PORT"`
	DBUsername    string `mapstructure:"DB_USERNAME"`
	DBPassword    string `mapstructure:"DB_PASSWORD"`
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB       int    `mapstructure:"REDIS_DB"`
}

func Load() *Config {
	var conf Config
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't find the file .env : ", err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		log.Fatal("Environment can't be loaded: ", err)
	}

	if conf.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &conf
}
