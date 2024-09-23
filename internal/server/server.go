package server

import (
	"fmt"
	stdlog "log"
	"net/http"
	"time"

	"github.com/kcharymyrat/e-commerce/internal/app"
	"github.com/kcharymyrat/e-commerce/internal/routes"
	"github.com/rs/zerolog"
)

type loggerAdapter struct {
	logger *zerolog.Logger
}

func (l loggerAdapter) Write(p []byte) (n int, err error) {
	l.logger.Error().Msg(string(p))
	return len(p), nil
}

func Serve(app *app.Application) error {
	mux := routes.Routes(app)

	stdLog := stdlog.New(loggerAdapter{app.Logger}, "", 0)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Config.Port),
		Handler:      mux,
		ErrorLog:     stdLog,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Info().Msg(fmt.Sprintf("starting %s server on %s", app.Config.Env, srv.Addr))

	return srv.ListenAndServe()
}
