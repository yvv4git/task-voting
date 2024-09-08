package infrastructure

import (
	"log/slog"
	"os"
)

func NewDefaultLogger() *slog.Logger {
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, nil),
	)

	return logger
}
