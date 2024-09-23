package app

import (
	"github.com/go-redis/redis_rate/v10"
	"github.com/kcharymyrat/e-commerce/internal/config"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Application struct {
	Config  config.Config
	Logger  *zerolog.Logger
	Models  data.Models
	RDB     *redis.Client
	Limiter *redis_rate.Limiter
}

func NewApplication(
	cfg config.Config,
	logger *zerolog.Logger,
	models data.Models,
	rdb *redis.Client,
	limiter *redis_rate.Limiter,
) *Application {
	return &Application{
		Config:  cfg,
		Logger:  logger,
		Models:  models,
		RDB:     rdb,
		Limiter: limiter,
	}
}
