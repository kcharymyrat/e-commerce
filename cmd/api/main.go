package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/config"
	"github.com/kcharymyrat/e-commerce/internal/repository"
	"github.com/kcharymyrat/e-commerce/internal/server"
	"github.com/kcharymyrat/e-commerce/internal/validation"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

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

	app := app.NewApplication(
		cfg,
		&logger,
		repository.NewRepositories(db),
		rdb,
		limiter,
		validator,
		valUniTrans,
	)

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
