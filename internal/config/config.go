package config

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Config struct {
	Env  string `mapstructure:"env"`
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
	DB   DB     `mapstructure:"database"`
}

type DB struct {
	MaxIdleCons    int  `mapstructure:"maxIdleCons"`
	MaxOpenCons    int  `mapstructure:"maxOpenCons"`
	ConMaxIdleTime int  `mapstructure:"conMaxIdleTime"`
	ConMaxLifetime int  `mapstructure:"conMaxLifeTime"`
	replica        PSQL `mapstructure:"replica"`
	Master         PSQL `mapstructure:"master"`
}

type PSQL struct {
	DBName   string `mapstructure:"dbName"`
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	Schema   string `mapstructure:"schema"`
	User     string `mapstructure:"user"`
	Debug    bool   `mapstructure:"debug"`
}

func (p PSQL) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s search_path=%s sslmode=%s",
		p.Host, strconv.Itoa(p.Port), p.User, p.Password, p.DBName, p.Schema, "disable")
}

func Load() *Config {
	v := viper.New()

	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetConfigName("config")

	if err := v.ReadInConfig(); err != nil {
		log.Panic("error read config")
	}

	conf := Config{}
	if err := v.Unmarshal(&conf); err != nil {
		log.Panic("error read config")
	}

	return &conf
}

func InitializeDB(conf *DB) (master *sqlx.DB, replica *sqlx.DB) {
	var err error
	ctx := context.Background()
	master, err = sqlx.Open("postgres", conf.Master.ConnectionString())
	if err != nil {
		log.Fatal(ctx, "Can't connect to master DB %+v", err)
	}

	master.SetMaxIdleConns(conf.MaxIdleCons)
	master.SetMaxOpenConns(conf.MaxOpenCons)
	master.SetConnMaxLifetime(time.Duration(conf.ConMaxLifetime) * time.Millisecond)
	master.SetConnMaxIdleTime(time.Duration(conf.ConMaxIdleTime) * time.Millisecond)

	replica, err = sqlx.Connect("postgres", conf.replica.ConnectionString())
	if err != nil {
		log.Fatal(ctx, "Can't connect to replica DB %+v", err)
	}

	replica.SetMaxIdleConns(conf.MaxIdleCons)
	replica.SetMaxOpenConns(conf.MaxOpenCons)
	replica.SetConnMaxLifetime(time.Duration(conf.ConMaxLifetime) * time.Millisecond)
	replica.SetConnMaxIdleTime(time.Duration(conf.ConMaxIdleTime) * time.Millisecond)

	return
}
