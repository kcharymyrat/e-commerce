package app

import (
	"sync"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis_rate/v10"
	"github.com/kcharymyrat/e-commerce/internal/config"
	"github.com/kcharymyrat/e-commerce/internal/repository"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Application struct {
	Config       config.Config
	Logger       *zerolog.Logger
	Repositories repository.Repositories
	RDB          *redis.Client
	Limiter      *redis_rate.Limiter
	Validator    *validator.Validate
	ValUniTrans  *ut.UniversalTranslator
	I18nBundle   *i18n.Bundle
	Wg           *sync.WaitGroup
}

func NewApplication(
	cfg config.Config,
	logger *zerolog.Logger,
	repositories repository.Repositories,
	rdb *redis.Client,
	limiter *redis_rate.Limiter,
	validator *validator.Validate,
	uniTrans *ut.UniversalTranslator,
	i18nBundle *i18n.Bundle,
	wg *sync.WaitGroup,
) *Application {
	return &Application{
		Config:       cfg,
		Logger:       logger,
		Repositories: repositories,
		RDB:          rdb,
		Limiter:      limiter,
		Validator:    validator,
		ValUniTrans:  uniTrans,
		I18nBundle:   i18nBundle,
		Wg:           wg,
	}
}
