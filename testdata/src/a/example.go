package a

import (
	"log/slog"
)

func Test() {
	slog.Info("Starting server")    // want "log message should start with a lowercase letter"
	slog.Info("server started")     // ✅
	slog.Info("user password: 123") // want "log message should not contain sensitive data"
	slog.Info("запуск сервера")     // want "log message should be in English"
	slog.Info("server started! 🚀")  // want "log message should not contain special symbols or emojis"
}
