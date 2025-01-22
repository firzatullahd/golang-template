package config

import (
	"context"
	"fmt"

	"github.com/firzatullahd/golang-template/utils/logger"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

func (c *Config) PsqlConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUsername, c.DBPassword, c.DBName, "disable")
}

func NewPsqlClient(conf *Config) (*sqlx.DB, error) {
	conn, err := sqlx.Open("postgres", conf.PsqlConnectionString())
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("Connection to database failed %w", err)
	}

	logger.Log.Info("Connection to database established")

	return conn, nil
}

// TODO: separate redis connection for different purpose
func NewRedisClient(ctx context.Context, conf *Config) (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr:                  conf.RedisAddress,
		Password:              conf.RedisPassword,
		DB:                    conf.RedisDB,
		Protocol:              2,
		ContextTimeoutEnabled: true,
	})

	if err := conn.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("Connection to redis failed %w", err)
	}

	logger.Log.Info("Connection to redis established")
	return conn, nil

}
