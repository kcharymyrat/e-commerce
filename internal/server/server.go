package server

import (
	"context"
	"errors"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		app.Logger.Info().
			Str("signal", s.String()).
			Msg(fmt.Sprintf("shutting down server. Signal - %s", s.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		app.Logger.Info().
			Str("addr", srv.Addr).
			Str("env", app.Config.Env).
			Msg(fmt.Sprintf("completing background tasks of %s server on %s", app.Config.Env, srv.Addr))

		app.Wg.Wait()
		shutdownError <- nil
	}()

	app.Logger.Info().
		Str("addr", srv.Addr).
		Str("env", app.Config.Env).
		Msg(fmt.Sprintf("starting %s server on %s", app.Config.Env, srv.Addr))

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.Logger.Info().
		Str("addr", srv.Addr).
		Msg("server stopped")

	return nil
}
