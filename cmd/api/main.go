package main

import (
	"context"
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/config"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/kcharymyrat/e-commerce/internal/routes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type loggerAdapter struct {
	logger zerolog.Logger
}

func (l loggerAdapter) Write(p []byte) (n int, err error) {
	l.logger.Error().Msg(string(p))
	return len(p), nil
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	cfg := config.Config{}

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

	app := &app.Application{
		Config: cfg,
		Logger: &logger,
		Models: data.NewModels(db),
	}

	mux := routes.Routes(app)

	stdLog := stdlog.New(loggerAdapter{logger}, "", 0)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ErrorLog:     stdLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Info().Msg(fmt.Sprintf("starting %s server on %s", cfg.Env, srv.Addr))

	err = srv.ListenAndServe()
	logger.Fatal().Stack().Err(err).Msg("Server failed to start")

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
