package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/config"
	"github.com/kcharymyrat/e-commerce/internal/repository"
	"github.com/kcharymyrat/e-commerce/internal/server"
	"github.com/kcharymyrat/e-commerce/internal/validation"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"golang.org/x/text/language"

	_ "github.com/kcharymyrat/e-commerce/docs"
)

// @title E-commerce
// @version 1.0
// @description Go lang E-commerce

// @contact.name Charymyrat Garryyev
// @contact.email kcharymyrat@gmail.com

// @license.name Apache 2.0
// @license.url ""

// @host localhost:4000
// @BasePath /

// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	cfg := config.Config{}

	loadEnv(&logger)

	port := viper.GetInt("APP_PORT")
	env := viper.GetString("ENV")
	dbDsn := viper.GetString("POSTGRES_DSN")

	poolMaxConns := viper.GetInt32("POOL_MAX_CONNS")
	poolMinConns := viper.GetInt32("POOL_MIN_CONNS")
	poolMaxConnLifetimeHours := viper.GetInt("POOL_MAX_CONN_LITETIME_HOURS")
	poolMaxConnIdleTimeMinutes := viper.GetInt("POOL_MAX_CONN_IDLE_TIME_MINUTES")
	poolHealthCheckPeriodMinutes := viper.GetInt("POOL_HEALTH_CHECK_PERIOD_MINUTES")
	poolConnectTimeoutSeconds := viper.GetInt("POOL_CONNECT_TIMEOUT_SECONDS")

	redisAddr := viper.GetString("REDIS_ADDR")
	redisPort := viper.GetInt("REDIS_PORT")

	flag.IntVar(&cfg.Port, "port", port, "API server port")
	flag.StringVar(&cfg.Env, "env", env, "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.DSN, "db-dsn", dbDsn, "PostgreSQL DSN")

	cfg.DB.MaxConns = poolMaxConns
	cfg.DB.MinConns = poolMinConns
	cfg.DB.MaxConnLifetime = time.Duration(poolMaxConnLifetimeHours) * time.Hour
	cfg.DB.MaxConnIdleTime = time.Duration(poolMaxConnIdleTimeMinutes) * time.Minute
	cfg.DB.HealthCheckPeriod = time.Duration(poolHealthCheckPeriodMinutes) * time.Minute
	cfg.DB.ConnectTimeout = time.Duration(poolConnectTimeoutSeconds) * time.Second

	secretKey := viper.GetString("SECRET_KEY")
	if len(secretKey) < 32 {
		log.Fatal().Msg("secret key must be at least 32 characters long")
	}
	cfg.SecretKey = []byte(secretKey)

	flag.Parse()

	db, err := openDB(&cfg)
	if err != nil {
		logger.Error().Stack().Err(err).Msg("DB connection failed")
	}
	defer db.Close()
	log.Info().Str("env", cfg.Env).Msg("database connection pool established")

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", redisAddr, redisPort),
	})
	limiter := redis_rate.NewLimiter(rdb)
	log.Info().Str("env", cfg.Env).Msg("redis connection and limiter established")

	validator := validation.NewValidator()
	valUniTrans := validation.NewUniversalTranslator()
	i18nBundle := loadTranslations()
	wg := sync.WaitGroup{}

	app := app.NewApplication(
		cfg,
		&logger,
		repository.NewRepositories(db),
		rdb,
		limiter,
		validator,
		valUniTrans,
		i18nBundle,
		&wg,
	)

	// Get the translator for each language (for example, English)
	enTrans := validation.GetTranslator(valUniTrans, "en")
	ruTrans := validation.GetTranslator(valUniTrans, "ru_RU")
	tkTrans := validation.GetTranslator(valUniTrans, "tk_TM")

	// Register default translations for each language
	err = validation.RegisterTranslations(app, enTrans, "en")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to register English translations")
	}

	err = validation.RegisterTranslations(app, ruTrans, "ru_RU")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to register Russian translations")
	}

	err = validation.RegisterTranslations(app, tkTrans, "tk_TM")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to register Turkmen translations")
	}

	// Register custom translations (you can reuse the same translators)
	validation.RegisterCustomEnTranslations(app, enTrans)
	validation.RegisterCustomRuTranslations(app, ruTrans)
	validation.RegisterCustomTkTranslations(app, tkTrans)

	err = server.Serve(app)
	if err != nil {
		app.Logger.Fatal().Stack().Err(err).Msg("Server failed to start")
	}
}

func loadEnv(logger *zerolog.Logger) {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Fatal().Stack().Err(err).Msg("Config file .env not found")
			panic(err)
		} else {
			logger.Fatal().Stack().Err(err).Msg("Reading .env file failed")
			panic(err)
		}
	}
}

func openDB(cfg *config.Config) (*pgxpool.Pool, error) {
	// Parse the connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.DB.MaxConns
	poolConfig.MinConns = cfg.DB.MinConns
	poolConfig.MaxConnLifetime = cfg.DB.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.DB.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.DB.HealthCheckPeriod
	poolConfig.ConnConfig.ConnectTimeout = cfg.DB.ConnectTimeout

	// Create the connection pool
	dbConnPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbConnPool.Ping(ctx)
	if err != nil {
		dbConnPool.Close()
		return nil, err
	}

	return dbConnPool, nil

}

func loadTranslations() *i18n.Bundle {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load translation files
	bundle.MustLoadMessageFile("internal/translations/en.json")
	bundle.MustLoadMessageFile("internal/translations/ru.json")
	bundle.MustLoadMessageFile("internal/translations/tk.json")

	return bundle
}
