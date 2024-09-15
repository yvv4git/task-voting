package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type GinServerWrapper struct {
	server *http.Server
	log    *slog.Logger
}

func NewWebServer(log *slog.Logger, router *gin.Engine, addr string) *GinServerWrapper {
	return &GinServerWrapper{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		log: log,
	}
}

func (w *GinServerWrapper) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second) // TODO: fix hardcode
		defer cancel()
		if err := w.server.Shutdown(ctxShutdown); err != nil {
			w.log.Error("server shutdown", slog.Any("error", err))
		}
	}()

	w.log.Info(fmt.Sprintf("web server starting on addr %s", w.server.Addr))
	if err := w.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
