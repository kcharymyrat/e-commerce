package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const version = "v1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn               string
		maxConns          int32
		minConns          int32
		maxConnLifetime   time.Duration
		maxConnIdleTime   time.Duration
		healthCheckPeriod time.Duration
		connectTimeout    time.Duration
	}
}

type application struct {
	config config
	logger *zerolog.Logger
	models data.Models
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	var cfg config

	loadEnv()

	port := viper.GetInt("APP_PORT")
	env := viper.GetString("ENV")
	dbDsn := viper.GetString("POSTGRES_DSN")

	poolMaxConns := viper.GetInt32("POOL_MAX_CONNS")
	poolMinConns := viper.GetInt32("POOL_MIN_CONNS")
	poolMaxConnLifetimeHours := viper.GetInt("POOL_MAX_CONN_LITETIME_HOURS")
	poolMaxConnIdleTimeMinutes := viper.GetInt("POOL_MAX_CONN_IDLE_TIME_MINUTES")
	poolHealthCheckPeriodMinutes := viper.GetInt("POOL_HEALTH_CHECK_PERIOD_MINUTES")
	poolConnectTimeoutSeconds := viper.GetInt("POOL_CONNECT_TIMEOUT_SECONDS")

	flag.IntVar(&cfg.port, "port", port, "API server port")
	flag.StringVar(&cfg.env, "env", env, "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", dbDsn, "PostgreSQL DSN")

	cfg.db.maxConns = poolMaxConns
	cfg.db.minConns = poolMinConns
	cfg.db.maxConnLifetime = time.Duration(poolMaxConnLifetimeHours) * time.Hour
	cfg.db.maxConnIdleTime = time.Duration(poolMaxConnIdleTimeMinutes) * time.Minute
	cfg.db.healthCheckPeriod = time.Duration(poolHealthCheckPeriodMinutes) * time.Minute
	cfg.db.connectTimeout = time.Duration(poolConnectTimeoutSeconds) * time.Second

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.Error().Stack().Err(err).Msg("some message")
	}
	defer db.Close()

	log.Info().Str("env", cfg.env).Msg("database connection pool established")

	app := &application{
		config: cfg,
		logger: &logger,
		models: data.NewModels(db),
	}

	mux := app.routes()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info().Msg(fmt.Sprintf("starting %s server on %s", cfg.env, srv.Addr))

	err = srv.ListenAndServe()
	logger.Fatal().Stack().Err(err).Msg("other")

}

func loadEnv() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(err)
		} else {
			panic(err)
		}
	}
}

func openDB(cfg config) (*pgxpool.Pool, error) {
	// Parse the connection pool configuration
	poolConfig, err := pgxpool.ParseConfig(cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = cfg.db.maxConns
	poolConfig.MinConns = cfg.db.minConns
	poolConfig.MaxConnLifetime = cfg.db.maxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.db.maxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.db.healthCheckPeriod
	poolConfig.ConnConfig.ConnectTimeout = cfg.db.connectTimeout

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
