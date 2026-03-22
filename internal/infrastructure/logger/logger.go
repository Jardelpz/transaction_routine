package logger

import (
	"log/slog"
	"os"
)

func NewSlog() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
