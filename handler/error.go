package handler

import (
	"log/slog"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request) {
	slog.Error("ErrorHandler", slog.String("err", "500 Internal Server Error"))
	// 500エラーを発生させる
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("500 Internal Server Error"))
	return
}
