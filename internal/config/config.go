package config

import (
	"fmt"
	"log"
	"time"

	"github.com/firzatullahd/cats-social-api/internal/utils/constant"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	*viper.Viper
	DBName       string
	Host         string
	Password     string
	Port         string
	Schema       string
	User         string
	Params       string
	BcryptSalt   string
	JWTSecretKey string
}

func (p *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.DBName, p.Schema, "disable")
}

func Load() *Config {
	v := viper.New()
	v.AutomaticEnv()

	return &Config{
		Viper:        v,
		DBName:       v.GetString("DB_NAME"),
		Host:         v.GetString("DB_HOST"),
		Password:     v.GetString("DB_PASSWORD"),
		Port:         v.GetString("DB_PORT"),
		Schema:       "public",
		User:         v.GetString("DB_USERNAME"),
		Params:       v.GetString("DB_PARAMS"),
		JWTSecretKey: v.GetString("JWT_SECRET"),
		BcryptSalt:   v.GetString("BCRYPT_SALT"),
	}
}

func (c *Config) InitializeDB() (master *sqlx.DB, replica *sqlx.DB) {
	var err error
	master, err = sqlx.Open("postgres", c.ConnectionString())
	if err != nil {
		log.Fatal("Can't connect to master DB", err)
	}

	log.Println("Successfully connect to master DB")

	master.SetMaxIdleConns(constant.MaxIdleCons)
	master.SetMaxOpenConns(constant.MaxOpenCons)
	master.SetConnMaxLifetime(time.Duration(constant.ConMaxLifeTime) * time.Millisecond)
	master.SetConnMaxIdleTime(time.Duration(constant.ConMaxIdleTime) * time.Millisecond)

	replica, err = sqlx.Connect("postgres", c.ConnectionString())
	if err != nil {
		log.Fatal("Can't connect to replica DB", err)
	}

	replica.SetMaxIdleConns(constant.MaxIdleCons)
	replica.SetMaxOpenConns(constant.MaxOpenCons)
	replica.SetConnMaxLifetime(time.Duration(constant.ConMaxLifeTime) * time.Millisecond)
	replica.SetConnMaxIdleTime(time.Duration(constant.ConMaxIdleTime) * time.Millisecond)

	log.Println("Successfully connect to replica DB")

	return
}
