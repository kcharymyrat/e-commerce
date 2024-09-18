package config

import (
	"time"
)

type Config struct {
	Port int
	Env  string
	DB   struct {
		DSN               string
		MaxConns          int32
		MinConns          int32
		MaxConnLifetime   time.Duration
		MaxConnIdleTime   time.Duration
		HealthCheckPeriod time.Duration
		ConnectTimeout    time.Duration
	}
}
