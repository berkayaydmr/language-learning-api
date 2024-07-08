package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/berkayaydmr/language-learning-api/common"
	"github.com/berkayaydmr/language-learning-api/pkg/storage"
	"github.com/berkayaydmr/language-learning-api/transport"
	"github.com/berkayaydmr/language-learning-api/transport/middleware/authmiddleware"
)

func main() {
	// storage olusturulacak
	// http handler'lar olusturulacak
	// http server olusturulacak
	// http server baslatilacak

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := storage.New()
	err := storage.Open(ctx, common.DSN)
	if err != nil {
		panic(err)
	}

	logger := slog.Default()

	authMiddleware := authmiddleware.NewAuthMiddleware("a", nil)

	handler := transport.MakeHTTPHandler(logger, storage, authMiddleware)

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	logger.Info("Server Started Running on :8080")
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
