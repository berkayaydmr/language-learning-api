package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/berkayaydmr/language-learning-api/pkg/storage"
	"github.com/berkayaydmr/language-learning-api/pkg/transport"
	"github.com/berkayaydmr/language-learning-api/pkg/transport/middleware/authmiddleware"
)

func main() {
	// storage olusturulacak
	// http handler'lar olusturulacak
	// http server olusturulacak
	// http server baslatilacak

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	storage := storage.New()
	err := storage.Open(ctx, "../../words.db")
	if err != nil {
		panic(err)
	}

	logger := slog.Default()

	var handler http.Handler

	authMiddleware := authmiddleware.NewAuthMiddleware("a", handler)

	handler = transport.MakeHTTPHandler(logger, storage, authMiddleware)

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
