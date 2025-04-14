package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/70ekanetugu/webdemo/handler"
	"github.com/70ekanetugu/webdemo/middleware"
	"github.com/70ekanetugu/webdemo/repository"
)

func main() {
	if err := repository.Initialize(os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD")); err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr: "0.0.0.0:8080",
	}

	shutdownBus := make(chan int, 1)
	defer close(shutdownBus)

	listenSignal(server, shutdownBus)

	http.HandleFunc("GET /hello", middleware.Logging(handler.HelloWorld))
	http.HandleFunc("GET /error", middleware.Logging(handler.ErrorHandler))
	http.HandleFunc("GET /todos", middleware.Logging(handler.GetTodoList))
	http.HandleFunc("GET /todos/{id}", middleware.Logging(handler.GetTodoById))
	http.HandleFunc("POST /todos", middleware.Logging(handler.SaveTodo))
	http.HandleFunc("PUT /todos/{id}", middleware.Logging(handler.SaveTodo))
	http.HandleFunc("PATCH /todos/{id}/status", middleware.Logging(handler.SaveStatus))
	http.HandleFunc("GET /images/{name}", middleware.Logging(handler.DownloadImage))

	slog.Info("Starting server...")
	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("HTTP server ListenAndServe", slog.String("err", err.Error()))
	}

	<-shutdownBus
}

// SIGINT, SIGTERMを監視するゴルーチンを起動する
func listenSignal(server *http.Server, shutdownBus chan<- int) {
	go func() {
		// シグナル監視用のチャネル
		sigs := make(chan os.Signal, 1)
		defer close(sigs)

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigs
		slog.Info("Got signal", slog.String("Signal", s.String()))

		// シグナルを受け取ったらシャットダウンを実行
		if err := server.Shutdown(context.Background()); err != nil {
			slog.Error("HTTP server Shutdown", slog.String("err", err.Error()))
			shutdownBus <- 1
			return
		}
		slog.Info("Server was shutdown")
		shutdownBus <- 0
	}()
}
