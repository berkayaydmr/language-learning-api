package utils

import (
	"encoding/json"
	"net/http"

	customerr "github.com/berkayaydmr/language-learning-api/pkg/error"
)

func RespondWithError(w http.ResponseWriter, err error) {
	resp, code := customerr.NewErrorResponse(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func Respond(w http.ResponseWriter, data interface{}, code int) {
	resp, err := json.Marshal(data)
	if err != nil {
		RespondWithError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
