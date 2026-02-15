package a

import (
	"log/slog"
)

func Test() {
	slog.Info("Starting server") // want "log message should start with a lowercase letter"
	slog.Info("server started")  // ✅
}
