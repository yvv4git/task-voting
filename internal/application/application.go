package application

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

type Application interface {
	start(ctx context.Context) error
}

type application struct {
	log *slog.Logger
	app Application
}

func (a *application) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		<-quit
		a.log.Info("received shutdown signal")
		cancel()
	}()

	if err := a.app.start(ctx); err != nil {
		return err
	}

	return nil
}
