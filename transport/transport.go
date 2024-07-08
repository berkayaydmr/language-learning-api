package transport

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/berkayaydmr/language-learning-api/pkg/storage"
	"github.com/berkayaydmr/language-learning-api/pkg/utils"
	"github.com/berkayaydmr/language-learning-api/transport/middleware"
)

func MakeHTTPHandler(logger *slog.Logger, storages storage.Storage, authMiddleware middleware.Middleware) http.Handler {
	handler := http.NewServeMux()

	handler.Handle("GET /health-check", makeHealthCheckHandler(logger))
	handler.Handle("GET /words", makeListHandler(logger, storages))
	handler.Handle("POST /words", authMiddleware(makeCreateHandler(logger, storages)))
	handler.Handle("PATCH /words/{id}", authMiddleware(makeUpdateHandler(logger, storages)))
	handler.Handle("DELETE /words/{id}", authMiddleware(makeDeleteHandler(logger, storages)))

	return handler
}

func makeHealthCheckHandler(logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func makeListHandler(logger *slog.Logger, storages storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("List Handler Called")
		context, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		words, err := storages.List(context)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		utils.RespondWithJSON(w, words, http.StatusOK)
	})
}

type Create struct {
	Word            string `json:"word"`
	Translation     string `json:"translation"`
	Language        string `json:"language"`
	ExampleSentence string `json:"exampleSentence"`
}

func makeCreateHandler(logger *slog.Logger, storages storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := Create{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		word := storage.Word{
			Word:            request.Word,
			Translation:     request.Translation,
			Language:        request.Language,
			ExampleSentence: request.ExampleSentence,
		}

		context, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		id, err := storages.Create(context, word)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		utils.RespondWithJSON(w, map[string]any{"id": &id}, http.StatusCreated)
	})
}

func makeUpdateHandler(logger *slog.Logger, storages storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := storage.Update{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		idStr, err := utils.GetUrlParam(r, utils.UrlParamKeyID)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		id, err := utils.ParseStrToInt(*idStr)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		context, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		err = storages.Update(context, storage.PrimaryKey(id), request)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		utils.RespondWithJSON(w, nil, http.StatusNoContent)
	})
}

func makeDeleteHandler(logger *slog.Logger, storages storage.Storage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idStr, err := utils.GetUrlParam(r, utils.UrlParamKeyID)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		id, err := utils.ParseStrToInt(*idStr)
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		context, cancel := context.WithTimeout(context.Background(), time.Minute*2)
		defer cancel()

		err = storages.Delete(context, storage.PrimaryKey(id))
		if err != nil {
			logger.Error(err.Error())
			utils.RespondWithError(w, err)
			return
		}

		utils.RespondWithJSON(w, nil, http.StatusNoContent)
	})
}
