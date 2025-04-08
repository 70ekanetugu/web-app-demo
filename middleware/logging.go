package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logging(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// ログ出力
		slog.Info("Request-start",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
		)

		start := time.Now()

		// リクエストの処理
		next.ServeHTTP(w, r)

		// レスポンスの処理時間を計測
		duration := time.Since(start)

		// ログ出力
		slog.Info("Request-end",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("remote_addr", r.RemoteAddr),
			slog.Duration("duration", duration),
		)
	})
}
