package app

import (
	"github.com/kcharymyrat/e-commerce/internal/config"
	"github.com/kcharymyrat/e-commerce/internal/data"
	"github.com/rs/zerolog"
)

type Application struct {
	Config config.Config
	Logger *zerolog.Logger
	Models data.Models
}
