package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/berkayaydmr/language-learning-api/pkg/storage"
	"github.com/berkayaydmr/language-learning-api/pkg/transport"
	"github.com/berkayaydmr/language-learning-api/pkg/transport/middleware/authmiddleware"
)

const (
	dsn    = "cmd/server/words.db"
	apiKey = "a"

	addr            = ":8080"
	readTimeout     = 5 * time.Second
	writeTimeout    = 10 * time.Second
	idleTimeout     = 60 * time.Second
	maxHeaderBytes  = 1 * 1024 * 1024 // 1MB
	shutdownTimeout = 15 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	s := storage.New()
	err := s.Open(ctx, dsn)
	if err != nil {
		slog.Error("failed to open db", "error", err)
		return
	}

	authMiddleware := authmiddleware.NewAuthMiddleware(apiKey)

	handler := transport.MakeHTTPHandler(logger, s, authMiddleware)

	server := &http.Server{
		Addr:           addr,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		MaxHeaderBytes: maxHeaderBytes,
		Handler:        handler,
	}

	go func() {
		logger.Info("server is listening", "addr", addr)
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "error", err)
			cancel()
		}
	}()

	<-ctx.Done()

	ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancelShutdown()

	err = server.Shutdown(ctxShutdown)
	if err != nil {
		logger.Error("failed to close server", "error", err)
	}

	err = s.Close()
	if err != nil {
		logger.Error("failed to close storage", "error", err)
	}
}
