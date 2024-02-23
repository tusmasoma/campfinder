package config

import (
	"context"

	"github.com/sethvargo/go-envconfig"
)

const (
	dbPrefix    = "MYSQL_"
	cachePrefix = "REDIS_"
)

type DBConfig struct {
	Host     string `env:"HOST, required"`
	Port     string `env:"PORT, required"`
	User     string `env:"USER, required"`
	Password string `env:"PASSWORD, required"`
	DBName   string `env:"DB_NAME, required"`
}

type CacheConfig struct {
	Addr     string `env:"ADDR, required"`
	Password string `env:"PASSWORD, required"`
	Db       int    `env:"DB, required"`
}

func NewDBConfig(ctx context.Context) (*DBConfig, error) {
	conf := &DBConfig{}
	pl := envconfig.PrefixLookuper(dbPrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		return nil, err
	}
	return conf, nil
}

func NewCacheConfig(ctx context.Context) (*CacheConfig, error) {
	conf := &CacheConfig{}
	pl := envconfig.PrefixLookuper(cachePrefix, envconfig.OsLookuper())
	if err := envconfig.ProcessWith(ctx, conf, pl); err != nil {
		return nil, err
	}
	return conf, nil
}
